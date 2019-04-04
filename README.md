# colly-bolt-storage ![Tag](https://img.shields.io/github/tag/earlzo/colly-bolt-storage.svg?style=flat-square) [![Go Report Card](https://goreportcard.com/badge/github.com/earlzo/colly-bolt-storage)](https://goreportcard.com/report/github.com/earlzo/colly-bolt-storage?style=flat-square) [![Build Status](https://img.shields.io/travis/earlzo/colly-bolt-storage.svg?style=flat-square)](https://travis-ci.org/earlzo/colly-bolt-storage) ![Coverage](https://img.shields.io/codecov/c/github/earlzo/colly-bolt-storage.svg?style=flat-square&token=3fb8c6d7-8912-4083-9c19-38c577228b70)

Simple and fast storage for colly, built on top of [bolt](https://github.com/etcd-io/bbolt#project-status).

# Features

Implemented interfaces:

- Storage in [github.com/gocolly/colly/queue](https://github.com/gocolly/colly/blob/master/queue/queue.go)
- Storage in [github.com/gocolly/colly/storage](https://github.com/gocolly/colly/blob/master/storage/storage.go)

## Comparison

| Projects                      | Persistence | Queue | No Service Dependency |
|-------------------------------|-------------|-------|-----------------------|
| earlzo/colly-bolt-storage     | Yes         | Yes   | Yes                   |
| velebak/colly-sqlite3-storage | Yes(actually is unusable, see https://github.com/velebak/colly-sqlite3-storage/pull/3)         | No    | Yes  |
| gocolly/redisstorage          | Yes         | Yes   | No                    |
| zolamk/colly-mongo-storage    | Yes         | No    | No                    |

# Install

```bash
go get github.com/earlzo/colly-bolt-storage/...
```

# Example

```go
package main

import (
    "log"

    "github.com/gocolly/colly"
    "github.com/gocolly/colly/queue"
    "github.com/earlzo/colly-bolt-storage/colly/bolt"
)

func main() {
    urls := []string{
        "http://httpbin.org/",
        "http://httpbin.org/ip",
        "http://httpbin.org/cookies/set?a=b&c=d",
        "http://httpbin.org/cookies",
    }

    c := colly.NewCollector()

    // create the storage
    storage := &bolt.Storage{
            Path: "test.boltdb",
    }

    // add storage to the collector
    err := c.SetStorage(storage)
    if err != nil {
        panic(err)
    }

    // close
    defer storage.DB.Close()

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

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage?ref=badge_large)
