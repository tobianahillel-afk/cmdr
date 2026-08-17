---
id: investigate-memory-forensics-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Source migration

| Legacy source | Useful needs migrated | Replacement | Status |
|---|---|---|---|
| `07-investigate/modules/memory-forensics/README.md` | processes, injection candidates, memory network, kernel, secret masking | CAP-INV-347..362 and module README | deprecated pointer |
| `07-investigate/modules/memory-forensics/plugin-contract.md` | producer/version/parameters, partial failure, provenance and sensitive masking | Studio Tool boundaries, CAP-INV-348/350/356/361 | deprecated pointer |
| `07-investigate/modules/memory-forensics/screens/memory-forensics.md` | Workbench entry, compare, extract, Inspector and states | Screen remains active; capability map updated | active screen |

## Acceptance
No competing functional architecture remains; useful needs are preserved; the active screen is not deprecated; date and replacement are explicit; links point to the canonical module.
