package sessions

import (
	"time"
)

type MemoryStore struct {
	dataCache   map[string]StoreValue
	expireCache map[string]time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		dataCache:   make(map[string]StoreValue),
		expireCache: make(map[string]time.Time),
	}
}

func (store *MemoryStore) Get(key string) (StoreValue, error) {
	//已过期，删除
	if t := store.expireCache[key]; t.Before(time.Now()) {
		store.Delete(key)
		return nil, DoesNotExists
	}
	value, exists := store.dataCache[key]
	if !exists {
		return nil, DoesNotExists
	}
	return value, nil
}

func (store *MemoryStore) Set(key string, value StoreValue, expire int32) error {
	if expire <= 0 {
		store.Delete(key)
		return nil
	}
	now := time.Now()
	store.dataCache[key] = value
	store.expireCache[key] = now.Add(time.Second * time.Duration(expire))
	return nil
}

func (store *MemoryStore) Delete(key string) {
	delete(store.dataCache, key)
	delete(store.expireCache, key)
}

func (store *MemoryStore) Exists(key string) bool {
	_, exists := store.dataCache[key]
	return exists
}

func (store *MemoryStore) GetExpireTime(key string) time.Time {
	return store.expireCache[key]
}

func (store *MemoryStore) SetExpireTime(key string, expire time.Time) {
	store.expireCache[key] = expire
}
