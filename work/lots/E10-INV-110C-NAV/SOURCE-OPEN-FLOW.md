# E10-INV-110C-NAV — Open a synthetic owner source from the Case Timeline

## Delivery

The existing flow already provided authenticated Case discovery, filtered Timeline entries and a source-navigation draft with per-destination authorization. This increment closes the last step of the **loopback-only demonstration** without inventing a production owner store:

1. The analyst obtains Case and Timeline entries through the existing session, catalog, owner and typed source permissions.
2. The analyst chooses an authorized Timeline Entry; `/case-source-navigation` rechecks access and returns a non-authoritative owner draft plus presentation-only return selection.
3. The analyst browser calls a **separate** `GET /demo-source-open?case_id=...&entry_id=...`, using the ephemeral bearer again. The owner handler independently re-evaluates Case, selected Timeline Entry and canonical typed source read, and retrieves a read-only source preview from the same process-local owner snapshot.
4. The browser checks that the opened object's Case, Entry, kind and canonical ID match the prior draft. It displays the returned preview with `textContent`, without a redirect, HTML interpretation, source-ref construction, new permission token, credential storage or mutation.
5. Returning to the Case reruns the timeline API authorization; copying the local return link carries only Case selection and filters, never the bearer.

## Technology choice

Reuse the existing Go `MemoryTimelineData` owner adapter, `DraftCaseSourceNavigation`, `AccessChecker`, canonical typed source permission map, `net/http` and the same-origin browser UI. No new generic router, custom auth scheme, database, source index, JavaScript dependency, API vendor or external authorization provider. The synthetic `Preview` belongs to `MemoryTimelineSource` and is **not copied into the Shared Timeline Entry or Timeline projection**.

## Security and acceptance

The new owner-open endpoint accepts only bounded `case_id` and `entry_id`, never `object_id`, `tenant`, `source_ref`, a return URL, or caller ACL flags. It rejects unauthenticated requests, denied/hidden or cross-tenant sources with neutral 404 and returns no private identifiers in error bodies. It rechecks the typed owner's canonical read permission immediately before returning its static in-memory preview; responses are `no-store` and vary on `Authorization`. Preview data is validated on construction (valid UTF-8; at most 2048 bytes), returned as JSON, and rendered in the browser only as inert text.

Tests: Go HTTP end-to-end for both demo tenants, authorized Event/Evidence/Incident owner previews and denial paths; construction bounds and owner authorization tests; Node DOM test for the actual two-request destination draft -> owner opening plus denied opening and safe return state.

## Strict limitations

This is a **runnable functional test/demonstration**, not real-world source ingestion or an integrated persistent EDR/SIEM backend. The bearer and owner data are synthetic; static reference data and permissions never change after construction. It does NOT implement durable Case or Timeline storage, vendor connectors, audit ingestion, TLS ingress, production IAM, a final UI, a production SLO or a cross-product source router. A real host must independently wire the owner application and enforce authorization on every source open. `OPEN-013` and the immutable Product Spec remain unchanged.

Do not mark all of CAP-INV-110 delivered merely because this demo journey passes; C-NAV can be verified within this **bounded scope** only after exact-SHA CI and critical negative tests. A distinct independent L2 must inspect the published functional milestone, not any in-chat-only draft.
