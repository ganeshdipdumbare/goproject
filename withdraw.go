package ledger

// Withdraw subtracts amount from the value stored under key and reports
// whether the key was present. A missing key is left unchanged. A stored
// zero is present, so the debit is applied.
func (s *Store) Withdraw(key string, amount int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.m[key]
	if !ok {
		return false
	}
	s.m[key] = v - amount
	return true
}
