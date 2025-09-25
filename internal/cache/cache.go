package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Get(ctx context.Context, client redis.Client, key string) (string, error) {
	data := client.Get(ctx, key)
	return data.Result()
}

func Set(ctx context.Context, client redis.Client, key string, val []byte, ttl time.Duration) error {
	if err := client.Set(
		ctx,
		key,
		val,
		ttl,
	).Err(); err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, client redis.Client, key string) error {
	return client.Del(ctx, key).Err()
}
