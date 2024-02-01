package cache

import "errors"

type Cacher[K comparable, V any] interface {
	Get(key K) (value V, err error)
	Put(key K, value V) (err error)
}

// Concrete LRU cache
type lruCache[K comparable, V any] struct {
	size      int
	remaining int
	cache     map[K]V
	queue     []K
}

// Constructor
func NewCacher[K comparable, V any](size int) Cacher[K, V] {
	return &lruCache[K, V]{size: size, remaining: size, cache: make(map[K]V), queue: make([]K, 0)}
}

// Fetch function, fetchs value and moves to tail of queue
func (c *lruCache[K, V]) Get(key K) (value V, err error) {
	v, ok := c.cache[key]
	if !ok {
		return v, errors.New("key not found")
	}

	// Move the key to the tail of the queue
	c.deleteFromQueue(key)
	c.queue = append(c.queue, key)

	return v, nil
}

// Put Function, puts value in cache and evicts if necessarys
func (c *lruCache[K, V]) Put(key K, value V) (err error) {
	// Check if key already exists
	if _, ok := c.cache[key]; ok {
		// Update the value
		c.cache[key] = value
		// Move key to the tail of the queue
		c.deleteFromQueue(key)
		c.queue = append(c.queue, key)
		return nil
	}

	// Check capacity and evict if needed
	if c.remaining == 0 {
		// Evict the least recently used element
		evictKey := c.queue[0]
		delete(c.cache, evictKey)
		c.deleteFromQueue(evictKey)
	} else {
		c.remaining--
	}

	// Add the new key-value pair
	c.cache[key] = value
	c.queue = append(c.queue, key)
	return nil
}

// Helper method to delete all occurrences of a key from the queue
func (c *lruCache[K, V]) deleteFromQueue(key K) {
	newQueue := make([]K, 0, c.size)
	for _, k := range c.queue {
		if k != key {
			newQueue = append(newQueue, k)
		}
	}
	c.queue = newQueue
}
