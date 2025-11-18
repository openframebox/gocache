package gocache

import (
	"encoding/json"
	"sync"
	"time"
)

type InMemoryDriver struct {
	mu    sync.RWMutex
	items map[string]item
}

type item struct {
	Value      []byte
	Expiration int64
}

func NewInMemoryDriver() *InMemoryDriver {
	d := &InMemoryDriver{
		items: make(map[string]item),
	}
	go d.startGC()
	return d
}

func (d *InMemoryDriver) Get(key string, dst any) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	it, ok := d.items[key]
	if !ok || it.isExpired() {
		return nil
	}

	return json.Unmarshal(it.Value, dst)
}

func (d *InMemoryDriver) Set(key string, value any, ttl time.Duration) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var expiration int64
	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	d.items[key] = item{
		Value:      val,
		Expiration: expiration,
	}
	return nil
}

func (d *InMemoryDriver) Delete(key string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.items, key)
	return nil
}

func (d *InMemoryDriver) Exists(key string) (bool, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	it, ok := d.items[key]
	if !ok || it.isExpired() {
		return false, nil
	}

	return true, nil
}

func (d *InMemoryDriver) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.items = make(map[string]item)
	return nil
}

func (i item) isExpired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

func (d *InMemoryDriver) startGC() {
	for {
		time.Sleep(1 * time.Minute)
		d.mu.Lock()
		for key, it := range d.items {
			if it.isExpired() {
				delete(d.items, key)
			}
		}
		d.mu.Unlock()
	}
}
