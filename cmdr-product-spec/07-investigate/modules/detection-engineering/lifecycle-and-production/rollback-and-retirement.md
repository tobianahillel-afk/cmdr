---
id: investigate-detection-rollback-retirement
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Rollback, deactivation, retirement and replacement

CAP-INV-433 prepares rollback and recovery evidence; CAP-INV-434 prepares deactivation, retirement and replacement. Investigate never executes the runtime change.

Rollback planning names the return version, targets, authority, validations, risks, preserved data and recovery checks. `rollback-requested`, `rolled-back` and `full-recovery` remain distinct. Per-target partiality is never collapsed.

Retirement planning names consumers, historical Signals/Alerts/Incidents, dependencies, affected coverage, gaps, replacement candidates, coexistence, transition and rollback. `inactive`, `retired` and `deleted` remain distinct. Historical versions and provenance are preserved; replacement never implies equivalent coverage.
