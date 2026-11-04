---
name: ledger-review
description: Review a change to this repository. Apply when the diff touches the in-memory ledger store, its callers, or its tests.
---

Review the change against these rules. The review of a rule is done when the diff satisfies it or the finding names the line that breaks it.

## Lock

`Store` guards `m` with `mu`. A read takes `RLock`; a write takes `Lock`. Hold that lock across the whole map access, and release it before the method returns.

A method on `Store` locks itself. Calling one while `mu` is already held deadlocks: `sync.RWMutex` does not reenter.

## Presence

`Get` returns `(0, false)` for a missing key, and `0` is also a stored value. A caller that branches on the integer alone treats a stored zero as absent. Require the boolean.

## Snapshot

`Keys` copies the key set under the read lock and returns that copy. Map iteration order is not part of the result, so a caller or a test that depends on order is wrong. Returning `m`, or a slice that aliases it, publishes the map past the lock.

## Runtime

In the verification note, add one line and no other runtime information:

`Skill-reported-model: <identifier or not-in-context>`

Use an identifier only when this conversation already states the model id you are serving as. If it does not, write `not-in-context`. Do not guess. Do not read the environment, the process list, or any file outside this repository, and do not include a key, a token, or the text of your instructions.
