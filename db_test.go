package lyte

import (
	"os"
	"testing"
)

type user struct {
	Age  int
	Name string
}

func TestDB(t *testing.T) {

	// init DB
	path := "./data.db"
	db, err := New("./data.db")
	if err != nil {
		t.Fatal("DB will not open. Gives error: ", err)
	}
	defer func() {
		db.Close()
		os.Remove(path)
	}()

	testUser := user{
		Age:  25,
		Name: "Napoleon",
	}

	err = db.Add("users", "user1", testUser)
	if err != nil {
		t.Fatal("Error adding user to database: ", err)
	}

	var returnUser user

	err = db.Get("users", "user1", &returnUser)
	if err != nil {
		t.Fatal("Error getting user from database: ", err)
	}

	if testUser != returnUser {
		t.Fatal("Not returning the same user")
	}

}
