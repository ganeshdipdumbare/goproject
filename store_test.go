package ledger

import (
	"sync"
	"testing"
)

func TestStoreSum(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]int
		want    int
	}{
		{
			name:    "empty store",
			entries: nil,
			want:    0,
		},
		{
			name:    "single entry",
			entries: map[string]int{"a": 5},
			want:    5,
		},
		{
			name:    "multiple entries",
			entries: map[string]int{"a": 1, "b": 2, "c": 3},
			want:    6,
		},
		{
			name:    "mixed positive and negative",
			entries: map[string]int{"a": 10, "b": -4, "c": -6},
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for k, v := range tt.entries {
				s.Set(k, v)
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStoreSumConcurrentReads(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)
	s.Set("c", 3)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := s.Sum(); got != 6 {
				t.Errorf("Sum() = %d, want 6", got)
			}
			_ = s.Keys()
			_, _ = s.Get("a")
		}()
	}
	wg.Wait()
}
