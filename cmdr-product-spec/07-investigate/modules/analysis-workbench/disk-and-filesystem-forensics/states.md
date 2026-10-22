---
id: investigate-disk-filesystem-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Functional states

| Area | States |
|---|---|
| Disk Intake | draft, incomplete, ready, partial-image, corrupted, encrypted, restricted, unsupported, filesystem-required, tool-unavailable, policy-blocked |
| Disk Session | draft, ready, active, paused, blocked, partial, completed, failed, archived, superseded |
| Integrity Review | unverified, verifying, verified, partially-verified, inconsistent, incomplete, corrupted, disputed, restricted |
| Filesystem Identification | proposed, selected, confirmed-by-analyst, ambiguous, conflicting, unsupported, superseded |
| File Entry | active, deleted-candidate, inaccessible, partial, inconsistent, recovered-candidate, superseded |
| Recovery Result | proposed, processing, available, partial, invalid, unattributed, restricted, superseded, withdrawn-from-use |
| Persistence Candidate | candidate, weakly-supported, supported, contradicted, inconclusive, disputed, superseded, withdrawn |
| Reproducibility | reproducible, partially-reproducible, not-reproducible, missing-image, incomplete-image, missing-volume, missing-filesystem-support, missing-tool, missing-version, restricted-content, policy-blocked, disputed |

These are functional states, not final object state machines. Loading, Empty, Partial, Error, Offline, Permission denied and Stale follow the Design System.
