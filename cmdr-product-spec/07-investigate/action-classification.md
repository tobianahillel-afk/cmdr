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
# Action classification — Investigate through Phase 4B.3B.1

| Class | Meaning | Detection Engineering examples | Threat Intelligence examples | Local execution |
|---:|---|---|---|---|
| 0 | observation | consult/read candidate, readiness, version, health and feedback | consult sources, candidates, relations, Sightings, versions and provenance | allowed by permission |
| 1 | bounded analytical execution | validation, replay, comparison, health/drift/performance assessment | bounded extraction, normalization, comparison, duplicate search, assessment and authorized export | explicit, bounded and attributed |
| 2 | reversible analytical mutation | Drafts, reviews, plans, proposals and improvement packages | Requirements, Projects, candidates, relations, assessments, versions, supersession, expiry/revocation proposals and handoffs | Investigate-owned reversible concepts; OPEN-013 |
| 3 | governed production/external action | promotion, activation, rollback and active exceptions | external publication/sharing, active Indicator/watchlist, block or source-admin change | **not executed by Investigate**; Govern/owner |
| 4 | destructive/irreversible | deletion of active content/provenance | irreversible knowledge/history deletion or exchange without rollback | denied by default or strictly governed |

CAP-INV-001..518 expose classes 0–2 locally. Any class-3/4 action is blocked, prepared as a future request or routed to its owner. No Signal, Alert, Incident, source, Indicator, watchlist or provenance is silently mutated or deleted.
