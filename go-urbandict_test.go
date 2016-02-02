// Package urbandict provides a Go wrapper for the Urban Dictionary REST API.
package urbandict

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestSearchForWotd(t *testing.T) {
	doctree, err := html.Parse(strings.NewReader(homePage))
	if err != nil {
		t.Error(err)
	}

	wotd, err := searchForWotd(doctree)
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(wotd, "Icicle Fingers") != 0 {
		t.Errorf("expected \"Icicle Fingers\", got \"%s\"", wotd)
	}
}

func TestSearchForTrending(t *testing.T) {
	doctree, err := html.Parse(strings.NewReader(randomPage))
	if err != nil {
		t.Error(err)
	}

	trendingList, err := searchForTrending(doctree)
	if err != nil {
		t.Error(err)
	}
	for i, _ := range trendingList {
		if strings.Compare(trendingList[i], expectedTrendingList[i]) != 0 {
			t.Errorf("expected \"%s\", got \"%s\"",
				expectedTrendingList[i],
				trendingList[i])
		}
	}
}

var homePage string = `
<!DOCTYPE html>
<html>
<head>
<meta charset='UTF-8'>
<title>Urban Dictionary, January 29: Icicle Fingers</title>
</head>
<body></body>
</html>`

var randomPage string = `
<!DOCTYPE html>
<html>
<head>
<meta charset='UTF-8'>
<title>Urban Dictionary, January 29: Icicle Fingers</title>
</head>
<body>
<div class='panel'>
<ul class='no-bullet trending'>
<li><a href="/define.php?term=netflix+and+chill">netflix and chill</a></li>
<li><a href="/define.php?term=cleveland+steamer">cleveland steamer</a></li>
<li><a href="/define.php?term=tubgirl">tubgirl</a></li>
<li><a href="/define.php?term=rimjob">rimjob</a></li>
<li><a href="/define.php?term=dabbing">dabbing</a></li>
<li><a href="/define.php?term=dirty+sanchez">dirty sanchez</a></li>
<li><a href="/define.php?term=alabama+hot+pocket">alabama hot pocket</a></li>
<li><a href="/define.php?term=donkey+punch">donkey punch</a></li>
<li><a href="/define.php?term=blumpkin">blumpkin</a></li>
<li><a href="/define.php?term=dabbin%27">dabbin&#39;</a></li>
</ul>
</div>
</body>
</html>`

var expectedTrendingList []string = []string{
	"netflix and chill",
	"cleveland steamer",
	"tubgirl",
	"rimjob",
	"dabbing",
	"dirty sanchez",
	"alabama hot pocket",
	"donkey punch",
	"blumpkin",
	"dabbin'"}
