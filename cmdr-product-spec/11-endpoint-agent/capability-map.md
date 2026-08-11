---
id: endpoint-capability-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint Capability Map

## EPT-1 foundations
| Range | Family | Capabilities | Boundary |
|---|---|---:|---|
| CAP-EPT-001..003 | identity/enrollment/scope | 3 | local Agent state; Settings admin remains external |
| CAP-EPT-004..007 | platform/version/inventory/freshness | 4 | observed facts; no support or global CMDB claim |
| CAP-EPT-008..010 | health/heartbeat/operational state | 3 | technical state only; health != security posture |
| CAP-EPT-011..014 | capability/Fleet/Policy/provenance | 4 | projections/handoffs; no foreign ownership transfer |

Total EPT-1: **14 capabilities / 378 sections / 84 mandatory tables**.

## Future capability families
EPT-2 Telemetry; EPT-3 Detection/Investigation; EPT-4 Collection/Live Response; EPT-5 Containment/Verification; EPT-6 Updates/Resilience/Security. All remain **NOT STARTED**.

## Non-equivalence
Endpoint Agent != Automation Agent; Endpoint != Agent automatically; Device != Agent Instance; Host Identity != Agent Identity; inventory != Fleet; health != security posture; advertised != authorized/available; Endpoint Policy != effective local state; technical output != Govern Result; execution != Response Run/Tool Call.

---

## EPT-2 telemetry/observation capability layer — build
The preceding EPT-2 NOT STARTED statement is historical pre-EPT-2 evidence.

| Range | Family | Capabilities | Boundary |
|---|---|---:|---|
| CAP-EPT-015..016 | source/observation and Shared event projection | 2 | Endpoint local facts; Shared owns generic event envelope |
| CAP-EPT-017..022 | process/file/network/auth/system/sensor observations | 6 | technical facts only; no Detection/Finding/Evidence |
| CAP-EPT-023..026 | normalization/quality/rate/privacy | 4 | no final schema, storage, event bus or RBAC |
| CAP-EPT-027..030 | declaration/availability/handoff/provenance | 4 | CAP-EPT-011 summary preserved; no Tool/Result ownership transfer |

EPT-2: **16 capabilities / 432 sections / 96 mandatory tables**. Endpoint cumulative: **30 / 810 / 180**. `OPEN-008` remains open. Endpoint Screen IDs remain 0. EPT-3..EPT-6 remain NOT STARTED.

---

## EPT-3 local detection/investigation capability layer — build
The preceding EPT-3 NOT STARTED statement is historical pre-EPT-3 evidence.

| Range | Family | Capabilities | Boundary |
|---|---|---:|---|
| CAP-EPT-031..033 | Detection Content consumption, local evaluation/match, local signal-candidate | 3 | Investigate owns content; Command owns canonical Detection/Signal |
| CAP-EPT-034..036 | context/severity/confidence, grouping/suppression projection, coverage/health/gaps | 3 | no impact/Finding/response promotion |
| CAP-EPT-037..041 | process/file/network/user-session/system local investigation | 5 | existing data only; no acquisition |
| CAP-EPT-042..044 | local timeline/correlation and pivots/context expansion | 3 | Shared generic mechanisms retained; pivot != Collection |
| CAP-EPT-045..046 | summary/handoff and provenance/cross-product contracts | 2 | summary != Finding/Evidence/Case/Result |

EPT-3: **16 capabilities / 432 sections / 96 mandatory tables**. Endpoint cumulative: **46 / 1242 / 276**. `OPEN-008` and `OPEN-017` remain open. Endpoint Screen IDs remain 0. EPT-4..EPT-6 remain NOT STARTED.