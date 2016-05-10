# go-gen-import

this is my test project

# How to Use


### before

```go
package main

import (
	"fmt"

	//+imports

	"not-edit"
	//+imports-end
)

var ()

//go:generate go-gen-import
//go:generate go fmt main.go
func main() {
	fmt.Println("this is sample")
	log.Println("this is sample by log")
}
```

### after

```go
package main

import (
	"fmt"

	//+imports

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//+imports-end
)

var ()

//go:generate go-gen-import
//go:generate go fmt main.go
func main() {
	fmt.Println("this is sample")
	log.Println("this is sample by log")
}
```


# TODO

- input/ouput file
- update import area parser
