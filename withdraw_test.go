package ledger

import "testing"

func TestWithdrawDebitsStoredValue(t *testing.T) {
	s := NewStore()
	s.Set("acct", 10)
	if !s.Withdraw("acct", 4) {
		t.Fatal("withdraw")
	}
	got, ok := s.Get("acct")
	if !ok || got != 6 {
		t.Fatalf("balance = %d, %v", got, ok)
	}
}

func TestWithdrawAppliesToStoredZero(t *testing.T) {
	s := NewStore()
	s.Set("acct", 0)
	if !s.Withdraw("acct", 1) {
		t.Fatal("stored zero is present")
	}
	got, ok := s.Get("acct")
	if !ok || got != -1 {
		t.Fatalf("balance = %d, %v", got, ok)
	}
}

func TestWithdrawMissingKey(t *testing.T) {
	s := NewStore()
	if s.Withdraw("acct", 1) {
		t.Fatal("missing key")
	}
	if _, ok := s.Get("acct"); ok {
		t.Fatal("missing key was created")
	}
}
