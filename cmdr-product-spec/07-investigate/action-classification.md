---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-019]
---
# Action classification — Investigate through Phase 4B.3B.2

| Class | Meaning | Threat Intelligence Analysis/Products examples | Local execution |
|---:|---|---|---|
| 0 | observation | read/filter/search/compare Sessions, hypotheses, assessments, products, publications, watchlist projections, feedback and provenance | allowed by permission |
| 1 | bounded analytical execution | structured analysis, source fusion, comparison, deterministic quality checks, monitoring request/run, package generation and authorized export | explicit, bounded, attributed and reversible |
| 2 | reversible analytical mutation | create/update/review/version Session, hypothesis, Product Draft, Release Recommendation, Dissemination Plan, Internal Publication Record, Watchlist Definition, handoff, feedback, correction/retraction and improvement package | Investigate-owned concepts; OPEN-013/019 |
| 3 | governed production/external action | external sharing, client/public publication, watchlist activation, Indicator deployment, runtime change, global revocation/recall | not executed by Investigate; Govern and runtime/destination owner |
| 4 | destructive/irreversible | provenance/history destruction or non-revocable sensitive disclosure | denied by default or strictly governed |

CAP-INV-001..537 expose only classes 0–2 locally. Internal publication in CAP-INV-529 is policy-bound, reversible and tenant-scoped; it is not class-3 external sharing. Any class-3/4 action is blocked or prepared as an Action Request/handoff. No Signal, Alert, Incident, rule, active watchlist, deployed Indicator, source, permission or provenance is silently mutated.
