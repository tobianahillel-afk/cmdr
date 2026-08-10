---
id: studio-std3-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-10
source-of-truth: quality-report
---
# Studio STD-3 — Source Audit

Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint** (`roadmap-phase-5-studio-and-endpoint`). Execution lot: **STD-3 — Agents, Human Gates & Runtime Control**.

## Git baseline audited
- repository: `tobianahillel-afk/cmdr`;
- branch: `docs/cmdr-product-spec-foundation`;
- expected and observed baseline: `c472b055ce00fd33efd96ac920b0add5f65ab8f7`;
- PR #2: open, Draft, unmerged, base `main`;
- branch/main README: exact `# cmdr`, same blob;
- `main` unchanged at the pre-STD-3 baseline;
- STD-2 chain: five functional commits plus two post-publication documentary corrections preserved;
- `CAP-STD-034` and `CAP-STD-051` searches returned no result; complete PR manifest showed no existing STD-3 capability file;
- `CAP-EPT-*` search returned no result.

## Directly reread — owner and runtime sources
Governance/source-of-truth, ownership, Capability Register, unresolved decisions, capability template, AI constraints, Studio action classification.

Studio:
- `automation-agents/README.md`, `automation-agent-model.md`, `tool-access.md`, `memory-and-context.md`;
- `agent-teams/README.md`, `coordination.md`, `teams-of-agents.md`;
- `human-gates/README.md`, `gate-model.md`, `approval-patterns.md`;
- `control-room/README.md`, `runtime-supervision.md`, `agent-team-operations.md`, `deployment-monitoring.md`;
- Automation Agent, Agent Team and Human Gate canonical object files;
- Studio README, Capability Map, Object Consumption Map, Automation/AI Model, Cross-product Links and Screen Capability Map.

Implementation/shared:
- Automation Agent, Agent Team, Human Gate and Response Run contracts;
- idempotency contract;
- Shared Background Jobs, Shared README and Notification Center.

Govern:
- execution boundaries;
- Runs & Rollback README;
- CAP-GOV-022, CAP-GOV-023, CAP-GOV-025, CAP-GOV-026, CAP-GOV-027 and CAP-GOV-032.

Settings:
- Users & Roles;
- Models & Providers;
- Secrets & Connections;
- Tenants & Environments;
- Health.

Endpoint boundary:
- Endpoint Agent README;
- Endpoint queueing-and-retry.

Screens:
- all eleven active Studio Screen IDs `STD-AGT-001`, `STD-ATM-001`, `STD-ASR-001`, `STD-CTL-001`, `STD-WFL-001`, `STD-BLD-001`, `STD-LIB-001`, `STD-SKL-001`, `STD-SIM-001`, `STD-EVL-001`, `STD-DEP-001`.

Quality/history:
- STD-1 capability conformance report;
- STD-2 build-time capability conformance report;
- STD-2 post-publication verification;
- current STATUS and Delivery Roadmap Phase 5.

## Complete-manifest structural audit
The complete PR changed-file manifest was inspected to identify all Studio agent/team/Human Gate/Control Room/runtime, governance, Settings, Shared, Endpoint, screen, migration, report and capability paths. Manifest presence was used for discovery and collision detection only; it was not treated as evidence of implemented runtime behavior.

## Critical findings
1. Automation Agent, Agent Team and Human Gate are canonical Studio-owned objects.
2. Automation Run is Studio-owned in the Ownership Register, but its canonical physical object file remains deferred to Phase 7.
3. Human Gate historical object/docs contain approval-like wording/states; STD-3 must preserve the object while enforcing functional non-equivalence with Govern Approval/Decision.
4. Shared owns generic Background Jobs/queue mechanics; Studio owns Automation Run business semantics.
5. Govern owns Response Run, execution authority, verification/rollback governance and Result.
6. Settings owns principals/roles, providers, integrations, raw secrets, tenant/environment and runtime administration/health.
7. Endpoint owns technical endpoint execution, local state and its own queue/retry primitives.
8. Evaluation, Simulation, Assurance and Deployment are read only as STD-4 boundaries.
9. Permission namespaces `perm.studio.*` and `perm.cmdr-studio.*` coexist and remain unresolved.
10. OPEN-007, OPEN-013 and OPEN-015 remain open; no new OPEN is justified by this source audit.

## ID disposition
`CAP-STD-001..033` remain allocated and immutable. No collision/reservation was found for `CAP-STD-034..051`; exactly eighteen STD-3 capability IDs are therefore allocated by the execution lot. No Endpoint ID is allocated.

## Verdict
**SOURCE PREFLIGHT PASS for STD-3 capability creation.** This audit authorizes only the functional documentary scope of STD-3; it selects no agent framework, model/provider, scheduler, runtime, API, protocol, JSON Schema, final RBAC/ABAC or Endpoint capability.
