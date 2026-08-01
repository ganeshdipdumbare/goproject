package ledger

import (
	"strconv"
	"sync"
	"testing"
)

func TestResetClearsPopulatedStore(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)
	s.Set("c", 3)

	s.Reset()

	if keys := s.Keys(); len(keys) != 0 {
		t.Fatalf("expected no keys after Reset, got %v", keys)
	}
	if _, ok := s.Get("a"); ok {
		t.Fatalf("expected key %q to be absent after Reset", "a")
	}
}

func TestResetEmptyStoreIsUsable(t *testing.T) {
	s := NewStore()

	// Resetting an empty store must not panic.
	s.Reset()

	if keys := s.Keys(); len(keys) != 0 {
		t.Fatalf("expected no keys, got %v", keys)
	}

	// The store must remain usable afterwards.
	s.Set("x", 42)
	if v, ok := s.Get("x"); !ok || v != 42 {
		t.Fatalf("expected (42, true) after Set following Reset, got (%d, %v)", v, ok)
	}
}

func TestResetConcurrentAccess(t *testing.T) {
	s := NewStore()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func(n int) {
			defer wg.Done()
			s.Set(strconv.Itoa(n), n)
		}(i)
		go func() {
			defer wg.Done()
			s.Reset()
		}()
		go func() {
			defer wg.Done()
			_ = s.Keys()
		}()
	}
	wg.Wait()
}
