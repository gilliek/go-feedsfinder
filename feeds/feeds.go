// Copyright 2014 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feeds

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.google.com/p/go.net/html"
)

const (
	rssMIME  = "application/rss+xml"
	atomMIME = "application/atom+xml"
)

type Link struct {
	URL  string
	Type string
}

func Find(b []byte) ([]Link, error) {
	var links []Link

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return links, err
	}

	parse(doc, &links)

	return links, nil
}

func FindFromURL(url string) ([]Link, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []Link{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return []Link{}, errors.New(fmt.Sprintf("Invalid URL: HTTP status %s", resp.Status))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Link{}, err
	}

	return Find(b)
}

func FindFromFile(filePath string) ([]Link, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []Link{}, err
	}

	return Find(b)
}

func parse(n *html.Node, links *[]Link) {
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
			*links = append(*links, Link{URL: hrefAttr, Type: typeAttr})
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parse(c, links)
	}
}
