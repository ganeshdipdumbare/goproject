# goproject

A small concurrency-safe in-memory key/value store (`ledger.Store`).

## API

- `NewStore() *Store` — create an empty store.
- `Get(key string) (int, bool)` — read a value and whether the key was present.
- `Set(key string, val int)` — store a value, overwriting any previous one.
- `Remove(key string)` — delete a key; a no-op when the key is absent.
- `Keys() []string` — snapshot of the keys currently in the store.
