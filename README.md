# goproject

A concurrency-safe in-memory integer store (`package ledger`).

## Store API

- `NewStore() *Store` — returns an empty store ready for use.
- `Get(key string) (int, bool)` — returns the stored value and whether the key was present.
- `Set(key string, val int)` — stores a value, overwriting any previous one.
- `Keys() []string` — returns a snapshot of the keys currently in the store.
- `Sum() int` — returns the sum of all values currently in the store.
