package gocache

import "time"

type GoCache struct {
	driver Driver
}

func New(driver Driver) *GoCache {
	gc := GoCache{
		driver: driver,
	}

	return &gc
}

func (gc *GoCache) Get(key string, dst any) error {
	return gc.driver.Get(key, dst)
}

func (gc *GoCache) Set(key string, value any, ttl time.Duration) error {
	return gc.driver.Set(key, value, ttl)
}

func (gc *GoCache) Delete(key string) error {
	return gc.driver.Delete(key)
}

func (gc *GoCache) Exists(key string) (bool, error) {
	return gc.driver.Exists(key)
}

func (gc *GoCache) Clear() error {
	return gc.driver.Clear()
}

func (gc *GoCache) Remember(key string, ttl time.Duration, dst any, fn func() (any, error)) error {
	err := gc.Get(key, dst)
	if err == nil {
		return nil
	}

	val, err := fn()
	if err != nil {
		return err
	}

	return gc.Set(key, val, ttl)
}
