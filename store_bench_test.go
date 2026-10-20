package ledger

import (
	"strconv"
	"testing"
)

// newPopulatedStore returns a Store holding n keys and the slice of those keys,
// so benchmarks can build their fixtures outside the timed loop.
func newPopulatedStore(n int) (*Store, []string) {
	s := NewStore()
	keys := makeKeys(n)
	for i, k := range keys {
		s.Set(k, i)
	}
	return s, keys
}

// makeKeys pre-generates n distinct keys.
func makeKeys(n int) []string {
	keys := make([]string, n)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}
	return keys
}

// Package-level sinks keep the compiler from optimising away benchmarked calls.
var (
	sinkStore *Store
	sinkKeys  []string
)

func BenchmarkNewStore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkStore = NewStore()
	}
}

func BenchmarkSet(b *testing.B) {
	s := NewStore()
	keys := makeKeys(1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(keys[i%len(keys)], i)
	}
}

func BenchmarkGet(b *testing.B) {
	s, keys := newPopulatedStore(1024)

	b.Run("hit", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = s.Get(keys[i%len(keys)])
		}
	})

	b.Run("miss", func(b *testing.B) {
		missing := "missing-key"
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = s.Get(missing)
		}
	})
}

func BenchmarkKeys(b *testing.B) {
	for _, size := range []int{10, 1000, 100000} {
		b.Run("size="+strconv.Itoa(size), func(b *testing.B) {
			s, _ := newPopulatedStore(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sinkKeys = s.Keys()
			}
		})
	}
}

func BenchmarkGetParallel(b *testing.B) {
	s, keys := newPopulatedStore(1024)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = s.Get(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkSetParallel(b *testing.B) {
	s := NewStore()
	keys := makeKeys(1024)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Set(keys[i%len(keys)], i)
			i++
		}
	})
}

// BenchmarkMixedParallel models a read-heavy workload: 90% Get, 10% Set.
func BenchmarkMixedParallel(b *testing.B) {
	s, keys := newPopulatedStore(1024)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			k := keys[i%len(keys)]
			if i%10 == 0 {
				s.Set(k, i)
			} else {
				_, _ = s.Get(k)
			}
			i++
		}
	})
}
