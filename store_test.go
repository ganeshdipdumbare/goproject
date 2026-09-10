package goproject_test

import (
	"sort"
	"sync"
	"testing"

	"goproject"
)

func TestStoreRemove(t *testing.T) {
	tests := []struct {
		name     string
		seed     map[string]int
		remove   string
		wantKeys []string
	}{
		{
			name:     "removes an existing key",
			seed:     map[string]int{"a": 1, "b": 2},
			remove:   "a",
			wantKeys: []string{"b"},
		},
		{
			name:     "missing key is a no-op",
			seed:     map[string]int{"a": 1, "b": 2},
			remove:   "c",
			wantKeys: []string{"a", "b"},
		},
		{
			name:     "empty store is a no-op",
			seed:     nil,
			remove:   "a",
			wantKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := goproject.NewStore()
			for k, v := range tt.seed {
				s.Set(k, v)
			}

			s.Remove(tt.remove)

			if _, ok := s.Get(tt.remove); ok {
				t.Errorf("Get(%q) after Remove: got ok = true, want false", tt.remove)
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

func TestStoreRemoveThenSet(t *testing.T) {
	s := goproject.NewStore()
	s.Set("a", 1)
	s.Remove("a")
	s.Set("a", 2)

	v, ok := s.Get("a")
	if !ok || v != 2 {
		t.Fatalf("Get(\"a\") = (%d, %t), want (2, true)", v, ok)
	}
	if got := len(s.Keys()); got != 1 {
		t.Fatalf("len(Keys()) = %d, want 1", got)
	}
}

func TestStoreRemoveConcurrent(t *testing.T) {
	s := goproject.NewStore()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			s.Set("a", 1)
		}()
		go func() {
			defer wg.Done()
			s.Remove("a")
		}()
		go func() {
			defer wg.Done()
			s.Get("a")
		}()
	}
	wg.Wait()
}
