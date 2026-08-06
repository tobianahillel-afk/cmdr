---
id: investigate-threat-intelligence-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-INV-006
open_decisions:
  - OPEN-018
---
# Capability map — Threat Intelligence Foundations

| ID | Capability | Primary role | Local classes |
|---|---|---|---|
| CAP-INV-501 | Threat Intelligence Intake and Preconditions | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-502 | Intelligence Requirements and Collection Priority Management | Intelligence Manager | 0,1,2 |
| CAP-INV-503 | Threat Intelligence Workspace and Knowledge Project Management | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-504 | Intelligence Source Catalog and Access Context | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-505 | Source Reliability and Information Credibility Assessment | Intelligence Reviewer | 0,1,2 |
| CAP-INV-506 | Intelligence Material Intake, Parsing and Normalization | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-507 | Observable and Indicator Candidate Management | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-508 | Threat Entity and Identity Candidate Management | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-509 | Malware, Tool and Capability Knowledge Management | Malware Analyst | 0,1,2 |
| CAP-INV-510 | Infrastructure and Resource Knowledge Management | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-511 | Campaign, Activity Cluster and Intrusion Set Candidate Management | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-512 | Technique, Behavior and TTP Mapping | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-513 | Sighting, Occurrence and Observation Management | SOC Analyst | 0,1,2 |
| CAP-INV-514 | Intelligence Relationship Graph and Evidence-Backed Linking | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-515 | Intelligence Confidence, Assessment and Contradiction Management | Intelligence Reviewer | 0,1,2 |
| CAP-INV-516 | Knowledge Deduplication, Versioning and Supersession | Threat Intelligence Analyst | 0,1,2 |
| CAP-INV-517 | Intelligence Expiration, Revocation and Lifecycle Management | Intelligence Manager | 0,1,2 |
| CAP-INV-518 | Threat Intelligence Provenance and Analysis Handoff | Threat Intelligence Analyst | 0,1,2 |

The eighteen capabilities separate intake, requirements, workspace, source context, reliability/credibility, material handling, candidate knowledge, Sightings, relationships, confidence, versioning/lifecycle and provenance/handoff. No 4B.3B.2, watchlist, external sharing, Report, attribution or operational action is created.
