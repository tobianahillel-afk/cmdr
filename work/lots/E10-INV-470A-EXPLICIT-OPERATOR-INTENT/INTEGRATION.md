# E10-INV-470A — Fail-closed operator intent for private forensic output and byte verification

## Material security flaw in the consolidated investigator delivery

The integration PR #87 combines three real Case investigation paths: bounded imported source inventory, headless first-party Wazuh+generic Event Search/temporal leads, and independent saved Case comparison. All accept an optional `-private-report ABSOLUTE_NEW_PATH`, and saved Case resume supports optional `-expect-snapshot-sha256 DIGEST`.

Previously, the CLI interpreted the **empty string** as if the option had not been supplied, even when the operator explicitly supplied `-private-report ""` or `-private-report=`. A shell environment variable unintentionally unset or empty could therefore silently **print an authorized but sensitive investigation report to stdout** instead of publishing it privately. Similarly `-resume-import-search-case ... -expect-snapshot-sha256 ""` could silently become an **unchecked** Case resume that launches a browser listener and issues a bearer. The former is a material disclosure/operational safety risk; the latter silently removes an operator-requested integrity check.

The correction is intentionally located at the ONE CLI trust boundary: after `flag.Parse`, inspect `flags.Visit` for **explicit operator presence** of either security option and refuse any explicitly present but empty value. An ABSENT optional flag retains its current documented legacy behavior; a NONEMPTY option retains exactly the original owner-authorization, hash-checking, and atomic private-file logic. This adds no new authorization layer, duplicate validator subsystem, dependency, database, network listener, or paid API.

## Acceptance and limitations

- Actual Case owners containing mixed Wazuh+normalized JSONL sources exercise **all three** originally integrated authorized CLI report modes. Explicit `-private-report ""` or `-private-report=` must return an error **before any stdout bytes**. Previously, normal modes fell back to stdout. Existing behavior with a valid path must remain unchanged.
- Actual saved Case restart, offline inventory and headless hunt must reject explicitly present `-expect-snapshot-sha256 ""` (also the equals syntax) **before HTTP listener, bearer or source result emission**.
- Operators who deliberately OMIT `-private-report` continue to receive the intentionally public CLI stdout receipt; operators who omit SHA on the optional resume path continue their explicitly backward-compatible unverified synthetic session. The new offline verify and hunt modes still require nonempty SHA as before.
- A positive test validates an actual authorized Case with the correct operator-held digest and a private report destination, preserving report production and no stdout.
- This does **not** create production tenant IAM, signed forensic custody, trusted timestamp evidence, a transactional case store or continuous ingestion. It prevents accidental privacy/integrity downgrades caused specifically by explicitly empty operator flags.

## Branch and security gates

Release-candidate branch forks **immutable PR #87** SHA `0d70dd9eacd31e21f43cdf6ab51d72622b02e73d` without changing it. Integration target remains the immutably reviewed **PR #81** SHA `c173c34dbe224b6427cfb37ed880da9266504e1c`, not main, for a consolidated, full-scope security gate. The older #87 broad-scope `engineering-kernel` failed at **88.60% < mandatory 90.00%**; this new candidate MUST satisfy its own exact-HEAD Go+Node Linux/macOS/Windows tests and unmodified mandatory 90% changed security-critical floor, plus native review and upstream release requirements. No merge while red and no waiver/threshold change. L2 should review only the material new CLI trust-boundary delta and its release implications, not repeat all prior independent reviews.
