---
id: dependency-register-studio-std2
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-10
source-of-truth: registry
---
# Dependency Register — Studio STD-2

Additive evidence only. It does not replace the global Dependency Register or STD-1 dependency evidence.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN |
|---|---|---|---|---|---|---|---|---|
| DEP-STD2-001 | CAP-STD-017/018 | CAP-STD-019..033 | definition/editing | exact Workflow/base-version context | active | Studio | yes | REQ-PROD-016 |
| DEP-STD2-002 | CAP-STD-003..014 | CAP-STD-021/023/024/030/033 | STD-1 consumption | Tool/Skill definitions, I/O and versions are referenced, never redefined | active | Studio | yes | REQ-PROD-016 |
| DEP-STD2-003 | canonical Workflow/Version objects | CAP-STD-017..033 | object | preserve Studio object ownership and version lifecycle | active | Studio | yes | REQ-OBJ-009 |
| DEP-STD2-004 | Platform Settings secret/provider/runtime refs | CAP-STD-019/023/025/030/033 | admin projection | no raw secret/provider ownership transfer | partial | Settings | where used | REQ-SEC-001 |
| DEP-STD2-005 | Shared Versioning/Trace/Activity/Jobs/Search/Recovery | CAP-STD-018/024/025/026/031/033 | shared | generic mechanisms remain Shared | partial | Shared | where used | REQ-PROD-019 |
| DEP-STD2-006 | Govern Approval/Decision/Playbook/Response Run | CAP-STD-017/022/027/028/029/030/033 | authority | Workflow/Human Gate/compensation never replace Govern authority | active/partial | Govern | effectful paths | OPEN-007/013/015 |
| DEP-STD2-007 | CAP-GOV-025 | CAP-STD-029/033 | handoff | future execution correlation preserves Decision/Response Run ownership | active | Govern/Studio | yes when used | OPEN-007/015 |
| DEP-STD2-008 | Workflow implementation contract | CAP-STD-020/022/025/026/027/028/029 | technical boundary | graph/retry/compensation/Human Gate semantics remain implementation-agnostic | partial | Product Architecture | before implementation | REQ-PROD-019 |
| DEP-STD2-009 | idempotency contract | CAP-STD-027 | technical boundary | retry-safe expectations without exactly-once claim | partial | Product Architecture | before runtime | OPEN-015 |
| DEP-STD2-010 | Security permission model | CAP-STD-017..033 | permission | functional permission needs only; namespace anomaly preserved | partial | Security | yes | REQ-SEC-001/002; OPEN-013 |
| DEP-STD2-011 | future STD-3 runtime | CAP-STD-025..030/033 | future boundary | scheduling, Automation Run, pause/stop and runtime lifecycle are excluded from STD-2 | planned | Studio | before runtime | OPEN-015 |
| DEP-STD2-012 | future STD-4 assurance/deployment | CAP-STD-030..032 | future boundary | validation/review candidate does not publish/deploy | planned | Studio | before promotion | REQ-PROD-016 |
| DEP-STD2-013 | Endpoint future technical capability | CAP-STD-033 | future boundary | Workflow may reference, never own, Endpoint primitive | planned | Endpoint Agent | before endpoint execution | OPEN-008 |
| DEP-STD2-014 | Command/Investigate consumer refs | CAP-STD-017/033 | cross-product | consumers retain Incident/Case/Evidence/Finding ownership and return-origin | active | source owners | yes for handoff | REQ-PROD-013/014 |
| DEP-STD2-015 | OPEN-007/013/015 | STD-2 | decision | non-equivalence fixed; unresolved policy/bridge choices remain open | open | Product Architecture/Security/Govern/Studio | before final runtime policy | OPEN-007/013/015 |

No dependency selects an orchestration runtime, language, API, protocol, scheduler, provider or endpoint implementation.
