package sessions

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (store *RedisStore) Get(key string) (StoreValue, error) {
	result := store.client.Get(key)
	if err := result.Err(); err != nil {
		if err == redis.Nil {
			return nil, DoesNotExists
		}
		return nil, err
	}
	value := StoreValue{}
	data, err := result.Bytes()
	if err != nil {
		return value, err
	}
	err = Decode(data, &value)
	return value, InvalidData
}

func (store *RedisStore) Set(key string, value StoreValue, expire int32) error {
	if expire <= 0 {
		return nil
	}
	data, err := Encode(value)
	if err != nil {
		return err
	}
	store.client.Set(key, data, time.Duration(expire)*time.Second)
	return nil
}

func (store *RedisStore) Delete(key string) {
	store.client.Del(key)
}

func (store *RedisStore) Exists(key string) bool {
	result := store.client.Exists(key)
	return result.Val() == 1
}

func (store *RedisStore) GetExpireTime(key string) time.Time {
	result := store.client.PTTL(key)
	now := time.Now()
	return now.Add(result.Val())
}

func (store *RedisStore) SetExpireTime(key string, t time.Time) {
	store.client.ExpireAt(key, t)
}
