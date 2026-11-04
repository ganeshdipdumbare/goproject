package ledger

import "testing"

func TestGetDistinguishesStoredZero(t *testing.T) {
	s := NewStore()
	s.Set("acct", 0)
	got, ok := s.Get("acct")
	if !ok || got != 0 {
		t.Fatalf("stored zero: got=%d ok=%v", got, ok)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("missing key reported present")
	}
}
