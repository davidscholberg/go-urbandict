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
const homepageUrl = "http://www.urbandictionary.com"

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

// WordOfTheDay returns the definition for Urban Dictionary's word of the day.
func WordOfTheDay() (*Definition, error) {
    response, err := http.Get(homepageUrl)
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
