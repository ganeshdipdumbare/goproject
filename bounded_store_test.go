package ledger

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestNewBoundedStore(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{name: "negative", size: -1, wantErr: true},
		{name: "zero", size: 0, wantErr: true},
		{name: "one", size: 1},
		{name: "three", size: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewBoundedStore(tt.size)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewBoundedStore(%d) error = nil, want error", tt.size)
				}
				if s != nil {
					t.Fatalf("NewBoundedStore(%d) store = %v, want nil", tt.size, s)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewBoundedStore(%d) unexpected error: %v", tt.size, err)
			}
			if s == nil {
				t.Fatalf("NewBoundedStore(%d) returned nil store", tt.size)
			}
			if got := s.Len(); got != 0 {
				t.Errorf("Len() = %d, want 0", got)
			}
			if got := s.Keys(); len(got) != 0 {
				t.Errorf("Keys() = %v, want empty", got)
			}
		})
	}
}

type boundedOp struct {
	get bool // false: Set(key, val); true: Get(key)
	key string
	val int
}

func opSet(key string, val int) boundedOp { return boundedOp{key: key, val: val} }
func opGet(key string) boundedOp          { return boundedOp{get: true, key: key} }

func TestBoundedStoreOperations(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		ops      []boundedOp
		wantKeys []string       // least to most recently used
		present  map[string]int // key -> expected value
		missing  []string
	}{
		{
			name:     "basic set and get",
			size:     3,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2)},
			wantKeys: []string{"a", "b"},
			present:  map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "overwrite without eviction",
			size:     2,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2), opSet("a", 10)},
			wantKeys: []string{"b", "a"},
			present:  map[string]int{"a": 10, "b": 2},
		},
		{
			name:     "eviction with no reads",
			size:     2,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2), opSet("c", 3)},
			wantKeys: []string{"b", "c"},
			present:  map[string]int{"b": 2, "c": 3},
			missing:  []string{"a"},
		},
		{
			name:     "get refreshes recency",
			size:     2,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2), opGet("a"), opSet("c", 3)},
			wantKeys: []string{"a", "c"},
			present:  map[string]int{"a": 1, "c": 3},
			missing:  []string{"b"},
		},
		{
			name:     "set on existing key refreshes recency",
			size:     3,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2), opSet("c", 3), opSet("a", 4), opSet("d", 5)},
			wantKeys: []string{"c", "a", "d"},
			present:  map[string]int{"a": 4, "c": 3, "d": 5},
			missing:  []string{"b"},
		},
		{
			name:     "size one",
			size:     1,
			ops:      []boundedOp{opSet("a", 1), opSet("b", 2)},
			wantKeys: []string{"b"},
			present:  map[string]int{"b": 2},
			missing:  []string{"a"},
		},
		{
			name:     "missing key versus stored zero",
			size:     2,
			ops:      []boundedOp{opSet("zero", 0), opGet("nope")},
			wantKeys: []string{"zero"},
			present:  map[string]int{"zero": 0},
			missing:  []string{"nope"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewBoundedStore(tt.size)
			if err != nil {
				t.Fatalf("NewBoundedStore(%d): %v", tt.size, err)
			}
			for _, op := range tt.ops {
				if op.get {
					s.Get(op.key)
				} else {
					s.Set(op.key, op.val)
				}
			}

			// Check order and length before any Get below changes recency.
			if got := s.Keys(); !reflect.DeepEqual(got, tt.wantKeys) {
				t.Errorf("Keys() = %v, want %v", got, tt.wantKeys)
			}
			// Keys must not change recency.
			if got := s.Keys(); !reflect.DeepEqual(got, tt.wantKeys) {
				t.Errorf("second Keys() = %v, want %v", got, tt.wantKeys)
			}
			if got := s.Len(); got != len(tt.wantKeys) {
				t.Errorf("Len() = %d, want %d", got, len(tt.wantKeys))
			}

			for k, want := range tt.present {
				got, ok := s.Get(k)
				if !ok || got != want {
					t.Errorf("Get(%q) = (%d, %v), want (%d, true)", k, got, ok, want)
				}
			}
			for _, k := range tt.missing {
				if got, ok := s.Get(k); ok || got != 0 {
					t.Errorf("Get(%q) = (%d, %v), want (0, false)", k, got, ok)
				}
			}
		})
	}
}

func TestBoundedStoreConcurrent(t *testing.T) {
	const (
		size       = 8
		goroutines = 16
		iterations = 500
	)
	s, err := NewBoundedStore(size)
	if err != nil {
		t.Fatalf("NewBoundedStore(%d): %v", size, err)
	}

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				key := fmt.Sprintf("k%d", (g+i)%(size*2))
				s.Set(key, i)
				s.Get(key)
				if n := s.Len(); n > size {
					t.Errorf("Len() = %d, exceeds size %d", n, size)
				}
				if n := len(s.Keys()); n > size {
					t.Errorf("len(Keys()) = %d, exceeds size %d", n, size)
				}
			}
		}(g)
	}
	wg.Wait()

	if n := s.Len(); n > size {
		t.Errorf("final Len() = %d, exceeds size %d", n, size)
	}
	if n := len(s.Keys()); n > size {
		t.Errorf("final len(Keys()) = %d, exceeds size %d", n, size)
	}
}
