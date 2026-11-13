package ledger

import "testing"

func TestLenEmpty(t *testing.T) {
	s := NewStore()
	if got := s.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}

func TestLenTwoKeys(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
}
