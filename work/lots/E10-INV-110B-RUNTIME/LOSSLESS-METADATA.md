# E10-INV-110B-RUNTIME — lossless forensic metadata boundary

## Scope and technology choice

A Go JSON encoder replaces malformed UTF-8 bytes in string fields with the Unicode replacement character, which could silently alter an observed source timestamp, declared timezone, provenance string or prior-version reference. Because these are forensic *evidence descriptors*, the bounded first-party projection now rejects **authorized visible** entries that contain malformed UTF-8 in those four metadata fields. Existing `invalidScope` already rejects malformed tenant/environment/entry ID/linked-source selectors. No normalization, character replacement, third-party Unicode library, database, or new provider is introduced.

Reuse: Go stdlib `unicode/utf8.ValidString` is a direct predicate; it is simpler than introducing an external encoding library or canonicalizing source evidence. This is deliberately not a general character whitelist: valid multi-script Unicode and parseable or unknown-but-UTF-8 timestamps keep their original bytes and uncertain status.

## Security tests

Negative examples inject raw 0xff into each permitted metadata field; every projected result must be atomic deny with zero entries, never a lossy JSON replacement. Complementary negative cases show that the same malformed metadata in source-denied, unknown, missing and foreign-tenant entries does not influence an otherwise valid projection. A positive example retains the exact multi-script provenance and unknown-clock strings unchanged and reports time uncertainty.

The pure `Project` function is not an authoritative permission engine. Trusted-owner and legacy adapters must supply actual permission outcomes; no production provider is bundled. This changes neither Product Spec nor OPEN-013. The work graph remains B READY until independent L2, exact-SHA CI and final release gates are reconciled. Historical on-demand synthetic performance measurements belong to their own SHA and are not re-attributed here.
