---
id: endpoint-object-consumption-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint Object Consumption Map — EPT-1

Functional map only; it defines no physical schema.

| Concept | Canonical owner | Endpoint usage | Local operations | Gap | Future lot |
|---|---|---|---|---|---|
| Endpoint Agent | Endpoint Agent | local instance identity/state | read/update local state | canonical object exists | EPT-1 |
| Endpoint | shared model, Settings-administered | target/association ref | read/link | Phase-7 canonical file deferred | future model work |
| Device | no canonical owner established | contextual target term | reference only | no canonical object | future object work |
| Host Identity | no standalone canonical owner | observed host context | observe/reference | no canonical object | future object work |
| Agent Identity | Endpoint Agent | instance identity semantics | derive/validate | represented by endpoint-agent | EPT-1 |
| Registration State | Endpoint Agent | local lifecycle fact | derive/transition | no standalone object | EPT-1 |
| Enrollment State | Endpoint Agent local / Settings admin source | local effect of admin handoff | derive/transition | admin record remains Settings | EPT-1 |
| Tenant Reference | Platform Settings | scope | read/validate | none | EPT-1 |
| Environment Reference | Platform Settings | scope | read/validate | none | EPT-1 |
| Platform Fact | Endpoint Agent | observed fact | observe/normalize | support decision OPEN-008 | EPT-1 |
| OS Fact | Endpoint Agent | observed fact | observe/normalize | support decision OPEN-008 | EPT-1 |
| Architecture Fact | Endpoint Agent | observed fact | observe/normalize | support decision OPEN-008 | EPT-1 |
| Agent Version | Endpoint Agent | observed fact | observe/compare | approval/support external | EPT-1 |
| Agent Build | Endpoint Agent | observed fact | observe/compare | trust/approval external | EPT-1 |
| Inventory Snapshot | Endpoint Agent | point-in-time local record | create/read/supersede | no global CMDB object | EPT-1 |
| Inventory Change | Endpoint Agent | derived comparison | derive/read | not Detection | EPT-1 |
| Agent Health | Endpoint Agent | technical health assessment | derive/read | no standalone object | EPT-1 |
| Self-Check Assessment | Endpoint Agent | bounded health evidence | derive/read | not integrity proof | EPT-1 |
| Heartbeat | Endpoint Agent | observation/fact | read/derive state | no wire format | EPT-1 |
| Connectivity State | Endpoint Agent | derived local fact | derive/read | no transport design | EPT-1 |
| Last Seen | Endpoint Agent | timestamp fact | read/derive freshness | distinct from heartbeat | EPT-1 |
| Freshness State | Endpoint Agent | derived source age | derive/read | thresholds source-owned | EPT-1 |
| Endpoint Operational State | Endpoint Agent | online/degraded/offline/stale etc. | derive/transition | no recovery engine | EPT-1 |
| Technical Capability Advertisement | Endpoint Agent | local technical metadata | derive/read | not Studio Tool | EPT-1 |
| Capability Availability | Endpoint Agent | local availability state | derive/read | not authorization/support | EPT-1 |
| Endpoint Agent Fleet | Platform Settings | aggregation destination/context | read/project only | no Endpoint admin | EPT-1 boundary |
| Endpoint Policy | Platform Settings | assigned policy ref | read projection | no Endpoint Policy mutation | EPT-1 boundary |
| Effective Policy State | Endpoint Agent | local observation | derive/read | no standalone canonical object | EPT-1 |
| Local Audit Event | Endpoint Agent | provenance event | create/read | local audit != Shared Trace | EPT-1 |
| Secret Reference | Platform Settings | metadata/reference only | reference, never reveal value | raw secret external | future execution |
| Provider Reference | Platform Settings | context only | reference | provider choice unresolved | future |
| Tool | Studio | distinct orchestration/execution concept | reference only | Tool canonical file deferred Phase 7 | Studio |
| Automation Agent | Studio | distinct agentic concept | reference only | canonical object exists | Studio |
| Automation Run | Studio | distinct run concept | reference only | Phase-7 object deferred | Studio |
| Decision | Govern | authority context | read/reference | no local Decision | Govern |
| Response Run | Govern | governed execution context | read/reference | no local Response Run | Govern |
| Result | Govern | canonical response outcome | read/reference | technical output remains distinct | Govern |
| Job | Shared generic mechanisms | generic background context | consume | generic object not finalized | Shared |
| Trace | Shared | generic trace context | link/reference | local audit distinct | Shared |
| Activity | Shared | generic activity context | link/reference | local audit distinct | Shared |

No row transfers ownership or grants permission.