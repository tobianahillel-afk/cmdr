---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.3A.1

| Class | Meaning | Detection Engineering Authoring examples |
|---:|---|---|
| 0 | observation | consult, navigate, filter, search, compare, read data/schema/validation/replay/coverage projections |
| 1 | bounded analytical execution | validation, controlled test, authorized historical replay, comparison, permitted export |
| 2 | reversible analytical mutation | create/update Project, Detection Hypothesis, Draft, metadata, logic, scenario, Expected Outcome, candidate match disposition, coverage/gap and review package |
| 3/4 | production, authority or destructive action | excluded; promotion, deployment, activation, deactivation, rollback, active exception/suppression and runtime object changes are future 4B.3A.2/Govern/owner actions |

CAP-INV-401..417 use only classes 0, 1 and 2. OPEN-013 remains open. No action changes a runtime Detection, Signal, Alert or Incident.
