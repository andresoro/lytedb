# lytedb

lytedb is an easy to use on-disk database for Go applications. Data is divided into collections and accessed by a key. The main purpose 
of this database is to store Go data structures on disk.

# Usage

## Open DB

Simply supply a path to a file. If no file exists, it will be created automatically. Remember that this will obtain a file lock on the specified file so no other process can open it. 

``` go 

db, err := lyte.New("/path/to/file.db")
if err != nil {
    // handle
}
defer db.Close()

```

## Add/Get data

lytedb uses collections to seperate data but also uses a "main" collection for general storage. Remember that the db encodes values so only exported values get written.
``` go
type User struct {
    ID string
    Name string
    Age int
}

user := User{
    ID: "Hello World",
    Name: "Napoleon",
    Age: 25,
}

db.Add("napoleon", user)
```

Similarly, getting values writes back onto a struct. This is why values need to be exported so that the underlying encoding package (encoding/gob) can read/write from structs.

The Get method must recieve a pointer to the struct to write onto.

``` go
var user User
db.Get("napoleon", &user)
```

## Collections

Data can be divided into different groups called collections. Collections are a convenient way to seperate data while maintaining unique keys. Meaning you can have entries with the same key as long as they are in different collections.

``` go
userDB, err := db.Collection("users")

userDB.Add("user1", user)
```


