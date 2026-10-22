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
# Capability map — Analysis Workbench through Phase 4B.2

## Canonical families
| Range | Family | Count | Relationship |
|---|---|---:|---|
| CAP-INV-301..313 | Workbench Foundation and Static Analysis | 13 | canonical Workbench foundation |
| CAP-INV-314..328 | Dynamic Sandbox and Behavioral Analysis | 15 | separate module consumed by handoff |
| CAP-INV-329..346 | Reverse Engineering and Debugger | 18 | canonical Workbench family |
| CAP-INV-347..362 | Memory Forensics | 16 | canonical Workbench family |
| CAP-INV-363..379 | Disk and Filesystem Forensics | 17 | canonical Workbench family |
| CAP-INV-380..397 | Network Forensics | 18 | canonical Workbench family |

The Analysis Workbench contains **82** canonical capabilities excluding the separate Dynamic Sandbox module, and **97 CAP-INV-3xx** capabilities when Dynamic Sandbox is included in Phase 4B.2 analysis coverage.

## Network Forensics
CAP-INV-380..397 separate intake/session, source integrity, coverage/timebase, packet inspection, flow/session reconstruction, conversations/protocols, DNS, application transactions, encrypted metadata, transfers, Entity relations, behavior candidates, Network Timeline, multi-capture comparison, extraction, provenance and owner handoffs.

No engine, protocol, API, command, detailed screen, Detection Engineering capability or Intelligence capability is created.
