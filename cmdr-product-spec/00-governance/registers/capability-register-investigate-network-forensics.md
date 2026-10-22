---
id: capability-register-investigate-network-forensics
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-05
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
---
# Capability Register — Investigate Network Forensics

All IDs are immutable, unique, `draft`, `defined` and `planned`. Investigate owns the analytical context; source, Studio, Settings, Shared, Command and Govern ownership remains unchanged.

| ID | Name | Owner product | Owner module | Status | Delivery status | Delivery mode | Canonical file | Roles | Primary objects or concepts | Consumers | Requirement IDs | OPEN | Dependencies | Supersedes | Date |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| CAP-INV-380 | Network Forensics Intake and Preconditions | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-forensics-intake-and-preconditions.md` | DFIR Analyst; Case Analyst | Network Intake Assessment; Capture Artifact relation | Case; Workbench; CAP-INV-381 | REQ-INV-001; REQ-PROD-014,020; REQ-AI-002; REQ-SEC-001 | 005,008,013,014,015 | CAP-INV-105,208,212..214; Studio; Settings; Shared | none | 2026-08-05 |
| CAP-INV-381 | Network Forensics Session Management | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-forensics-session-management.md` | Network Forensics Analyst; Investigation Lead | Network Forensics Session; session membership | Workbench; Case; comparison | REQ-INV-001; REQ-PROD-014,020; REQ-AI-002; REQ-SEC-001 | 005,008,013,014,015 | Case; Capture Artifact; Studio; Shared | none | 2026-08-05 |
| CAP-INV-382 | Capture Integrity, Scope and Acquisition Context Review | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/capture-integrity-scope-and-acquisition-context-review.md` | Evidence Reviewer; Network Analyst | Capture Integrity Review; acquisition context | Intake; Evidence; CAP-INV-383 | REQ-INV-001; REQ-PROD-014,020; REQ-SEC-001 | 005,008,013,014,015 | CAP-INV-203,208,212..214; custody | none | 2026-08-05 |
| CAP-INV-383 | Capture Coverage, Timebase, Interface and Direction Context | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/capture-coverage-timebase-interface-and-direction-context.md` | Network Analyst; Sensor Reviewer | Capture Coverage Assessment; Sensor Context | Packet inspection; timeline; comparison | REQ-INV-001; REQ-PROD-014,020; REQ-UX-006 | 005,008,013,015 | Settings health/timebase; CAP-INV-382 | none | 2026-08-05 |
| CAP-INV-384 | Packet and Frame Inspection | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/packet-and-frame-inspection.md` | Packet Analyst | Packet Observation; Frame Observation | Flow reconstruction; extraction | REQ-INV-001; REQ-PROD-014,020; REQ-SEC-001 | 005,008,013,015 | CAP-INV-383; Studio Tool; payload permissions | none | 2026-08-05 |
| CAP-INV-385 | Flow and Session Reconstruction | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/flow-and-session-reconstruction.md` | Network Forensics Analyst | Flow Observation; Network Session Candidate | Conversation; timeline; behavior | REQ-INV-001; REQ-PROD-014,020 | 005,008,013,015 | CAP-INV-384; Studio; Shared Timeline | none | 2026-08-05 |
| CAP-INV-386 | Protocol Identification and Conversation Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/protocol-identification-and-conversation-analysis.md` | Protocol Analyst | Protocol Candidate; Conversation Candidate | DNS; transactions; transfers | REQ-INV-001; REQ-PROD-014,020 | 005,008,013,015 | CAP-INV-385; Studio decoder Tools | none | 2026-08-05 |
| CAP-INV-387 | DNS and Name Resolution Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/dns-and-name-resolution-analysis.md` | Name Resolution Analyst | DNS Observation; indicator candidate | Entity; Hypothesis; Evidence handoff | REQ-INV-001; REQ-PROD-014,020 | 005,008,013,015 | CAP-INV-386; CAP-INV-104; Shared Entity | none | 2026-08-05 |
| CAP-INV-388 | Web and Application Transaction Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/web-and-application-transaction-analysis.md` | Application Traffic Analyst | Application Transaction | Hypothesis; extraction; timeline | REQ-INV-001; REQ-PROD-014,020; REQ-SEC-001 | 005,008,013,015 | CAP-INV-386; payload permissions | none | 2026-08-05 |
| CAP-INV-389 | Encrypted Traffic Metadata and Certificate Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/encrypted-traffic-metadata-and-certificate-analysis.md` | Encrypted Traffic Analyst | Encrypted Traffic Observation; Certificate Observation | Entity; Hypothesis; handoff | REQ-INV-001; REQ-PROD-014,020; REQ-SEC-001 | 005,008,013,015 | CAP-INV-386; Security permissions | none | 2026-08-05 |
| CAP-INV-390 | File, Object and Content Transfer Reconstruction | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/file-object-and-content-transfer-reconstruction.md` | Network Artifact Analyst | Transfer Observation; Reconstructed Object Candidate | CAP-INV-395; Static; Reverse | REQ-INV-001; REQ-PROD-014,020; REQ-OBJ-003 | 005,008,013,014,015 | CAP-INV-386/388; CAP-INV-311 | none | 2026-08-05 |
| CAP-INV-391 | Network Entity, Endpoint and Relationship Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-entity-endpoint-and-relationship-analysis.md` | Entity Analyst | Network Relationship; Entity Candidate | Case; Graph; timeline | REQ-INV-001; REQ-PROD-014,020 | 005,008,013,015 | CAP-INV-104; Shared Entity/Graph | none | 2026-08-05 |
| CAP-INV-392 | Network Behavior, Periodicity and Anomaly Analysis | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-behavior-periodicity-and-anomaly-analysis.md` | Network Behavior Analyst | Network Anomaly Candidate | Hypothesis; Evidence candidate | REQ-INV-001; REQ-PROD-014,020; REQ-AI-002 | 005,008,013,015 | CAP-INV-383,385,391; CAP-INV-321 boundary | none | 2026-08-05 |
| CAP-INV-393 | Network Timeline and Cross-Source Correlation | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-timeline-and-cross-source-correlation.md` | Timeline Analyst | Network Timeline Entry; Correlation Relation | Case; Memory; Disk; Event Search | REQ-INV-001; REQ-PROD-014,020; REQ-UX-006 | 005,008,013,015 | Shared Timeline; CAP-INV-002,355,374 | none | 2026-08-05 |
| CAP-INV-394 | Multi-Capture, Sensor and Time-Window Comparison | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/multi-capture-sensor-and-time-window-comparison.md` | Comparison Analyst | Network Comparison Result | Timeline; Finding Draft | REQ-INV-001; REQ-PROD-014,020; REQ-UX-006 | 005,008,013,015 | CAP-INV-381,383,393; Shared comparison | none | 2026-08-05 |
| CAP-INV-395 | Network Artifact Extraction and Derived Artifact Management | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-artifact-extraction-and-derived-artifact-management.md` | Network Artifact Analyst | Extraction Result; Derived Artifact | Static; Reverse; Sandbox | REQ-INV-001; REQ-PROD-014,020; REQ-OBJ-003 | 005,008,013,014,015 | CAP-INV-311; Studio; payload permissions | none | 2026-08-05 |
| CAP-INV-396 | Network Forensics Provenance and Reproducibility | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-forensics-provenance-and-reproducibility.md` | Audit Analyst; Quality Reviewer | Reproducibility Assessment; Provenance relation | Evidence; Reporting; QA | REQ-INV-001; REQ-PROD-014,020; REQ-AI-002; REQ-SEC-001 | 005,008,013,014,015 | Studio; Settings; Shared Trace/Activity/Versioning | none | 2026-08-05 |
| CAP-INV-397 | Network Forensics Handoff to Evidence, Findings, Detection Engineering and Intelligence | Investigate | Analysis Workbench | draft | defined | planned | `network-forensics/capabilities/network-forensics-handoff.md` | Investigation Lead; Evidence Reviewer | Evidence Candidate Package; Finding Draft; Future Handoff Package | CAP-INV-107/108/109; future 4B.3 | REQ-INV-001; REQ-PROD-014,020; REQ-OBJ-004 | 005,008,013,014,015 | CAP-INV-380..396; future Phase 4B.3 | none | 2026-08-05 |

## Measures
- capabilities: **18**;
- sections: **486**;
- mandatory tables: **108**;
- delivery: **18 defined / 18 planned**;
- duplicate or recycled IDs: **0**;
- CAP-INV-4xx/5xx: **0**.
