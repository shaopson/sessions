# gin-session
A simple gin session middleware, support memory, file, redis, database store backend.

## Install

```
go get github.com/dev-shao/sessions
```
## Quick start

```go
package main

import (
    "github.com/dev-shao/sessions"
    "github.com/gin-gonic/gin"
)

func main() { 
    //gin router
    router := gin.Default()
    //file store
    store := sessions.NewFileStore("./") 
    //use middleware
    router.Use(sessions.Middleware(store))
    router.GET("/test", func(context *gin.Context) {
        //get instance with context
        session, _ := sessions.Get(context)
        //set
        session.Set("key", "hello")
        //get
        value := session.Get("key")
        //delete
        session.Delete("key")
        
        context.String(200, value.(string))
    })
    router.Run()
}
```

## Store Backend
### FileStore
```go
//session file store path
path := "./"
store := sessions.NewFileStore(path)
```

### RedisStore
```go
package main

import (
    "github.com/dev-shao/sessions"
    "github.com/go-redis/redis"
)

func main() {
    opt := &redis.Options{
        Network: "tcp",
        Addr: "127.0.0.1:6379",
        DB: 0,
    }
    client := redis.NewClient(opt)
    store := sessions.NewRedisStore(client)
    middleware := sessions.Middleware(store)
}
```

### DBStore
#### sqlite3
install sqlite3 driver

```shell
go get github.com/mattn/go-sqlite
```

```go
package main

import (
    "database/sql"
    "github.com/dev-shao/sessions"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        panic(err)
    }
    store := sessions.NewDBStore("session", db)
    middleware := sessions.Middleware(store)
    ...
}
```
#### mysql

install mysql driver
```shell
go get github.com/go-sql-driver/mysql
```

```go
package main

import (
    "database/sql"
    "github.com/dev-shao/sessions"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err)
    }
    store := sessions.NewDBStore("session", db)
    middleware := sessions.Middleware(store)
    ...
}
```


### MemoryStore
only used in test or development environment
```go
middleware := sessions.Middleware(sessions.NewMemoryStore())
```

