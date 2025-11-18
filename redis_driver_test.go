package gocache_test

import (
	"os"
	"testing"
	"time"

	"github.com/openframebox/gocache"
	"github.com/stretchr/testify/assert"
)

func TestRedisDriver(t *testing.T) {
	if os.Getenv("REDIS_SERVER") == "" {
		t.Skip("Skipping test; REDIS_SERVER not set")
	}
	redisPort := "6379"
	if os.Getenv("REDIS_PORT") != "" {
		redisPort = os.Getenv("REDIS_PORT")
	}
	driver := gocache.NewRedisDriver(gocache.RedisDriverConfig{
		Host:     os.Getenv("REDIS_SERVER"),
		Port:     redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
	})
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
