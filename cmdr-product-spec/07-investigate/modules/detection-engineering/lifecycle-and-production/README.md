---
id: investigate-detection-lifecycle-production
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
  - OPEN-017
---
# Detection lifecycle and production

This submodule consumes CAP-INV-401..417 and defines CAP-INV-418..435. It begins with the Review Package and ends with a Continuous Improvement Package.

It covers review, Release Candidate, readiness, promotion planning, Govern handoff, shadow, canary, deployment coordination, version reconciliation, health, Command feedback, production review, tuning, bounded suppression/exception proposals, drift, performance, rollback, retirement and lifecycle provenance.

It does not define a rule language, engine, adapter, API, package format, deployment command, storage model or implementation. It performs no real production change.
