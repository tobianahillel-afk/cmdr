---
id: investigate-disk-filesystem-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
---
# Capability map — Disk and Filesystem Forensics

| ID | Capability | Family | Primary local concepts | Classes |
|---|---|---|---|---|
| CAP-INV-363 | Disk Forensics Intake and Preconditions | intake | Disk Forensics Intake; Disk Forensics Session relation | 0,1,2 |
| CAP-INV-364 | Disk Forensics Session Management | session | Disk Forensics Session; Session membership | 0,1,2 |
| CAP-INV-365 | Disk Image Integrity and Acquisition Context Review | integrity | Disk Image Integrity Review; Acquisition context relation | 0,2 |
| CAP-INV-366 | Partition, Volume and Filesystem Identification | structure | Partition / Volume / Filesystem Candidate | 0,1,2 |
| CAP-INV-367 | Filesystem Navigation and Metadata Inspection | navigation | Filesystem Entry; Bookmark / selection | 0,1,2 |
| CAP-INV-368 | File Identity, Content and Relationship Analysis | file-analysis | File Observation; Content identity relation | 0,1,2 |
| CAP-INV-369 | Deleted and Unallocated Data Analysis | deleted-unallocated | Deleted Entry Candidate / Unallocated Region | 0,1,2 |
| CAP-INV-370 | Filesystem Journal and Change History Analysis | journal | Journal Record; Change-history relation | 0,1,2 |
| CAP-INV-371 | Operating System and Configuration Artifact Analysis | system-config | System Artifact Observation | 0,1,2 |
| CAP-INV-372 | User Activity and Application Artifact Analysis | user-activity | User / Application Artifact | 0,1,2 |
| CAP-INV-373 | Startup, Persistence and Execution Artifact Analysis | persistence | Persistence Candidate | 0,1,2 |
| CAP-INV-374 | Disk Timeline and Temporal Reconstruction | timeline | Disk Timeline Entry | 0,1,2 |
| CAP-INV-375 | File Carving and Recovery | recovery | Recovery Result; Derived Artifact | 0,1,2 |
| CAP-INV-376 | Encrypted, Compressed and Restricted Content Handling | restricted-content | Restricted Content Record | 0,1,2 |
| CAP-INV-377 | Multi-Image, Volume and Snapshot Comparison | comparison | Disk Comparison Result | 0,1,2 |
| CAP-INV-378 | Disk Forensics Provenance and Reproducibility | provenance | Reproducibility Assessment | 0,1,2 |
| CAP-INV-379 | Disk Forensics Handoff to Evidence, Findings and Future Analysis | handoff | Evidence Candidate Package / Finding Draft | 0,1,2 |

No full Network Forensics capability is included.
