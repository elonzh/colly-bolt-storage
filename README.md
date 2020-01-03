# colly-bolt-storage ![Tag](https://img.shields.io/github/tag/elonzh/colly-bolt-storage.svg?style=flat-square) [![Go Report Card](https://goreportcard.com/badge/github.com/elonzh/colly-bolt-storage)](https://goreportcard.com/report/github.com/elonzh/colly-bolt-storage?style=flat-square) [![Build Status](https://img.shields.io/travis/elonzh/colly-bolt-storage.svg?style=flat-square)](https://travis-ci.org/elonzh/colly-bolt-storage) ![Coverage](https://img.shields.io/codecov/c/github/elonzh/colly-bolt-storage.svg?style=flat-square)

Simple and fast storage for colly, built on top of [bolt](https://github.com/etcd-io/bbolt#project-status).

# Features

Implemented interfaces:

- Storage in [github.com/gocolly/colly/queue](https://github.com/gocolly/colly/blob/master/queue/queue.go)
- Storage in [github.com/gocolly/colly/storage](https://github.com/gocolly/colly/blob/master/storage/storage.go)

## Comparison

| Projects                      | Persistence | Queue | No Service Dependency |
|-------------------------------|-------------|-------|-----------------------|
| elonzh/colly-bolt-storage     | Yes         | Yes   | Yes                   |
| velebak/colly-sqlite3-storage | Yes         | No    | Yes                   |
| gocolly/redisstorage          | Yes         | Yes   | No                    |
| zolamk/colly-mongo-storage    | Yes         | No    | No                    |

# Install

```bash
go get github.com/elonzh/colly-bolt-storage/...
```

# Example

```go
package main

import (
    "log"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/queue"
    "github.com/elonzh/colly-bolt-storage/colly/bolt"
	"go.etcd.io/bbolt"
)

func main() {
    urls := []string{
        "http://httpbin.org/",
        "http://httpbin.org/ip",
        "http://httpbin.org/cookies/set?a=b&c=d",
        "http://httpbin.org/cookies",
    }

    c := colly.NewCollector()
    path := "colly_storage.boltdb"
    var (
        db *bbolt.DB
        err error
    )
    if db, err = bbolt.Open(path, 0666, nil); err != nil {
		panic(err)
	}
    // create the storage
    storage := bolt.NewStorage(db)

    // add storage to the collector
    err = c.SetStorage(storage)
    if err != nil {
        panic(err)
    }

    // close
    defer db.Close()

    // create a new request queue
    q, _ := queue.New(2, storage)

    c.OnResponse(func(r *colly.Response) {
        log.Println("Cookies:", c.Cookies(r.Request.URL.String()))
    })

    // add URLs to the queue
    for _, u := range urls {
        q.AddURL(u)
    }
    // consume requests
    q.Run(c)
}
```

# License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Felonzh%2Fcolly-bolt-storage.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Felonzh%2Fcolly-bolt-storage?ref=badge_large)
