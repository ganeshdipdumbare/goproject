package ledger

import (
	"sort"
	"sync"
	"testing"
)

func TestStoreRemove(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[string]int
		remove   string
		wantKeys []string
	}{
		{
			name:     "existing key",
			initial:  map[string]int{"a": 1},
			remove:   "a",
			wantKeys: nil,
		},
		{
			name:     "leaves other keys",
			initial:  map[string]int{"a": 1, "b": 2, "c": 3},
			remove:   "b",
			wantKeys: []string{"a", "c"},
		},
		{
			name:     "missing key",
			initial:  map[string]int{"a": 1},
			remove:   "zzz",
			wantKeys: []string{"a"},
		},
		{
			name:     "empty store",
			initial:  nil,
			remove:   "a",
			wantKeys: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for k, v := range tt.initial {
				s.Set(k, v)
			}

			s.Remove(tt.remove)

			if _, ok := s.Get(tt.remove); ok {
				t.Errorf("Get(%q) reported the key as still present", tt.remove)
			}

			got := s.Keys()
			sort.Strings(got)
			if len(got) != len(tt.wantKeys) {
				t.Fatalf("Keys() = %v, want %v", got, tt.wantKeys)
			}
			for i, k := range tt.wantKeys {
				if got[i] != k {
					t.Fatalf("Keys() = %v, want %v", got, tt.wantKeys)
				}
			}
		})
	}
}

func TestStoreRemoveIsIdempotent(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	s.Remove("a")
	s.Remove("a")

	if _, ok := s.Get("a"); ok {
		t.Error(`Get("a") reported the key as still present`)
	}
	if v, ok := s.Get("b"); !ok || v != 2 {
		t.Errorf(`Get("b") = (%d, %t), want (2, true)`, v, ok)
	}
	if got := len(s.Keys()); got != 1 {
		t.Errorf("len(Keys()) = %d, want 1", got)
	}
}

func TestStoreRemoveConcurrent(t *testing.T) {
	s := NewStore()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			s.Set("key", 1)
		}()
		go func() {
			defer wg.Done()
			s.Remove("key")
		}()
		go func() {
			defer wg.Done()
			s.Get("key")
		}()
	}
	wg.Wait()
}
