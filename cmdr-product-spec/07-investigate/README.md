---
id: investigate-product
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-INV-006
open_decisions:
  - OPEN-018
---
# Investigate

Investigate owns investigation work, analyst reasoning, Evidence/Finding preparation, Detection Engineering business context and Threat Intelligence foundations/knowledge management.

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
- Detection Engineering Foundations, Authoring and Validation — CAP-INV-401..417.
- Detection Review, Promotion, Runtime Performance and Lifecycle — CAP-INV-418..435.
- Threat Intelligence Foundations and Knowledge Management — CAP-INV-501..518.

## Counts
**187 capabilities:** 186 defined and one proposed (`CAP-INV-106`); all delivery modes are planned. Phase 4B.1 contributes 22, Phase 4B.2 contributes 112, Detection Engineering contributes 35 and Threat Intelligence Foundations contributes 18.

## Boundaries
Command owns runtime Detection, Signal, Alert, Incident and operational coordination. Settings owns Data Sources, providers, connectors, schemas, environments, health, retention, Fleet and policies. Endpoint Agent owns declared telemetry capabilities and local runtime behavior. Govern owns Action Request, Decision, Approval, external release authority, Response Run and Result. Studio owns Tool, Tool Call, Workflow, Evaluation and Automation Run. Shared owns Entity, Graph, Query, Timeline, Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Comparison and collaboration.

Threat Intelligence creates functional candidates and assessments; it does not silently confirm attribution, merge identity, deploy Indicator, activate watchlist, create Detection Content, block, respond, publish or share externally.

## Maturity
Phase 4B.3A is PASS. Phase 4B.3B.1 is PASS after publication verification. Phase 4B.3B, Phase 4B.3 and Phase 4B remain PARTIAL because 4B.3B.2 is not started. No code, API, protocol, selected standard/provider, actual collection/deployment/sharing, object schema, atomic permission matrix or detailed screen rewrite is claimed.
