package ledger

import (
	"sync"
	"testing"
)

func TestRemoveExistingKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)

	s.Remove("a")

	if _, ok := s.Get("a"); ok {
		t.Fatalf("expected key %q to be absent after Remove", "a")
	}
	if got := len(s.Keys()); got != 0 {
		t.Fatalf("expected empty store after Remove, got %d keys", got)
	}
}

func TestRemoveNonExistentKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)

	// Removing a key that was never stored must be a safe no-op.
	s.Remove("missing")

	if v, ok := s.Get("a"); !ok || v != 1 {
		t.Fatalf("expected untouched key a=1, got (%d, %t)", v, ok)
	}
	if got := len(s.Keys()); got != 1 {
		t.Fatalf("expected 1 key remaining, got %d", got)
	}
}

func TestRemoveFromEmptyStore(t *testing.T) {
	s := NewStore()

	// Must not panic on an empty store.
	s.Remove("a")

	if got := len(s.Keys()); got != 0 {
		t.Fatalf("expected empty store, got %d keys", got)
	}
}

func TestRemoveConcurrent(t *testing.T) {
	s := NewStore()
	const n = 100

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Set("k", 1)
			s.Remove("k")
			s.Get("k")
		}()
	}
	wg.Wait()
}
