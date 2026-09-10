package ledger

import (
	"fmt"
	"sync"
	"testing"
)

func TestStore_Sum(t *testing.T) {
	tests := []struct {
		name   string
		values []struct {
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
			name: "single value",
			values: []struct {
				key string
				val int
			}{{"a", 5}},
			want: 5,
		},
		{
			name: "multiple values",
			values: []struct {
				key string
				val int
			}{{"a", 1}, {"b", 2}, {"c", 3}},
			want: 6,
		},
		{
			name: "negative and zero values",
			values: []struct {
				key string
				val int
			}{{"a", -4}, {"b", 0}, {"c", 10}},
			want: 6,
		},
		{
			name: "overwritten key counts once",
			values: []struct {
				key string
				val int
			}{{"a", 1}, {"a", 7}, {"b", 2}},
			want: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for _, v := range tt.values {
				s.Set(v.key, v.val)
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestStore_SumConcurrent(t *testing.T) {
	const writers = 50

	s := NewStore()
	var wg sync.WaitGroup

	want := 0
	for i := 1; i <= writers; i++ {
		want += i
	}

	for i := 1; i <= writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Set(fmt.Sprintf("key-%d", i), i)
		}(i)
	}

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = s.Sum()
			_, _ = s.Get(fmt.Sprintf("key-%d", i))
			_ = s.Keys()
		}(i)
	}

	wg.Wait()

	if got := s.Sum(); got != want {
		t.Errorf("Sum() = %d, want %d", got, want)
	}
}
