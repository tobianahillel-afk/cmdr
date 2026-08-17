---
id: investigate-memory-forensics-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Functional states

| Work item | States |
|---|---|
| Intake | draft, incomplete, ready, partial-image, corrupted, restricted, unsupported, profile-required, tool-unavailable, policy-blocked |
| Session | draft, ready, active, paused, blocked, partial, completed, failed, archived, superseded |
| Integrity Review | unverified, verifying, verified, partially-verified, inconsistent, incomplete, corrupted, disputed, restricted |
| Platform/Profile | proposed, selected, confirmed-by-analyst, ambiguous, conflicting, unsupported, superseded |
| Reconstruction Result | queued, processing, partial, available, failed, incompatible, disputed, superseded |
| Sensitive Candidate | detected, masked, restricted, under-review, invalid, confirmed-exposure, disputed, superseded |
| Memory Anomaly | candidate, weakly-supported, supported, contradicted, inconclusive, disputed, superseded, withdrawn |
| Derived Artifact | proposed, extracting, available, partial, invalid, restricted, superseded, withdrawn-from-use |
| Reproducibility | reproducible, partially-reproducible, not-reproducible, missing-image, incomplete-image, missing-profile, missing-tool, missing-version, unsupported-platform, policy-blocked, disputed |

Ces états sont fonctionnels et ne remplacent aucune machine d’état objet future.
