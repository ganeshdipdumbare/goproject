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
