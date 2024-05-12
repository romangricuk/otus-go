package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type queueItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cacheItem, isExist := cache.items[key]

	if isExist {
		cacheItem.Value = queueItem{
			key:   key,
			value: value,
		}
		cache.queue.MoveToFront(cacheItem)
	} else {
		cache.items[key] = cache.queue.PushFront(queueItem{
			key:   key,
			value: value,
		})
	}

	if cache.queue.Len() > cache.capacity {
		back := cache.queue.Back()
		cache.queue.Remove(back)
		delete(cache.items, back.Value.(queueItem).key)
	}

	return isExist
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cacheItem, isExist := cache.items[key]

	if isExist {
		cache.queue.MoveToFront(cacheItem)

		return cacheItem.Value.(queueItem).value, isExist
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
