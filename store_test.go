package ledger

import (
	"strconv"
	"sync"
	"testing"
)

func TestResetPopulatedStore(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	s.Reset()

	if got := len(s.Keys()); got != 0 {
		t.Fatalf("expected 0 keys after Reset, got %d", got)
	}
	if _, ok := s.Get("a"); ok {
		t.Fatalf("expected key %q to be absent after Reset", "a")
	}
}

func TestResetEmptyStore(t *testing.T) {
	s := NewStore()

	s.Reset()

	if got := len(s.Keys()); got != 0 {
		t.Fatalf("expected empty store to remain empty, got %d keys", got)
	}

	// A store must remain usable after being reset.
	s.Set("a", 1)
	if v, ok := s.Get("a"); !ok || v != 1 {
		t.Fatalf("expected (1, true) after Set following Reset, got (%d, %v)", v, ok)
	}
}

func TestResetConcurrent(t *testing.T) {
	s := NewStore()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := strconv.Itoa(n)
				s.Set(key, j)
				s.Get(key)
				s.Keys()
				s.Reset()
			}
		}(i)
	}
	wg.Wait()
}
