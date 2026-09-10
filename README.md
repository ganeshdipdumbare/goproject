# goproject

## ledger.Store

A concurrency-safe in-memory integer store.

- `Get(key string) (int, bool)` — read a value and whether the key exists.
- `Set(key string, val int)` — store a value, overwriting any previous one.
- `Keys() []string` — snapshot of the keys currently in the store.
- `Sum() int` — sum of all values currently in the store.
