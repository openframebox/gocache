package gocache_test

import (
	"testing"
	"time"

	"github.com/rizkykurniawan-personal/gocache"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDriver_WithTTL(t *testing.T) {
	driver := gocache.NewInMemoryDriver()
	cache := gocache.New(driver)

	type User struct {
		ID   int
		Name string
	}

	// Test Set and Get with TTL
	user := User{ID: 1, Name: "John Doe"}
	err := cache.Set("user:1", user, 1*time.Second)
	assert.NoError(t, err)

	var cachedUser User
	err = cache.Get("user:1", &cachedUser)
	assert.NoError(t, err)
	assert.Equal(t, user, cachedUser)

	// Wait for the TTL to expire
	time.Sleep(2 * time.Second)

	// Test that the value is no longer in the cache
	exists, err := cache.Exists("user:1")
	assert.NoError(t, err)
	assert.False(t, exists)

	var expiredUser User
	err = cache.Get("user:1", &expiredUser)
	assert.NoError(t, err)
	assert.Empty(t, expiredUser)
}
