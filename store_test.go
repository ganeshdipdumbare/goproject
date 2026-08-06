package ledger

import (
	"sync"
	"testing"
)

func TestSumEmpty(t *testing.T) {
	s := NewStore()
	if got := s.Sum(); got != 0 {
		t.Fatalf("Sum() on empty store = %d, want 0", got)
	}
}

func TestSumPopulated(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)
	s.Set("c", 3)
	s.Set("neg", -10)

	if got, want := s.Sum(), -4; got != want {
		t.Fatalf("Sum() = %d, want %d", got, want)
	}
}

func TestSumConcurrent(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() { defer wg.Done(); _ = s.Sum() }()
		go func() { defer wg.Done(); _ = s.Keys() }()
		go func(n int) { defer wg.Done(); s.Set("a", n) }(i)
	}
	wg.Wait()
}
