---
id: investigate-product
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-014
---
# Investigate

Investigate owns investigation work, analyst reasoning, provenance and the path from source data to Evidence and Finding.

## Canonical modules
- Signals and Hunt — CAP-INV-001..008.
- Cases and Evidence — CAP-INV-101..114.
- Collection and Live Response — CAP-INV-201..215.
- Analysis Workbench and Static Analysis — CAP-INV-301..313.
- Dynamic Sandbox and Behavioral Analysis — CAP-INV-314..328.
- Reverse Engineering and Debugger — CAP-INV-329..346.
- Memory Forensics — CAP-INV-347..362.
- Disk and Filesystem Forensics — CAP-INV-363..379.
- Network Forensics — CAP-INV-380..397.

## Counts
**134 capabilities:** 133 defined and one proposed (`CAP-INV-106`); all delivery modes are planned. Phase 4B.1 contributes 22 and Phase 4B.2 contributes 112.

## Boundaries
Command owns Detection, Signal, Alert, Incident and operational coordination. Settings owns Fleet, policies, sensors, storage, retention, health, time synchronization and administered environments. Endpoint Agent and Collection contribute authorized acquisition results. Govern owns Decision, Approval, Response Run, Result and real-target authority. Studio owns Tool, Tool Call, Workflow and Automation Run. Shared owns Entity, Graph, Timeline and generic mechanisms.

## Maturity
Phase 4B.2 is PASS after Network Forensics and cross-phase closure. Phase 4B remains PARTIAL because Phase 4B.3 Detection Engineering and Intelligence is not started. No code, API, protocol, command, engine, object schema, atomic permission matrix or detailed screen rewrite is claimed.
