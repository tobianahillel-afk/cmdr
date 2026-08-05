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

## Counts
**83 capabilities:** 82 defined and one proposed (`CAP-INV-106`); all delivery modes are planned.

## Boundaries
Command owns operational coordination. Settings owns Fleet, policies and administered execution/sandbox environments. Endpoint Agent executes authorized endpoint operations. Govern owns Decision, Approval, Response Run and Result. Studio owns Tool, Tool Call, Workflow and Automation Run. Shared owns generic mechanisms.

Phase 4B.2B.2 is PASS after Dynamic Sandbox and Reverse/Debugger. Phase 4B.2B remains PARTIAL because Forensics is not started. No code, API, protocol, command, engine, debugger product, object schema, atomic permission matrix or detailed screen rewrite is claimed.
