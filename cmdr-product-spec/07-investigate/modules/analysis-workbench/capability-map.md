---
id: analysis-workbench-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
---
# Capability map — Analysis Workbench

## Static Analysis foundation
CAP-INV-301..313 remain canonical and unchanged.

## Dynamic Sandbox
CAP-INV-314..328 remain a separate module consumed through explicit handoffs.

## Reverse Engineering and Debugger
CAP-INV-329..346 remain canonical and unchanged.

## Memory Forensics
| ID | Capability | Family | Primary local concepts | Classes |
|---|---|---|---|---|
| CAP-INV-347 | Memory Forensics Intake and Preconditions | intake | Memory Forensics Intake; Memory Forensics Session relation | 0,1,2 |
| CAP-INV-348 | Memory Forensics Session Management | session | Memory Forensics Session; Session membership | 0,1,2 |
| CAP-INV-349 | Memory Image Integrity and Acquisition Context Review | integrity | Image Integrity Review; Memory Image relation | 0,2 |
| CAP-INV-350 | Platform and Analysis Profile Identification | profile | Platform/Profile Identification; Selected profile relation | 0,1,2 |
| CAP-INV-351 | Process and Thread Reconstruction | reconstruction | Process Observation; Thread Observation | 0,1,2 |
| CAP-INV-352 | Memory Region, Mapping and Protection Analysis | memory-layout | Memory Region; Mapping | 0,1,2 |
| CAP-INV-353 | Loaded Modules, Images and Driver Analysis | modules-drivers | Module Observation; Driver Observation | 0,1,2 |
| CAP-INV-354 | Handles, System Objects and IPC Analysis | objects-ipc | Handle Observation; System Object Observation | 0,1,2 |
| CAP-INV-355 | Network State and Connection Artifact Reconstruction | memory-network | Network State Observation; Connection candidate | 0,1,2 |
| CAP-INV-356 | Sensitive Material and Credential Exposure Assessment | sensitive-data | Sensitive Material Candidate; Access event | 0,1,2 |
| CAP-INV-357 | Injection, Hollowing and Memory Anomaly Analysis | anomaly | Memory Anomaly; Hypothesis relation | 0,1,2 |
| CAP-INV-358 | Kernel State and Rootkit Indicator Analysis | kernel | Kernel State Observation; Rootkit indicator candidate | 0,1,2 |
| CAP-INV-359 | Memory Timeline and Cross-Source Correlation | timeline | Memory Timeline projection; Correlation relation | 0,1,2 |
| CAP-INV-360 | Memory Artifact Extraction and Carving | extraction | Derived Artifact; Extraction Result | 0,1,2 |
| CAP-INV-361 | Memory Forensics Provenance and Reproducibility | provenance | Reproducibility Assessment; Provenance relation | 0,1,2 |
| CAP-INV-362 | Memory Forensics Handoff to Evidence, Findings and Detection Engineering | handoff | Evidence candidate package; Finding Draft | 0,1,2 |

The sixteen capabilities separate intake/session, source trust, profile, reconstruction domains, sensitive data, anomaly/kernel review, timeline, extraction, provenance and handoff. No Disk, Filesystem or full Network Forensics capability is included.

The Analysis Workbench now contains 47 capabilities: 13 Static, 18 Reverse/Debugger and 16 Memory Forensics. Dynamic Sandbox remains separate. No Disk, Filesystem or full Network Forensics capability is included.
