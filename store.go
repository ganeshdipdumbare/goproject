package ledger

import "sync"

// Store is a concurrency-safe in-memory integer store. All access to the
// underlying map goes through the mutex, so it is safe for concurrent use by
// multiple goroutines.
type Store struct {
	mu sync.RWMutex
	m  map[string]int
}

// NewStore returns an empty Store ready for use.
func NewStore() *Store {
	return &Store{m: make(map[string]int)}
}

// Get returns the value stored under key and whether the key was present.
// A missing key returns (0, false); callers must check the boolean before
// trusting the returned value, since 0 is also a legitimate stored value.
func (s *Store) Get(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

// Set stores val under key, overwriting any previous value.
func (s *Store) Set(key string, val int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = val
}

// Remove deletes key from the store. Removing a key that is not present is a
// safe no-op, mirroring the behaviour of Go's built-in delete.
func (s *Store) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

// Keys returns a snapshot of the keys currently in the store.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.m))
	for k := range s.m {
		out = append(out, k)
	}
	return out
}
