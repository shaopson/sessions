package test

import (
	"database/sql"
	"github.com/dev-shao/sessions"
	_ "github.com/mattn/go-sqlite3"
	"reflect"
	"testing"
)

var db *sql.DB
var dbStore sessions.Store

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	dbStore = sessions.NewDBStore("sqlite3", db)

}

func TestDBStoreSet(t *testing.T) {
	key := "test-key"
	value := sessions.StoreValue{
		"string": "string",
		"int":    99,
		"float":  3.14,
	}
	if err := dbStore.Set(key, value, 60); err != nil {
		t.Fatal(err)
	}
	if v, err := dbStore.Get(key); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(v, value) {
		t.Fatal("Set value error")
	}
}

func TestDBStoreGet(t *testing.T) {
	value := sessions.StoreValue{
		"string": "string",
		"int":    99,
		"float":  3.14,
	}
	//not exists
	if _, err := dbStore.Get("not-exists-key"); err != sessions.DoesNotExists {
		t.Error(err)
	}
	//expired
	if err := dbStore.Set("expired-key", value, -60); err != nil {
		t.Error(err)
	}
	if v, err := dbStore.Get("expired-key"); err != sessions.DoesNotExists || v != nil {
		t.Error("get expired data")
	}

}
