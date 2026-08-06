package ledger

import (
	"sync"
	"testing"
)

func TestRemove(t *testing.T) {
	tests := []struct {
		name    string
		seed    map[string]int
		remove  string
		wantLen int
	}{
		{
			name:    "existing key",
			seed:    map[string]int{"a": 1, "b": 2},
			remove:  "a",
			wantLen: 1,
		},
		{
			name:    "absent key",
			seed:    map[string]int{"a": 1},
			remove:  "missing",
			wantLen: 1,
		},
		{
			name:    "empty store",
			seed:    nil,
			remove:  "a",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for k, v := range tt.seed {
				s.Set(k, v)
			}

			s.Remove(tt.remove)

			if _, ok := s.Get(tt.remove); ok {
				t.Fatalf("Get(%q) still present after Remove", tt.remove)
			}
			if got := len(s.Keys()); got != tt.wantLen {
				t.Fatalf("len after Remove = %d, want %d", got, tt.wantLen)
			}
		})
	}
}

func TestRemoveConcurrent(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		key := "k"
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.Set(key, 1)
		}()
		go func() {
			defer wg.Done()
			s.Remove(key)
		}()
	}
	wg.Wait()
}
