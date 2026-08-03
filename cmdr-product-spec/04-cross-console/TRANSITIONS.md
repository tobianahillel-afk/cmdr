# Cross-Console Transitions

## Command Center → Investigation Lab

Trigger: an incident needs evidence-led analysis. The transition creates or opens a Case, carries Incident linkage, tenant, time window, selected entities, alerts and operational summary. Command remains the source of operational status; Investigation owns analytical work.

## Investigation Lab → Response & Governance

Trigger: a confirmed Finding or documented emergency basis supports an action recommendation. The transition creates a Response request containing action, targets, justification, evidence package, confidence, expected impact, risk of action, risk of inaction, proposed conditions and rollback plan.

## Response & Governance → Command Center

Trigger: a Decision or Run changes the operational situation. Command receives decision outcome, conditions, execution state, rollback availability, affected entities, business impact and residual risk. It does not duplicate the full governance record.

## Invariants

- Source and destination objects are bidirectionally linked.
- Context is preserved but authority boundaries are not.
- Transitions are auditable and idempotent.
- Re-entering an existing destination object is preferred to creating duplicates.
- Users can return to the exact source context.
