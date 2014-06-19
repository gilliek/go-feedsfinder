// Copyright 2014 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feeds

import (
	"io"
	"net/http"
	"net/http/httptest"
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
	links, err := Find([]byte(htmlCode), "")
	if err != nil {
		t.Fatal(links)
	}

	testResults(t, links)
}

func TestFindFromURL(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, htmlCode)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	links, err := FindFromURL(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	testResults(t, links)
}

func TestFindFromFile(t *testing.T) {
	links, err := FindFromFile(os.Getenv("GOPATH")+"/src/github.com/gilliek/go-feedsfinder/testdata/index.html", "")
	if err != nil {
		t.Fatal(links)
	}

	testResults(t, links)
}

func TestFormatLink(t *testing.T) {
	href1 := "https://www.foo.com/atom.xml"
	baseURL1 := "http://www.foo.com"

	if res, _ := formatLink(href1, baseURL1); res != href1 {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", href1, res)
	}

	const expected = "http://www.foo.com/atom.xml"

	href2 := "http://www.foo.com/atom.xml"
	baseURL2 := "http://www.foo.com"

	if res, _ := formatLink(href2, baseURL2); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

	href3 := "/atom.xml"
	baseURL3 := "http://www.foo.com/"

	if res, _ := formatLink(href3, baseURL3); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

	href4 := "/atom.xml"
	baseURL4 := "http://www.foo.com"

	if res, _ := formatLink(href4, baseURL4); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

	href5 := "atom.xml"
	baseURL5 := "http://www.foo.com"

	if res, _ := formatLink(href5, baseURL5); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

	href6 := "atom.xml"
	baseURL6 := "http://www.foo.com/index.html"

	if res, _ := formatLink(href6, baseURL6); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

	href7 := "atom.xml"
	baseURL7 := "http://www.foo.com/bar/index"

	if res, _ := formatLink(href7, baseURL7); res != "http://www.foo.com/bar/atom.xml" {
		t.Errorf("Invalid feed link: expected 'http://www.foo.com/bar/atom.xml', found '%s'", res)
	}

	href8 := "/atom.xml"
	baseURL8 := "http://www.foo.com/bar/index"

	if res, _ := formatLink(href8, baseURL8); res != expected {
		t.Errorf("Invalid feed link: expected '%s', found '%s'", expected, res)
	}

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
