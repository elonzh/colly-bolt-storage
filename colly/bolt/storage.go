package bolt

import (
	"encoding/binary"
	"fmt"
	"net/url"

	bolt "go.etcd.io/bbolt"
)

var requestBucketName = []byte("request")
var cookieBucketName = []byte("cookie")
var queueBucketName = []byte("queue")

func uint64toByteArray(n uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, n)
	return bs
}

// Storage is a implementation for colly/queue and colly/storage
type Storage struct {
	db *bolt.DB
}

func NewStorage(db *bolt.DB) *Storage {
	return &Storage{db: db}
}

// Init initializes the storage
func (s *Storage) Init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucketName := range [][]byte{
			requestBucketName,
			cookieBucketName,
			queueBucketName,
		} {
			if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
				return err
			}
		}
		return nil
	})
}

// Visited receives and stores a request ID that is visited by the Collector{}
func (s *Storage) Visited(requestID uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		requestBucket := tx.Bucket(requestBucketName)
		return requestBucket.Put(uint64toByteArray(requestID), []byte{})
	})
}

// IsVisited returns true if the request was visited before IsVisited{}
// is called{}
func (s *Storage) IsVisited(requestID uint64) (bool, error) {
	var isVisited bool
	err := s.db.View(func(tx *bolt.Tx) error {
		requestBucket := tx.Bucket(requestBucketName)
		isVisited = requestBucket.Get(uint64toByteArray(requestID)) != nil
		return nil
	})
	return isVisited, err
}

// Cookies retrieves stored cookies for a given host{}
func (s *Storage) Cookies(u *url.URL) string {
	var cookies string
	err := s.db.View(func(tx *bolt.Tx) error {
		cookieBucket := tx.Bucket(cookieBucketName)
		cookies = string(cookieBucket.Get([]byte(u.String())))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return cookies
}

// SetCookies stores cookies for a given host{}
func (s *Storage) SetCookies(u *url.URL, cookies string) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		cookieBucket := tx.Bucket(cookieBucketName)
		return cookieBucket.Put([]byte(u.String()), []byte(cookies))
	})
	if err != nil {
		panic(err)
	}
}

// AddRequest adds a serialized request to the queue
func (s *Storage) AddRequest(request []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		queueBucket := tx.Bucket(queueBucketName)
		n, err := queueBucket.NextSequence()
		if err != nil {
			return err
		}
		key := uint64toByteArray(n)
		return queueBucket.Put(key, request)
	})
	return err
}

// GetRequest pops the next request from the queue
// or returns error if the queue is empty
func (s *Storage) GetRequest() ([]byte, error) {
	var request []byte
	err := s.db.Update(func(tx *bolt.Tx) error {
		queueBucket := tx.Bucket(queueBucketName)
		if queueBucket.Stats().KeyN == 0 {
			return fmt.Errorf("the queue is empty")
		}
		c := queueBucket.Cursor()
		_, request = c.First()
		return c.Delete()
	})
	return request, err
}

// QueueSize returns with the size of the queue
func (s *Storage) QueueSize() (int, error) {
	var queueSize int
	err := s.db.View(func(tx *bolt.Tx) error {
		queueSize = tx.Bucket(queueBucketName).Stats().KeyN
		return nil
	})
	return queueSize, err
}
