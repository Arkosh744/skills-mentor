package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	storage map[string]CacheItem
	mu      sync.Mutex
}

type CacheItem struct {
	value      any
	expiration time.Time
}

func New() *Cache {
	return &Cache{
		storage: map[string]CacheItem{},
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl)
	c.storage[key] = CacheItem{
		value:      value,
		expiration: expiration,
	}
}

func (c *Cache) Get(key string) (value any, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.storage[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		// If the item has expired, delete it and return false
		delete(c.storage, key)

		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.storage, key)
}

func main() {
	mainRace()
}

func mainRace() {
	cache := New()
	wg := sync.WaitGroup{}

	cache.Set("name", "Alex", time.Second*1)
	cache.Set("hobby", "BJJ", time.Second*2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(cache.Get("name"))
		fmt.Println(cache.Get("hobby"))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.Delete("hobby")
		fmt.Println(cache.Get("hobby"))
		cache.Set("hobby", "DoDo", time.Second*5)
		fmt.Println(cache.Get("hobby"))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 3)
		fmt.Println(cache.Get("name"))
		fmt.Println(cache.Get("hobby"))
	}()

	wg.Wait()
}
