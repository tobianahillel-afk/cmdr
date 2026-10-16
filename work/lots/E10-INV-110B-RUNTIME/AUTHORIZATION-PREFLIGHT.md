# E10-INV-110B — two-pass Timeline permission preflight

## Risk and chosen approach

The trusted owner path formerly checked an ID and immediately resolved its owner snapshot and linked source. If a later selected Timeline Entry authorization errored, the result was denied atomically but earlier owner/source queries had already happened. That created unnecessary external authorization/metadata I/O before discovering that the requested selection could not be fully authorized.

Reuse existing `AccessChecker`, immutable `AccessRequest`, `EntryResolver` and standard Go slices; perform two bounded passes: (1) Case permission then **all selected Timeline Entry access checks**, recording only the four allowed decision labels; fail atomically on any provider error/invalid label/cancellation; (2) resolve and check source for *allow* entries only, preserving selection order and producing opaque filtered/missing entries for other outcomes. No extra round-trip per ID, backend, store, schema, ACL cache, dependency or paid API.

Compared alternatives: caching external permissions or batching requires a provider-specific consistency/security contract, so it was not added. Index/event-sourcing layers cannot solve permission ordering. This change does **not** provide an atomic ACL snapshot, principal validation, Case membership, source ownership mapping or proof of the external authorizer; the injected trusted providers remain responsible. `OPEN-013` is unchanged and B remains READY.

Tests assert exact Case → all Timeline rights → owner/source order with two allowed IDs; a provider failure on the second Timeline right must cause **zero** owner resolutions/source checks; a denied later ID must remain opaque without source lookup. Existing single-entry, denied and source-budget tests remain applicable.

GitHub CI and separate L2 verdict must be checked on the immutable published SHA before qualifying this gate; historical measurements cannot be reused as performance evidence for this changed path.
