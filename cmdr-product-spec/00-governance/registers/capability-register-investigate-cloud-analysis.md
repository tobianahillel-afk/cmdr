---
id: capability-register-investigate-cloud-analysis
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: registry
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Capability Register — Investigate Cloud Analysis

Phase canonique : **4B.4A — Cloud Analysis Foundations and Cloud Investigation**.

Ce shard enregistre les capabilities fonctionnelles Cloud. `defined` / `planned` signifie documenté, pas implémenté. Aucun provider, API, protocole, query language, schema, connector, command, target mutation ou Mobile Forensics n’est livré.

| ID | Title | Primary role | Canonical path | Main local concepts | Action classes | Open decisions | Status | Delivery |
|---|---|---|---|---|---|---|---|---|
| CAP-INV-601 | Cloud Analysis Intake and Preconditions | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-analysis-intake-and-preconditions.md` | Cloud Analysis Intake, Cloud Preconditions Assessment, Collection or access gap proposal | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-602 | Cloud Investigation Session and Workspace Management | Investigation Lead | `07-investigate/modules/cloud-analysis/capabilities/cloud-investigation-session-and-workspace-management.md` | Cloud Investigation Session, Session membership and view state, Session lifecycle event | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-603 | Cloud Scope, Organization, Tenant and Account Context | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-scope-organization-tenant-and-account-context.md` | Cloud Scope Assessment, Scope inclusion/exclusion set, Cross-account relation observation | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-604 | Cloud Provider, Service and Resource Inventory Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-provider-service-and-resource-inventory-analysis.md` | Cloud Resource Observation, Cloud Inventory Snapshot Assessment, Unknown or inventory gap | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-605 | Cloud Identity, Principal and Role Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-identity-principal-and-role-analysis.md` | Cloud Identity Observation, Role Observation, Identity relationship candidate | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-606 | Cloud IAM Policy and Effective Permission Path Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-iam-policy-and-effective-permission-path-analysis.md` | Permission Observation, Effective Permission Candidate, Permission Path Assessment | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-607 | Cloud Audit Event and Activity Log Analysis | DFIR Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-audit-event-and-activity-log-analysis.md` | Cloud Activity Observation, Audit coverage and gap assessment, Event correlation candidate | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-608 | Cloud Resource Configuration and State Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-resource-configuration-and-state-analysis.md` | Configuration Observation, Configuration Drift Assessment, Misconfiguration Candidate | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-609 | Cloud Compute and Workload Analysis | DFIR Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-compute-and-workload-analysis.md` | Compute Observation, Workload Observation, Endpoint/Disk/Memory handoff package | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-610 | Container and Orchestration Workload Analysis | DFIR Analyst | `07-investigate/modules/cloud-analysis/capabilities/container-and-orchestration-workload-analysis.md` | Container Observation, Orchestration Workload Observation, Container/Endpoint handoff package | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-611 | Serverless and Managed Execution Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/serverless-and-managed-execution-analysis.md` | Serverless Observation, Managed Execution Change Assessment, Artifact/Detection handoff proposal | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-612 | Cloud Network, Exposure and Connectivity Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-network-exposure-and-connectivity-analysis.md` | Cloud Network Observation, Network Exposure Observation, Connectivity Assessment | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-613 | Cloud Storage and Data Access Analysis | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-storage-and-data-access-analysis.md` | Storage Observation, Data Access Observation, Data Movement Candidate | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-614 | Cloud Secrets, Keys and Sensitive Material Assessment | Evidence Reviewer | `07-investigate/modules/cloud-analysis/capabilities/cloud-secrets-keys-and-sensitive-material-assessment.md` | Sensitive Material Candidate, Sensitive access record, Evidence Candidate context | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-615 | Cloud Anomaly, Misconfiguration and Hypothesis Management | Cloud Security Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-anomaly-misconfiguration-and-hypothesis-management.md` | Cloud Anomaly, Cloud Hypothesis, Misconfiguration Candidate disposition | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-616 | Cloud Timeline and Cross-Source Correlation | DFIR Analyst | `07-investigate/modules/cloud-analysis/capabilities/cloud-timeline-and-cross-source-correlation.md` | Cloud Timeline, Cross-source Correlation Candidate, Time quality assessment | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |
| CAP-INV-617 | Cloud Analysis Handoff to Evidence, Findings and Detection | Investigation Lead | `07-investigate/modules/cloud-analysis/capabilities/cloud-analysis-handoff-to-evidence-findings-and-detection.md` | Evidence Candidate Package, Finding Draft, Cloud Detection Gap / Engineering Package, TI / Collection / future response handoff | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015/OPEN-017 | defined | planned |
| CAP-INV-618 | Cloud Analysis Provenance and Reproducibility | Evidence Reviewer | `07-investigate/modules/cloud-analysis/capabilities/cloud-analysis-provenance-and-reproducibility.md` | Cloud Reproducibility Assessment, Cloud Analysis Provenance Package, Correction or gap proposal | 0,1,2 | OPEN-008/OPEN-012/OPEN-013/OPEN-014/OPEN-015 | defined | planned |

## Totals
- capabilities: **18**;
- numbered sections: **486**;
- mandatory tables: **108**;
- delivery: **18 defined / 18 planned**;
- IDs: `CAP-INV-601..618`, continuous and unique;
- implementation claims: **0**.

## Ownership
Investigate owns local Cloud analytical concepts only. Platform Settings retains providers, connectors, credentials, configured scopes, ingestion, schemas, parsers, health, retention, storage and policies. Command retains runtime operational objects. Govern retains target mutations and authority. Studio retains Tools and Runs. Shared retains generic Entity, Graph, Timeline, Linking, Trace, Export, Reporting, Collaboration and Recovery.
