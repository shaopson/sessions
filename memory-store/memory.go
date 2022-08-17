package memory_store


import "time"

type storeValue map[string]interface{}

type MemoryStore struct {
	cache map[string]storeValue
	cacheExpire map[string]time.Time
}

func New() *MemoryStore {
	return &MemoryStore{
		cache: make(map[string]storeValue),
		cacheExpire: make(map[string]time.Time),
	}
}

func (self *MemoryStore) Get(key string) (map[string]interface{}, error) {
	now := time.Now()
	if expireTime := self.cacheExpire[key]; now.After(expireTime) {
		self.Delete(key)
		return nil, nil
	}
	value, exists := self.cache[key]
	if !exists{
		return nil, nil
	}
	return value, nil
}

func (self *MemoryStore) Set(key string, data map[string]interface{}, expire time.Duration) error {
	if expire <= 0 {
		self.Delete(key)
		return nil
	}
	self.cacheExpire[key] = time.Now().Add(expire)
	self.cache[key] = data
	return nil
}

func (self *MemoryStore) Delete(key string) {
	delete(self.cacheExpire, key)
	delete(self.cache, key)
}

func (self *MemoryStore) Exists(key string) bool {
	_, exists := self.cache[key]
	return exists
}

func (self *MemoryStore) GetExpireTime(key string) time.Time {
	return self.cacheExpire[key]
}

func (self *MemoryStore) SetExpireTime(key string, duration time.Duration) {
	if duration <= 0 {
		self.Delete(key)
	}
	self.cacheExpire[key] = time.Now().Add(duration)
}

func (self *MemoryStore) ClearExpired() {
	now := time.Now()
	for k, expireTime := range self.cacheExpire {
		if now.After(expireTime) {
			self.Delete(k)
		}
	}
}
