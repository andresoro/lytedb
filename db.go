package lyte

import (
	"bytes"
	"encoding/gob"
	"errors"

	"go.etcd.io/bbolt"
)

// DB is made up of 'collections' which hold certain 'types' which are stored on disk.
// It is up to the user to define these collections and add the relevant data
type DB struct {
	path        string
	collections map[string]*bbolt.Bucket
	bolt        *bbolt.DB
}

// New returns a new db object
func New(path string) (*DB, error) {
	bolt, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, errors.New("Unable to open DB file")
	}

	db := &DB{
		path:        path,
		collections: make(map[string]*bbolt.Bucket),
		bolt:        bolt,
	}

	return db, nil
}

// Close underlying DB
func (db *DB) Close() error {
	return db.bolt.Close()
}

// Add data specified by a key to a given collection
func (db *DB) Add(col string, key string, data interface{}) error {

	// encode key/values
	k := []byte(key)
	var v bytes.Buffer
	gob.NewEncoder(&v).Encode(data)

	err := db.bolt.Update(func(tx *bbolt.Tx) error {
		// if bucket doesnt exist, create it and add to map
		bucket, ok := db.collections[col]
		if !ok {
			b, err := tx.CreateBucket([]byte(col))
			if err != nil {
				return err
			}
			db.collections[col] = b

			// add key/val to bucket
			err = b.Put(k, v.Bytes())
			if err != nil {
				return nil
			}
			return nil
		}

		err := bucket.Put(k, v.Bytes())
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Get will decode the data under the collection/key onto the value passed
// will return an error if not found or otherwise
func (db *DB) Get(col, key string, value interface{}) error {

	err := db.bolt.View(func(tx *bbolt.Tx) error {
		bucket, ok := db.collections[col]
		if !ok {
			return errors.New("Collection does not exist")
		}

		c := bucket.Cursor()
		k, v := c.Seek([]byte(key))
		if k == nil || string(k) != key {
			return errors.New("Key not found")
		}

		data := bytes.NewReader(v)
		return gob.NewDecoder(data).Decode(value)
	})

	return err
}
