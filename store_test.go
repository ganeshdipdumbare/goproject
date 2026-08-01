package ledger

import "testing"

func TestRemoveExistingKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	s.Remove("a")

	if _, ok := s.Get("a"); ok {
		t.Fatalf("expected key %q to be removed", "a")
	}
	if got := len(s.Keys()); got != 1 {
		t.Fatalf("expected 1 key remaining, got %d", got)
	}
	if _, ok := s.Get("b"); !ok {
		t.Fatalf("expected key %q to still be present", "b")
	}
}

func TestRemoveNonExistentKey(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)

	s.Remove("missing")

	if _, ok := s.Get("a"); !ok {
		t.Fatalf("expected key %q to still be present", "a")
	}
	if got := len(s.Keys()); got != 1 {
		t.Fatalf("expected 1 key remaining, got %d", got)
	}
}

func TestRemoveFromEmptyStore(t *testing.T) {
	s := NewStore()

	s.Remove("anything")

	if got := len(s.Keys()); got != 0 {
		t.Fatalf("expected store to remain empty, got %d keys", got)
	}
}
