# gin-session
A simple gin session middleware, currently supports redis, gorm store backend.

## Install

```
go get github.com/dev-shao/gin-session
```
## Quick start
```go
package main

import (
    "github.com/dev-shao/gin-session"
    "github.com/gin-gonic/gin"
)

func main(){
    //gin router
    router := gin.Default()
    //default sqlite3 database store
    router.Use(gin_session.Default())
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


## Store Backend
### GormStore
default: sqlite3
```go
//default GormStore, sqlite3 database
middlware := gin_session.Defalult()

//等同于
import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "github.com/dev-shao/gin-session"
    "github.com/dev-shao/gin-session/gorm-store"
)

func main(){
    db, _ := gorm.Open(sqlite.Open("./sessions.db"))
    store := gorm_store(db)
    middleware := gin_session.Middleware(store)
}
```
mysql
```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "github.com/dev-shao/gin-session"
    "github.com/dev-shao/gin-session/gorm-store"
)

func main(){
    dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn))
    store := gorm_store(db)
    middleware := gin_session.Middleware(store)
}

```
more database see: [gorm document](https://gorm.io/docs/connecting_to_the_database.html)

### RedisStore
```go
import (
    "github.com/dev-shao/gin-session"
    "github.com/dev-shao/gin-session/redis-store"
    "github.com/go-redis/redis"
)

func main() {
    opt := redis.Options{
        Addr: "127.0.0.1:6379",
        Password: "",
        DB: 0,
	}
	store := redis_store.New(&opt)
	middleware := gin_session.Middleware(store)
}
```
### MemoryStore
only used in test or development environment
```go
import (
    "github.com/dev-shao/gin-session"
    "github.com/dev-shao/gin-session/memory-store"
)

func main() {
	store := memory_store.New()
	middleware := gin_session.Middleware(store)
}
```
