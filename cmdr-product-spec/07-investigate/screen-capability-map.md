---
id: investigate-screen-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-010
---
# Screen capability map — Investigate through Phase 4B.2A

| Écran actif | Capability principale | Capabilities secondaires | Objets | Défauts connus | Phase de réécriture |
|---|---|---|---|---|---|
| INV-TRI-001 Triage Desk | CAP-INV-001 | CAP-INV-102 | Signal, Incident, Case | template-level | Phase 6 |
| INV-EVS-001 Event Search | CAP-INV-002 | 003,004,006,007,008 | Event, Query, Search Job, Case | hunt/inspector incomplete | Phase 6 |
| INV-CAS-001 Case Workspace | CAP-INV-102 | 103..114, 201..215 as contextual actions | Case, Endpoint, Requests, Jobs, Sessions, Artifacts | no Endpoint/collection/live-response composition yet | Phase 6 |
| INV-EVD-001 Evidence Board | CAP-INV-107 | 105,108,109,203,212,213 | Artifact, Evidence, Operation Result, custody | collection/result handoff absent | Phase 6 |
| INV-HYP-001 Hypotheses & Findings | CAP-INV-103 | 108,109,113,215 | Hypothesis, Evidence, Finding, Action Request | endpoint containment package absent | Phase 6 |
| INV-ENT-001 Entity Graph | CAP-INV-104 | 002,107 | Entity, Event, Evidence | merge/temporal generic | Phase 6 |
| INV-TIM-001 Case Timeline | CAP-INV-110 | 008,112,203,209,212,214 | Timeline Entry, Jobs, Sessions, Results | collection/live-response events not composed | Phase 6 |
| INV-RPT-001 Investigation Report | CAP-INV-114 | 109,112,213,214 | Case, Evidence, Finding, provenance | reporting boundary generic | Phase 6 |
| SET-EAF-001 Endpoint Agent Fleet | Settings Fleet capability | CAP-INV-201 as transition only | Fleet, Endpoint Agent, Policy | Investigation deep link not specified | Settings/Phase 6 |
| SET-EP-001 Endpoint Policies | Settings Policy capability | CAP-INV-201/202 projection only | Policy, Fleet | no local mutation allowed | Settings/Phase 6 |
| SET-HLT-001 System Health | Settings Health capability | CAP-INV-201 projection only | Health, Agent/Fleet projection | freshness handoff generic | Settings/Phase 6 |
| GOV-ACT-001 Action Center | Govern | CAP-INV-113,215 producer context | Action Request, Decision | endpoint package consumer not detailed | Govern/Phase 6 |
| GOV-RUN-001 Runs and Rollback | Govern | CAP-INV-212,214 projection only | Decision, Response Run, Result | operation-result distinction to apply | Govern/Phase 6 |
| INV-STA-001 Static Analysis | future CAP-INV-3xx | CAP-INV-105,107,207 handoff only | Artifact, Memory Image, Evidence | Phase 4B.2B not treated | 4B.2B then Phase 6 |
| INV-SBX-001 Dynamic Sandbox | future CAP-INV-3xx | CAP-INV-105,107 | Artifact, Evidence | Phase 4B.2B not treated | 4B.2B then Phase 6 |
| INV-REV-001 Reverse Engineering | future CAP-INV-3xx | CAP-INV-105,107,109 | Artifact, Evidence, Finding | Phase 4B.2B not treated | 4B.2B then Phase 6 |
| INV-DBG-001 Debugger | future CAP-INV-3xx | CAP-INV-105,107 | Artifact, Evidence | Phase 4B.2B not treated | 4B.2B then Phase 6 |
| INV-MEM-001 Memory Forensics | future CAP-INV-3xx | CAP-INV-105,107,207 | Memory Image, Artifact, Evidence | only acquisition handoff defined | 4B.2B then Phase 6 |
| INV-DSK-001 Disk & Artifact Forensics | future CAP-INV-3xx | CAP-INV-105,107 | Disk Image, Artifact, Evidence | Phase 4B.2B not treated | 4B.2B then Phase 6 |

No detailed screen is modified or created. Endpoint Search, Collection, Live Session, Jobs and Terminal screens do not currently have canonical Screen IDs; Phase 4B.2A defines capabilities only.