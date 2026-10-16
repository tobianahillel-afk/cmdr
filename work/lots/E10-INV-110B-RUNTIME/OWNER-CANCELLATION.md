# E10-INV-110B — cancellation before owner resolution

The runtime's three-phase authorization preflight first checks Case and all selected Timeline Entry permissions, then resolves permitted owner entries, and finally checks source permissions. Between the final Timeline Entry `allow` result and the first `ResolveEntry` call the request can be cancelled. A cancelled request must not make an unnecessary owner metadata read.

Reuse first: a single Go `ctx.Err()` check at the transition into each owner iteration, using the existing `context.Context` and interfaces. No new background goroutine, third-party library, policy engine, storage or paid service. This guard does not make remote ACL decisions atomic, and a context cancellation that occurs after the check can still race with an in-flight resolver; the resolver remains obligated to honor the passed context.

The new deterministic test exposes a context whose `Err()` turns to `context.Canceled` exactly after the Case + Timeline checks and before owner resolution. It asserts denial, zero owner reads and precisely the preliminary Case + Timeline calls. The test fails on the parent implementation, which invokes `ResolveEntry` without a second cancellation check.

B remains READY pending completed gate qualification and independently accessible review; no production owner, ACL provider or SLO is claimed. This note is not a CI result: attach only observed exact-SHA GitHub Actions evidence.
