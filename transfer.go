package ledger

// Transfer moves amount from one key to another under a single write lock.
// It reports false and changes nothing when the keys are the same, amount is
// not positive, either key is missing, or the source value is smaller than
// amount. A stored zero is present, so it can receive a transfer.
func (s *Store) Transfer(from, to string, amount int) bool {
	if from == to || amount <= 0 {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	src, srcOK := s.m[from]
	dst, dstOK := s.m[to]
	if !srcOK || !dstOK || src < amount {
		return false
	}
	s.m[from] = src - amount
	s.m[to] = dst + amount
	return true
}
