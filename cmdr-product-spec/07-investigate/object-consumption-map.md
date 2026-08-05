---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-014
  - REQ-PROD-061
  - REQ-PROD-062
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Object consumption map — Investigate through Phase 4B.2 closure

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | context and return origin | read/link | final cardinality | Objects |
| Hypothesis | Investigate | analytical question and relation | read/link/prepare | final confidence model | Objects |
| Incident | Command | operational context | read/link only | cross-product contract | Command/Objects |
| Endpoint | shared / Endpoint Agent | acquisition source or correlation | read projection | stable identity and support | Objects/Endpoint |
| Collection Request | Investigate concept | requested network acquisition context | read/link only | final object schema | Objects |
| Collection Job | Investigate concept / execution producer | status, errors, window and result | read/link only | canonical object absent | Objects |
| Network Capture Artifact | Artifact projection under Investigate | immutable analytical source | read/link/review; never mutate | canonical subtype/status unresolved | Objects |
| Artifact | Investigate | source and lineage | read/link/version | OPEN-014 | Objects |
| Derived Artifact | Investigate concept | extracted network content | create/read/export/withdraw | object absent | Objects |
| Evidence | Investigate | qualified destination | candidate package only | qualification workflow remains CAP-INV-107/108 | Objects/Trust |
| Finding | Investigate | reasoned conclusion destination | prepare draft only | final state model | Objects |
| Entity | Shared Capabilities | candidate identity and relation | read/link/propose; no silent merge | resolution contract future | Shared/Objects |
| Network Forensics Session | Investigate concept | durable analytical context | create/update/pause/close/reopen/supersede | canonical object absent | Objects |
| Capture Integrity Review | Investigate concept | source integrity, scope and custody interpretation | create/review/dispute/supersede | object absent | Objects/Trust |
| Capture Coverage Assessment | Investigate concept | coverage, loss, truncation, direction and timebase | create/annotate/compare | object absent | Objects |
| Sensor Context | Platform Settings projection / producer | interface, point of observation, limitations and health | read/link only | sensor model and support OPEN-008 | Settings/Objects |
| Packet Observation | Investigate concept | captured unit interpretation | read/annotate/select | object absent; packet ≠ frame | Objects |
| Frame Observation | Investigate concept | link-layer projection when available | read/annotate/select | object absent | Objects |
| Flow Observation | Investigate concept | packet grouping candidate | create/read/compare/dispute | reconstruction contract absent | Objects |
| Network Session Candidate | Investigate concept | reconstructed transport/session candidate | create/read/compare/dispute | object absent; not application certainty | Objects |
| Conversation Candidate | Investigate concept | ordered message relationship candidate | create/read/annotate | object absent | Objects |
| Protocol Candidate | Investigate concept | protocol interpretation with confidence | propose/select/confirm-by-analyst/dispute | final support model absent | Objects |
| DNS Observation | Investigate concept | query/response/name-resolution interpretation | create/read/link/annotate | object absent; not confirmed IOC | Objects |
| Application Transaction | Investigate concept | request/response or operation candidate | create/read/compare/annotate | object absent | Objects |
| Certificate Observation | Investigate concept | certificate and chain metadata | create/read/link/annotate | identity semantics absent | Objects/Trust |
| Encrypted Traffic Observation | Investigate concept | handshake and authorized metadata | create/read/compare | payload access model future | Objects/Permissions |
| Transfer Observation | Investigate concept | file/object/content transfer candidate | create/read/compare | reassembly model absent | Objects |
| Network Relationship | Investigate relation using Shared Entity/Graph | analytical relation between observations and entities | propose/link/dispute/supersede | relationship schema absent | Shared/Objects |
| Network Anomaly | Investigate candidate concept | periodicity, scan, lateral or exfiltration candidate | propose/review/dispute/withdraw | not a Finding | Objects |
| Network Timeline Entry | Investigate projection using Shared Timeline | temporal reconstruction with timestamp quality | create/read/correlate/annotate | not canonical Timeline object | Shared/Objects |
| Network Comparison Result | Investigate concept | multi-capture/sensor/window differences | create/read/export/link | object absent | Objects |
| Reproducibility Assessment | Investigate concept | reproducibility state and missing context | create/review/dispute/supersede | object absent | Objects/Trust |
| Provenance Record | Shared records plus Investigate relation | source chain and disposition | read/link/emit events | final contracts absent | Shared/Trust |
| Tool | CMDR Studio | compatible processing selection | select/read only | final object/version contract absent | Studio/Objects |
| Tool Call | CMDR Studio | execution, parameters, status and output | invoke/read/link according to permission | final contract absent | Studio/Objects |
| Automation Run | CMDR Studio | optional orchestration and attribution | request/read/link | OPEN-015 | Studio/Objects |
| Fleet / Policy / Provider / Storage / Retention / Health / Time synchronization | Platform Settings | administrative projections | read only | support and SLOs unresolved | Settings/Objects |
| Background Job / Notification / Trace / Activity / Timeline / Graph / Search / Export / Reporting / Recovery | Shared mechanisms | progress, navigation, correlation, provenance and output | consume/emit semantics | contracts future | Shared/Trust |
| Action Request / Decision / Approval / Response Run / Result | Govern | real-target authority and returned outcomes | read/link/prepare only | future contracts | Govern |
| Attachment | open | documentary content reference | reference only | OPEN-014 | Objects |

No schema, JSON Schema, final cardinality, object state machine, capture format, packet field model, protocol support list or atomic permission is defined. Network Capture, Network Forensics Session, Network Flow and Network Observation remain functional concepts or projections until the Objects phase decides their canonical status.
