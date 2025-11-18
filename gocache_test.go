package gocache_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rizkykurniawan-personal/gocache"
)

func TestGoCache_Remember(t *testing.T) {
	if os.Getenv("REDIS_SERVER") == "" {
		t.Skip("Skipping test; REDIS_SERVER not set")
	}
	redisPort := "6379"
	if os.Getenv("REDIS_PORT") != "" {
		redisPort = os.Getenv("REDIS_PORT")
	}
	// Create the Redis driver
	// NOTE: Please configure the test with your Redis connection details.
	driver := gocache.NewRedisDriver(gocache.RedisDriverConfig{
		Host:     os.Getenv("REDIS_SERVER"),
		Port:     redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	// Create the GoCache instance
	cache := gocache.New(driver)

	// Create a struct to store the user data
	type User struct {
		ID   int
		Name string
	}

	// Use the Remember function to get the user data
	var user User
	err := cache.Remember("user:1", 5*time.Minute, &user, func() (any, error) {
		// This function will be called if the data is not in the cache
		t.Log("Fetching user from the database...")
		return User{ID: 1, Name: "John Doe"}, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("User: %+v\n", user)

	// Try to get the user data again. This time, it should be fetched from the cache
	var cachedUser User
	err = cache.Remember("user:1", 5*time.Minute, &cachedUser, func() (any, error) {
		t.Log("This should not be printed")
		return User{}, nil
	})

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Cached User: %+v\n", cachedUser)
}
