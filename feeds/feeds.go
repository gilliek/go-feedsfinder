// Copyright 2014 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package feeds provides functions for finding RSS/Atom feeds in a web page.
*/
package feeds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"code.google.com/p/go.net/html"
)

const (
	rssMIME  = "application/rss+xml"
	atomMIME = "application/atom+xml"
)

// Link holds information about a link to a RSS or Atom feed.
type Link struct {
	// URL contains the reference to the RSS/Atom feed.
	URL string

	// Type is the type of the feed. It can be either "rss" or "atom".
	Type string
}

// Find finds RSS/Atom feeds in a web page given as a byte slice.
//
// baseURL is the URL of the web page. This is used to deal with absolute path.
func Find(b []byte, baseURL string) ([]Link, error) {
	var links []Link

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return links, err
	}

	parse(doc, &links, baseURL)

	return links, nil
}

// FindFromURL finds RSS/Atom feeds in a web page given as an URL.
func FindFromURL(url string) ([]Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []Link{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return []Link{}, fmt.Errorf("Invalid URL: HTTP status %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Link{}, err
	}

	return Find(b, url)
}

// FindFromFile finds RSS/Atom feeds in a web page given as a file path.
//
// baseURL is the URL of the web page. This is used to deal with absolute path.
func FindFromFile(filePath string, baseURL string) ([]Link, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []Link{}, err
	}

	return Find(b, baseURL)
}

// parse recursively parses a HTML page.
func parse(n *html.Node, links *[]Link, baseURL string) {
	if n.Type == html.ElementNode && n.Data == "body" {
		return
	}

	if n.Type == html.ElementNode && n.Data == "link" {
		var hrefAttr, typeAttr string

		for _, attr := range n.Attr {
			switch key := attr.Key; key {
			case "type":
				if attr.Val == rssMIME {
					typeAttr = "rss"
				} else if attr.Val == atomMIME {
					typeAttr = "atom"
				}
			case "href":
				hrefAttr = attr.Val
			}
		}

		if hrefAttr != "" && typeAttr != "" {
			// TODO error handling
			url, _ := formatLink(hrefAttr, baseURL)

			*links = append(*links, Link{URL: url, Type: typeAttr})
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, links, baseURL)
	}
}

func formatLink(href, baseURL string) (string, error) {
	if len(href) == 0 {
		return baseURL, nil
	}

	if len(href) > 7 && (href[0:7] == "http://" || href[0:8] == "https://") {
		return href, nil
	}

	url, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if href[0] == '/' {
		return url.Scheme + "://" + url.Host + href, nil
	}

	if baseURL[len(baseURL)-1] == '/' {
		return baseURL + href, nil
	}

	path := filepath.Dir(url.Path)

	if path == "." || path == "/" {
		return url.Scheme + "://" + url.Host + "/" + href, nil
	}

	return url.Scheme + "://" + url.Host + path + "/" + href, nil
}
