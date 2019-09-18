package lyte

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"

	"go.etcd.io/bbolt"
)

// DB is made up of 'collections' which hold certain 'types' which are stored on disk.
// It is up to the user to define these collections and add the relevant data
type DB struct {
	path        string
	collections map[string]*Col
	bolt        *bbolt.DB
	mu          *sync.Mutex
}

// New returns a new db object
func New(path string) (*DB, error) {

	// underlying database
	bolt, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, errors.New("Unable to open DB file")
	}

	// main DB struct
	db := &DB{
		path:        path,
		collections: make(map[string]*Col),
		bolt:        bolt,
		mu:          &sync.Mutex{},
	}

	// main collection to hold general values
	err = db.createBucket("main")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Close underlying DB
func (db *DB) Close() error {
	return db.bolt.Close()
}

// Collection will return a reference to the underlying collection
// or create a new one if already exists
func (db *DB) Collection(name string) (*Col, error) {
	col, ok := db.collections[name]
	if ok {
		return col, nil
	}

	db.mu.Lock()

	col = &Col{
		name: name,
		db:   db,
	}
	db.collections[name] = col

	db.mu.Unlock()

	err := db.createBucket(name)

	return col, err
}

// Add data specified by a key to the "main" collection
func (db *DB) Add(key string, data interface{}) error {
	return db.add("main", key, data)
}

// Get will retreive the key from the "main" collection and write the data
// to the data interface
func (db *DB) Get(key string, data interface{}) error {
	return db.get("main", key, data)
}

// Delete an entry from DB
func (db *DB) Delete(key string) error {
	return db.delete("main", key)
}

func (db *DB) add(col, key string, data interface{}) error {
	// encode key/values
	k := []byte(key)
	var v bytes.Buffer
	gob.NewEncoder(&v).Encode(data)

	return db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(col))
		err := bucket.Put(k, v.Bytes())
		if err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) get(col, key string, value interface{}) error {

	return db.bolt.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(col))
		k := []byte(key)
		v := bucket.Get(k)

		data := bytes.NewReader(v)
		err := gob.NewDecoder(data).Decode(value)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *DB) delete(col, key string) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(col))
		k := []byte(key)

		err := bucket.Delete(k)
		return err
	})
}

func (db *DB) createBucket(name string) error {
	return db.bolt.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("Cannot create collection: %s", err)
		}

		return err
	})
}
