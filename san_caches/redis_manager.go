package san_caches

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

type redisCacheManager struct {
	c *redis.Client
}

func NewRedisCacheManager(c *redis.Client) CacheManager {
	return &redisCacheManager{c: c}
}

func (cm *redisCacheManager) Set(ctx context.Context, key CacheKey, value any, ttl time.Duration) error {
	keyStr, err := key.GetKey()
	if err != nil {
		return err
	}

	protoValue, ok := value.(proto.Message)
	if ok {
		bytes, err := proto.Marshal(protoValue)
		if err != nil {
			return err
		}
		return cm.c.Set(ctx, keyStr, bytes, ttl).Err()
	}

	return cm.c.Set(ctx, keyStr, value, ttl).Err()
}

func (cm *redisCacheManager) Get(ctx context.Context, key CacheKey, value any) error {
	keyStr, err := key.GetKey()
	if err != nil {
		return err
	}

	protoValue, ok := value.(proto.Message)
	if ok {
		rawBytes := []byte{}
		err = cm.c.Get(ctx, keyStr).Scan(&rawBytes)
		if err != nil {
			return err
		}
		return proto.Unmarshal(rawBytes, protoValue)
	}

	return cm.c.Get(ctx, keyStr).Scan(value)
}

func (cm *redisCacheManager) Del(ctx context.Context, key CacheKey) error {
	keyStr, err := key.GetKey()
	if err != nil {
		return err
	}
	return cm.c.Del(ctx, keyStr).Err()
}

func (cm *redisCacheManager) DelNamespace(ctx context.Context, namespace string) error {
	pattern := namespace + "*"
	var cursor uint64
	for {
		keys, next, err := cm.c.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := cm.c.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}
