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

Investigate owns investigation work: Signals and Hunt, Cases and Evidence, and Collection and Live Response. It preserves analyst reasoning, provenance and the path from source data to Evidence and Finding.

## Canonical modules
- Signals and Hunt — CAP-INV-001..008.
- Cases and Evidence — CAP-INV-101..114.
- Collection and Live Response — CAP-INV-201..215.

## Boundaries
Command owns Detection, Signal, Alert, Incident and operational Task. Platform Settings owns Endpoint Agent Fleet and Endpoint Policies. Endpoint Agent executes authorized local operations and reports health/results. Govern owns Action Request lifecycle, Decision, Approval, Response Run and Result. Studio owns Workflow, Tool, Tool Call and Automation Run. Shared owns generic jobs, linking, timeline, trace, notifications, export and collaboration mechanisms.

## Delivery
Thirty-seven Investigate capabilities are registered: 36 `defined`, one `proposed` (`CAP-INV-106`); all delivery modes are `planned`. Phase 4B.2B Analysis Workbench and Phase 4B.3 are not started.

## Non-goals
No Fleet administration, Policy mutation, protocol, API, command syntax, engine choice, product code, final object schema, atomic permission matrix or detailed screen rewrite is defined here.