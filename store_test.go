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
			entries: map[string]int{"a": 42},
			want:    42,
		},
		{
			name:    "multiple entries with negatives",
			entries: map[string]int{"a": 5, "b": -2, "c": 10},
			want:    13,
		},
		{
			name:    "zero values",
			entries: map[string]int{"a": 0, "b": 0},
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
				t.Fatalf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStoreSumAfterMutation(t *testing.T) {
	s := NewStore()
	s.Set("a", 5)
	s.Set("b", 7)
	if got := s.Sum(); got != 12 {
		t.Fatalf("Sum() = %d, want 12", got)
	}

	s.Set("a", 1)
	if got := s.Sum(); got != 8 {
		t.Fatalf("Sum() after overwrite = %d, want 8", got)
	}
}

func TestStoreSumConcurrent(t *testing.T) {
	s := NewStore()
	for i := 0; i < 10; i++ {
		s.Set(string(rune('a'+i)), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = s.Sum()
				_, _ = s.Get("a")
				_ = s.Keys()
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			s.Set("writer", j)
		}
	}()

	wg.Wait()
}
