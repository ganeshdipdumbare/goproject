# goproject

## Store API

- `NewStore()` — create an empty store.
- `Get(key)` — read a value and whether the key was present.
- `Set(key, val)` — store a value under a key.
- `Remove(key)` — delete a key; removing an absent key is a no-op.
- `Keys()` — snapshot of the keys currently stored.
