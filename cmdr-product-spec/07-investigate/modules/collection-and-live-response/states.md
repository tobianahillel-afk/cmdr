# Functional states

Endpoint: available, degraded, offline, stale, unsupported, restricted, policy-blocked, unknown.
Request: draft, incomplete, ready, submitted, awaiting-approval, accepted, rejected, superseded, cancelled.
Job: queued, dispatched, running, partial, completed, failed, cancelled, expired, endpoint-offline, awaiting-reconnect.
Session: requested, awaiting-approval, opening, active, idle, suspended, reconnecting, closing, closed, expired, failed, revoked.
Operation: draft, validation-required, ready, queued, running, interrupted, partial, completed, failed, cancelled, denied, policy-blocked.
Transfer: preparing, transferring, verifying, completed, partial, failed, cancelled, collision, restricted.
Result review: received, incomplete, under-review, verified, disputed, superseded, linked-to-artifact, linked-to-evidence.

These are functional dimensions, not final object state machines.
