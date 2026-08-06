---
id: investigate-detection-engineering
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-INV-006
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-AI-002
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
---
# Detection Engineering

## Mission
Transform a sourced need into governed Detection Content through two connected but distinct stages:

1. **CAP-INV-401..417** — foundations, authoring, validation, controlled test data, historical replay, candidate FP/FN review, coverage and Review Package.
2. **CAP-INV-418..435** — formal review, Release Candidate, readiness, promotion planning, Govern handoff, shadow/canary observation, deployment coordination, runtime state/health/quality assessment, tuning, bounded suppression/exception proposals, drift, performance, rollback, retirement and continuous improvement.

## Ownership
Investigate owns Detection Engineering business concepts and assessments. Command retains runtime Detection, Signal, Alert, Incident and operational feedback. Settings retains sources, parsers, schemas, environments, targets, configured runtimes, health, providers, secrets and retention. Endpoint retains declared capabilities and local execution projections. Govern retains Action Request, Decision, Approval, Response Run, Result and production authority. Studio retains Tools, Tool Calls, Workflows, Automation Runs and generic Evaluation/Dataset mechanisms. Shared retains generic Jobs, Trace, Activity, Versioning, Linking, Comparison, Search, Export, Reporting, Collaboration and Recovery.

## Decision boundary
`OPEN-017` is open for Detection runtime, target language and portability strategy. No engine, language, adapter, canonical executable representation or vendor option is selected.

## Safety
No real rule is deployed, activated, deactivated, suppressed, excepted, rolled back or retired by this documentation phase. No Signal or Alert is deleted. No API, protocol, code, command, package format, compiler, parser, AST, model or detailed screen is defined.

## Delivery
Phase 4B.3A.1: PASS. Phase 4B.3A.2: PASS after publication verification. Phase 4B.3A: PASS. Threat Intelligence remains not started.
