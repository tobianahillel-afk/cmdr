---
id: investigate-detection-engineering-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Workflows

## Authoring flow — 4B.3A.1
Finding/Hypothesis/Case/Incident/Hunt/technical handoff → Intake → Project → Detection Hypothesis → Data Readiness → Schema/Fields → Draft → Logic/Enrichment/Metadata → Validation → Test Scenario → Expected Outcome → Historical Replay → Match Review → Coverage/Gap → Review Package.

## Lifecycle flow — 4B.3A.2
Review Package → Release Candidate → Detection Review → Deployment Readiness → Promotion Plan → Detection Change Request Draft → Govern Action Request → Decision/Approval → Response Run/Result → Runtime Version Reconciliation → Health/Quality/Production Review → Tuning, Suppression/Exception, Drift/Performance, Rollback or Retirement → Lifecycle Provenance → Continuous Improvement Package → new CAP-INV-401/402/403/406 cycle.

## Shadow and canary branches
- Release Candidate → authorized Shadow Evaluation → non-alerting Shadow Assessment → Review.
- Release Candidate → Canary Plan → Govern path → per-target Canary Assessment → promote, stop or prepare rollback.

## Guarantees
Every transition transmits source, trigger, destination, owner, tenant/environment/target, immutable version, restrictions, authority, errors, return origin and provenance. It never expands permission, invents active state or collapses partial target results.
