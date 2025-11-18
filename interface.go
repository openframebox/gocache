package gocache

import "time"

type Driver interface {
	Get(key string, dst any) error
	Set(key string, value any, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Clear() error
}

type RedisDriverConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}
