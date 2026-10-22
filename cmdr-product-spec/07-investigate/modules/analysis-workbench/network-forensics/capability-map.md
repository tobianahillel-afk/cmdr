---
id: investigate-network-forensics-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-012
  - REQ-PROD-014
---
# Capability map — Network Forensics

| ID | Capability | Primary role | Primary local concepts | Classes |
|---|---|---|---|---|
| CAP-INV-380 | Network Forensics Intake and Preconditions | DFIR Analyst | Network Intake Assessment | 0,1,2 |
| CAP-INV-381 | Network Forensics Session Management | Network Forensics Analyst | Network Forensics Session | 0,1,2 |
| CAP-INV-382 | Capture Integrity, Scope and Acquisition Context Review | Evidence Reviewer | Capture Integrity Review | 0,1,2 |
| CAP-INV-383 | Capture Coverage, Timebase, Interface and Direction Context | Network Forensics Analyst | Capture Coverage Assessment | 0,1,2 |
| CAP-INV-384 | Packet and Frame Inspection | Packet Analyst | Packet Observation and Frame Observation | 0,1,2 |
| CAP-INV-385 | Flow and Session Reconstruction | Network Forensics Analyst | Flow Observation and Network Session Candidate | 0,1,2 |
| CAP-INV-386 | Protocol Identification and Conversation Analysis | Protocol Analyst | Protocol Candidate and Conversation Candidate | 0,1,2 |
| CAP-INV-387 | DNS and Name Resolution Analysis | Name Resolution Analyst | DNS Observation | 0,1,2 |
| CAP-INV-388 | Web and Application Transaction Analysis | Application Traffic Analyst | Application Transaction | 0,1,2 |
| CAP-INV-389 | Encrypted Traffic Metadata and Certificate Analysis | Encrypted Traffic Analyst | Encrypted Traffic Observation and Certificate Observation | 0,1,2 |
| CAP-INV-390 | File, Object and Content Transfer Reconstruction | Network Artifact Analyst | Transfer Observation and Reconstructed Object Candidate | 0,1,2 |
| CAP-INV-391 | Network Entity, Endpoint and Relationship Analysis | Entity Analyst | Network Relationship and Entity Candidate | 0,1,2 |
| CAP-INV-392 | Network Behavior, Periodicity and Anomaly Analysis | Network Behavior Analyst | Network Anomaly Candidate | 0,1,2 |
| CAP-INV-393 | Network Timeline and Cross-Source Correlation | Timeline Analyst | Network Timeline Entry and Correlation Relation | 0,1,2 |
| CAP-INV-394 | Multi-Capture, Sensor and Time-Window Comparison | Comparison Analyst | Network Comparison Result | 0,1,2 |
| CAP-INV-395 | Network Artifact Extraction and Derived Artifact Management | Network Artifact Analyst | Derived Artifact and Extraction Result | 0,1,2 |
| CAP-INV-396 | Network Forensics Provenance and Reproducibility | Audit Analyst | Reproducibility Assessment and Provenance Record | 0,1,2 |
| CAP-INV-397 | Network Forensics Handoff to Evidence, Findings, Detection Engineering and Intelligence | Investigation Lead | Evidence Candidate Package, Finding Draft and Future Handoff Package | 0,1,2 |

The eighteen capabilities separate intake/session, source trust, coverage/timebase, packet inspection, reconstruction, protocol-specific analysis, entities, behavior, timeline, comparison, extraction, provenance and owner handoffs. No Detection Engineering or Intelligence capability is created.
