package ledger

import (
	"strconv"
	"sync"
	"testing"
)

func TestStoreSum(t *testing.T) {
	tests := []struct {
		name    string
		entries []struct {
			key string
			val int
		}
		want int
	}{
		{
			name: "empty store",
			want: 0,
		},
		{
			name: "single entry",
			entries: []struct {
				key string
				val int
			}{
				{"a", 42},
			},
			want: 42,
		},
		{
			name: "multiple entries including negatives",
			entries: []struct {
				key string
				val int
			}{
				{"a", 5},
				{"b", -2},
				{"c", 10},
			},
			want: 13,
		},
		{
			name: "overwritten key counts once",
			entries: []struct {
				key string
				val int
			}{
				{"a", 5},
				{"b", 3},
				{"a", 1},
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for _, e := range tt.entries {
				s.Set(e.key, e.val)
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStoreSumConcurrentAccess(t *testing.T) {
	s := NewStore()
	for i := 0; i < 10; i++ {
		s.Set(strconv.Itoa(i), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = s.Sum()
				_, _ = s.Get("1")
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
