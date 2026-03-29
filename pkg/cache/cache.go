package cache

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"gin-template/pkg/errs"
	"gin-template/pkg/timex"
)

// ErrLoaderNotConfigured 表示当前缓存实例未配置全量加载器，
// 因此无法执行 Refresh 或 Run 的初始化加载逻辑。
var ErrLoaderNotConfigured = errors.New("缓存加载器未配置")

// Loader 定义全量缓存加载器。
// 调用方需要返回一份完整快照，Refresh 会用它整体替换当前缓存内容。
type Loader[K ~string, V any] func(ctx context.Context) (map[K]V, error)

// Options 定义缓存组件的初始化参数。
// KeyPrefix 用于在内部形成统一命名空间，Refresh 时也会用它过滤加载结果。
type Options[K ~string, V any] struct {
	KeyPrefix       string
	RefreshInterval time.Duration
	Loader          Loader[K, V]
}

// Cache 维护一份线程安全的内存快照，并支持按需全量刷新。
// 组件对外始终暴露业务 key，内部会在需要时自动应用 key 前缀。
type Cache[K ~string, V any] struct {
	keyPrefix       string
	refreshInterval time.Duration
	loader          Loader[K, V]

	cacheMu sync.RWMutex
	items   map[K]V

	// refreshMu 串行化 Refresh / Set / Del，避免刷新替换与写操作互相覆盖。
	refreshMu sync.Mutex
	// startMu 与 started 共同保证后台循环只会启动一次。
	startMu sync.Mutex
	started bool
}

// New 创建一个最小可用的缓存实例。
// 构造过程不会触发 loader，避免在包初始化阶段产生副作用。
func New[K ~string, V any](opts Options[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		keyPrefix:       opts.KeyPrefix,
		refreshInterval: opts.RefreshInterval,
		loader:          opts.Loader,
		items:           make(map[K]V),
	}
}

// Run 启动缓存生命周期。
// 首次调用会先同步刷新一次；若配置了刷新间隔，则继续启动后台定时刷新。
func (c *Cache[K, V]) Run(ctx context.Context) error {
	c.startMu.Lock()
	defer c.startMu.Unlock()

	if c.started {
		return nil
	}
	if err := c.Refresh(ctx); err != nil {
		return err
	}

	c.started = true
	if c.refreshInterval <= 0 {
		return nil
	}

	go c.runLoop(ctx)
	return nil
}

// Get 读取单个缓存项。
// 返回 false 表示当前快照中没有命中目标 key。
func (c *Cache[K, V]) Get(key K) (V, bool) {
	lookupKey := c.buildKey(key)

	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	value, ok := c.items[lookupKey]
	return value, ok
}

// GetAll 返回当前缓存快照的副本。
// 调用方拿到的是业务 key 视图，修改返回值不会影响内部缓存。
func (c *Cache[K, V]) GetAll() map[K]V {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	cloned := make(map[K]V, len(c.items))
	for key, value := range c.items {
		cloned[c.exposeKey(key)] = value
	}
	return cloned
}

// Set 将单个缓存项写入当前快照。
// 写入只影响内存态，后续 Refresh 会以 loader 返回的全量快照为准。
func (c *Cache[K, V]) Set(key K, value V) {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	c.cacheMu.Lock()
	c.items[c.buildKey(key)] = value
	c.cacheMu.Unlock()
}

// Refresh 执行一次全量刷新，并以新的快照整体替换当前缓存内容。
// 当配置了 keyPrefix 时，只会保留匹配此前缀的加载结果。
func (c *Cache[K, V]) Refresh(ctx context.Context) error {
	if c.loader == nil {
		return errs.WithStack(ErrLoaderNotConfigured)
	}

	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	items, err := c.loader(ctx)
	if err != nil {
		return err
	}

	nextItems := make(map[K]V, len(items))
	for key, value := range items {
		if c.keyPrefix != "" && !strings.HasPrefix(string(key), c.keyPrefix) {
			continue
		}
		nextItems[key] = value
	}

	c.cacheMu.Lock()
	c.items = nextItems
	c.cacheMu.Unlock()
	return nil
}

// Del 从当前快照中删除指定 key。
// 删除只影响内存态，后续 Refresh 仍可能把该 key 从 loader 结果重新带回。
func (c *Cache[K, V]) Del(key K) {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	c.cacheMu.Lock()
	delete(c.items, c.buildKey(key))
	c.cacheMu.Unlock()
}

func (c *Cache[K, V]) runLoop(ctx context.Context) {
	defer errs.Recover()

	timex.SetInterval(ctx, c.refreshInterval, func(runCtx context.Context) {
		_ = c.Refresh(runCtx)
	})
}

func (c *Cache[K, V]) buildKey(key K) K {
	if c.keyPrefix == "" {
		return key
	}
	return K(c.keyPrefix + string(key))
}

func (c *Cache[K, V]) exposeKey(key K) K {
	if c.keyPrefix == "" {
		return key
	}
	return K(strings.TrimPrefix(string(key), c.keyPrefix))
}
