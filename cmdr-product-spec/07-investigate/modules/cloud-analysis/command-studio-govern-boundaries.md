---
id: investigate-cloud-analysis-product-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-013, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016]
open_decisions: [OPEN-012, OPEN-013, OPEN-015, OPEN-017]
---
# Command, Studio and Govern boundaries

## Command
Command owns runtime Detection, Signal, Alert, Incident, priority and operational dispositions. A Cloud anomaly, audit observation or correlation does not become a Signal or Incident automatically.

## Detection Engineering
Detection Engineering owns Detection Content, Hypothesis, Gap, Coverage, validation, promotion and lifecycle. Cloud Analysis prepares a package only.

## CMDR Studio
Studio owns Tool, Tool Call, Workflow, Automation Run, Automation Agent and Human Gate. Cloud Analysis links execution provenance; it does not redefine or schedule Tools.

## Govern
Govern owns Decision, Approval, Action Request, Response Run, Result and all real target mutation: permission changes, credential revocation, isolation, shutdown, deletion, blocking and remediation.

## Return contract
Every handoff carries tenant, environment, scope, source versions, permissions, restrictions, evidence links, errors, uncertainty, requested action, destination owner and return origin.
