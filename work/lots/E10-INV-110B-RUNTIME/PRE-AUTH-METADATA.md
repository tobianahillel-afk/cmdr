# E10-INV-110B-RUNTIME — lossless metadata before source authorization

## Finding and minimal correction

The existing owner-resolver and legacy access-checker adapters call a source-provider only after `validTimelineEntryBeforeSourceLookup` has checked entry structure. Following the published UTF-8 defense at `fad5b087bc57ff69d44dd933b3ac9eb921da0893`, the final `Project` rejects malformed UTF-8 in the visible `source_time`, `timezone`, `provenance` and `supersedes` strings. However, that same check was absent from the **pre-source-provider** gate; an invalid entry could cause an unnecessary external authorization call before eventually being denied.

Reuse the existing validator and Go standard-library `unicode/utf8.ValidString`, rather than adding a duplicate adapter, external JSON framework or new dependency. Denied, unknown, missing and foreign Timeline Entry snapshots remain opaque and must never be resolved to inspect their metadata. Parseable timestamps are not mandated; unknown timestamps remain uncertain.

## Proof and limits

Unit tests verify that all four malformed fields fail closed *after* Case and Timeline Entry authorization, but *before* the third, external-source permission request, for both owner-resolved and caller-snapshot paths. Tests also assert denied/missing entries never reach the source provider and legitimate multi-script provenance and unknown clocks remain unchanged.

This is an in-process trust-boundary correction, not an implemented production AccessChecker, EntryResolver or source provider; it adds no storage, network service, canonical schema or SLO. `OPEN-013` remains open; B stays READY pending exact-SHA CI and independent review. The previous 4.700ms synthetic performance measurement belongs to its older SHA and must not be reattributed to this change.
