package redis

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Cache struct {
	rDB *redis.Client
}

func NewRedisCache(cfg Config) *Cache {
	return &Cache{
		rDB: redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
			Password: cfg.Password,
			DB:       cfg.Db,
			PoolSize: cfg.PoolSize,
		})}
}

func (cache *Cache) Ping(ctx context.Context) error {
	result, err := cache.rDB.Ping(ctx).Result()
	if err != nil {
		return errors.Wrap(err, "failed to ping redis")
	}
	if result != "PONG" {
		return errors.New("redis ping failed")
	}
	return nil
}

func (cache *Cache) Disconnect() error {
	if cache.rDB != nil {
		return cache.rDB.Close()
	}
	return nil
}

func (cache *Cache) SetJson(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	result := cache.rDB.Set(ctx, key, bytes, expiration)
	return result.Err()
}

func (cache *Cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	result := cache.rDB.Set(ctx, key, value, expiration)
	return errors.Wrapf(result.Err(), "failed to set key %s", key)
}

func (cache *Cache) Get(ctx context.Context, key string) (string, error) {
	result, err := cache.rDB.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", ErrKeyNotFound
	case err != nil:
		return "", err
	case result == "":
		return "", ErrKeyNotFound
	}
	return result, nil
}

func (cache *Cache) GetJSON(ctx context.Context, key string, value interface{}) error {
	result := cache.rDB.Get(ctx, key)
	storedBytes, err := result.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(storedBytes, &value)
}

func (cache *Cache) Delete(ctx context.Context, key string) error {
	return cache.rDB.Del(ctx, key).Err()
}
