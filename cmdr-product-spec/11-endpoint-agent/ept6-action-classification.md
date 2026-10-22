---
id: endpoint-ept6-action-classification
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
open_decisions: [OPEN-008, OPEN-013]
---
# EPT-6 Action Classification

## Class 0 — inspect
- inspect assignment/current version/package context;
- inspect download/stage/update state;
- inspect buffer/replay/recovery/degradation state;
- inspect self-protection/privilege/security-state;
- inspect local audit/provenance.

## Class 1 — deterministic assessment
- compatibility/readiness;
- post-update health/verification;
- queue freshness/pressure;
- replay/duplicate candidate assessment;
- dependency/capability reassessment;
- privilege/security-state normalization.

## Class 2 — not globally assumed
OPEN-013 remains open. Bounded refresh, diagnostics, replay/status or preparatory operations may be Class 2 only where the canonical source/policy permits. EPT-6 does not establish a universal Class-2 default.

## Class 3 — effectful local operation candidates
Where policy/authority requires, significant local effects include update install/activation, previous-version reversion, self-protection mutation if separately sourced, and any disable/uninstall capability if a future source introduces it. EPT-6 does not grant authority by classification.

## Class 4
No autonomous destructive action and no provenance destruction are introduced.

## Invariants
Read != execute; privilege != permission; assignment != authorization; retry != renewed authority; update reversion != Govern rollback; local recovery != response rollback. Any effectful action rechecks tenant, target, policy, permission, authority where applicable and current technical preconditions.