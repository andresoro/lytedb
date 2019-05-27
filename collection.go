package lyte

// Col refers to a collection for the DB
type Col struct {
	name string
	db   *DB
}

// Add a key/value pair to DB under this collection
func (c *Col) Add(key string, value interface{}) error {
	return c.db.add(c.name, key, value)
}

// Get will write the value associated with the key to the given interface
func (c *Col) Get(key string, value interface{}) error {
	return c.db.get(c.name, key, value)
}
