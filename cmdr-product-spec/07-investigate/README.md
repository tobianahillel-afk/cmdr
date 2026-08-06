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
---
# Investigate

Investigate owns investigation work, analyst reasoning, Evidence/Finding preparation and Detection Engineering business context.

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

## Counts
**169 capabilities:** 168 defined and one proposed (`CAP-INV-106`); all delivery modes are planned. Phase 4B.1 contributes 22, Phase 4B.2 contributes 112 and complete Detection Engineering contributes 35.

## Boundaries
Command owns runtime Detection, Signal, Alert, Incident and operational coordination. Settings owns Data Sources, parsers, schemas, environments, targets, configured runtimes, health, retention, Fleet and policies. Endpoint Agent owns declared telemetry capabilities and local runtime behavior. Govern owns Action Request, Decision, Approval, Response Run, Result and production authority. Studio owns Tool, Tool Call, Workflow, Dataset/Evaluation mechanisms and Automation Run. Shared owns Query, Telemetry Event, Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Comparison and collaboration.

Investigate prepares and assesses; it does not silently deploy, activate, deactivate, suppress, except, roll back or retire runtime content.

## Maturity
Phase 4B.3A.1 and Phase 4B.3A.2 are PASS; Phase 4B.3A is PASS. Phase 4B.3 and Phase 4B remain PARTIAL because Threat Intelligence is not started. No code, API, protocol, selected engine/language, actual deployment, active rule change, Intelligence capability, object schema, atomic permission matrix or detailed screen rewrite is claimed.
