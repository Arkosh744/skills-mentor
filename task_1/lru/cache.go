package main

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	storage map[string]*Node

	head     *Node
	tail     *Node
	capacity int
}

type CacheItem struct {
	key        string
	value      any
	expiration time.Time
}

func New(capacity int) *Cache {
	return &Cache{
		storage:  make(map[string]*Node),
		capacity: capacity,
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl)

	if node, ok := c.storage[key]; ok {
		node.value.value = value
		node.value.expiration = expiration
		c.moveToFront(node)
		return
	}

	newNode := &Node{
		value: &CacheItem{
			key:        key,
			value:      value,
			expiration: expiration,
		},
	}

	if len(c.storage) >= c.capacity {
		delete(c.storage, c.tail.value.key)
		c.removeNode(c.tail)
	}

	c.addNode(newNode)
	c.storage[key] = newNode
}

func (c *Cache) Get(key string) (value any, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.storage[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(node.value.expiration) {
		delete(c.storage, key)
		c.removeNode(node)
		return nil, false
	}

	c.moveToFront(node)
	return node.value.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.storage[key]
	if ok {
		delete(c.storage, key)
		c.removeNode(node)
	}
}
