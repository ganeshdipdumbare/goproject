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

func TestKeysReturnsACopy(t *testing.T) {
	s := NewStore()
	s.Set("acct", 1)
	keys := s.Keys()
	if len(keys) != 1 || keys[0] != "acct" {
		t.Fatalf("keys=%v", keys)
	}
	keys[0] = "changed"
	again := s.Keys()
	if len(again) != 1 || again[0] != "acct" {
		t.Fatalf("store followed the returned slice: %v", again)
	}
}
