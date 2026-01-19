package jsondb

import "testing"

type user struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestAddUpdate(t *testing.T) {
	db, err := NewDatabase[user]("users.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	db.AddOrUpdate(func(u user) bool {
		return u.Id == 21
	}, user{Id: 21, Name: "Jason", Age: 46})

	db.AddOrUpdate(func(u user) bool {
		return u.Id == 22
	}, user{Id: 22, Name: "Erin", Age: 43})

	db.AddOrUpdate(func(u user) bool {
		return u.Id == 100
	}, user{Id: 100, Name: "Remove me", Age: 100})

	db.Store()
}

func TestExists(t *testing.T) {
	db, err := NewDatabase[user]("users.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	exists := db.Exists(func(u user) bool {
		return u.Id == 21
	})

	if !exists {
		t.Fail()
	}
	t.Log("exists", exists)
}

func TestFilter(t *testing.T) {
	db, err := NewDatabase[user]("users.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	age46 := db.Filter(func(u user) bool {
		return u.Age == 46
	})

	t.Log(age46)
}

func TestRemove(t *testing.T) {
	db, err := NewDatabase[user]("users.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	elem, found := db.Remove(func(u user) bool {
		return u.Id == 100
	})
	if !found {
		t.Fail()
	}
	t.Log("removed", elem)
}
