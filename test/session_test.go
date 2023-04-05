package test

import (
	"github.com/dev-shao/sessions"
	"testing"
)

func TestSession(t *testing.T) {
	store := sessions.NewFileStore("./")
	// key not exists
	session, err := sessions.Load("key123", store)
	if err != nil {
		t.Error(err)
	}
	//set
	session.Set("name", "ling")
	session.Set("age", 18)
	//get
	v := session.Get("age")
	if i := v.(int); i != 18 {
		t.Error("session get value error")
	}
	//save
	if err := session.Save(); err != nil {
		t.Error(err)
	}
	// load
	s, err := sessions.Load(session.Key(), store)
	if err != nil {
		t.Error(err)
	}
	v = s.Get("name")
	if name := v.(string); name != "ling" {
		t.Error("load session error")
	}

}
