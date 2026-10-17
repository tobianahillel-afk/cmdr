# E10-INV-120A — real-owner text hunt and dependable analyst pivot

## Functional delta

PR #42 now uses the existing authenticated Case owner/Timeline/source-ACL read flow for text-hunt results, and makes the pivot faithfully reproduce the **unfiltered** cross-Case search: selecting a matching Case clears incompatible existing event-kind and local time-range filters before setting the searched text and submitting the existing Case-specific form. The server still independently checks access on that fresh request.

The previous feature grouped identical `source_ref` values across Cases as if a canonical source object could belong to several Cases. The actual Go memory provider binds each source identity to one Case, so that grouping was unreachable on real valid owner snapshots and risked misleading analyst expectations. The feature now reports the **cooccurrence of the searched text** across separately authorized Cases, with total observed and inferred counts, and no fabricated cross-Case source identity.

A new Go integration test writes an actual two-Case operator owner dataset to disk, loads it, then uses live `httptest` HTTP routes to confirm that the same text is independently searchable in two authorized Cases, with **distinct typed source owners**. Denied source text, cross-tenant Case, anonymous access and cross-Case source preview remain unavailable. The browser VM tests check the pivot clears stale time/kind filters and uses the existing authenticated search (not a new privileged correlation API).

## Reuse vs alternatives

Chosen: existing Go memory owner provider, case discovery, paginated Timeline, typed-source authorization, browser `textContent`, standard DOM events, Go `httptest` and Node's built-in test runner. Rejected: pretending source-ID equality models a global IOC, a new privileged cross-tenant join API, extra JS framework or commercial SIEM connector. True normalized indicator matching and durable provider ingestion remain open architectural/production work, not a string-substring match.

## Limits and handoff

- 8 authorized Cases, at most 5000 visible search matches each, up to 200 paged requests, no atomic snapshot and no claim of causal or security attribution.
- The demo principal and operator-local owner file are synthetic, not production IAM or trusted ingestion. No real forensic chain of custody, source index, or retention guarantees.
- L2 must review this **functional** PR #42 snapshot in read-only mode. GitHub CI must pass at the immutable commit; no approval is inferred from a past SHA.
- No changes to canonical Product Spec or `OPEN-013`, no force-push, no merge into base or main. The parent E10-INV-120A-SNAPSHOT stays READY pending reconciled acceptance.
