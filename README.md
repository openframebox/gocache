# GoCache

GoCache is a flexible and extensible caching library for Go. It provides a simple and consistent API for interacting with various caching backends, while allowing you to easily create your own custom cache drivers.

## Features

-   **Simple API**: Easy-to-use methods for common caching operations like `Get`, `Set`, `Delete`, `Exists`, and `Clear`.
-   **Extensible**: Easily create your own custom cache drivers by implementing a simple `Driver` interface.
-   **Ready-to-use Drivers**: Comes with several ready-to-use drivers for popular caching backends.
-   **Built-in TTL**: Support for Time-To-Live (TTL) on cached items.

## Installation

```bash
go get github.com/openframebox/gocache
```

## Usage

Here's a simple example of how to use GoCache with the in-memory driver:

```go
package main

import (
	"fmt"
	"time"

	"github.com/openframebox/gocache"
)

func main() {
	// Create the in-memory driver
	driver := gocache.NewInMemoryDriver()

	// Create the GoCache instance
	cache := gocache.New(driver)

	// Create a struct to store user data
	type User struct {
		ID   int
		Name string
	}

	// Set a value in the cache
	user := User{ID: 1, Name: "John Doe"}
	err := cache.Set("user:1", user, 5*time.Minute)
	if err != nil {
		panic(err)
	}

	// Get a value from the cache
	var cachedUser User
	err = cache.Get("user:1", &cachedUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Cached User: %+v\n", cachedUser)
}
```

## GoCache Functions

The `gocache.Cache` struct provides the following functions:

-   **`New(driver Driver) *Cache`**: Creates a new `Cache` instance with the specified driver.
-   **`Remember(key string, ttl time.Duration, dest any, callback func() (any, error)) error`**: Retrieves an item from the cache. If the item does not exist, the `callback` function is executed, its result is stored in the cache, and then returned.
-   **`Get(key string, dest any) error`**: Retrieves a value from the cache and unmarshals it into `dest`.
-   **`Set(key string, value any, ttl time.Duration) error`**: Stores a value in the cache with a specified Time-To-Live (TTL).
-   **`Delete(key string) error`**: Removes a value from the cache.
-   **`Exists(key string) (bool, error)`**: Checks if a key exists in the cache.
-   **`Clear() error`**: Clears the entire cache.

## Initialization

To use GoCache, you first need to initialize a driver. Then, you can create a new GoCache instance with the driver.

```go
// Create a driver (e.g., in-memory)
driver := gocache.NewInMemoryDriver()

// Create the GoCache instance
cache := gocache.New(driver)
```

## Available Drivers

GoCache comes with the following ready-to-use drivers:

### In-Memory

The in-memory driver stores cache items in memory. This is useful for single-process applications or for testing.

**Initialization:**

```go
driver := gocache.NewInMemoryDriver()
```

### Redis

The Redis driver uses a Redis server as the caching backend.

**Initialization:**

```go
driver := gocache.NewRedisDriver(gocache.RedisDriverConfig{
    Host:     "localhost",
    Port:     "6379",
    Password: "", // your redis password
    DB:       0,  // your redis db
})
```

### Memcached

The Memcached driver uses a Memcached server as the caching backend.

**Initialization:**

```go
driver := gocache.NewMemcachedDriver("localhost:11211")
```

## Extending with Custom Drivers

You can easily extend GoCache by creating your own custom driver. A custom driver must implement the `gocache.Driver` interface.

**Driver Interface:**

```go
type Driver interface {
	Get(key string, dst any) error
	Set(key string, value any, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Clear() error
}
```

**Example Custom Driver:**

Here's an example of a simple file-based cache driver:

```go
package gocache

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type FileDriver struct {
	dir string
}

func NewFileDriver(dir string) *FileDriver {
	os.MkdirAll(dir, 0755)
	return &FileDriver{dir: dir}
}

func (d *FileDriver) Get(key string, dst any) error {
	data, err := ioutil.ReadFile(filepath.Join(d.dir, key))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func (d *FileDriver) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(d.dir, key), data, 0644)
}

func (d *FileDriver) Delete(key string) error {
	return os.Remove(filepath.Join(d.dir, key))
}

func (d *FileDriver) Exists(key string) (bool, error) {
	_, err := os.Stat(filepath.Join(d.dir, key))
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (d *FileDriver) Clear() error {
	return os.RemoveAll(d.dir)
}
```
