# E10-INV-120A — functional integration: operator-defined Case analyst

## Why this increment matters

The prior on-disk owner snapshot already fed a complete loopback Case picker and Timeline, but the CLI still accepted only the two hard-coded demo principal IDs. This prevented an operator-supplied owner dataset from selecting its own synthetic analyst, even when it represented the correct tenant, Case membership and typed source ACLs.

The changed CLI now accepts a synthetic principal **only when the supplied local owner dataset identifies that principal as a reader of a Case**. The existing `NewMemoryCaseTimelineHTTPHandlers` validates the full owner snapshot (duplicate principals, tenant consistency, source ownership, Case/Timeline/source ACL) before the listener starts. It still takes NO principal or permission from a network request. Without a supplied snapshot, the built-in two-principal fixture remains the only possible authority.

Operator-defined Case IDs are percent-escaped when embedded in the console URL, and operator-controlled IDs are quoted in terminal informational text. A live loopback HTTP integration test covers an owner snapshot with a nondefault analyst and a Case ID containing spaces, `&` and `#`, exercising authenticated Case discovery, Timeline, navigation, source preview, anonymous denial, foreign-tenant denial, denied-source suppression and cancellation/shutdown.

**Reuse-first:** use the existing in-memory owner provider, Case/Timeline/source API, Go `net/url` and `encoding/json`. No new framework, storage engine, IAM, network provider, browser-side authorization or write path.

**Strict limits:** the snapshot remains an operator-selected local JSON demo and the selected analyst remains a **synthetic claim**, not a production authentication identity. The bearer only protects the local loopback process. No live input feed, durable Case repository, production source adapter, audit backend, multi-user authorization integration, or Shared Timeline Entry mutation is claimed. `OPEN-013` remains unresolved. Proof requires exact-SHA GitHub CI and a single risk-based independent L2 for this functional lot.

The HTTP integration fixture also changes the source owner's preview to a unique string present **only in the on-disk JSON** and verifies it through the authenticated `/demo-source-open` route. This distinguishes actual disk-to-source data integration from a test that accidentally falls back to the built-in Go fixture. The Go JSON encoder may escape `&` as `\\u0026`; Case discovery assertions decode the JSON rather than comparing its raw representation.
