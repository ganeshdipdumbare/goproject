package ledger

// Credit adds amount to an existing key under one write lock.
// It reports false and changes nothing when amount is not positive, the key
// is missing, or the sum would overflow. A stored zero is present, so it
// can be credited.
func (s *Store) Credit(key string, amount int) bool {
	if amount <= 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.m[key]
	if !ok {
		return false
	}
	const maxInt = int(^uint(0) >> 1)
	if cur > maxInt-amount {
		return false
	}
	s.m[key] = cur + amount
	return true
}
