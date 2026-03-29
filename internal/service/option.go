package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	appLog "gin-template/internal/app/log"
	sysoptionStore "gin-template/internal/store/sysoption"
	"gin-template/pkg/errs"
	"gin-template/pkg/timex"
	"go.uber.org/zap"
)

var errOptionNotFound = errors.New("配置项不存在")

// optionService 负责维护系统配置项的内存快照，并按固定周期刷新。
// 这样业务读取常用配置时可以直接走内存，避免每次都访问数据库。
type optionService struct {
	logger appLog.Logger

	// refreshInterval 定义后台刷新缓存的周期，默认 1 分钟。
	refreshInterval time.Duration
	// loader 负责从存储层读取全部配置，测试时可替换为自定义实现。
	loader func(ctx context.Context) ([]sysoptionStore.Model, error)
	// now 负责生成当前时间，测试时可替换以便断言缓存刷新时刻。
	now func() time.Time

	// cacheMu 保护 cache 与 loadedAt，确保并发读写安全。
	cacheMu  sync.RWMutex
	cache    map[string]sysoptionStore.Model
	loadedAt time.Time

	// startMu 用于保证后台刷新协程只会启动一次。
	startMu sync.Mutex
	started bool
}

var defaultOptionService = newOptionService(time.Minute)

func newOptionService(refreshInterval time.Duration) *optionService {
	return &optionService{
		logger:          appLog.Component("service.option"),
		refreshInterval: refreshInterval,
		loader:          sysoptionStore.GetAll,
		now:             time.Now,
		cache:           make(map[string]sysoptionStore.Model),
	}
}

// StartOptionAutoRefresh 启动系统配置缓存的后台刷新任务。
// 启动时会先同步加载一次，确保服务对外提供能力前缓存已经可用。
func StartOptionAutoRefresh(ctx context.Context) error {
	return defaultOptionService.StartAutoRefresh(ctx)
}

// RefreshOptions 会立即从数据库重新加载全部配置并替换内存快照。
// 更新配置后可以主动调用它，避免等待下一个定时周期。
func RefreshOptions(ctx context.Context) error {
	return defaultOptionService.Refresh(ctx)
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
	return defaultOptionService.Update(ctx, key, req)
}

// GetOptionString 根据 key 获取字符串配置值。
// 当缓存尚未初始化时，会先同步触发一次加载。
func GetOptionString(ctx context.Context, key string) (string, error) {
	return defaultOptionService.GetString(ctx, key)
}

// GetOptionJSON 根据 key 读取 JSON 配置，并反序列化为目标结构体。
// 当前只支持 JSON 文本格式，其他序列化格式后续可按需扩展。
func GetOptionJSON[T any](ctx context.Context, key string) (T, error) {
	return getOptionJSON[T](ctx, defaultOptionService, key)
}

// ResetOptionServiceForTest 用于测试场景下重置全局单例，避免测试间互相污染。
func ResetOptionServiceForTest() {
	defaultOptionService = newOptionService(time.Minute)
}

// StartAutoRefresh 启动当前实例的后台刷新循环。
// 如果已经启动过，则直接复用现有循环，不会重复创建 goroutine。
func (s *optionService) StartAutoRefresh(ctx context.Context) error {
	s.startMu.Lock()
	defer s.startMu.Unlock()

	if s.started {
		return nil
	}
	if err := s.Refresh(ctx); err != nil {
		return errs.Wrap(err, "启动配置自动刷新前初始化缓存失败")
	}

	s.started = true
	go s.refreshLoop(ctx)
	return nil
}

// Refresh 从数据库读取全部配置，并以整体替换的方式更新缓存。
// 这里不做增量更新，目的是让内存视图始终与一次完整查询保持一致。
func (s *optionService) Refresh(ctx context.Context) error {
	items, err := s.loader(ctx)
	if err != nil {
		return errs.Wrap(err, "加载系统配置失败")
	}

	nextCache := make(map[string]sysoptionStore.Model, len(items))
	for _, item := range items {
		nextCache[item.OptionKey] = item
	}

	s.cacheMu.Lock()
	s.cache = nextCache
	s.loadedAt = s.now()
	s.cacheMu.Unlock()

	s.logger.InfoCtx(ctx, "option cache refreshed", zap.Int("count", len(nextCache)))
	return nil
}

// Update 会先更新数据库中的配置记录，再刷新内存缓存。
// 这样可以保证同一进程内后续读取立即拿到最新值。
func (s *optionService) Update(ctx context.Context, key string, req UpdateOptionRequest) (*sysoptionStore.Model, error) {
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
// 如果缓存尚未加载，则会先尝试进行一次同步刷新。
func (s *optionService) GetString(ctx context.Context, key string) (string, error) {
	if err := s.ensureLoaded(ctx); err != nil {
		return "", errs.Wrap(err, "确保配置缓存已加载失败")
	}

	s.cacheMu.RLock()
	item, ok := s.cache[key]
	s.cacheMu.RUnlock()
	if !ok {
		return "", errs.WithStack(fmt.Errorf("%w: %s", errOptionNotFound, key))
	}
	return item.OptionValue, nil
}

// getOptionJSON 读取 JSON 文本配置并反序列化到泛型结构体。
// 这里使用顶层泛型函数封装，避免在结构体方法上使用类型参数。
func getOptionJSON[T any](ctx context.Context, s *optionService, key string) (T, error) {
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

// ensureLoaded 确保缓存至少被加载过一次。
// 这样就算调用方没有显式启动后台刷新，也能在首次读取时正常工作。
func (s *optionService) ensureLoaded(ctx context.Context) error {
	s.cacheMu.RLock()
	loaded := !s.loadedAt.IsZero()
	s.cacheMu.RUnlock()
	if loaded {
		return nil
	}
	return s.Refresh(ctx)
}

// refreshLoop 负责按固定周期刷新缓存。
// 如果某次刷新失败，会记录日志并在下一轮继续重试，不影响已有缓存使用。
func (s *optionService) refreshLoop(ctx context.Context) {
	timex.SetInterval(ctx, s.refreshInterval, func(runCtx context.Context) {
		if err := s.Refresh(context.Background()); err != nil {
			s.logger.ErrorCtx(runCtx, "refresh option cache failed", zap.Error(err))
		}
	})
}
