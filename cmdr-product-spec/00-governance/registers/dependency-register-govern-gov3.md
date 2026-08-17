---
id: dependency-register-govern-gov3
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: registry-shard
---
# Dependency Register — Govern GOV-3

This additive shard preserves the historical global Dependency Register unchanged.

| ID | Source | Dependent | Type | Reason | Status | Owner | Requirement / OPEN |
|---|---|---|---|---|---|---|---|
| DEP-GOV-021 | GOV-1/GOV-2 provenance | CAP-GOV-034..038 | audit source | audit consumes historical source objects/events without mutation | active | Govern/source owners | REQ-PROD-020; OPEN-015 |
| DEP-GOV-022 | Shared Trace/Activity/Search | CAP-GOV-034..038 | shared mechanism | Govern owns semantics, Shared owns generic infrastructure | partial | Shared/Govern | REQ-PROD-019,020 |
| DEP-GOV-023 | Settings retention/storage/tenant context | CAP-GOV-034..038/046 | availability/privacy | retention and storage remain Settings-owned | partial | Settings | REQ-SEC-001,002; OPEN-008/019 |
| DEP-GOV-024 | Security audit/privacy/integrity/export policy | CAP-GOV-034..038 | security | completeness/integrity claims remain bounded by security contracts | partial | Security | REQ-SEC-001,002 |
| DEP-GOV-025 | Shared Metrics Engine | CAP-GOV-039..046 | shared mechanism | Govern owns metric meaning, Shared owns generic calculation/version/privacy mechanism | partial | Shared/Govern | REQ-PROD-019,020 |
| DEP-GOV-026 | Shared Reporting/Export | CAP-GOV-038..047 | shared mechanism | rendering/export cannot widen source visibility | partial | Shared | REQ-PROD-019; OPEN-019 |
| DEP-GOV-027 | GOV-1 Policy/Authority/Approval/Decision | CAP-GOV-039..041/045 | metric source | source lifecycle remains GOV-1-owned | active | Govern | REQ-PROD-004,015 |
| DEP-GOV-028 | GOV-2 Run/Verification/Rollback/Result | CAP-GOV-042..045 | metric source | source lifecycle remains GOV-2-owned | active | Govern | REQ-OBJ-007; OPEN-015 |
| DEP-GOV-029 | product-owned metrics/context | CAP-GOV-044/046/047 | projection | Command/Investigate/Studio/Endpoint keep their metric ownership | partial | product owners | REQ-PROD-013..020 |
| DEP-GOV-030 | governance/quality/roadmap evidence | CAP-GOV-047 | closure | documentary closure requires exact registry/report/Git evidence | active | QA/Product Architecture | REQ-PROD-006,020 |

No dependency selects an audit engine, metrics engine implementation, warehouse, storage schema, API, protocol, query language, provider or runtime.