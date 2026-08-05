---
id: investigate-detection-engineering-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Functional states

| Family | States |
|---|---|
| Intake | draft, incomplete, ready, blocked, insufficient-evidence, data-source-required, permission-blocked, out-of-scope, superseded |
| Project | draft, ready, active, paused, blocked, partial, completed-authoring, failed, archived, superseded |
| Detection Content Draft | draft, incomplete, structurally-valid, valid-with-warnings, semantically-reviewed, test-ready, replay-ready, review-ready, blocked, invalid, superseded, withdrawn |
| Data readiness | unknown, unavailable, partially-available, available, stale, degraded, coverage-limited, retention-insufficient, unsupported, policy-blocked, disputed |
| Validation | not-run, queued, validating, valid, valid-with-warnings, invalid, blocked, partial, failed, superseded |
| Test scenario | draft, ready, restricted, incomplete, executed, partial, failed, superseded, withdrawn |
| Replay | draft, queued, running, partial, completed, failed, cancelled, retention-limited, data-gap, permission-blocked, superseded |
| Match review | unreviewed, candidate-true-positive, candidate-false-positive, candidate-false-negative, expected-non-match, ambiguous, insufficient-context, disputed, superseded |
| Coverage | unknown, no-coverage, planned, partial, conditionally-covered, test-covered, historically-observed, gap-identified, disputed, superseded |

These are functional states, not final object state machines. No Detection Content state is `deployed`, `active` or `production`.
