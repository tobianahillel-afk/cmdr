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
  - REQ-INV-002
  - REQ-INV-003
  - REQ-INV-004
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Analysis Workbench

## Mission
Provide one Technical Workbench context for Artifact intake, static analysis, Reverse Engineering, isolated Debugger work, analytical knowledge, provenance and explicit Evidence/Finding handoff.

## Canonical capability ranges
- CAP-INV-301..313 — Workbench Foundation and Static Analysis.
- CAP-INV-329..346 — Reverse Engineering and Debugger.
- CAP-INV-314..328 remains the separate Dynamic Sandbox module and is consumed through handoffs.

## Ownership
Investigate owns Case, Artifact, analytical sessions, annotations, interpretations and candidates/drafts. Studio owns Tool, Tool Call, Workflow and Automation Run. Platform Settings owns providers, secrets and execution environments. Govern owns real-target authority. Shared owns generic mechanisms.

## Workbench constraints
One main canvas, one Inspector, at most two auxiliary panels, at most six visible technical tabs, optional bottom Console, contextual Artifact Explorer, Automation Tray closed by default, accessible resizers, preserved history/return/focus and structured alternatives to graphs.

## Delivery
Thirty-one Analysis Workbench capabilities are defined and planned: thirteen static and eighteen Reverse/Debugger. No engine, third-party product, API, protocol, command, code, detailed screen rewrite, final object schema or atomic permission matrix is claimed.
