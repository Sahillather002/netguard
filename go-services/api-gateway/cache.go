package main

import (
	"sync"
	"time"
)

// CacheItem represents a cached item
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// Cache represents an in-memory cache
type Cache struct {
	items map[string]*CacheItem
	mu    sync.RWMutex
}

var cache = &Cache{
	items: make(map[string]*CacheItem),
}

// Set stores a value in cache with expiration
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration)
	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all items from cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*CacheItem)
}

// CleanExpired removes expired items
func (c *Cache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}

// GetStats returns cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	expired := 0
	now := time.Now()
	for _, item := range c.items {
		if now.After(item.Expiration) {
			expired++
		}
	}

	return map[string]interface{}{
		"total_items":   len(c.items),
		"expired_items": expired,
		"active_items":  len(c.items) - expired,
	}
}

// Start cache cleanup routine
func startCacheCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			cache.CleanExpired()
		}
	}()
}

// Invalidate cache for specific resource
func invalidateCache(pattern string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for key := range cache.items {
		if len(key) >= len(pattern) && key[:len(pattern)] == pattern {
			delete(cache.items, key)
		}
	}
}
