---
id: investigate-detection-engineering-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Scope

## Included
- Intake, Project, Detection Hypothesis, data readiness, schema/field review and authoring.
- Conditions, correlation, sequence, threshold, windows, enrichment, metadata and documentation.
- Structural/semantic validation, test scenarios/datasets, Expected Outcomes, historical replay, candidate FP/FN review, coverage, gaps and Review Package.
- Release Candidate review, readiness, target/promotion planning and Govern handoff.
- Shadow non-alerting observation, canary planning and deployment coordination projections.
- Runtime version reconciliation, health, Command feedback, production match review, tuning proposals, bounded suppression/exception proposals, drift, performance, rollback, retirement and lifecycle provenance.

## Excluded
- Selecting a runtime, engine, rule language, adapter, vendor syntax or executable format.
- Direct administration of environments, targets, sources, parsers, schemas, providers or secrets.
- Direct class-3/4 production execution by Investigate.
- Any actual deployment, activation, deactivation, suppression, exception, rollback, retirement or history deletion during documentation.
- Threat Intelligence, CAP-INV-5xx, Cloud Analysis, Mobile Forensics, APIs, protocols, commands, code and detailed screen design.
