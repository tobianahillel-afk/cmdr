---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---

# Cross-product links — Investigate

| Transition | Source | Destination | Contexte | Ownership et retour |
|---|---|---|---|---|
| Signal → investigation | Command Signal | Triage, Search ou Case | tenant, env, source, Detection, Events, Incident | Signal reste Command ; retour exact |
| Incident → Case | Command Incident | Case Lifecycle | Incident, service, Signal/Alert, période, objectif | Case Investigate ; retour Incident Detail |
| Search → Case | Search Job | Case Workspace | Query/version, run, résultats, période, annotations | Events/Query Shared ; retour Search |
| Hunt → Case | Hunt workspace | Case Lifecycle | question, scope, queries, results, Hypotheses | Case Investigate ; Hunt conservé |
| Artifact → Evidence | Artifact | Evidence Creation | source, acquisition, Case, version, transformations, reason | deux objets Investigate distincts |
| Evidence → Finding | Evidence/qualification | Finding Management | Evidence pour/contre, reviewer, uncertainty | Finding distinct |
| Finding → Action Request | Finding | Action Request Preparation | Finding, Evidence, target, impact, urgency, alternatives, rollback | Request Govern ; Finding/Evidence Investigate |
| Action Request → Govern | Investigate producer | Govern | request, refs, requester, return origin | Govern lifecycle ; retour Case/Finding |
| Result → Case | Govern Result | Case Workspace | Request, Decision, Run, Result, residual risk | Result Govern ; relation Investigate |

Une transition échouée conserve le workspace source, les drafts récupérables et une correlation ID.
