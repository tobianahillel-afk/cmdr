# End-to-End Operating Model

## Operating sequence

1. Signals enter through integrations and become searchable events, alerts or observations.
2. Command Center groups and prioritises operational work into incidents and queues.
3. An analyst opens or links an Investigation Lab case when deeper reasoning is required.
4. Investigators collect, preserve and analyse evidence; they maintain timelines, entities and hypotheses.
5. Confirmed analysis becomes a finding with confidence, provenance and supporting evidence.
6. A recommendation that requires controlled action creates a response request.
7. Response & Governance evaluates impact, policy gates, authority and rollback.
8. An authorised decision launches a playbook or action run.
9. Execution telemetry, failures and rollback status are audited and returned to the incident.
10. Command Center updates status, residual risk, business impact and stakeholder reporting.

## Core invariants

- Every transition preserves provenance.
- A finding references evidence; a decision references findings or an explicit emergency basis.
- An execution references an approved decision except where a policy-defined pre-authorisation applies.
- Every mutation is attributable to a user or service identity.
- Tenant and legal-hold boundaries apply throughout the chain.
