## go-urbandict

go-urbandict is a [Go](https://golang.org/) library for accessing the Urban Dictionary REST API.

### Get

Fetch and build go-urbandict:

```
go get github.com/davidscholberg/go-urbandict
```

### Library overview

```
func Define(term string) (*Definition, *Err)
func Random() (*Definition, *Err)
func DefineRaw(term string) (*DefinitionResponse, *Err)
func RandomRaw() (*DefinitionResponse, *Err)
type Definition struct { ... }
type DefinitionResponse struct { ... }
type Err struct { ... }
```

### Usage

```golang
package main

import (
    "fmt"
    "os"
    urbandict "github.com/davidscholberg/go-urbandict"
)

func main () {
    // get top definition of "1337"
    def, err := urbandict.Define("1337")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of \"1337\" definition:\n%s\n\n", def)
    fmt.Printf("Accessing individual elements:\ndef: %s\nexample: %s\n\n",
        def.Definition,
        def.Example)

    // get a random definition
    def, err = urbandict.Random()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of random definition:\n%s\n\n", def)

    // get raw response object for a search query
    defRaw, err := urbandict.DefineRaw("w00t")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of raw response to \"w00t\" query:\n%s\n\n",
        defRaw)

    // get raw response object for random word query
    defRaw, err = urbandict.RandomRaw()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of raw response to random query:\n%s\n\n",
        defRaw)
}
```
