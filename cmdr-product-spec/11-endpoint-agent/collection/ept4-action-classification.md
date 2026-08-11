---
id: endpoint-ept4-action-classification
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-4 Action Classification

## Class 0 — inspect
Inspect Collection Request/plan/progress/output/session/execution state/provenance. No target effect.

## Class 1 — deterministic no-effect assessment
Eligibility, scope validation, completeness calculation, transfer-status validation, provenance reconstruction and readiness checks.

## Class 2 — bounded technical controls
Potentially bounded low-impact collection; read-only fresh snapshot; open authorized diagnostic session; non-mutating diagnostic invocation; request cancel/stop/reconnect/close. Actual classification depends on impact/policy and preserves OPEN-013.

## Class 3 — effectful/sensitive execution
Effectful command/script, temporary file upload/collision override, invasive/high-impact acquisition or other target mutation requires Govern authority context. EPT-4 can execute the technical primitive only after bounded authority; it does not create Decision/Response Run/Result.

## Class 4
EPT-4 creates no autonomous irreversible/destructive authority. Destructive remediation/containment is out of scope and stops at EPT-5/Govern.

## Mandatory distinctions
Eligible != authorized != started; cancel != rollback; stop request != stopped; session close != state restoration; technical success != Response Run success; technical output != Result; operator != Govern approver.