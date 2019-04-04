# colly-bolt-storage [![license](https://img.shields.io/github/license/earlzo/colly-bolt-storage.svg?style=flat-square)](https://github.com/earlzo/colly-bolt-storage/blob/master/LICENSE)[![Build Status](https://img.shields.io/travis/earlzo/colly-bolt-storage.svg?style=flat-square)](https://travis-ci.org/earlzo/colly-bolt-storage)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage?ref=badge_shield)

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

# Usage

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


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fearlzo%2Fcolly-bolt-storage?ref=badge_large)