---
id: capability-register-investigate-threat-intelligence-analysis-and-products
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-INV-006
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-018
  - OPEN-019
---
# Capability Register — Investigate Threat Intelligence Analysis and Products

IDs are immutable, unique and continue CAP-INV-501..518 without recycling. All capabilities are documentary `draft`, delivery status `defined`, delivery mode `planned`, and owned by Investigate / Threat Intelligence.

| ID | Capability | Primary role | Canonical file | Local concepts | Classes | OPEN |
|---|---|---|---|---|---|---|
| CAP-INV-519 | Intelligence Analysis Intake and Session Management | Threat Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/intelligence-analysis-intake-and-session-management.md` | Intelligence Analysis Session | 0,1,2 | OPEN-013/014/015/018 |
| CAP-INV-520 | Intelligence Questions, Hypotheses and Competing Assessments | Senior Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/intelligence-questions-hypotheses-and-competing-assessments.md` | Intelligence Question, Intelligence Hypothesis, Competing Assessment | 0,1,2 | OPEN-013/018 |
| CAP-INV-521 | Multi-Source Fusion, Corroboration and Structured Analysis | Senior Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/multi-source-fusion-corroboration-and-structured-analysis.md` | Source Fusion Assessment | 0,1,2 | OPEN-013/018 |
| CAP-INV-522 | Threat Actor and Attribution Assessment | Senior Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/threat-actor-and-attribution-assessment.md` | Threat Actor Assessment, Attribution Assessment | 0,1,2 | OPEN-013/018 |
| CAP-INV-523 | Campaign, Intrusion Set and Activity Assessment | Threat Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/campaign-intrusion-set-and-activity-assessment.md` | Campaign, Intrusion Set and Activity Assessments | 0,1,2 | OPEN-013/018 |
| CAP-INV-524 | Malware, Tool, Infrastructure and Capability Assessment | Malware Analyst | `07-investigate/modules/threat-intelligence/capabilities/malware-tool-infrastructure-and-capability-assessment.md` | Malware, Tool, Infrastructure and Capability Assessments | 0,1,2 | OPEN-013/014/018 |
| CAP-INV-525 | Intelligence Product Planning and Audience Definition | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/intelligence-product-planning-and-audience-definition.md` | Product Plan, Audience Definition | 0,1,2 | OPEN-013/018/019 |
| CAP-INV-526 | Intelligence Product Authoring and Structured Assessment | Threat Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/intelligence-product-authoring-and-structured-assessment.md` | Product Draft, Product Version | 0,1,2 | OPEN-013/015/018/019 |
| CAP-INV-527 | Intelligence Review, Quality Control and Release Recommendation | Consumer Reviewer | `07-investigate/modules/threat-intelligence/capabilities/intelligence-review-quality-control-and-release-recommendation.md` | Intelligence Review, Release Recommendation | 0,1,2 | OPEN-013/019 |
| CAP-INV-528 | Handling, Releasability and Dissemination Planning | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/handling-releasability-and-dissemination-planning.md` | Releasability Assessment, Dissemination Plan | 0,1,2 | OPEN-013/018/019 |
| CAP-INV-529 | Internal Intelligence Publication and Consumer Access | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/internal-intelligence-publication-and-consumer-access.md` | Internal Publication Record, Consumer Access Record | 0,1,2 | OPEN-013/019 |
| CAP-INV-530 | Intelligence Watchlist Definition and Lifecycle | Threat Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/intelligence-watchlist-definition-and-lifecycle.md` | Watchlist Definition | 0,1,2 | OPEN-013/017/018/019 |
| CAP-INV-531 | Indicator Operationalization and Detection or Command Handoff | Detection Engineer | `07-investigate/modules/threat-intelligence/capabilities/indicator-operationalization-and-detection-or-command-handoff.md` | Indicator Operationalization Package | 0,1,2 | OPEN-008/013/017/018/019 |
| CAP-INV-532 | Intelligence Monitoring, Sighting Updates and Change Notification | Threat Intelligence Analyst | `07-investigate/modules/threat-intelligence/capabilities/intelligence-monitoring-sighting-updates-and-change-notification.md` | Monitoring Plan, Change Assessment | 0,1,2 | OPEN-008/013/015/018/019 |
| CAP-INV-533 | External Sharing and Exchange Preparation | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/external-sharing-and-exchange-preparation.md` | External Sharing Package | 0,1,2 | OPEN-013/015/018/019 |
| CAP-INV-534 | Consumer Feedback and Intelligence Effectiveness Assessment | Consumer Reviewer | `07-investigate/modules/threat-intelligence/capabilities/consumer-feedback-and-intelligence-effectiveness-assessment.md` | Feedback Record, Effectiveness Assessment | 0,1,2 | OPEN-013/018/019 |
| CAP-INV-535 | Intelligence Requirement Satisfaction and Collection Feedback | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/intelligence-requirement-satisfaction-and-collection-feedback.md` | Requirement Satisfaction Assessment, Collection Feedback Package | 0,1,2 | OPEN-008/013/015/018 |
| CAP-INV-536 | Intelligence Correction, Retraction and Product Supersession | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/intelligence-correction-retraction-and-product-supersession.md` | Correction Record, Retraction Record, Product Supersession | 0,1,2 | OPEN-013/018/019 |
| CAP-INV-537 | Intelligence Lifecycle Provenance, Closure and Continuous Improvement | Intelligence Manager | `07-investigate/modules/threat-intelligence/capabilities/intelligence-lifecycle-provenance-closure-and-continuous-improvement.md` | Continuous Improvement Package | 0,1,2 | OPEN-013/014/015/018/019 |

## Totals
- Capabilities: **19/19**.
- Numbered sections: **513/513**.
- Mandatory S8/S9/S10/S13/S16/S17 tables: **114/114**.
- Defined / proposed / planned: **19 / 0 / 19**.
- Duplicate IDs, concurrent owners, empty mandatory tables and implemented claims: **0**.
