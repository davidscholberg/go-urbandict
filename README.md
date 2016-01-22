## go-urbandict

go-urbandict is a [Go](https://golang.org/) library for accessing the Urban Dictionary REST API.

### Get

Fetch and build go-urbandict:

```
go get github.com/davidscholberg/go-urbandict
```

### Library overview

```
func Trending() ([]string, error)
func Define(term string) (*Definition, error)
func Random() (*Definition, error)
func WordOfTheDay() (*Definition, error)
func DefineRaw(term string) (*DefinitionResponse, error)
func RandomRaw() (*DefinitionResponse, error)
type Definition struct { ... }
type DefinitionResponse struct { ... }
```

**NOTE**: The Trending and WordOfTheDay functions scrape Urban Dictionary's website since these features do not appear to be exposed in their REST API.

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

    // get the word of the day
    def, err = urbandict.WordOfTheDay()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of the word of the day:\n%s\n\n", def)

    // get a random definition
    def, err = urbandict.Random()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("JSON representation of random definition:\n%s\n\n", def)

    // get trending words
    trending, err := urbandict.Trending()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    fmt.Printf("Trending words: %v\n\n", trending)

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
