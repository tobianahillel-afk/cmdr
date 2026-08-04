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

# Screen capability map — Investigate

| Écran actif | Capability principale | Capabilities secondaires | Objets | Défauts connus | Phase de réécriture |
|---|---|---|---|---|---|
| INV-TRI-001 Triage Desk | CAP-INV-001 | CAP-INV-102 | Signal, Incident, Case | template-level, objectif ancien intake | pilote/Phase 6 |
| INV-EVS-001 Event Search | CAP-INV-002 | CAP-INV-003,004,006,007,008 | Event, Query, Search Job, Case | pas d’Event Inspector/Hunt détaillé | pilote/Phase 6 |
| INV-CAS-001 Case Workspace | CAP-INV-102 | CAP-INV-103..114 | Case, Hypothesis, Artifact, Evidence, Finding | sections génériques | pilote/Phase 6 |
| INV-EVD-001 Evidence Board | CAP-INV-107 | CAP-INV-105,108,109 | Artifact, Evidence, Finding | qualification dimensions absente | rollout/Phase 6 |
| INV-HYP-001 Hypotheses & Findings | CAP-INV-103 | CAP-INV-108,109,113 | Hypothesis, Evidence, Finding | review/AI boundaries génériques | rollout/Phase 6 |
| INV-ENT-001 Entity Graph | CAP-INV-104 | CAP-INV-002,107 | Entity, Event, Evidence | merge/temporalité génériques | rollout/Phase 6 |
| INV-TIM-001 Case Timeline | CAP-INV-110 | CAP-INV-008,112 | Timeline Entry, Event, Evidence, Finding | contenu métier incomplet | rollout/Phase 6 |
| INV-RPT-001 Investigation Report | CAP-INV-114 | CAP-INV-109,112 | Case, Evidence, Finding | Reporting Engine boundary à appliquer | rollout/Phase 6 |
| INV-STA-001 Static Analysis | future CAP-INV-3xx | CAP-INV-105,107 | Artifact, Evidence | 4B.2 non traité | 4B.2 puis Phase 6 |
| INV-SBX-001 Dynamic Sandbox | future CAP-INV-3xx | CAP-INV-105,107 | Artifact, Evidence | 4B.2 non traité | 4B.2 puis Phase 6 |
| INV-REV-001 Reverse Engineering | future CAP-INV-3xx | CAP-INV-105,107,109 | Artifact, Evidence, Finding | 4B.2 non traité | 4B.2 puis Phase 6 |
| INV-DBG-001 Debugger | future CAP-INV-3xx | CAP-INV-105,107 | Artifact, Evidence | 4B.2 non traité | 4B.2 puis Phase 6 |
| INV-MEM-001 Memory Forensics | future CAP-INV-3xx | CAP-INV-105,107 | Memory Image, Artifact, Evidence | 4B.2 non traité | 4B.2 puis Phase 6 |
| INV-DSK-001 Disk & Artifact Forensics | future CAP-INV-3xx | CAP-INV-105,107,110 | Disk Image, Artifact, Evidence | 4B.2 non traité | 4B.2 puis Phase 6 |

Aucun écran n’est modifié ni réécrit en Phase 4B.1. Une capability ne crée pas automatiquement un écran.
