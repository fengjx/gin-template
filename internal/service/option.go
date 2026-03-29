package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	appLog "gin-template/internal/app/log"
	sysoptionStore "gin-template/internal/store/sysoption"
	pkgcache "gin-template/pkg/cache"
	"gin-template/pkg/errs"
)

var errOptionNotFound = errors.New("配置项不存在")

// option 负责维护系统配置项的内存快照，并按固定周期刷新。
// 这样业务读取常用配置时可以直接走内存，避免每次都访问数据库。
type option struct {
	logger appLog.Logger
	cache  *pkgcache.Cache[string, sysoptionStore.Model]
}

var defaultOption = newOption(time.Minute)

func newOption(refreshInterval time.Duration) *option {
	return newOptionWithLoader(refreshInterval, sysoptionStore.GetAll)
}

func newOptionWithLoader(
	refreshInterval time.Duration,
	loader func(ctx context.Context) ([]sysoptionStore.Model, error),
) *option {
	logger := appLog.Component("service.option")

	return &option{
		logger: logger,
		cache: pkgcache.New(pkgcache.Options[string, sysoptionStore.Model]{
			RefreshInterval: refreshInterval,
			Loader:          newOptionCacheLoader(logger, loader),
		}),
	}
}

// StartOptionAutoRefresh 启动系统配置缓存的后台刷新任务。
// 启动时会先同步加载一次，确保服务对外提供能力前缓存已经可用。
func StartOptionAutoRefresh(ctx context.Context) error {
	return defaultOption.StartAutoRefresh(ctx)
}

// RefreshOptions 会立即从数据库重新加载全部配置并替换内存快照。
// 更新配置后可以主动调用它，避免等待下一个定时周期。
func RefreshOptions(ctx context.Context) error {
	return defaultOption.Refresh(ctx)
}

// UpdateOptionRequest 定义了更新系统配置时需要写入的字段。
type UpdateOptionRequest struct {
	Value       string
	Description string
	IsPublic    bool
}

// UpdateOption 会更新数据库中的配置项，并在成功后立即刷新内存缓存。
// 这样调用方只需要关注业务写入，不需要自行处理缓存一致性。
func UpdateOption(ctx context.Context, key string, req UpdateOptionRequest) (*sysoptionStore.Model, error) {
	return defaultOption.Update(ctx, key, req)
}

// GetOptionString 根据 key 获取字符串配置值。
// 当缓存尚未初始化时，会先同步触发一次加载。
func GetOptionString(ctx context.Context, key string) (string, error) {
	return defaultOption.GetString(ctx, key)
}

// GetOptionJSON 根据 key 读取 JSON 配置，并反序列化为目标结构体。
// 当前只支持 JSON 文本格式，其他序列化格式后续可按需扩展。
func GetOptionJSON[T any](ctx context.Context, key string) (T, error) {
	return getOptionJSON[T](ctx, defaultOption, key)
}

// ResetOptionForTest 用于测试场景下重置全局单例，避免测试间互相污染。
func ResetOptionForTest() {
	defaultOption = newOption(time.Minute)
}

// StartAutoRefresh 启动当前实例的后台刷新循环。
// 如果已经启动过，则直接复用现有循环，不会重复创建 goroutine。
func (s *option) StartAutoRefresh(ctx context.Context) error {
	if err := s.cache.Run(ctx); err != nil {
		return errs.Wrap(err, "启动配置自动刷新前初始化缓存失败")
	}
	return nil
}

// Refresh 从数据库读取全部配置，并以整体替换的方式更新缓存。
// 这里不做增量更新，目的是让内存视图始终与一次完整查询保持一致。
func (s *option) Refresh(ctx context.Context) error {
	return s.cache.Refresh(ctx)
}

// Update 会先更新数据库中的配置记录，再刷新内存缓存。
// 这样可以保证同一进程内后续读取立即拿到最新值。
func (s *option) Update(ctx context.Context, key string, req UpdateOptionRequest) (*sysoptionStore.Model, error) {
	item, err := sysoptionStore.ByKey(ctx, key)
	if err != nil {
		return nil, errs.Wrap(err, "查询待更新配置项失败")
	}

	item.OptionValue = req.Value
	item.Description = req.Description
	item.IsPublic = req.IsPublic
	if err := sysoptionStore.Save(ctx, item); err != nil {
		return nil, errs.Wrap(err, "保存配置项失败")
	}
	if err := s.Refresh(ctx); err != nil {
		return nil, errs.Wrap(err, "刷新配置缓存失败")
	}
	return item, nil
}

// GetString 从缓存中读取字符串值。
// 如果当前快照未命中目标 key，则会同步刷新一次后再重试。
func (s *option) GetString(ctx context.Context, key string) (string, error) {
	item, ok := s.cache.Get(key)
	if !ok {
		if err := s.Refresh(ctx); err != nil {
			return "", errs.Wrap(err, "刷新配置缓存失败")
		}
		item, ok = s.cache.Get(key)
	}
	if !ok {
		return "", errs.WithStack(fmt.Errorf("%w: %s", errOptionNotFound, key))
	}
	return item.OptionValue, nil
}

// getOptionJSON 读取 JSON 文本配置并反序列化到泛型结构体。
// 这里使用顶层泛型函数封装，避免在结构体方法上使用类型参数。
func getOptionJSON[T any](ctx context.Context, s *option, key string) (T, error) {
	var target T

	value, err := s.GetString(ctx, key)
	if err != nil {
		return target, errs.Wrap(err, "读取 JSON 配置项失败")
	}
	if err := json.Unmarshal([]byte(value), &target); err != nil {
		return target, errs.Wrapf(err, "解析配置项 %s 的 JSON 失败", key)
	}
	return target, nil
}

func newOptionCacheLoader(
	logger appLog.Logger,
	loader func(ctx context.Context) ([]sysoptionStore.Model, error),
) pkgcache.Loader[string, sysoptionStore.Model] {
	return func(ctx context.Context) (map[string]sysoptionStore.Model, error) {
		items, err := loader(ctx)
		if err != nil {
			logger.ErrorCtx(ctx, "load option cache failed", zap.Error(err))
			return nil, errs.Wrap(err, "加载系统配置失败")
		}

		nextCache := make(map[string]sysoptionStore.Model, len(items))
		for _, item := range items {
			nextCache[item.OptionKey] = item
		}

		logger.InfoCtx(ctx, "option cache refreshed", zap.Int("count", len(nextCache)))
		return nextCache, nil
	}
}
