# E10-INV-120A-SNAPSHOT — opt-in local forensic summary report

## Functional value and reuse-first choice

The operator may now deliberately save an **unsigned, local JSON investigation summary** after an authenticated full read of the existing Case timeline pages. This bridges the local Case picker → verified tenant/source ACL page fetches → temporal and source correlation → portable human-readable triage output. Reuses the existing browser summary (counts, known timestamps and top 3 authorized source groups) and native Blob/download support; **no** new endpoint, database, IAM system, frontend SDK, secondary authorization logic or event-index engine.

## Trust and data-lifecycle boundaries

- The download button starts disabled; it is enabled **only** after every authorized page succeeds with consistent totals, distinct IDs and acceptable source/time metadata.
- On token, Case, kind, search or time-window change, or on failed/aborted/revoked/partial aggregate, the previous report is cleared and download disabled. The browser makes no persistence until a click.
- The report includes selection, creation timestamp, per-kind counts, uncertain time count, certain time span, maximum *visible* gap, number of unique source refs and only top three authorized source groups.
- It excludes raw event rows, typed source owner previews, ephemeral bearer token, URLs and unknown/unreadable sources. Saving is an **explicit operator decision** and the local file must be handled as sensitive; it may contain authorized source identifiers.
- The snapshot is synthetic and owner-selected at CLI startup. Reports are not signed, not cryptographic custody proofs, not transactionally isolated across pages, not an audited export, not a production IAM/persistence or evidence store. Do **not** run against customer production data. OPEN-013 is unchanged.

## Tests and acceptance

Node VM browser tests exercise 205 authorized events across two requests, report schema and counts, absence of raw entry IDs/bearer/URLs, download only on explicit click, Blob disposal, clearing on token/Case change, and failure to export incomplete pagination. Existing Go integration confirms real disk JSON → Case/source ACL HTTP; new UI files remain covered by the embedded Go build and Go-run Node suite when Node is installed. Exact GitHub Actions and any independent L2 must be checked at the immutable published SHA before declaring PASS.

Next useful work: decide a production-grade principal/owner/persistence integration independently of this offline demonstration; avoid pretending that more local-JSON variants close that gap.
