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
		{
			name:   "empty store returns zero",
			values: nil,
			want:   0,
		},
		{
			name:   "positive values",
			values: map[string]int{"a": 1, "b": 2, "c": 3},
			want:   6,
		},
		{
			name:   "negative and zero values",
			values: map[string]int{"a": -5, "b": 0, "c": 10},
			want:   5,
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

// TestSumConcurrent exercises Sum alongside other read-locked accessors to
// confirm there is no data race under `go test -race`.
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
