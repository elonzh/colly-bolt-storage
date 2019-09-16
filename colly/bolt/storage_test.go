package bolt

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gocolly/colly/queue"
	"github.com/gocolly/colly/storage"
	"go.etcd.io/bbolt"
)

func TestStorage(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test_colly_storage.boltdb")
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		panic(err)
	}
	s := NewStorage(db)
	var _ queue.Storage = s
	var _ storage.Storage = s

	defer func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}()

	if err := s.Init(); err != nil {
		t.Error("failed to initialize client: " + err.Error())
	}
	// test visit
	var requestID uint64 = 1
	if isVisited, err := s.IsVisited(requestID); isVisited != false || err != nil {
		t.Fatal("unexpected result", isVisited, err)
	}
	if err := s.Visited(requestID); err != nil {
		t.Fatal("unexpected result", err)
	}
	if isVisited, err := s.IsVisited(requestID); isVisited != true || err != nil {
		t.Fatal("unexpected result", isVisited, err)
	}

	urls := []*url.URL{
		{Scheme: "http", Host: "go-colly.org"},
		{Scheme: "http", Host: "example.com"},
		{Scheme: "http", Host: "xx.yy", Path: "/zz"},
	}

	// test cookie
	if s.Cookies(urls[0]) != "" {
		t.Fatal("unexpected result", s.Cookies(urls[0]))
	}
	cookies := "fake cookie"
	s.SetCookies(urls[0], cookies)
	if s.Cookies(urls[0]) != cookies {
		t.Fatal("unexpected result", cookies)
	}

	// test queue
	for _, u := range urls {
		if err := s.AddRequest([]byte(u.String())); err != nil {
			t.Fatal("failed to add request: " + err.Error())
		}
	}
	if size, err := s.QueueSize(); size != 3 || err != nil {
		t.Fatal("invalid queue size")
	}
	for _, u := range urls {
		if r, err := s.GetRequest(); err != nil || string(r) != u.String() {
			t.Fatal("failed to get request: ", err)
		}
	}
}
