package san_caches

import (
	"context"
	"time"
)

type CacheKey interface {
	GetKey() (string, error)
}

type CacheManager interface {
	Set(ctx context.Context, key CacheKey, value any, ttl time.Duration) error
	Get(ctx context.Context, key CacheKey, value any) error
	Del(ctx context.Context, key CacheKey) error
	// DelNamespace deletes every cached key that starts with namespace.
	DelNamespace(ctx context.Context, namespace string) error
}
