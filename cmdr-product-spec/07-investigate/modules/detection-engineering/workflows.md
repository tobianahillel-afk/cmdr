---
id: investigate-detection-engineering-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Workflows

## Canonical authoring flow
Finding/Hypothesis/Case/Incident/Hunt/technical handoff → Intake → Project → Detection Hypothesis → Data Readiness → Schema/Fields → Draft → Conditions/Sequence/Enrichment → Metadata → Validation → Test Scenario → Expected Outcome → Historical Replay → Match Review → Coverage/Gap → Provenance/Review Package.

## Transition guarantees
Every transition transmits source, trigger, destination, ownership, tenant/environment, version, restrictions, errors, context, return origin and provenance. It never expands permissions or creates runtime authority.

## Future-only transition
CAP-INV-417 may submit a package to Phase 4B.3A.2. No Approval, promotion, deployment or activation occurs here.
