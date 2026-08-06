---
id: investigate-detection-engineering-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-INV-006
  - REQ-PROD-012
  - REQ-PROD-014
---
# Capability map — Complete Detection Engineering

| ID | Capability | Primary role | Primary local concepts | Classes |
|---|---|---|---|---|
| CAP-INV-401 | Detection Engineering Intake and Preconditions | Detection Engineer | Detection Engineering Intake | 0,1,2 |
| CAP-INV-402 | Detection Engineering Project and Workspace Management | Detection Engineering Lead | Detection Engineering Project | 0,1,2 |
| CAP-INV-403 | Detection Hypothesis and Objective Management | Detection Engineer | Detection Hypothesis | 0,1,2 |
| CAP-INV-404 | Data Source and Telemetry Readiness Assessment | Detection Engineer | Telemetry Readiness Assessment | 0,1,2 |
| CAP-INV-405 | Event Schema, Field and Normalization Mapping Review | Detection Engineer | Field Mapping Review | 0,1,2 |
| CAP-INV-406 | Detection Content Authoring | Detection Engineer | Detection Content Draft | 0,1,2 |
| CAP-INV-407 | Detection Conditions and Correlation Logic Design | Detection Engineer | Detection Condition and Correlation Definition | 0,1,2 |
| CAP-INV-408 | Sequence, Threshold and Time-Window Design | Detection Engineer | Sequence and Threshold Definition | 0,1,2 |
| CAP-INV-409 | Detection Enrichment and Context Requirement Design | Detection Engineer | Enrichment Requirement | 0,1,2 |
| CAP-INV-410 | Detection Metadata, Ownership and Documentation | Detection Content Owner | Detection Metadata Record | 0,1,2 |
| CAP-INV-411 | Detection Content Structural and Semantic Validation | Detection Validator | Validation Result | 0,1,2 |
| CAP-INV-412 | Detection Test Scenario and Dataset Management | Detection Test Designer | Detection Test Scenario and Dataset | 0,1,2 |
| CAP-INV-413 | Expected Outcome and Test Oracle Management | Detection Test Reviewer | Expected Outcome | 0,1,2 |
| CAP-INV-414 | Historical Replay and Backtesting | Detection Engineer | Replay Result | 0,1,2 |
| CAP-INV-415 | Match Review and False Positive/False Negative Analysis | Detection Reviewer | Match Review | 0,1,2 |
| CAP-INV-416 | Detection Coverage Mapping and Gap Analysis | Detection Coverage Analyst | Coverage Assessment and Detection Gap | 0,1,2 |
| CAP-INV-417 | Detection Authoring Provenance and Review Handoff | Detection Engineering Lead | Authoring Provenance and Review Package | 0,1,2 |
| CAP-INV-418 | Detection Review Queue and Release Candidate Management | Detection Reviewer | Release Candidate, Detection Review | 0,1,2; 3/4 owner-only |
| CAP-INV-419 | Deployment Readiness and Production Preconditions | Detection Reviewer | Deployment Readiness Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-420 | Detection Change Request and Approval Handoff | Detection Owner | Detection Change Request Draft | 0,1,2; 3/4 owner-only |
| CAP-INV-421 | Environment Promotion Planning and Target Selection | Platform Operator | Promotion Plan, Deployment Target selection | 0,1,2; 3/4 owner-only |
| CAP-INV-422 | Shadow Evaluation and Non-Alerting Observation | Detection Engineer | Shadow Evaluation Plan, Shadow Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-423 | Canary and Phased Rollout Planning | Detection Owner | Canary Plan, Canary Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-424 | Detection Deployment and Activation Coordination | Detection Owner | Deployment coordination projection | 0,1,2; 3/4 owner-only |
| CAP-INV-425 | Runtime Detection Version and State Reconciliation | Platform Operator | Runtime Version Observation | 0,1,2; 3/4 owner-only |
| CAP-INV-426 | Detection Health and Execution Monitoring | Platform Operator | Detection Health Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-427 | Runtime Signal Quality and Operational Feedback | Detection Owner | Runtime Quality Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-428 | Production Match Review and Outcome Reconciliation | Detection Reviewer | Production Match Review | 0,1,2; 3/4 owner-only |
| CAP-INV-429 | Detection Tuning Proposal Management | Detection Engineer | Tuning Proposal | 0,1,2; 3/4 owner-only |
| CAP-INV-430 | Suppression and Exception Proposal Management | Detection Owner | Suppression Proposal, Exception Proposal | 0,1,2; 3/4 owner-only |
| CAP-INV-431 | Data, Schema and Dependency Drift Assessment | Detection Engineer | Drift Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-432 | Detection Performance and Resource Impact Assessment | Platform Operator | Performance Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-433 | Detection Rollback and Recovery Planning | Detection Owner | Rollback Plan, Recovery Assessment | 0,1,2; 3/4 owner-only |
| CAP-INV-434 | Detection Deactivation, Retirement and Replacement Management | Detection Owner | Retirement Proposal, Replacement Relation | 0,1,2; 3/4 owner-only |
| CAP-INV-435 | Detection Lifecycle Provenance and Continuous Improvement Handoff | Detection Owner | Lifecycle Provenance Assessment, Continuous Improvement Package | 0,1,2; 3/4 owner-only |
## Totals and boundary
- CAP-INV-401..417: 17 authoring/validation capabilities.
- CAP-INV-418..435: 18 review/runtime-lifecycle capabilities.
- Total Detection Engineering: 35 capabilities, 945 numbered sections and 210 mandatory tables.
- Investigate executes classes 0–2 only. Class-3/4 changes remain Govern/Settings/runtime-owner actions.
- No CAP-INV-5xx or Threat Intelligence capability is created.
