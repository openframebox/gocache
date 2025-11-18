package gocache

import (
	"encoding/json"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedDriver struct {
	client *memcache.Client
}

func NewMemcachedDriver(servers ...string) *MemcachedDriver {
	return &MemcachedDriver{
		client: memcache.New(servers...),
	}
}

func (d *MemcachedDriver) Get(key string, dst any) error {
	item, err := d.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil
		}
		return err
	}

	return json.Unmarshal(item.Value, dst)
}

func (d *MemcachedDriver) Set(key string, value any, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return d.client.Set(&memcache.Item{
		Key:        key,
		Value:      val,
		Expiration: int32(ttl.Seconds()),
	})
}

func (d *MemcachedDriver) Delete(key string) error {
	return d.client.Delete(key)
}

func (d *MemcachedDriver) Exists(key string) (bool, error) {
	_, err := d.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *MemcachedDriver) Clear() error {
	return d.client.FlushAll()
}
