package ledger

import (
	"sort"
	"sync"
	"testing"
)

func TestNewStoreEmpty(t *testing.T) {
	s := NewStore()

	if v, ok := s.Get("x"); ok || v != 0 {
		t.Fatalf("Get on empty store = (%d, %t), want (0, false)", v, ok)
	}

	keys := s.Keys()
	if keys == nil {
		t.Fatal("Keys() = nil, want non-nil empty slice")
	}
	if len(keys) != 0 {
		t.Fatalf("Keys() len = %d, want 0", len(keys))
	}
}

func TestSetGetRoundTrip(t *testing.T) {
	s := NewStore()
	s.Set("a", 42)

	v, ok := s.Get("a")
	if !ok || v != 42 {
		t.Fatalf("Get(\"a\") = (%d, %t), want (42, true)", v, ok)
	}
}

func TestSetOverwrite(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("a", 2)

	v, ok := s.Get("a")
	if !ok || v != 2 {
		t.Fatalf("Get(\"a\") after overwrite = (%d, %t), want (2, true)", v, ok)
	}
}

func TestGetZeroValueDisambiguation(t *testing.T) {
	s := NewStore()
	s.Set("z", 0)

	if v, ok := s.Get("z"); !ok || v != 0 {
		t.Fatalf("Get(\"z\") = (%d, %t), want (0, true)", v, ok)
	}

	if v, ok := s.Get("missing"); ok || v != 0 {
		t.Fatalf("Get(\"missing\") = (%d, %t), want (0, false)", v, ok)
	}
}

func TestKeysSnapshot(t *testing.T) {
	s := NewStore()
	want := []string{"a", "b", "c"}
	for i, k := range want {
		s.Set(k, i)
	}

	got := s.Keys()
	sort.Strings(got)

	if len(got) != len(want) {
		t.Fatalf("Keys() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Keys() = %v, want %v", got, want)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := NewStore()
	const goroutines = 50
	const iterations = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(g int) {
			defer wg.Done()
			key := string(rune('a' + g%26))
			for i := 0; i < iterations; i++ {
				s.Set(key, i)
				s.Get(key)
				_ = s.Keys()
			}
		}(g)
	}
	wg.Wait()

	if v, ok := s.Get("a"); !ok || v != iterations-1 {
		t.Fatalf("Get(\"a\") after concurrent access = (%d, %t), want (%d, true)", v, ok, iterations-1)
	}
}
