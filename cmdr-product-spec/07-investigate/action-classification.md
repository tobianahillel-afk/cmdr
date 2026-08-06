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
# Action classification — Investigate through Phase 4B.3A

| Class | Meaning | Detection Engineering examples | Local execution |
|---:|---|---|---|
| 0 | observation | consult, navigate, filter, compare, read candidate/readiness/Decision/Result/version/health/Signal feedback | allowed by permission |
| 1 | bounded analytical execution | deterministic validation, shadow assessment request where authorized, comparison, health/drift/performance assessment, permitted export/reproduction | explicit, bounded and attributed |
| 2 | reversible analytical mutation | create/update Release Candidate review, readiness, Promotion Plan, Change Request Draft, Tuning/Suppression/Exception Proposal, Drift Assessment, Rollback Plan, Retirement Proposal and improvement package | Investigate-owned reversible concepts; OPEN-013 |
| 3 | governed production change | promotion, activation, deactivation, active suppression/exception, rollback, retirement and modification of an active version | **not executed by Investigate**; Govern/Settings/runtime owner |
| 4 | destructive/irreversible | deletion of active content or provenance and change without rollback | denied by default or strictly governed; never recommended |

CAP-INV-401..435 expose classes 0–2 locally. Any class-3 projection is prepared or observed only. No Signal, Alert, Incident or provenance is silently mutated or deleted.
