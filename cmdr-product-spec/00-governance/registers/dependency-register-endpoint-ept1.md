---
id: dependency-register-endpoint-ept1
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-1

| ID | Source | Dependent | Type | Ownership / rule | Blocking |
|---|---|---|---|---|---|
| DEP-EPT1-001 | canonical Endpoint Agent object | CAP-EPT-001..014 | object | Endpoint owns local Agent identity/state | yes |
| DEP-EPT1-002 | Platform Settings enrollment/tenant/environment | CAP-EPT-002..004/012..014 | administrative projection | Settings administers; Endpoint consumes local/effective refs | yes where referenced |
| DEP-EPT1-003 | Platform support + OPEN-008 | CAP-EPT-004/005/011 | open platform decision | no platform/version support claim | before delivery, not documentary definition |
| DEP-EPT1-004 | Settings Fleet/capability inventory/upgrade management | CAP-EPT-005/011/012/014 | administrative projection | no Fleet or upgrade ownership transfer | yes for admin projection |
| DEP-EPT1-005 | Settings Endpoint Policy/policy assignment | CAP-EPT-013/014 | policy projection | Policy/assignment remain Settings-owned | yes where assigned context is used |
| DEP-EPT1-006 | Endpoint health monitoring | CAP-EPT-008..011 | technical source | local health/heartbeat/degradation only | yes for declared health facts |
| DEP-EPT1-007 | Security permission/tenant isolation | CAP-EPT-001..014 | security | no final atomic RBAC; no cross-tenant fallback | yes |
| DEP-EPT1-008 | Govern Decision/Response Run/Result/rollback | CAP-EPT-014 and future execution | boundary | technical state/result != Govern canonical lifecycle | no EPT-1 execution |
| DEP-EPT1-009 | Studio Tool/Automation Agent/Automation Run | CAP-EPT-011/014 | boundary | Endpoint capability/Agent/execution remain distinct | no |
| DEP-EPT1-010 | Shared Jobs/Trace/Activity/Search/Reporting/Export/Versioning/Recovery | CAP-EPT-014 | shared mechanisms | no competing Shared engine | no |
| DEP-EPT1-011 | Investigate/Command consumers | CAP-EPT-006..014 | projection/handoff | facts do not become Evidence/Finding/Incident ownership | no |
| DEP-EPT1-012 | OPEN-013/015 | class-2/provenance bridge | open decision | no default authority or run bridge selected | before affected implementation |

EPT-1 introduces no protocol, API, PKI, port, certificate/token format, code or EPT-2+ technical execution dependency.