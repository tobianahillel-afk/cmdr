---
id: source-product-boundaries
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# Product Boundaries

## Command

Owns prioritization, operational coordination, Incident, Work Queue, assignment, SLA, situation, handover, service impact and readiness. It does not own forensic analysis, rule editing, reverse engineering, endpoint terminal operations, detailed approval policy or agent building.

## Investigate

Owns Case, Hypothesis, Artifact, Evidence, Finding, search, Hunt, collection from the investigation perspective, technical workbench, Detection Engineering and Intelligence. It does not independently authorize high-risk response.

## Govern

Owns Action Request processing, Decision, Approval, policy evaluation, Response Run, verification, rollback and audit ledger. It is not a second general Work Queue or Case workspace.

## CMDR Studio

Owns Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate, Automation Run, evaluation, simulation, deployment, costs, usage and Control Room. It does not replace operational products.

## Platform Settings

Owns users, roles, tenants, environments, secrets, integrations, model providers, Endpoint Agent Fleet, endpoint policies, data sources, retention, health and administrative audit. It is not an investigation workspace.

## Endpoint Agent

A distinct product component and a planned complete native EDR. Settings administers fleet and policies; Investigate uses endpoints for inspection and collection; Govern controls risky actions.

## Shared Capabilities

Owns Reporting Engine, generic Saved Views and other cross-product services. Product-specific projections do not redefine these capabilities.