---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: canonical
---
# Requirements Traceability Matrix — current capability-specification evidence

The **122 source Requirement IDs remain unchanged**. `conform` records documentary evidence only, never implementation. This matrix preserves the verified Command and Investigate evidence and adds GOV-1 documentary evidence without changing any global Requirement state.

| State | Current |
|---|---:|
| conform | 99 |
| partial | 20 |
| absent | 3 |
| contradictory | 0 |
| **Total** | **122** |

## Capability Specification Phase 4A — Command evidence

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-CMD-001..006 | Mission Control: situation, priority, timeline, handover, blockers and recent Results | REQ-PROD-003,005,006,008,010,013,021; REQ-UX-007 and supporting requirements | none | functional coordination is defined; implementation, final objects, permissions and screens remain future |
| CAP-CMD-101..110 | Incidents and Work Queue: one queue, assignment, ownership, priority/severity, SLA, Incident/Task coordination, bulk, freshness and escalation | REQ-OBJ-001,012; REQ-PROD-003,004,005,006,008,009,013,021,053; REQ-SEC-001; REQ-UX-005,008,009 | none | runtime, atomic permissions and final Incident/Task state machines remain future |
| CAP-CMD-201..205 | Risk and Coverage: Service, Exposure, Coverage, business impact and risk-prioritization context | REQ-PROD-003,005,006,010,013,021,032,037 | none | Service/Exposure detailed objects and final risk engine remain future; no scanner or opaque score is claimed |
| CAP-CMD-301..305 | Readiness and Operations: readiness, exercises, improvement Tasks, plans and capability readiness | REQ-OBJ-012; REQ-PROD-005,009,012,013,019,021,057 | none | assessments are functional records only; Studio assurance, platform health, final objects and permissions remain source-owned/future |
| CAP-CMD-401 | Customers and Delivery proposal | REQ-PROD-013,019,033,053 | none | remains `proposed`, `planned`, deployment-dependent and governed by OPEN-006 |

Command evidence preserves these invariants: Incident/operational Task coordination remains Command-owned; Case/Evidence/Finding remain Investigate-owned; Action Request/Decision/Response Run/Result remain Govern-owned; Studio automation objects, Settings administration and Shared engines are consumed but not redefined. The six Work Queue system views remain configurations of one workspace, not separate pages.

## Govern GOV-1 evidence

Parent roadmap: **Delivery Roadmap Phase 4 — Govern** (`roadmap-phase-4-govern`). Execution lot: **GOV-1**. GOV-1 is not a roadmap phase and `Phase 4C Govern` does not exist.

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-GOV-001..003 | Govern intake, Response Inbox and Action Request processing lifecycle | REQ-PROD-004,006,008,015,019,020; REQ-SEC-001,002; REQ-UX-008,009 | none | received/submitted/complete/authorized remain distinct; no general Command Work Queue duplication or execution |
| CAP-GOV-004..006 | context/scope/target, impact/risk/reversibility and completeness/Evidence-context review | REQ-PROD-003,004,008,014,015,020; REQ-SEC-001,002 | none | target verification is not authorization; Govern does not requalify Evidence/Finding; no opaque universal risk score |
| CAP-GOV-007..008 | Policy applicability/evaluation, conflicts and Exception Candidates | REQ-PROD-004,015,019,020; REQ-SEC-001,002 | none | Policy outcome is not Decision; conflict is not automatic rejection; Exception Candidate is not active exception |
| CAP-GOV-009..013 | authority requirements/context, approver eligibility/SoD, Approval lifecycle, delegation/escalation and emergency governance | REQ-PROD-004,015,020; REQ-SEC-001,002; REQ-AI-004 | none | authority remains distinct from Role/permission; Human Gate is not Approval; no self-approval under applicable SoD or uncontrolled emergency bypass |
| CAP-GOV-014..016 | Decision preparation, Decision disposition/conditions/expiration and no-effect Execution Handoff Package | REQ-PROD-004,008,015,016,020; REQ-AI-002; REQ-SEC-001,002 | none | no auto-Decision; approve is not execution; package is not Response Run; GOV-2/GOV-3 remain not started |

GOV-1 documentary evidence strengthens the functional definition of Govern but does not promote any Requirement from partial/absent to conform because final object schemas, atomic permissions, policy/authority engines, detailed screens, response execution, technique and implementation remain future.

