# gin-session
A simple gin session middleware, currently supports redis, gorm store backend.

## 快速入门

```
go get github.com/dev-shao/gin-session
```


```go
package main

import (
	gin_session "github.com/Dev-Shao/gin-session"
	"github.com/gin-gonic/gin"
	"github.com/Dev-Shao/gin-session/redis-store"
	"github.com/go-redis/redis"
)

func main(){
	//redis options
	opt := redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	}
	//new redis store
	store := redis_store.New(&opt)
	middleware := gin_session.Middleware(store)
	//gin router
	router := gin.Default()
	router.Use(middleware)
	router.GET("/test", func(context *gin.Context) {
		//init session from context
		session := gin_session.From(context)
		//set
		session.Set("key","value")
		//get
		session.Get("key")
		//delete
		session.Delete("key")
		//save
		session.Save()
	})
	
}

```

