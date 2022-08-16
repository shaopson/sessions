package gin_session

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)


var SessionKey = "Session"
var SessionCookieName = "SessionId"
var SessionExpireAge time.Duration = 604800 //seconds, default 7 days


func Middleware(store Store) gin.HandlerFunc {
	return func(context *gin.Context) {
		key, _ := context.Cookie(SessionCookieName)
		session := New(key, SessionExpireAge, store)
		if err := session.Load(); err != nil {
			panic(err)
		}
		context.Set(SessionKey, session)
		context.Next()
		if err := session.Save(); err != nil {
			panic(err)
		}
		context.SetCookie(SessionCookieName, session.key, int(SessionExpireAge), "/", "", true, false)
	}
}


type Store interface {
	Get(key string) (map[string]interface{}, error)
	Set(key string, value map[string]interface{}, expire time.Duration) error
	Delete(key string)
	Exists(key string) bool
	GetExpireTime(key string) time.Time
	SetExpireTime(key string, duration time.Duration)
}


type Session struct {
	key string
	data map[string]interface{}
	expireAge time.Duration
	store Store
}

func From(ctx *gin.Context) *Session {
	value, exists := ctx.Get(SessionKey)
	if !exists {
		return nil
	}
	return value.(*Session)
}


func New(key string, expire time.Duration, store Store) *Session {
	return &Session{
		key: key,
		data: make(map[string]interface{}),
		expireAge: expire * time.Second,
		store: store,
	}
}

func (self *Session) Load() error {
	if self.key == "" {
		return nil
	}
	data, err := self.store.Get(self.key)
	if data == nil || err != nil {
		return err
	}
	for k, v := range data {
		self.data[k] = v
	}
	return nil
}

func (self *Session) Get(key string) interface{} {
	return self.data[key]
}

func (self *Session) DefaultGet(key string, defaultValue interface{}) interface{} {
	if value, exists := self.data[key]; !exists {
		return defaultValue
	} else {
		return value
	}
}

func (self *Session) Set(key string, value interface{}) {
	self.data[key] = value
}

func (self *Session) Delete(key string) {
	delete(self.data, key)
}

func (self *Session) Clear() {
	for k, _ := range self.data {
		delete(self.data, k)
	}
}

func (self *Session) Save() error {
	if self.key == "" {
		self.key = self.newKey()
	}
	return self.store.Set(self.key, self.data, self.expireAge)
}

// refresh key
func (self *Session) Refresh() {
	self.key = self.newKey()
}

func (self *Session) GetExpireTime() time.Time {
	return self.store.GetExpireTime(self.key)
}

func (self *Session) SetExpireAge(expire time.Duration) {
	self.expireAge = expire
}

var chars = "1234567890abcdefghijklnmopqrstuvwxyzABCDEFGHIJKLNMOPQRSTUVWSYZ#$*"

func (self *Session) newKey() string {
	rand.Seed(time.Now().UnixNano())
	length := len(chars)
	buf := strings.Builder{}
	for true {
		for i := 0; i < 32; i++ {
			x := rand.Intn(length)
			c := chars[x:x+1]
			buf.WriteString(c)
		}
		key := buf.String()
		if !self.store.Exists(key) {
			return key
		}
	}
	return ""
}


