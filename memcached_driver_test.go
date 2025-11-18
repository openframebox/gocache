package gocache_test

import (
	"os"
	"testing"
	"time"

	"github.com/rizkykurniawan-personal/gocache"
	"github.com/stretchr/testify/assert"
)

func TestMemcachedDriver(t *testing.T) {
	server := os.Getenv("MEMCACHED_SERVER")
	if server == "" {
		t.Skip("Skipping test; MEMCACHED_SERVER not set")
	}

	driver := gocache.NewMemcachedDriver(server)
	cache := gocache.New(driver)

	type User struct {
		ID   int
		Name string
	}

	// Test Set and Get
	user := User{ID: 1, Name: "John Doe"}
	err := cache.Set("user:1", user, 5*time.Minute)
	assert.NoError(t, err)

	var cachedUser User
	err = cache.Get("user:1", &cachedUser)
	assert.NoError(t, err)
	assert.Equal(t, user, cachedUser)

	// Test Exists
	exists, err := cache.Exists("user:1")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test Delete
	err = cache.Delete("user:1")
	assert.NoError(t, err)

	exists, err = cache.Exists("user:1")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test Clear
	err = cache.Set("user:2", user, 5*time.Minute)
	assert.NoError(t, err)

	err = cache.Clear()
	assert.NoError(t, err)

	exists, err = cache.Exists("user:2")
	assert.NoError(t, err)
	assert.False(t, exists)
}
