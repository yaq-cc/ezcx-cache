package cache

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound   = errors.New("key not found")
	ErrEndOfIterator = errors.New("end of iterator")
)

type Cache[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
	l  int
}

type Cacheable[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Delete(key K)
	Pop(key K) (V, bool)
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		m: make(map[K]V),
	}
}

func (c *Cache[K, V]) set(key K, value V) {
	c.m[key] = value
	c.l++
}

func (c *Cache[K, V]) get(key K) (V, bool) {
	value, ok := c.m[key]
	return value, ok
}

func (c *Cache[K, V]) del(key K) bool {
	_, ok := c.get(key)
	if ok {
		delete(c.m, key)
		c.l--
	}
	return ok
	
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	c.set(key, value)
	c.mu.Unlock()
}


func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	value, ok := c.get(key)
	c.mu.RUnlock()
	return value, ok
}

func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	ok := c.del(key)
	c.mu.Unlock()
	return ok
}

func (c *Cache[K, V]) Pop(key K) (V, bool) {
	c.mu.Lock()
	value, ok := c.get(key)
	c.del(key)
	c.mu.Unlock()
	return value, ok
}

func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.l)
	i := 0
	c.mu.RLock()
	for key := range c.m {
		keys[i] = key
		i++
	}
	c.mu.RUnlock()
	return keys
}

func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.l)
	i := 0
	c.mu.RLock()
	for _, value := range c.m {
		values[i] = value
		i++
	}
	c.mu.RUnlock()
	return values
}

// // Avoided.  There may be cases where the value of type V
// // is not comparable.
// func (c *Cache[K, V]) Update(key K, value V) {
// 	c.mu.Lock()
// 	old, ok := c.get(key)
// 	if !ok || old != value {
// 		c.set(key, value)
// 	}
// 	c.mu.Unlock()
// }

// // These methods are skipped 
// func (c *Cache[K, V]) Items() *Items[K, V] {
// 	return &Items[K, V]{
// 		c:    c,
// 		keys: c.Keys(),
// 	}
// }

// type Items[K comparable, V any] struct {
// 	c    *Cache[K, V]
// 	keys []K
// 	i    int
// }

// func (it *Items[K, V]) Next() (*Item[K, V], error) {
// 	if it.i < len(it.keys) {
// 		key := it.keys[it.i]
// 		value, ok := it.c.Get(key)
// 		if !ok {
// 			return nil, ErrKeyNotFound
// 		}
// 		it.i++
// 		return &Item[K, V]{Key: key, Value: value}, nil
// 	} else {
// 		return nil, ErrEndOfIterator
// 	}
// }

// type Item[K comparable, V any] struct {
// 	Key   K
// 	Value V
// }
