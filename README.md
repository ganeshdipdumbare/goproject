# goproject

A small concurrency-safe in-memory key/value store.

## Store

`Store` is a concurrency-safe in-memory integer store; all access to the
underlying map is guarded by a mutex, so it is safe for use by multiple
goroutines.

| Method | Description |
| --- | --- |
| `NewStore() *Store` | Returns an empty store ready for use. |
| `Get(key string) (int, bool)` | Returns the value stored under `key` and whether it was present. |
| `Set(key string, val int)` | Stores `val` under `key`, overwriting any previous value. |
| `Remove(key string) bool` | Deletes `key` from the store. No-op when the key is absent; returns whether a key was removed. |
| `Keys() []string` | Returns a snapshot of the keys currently in the store. |
