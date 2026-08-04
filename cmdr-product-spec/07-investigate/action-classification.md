---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.2B.2A

| Class | Meaning | Dynamic examples |
|---:|---|---|
| 0 | observation | read timeline, process tree, changes, network and comparisons |
| 1 | bounded isolated processing | explicit Run start, Runtime Artifact capture, controlled export/reproduction |
| 2 | reversible mutation | session/profile update, stop/retry, annotation, relations and handoffs |
| 3/4 | real-target response/destructive action | excluded; route to Collection/Live Response and Govern |

CAP-INV-314..328 use only classes 0, 1 and 2. OPEN-013 remains open. A sandbox safety stop is not Endpoint containment.
