package ledger

import (
	"sync"
	"testing"
)

func TestRemoveExistingKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)

	if removed := s.Remove("a"); !removed {
		t.Fatalf("Remove(\"a\") = false, want true for a present key")
	}
	if _, ok := s.Get("a"); ok {
		t.Fatalf("Get(\"a\") still present after Remove")
	}
	if got := len(s.Keys()); got != 0 {
		t.Fatalf("len(Keys()) = %d, want 0 after removing the only key", got)
	}
}

func TestRemoveMissingKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)

	if removed := s.Remove("b"); removed {
		t.Fatalf("Remove(\"b\") = true, want false for an absent key")
	}
	if v, ok := s.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(\"a\") = (%d, %t), want (1, true); Remove must not disturb other keys", v, ok)
	}
	if got := len(s.Keys()); got != 1 {
		t.Fatalf("len(Keys()) = %d, want 1 after removing an absent key", got)
	}
}

func TestRemoveFromEmptyStore(t *testing.T) {
	s := NewStore()

	if removed := s.Remove("a"); removed {
		t.Fatalf("Remove(\"a\") = true, want false on an empty store")
	}
	if got := len(s.Keys()); got != 0 {
		t.Fatalf("len(Keys()) = %d, want 0 for an empty store", got)
	}
}

func TestRemoveConcurrent(t *testing.T) {
	s := NewStore()
	const workers = 8

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				s.Set("k", j)
				s.Remove("k")
				s.Get("k")
				_ = s.Keys()
			}
		}()
	}
	wg.Wait()
}
