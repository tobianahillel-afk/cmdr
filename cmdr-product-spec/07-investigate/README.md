---
id: investigate-product
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-INV-001
  - REQ-INV-006
open_decisions:
  - OPEN-012
  - OPEN-018
  - OPEN-019
---
# Investigate

Investigate owns investigation work, analyst reasoning, Evidence and Finding preparation, specialized technical analysis, Detection Engineering, Threat Intelligence and provider-neutral Cloud analytical concepts. It does not own platform administration, production response authority or Shared engines.

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
- Intelligence Analysis, Products, Dissemination and Operationalization — CAP-INV-519..537.
- Cloud Analysis Foundations and Cloud Investigation — CAP-INV-601..618.

## Counts after Cloud closure content
**224 capabilities:** 223 defined and one proposed (`CAP-INV-106`); all delivery modes are planned. Investigate contains 6048 numbered capability sections and 1344 mandatory tables.

## Boundaries
Command owns runtime Detection, Signal, Alert, Incident and operational coordination. Settings owns Data Sources, providers, connectors, credentials, secrets, configured Cloud scopes, schemas, environments, health, retention, Fleet and policies. Endpoint Agent owns declared local capabilities and execution projections. Govern owns Action Request, Decision, Approval, external release authority, Response Run, Result and real-target authority. Studio owns Tool, Tool Call, Workflow, Evaluation and Automation Run. Shared owns Entity, Graph, Query, Timeline, Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Comparison and collaboration.

Threat Intelligence creates functional candidates, assessments and products; it does not silently confirm attribution, merge identity, deploy an Indicator, activate a watchlist, create Detection Content, block, respond, publish externally or share externally.

Cloud Analysis creates scoped Sessions, observations, candidates, Hypotheses, correlations and handoff packages. It does not configure a provider, use a secret, scan a target, execute a command, qualify Evidence, confirm a Finding, deploy Detection Content or perform response.

## Maturity
Phase 4B.3 is PASS. Phase 4B.4A remains `PENDING POST-PUBLICATION VERIFICATION` until the fifth functional commit and remote checks are complete. Phase 4B.4, Phase 4B, Phase 4 and global maturity remain PARTIAL. Mobile Forensics is NOT STARTED. No code, API, protocol, selected provider, connector, actual collection, deployment, target mutation, complete object schema, atomic permission matrix or detailed Cloud screen rewrite is claimed.
