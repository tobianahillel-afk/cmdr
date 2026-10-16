# E10-INV-110B — negative-zero offset retains forensic uncertainty

This bounded read-only Go projection now preserves local-offset uncertainty if either the original RFC3339 source timestamp ends with `-00:00` or the explicit source offset field is `-00:00`. The Go standard-library `time.Parse` resolves the numeric value to a zero-second offset, so a numerical offset comparison alone loses the unknown-offset distinction.

Use the **existing** `strings.HasSuffix` and first-party `time.Parse` integration. No vendor temporal engine, timezone database, new dependency, persisted schema or migration. Preserve the original timestamp and label exactly; move uncertain entries behind those with known instants in the existing deterministic sort. RFC 3339 section 4.3 originally identified `-00:00` as unknown local offset; RFC 9557 updates the preferred notation to `Z`. This patch avoids upgrading the legacy negative-zero convention while retaining the existing CMDR contract of an **explicit trusted timezone/offset label** for `Z` timestamps.

References:
- https://www.rfc-editor.org/rfc/rfc3339.html#section-4.3
- https://www.rfc-editor.org/rfc/rfc9557.html#section-2.2

New tests separately assert legacy negative-zero input and metadata, positive-zero/explicit Z baselines, unchanged forensic source strings, and stable position of unknown entries. No production ACL-provider integration or end-to-end production chronology claim; `OPEN-013` remains unresolved. Status is pending exact-SHA CI and independent L2 assessment.
