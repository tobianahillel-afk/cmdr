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
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Investigate through Phase 4B.2A

## Phase 4B.1 transitions retained
Signal→Investigation; Incident→Case; Search/Hunt→Case; Artifact→Evidence; Evidence→Finding; Finding→Action Request; Action Request→Govern; Result→Case; Investigate→Studio remain canonical and ownership-preserving.

## Collection and Live Response
| Transition | Source / trigger | Destination | Context transmitted | Ownership | Return and errors |
|---|---|---|---|---|---|
| Case → Endpoint Context | open linked Endpoint | CAP-INV-201 | tenant, environment, Case, Incident, Endpoint, objectives, return origin | Case Investigate; Fleet/Policy Settings; Agent Endpoint | same Case; unresolved/offline explicit |
| Endpoint Context → Collection Request | prepare collection | CAP-INV-202 | target, Case, scope, policy, capabilities, class, permissions | Request Investigate | Endpoint Context/Case; draft preserved |
| Collection Request → Collection Job | submit/accept | CAP-INV-203 | request/version, target, initiator, authority, scope, limits, provenance | business context Investigate; generic Job Shared | Request/Case; no false start |
| Collection Job → Artifact | successful item | CAP-INV-105 | result, source, acquisition, timestamps, transformations, errors, Case | Artifact Investigate; Job distinct | Job/Case; partials retained |
| Artifact → Evidence | human qualification | CAP-INV-107 | existing 4B.1 contract | both Investigate, distinct objects | Artifact; no automatic conversion |
| Case → Live Session | request session | CAP-INV-209 | Case, Endpoint, policy, permissions, reason, class, expiry, participants, return origin | session context Investigate; local execution Agent | Case/Endpoint Context; offline/conflict explicit |
| Live Session → Endpoint Operation | explicit operation | CAP-INV-210 | session, operation category, target, scope, initiator, authority, trace | operation context Investigate; Agent Command Endpoint | session; denied/policy-blocked explicit |
| Endpoint Operation → Operation Result | local completion/error | CAP-INV-212 | operation, target, output, errors, files, timestamps, status, provenance | local result context Investigate; not Govern Result | session/Case |
| Operation Result → Case/Artifact | review/link | CAP-INV-102/105 | result, Artifact, next action, verification, trace | Case/Artifact Investigate | result/session |
| Finding → Containment Request | prepare endpoint containment | CAP-INV-215 then 113 | Finding, Evidence, Endpoint, action, impact, urgency, alternatives, rollback | Finding/Evidence Investigate; Request lifecycle Govern | Case |
| Containment Request → Govern | submit | Govern Review Queue | Action Request, class, target, risks, permissions, provenance | Govern owns Decision/Run/Result | same draft/Case on return |
| Govern Result → Case | verified result | Case Workspace | Decision, Response Run, Result, verification, effect on Finding, next action | Result remains Govern | Result/Incident/Case |
| Investigate → Studio | request deployed assistance | Workflow/Automation Run | Case/Endpoint/scope, workflow version, permissions, limits, sources | Studio owns Workflow/Run/Tools | source workspace with proposal |
| Settings → Investigate projection | open from Fleet/Policy/Health | CAP-INV-201 | Endpoint/Agent identity, status, freshness, capabilities, policy reference | Settings owns administration | Settings source or Case |

## Invariants
Tenant, environment, selection and return origin are preserved; permission is reevaluated; projections never transfer ownership; errors preserve drafts with correlation ID; class 3/4 always route through Govern; no endpoint operation output is presented as Govern Result.