package redis_store

import (
	"bytes"
	"encoding/gob"
	"github.com/go-redis/redis"
	"io"
	"strings"
	"time"
)


func New(opt *redis.Options) *RedisStore {
	store := &RedisStore{}
	store.client = redis.NewClient(opt)
	return store
}


type RedisStore struct {
	client *redis.Client
}

func (self *RedisStore) Get(key string) (map[string]interface{}, error) {
	result := self.client.Get(key)
	value, err := result.Result()
	if err != nil {
		return nil, err
	}
	return self.decode(value)
}

func (self *RedisStore) Set(key string, value map[string]interface{}, expiration time.Duration) error {
	s, err := self.encode(value)
	if err != nil {
		return err
	}
	result := self.client.Set(key, s, expiration)
	return result.Err()
}

func (self *RedisStore) Delete(key string) {
	self.client.Del(key)
}

func (self *RedisStore) Exists(key string) bool {
	result := self.client.Exists(key)
	if result.Val() == 1 {
		return true
	}
	return false
}

func (self *RedisStore) GetExpireTime(key string) time.Time {
	result := self.client.TTL(key)
	ttl := result.Val()
	if ttl < 0 {
		return time.Time{}
	}
	return time.Now().Add(ttl)
}

func (self *RedisStore) SetExpireTime(key string, duration time.Duration) {
	self.client.Expire(key, duration)
}

func (self *RedisStore) ClearExpired() {}

func (self *RedisStore) encode(data map[string]interface{}) (string, error) {
	writer := bytes.Buffer{}
	encoder := gob.NewEncoder(&writer)
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}

func (self *RedisStore) decode(src string) (map[string]interface{}, error) {
	reader := strings.NewReader(src)
	decoder := gob.NewDecoder(reader)
	data := make(map[string]interface{})
	err := decoder.Decode(&data)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF{
		return nil, err
	}
	return data, nil
}
