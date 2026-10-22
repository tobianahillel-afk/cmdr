---
id: reverse-debugger-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Functional states

| Work item | States |
|---|---|
| Reverse Intake | draft, incomplete, ready, unsupported, restricted, tool-unavailable, environment-unavailable, policy-blocked |
| Reverse Analysis Session | draft, ready, active, paused, blocked, partial, completed, failed, archived, superseded |
| Function Interpretation | detected, proposed, confirmed-by-analyst, ambiguous, conflicting, superseded, withdrawn |
| Decompilation Result | queued, generating, partial, available, failed, ambiguous, incompatible, superseded |
| Debugger Session | requested, validating, ready, opening, active, paused, stepping, running, stopping, closed, failed, crashed, timed-out, environment-unavailable, revoked |
| Breakpoint | proposed, enabled, disabled, reached, not-reached, invalid, unresolved, superseded, removed-from-session |
| Runtime Snapshot | capturing, available, partial, stale, invalid, restricted, superseded |
| Patch Hypothesis | draft, ready-for-experiment, applied-in-isolated-copy, observed, supported, contradicted, inconclusive, reverted, superseded, rejected |
| Reproducibility | reproducible, partially-reproducible, not-reproducible, missing-tool, missing-version, missing-input, missing-environment, state-not-restorable, behavior-not-reproduced, disputed |

These are functional states only; final object state machines remain deferred.
