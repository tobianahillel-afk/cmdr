---
id: investigate-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
---
# Capability map — Investigate through Phase 4B.2A

| ID | Capability | Module | Owner | Status | Delivery status | Delivery mode | OPEN |
|---|---|---|---|---|---|---|---|
| CAP-INV-001 | Signal Triage | signals-and-hunt | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-002 | Event Search | signals-and-hunt | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-003 | Query Authoring and Assistance | signals-and-hunt | Investigate Product Lead | draft | defined | planned | OPEN-015 |
| CAP-INV-004 | Event Inspection and Pivot | signals-and-hunt | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-005 | Hunt Management | signals-and-hunt | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-006 | Saved Searches and Query Assets | signals-and-hunt | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-007 | Search Result Organization | signals-and-hunt | Investigate Product Lead | draft | defined | planned | OPEN-013 |
| CAP-INV-008 | Search and Hunt Provenance | signals-and-hunt | Investigate Product Lead | draft | defined | planned | OPEN-015 |
| CAP-INV-101 | Case Queue | cases-and-evidence | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-102 | Case Lifecycle and Coordination | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013 |
| CAP-INV-103 | Hypothesis Management | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013, OPEN-015 |
| CAP-INV-104 | Entity and Relationship Management | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013 |
| CAP-INV-105 | Artifact Management | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-014 |
| CAP-INV-106 | Attachment Handling | cases-and-evidence | Investigate Product Lead | draft | proposed | planned | OPEN-014 |
| CAP-INV-107 | Evidence Creation and Management | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013 |
| CAP-INV-108 | Evidence Review and Qualification | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013 |
| CAP-INV-109 | Finding Management | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013, OPEN-015 |
| CAP-INV-110 | Investigation Timeline | cases-and-evidence | Investigate Product Lead | draft | defined | planned | — |
| CAP-INV-111 | Case Collaboration and Investigation Notes | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013, OPEN-014 |
| CAP-INV-112 | Case Replay and Investigation Review | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-013, OPEN-015 |
| CAP-INV-113 | Action Request Preparation | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-007, OPEN-013, OPEN-015 |
| CAP-INV-114 | Case Reporting Preparation | cases-and-evidence | Investigate Product Lead | draft | defined | planned | OPEN-014 |
| CAP-INV-201 | Endpoint Investigation Context | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008 |
| CAP-INV-202 | Collection Scope and Request Preparation | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-203 | Collection Job Management | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008 |
| CAP-INV-204 | Triage Package Collection | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-205 | File and Directory Acquisition | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-206 | Process and System Inspection | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-207 | Memory Acquisition Request | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-005, OPEN-008, OPEN-013 |
| CAP-INV-208 | Network Capture Request | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-209 | Live Session Management | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-007, OPEN-008, OPEN-013, OPEN-015 |
| CAP-INV-210 | Interactive Endpoint Operations | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-007, OPEN-008, OPEN-013, OPEN-015 |
| CAP-INV-211 | Session File Transfer | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-008, OPEN-013 |
| CAP-INV-212 | Endpoint Operation Result Handling | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-013, OPEN-015 |
| CAP-INV-213 | Collection Integrity and Custody | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-014 |
| CAP-INV-214 | Collection and Live Response Provenance | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-007, OPEN-015 |
| CAP-INV-215 | Containment Request Preparation | collection-and-live-response | Investigate Product Lead | draft | defined | planned | OPEN-007, OPEN-008, OPEN-013, OPEN-015 |

## Comptage
- Signals and Hunt : **8**.
- Cases and Evidence : **14**.
- Collection and Live Response : **15**.
- Total Investigate : **37 capabilities**.
- Delivery status : **36 defined**, **1 proposed** (`CAP-INV-106`).
- Delivery mode : **37 planned** ; aucune capability n’est native, integrated ou implemented.

## Plages réservées
| Plage | Sous-phase | État |
|---|---|---|
| CAP-INV-301..399 | Phase 4B.2B Analysis Workbench | non commencée ; aucun ID attribué |
| CAP-INV-401..499 | Phase 4B.3 Detection Engineering | non commencée ; aucun ID attribué |
| CAP-INV-501..599 | Phase 4B.3 Intelligence | non commencée ; aucun ID attribué |

`CAP-INV-215` spécialise la préparation Endpoint et remet le lifecycle à `CAP-INV-113`/Govern ; il ne crée pas de second Action Request.