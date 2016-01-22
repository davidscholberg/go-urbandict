// Package urbandict provides a Go wrapper for the Urban Dictionary REST API.
package urbandict

import (
    "encoding/json"
    "fmt"
    "golang.org/x/net/html"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const apiUrlFmtDefine = "http://api.urbandictionary.com/v0/define?%s"
const apiUrlRand = "http://api.urbandictionary.com/v0/random"
const wwwUrlHome = "http://www.urbandictionary.com"
const wwwUrlRand = "http://www.urbandictionary.com/random.php"

// DefinitionResponse represents the JSON response from urban dictionary.
type DefinitionResponse struct {
    List []Definition   `json:"list"`
    Result_type string  `json:"result_type"`
    Sounds []string     `json:"sounds"`
    Tags []string       `json:"tags"`
}

func (d *DefinitionResponse) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

// Definition represents a single urban dictionary definition.
type Definition struct {
    Author string       `json:"author"`
    Current_vote string `json:"current_vote"`
    Defid int           `json:"defid"`
    Definition string   `json:"definition"`
    Example string      `json:"example"`
    Permalink string    `json:"permalink"`
    Thumbs_down int     `json:"thumbs_down"`
    Thumbs_up int       `json:"thumbs_up"`
    Word string         `json:"word"`
}

func (d *Definition) String() string {
    str, err := json.MarshalIndent(d, "", "    ")
    if err != nil {
        return err.Error()
    }
    return string(str)
}

// Define gets the top definition for a search term.
func Define(term string) (*Definition, error) {
    defs, err := DefineRaw(term)
    if err != nil {
        return nil, err
    }

    if len(defs.List) == 0 {
        return nil, fmt.Errorf("no definitions for '%s' returned", term)
    }

    return &defs.List[0], nil
}

// DefineRaw gets the full response object for a search query.
func DefineRaw(term string) (*DefinitionResponse, error) {
    q := url.Values{}
    q.Add("term", term)
    apiUrl := fmt.Sprintf(apiUrlFmtDefine, q.Encode())

    return get(apiUrl)
}

// Random gets a random definition.
func Random() (*Definition, error) {
    randDefs, err := RandomRaw()
    if err != nil {
        return nil, err
    }

    if len(randDefs.List) == 0 {
        return nil, fmt.Errorf("no random definitions returned")
    }

    return &randDefs.List[0], nil
}

// RandomRaw gets a full response object for a random definition api call.
func RandomRaw() (*DefinitionResponse, error) {
    return get(apiUrlRand)
}

// Trending returns Urban Dictionary's currently trending words.
func Trending() ([]string, error) {
    response, err := http.Get(wwwUrlRand)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    doctree, err := html.Parse(response.Body)
    if err != nil {
        return nil, err
    }

    trending, err := searchForTrending(doctree)
    if err != nil {
        return nil, err
    }
    if len(trending) == 0 {
        return nil, fmt.Errorf("no trending words found")
    }

    return trending, nil
}

// WordOfTheDay returns the definition for Urban Dictionary's word of the day.
func WordOfTheDay() (*Definition, error) {
    response, err := http.Get(wwwUrlHome)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    doctree, err := html.Parse(response.Body)
    if err != nil {
        return nil, err
    }

    wotd, err := searchForWotd(doctree)
    if err != nil {
        return nil, err
    }
    if len(wotd) == 0 {
        return nil, fmt.Errorf("word of the day not found")
    }

    return Define(wotd)
}

// findChild searches html child nodes for the given type and matching data.
func findChild(n *html.Node, t html.NodeType, f func(string)bool) *html.Node {
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type == t && f(c.Data) {
            return c
        }
    }
    return nil
}

// get performs the urban dictionary api call and json parsing.
func get(apiUrl string) (*DefinitionResponse, error) {
    response, err := http.Get(apiUrl)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    defs := DefinitionResponse{}
    err = json.Unmarshal(body, &defs)
    if err != nil {
        return nil, err
    }

    return &defs, nil
}

// searchTrendingList searches the html trending list for trending words.
func searchTrendingList(n *html.Node) ([]string, error) {
    var trending []string

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        if c.Type == html.ElementNode && c.Data == "li" {
            matchFunc := func(s string)bool{return s == "a"}
            a := findChild(c, html.ElementNode, matchFunc)
            if a == nil {
                break
            }

            matchFunc = func(s string)bool{return true}
            t := findChild(a, html.TextNode, matchFunc)
            if t == nil {
                break
            }

            trending = append(trending, t.Data)
        }
    }

    return trending, nil
}

// searchForTrending searches the parsed html document for trending words.
func searchForTrending(n *html.Node) ([]string, error) {
    var trending []string
    if n.Type == html.ElementNode && n.Data == "ul" {
        for _, attr := range n.Attr {
            if attr.Key == "class" && strings.Contains(attr.Val, "trending") {
                return searchTrendingList(n)
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        trending, err := searchForTrending(c)
        if len(trending) > 0 || err != nil {
            return trending, err
        }
    }
    return trending, nil
}

// searchForWotd searches the parsed html document for the word of the day.
func searchForWotd(n *html.Node) (string, error) {
    if n.Type == html.ElementNode && n.Data == "title" {
        if n.FirstChild.Type != html.TextNode {
            err := fmt.Errorf("child of title not TextNode type")
            return "", err
        }
        parsedTitle := strings.Split(n.FirstChild.Data, ": ")
        if len(parsedTitle) != 2 || len(parsedTitle[1]) == 0 {
            err := fmt.Errorf("title text could not be parsed")
            return "", err
        }
        return parsedTitle[1], nil
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        wotd, err := searchForWotd(c)
        if len(wotd) > 0 || err != nil {
            return wotd, err
        }
    }
    return "", nil
}
