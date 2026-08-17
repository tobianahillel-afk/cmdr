---
id: studio-std3-migration-audit
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
---
# STD-3 Migration Audit

The complete PR path manifest and current Studio/Govern/Shared/Settings/Endpoint owner sources were audited for historical agent, automation, autonomous execution, runtime, run, Human Gate, approval-like, queue/scheduler, intervention, monitoring, Agent Team and run-history concepts.

## Findings
- Automation Agent, Agent Team and Human Gate already have canonical Studio objects; no duplicate object is created.
- Automation Run is Studio-owned in the Ownership Register; its physical canonical object file remains deferred to Phase 7. STD-3 adds functional semantics only.
- Historical Human Gate `approved/rejected` wording remains in pre-existing object/docs. STD-3 does not rewrite that physical source; capabilities define `accepted-for-workflow` / `rejected-for-workflow` and preserve OPEN-007.
- Control Room already exists; STD-3 adds capability links/semantics only, no new Screen ID.
- Generic Job/queue/scheduling remains Shared; Response Run/Result remains Govern; Endpoint Agent/primitives remain Endpoint.
- No real competing canonical Studio runtime source requires deprecation.

Deprecated pointers created: **0**. Workflow, Tool, Skill, Govern Approval/Response Run, Endpoint Agent and Shared Job remain active and distinct.