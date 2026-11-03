package ledger

import "testing"

func TestTransferMovesAmount(t *testing.T) {
	s := NewStore()
	s.Set("from", 10)
	s.Set("to", 3)
	if !s.Transfer("from", "to", 4) {
		t.Fatal("transfer")
	}
	from, fromOK := s.Get("from")
	to, toOK := s.Get("to")
	if !fromOK || from != 6 || !toOK || to != 7 {
		t.Fatalf("from=%d,%v to=%d,%v", from, fromOK, to, toOK)
	}
}

func TestTransferLeavesStoredZeroPresent(t *testing.T) {
	s := NewStore()
	s.Set("from", 4)
	s.Set("to", 0)
	if !s.Transfer("from", "to", 4) {
		t.Fatal("transfer onto stored zero")
	}
	from, fromOK := s.Get("from")
	to, toOK := s.Get("to")
	if !fromOK || from != 0 || !toOK || to != 4 {
		t.Fatalf("from=%d,%v to=%d,%v", from, fromOK, to, toOK)
	}
}

func TestTransferRefuses(t *testing.T) {
	s := NewStore()
	s.Set("from", 5)
	s.Set("to", 1)
	s.Set("zero", 0)

	cases := []struct {
		from, to string
		amount   int
	}{
		{"from", "from", 1},
		{"from", "to", 0},
		{"from", "to", -2},
		{"missing", "to", 1},
		{"from", "missing", 1},
		{"from", "to", 6},
		{"zero", "to", 1},
	}
	for _, tc := range cases {
		if s.Transfer(tc.from, tc.to, tc.amount) {
			t.Fatalf("accepted %s -> %s by %d", tc.from, tc.to, tc.amount)
		}
	}
	from, _ := s.Get("from")
	to, _ := s.Get("to")
	zero, zeroOK := s.Get("zero")
	if from != 5 || to != 1 || !zeroOK || zero != 0 {
		t.Fatalf("balances changed: from=%d to=%d zero=%d,%v", from, to, zero, zeroOK)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("missing key was created")
	}
}
