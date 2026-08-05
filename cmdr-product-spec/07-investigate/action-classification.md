---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.2 closure

| Class | Meaning | Network Forensics examples |
|---:|---|---|
| 0 | observation | consult, navigate, filter, search, group, compare views, read metadata, inspect authorized packets/flows/conversations and view Network Timeline |
| 1 | bounded isolated processing or extraction | deterministic flow/session reconstruction, authorized decode, comparison, Derived Artifact extraction, permitted export and reproduction |
| 2 | reversible analytical mutation | session update, interpretation selection, annotation, dispute, Entity relation proposal, withdrawal from active use, access request and handoff preparation |
| 3/4 | real-target, destructive or authority action | excluded; packet generation/injection, active scanning/interception, replay, target interaction and any real-system action are blocked or routed to owner capabilities and Govern |

CAP-INV-380..397 use only classes 0, 1 and 2. OPEN-013 remains open. Network analysis is not acquisition, live monitoring, active response, packet crafting or Detection Engineering.
