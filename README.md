# lytedb

lytedb is an easy to use on-disk database for Go applications. Data is divided into collections and accessed by a key. The main purpose 
of this database is to store Go data structures on disk.

## Example

``` go
import "github.com/andresoro/lytedb"

// Only exported fields are encoded
type Example struct {
    Data    string
    Num     int
}

func main() {

    db := lyte.New("/path/to/database/")

    data := Example{
        Data: "Hello World",
        Num: 42,
    }

    db.Add("examples", "example1", data)

    // Get from DB into the new struct
    var new Example
    db.Get("examples", "example1", &new)

}
```
