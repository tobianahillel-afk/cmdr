---
id: endpoint-ept3-object-consumption-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Object / Concept Consumption Map

| Concept | Canonical owner | Endpoint usage | Local operation | Boundary/gap | Future lot |
|---|---|---|---|---|---|
| Detection Content / Version | Investigate Detection Engineering | consume/reference | eligibility/evaluation input | no authoring/runtime format | — |
| Detection Eligibility | Endpoint Agent | own technical state | derive | eligible ≠ enabled | — |
| Local Detection Evaluation / Match | Endpoint Agent | own technical state | derive | not canonical Command Detection | — |
| Local Detection Signal Candidate | Endpoint Agent | own technical state | derive/group/forward | not canonical Command Signal | — |
| canonical Detection / Signal / Incident | Command | destination projection | read/handoff only | no local lifecycle takeover | — |
| Detection Context/Coverage/Gap | Endpoint Agent | own technical state | derive | not Finding/impact | — |
| Process/File/Network/User/System Context | Endpoint Agent | local contextual projection | derive/read | no acquisition | EPT-4 if new data needed |
| Endpoint Contextual Timeline | Endpoint Agent | local projection | derive | Shared owns generic Timeline | — |
| telemetry-event / Timeline/Linking/Search/Trace | Shared | consume generic refs | read/link | no owner transfer | — |
| Evidence / Finding / Case | Investigate | consumer destination | read/link only | Endpoint creates none | Investigate workflow |
| Tool / Automation Run / Skill | Studio | provenance/reference only | none | Endpoint capability ≠ Tool | — |
| Decision / Approval / Response Run / Result | Govern | context/handoff only | none | no response authority | EPT-5/Govern |
| Endpoint Policy / assignments / configured runtime | Platform Settings | consume projection | read only | no admin takeover | EPT-6/update where relevant |
| local-audit-event | Endpoint Agent | local audit | append/reference | distinct from Shared Trace | EPT-6 closure |

No physical schema is introduced.