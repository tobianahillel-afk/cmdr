# E10-INV-120A — two authorized Cases in one tenant

Functional enhancement to the operator-selected offline Case journey: a synthetic `investigator-42` now has **two** independently owned Cases in tenant A, and `investigator-99` sees a separate restricted Case in tenant B. The `-case=case-incident-43` startup choice selects a Case that is NOT the default first Case. A real loopback integration test follows that selected Case to its own chronological entry, typed source navigation and visible source preview. The original Case remains separately accessible using the Case picker; tenant B and unauthenticated requests remain denied.

This adds **only synthetic fixtures and targeted end-to-end regression tests**, not a new provider/permission layer, workflow engine or external dependency. Go's existing Case reader ACLs, typed-source checks, JSON snapshot loader and HTTP handlers are reused. The negative startup test now removes all owner Cases for the principal, so it remains a falsifiable rejection after adding a second owned Case.

Evidence must be tied to the exact commit's completed CI; no production SLO, authoritative identity, durable case store, real malware sample or evidence preservation is claimed. L2 remains one review of the materially integrated functional lot; no per-commit task.
