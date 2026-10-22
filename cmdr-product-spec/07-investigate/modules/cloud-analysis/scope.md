---
id: investigate-cloud-analysis-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-014, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-012]
---
# Cloud Analysis scope

## In scope
- intake from Case, Incident, Finding, Hunt, Signal, source Cloud or Artifact;
- organization, tenant, account, subscription, project, folder, region and environment context;
- provider/service/resource inventory observations;
- identities, principals, roles, policies, grants, denies, inheritance and permission paths;
- audit/activity observations;
- configuration and state observations;
- compute, workloads, containers, orchestration, serverless and managed execution;
- Cloud network, exposure and connectivity observations;
- storage, data-access and movement candidates;
- sensitive-material candidates and permission-aware access;
- anomalies, misconfiguration candidates and Hypotheses;
- Cloud Timeline and cross-source correlations;
- Evidence, Finding, Detection, TI, Collection and future-response handoff preparation;
- provenance and reproducibility.

## Out of scope
- provider selection or provider-specific product design;
- integrations, APIs, protocols, query languages, parsers, schemas or connectors;
- active discovery, scanning, probing or collection;
- credential validation, secret use, permission grant or revocation;
- resource, policy, network or workload mutation;
- response, containment, shutdown, isolation or deletion;
- automatic Evidence, Finding, attribution, compromise or vulnerability confirmation;
- detailed screens, physical object schemas, final permission matrix and product code;
- Mobile Forensics.

## Scope rule
An authorized projection can be analyzed only within its declared tenant, environment, source, time range and permissions. Configured does not imply accessible; accessible does not imply complete; observed does not imply current.
