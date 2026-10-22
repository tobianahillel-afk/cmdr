# Endpoint EPT-2 — Object Consumption Map

| Concept | Canonical owner | Endpoint use | Operations | Gap | Future lot |
|---|---|---|---|---|---|
| endpoint-agent | Endpoint | source identity/context | read | none | — |
| Observation Source / Sensor Reference / Local Observation | Endpoint | produce local facts | create/read | physical schema intentionally absent | implementation |
| telemetry-event | Shared | project Endpoint observations into shared envelope | project/read | physical schema out of scope | Shared/implementation |
| Process/File/Network/User/System/Sensor observations | Endpoint | endpoint-specific fact semantics | create/read | platform/source coverage OPEN-008 | EPT-3 consumes |
| Normalized Observation / Classification | Endpoint + Shared boundary | endpoint mapping context | derive/read | final external standard unselected | Shared/implementation |
| Telemetry Gap / Duplicate Candidate / Freshness / Sampling / Backpressure | Endpoint | quality facts | derive/read | runtime implementation absent | implementation |
| Technical Capability Declaration / Dependency / Availability | Endpoint | detailed local technical declaration | derive/read | support decision OPEN-008 | Settings consumes |
| Endpoint Inventory / Health | Endpoint | EPT-1 foundation consumed | read | none | — |
| Endpoint Policy / Fleet / Provider Reference | Platform Settings | read administrative projection | read | no ownership transfer | Settings |
| Evidence / Finding / Case | Investigate | consumer-owned interpretation only | reference | no auto-promotion | EPT-3 |
| Tool / Automation Run | Studio | consumer refs only | reference | capability != Tool | Studio |
| Decision / Response Run / Result | Govern | technical context only | reference | observation != Result | EPT-5 |
| Trace / Activity / Job / Report / Export | Shared | generic mechanisms only | reference | no duplication | Shared |
