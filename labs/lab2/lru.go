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
// Get retrieves a value from the cache using a specific key. If the key is found, it moves the key to the tail of the queue to mark it as recently used.
func (c *lruCache[K, V]) Get(key K) (value V, err error) {
    // Attempt to find the value by key in the cache map
    v, ok := c.cache[key]
    // If the key is not found, return an error indicating the absence of the key
    if !ok {
        return v, errors.New("key not found")
    }

    // If key is found, update its position to the tail of the queue to mark it as recently used
    // First, remove the key from its current position in the queue
    c.deleteFromQueue(key)
    // Then, append it to the tail of the queue
    c.queue = append(c.queue, key)

    // Return the found value without error
    return v, nil
}

// Put inserts or updates a value in the cache under a specific key. If the cache is full, it evicts the least recently used item before adding the new item.
func (c *lruCache[K, V]) Put(key K, value V) (err error) {
    // Check if the key already exists in the cache
    if _, ok := c.cache[key]; ok {
        // If it does, update the value associated with the key
        c.cache[key] = value
        // Move the key to the tail of the queue to mark it as recently used
        c.deleteFromQueue(key)
        c.queue = append(c.queue, key)
        return nil
    }

    // If the cache is at capacity (remaining slots are 0), evict the least recently used item
    if c.remaining == 0 {
        // Evict the least recently used element, which is at the front of the queue
        evictKey := c.queue[0]
        // Remove the evicted key from the cache map and the queue
        delete(c.cache, evictKey)
        c.deleteFromQueue(evictKey)
    } else {
        // If there's space, decrement the counter of remaining slots
        c.remaining--
    }

    // Add the new key-value pair to the cache and the queue
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
