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
	s.Set("c", 3)

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
		t.Fatalf("expected 0 keys after Reset on empty store, got %d", got)
	}

	s.Set("x", 42)
	if v, ok := s.Get("x"); !ok || v != 42 {
		t.Fatalf("expected store reusable after Reset, got (%d, %t)", v, ok)
	}
}

func TestResetConcurrent(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s.Set(strconv.Itoa(i)+"-"+strconv.Itoa(j), j)
			}
		}(i)
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s.Reset()
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = s.Keys()
			}
		}()
	}

	wg.Wait()
}
