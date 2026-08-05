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
# Capability map — Detection Engineering Authoring

| ID | Capability | Primary role | Primary local concepts | Classes |
| --- | --- | --- | --- | --- |
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
| CAP-INV-416 | Detection Coverage Mapping and Gap Analysis | Detection Coverage Analyst | Detection Coverage Assessment and Detection Gap | 0,1,2 |
| CAP-INV-417 | Detection Authoring Provenance and Review Handoff | Detection Engineering Lead | Detection Authoring Provenance and Review Package | 0,1,2 |

The seventeen capabilities separate intake/project, Detection Hypothesis, data/schema readiness, authoring logic and metadata, validation/test/oracle/replay, candidate match review, coverage/gaps and provenance/handoff. No deployment or Threat Intelligence capability is created.
