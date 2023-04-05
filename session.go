package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

type Session struct {
	key    string
	data   map[string]interface{}
	expire int32
	store  Store
	mutex  sync.RWMutex
}

var SessionName = "sessions.Session"
var SessionCookieName = "SessionId"
var SessionExpire int32 = 604800 //second; 604800 = 7 days
var CookieSecure = false
var CookieHttpOnly = false
var CookiePath = "/"
var CookieDomain = ""

func Middleware(store Store) gin.HandlerFunc {
	return func(context *gin.Context) {
		sessionId, _ := context.Cookie(SessionCookieName)
		session, err := Load(sessionId, store)
		if err != nil {
			context.Error(err)
			return
		}
		context.Set(SessionName, session)
		context.SetCookie(SessionCookieName, session.Key(), int(session.expire), CookiePath, CookieDomain, CookieSecure, CookieHttpOnly)
		context.Next()
		if err = session.Save(); err != nil {
			context.Error(err)
		}
	}
}

func Get(context *gin.Context) (*Session, error) {
	if value, exists := context.Get(SessionName); exists {
		if session, ok := value.(*Session); ok {
			return session, nil
		} else {
			return nil, fmt.Errorf("Get Session fail, Please do not modify the value of gin.Context key('%s')\r\n", SessionName)
		}
	}
	return nil, fmt.Errorf("Session not found. Did you forget to register the session middleware or modify the value of gin.Context key('%s')\r\n", SessionName)
}

// get a session from the store by key.
// If the key is invalid or the session does not exist, return the session of the new key
func Load(key string, store Store) (session *Session, err error) {
	session = &Session{
		data:   map[string]interface{}{},
		expire: SessionExpire,
		store:  store,
	}
	if len(key) != 32 {
		return
	}
	var data StoreValue
	if data, err = store.Get(key); err != nil {
		if err == DoesNotExists || err == InvalidData {
			err = nil
		}
		return
	} else {
		session.data = data
	}
	return
}

func New(store Store) *Session {
	return &Session{
		key:    newKey(store),
		data:   map[string]interface{}{},
		expire: SessionExpire,
		store:  store,
	}
}

func (session *Session) Key() string {
	if session.key == "" {
		session.key = newKey(session.store)
	}
	return session.key
}

func (session *Session) Get(key string) interface{} {
	session.mutex.RLock()
	defer session.mutex.RUnlock()
	return session.data[key]
}

func (session *Session) Set(key string, value interface{}) {
	session.mutex.Lock()
	defer session.mutex.Unlock()
	session.data[key] = value
}

func (session *Session) Delete(key string) {
	delete(session.data, key)
}

func (session *Session) Clear() {
	session.data = make(map[string]interface{})
}

func (session *Session) Save() error {
	session.mutex.RLock()
	defer session.mutex.RUnlock()
	return session.store.Set(session.Key(), session.data, session.expire)
}

// clear data and refresh key
func (session *Session) Refresh(context *gin.Context) {
	session.Clear()
	if session.key != "" {
		session.store.Delete(session.key)
	}
	session.key = newKey(session.store)
	context.SetCookie(SessionCookieName, session.Key(), int(session.expire), CookiePath, CookieDomain, CookieSecure, CookieHttpOnly)
}

func (session *Session) ExpireAge() int32 {
	return session.expire
}

func (session *Session) SetExpireAge(expire int32) {
	session.expire = expire
}

func newKey(store Store) (key string) {
	b := make([]byte, 32)
	for {
		rand.Read(b)
		key = base64.RawURLEncoding.EncodeToString(b)[:32]
		if !store.Exists(key) {
			return
		}
	}
}
