# E10-CMD-106M-RESULTREAD — bounded Govern Result consumption in Incident Detail

## Canonical behavior, not a new Result engine

The source of truth for Result is Govern (OBJ-RESULT and CAP-GOV-032). CAP-CMD-106 Command Incident Coordination and CMD-IWQ-006 Incident Detail are read-only consumers. A Result is a classification/provenance from Govern, NOT raw executor output, a Decision, confirmed finding or a Command Incident transition. CAP-GOV-032 permits Class-0 inspect with separately authorized `perm.govern.result.read`; OPEN-013/015 block mutations/finalization/provenance design outside this scope.

## Four owner boundaries

1. **Product:** Incident Command, Result Govern, object linking Shared; Work Queue selected Incident CAP-CMD-101. No duplicate Govern owner or policy engine in work queue.
2. **Engineering:** existing E10-CMD-106 epic; child E10-CMD-106M-RESULTREAD depends on L-ALERTSRC IMPLEMENTED (not ACCEPTED). No READY parent assumed; this is an explicitly bounded new Class-0 child.
3. **Code:** existing first-party Go Command work queue selected Incident Detail, narrow optional `IncidentResultReferencesOwner` and `GovernResultOwner` interfaces; existing local-browser synthetic host supplies two **separate** objects. No new URL, store or runtime dependency.
4. **Git:** new branch `feature/e10-cmd-106m-govern-result-readonly` starts from immutable PR #211 SHA `aabdf107824bee9c3561ccce440586ac1670fa3d`, whose independent L2 and ancestor chains remain unchanged. Git parent is not engineering or product parent.

## Real source-backed behavior in the local demonstrator

An authenticated GET `/command/incident-detail?incident=incident-a` first rechecks the Command selected Incident (Owner.List/CanRead/Project and existing selected-context owners) and Command Alert provenance; only then does the Shared relation source resolve at most eight Incident→Result references for that principal/tenant/env/Incident. Each Result is independently queried from Govern `ReadAuthorizedGovernResult`, which must itself enforce `perm.govern.result.read`; no browser-supplied Result ID, role or permissions. The returned narrow projection must match Result ID, tenant/env, valid Run ID, source/version and valid created-at; bounded lifecycle, outcome, verification and residual-risk are autoescaped via original html/template and existing strict headers. A verified state does not magically mean outcome success: synthetic example is `finalized`, `verification-failed`, `failed`, and residual risk unverified. A denied Result has no ID, count or existence in the HTML; missing source stays explicitly partial. Bad relation, duplicate, spoofed IDs, malformed text, owner error or cancellation fails before any HTML. Original 2-argument `NewIncidentDetailHTTPHandler` remains compatible.

## Proof and strict exclusions

New Go handler tests exercise positive HTTP plus denied/malformed/revoked/legacy and malformed timestamp cases, and real synthetic localhost Go HTTP tests the two Basic principals. This is *not* production IAM, remote Govern provider, persisted Response Run/Result, raw technical result access, Action Request, Decision, Case mutation, Incident state update, reporting, audit persistence, customer tenant isolation or fully accepted Incident Detail with all six UI states. This branch is not merged and remains Draft; CI must be verified at its exact HEAD. Any material new cross-owner trust decision needs independent L2 by a separate read-only context before acceptance. No API payante.
