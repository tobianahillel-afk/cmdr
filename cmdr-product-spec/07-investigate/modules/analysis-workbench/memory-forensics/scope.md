---
id: investigate-memory-forensics-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
open_decisions:
  - OPEN-005
  - OPEN-008
---
# Memory Forensics scope

## Included
Memory Image intake, session context, integrity/acquisition review, platform/profile identification, process/thread reconstruction, regions/mappings, modules/drivers, handles/IPC, memory-resident network state, sensitive material controls, anomalies, kernel indicators, timeline, extraction, provenance and handoff.

## Excluded
Acquisition mechanics; engines/plugins; offsets/signatures; commands; Disk, Filesystem or full Network Forensics; cloud/mobile; exploitation; credential extraction procedures; Detection Engineering rules; APIs, protocols and code.

## Boundaries
Collection produces the Memory Image. Debugger memory views remain session snapshots. Memory Forensics analyzes a complete or partial acquired image. Settings owns Fleet/Policies/storage/retention. Studio owns Tools/Runs. Govern owns real-target authority.