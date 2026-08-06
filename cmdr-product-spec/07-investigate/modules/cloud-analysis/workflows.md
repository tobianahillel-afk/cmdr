---
id: investigate-cloud-analysis-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-014, REQ-PROD-020, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-012, OPEN-013, OPEN-015]
---
# Cloud Analysis workflows

## Nominal flow
`Origin → Intake → Session → Scope → Inventory → Identity/IAM/Audit → Configuration and workload families → Anomalies/Hypotheses → Timeline/Correlation → Handoffs → Provenance`.

## Intake and session
CAP-INV-601 qualifies origin, scope candidates, sources, periods, access and gaps. CAP-INV-602 manages the versioned Session and collaboration context.

## Scope and observations
CAP-INV-603 preserves hierarchy and exclusions. CAP-INV-604..614 create sourced observations and candidates without provider mutation.

## Interpretation
CAP-INV-615 manages anomaly and Hypothesis candidates. CAP-INV-616 reconstructs time and cross-source correlations without asserting causality.

## Handoff and closure
CAP-INV-617 prepares packages; destination owners accept, reject or return them. CAP-INV-618 preserves provenance and assesses reproducibility.

## Failure and recovery
A source, permission, parser, provider projection or Tool failure produces partial/blocked/stale states while preserving valid prior results. Conflict resolution uses diff, versions, human disposition and Shared Recovery.

## No-action rule
No workflow step grants Cloud access, executes a command, changes a target, validates a secret, deploys a rule or triggers a response.
