# Feeds Finder

[![Build Status](https://travis-ci.org/gilliek/go-feedsfinder.png?branch=master)](https://travis-ci.org/gilliek/go-feedsfinder)

Feeds finder is Go package that provides functions for finding RSS/Atom feeds in a web page.

## Installation

```go get github.com/gilliek/go-feedsfinder/feeds```

## Usage

Typical usage:

```go
package main

import (
	"fmt"
	"log"

	"github.com/gilliek/go-feedsfinder/feeds"
)

func main() {
	links, err := feeds.FindFromURL("http://www.example.com")
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println(links)

    //...
}
```

Output:

```
[{http://blog.golang.org/feed.atom atom}]
```

See the document for more details.

## Documentation

Document can be found on [GoWalker](https://gowalker.org/github.com/gilliek/go-feedsfinder/feeds) 
or [GoDoc](http://godoc.org/github.com/gilliek/go-feedsfinder/feeds)

## License

BSD 3-clauses
