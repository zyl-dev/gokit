package cache

import (
	"fmt"
	"time"
)

import (
	"context"
)

type CacheGenerator struct {
	cache    *Cache
	distLock *DistributedLock
}

func NewCacheGenerator(cache *Cache, distLock *DistributedLock) *CacheGenerator {
	return &CacheGenerator{
		cache:    cache,
		distLock: distLock,
	}
}

// 生成缓存数据（带分布式锁）
func (cg *CacheGenerator) GenerateWithLock(
	ctx context.Context,
	cacheKey string,
	lockKey string,
	lockTimeout time.Duration,
	cacheExpiry time.Duration,
	fn func() (string, error), // 业务逻辑函数
) (string, error) {
	// 1. 先从缓存中获取数据
	cachedData, err := cg.cache.Get(ctx, cacheKey)
	if err == nil {
		return cachedData, nil
	}

	// 2. 尝试获取分布式锁
	lockValue, acquired, err := cg.distLock.AcquireLock(ctx, lockKey, lockTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock: %v", err)
	}
	if !acquired {
		// 如果获取锁失败，等待一段时间后重试
		time.Sleep(100 * time.Millisecond)
		return cg.GenerateWithLock(ctx, cacheKey, lockKey, lockTimeout, cacheExpiry, fn)
	}

	// 3. 启动锁续期
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel()
	go cg.distLock.RenewLock(ctxCancel, lockKey, lockValue, lockTimeout)

	// 4. 确保锁最终被释放
	defer cg.distLock.ReleaseLock(ctx, lockKey, lockValue)

	// 5. 再次检查缓存中是否有数据，防止其他服务已经更新了缓存
	cachedData, err = cg.cache.Get(ctx, cacheKey)
	if err == nil {
		return cachedData, nil
	}

	// 6. 执行业务逻辑生成缓存数据
	data, err := fn()
	if err != nil {
		return "", fmt.Errorf("failed to generate cache data: %v", err)
	}

	// 7. 将数据写入缓存
	err = cg.cache.Set(ctx, cacheKey, data, cacheExpiry)
	if err != nil {
		return "", fmt.Errorf("failed to set cache data: %v", err)
	}
	return data, nil
}
