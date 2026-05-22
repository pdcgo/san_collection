package san_caches

import (
	"context"
	"errors"
	"time"
)

type skipCacheManager struct{}

func NewSkipCacheManager() CacheManager {
	return &skipCacheManager{}
}

func (s *skipCacheManager) Set(ctx context.Context, key CacheKey, value any, ttl time.Duration) error {
	return nil
}

func (s *skipCacheManager) Get(ctx context.Context, key CacheKey, value any) error {
	return errors.New("cache not found")
}

func (s *skipCacheManager) Del(ctx context.Context, key CacheKey) error {
	return nil
}
