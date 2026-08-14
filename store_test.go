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
			values: map[string]int{"a": 42},
			want:   42,
		},
		{
			name:   "multiple values",
			values: map[string]int{"a": 1, "b": 2, "c": 3},
			want:   6,
		},
		{
			name:   "negative values",
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

func TestStoreSumConcurrent(t *testing.T) {
	s := NewStore()
	const n = 100

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Set(string(rune('a'+i%26))+string(rune('0'+i/26)), 1)
			_ = s.Sum()
		}(i)
	}
	wg.Wait()

	if got := s.Sum(); got != n {
		t.Errorf("Sum() = %d, want %d", got, n)
	}
}
