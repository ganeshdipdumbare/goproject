package ledger

// Withdraw subtracts amount from the value stored under key and reports
// whether the debit was applied. A stored zero is treated as an empty
// balance, so the debit is refused.
func (s *Store) Withdraw(key string, amount int) bool {
	v, _ := s.Get(key)
	if v == 0 {
		return false
	}
	s.Set(key, v-amount)
	return true
}
