package test

import (
	"github.com/dev-shao/sessions"
	"reflect"
	"testing"
)

func TestFileStore_Set(t *testing.T) {
	key := "test-key"
	value := sessions.StoreValue{
		"string": "string",
		"int":    99,
		"float":  3.14,
	}
	store := sessions.NewFileStore("./")
	err := store.Set(key, value, 60)
	if err != nil {
		t.Error(err)
	}
	//get
	if v, err := store.Get(key); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, value) {
		t.Error("Get value error")
	}
}

func TestFileStore_Get(t *testing.T) {
	store := sessions.NewFileStore("./")
	//expired
	if v, err := store.Get("expired-key"); err != nil && v["int"] == 0 {
		t.Error(err)
	}
	//not exists
	if _, err := store.Get("not-exists-key"); err != sessions.DoesNotExists {
		t.Error(err)
	}
	//content error
	if _, err := store.Get("content-error-key"); err != sessions.InvalidData {
		t.Error(err)
	}
}

func TestFileStore_Delete(t *testing.T) {
	store := sessions.NewFileStore("./")
	key := "test-key"
	value := sessions.StoreValue{
		"string": "string",
		"int":    99,
		"float":  3.14,
	}
	err := store.Set(key, value, 60)
	if err != nil {
		t.Error(err)
	}

	if store.Exists(key) != true {
		t.Error("Set fail, value not exists")
	}
	store.Delete(key)
	if store.Exists(key) == true {
		t.Error("Delete fail")
	}

}
