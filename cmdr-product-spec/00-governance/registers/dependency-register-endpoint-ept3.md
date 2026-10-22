---
id: dependency-register-endpoint-ept3
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-3

| Family | EPT-3 scope | Upstream owner/source | Rule |
|---|---|---|---|
| DEP-EPT3-001 | EPT-1 identity/platform/health | Endpoint Agent CAP-EPT-001..014 | consume unchanged |
| DEP-EPT3-002 | EPT-2 observations/source/time | Endpoint Agent CAP-EPT-015..030 + Shared telemetry-event | consume unchanged; observation ≠ verdict |
| DEP-EPT3-003 | Detection Content authoring/lifecycle | Investigate Detection Engineering CAP-INV-401..435 | Endpoint references/version-checks only |
| DEP-EPT3-004 | canonical Detection/Signal/Alert/Incident | Command | local match/candidate are not competing canonical objects |
| DEP-EPT3-005 | Case/Evidence/Finding | Investigate | no automatic creation/qualification |
| DEP-EPT3-006 | generic Search/Timeline/Linking/Correlation/Trace | Shared | local projections only |
| DEP-EPT3-007 | sources/runtimes/targets/policies/assignments | Platform Settings | administrative ownership retained |
| DEP-EPT3-008 | Decision/Approval/Response Run/Result | Govern | no response authority in EPT-3 |
| DEP-EPT3-009 | Tool/Tool Call/Skill/Automation Run | Studio | Endpoint capability/evaluation != Tool/Run |
| DEP-EPT3-010 | privacy/tenant/least privilege | Security + Settings | cross-tenant denied; masking retained |
| DEP-EPT3-011 | runtime/language/portability | OPEN-017 | unresolved; no final engine/language |
| DEP-EPT3-012 | platform/source support | OPEN-008 | unresolved; supported/unsupported/unknown explicit |
| DEP-EPT3-013 | future Collection/Live Response | EPT-4 NOT STARTED | missing local data stops at boundary |
| DEP-EPT3-014 | future containment/response | EPT-5 NOT STARTED | no response action |
| DEP-EPT3-015 | update/resilience/security closure | EPT-6 NOT STARTED | detection-update source not implemented here |

Dependencies are functional/documentary and introduce no physical runtime, protocol or deployment topology.