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

	user1 := user{
		Name: "Queen Elizabeth",
		Age:  25,
	}

	err = db.Add("user1", user1)
	if err != nil {
		t.Fatalf("Add failed with error: %s", err)
	}

	var returnUser user
	err = db.Get("user1", &returnUser)
	if err != nil {
		t.Fatalf("Get failed with error: %s", err)
	}

	if returnUser != user1 {
		t.Fatal("Return struct not the same as input")
	}

	err = db.Delete("user1")
	if err != nil {
		t.Fatal("Error on delete")
	}

	err = db.Get("user1", user{})
	if err == nil {
		t.Fatal("db should throw error when key doesnt exist")
	}
}

func TestCollection(t *testing.T) {
	// init DB
	path := "./data.db"
	db, err := New(path)
	if err != nil {
		t.Fatal("DB will not open. Gives error: ", err)
	}
	defer func() {
		db.Close()
		os.Remove(path)
	}()

	col, err := db.Collection("users")
	if err != nil {
		t.Fatalf("Could not create collection: %s", err)
	}

	user1 := &user{
		Name: "Andres",
		Age:  21,
	}

	err = col.Add("user1", user1)
	if err != nil {
		t.Fatalf("Could not add to collection: %s", err)
	}

	var returnUser user
	err = col.Get("user1", &returnUser)
	if err != nil {
		t.Fatalf("Could not get from collection: %s", err)
	}
}

func TestUniqueKeysOnCollections(t *testing.T) {
	path := "./data.db"
	db, err := New(path)
	if err != nil {
		t.Fatal("DB will not open. Gives error: ", err)
	}
	defer func() {
		db.Close()
		os.Remove(path)
	}()

	col1, err := db.Collection("collection1")
	if err != nil {
		t.Fatal("Could not create collection: ", err)
	}

	col2, err := db.Collection("collection2")
	if err != nil {
		t.Fatal("Could not create second collection: ", err)
	}

	// should be able to use the same key in different collections

	err = col1.Add("key", "value")
	if err != nil {
		t.Fatal("Error adding key/value to store: ", err)
	}

	err = col2.Add("key", "value")
	if err != nil {
		t.Fatal("Error adding same key to different collection: ", err)
	}
}
