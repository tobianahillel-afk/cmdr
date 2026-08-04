---
id: investigate-product
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
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
- Analysis Workbench / Static Analysis — CAP-INV-301..313.

## Counts
Fifty Investigate capabilities are registered: **49 defined**, **1 proposed** (`CAP-INV-106`); all 50 delivery modes are `planned`.

## Boundaries
Command owns operational coordination. Platform Settings owns Fleet, policies and environment administration. Endpoint Agent executes authorized local operations. Govern owns Decision, Approval, Response Run and Result. Studio owns Tool, Tool Call, Workflow and Automation Run. Shared owns generic mechanisms.

## Current maturity
Phase 4B.2B.1 is functionally defined. Phase 4B.2B remains PARTIAL because Dynamic Sandbox, Reverse Engineering, Debugger and forensic work are not completed. No implementation, API, protocol, command, engine, schema, permission matrix or detailed screen rewrite is claimed.
