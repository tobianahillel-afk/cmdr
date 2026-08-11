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