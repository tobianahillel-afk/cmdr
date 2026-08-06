---
id: roadmap-phase-4b4a-cloud-analysis
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-014, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015]
---
# Phase 4B.4A — Cloud Analysis Foundations and Cloud Investigation

## Identifier adaptation
The incoming phase number `4B.4A` was provisional. The canonical roadmap audit found no existing Cloud Analysis subphase identifier. This document therefore records `4B.4A` as the canonical identifier and prevents a second concurrent phase.

## Objective
Define provider-neutral functional capabilities for Cloud investigation without selecting or implementing providers, connectors, APIs, query languages, schemas, engines, scanners, commands or target modifications.

## Scope
`CAP-INV-601..618`, from intake and scope through inventory, identity/IAM, audit, configuration, workloads, network, storage, sensitive material, anomalies, timeline, handoff and provenance.

## Parent grouping
Phase 4B.4 groups extended-environment analysis:
- 4B.4A — Cloud Analysis;
- future Mobile Forensics work remains unstarted and receives no capability or implementation from this phase.

## Exit criteria
- 18/18 capability files;
- 486/486 numbered sections;
- 108/108 mandatory tables;
- at least 170 quality gates;
- no duplicate ID, concurrent owner, active contradiction or implementation claim;
- OPEN-012 reviewed and retained open unless a separate approved decision resolves it;
- remote SHA and PR state verified.

## Status consequence
Cloud Analysis may pass independently. Phase 4B and global maturity remain PARTIAL until Mobile Forensics is completed or formally deferred by an approved roadmap decision.
