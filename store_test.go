package ledger

import (
	"fmt"
	"sync"
	"testing"
)

func TestStore_Sum(t *testing.T) {
	tests := []struct {
		name  string
		setup func(s *Store)
		want  int
	}{
		{
			name:  "empty store",
			setup: func(s *Store) {},
			want:  0,
		},
		{
			name:  "single entry",
			setup: func(s *Store) { s.Set("a", 42) },
			want:  42,
		},
		{
			name: "multiple entries with negatives",
			setup: func(s *Store) {
				s.Set("a", 5)
				s.Set("b", -2)
				s.Set("c", 10)
			},
			want: 13,
		},
		{
			name: "overwritten value",
			setup: func(s *Store) {
				s.Set("a", 1)
				s.Set("b", 2)
				s.Set("a", 10)
			},
			want: 12,
		},
		{
			name: "zero values do not affect the total",
			setup: func(s *Store) {
				s.Set("a", 0)
				s.Set("b", 7)
				s.Set("c", 0)
			},
			want: 7,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			tt.setup(s)
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStore_Sum_Concurrent(t *testing.T) {
	const n = 50

	s := NewStore()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Set(fmt.Sprintf("key-%d", i), i)
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			// The interleaved value is not deterministic; assert only that
			// the read accessors return without deadlocking or racing.
			_ = s.Sum()
			_, _ = s.Get("key-0")
			_ = s.Keys()
		}()
	}
	wg.Wait()

	want := n * (n - 1) / 2
	if got := s.Sum(); got != want {
		t.Errorf("Sum() after concurrent writes = %d, want %d", got, want)
	}
}
