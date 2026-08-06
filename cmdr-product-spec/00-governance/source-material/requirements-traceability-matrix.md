---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
source-of-truth: canonical
---
# Requirements Traceability Matrix — through Phase 4B.4A pre-publication closure

The 122 source Requirement IDs remain unchanged. `conform` records documentary evidence only, not implementation. Cloud Analysis adds functional evidence without promoting a requirement whose implementation, provider strategy, object model, permission model or detailed UX remains unresolved.

| State | After 4B.3B.2 | After 4B.4A documentary update |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| **Total** | **122** | **122** |

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

## Disposition
- Threat Intelligence retains complete functional documentary evidence from intake through lifecycle closure without claiming implementation.
- Cloud Analysis now has provider-neutral functional evidence from intake through provenance, but the fifth commit and remote post-publication verification are still required before the phase verdict can become final.
- REQ-INV-001 and REQ-INV-006 remain globally partial because implementation, final technical contracts and Mobile Forensics are absent.
- OPEN-012 remains open for Cloud provider, service, SaaS, inventory, evidence, cross-account, cross-tenant and delivery strategy.
- OPEN-018 records unresolved Threat Intelligence ontology, interoperability and exchange; OPEN-019 records dissemination, releasability, sharing and access policy.
- OPEN-017 remains Detection-only and unchanged.
- REQ-PROD-055 remains partial and dependent on OPEN-008/012; source and provider support is not invented.
- REQ-PROD-020 and REQ-OBJ-009 remain partial because final trace and run contracts are not defined.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015/019 and final permissions and object relations.
- REQ-UX-010 remains partial because no detailed Cloud Analysis screen is created or rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.
