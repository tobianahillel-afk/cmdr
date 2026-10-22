---
id: investigate-analysis-workbench-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
---
# Analysis Workbench capabilities

| Family | Range | Count | Canonical location |
|---|---|---:|---|
| Workbench Foundation and Static Analysis | CAP-INV-301..313 | 13 | `capabilities/` |
| Reverse Engineering and Debugger | CAP-INV-329..346 | 18 | `reverse-engineering-and-debugger/` |
| Memory Forensics | CAP-INV-347..362 | 16 | `memory-forensics/` |
| Disk and Filesystem Forensics | CAP-INV-363..379 | 17 | `disk-and-filesystem-forensics/` |
| Network Forensics | CAP-INV-380..397 | 18 | `network-forensics/` |
| **Analysis Workbench total** |  | **82** |  |

Dynamic Sandbox CAP-INV-314..328 remains a separate Investigate module consumed by explicit handoffs. Across Phase 4B.2 analysis families, CAP-INV-301..397 account for 97 capabilities.

No Network capability duplicates Event Search, Dynamic Network Behavior, Memory Network State, Disk persistent network artifacts, Shared Entity/Graph/Timeline, Detection Engineering or Intelligence.
