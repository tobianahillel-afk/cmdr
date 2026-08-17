---
id: investigate-mobile-forensics-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-020, REQ-INV-001]
open_decisions: [OPEN-011, OPEN-014]
---
# Mobile Forensics Concepts

These are functional concepts and projections, not final canonical object schemas.

| Concept | Functional meaning | Required distinction |
|---|---|---|
| Mobile Device Candidate | scoped technical device reference | identifier ≠ person; declared owner ≠ effective user |
| Mobile Platform Candidate | declared/inferred platform and version with source/confidence | candidate ≠ selected supported platform |
| Mobile Acquisition Request | Mobile-specific need prepared for Collection | request ≠ Collection Job |
| Mobile Evidence Package | bounded received evidence representation plus custody/limitations | package ≠ device, backup, extraction, Session or Evidence |
| Device Backup | backup representation with source/version/time | backup ≠ filesystem/logical extraction; timestamp ≠ current state |
| Logical Extraction | declared logical representation | extraction ≠ original device or complete state |
| Filesystem Extraction | declared filesystem representation | extraction ≠ certain physical image |
| Mobile Forensics Session | durable analytical context | Session ≠ Evidence |
| Integrity / Completeness / Accessibility Assessments | separate judgments on trust, coverage and access | integrity ≠ completeness ≠ relevance ≠ usability |
| Mobile Filesystem Observation | partition/filesystem/storage/entry observation | observation ≠ current device state |
| Application Observation | package/version/install metadata candidate | installed ≠ used; package present ≠ active behavior |
| Application Data Observation | sourced application record | data ≠ user intent |
| Communication / Call / Contact Observation | stored communication artifacts | stored message ≠ author; call record ≠ conversation; contact ≠ relationship |
| Media / Document Observation | content/metadata relation | file ≠ user-created; metadata ≠ certain truth |
| Location / Sensor Observation | device-origin location/activity record | device location ≠ user presence; sensor ≠ human action |
| Account Observation | account identifier/context | account ≠ person |
| Sensitive Material Candidate | token/key/certificate/session/secret candidate | candidate ≠ valid/usable; presence ≠ reveal/copy/export/use permission |
| Connectivity Observation | network/Wi-Fi/Bluetooth/NFC/interface relation | record ≠ successful/malicious interaction; not full Network Forensics |
| SIM/eSIM Observation | SIM/eSIM/carrier/number metadata | SIM/number ≠ person |
| Paired Device Observation | sourced pairing/association record | paired device ≠ same owner |
| Synchronization Observation | local/cloud/cross-device sync relation | synchronized/cloud-backed ≠ certain local state or full Cloud Analysis |
| Deleted Entry Candidate | deletion/residual record candidate | candidate ≠ user deletion intent |
| Recovery Result | bounded recovered/carved fragment | recovered ≠ complete original; carved ≠ certain attribution |
| Mobile Anomaly | unusual sourced observation | anomaly ≠ compromise |
| Mobile Hypothesis | reviewable explanation with support/contradiction | persistence candidate ≠ confirmed persistence; suspicious app ≠ malware |
| Mobile Timeline | mobile-local chronology semantics | Mobile Timeline ≠ Case Timeline |
| Derived Artifact | bounded derived material with lineage | Derived Artifact ≠ Evidence |
| Evidence Candidate Package | handoff material for Evidence owner workflow | candidate ≠ qualified Evidence |
| Finding Draft | sourced draft for Finding workflow | draft ≠ confirmed Finding |
| Reproducibility Assessment | conditions and gaps for analytical replay/review | reproducible reasoning ≠ guaranteed technical replay |

All relations preserve tenant, scope, source, permissions, classification, timestamps, versions, partiality, errors and return origin. No relation transfers ownership or grants source access.
