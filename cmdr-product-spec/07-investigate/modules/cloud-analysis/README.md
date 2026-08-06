---
id: investigate-cloud-analysis-readme
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-012, OPEN-013, OPEN-014, OPEN-015]
---
# Cloud Analysis

## Canonical phase
The canonical identifier introduced by this module is **Phase 4B.4A — Cloud Analysis Foundations and Cloud Investigation**.

The incoming identifier was provisional. The roadmap audit found no existing Cloud Analysis subphase identifier and no active competing module. `4B.4A` is therefore retained and recorded as the canonical adaptation. Phase 4B.4 groups extended-environment analysis: Cloud is 4B.4A and Mobile Forensics remains future scope, with no Mobile capability created here.

## Mission
Cloud Analysis gives authorized Cloud Security Analysts, Investigation Leads, SOC Analysts, DFIR Analysts and Evidence Reviewers a provider-neutral, evidence-first workspace for investigating Cloud scope, identities, permissions, events, configurations, workloads, networks, storage, sensitive material, anomalies, timelines and handoffs.

## Functional range
- `CAP-INV-601..618`;
- 18 capabilities;
- 486 numbered capability sections;
- 108 mandatory capability tables;
- documentary and functional only;
- `defined` / `planned`, never implemented or deployed.

## Boundaries
Investigate owns Cloud analytical Sessions, observations, assessments, candidates, Hypotheses, timelines and handoff packages. It does not own providers, connectors, credentials, Cloud target configuration, runtime Detections, Signals, Incidents, Decisions, Approvals, Tools, generic Graph/Timeline/Reporting or any response execution.

No provider, API, protocol, query language, schema, connector, scanner, command, script, configuration or product code is selected or created.

## Reading order
1. `scope.md`;
2. `concepts.md`;
3. `user-questions.md`;
4. `workflows.md`;
5. `states.md`;
6. `permissions.md`;
7. `sensitive-information.md`;
8. `automation-and-ai.md`;
9. `platform-settings-boundaries.md`;
10. `command-studio-govern-boundaries.md`;
11. `shared-capabilities.md`;
12. `source-migration.md`;
13. `capability-map.md`;
14. `capabilities/`.

## Source of truth
Capability files own local functional behavior. Platform Settings owns configured sources and Cloud administration. Shared owns generic mechanisms. Security owns final permission policy. Object schemas, screens and implementation contracts remain future.

## Non-goals
No Cloud integration, connector, inventory collection, active scan, credential validation, secret use, permission change, resource mutation, containment, Detection rule, response, Mobile Forensics or implementation.
