package cache

import "errors"

type Cacher[K comparable, V any] interface {
    Get(key K) (value V, err error)
    Put(key K, value V) (err error)
}

// Concrete LRU cache
type lruCache[K comparable, V any] struct {
    size int
    remaining int
    cache map[K]V
    queue []K
}

// Constructor
func NewCacher[K comparable, V any](size int) Cacher[K, V] {
    return &lruCache[K, V]{size: size, remaining: size, cache: make(map[K]V), queue: make([]K, 0)}
}

func (c *lruCache[K, V]) Get(key K) (value V, err error) {
	// Your code here ...

    // Hint - Move key to tail of queue (mark as recently used)

}

func (c *lruCache[K, V]) Put(key K, value V) (err error) {
	// Your code here ...

    // Hint - Check if key already exists
    

	// Hint - Check capacity, and evict if needed

 
    // Hint - Add new key-value pair


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
