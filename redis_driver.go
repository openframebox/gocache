package gocache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDriver struct {
	client *redis.Client
}

// NewRedisDriver creates a new Redis driver
func NewRedisDriver(config RedisDriverConfig) *RedisDriver {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
	return &RedisDriver{client: rdb}
}

// Get retrieves a value from the cache
func (d *RedisDriver) Get(key string, dest any) error {
	val, err := d.client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in the cache
func (d *RedisDriver) Set(key string, value any, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return d.client.Set(context.Background(), key, val, ttl).Err()
}

// Delete removes a value from the cache
func (d *RedisDriver) Delete(key string) error {
	return d.client.Del(context.Background(), key).Err()
}

// Exists checks if a key exists in the cache
func (d *RedisDriver) Exists(key string) (bool, error) {
	val, err := d.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

// Clear clears the entire cache
func (d *RedisDriver) Clear() error {
	return d.client.FlushDB(context.Background()).Err()
}
