package test

import (
	"github.com/dev-shao/sessions"
	"github.com/go-redis/redis"
	"reflect"
	"testing"
	"time"
)

var opt = &redis.Options{
	Network: "tcp",
	Addr:    "127.0.0.1:6379",
	DB:      9,
}
var client = redis.NewClient(opt)
var store = sessions.NewRedisStore(client)

func TestRedisStore_Set(t *testing.T) {
	key := "test-key"
	value := sessions.StoreValue{
		"string": "hello",
		"int":    124,
		"float":  3.14,
	}
	if err := store.Set(key, value, 60); err != nil {
		t.Error(err)
	}
	//get
	if v, err := store.Get(key); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, value) {
		t.Error("Get value error")
	}

}

func TestRedisStoreGet(t *testing.T) {
	// not exists
	if _, err := store.Get("not-exists-key"); err != sessions.DoesNotExists {
		t.Error(err)
	}
	// content error
	client.Set("content-error-key", "sdjfl23rfeds", time.Second*60)
	if _, err := store.Get("content-error-key"); err != sessions.InvalidData {
		t.Error(err)
	}

}
