---
id: investigate-analysis-workbench
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-INV-001
  - REQ-INV-002
  - REQ-INV-003
  - REQ-INV-004
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Analysis Workbench

## Mission
Provide one Technical Workbench context for Artifact intake, Static Analysis, Reverse Engineering, isolated Debugger work, Memory Forensics, Disk and Filesystem Forensics, Network Forensics, analytical knowledge, provenance and explicit Evidence/Finding handoff.

## Canonical capability ranges
- CAP-INV-301..313 — Workbench Foundation and Static Analysis.
- CAP-INV-329..346 — Reverse Engineering and Debugger.
- CAP-INV-347..362 — Memory Forensics.
- CAP-INV-363..379 — Disk and Filesystem Forensics.
- CAP-INV-380..397 — Network Forensics.
- CAP-INV-314..328 remains the separate Dynamic Sandbox module and is consumed through explicit handoffs.

## Ownership
Investigate owns Case, Artifact, Memory Image, Disk Image, analytical sessions and concepts, observations, annotations, interpretations, Derived Artifacts and candidate packages. Collection and Endpoint Agent contribute acquisition results. Studio owns Tool, Tool Call, Workflow and Automation Run. Platform Settings owns Fleet, sensors, policies, providers, secrets, storage, retention, health, time synchronization and environments. Command owns Detection, Signal, Alert and Incident. Govern owns real-target authority. Shared owns Entity, Graph, Timeline and generic mechanisms.

## Workbench constraints
One main canvas, one Inspector, at most two auxiliary panels, at most six visible technical tabs, optional bottom Console, contextual Artifact Explorer, Automation Tray closed by default, accessible resizers, preserved history/return/focus, protected sensitive values and structured alternatives to graphs.

## Delivery
Eighty-two Analysis Workbench capabilities are defined and planned, excluding the separate fifteen-capability Dynamic Sandbox module. No engine, third-party product, API, protocol, command, code, detailed screen rewrite, final object schema or atomic permission matrix is claimed.
