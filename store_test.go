package ledger

import (
	"sync"
	"testing"
)

func TestStoreSum(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]int
		want   int
	}{
		{
			name:   "empty store",
			values: nil,
			want:   0,
		},
		{
			name:   "single value",
			values: map[string]int{"a": 5},
			want:   5,
		},
		{
			name:   "multiple values",
			values: map[string]int{"a": 1, "b": 2, "c": 3},
			want:   6,
		},
		{
			name:   "mixed signs",
			values: map[string]int{"a": 10, "b": -4, "c": -6},
			want:   0,
		},
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

// TestStoreSumConcurrent exercises Sum alongside other read accessors under the
// race detector to confirm the read lock prevents data races.
func TestStoreSumConcurrent(t *testing.T) {
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
