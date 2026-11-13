package ledger

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

// BoundedStore is a concurrency-safe in-memory integer store with a fixed
// capacity. When it is full, setting a new key evicts the least recently used
// key. All access goes through the underlying lru.Cache, which holds its own
// lock, so it is safe for concurrent use by multiple goroutines.
type BoundedStore struct {
	c *lru.Cache[string, int]
}

// NewBoundedStore returns an empty BoundedStore that holds at most size keys.
// It returns an error if size is not positive.
func NewBoundedStore(size int) (*BoundedStore, error) {
	if size <= 0 {
		return nil, fmt.Errorf("ledger: bounded store size must be positive, got %d", size)
	}
	c, err := lru.New[string, int](size)
	if err != nil {
		return nil, err
	}
	return &BoundedStore{c: c}, nil
}

// Get returns the value stored under key and whether the key was present.
// A missing key returns (0, false); callers must check the boolean before
// trusting the returned value, since 0 is also a legitimate stored value.
// A successful Get marks key as the most recently used.
func (s *BoundedStore) Get(key string) (int, bool) {
	return s.c.Get(key)
}

// Set stores val under key, overwriting any previous value, and marks key as
// the most recently used. If key is new and the store is full, the least
// recently used key is evicted.
func (s *BoundedStore) Set(key string, val int) {
	s.c.Add(key, val)
}

// Keys returns a snapshot of the keys currently in the store, ordered from
// least to most recently used. It does not change recency.
func (s *BoundedStore) Keys() []string {
	return s.c.Keys()
}

// Len returns the number of keys currently in the store.
func (s *BoundedStore) Len() int {
	return s.c.Len()
}
