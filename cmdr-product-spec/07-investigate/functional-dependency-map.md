---
id: investigate-functional-dependency-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
open_decisions:
  - OPEN-005
  - OPEN-007
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Functional dependency map — Investigate through Phase 4B.2A

The Phase 4B.1 capability dependencies remain canonical. Phase 4B.2A adds:

| Capability | Capability dependencies | Product / Shared dependencies | Object dependencies | OPEN / future | Failure behavior |
|---|---|---|---|---|---|
| CAP-INV-201 | 102,105,110,112 | Settings Fleet/Policy/Health; Agent; Linking/Context | Case, Endpoint, Agent, Policy | OPEN-008 | offline/stale/unsupported explicit; no execution |
| CAP-INV-202 | 201,203,108,113,215 | Agent capabilities; Settings Policy; Govern | Collection Request, Case, Endpoint | OPEN-007/008/013 | draft preserved; missing authority visible |
| CAP-INV-203 | 202,204..208,212,214 | Background Jobs, Notifications, Agent queue/results | Request, Job concept, Artifact | OPEN-008 | partial per item; no false global success |
| CAP-INV-204 | 202,203,205,206 | Agent collection profiles; Jobs | Request, Job, Artifact | OPEN-008/013 | category failures visible; targeted retry |
| CAP-INV-205 | 202,203,213 | Agent file collection; preview; Jobs | Request, Artifact | OPEN-008/013 | missing/locked/restricted explicit |
| CAP-INV-206 | 201,202,212,215 | Agent host/process/network inspection | Agent Command, Operation Result | OPEN-008/013 | read-only snapshot; no containment |
| CAP-INV-207 | 202,203,213 | Agent memory capability; Jobs | Memory Image, Artifact | OPEN-005/008/013 | unsupported explicit; no engine promise |
| CAP-INV-208 | 202,203,213 | Agent network capability; Jobs | Capture Artifact | OPEN-008/013 | loss/error visible; no engine/format |
| CAP-INV-209 | 201,210,211,214 | Agent availability; Console; Notifications | Live Session concept | OPEN-007/008/013/015 | disconnect ≠ success; transcript retained |
| CAP-INV-210 | 209,212,214,215 | Agent command catalogue; dangerous-action pattern | Endpoint Operation, Agent Command | OPEN-007/008/013/015 | class 3/4 blocked and routed to Govern |
| CAP-INV-211 | 209,212,213 | Agent transfer; Jobs | Artifact, Attachment relation | OPEN-008/013/014 | collision/partial/cancel explicit |
| CAP-INV-212 | 203,210,211,105,107 | Agent result reporting; Inspector | Operation Result, Artifact, Evidence | OPEN-013/015 | local result never becomes Govern Result |
| CAP-INV-213 | 203,205,207,208,212 | Export, retention, Trust | Artifact, custody/provenance concepts | OPEN-014 | missing links produce partial/disputed |
| CAP-INV-214 | 008,110,112,203,209,212,213 | Trace, Activity, Timeline, Studio, Govern | all run/result references | OPEN-007/015 | source owner preserved; broken link visible |
| CAP-INV-215 | 109,113,201,212,214 | Govern Action Center/Policy Gates/Runs; Agent containment | Finding, Evidence, Action Request | OPEN-007/008/013/015 | draft only; no containment execution |

## Rules
- dependency/projection never transfers ownership;
- source, freshness, permission and version accompany every projection;
- Shared mechanisms are consumed, never duplicated;
- no model is required for an essential workflow;
- `CAP-INV-3xx` and forensic engines remain Phase 4B.2B and are not defined here.