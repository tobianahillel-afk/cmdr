---
id: studio-std3-runtime-boundaries
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-002, REQ-AI-007, REQ-OBJ-009]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# STD-3 Runtime Boundaries

## Scope
STD-3 owns functional Studio semantics for Automation Agents, Agent Teams, Studio Human Gates, Automation Runs, runtime control, Control Room, Studio runtime outcome and runtime provenance.

## Upstream contracts consumed
- STD-1 remains source for Tool, Tool Call, Skill, Library and access/eligibility foundations.
- STD-2 remains source for Workflow Definition/Version/Builder, graph/steps, conditions/branches, mappings, subworkflows, retry/idempotency/compensation definitions, Human Gate step definition, readiness and pre-publish lifecycle.

STD-3 consumes these definitions; it does not recreate them.

## Govern boundary
Govern retains Action Request, Approval, Decision, Response Playbook, Response Run, production-response authority, verification/rollback governance and canonical Result.

- Automation Run ≠ Response Run.
- Human Gate ≠ Approval.
- Human Gate ≠ Decision.
- Human Gate response ≠ production authority.
- Agent proposal/plan ≠ action/authorization.
- Tool Call output ≠ Result.
- Studio compensation ≠ Govern rollback.
- Studio runtime outcome ≠ Govern Result.

`OPEN-007` and `OPEN-015` remain open; only non-equivalence is fixed.

## Settings boundary
Platform Settings retains principal/role administration, providers, integrations, credentials, raw secrets, tenant/environment and runtime administrative configuration/health. Studio consumes authorized projections and opaque Secret References only.

## Shared boundary
Shared retains generic Jobs, queue/scheduling infrastructure, Trace, Activity, Notifications, Search, Versioning and Recovery. Automation Run ≠ Shared Job. Studio owns the business semantics of its Run, not generic engines.

## Endpoint boundary
Endpoint retains endpoint technical execution primitives, local agent state, device-side result and endpoint queue/retry mechanics. STD-3 creates no Endpoint capability and no endpoint command.

## STD-4 boundary
Assurance, full Evaluations/Simulations, regression framework, publishing/promotion/deployment lifecycle and release channels remain STD-4. Control Room runtime supervision in STD-3 does not transfer Deployment ownership into this lot.

## Runtime safety
- Run created ≠ started; queued/scheduled ≠ started; start requested ≠ running.
- paused ≠ stopped; stop requested ≠ stopped; cancellation ≠ rollback.
- retry attempt ≠ new Run; idempotency ≠ exactly-once.
- Tool Call/step success ≠ Run success; partial completion ≠ success.
- runtime status ≠ source business-object status.
- transient context ≠ canonical object or permanent memory.
- AI recommendation ≠ authority.

Every effectful start/resume/retry remains subject to current permission, eligibility, scope and Govern authority when applicable. No class-4 destructive authority is created by Studio.
