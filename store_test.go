package ledger

import (
	"sync"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int
		want   int
	}{
		{name: "empty store", values: nil, want: 0},
		{name: "single value", values: map[string]int{"a": 7}, want: 7},
		{name: "multiple values", values: map[string]int{"a": 1, "b": 2, "c": 3}, want: 6},
		{name: "negative values", values: map[string]int{"a": 10, "b": -4, "c": -6}, want: 0},
		{name: "mixed values", values: map[string]int{"a": -5, "b": 20, "c": -3}, want: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for k, v := range tt.values {
				s.Set(k, v)
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestSumConcurrent exercises Sum alongside other readers to confirm the read
// lock is used correctly. Run with -race to detect data races.
func TestSumConcurrent(t *testing.T) {
	s := NewStore()
	for i := 0; i < 100; i++ {
		s.Set(string(rune('a'+i%26)), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Sum()
			_ = s.Keys()
			_, _ = s.Get("a")
		}()
	}
	wg.Wait()
}
