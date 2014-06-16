// Copyright 2014 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feeds

import (
	"os"
	"testing"
)

const (
	atomURL  = "http://www.example.com/atom_feed.xml"
	rssURL   = "http://www.example.com/rss_feed.xml"
	htmlCode = `<!DOCTYPE html>
<html>
    <head>
        <title>Foobar</title>
        <link href="http://www.example.com/atom_feed.xml" type="application/atom+xml"/>
        <link href="http://www.example.com/rss_feed.xml" type="application/rss+xml"/>
    </head>
    <body>
    </body>
</html>`
)

func TestFind(t *testing.T) {
	links, err := Find([]byte(htmlCode))
	if err != nil {
		t.Fatal(links)
	}

	testResults(t, links)
}

func TestFindFromFile(t *testing.T) {
	links, err := FindFromFile(os.Getenv("GOPATH") + "/src/github.com/gilliek/go-feedsfinder/testdata/index.html")
	if err != nil {
		t.Fatal(links)
	}

	testResults(t, links)
}

func testResults(t *testing.T, links []Link) {
	if nbLinks := len(links); nbLinks != 2 {
		t.Fatalf("Invalid number of links: expected 2, found %d", nbLinks)
	}

	atom := links[0]
	rss := links[1]

	if atom.URL != atomURL {
		t.Errorf("Invalid Atom feed URL: expected '%s', found '%s'", atomURL, atom.URL)
	}

	if atom.Type != "atom" {
		t.Errorf("Invalid Atom feed type: expected 'atom', found '%s'", atom.Type)
	}

	if rss.URL != rssURL {
		t.Errorf("Invalid RSS feed URL: expected '%s', found '%s'", rssURL, rss.URL)
	}

	if rss.Type != "rss" {
		t.Errorf("Invalid RSS feed type: expected 'rss', found '%s'", rss.Type)
	}

}
