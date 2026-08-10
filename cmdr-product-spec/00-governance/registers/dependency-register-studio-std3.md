---
id: dependency-register-studio-std3
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-10
source-of-truth: registry
---
# Dependency Register — Studio STD-3

Additive STD-3 evidence. It does not replace historical dependency evidence.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN |
|---|---|---|---|---|---|---|---|---|
| DEP-STD3-001 | CAP-STD-001..016 | CAP-STD-034..051 | STD-1 consumption | Tools, Tool Calls, Skills, eligibility and secret/runtime boundaries are consumed, never redefined | active | Studio | yes where used | REQ-PROD-016 |
| DEP-STD3-002 | CAP-STD-017..033 | CAP-STD-035/038/040..048/051 | STD-2 consumption | Workflow/version/steps/retry/compensation/Human Gate definitions remain STD-2-owned | active | Studio | yes | REQ-PROD-016 |
| DEP-STD3-003 | Automation Agent / Agent Team / Human Gate canonical objects | CAP-STD-034..041 | object | preserve existing Studio ownership; no physical schema rewrite | active | Studio | yes | REQ-OBJ-009 |
| DEP-STD3-004 | Automation Run ownership entry | CAP-STD-042..051 | object boundary | functional semantics now defined; physical canonical object remains deferred Phase 7 | partial | Studio/Object programme | yes | OPEN-015 |
| DEP-STD3-005 | Security permission/tenant isolation | CAP-STD-034..051 | authorization | Agent role/objective/access never becomes permission; cross-tenant access independently checked | active | Security | yes | REQ-SEC-001/002; OPEN-013 |
| DEP-STD3-006 | Platform Settings principals/providers/secrets/environments/health | CAP-STD-034..049/051 | admin projection | Settings retains identity/provider/secret/runtime administration | active/partial | Settings | where used | REQ-SEC-001 |
| DEP-STD3-007 | Shared Jobs/queue/scheduling/Trace/Activity/Notifications/Recovery | CAP-STD-040/043..051 | shared mechanism | generic mechanisms remain Shared; Automation Run retains Studio semantics | partial | Shared | where used | REQ-PROD-019 |
| DEP-STD3-008 | Govern Approval/Decision/Response Run/Result | CAP-STD-039..051 | authority/handoff | Human Gate/Automation Run/outcome never replace Govern authority or canonical Result | active/partial | Govern | effectful paths | OPEN-007/013/015 |
| DEP-STD3-009 | CAP-GOV-025/026/027/032 | CAP-STD-042/047/050/051 | execution bridge | bounded handoff, technical reconciliation/error and Result boundaries preserved | active | Govern/Studio | yes when linked | OPEN-015 |
| DEP-STD3-010 | Endpoint Agent technical primitives | CAP-STD-036/045..051 | future executor boundary | Studio may reference future endpoint capability but creates none | future | Endpoint | no for Studio-only path | OPEN-008 |
| DEP-STD3-011 | STD-4 Assurance/Evaluation/Simulation/Deployment | CAP-STD-034/037/049/051 | later-lot boundary | runtime control does not claim assurance/promotion/deployment | future | Studio STD-4 | no for STD-3 definition | REQ-PROD-016 |

No dependency transfers object ownership, permission, authority or implementation ownership.