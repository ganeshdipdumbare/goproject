package ledger

import (
	"sync"
	"testing"
)

func TestSum(t *testing.T) {
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
			name: "multiple values",
			setup: func(s *Store) {
				s.Set("a", 1)
				s.Set("b", 2)
				s.Set("c", 3)
			},
			want: 6,
		},
		{
			name: "overwritten key",
			setup: func(s *Store) {
				s.Set("a", 1)
				s.Set("a", 5)
			},
			want: 5,
		},
		{
			name: "negative and zero values",
			setup: func(s *Store) {
				s.Set("a", -2)
				s.Set("b", 2)
				s.Set("c", 0)
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			tt.setup(s)
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	s := NewStore()
	if v, ok := s.Get("missing"); ok || v != 0 {
		t.Errorf("Get(missing) = (%d, %t), want (0, false)", v, ok)
	}
	s.Set("a", 7)
	if v, ok := s.Get("a"); !ok || v != 7 {
		t.Errorf("Get(a) = (%d, %t), want (7, true)", v, ok)
	}
}

func TestKeys(t *testing.T) {
	s := NewStore()
	if got := s.Keys(); len(got) != 0 {
		t.Errorf("Keys() = %v, want empty", got)
	}
	s.Set("a", 1)
	s.Set("b", 2)
	got := s.Keys()
	if len(got) != 2 {
		t.Errorf("Keys() len = %d, want 2", len(got))
	}
	seen := map[string]bool{}
	for _, k := range got {
		seen[k] = true
	}
	if !seen["a"] || !seen["b"] {
		t.Errorf("Keys() = %v, want to contain a and b", got)
	}
}

func TestSumConcurrentReads(t *testing.T) {
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
