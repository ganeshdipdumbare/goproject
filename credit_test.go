package ledger

import "testing"

func TestCreditAddsToPresentKey(t *testing.T) {
	s := NewStore()
	s.Set("acct", 10)
	if !s.Credit("acct", 4) {
		t.Fatal("credit")
	}
	got, ok := s.Get("acct")
	if !ok || got != 14 {
		t.Fatalf("acct=%d,%v", got, ok)
	}
}

func TestCreditLeavesStoredZeroPresent(t *testing.T) {
	s := NewStore()
	s.Set("acct", 0)
	if !s.Credit("acct", 3) {
		t.Fatal("credit onto stored zero")
	}
	got, ok := s.Get("acct")
	if !ok || got != 3 {
		t.Fatalf("acct=%d,%v", got, ok)
	}
}

func TestCreditRefuses(t *testing.T) {
	s := NewStore()
	s.Set("acct", 5)
	s.Set("zero", 0)
	s.Set("full", int(^uint(0)>>1))

	cases := []struct {
		key    string
		amount int
	}{
		{"acct", 0},
		{"acct", -2},
		{"missing", 1},
		{"full", 1},
	}
	for _, tc := range cases {
		if s.Credit(tc.key, tc.amount) {
			t.Fatalf("accepted %s by %d", tc.key, tc.amount)
		}
	}
	acct, acctOK := s.Get("acct")
	zero, zeroOK := s.Get("zero")
	full, fullOK := s.Get("full")
	if !acctOK || acct != 5 || !zeroOK || zero != 0 || !fullOK || full != int(^uint(0)>>1) {
		t.Fatalf("balances changed: acct=%d,%v zero=%d,%v full=%d,%v", acct, acctOK, zero, zeroOK, full, fullOK)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("missing key was created")
	}
}
