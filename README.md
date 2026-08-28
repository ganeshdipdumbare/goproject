# goproject

A concurrency-safe in-memory key/value store.

## API

| Method | Description |
| --- | --- |
| `NewStore() *Store` | Returns an empty store ready for use. |
| `Get(key string) (int, bool)` | Returns the value stored under `key` and whether the key was present. |
| `Set(key string, val int)` | Stores `val` under `key`, overwriting any previous value. |
| `Remove(key string)` | Deletes `key` from the store. Removing a key that is not present is a no-op. |
| `Keys() []string` | Returns a snapshot of the keys currently in the store. |
