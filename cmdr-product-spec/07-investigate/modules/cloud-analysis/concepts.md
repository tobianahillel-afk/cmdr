---
id: investigate-cloud-analysis-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-020]
open_decisions: [OPEN-012, OPEN-014]
---
# Cloud Analysis concepts

## Local analytical concepts
- Cloud Analysis Intake;
- Cloud Preconditions Assessment;
- Cloud Investigation Session;
- Cloud Scope Assessment;
- Cloud Resource Observation;
- Cloud Identity Observation;
- Role Observation;
- Permission Observation;
- Effective Permission Candidate;
- Permission Path Assessment;
- Cloud Activity Observation;
- Configuration Observation;
- Compute, Workload, Container, Serverless, Network, Storage and Data Access Observations;
- Sensitive Material Candidate;
- Cloud Anomaly and Cloud Hypothesis;
- Cloud Timeline and correlation candidate;
- Evidence Candidate Package, Finding Draft and Detection Engineering Package;
- Cloud Analysis Provenance Package and Reproducibility Assessment.

These are functional concepts, not final canonical object schemas.

## Mandatory distinctions
- Cloud Account ≠ tenant ≠ organization ≠ subscription ≠ project;
- provider configured ≠ provider accessible;
- accessible ≠ complete coverage;
- inventory ≠ current state certain;
- resource/configuration observed ≠ current;
- audit event ≠ confirmed human action;
- identity/principal ≠ person;
- role assignment/policy ≠ certain effective permission;
- permission path ≠ exploit path;
- privilege candidate ≠ abuse;
- public endpoint ≠ exploitable endpoint;
- exposed storage ≠ confirmed data exposure;
- network rule ≠ observed connection;
- metadata ≠ content;
- secret candidate ≠ valid credential;
- secret existence ≠ permission to reveal/use;
- image ≠ running workload;
- container observation ≠ full container forensics;
- serverless invocation ≠ malicious execution;
- anomaly/misconfiguration candidate ≠ Finding/vulnerability;
- Cloud Timeline ≠ Case Timeline;
- Tool result ≠ analyst conclusion;
- Evidence candidate ≠ qualified Evidence;
- Finding Draft ≠ confirmed Finding;
- investigation action ≠ response action.

## Ownership
Shared retains canonical Entity, Graph, Timeline, Linking and Versioning. Cloud Analysis records local observations and projections without creating competing global objects.
