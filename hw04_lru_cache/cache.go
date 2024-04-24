package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
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
	cacheItem, isExist := cache.items[key]
	if isExist {
		cacheItem.Value = value
		cache.queue.MoveToFront(cacheItem)
	} else {
		cache.items[key] = cache.queue.PushFront(value)
	}

	if cache.queue.Len() > cache.capacity {
		back := cache.queue.Back()
		cache.queue.Remove(back)

		// todo переделать на O(1)
		for cacheKey, item := range cache.items {
			if item == back {
				delete(cache.items, cacheKey)
				break
			}
		}
	}

	return isExist
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cacheItem, isExist := cache.items[key]

	if isExist {
		cache.queue.MoveToFront(cacheItem)

		return cacheItem.Value, isExist
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
