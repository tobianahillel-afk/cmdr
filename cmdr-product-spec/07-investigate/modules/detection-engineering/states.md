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
| Validation / Replay / Match / Coverage | retain CAP-INV-411..416 functional states |
| Release Candidate | draft, incomplete, review-requested, under-review, changes-requested, review-passed, review-failed, readiness-pending, blocked, superseded, withdrawn |
| Deployment Readiness | not-assessed, assessing, ready, ready-with-conditions, incomplete, incompatible, dependency-missing, target-unavailable, approval-required, blocked, disputed |
| Change Request | draft, incomplete, ready, submitted, awaiting-decision, approved, rejected, returned-for-information, cancelled, superseded |
| Promotion/Deployment projection | planned, scheduled, awaiting-approval, approved, deploying, partially-deployed, deployed, activating, active, activation-failed, deactivating, inactive, rollback-requested, rolling-back, rolled-back, failed, cancelled, superseded |
| Health | unknown, healthy, degraded, delayed, partially-executing, failing, source-unavailable, schema-incompatible, dependency-failed, disabled, suspended, stale, disputed |
| Runtime Quality | unknown, insufficient-data, under-review, acceptable, noisy, low-value, coverage-limited, candidate-regression, candidate-improvement, disputed |
| Drift | no-drift-observed, drift-candidate, confirmed-change, compatibility-unknown, compatible, partially-compatible, incompatible, remediation-required, blocked, disputed |
| Retirement | proposed, under-review, approved-projection, scheduled, deactivating, inactive, retired, replacement-pending, completed, rollback-required, cancelled, superseded |

These are functional projections, not final object state machines. Govern and runtime owners remain the source of actual authority and state.
