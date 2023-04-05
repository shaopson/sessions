package sessions

import (
	"errors"
	"time"
)

var DoesNotExists = errors.New("Object does not exists")
var InvalidData = errors.New("Invalid Store Data")

type StoreValue map[string]interface{}

type Store interface {
	//if returns error is not nil, that StoreValue will be a nil map pointer
	Get(key string) (StoreValue, error)

	//if expire <= 0, Store maybe not store the value
	Set(key string, value StoreValue, expire int32) error

	Delete(key string)

	Exists(key string) bool

	GetExpireTime(key string) time.Time

	SetExpireTime(key string, t time.Time)
}
