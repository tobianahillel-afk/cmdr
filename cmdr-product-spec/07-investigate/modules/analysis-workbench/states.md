---
id: analysis-workbench-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
---
# Functional states

| Scope | States |
|---|---|
| Artifact routing | recognized, ambiguous, unsupported, restricted, corrupted, encrypted, incomplete, ready-for-analysis |
| Analysis Session | draft, ready, active, paused, blocked, partial, completed, failed, cancelled, archived, superseded |
| Tool execution | validating, queued, running, partial, completed, failed, cancelled, unavailable, incompatible, permission-denied |
| Derived Artifact | proposed, producing, available, partial, invalid, restricted, superseded, withdrawn-from-use |
| Reproducibility | reproducible, partially-reproducible, not-reproducible, missing-tool, missing-version, missing-input, environment-unavailable, disputed |

Ces états sont fonctionnels ; les machines d’état d’objet sont reportées.
