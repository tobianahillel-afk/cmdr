---
id: dependency-register-endpoint-ept2
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-2

| Dependency | Owner | Consumed by | Contract | Blocking meaning |
|---|---|---|---|---|
| EPT2-SHARED-TELEMETRY | Shared | CAP-EPT-016/023/029/030 | `telemetry-event` envelope and generic normalization remain Shared-owned | blocks duplicate Endpoint event object |
| EPT2-SOURCE-FAMILIES | Endpoint | CAP-EPT-015..026 | 11 canonical Telemetry documents define source-backed observation families/limitations | no source fact may be invented |
| EPT2-SETTINGS | Platform Settings | CAP-EPT-015/025..029 | source administration, Fleet, Policy, tenant and configuration ownership | Endpoint may project facts, not take admin ownership |
| EPT2-INVESTIGATE | Investigate | CAP-EPT-029 | Evidence/Finding/Case interpretation remains Investigate-owned | telemetry never auto-promotes |
| EPT2-STUDIO | Studio | CAP-EPT-027..029 | Tool/Tool Call/Automation Run ownership preserved | Endpoint capability != Tool |
| EPT2-GOVERN | Govern | CAP-EPT-029/030 | Decision/Response Run/Result ownership preserved | observation != Result |
| EPT2-OPEN-008 | Product Architecture | CAP-EPT-015..030 | platform/source availability and support remains unresolved | no delivered-platform claim |
| EPT2-EPT1 | Endpoint | CAP-EPT-015..030 | CAP-EPT-001..014 are consumed intact | no EPT-1 rewrite |
