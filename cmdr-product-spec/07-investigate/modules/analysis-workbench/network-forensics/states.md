---
id: investigate-network-forensics-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Functional states

| Famille | États |
|---|---|
| Network Intake | draft, incomplete, ready, partial-capture, corrupted, truncated, restricted, unsupported, coverage-unknown, tool-unavailable, policy-blocked |
| Network Session | draft, ready, active, paused, blocked, partial, completed, failed, archived, superseded |
| Capture Review | unverified, verifying, verified, partially-verified, inconsistent, incomplete, corrupted, truncated, disputed, restricted |
| Coverage | complete-for-declared-scope, partial, gap-detected, loss-detected, truncation-detected, direction-ambiguous, timebase-uncertain, duplicated, unknown |
| Protocol | proposed, selected, confirmed-by-analyst, ambiguous, conflicting, unsupported, superseded |
| Reconstruction | queued, processing, partial, available, failed, missing-packets, truncated, incompatible, disputed, superseded |
| Network Anomaly | candidate, weakly-supported, supported, contradicted, inconclusive, coverage-limited, disputed, superseded, withdrawn |
| Derived Artifact | proposed, extracting, available, partial, invalid, restricted, superseded, withdrawn-from-use |
| Reproducibility | reproducible, partially-reproducible, not-reproducible, missing-capture, incomplete-capture, missing-sensor-context, missing-timebase, missing-tool, missing-version, restricted-payload, policy-blocked, disputed |

Ces états sont fonctionnels et non des machines d’état objet définitives.