## Investigate evidence preserved

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-001..397 | Investigate foundation, collection and analysis | REQ-INV-001..005 and supporting requirements | none | final engines, models and implementation remain future |
| CAP-INV-401..435 | complete Detection Engineering lifecycle | REQ-INV-006 and supporting product, AI, security and UX requirements | none | runtime, language and implementation remain future |
| CAP-INV-501..518 | Threat Intelligence intake, requirements, sources, materials, candidate knowledge, Sightings, relationships, confidence, lifecycle and analysis handoff | REQ-PROD-014,019,020,055,060,061,062; REQ-INV-006; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | final ontology, objects, permissions, screens, exchange and implementation remain future |
| CAP-INV-519..524 | sessions, competing hypotheses, fusion and actor, campaign, malware and infrastructure assessments | REQ-PROD-014,019,020; REQ-INV-006; REQ-AI-002; REQ-SEC-001,002; REQ-UX-010 | none | final objects, algorithms and implementation remain future |
| CAP-INV-525..529 | product planning, authoring, review, releasability and internal publication | REQ-PROD-014,019,020,055; REQ-INV-006; REQ-AI-002; REQ-SEC-001,002; REQ-UX-010 | none | OPEN-019 and final permissions, screens and delivery contracts remain future |
| CAP-INV-530..533 | watchlist definition, operationalization, monitoring and external-sharing preparation | REQ-PROD-014,019,020,055; REQ-INV-006; REQ-SEC-001,002 | none | no activation, deployment, runtime mutation or actual sharing |
| CAP-INV-534..537 | feedback, effectiveness, satisfaction, correction and lifecycle closure | REQ-PROD-014,019,020; REQ-INV-006; REQ-AI-002 | none | final metrics, contracts and implementation remain future |
| CAP-INV-601..603 | Cloud intake, preconditions, Session and provider-neutral scope hierarchy | REQ-PROD-006,012,014,019,020,055; REQ-INV-001,006; REQ-AI-002; REQ-SEC-001,002; REQ-UX-010 | none | OPEN-008/012 and final objects, permissions, screens and implementation remain future |
| CAP-INV-604..607 | Cloud inventory, identities, roles, IAM permission candidates and audit activity | REQ-PROD-014,019,020,055; REQ-INV-001,006; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-010 | none | provider schemas, connectors, effective-permission implementation and event contracts remain unresolved |
| CAP-INV-608..614 | Cloud configuration, workloads, containers, serverless, network, storage and sensitive-material assessment | REQ-PROD-014,019,020,055,060,061,062; REQ-INV-001..006; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | no command, scan, acquisition, secret use, provider implementation or target mutation |
| CAP-INV-615..618 | Cloud anomalies, Hypotheses, Timeline, correlations, Evidence/Finding/Detection handoffs and provenance | REQ-PROD-002,005,006,008,010,011,012,014,019,020; REQ-INV-001,006; REQ-AI-001..011; REQ-SEC-001,002; REQ-UX-006,010 | none | candidates and packages do not create canonical Evidence, Findings, rules, collection or response automatically |
| CAP-INV-701..705 | Mobile intake, Session, device/platform/scope, acquisition context and integrity/completeness/accessibility | REQ-PROD-006,012,014,019,020,055; REQ-INV-001; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-010 | none | OPEN-008/011/013/014/015 and final platform, acquisition, objects, permissions and screens remain future |
| CAP-INV-706..715 | Mobile filesystem/storage, applications/data, communications, media, location/sensors, sensitive material, connectivity, backups/sync and deleted/recovered data | REQ-PROD-014,019,020,055,060,061,062; REQ-INV-001..005; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | no device acquisition/mutation, unlock/bypass, secret use, final parsers/tools/platform support or implementation |
| CAP-INV-716..719 | Mobile Timeline/correlation, anomaly/Hypothesis, Derived Artifacts/handoffs and provenance/reproducibility | REQ-PROD-002,005,006,008,010,011,012,014,019,020; REQ-INV-001,006; REQ-AI-001..011; REQ-SEC-001,002; REQ-UX-006,010 | none | correlation/candidates/packages do not create canonical Evidence, Findings, rules, collection or response automatically |

## Current disposition

- Capability Specification Phase 4A — Command remains **PASS** and its five evidence ranges above are preserved unchanged.
- Capability Specification Phase 4B — Investigate remains **PASS**; no Investigate capability is created or modified by GOV-1.
- GOV-1 functional documentary scope is complete but remains **PENDING POST-PUBLICATION VERIFICATION** until the fifth functional commit and remote gates complete.
- Govern capability specification remains **PARTIAL** because GOV-2 and GOV-3 are NOT STARTED.
- Delivery Roadmap Phase 4 — Govern remains **PARTIAL/PENDING GOV-1 verification**, never `Phase 4C`.
- Capability Specification global maturity remains **PARTIAL** because final object schemas, atomic permissions, detailed screens, technique and implementation remain future.
- REQ-INV-001 and REQ-INV-006 remain globally partial because implementation, final technical contracts, object schemas, atomic permissions and detailed UX remain absent.
- OPEN-006/007/010/011/012/013/014/015/016/017/018/019 and all other current decisions retain their existing dispositions; GOV-1 closes none and creates none.
- Requirement IDs added: **0**; removed: **0**; state changes: **0**; active contradictions introduced: **0**.