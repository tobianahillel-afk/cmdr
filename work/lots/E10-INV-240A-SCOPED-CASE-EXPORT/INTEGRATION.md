# E10-INV-240A — Authorized Case-only durable export and reopen

## R1 root cause and functional resolution

The preceding draft PR #62 implemented explicit no-overwrite local snapshot persistence, but serialized the **entire operator-selected owner graph**. Although the read-only HTTP Case/Timeline/typed-source owner denied unauthorized reads, serialized bytes could include **other tenants, other Cases, and same-Case entries with denied source ACLs**. A file read by an operator would bypass HTTP authorization entirely. This is a material confidentiality defect, not a minor coverage gap.

This delivery replaces full-graph serialization with **one selected, explicitly authorized Case projection**. Before marshaling any bytes, the existing first-party owner validator verifies the source graph; the code resolves the selected principal's exact tenant and environment, requires the specific Case to list that reader, and constructs a fresh `MemoryTimelineData` containing only:
- the ONE selected principal, not other users or tenants;
- the ONE selected Case, not unrelated Cases, with reader set to that principal;
- only Case-catalog Timeline entries where BOTH the entry and its matching canonical typed Source grant access to that principal;
- only those typed Source owner objects referenced by accepted entries; all original denied or unrelated previews, names, metadata, and principals are omitted.

The resulting scoped Case/owner graph is independently revalidated using the **same** `newDemoHandlerWithData` first-party owner constructor. Saving still requires explicit CLI `-save-import-snapshot`, preserves the existing 1–4 normalized-feed identities and provenance, publishes only a complete bounded 0600 temporary file through atomic no-replace hard-link creation, and prints only a non-authenticated SHA256 digest/path.

**The reopened Case is also re-scoped in memory** before configuring the existing Case/Timeline/typed-source HTTP owner and `POST /demo-local-search`. This is defensive compatibility with old local operator snapshots, which could still contain other Cases; it does NOT retroactively erase them from old files. A new bearer is generated, and source reads are re-authorized through the same first-party owner. Saved principal reader lists are reduced, but the output is NOT a general trusted cross-tenant export authorization service: the input snapshot is operator-supplied synthetic lab data, not backed by production IAM.

A Case containing denied entries produces a precisely scoped **readable projection**, not a byte-for-byte export of the original Case. Inaccessible entry/source metadata is not saved. This must not be mistaken for a forensic full-evidence export or signed chain of custody.

## Security / acceptance

- Actual multi-tenant input owner, with a second principal/tenant/Case and sensitive source preview, AND same-tenant Case entry whose Timeline/source readers do not include the selected principal. Source graph remains valid.
- Explicit write to new private snapshot: saved bytes are scanned to confirm **zero** other-tenant identifiers, denied principal identifiers, Case/entry/source IDs or preview text. On-disk owner contains exactly 1 principal, 1 Case and only readable entries/sources. Canonical imported entry remains searchable.
- Existing PR #62 positive tests still exercise multi-file original ID collision, actual save/restart, recovered first-party Event Search Job, source drilldown, private POSIX mode, no overwrite/symlink and denied bearer; do not weaken source permissions to make tests green.
- Reopen of an older synthetic multitenant snapshot is narrowed to the selected authorized Case in memory; old on-disk operator copies remain a sensitive legacy risk and must be handled under operator file controls.
- No new IAM, new detector/validator subsystem, database, DSL, paid API, service daemon, or mutation of canonical OPEN-013 decisions. Simple dataset minimization at **the output trust boundary**.

## Merge / audit

Draft child of frozen PR #62 SHA `db639fb117fb490138b073fa916b3718f5102001`, no change to any earlier L2 audited SHA. Exact-head `investigation-timeline-runtime` Linux Go/Node plus macOS/Windows, and the unchanged mandatory `engineering-kernel` 90% changed security coverage floor must be verified; failed coverage blocks merge even if functional tests succeed. This is a material **new data-export security boundary** and requires ONE independent read-only L2 after verifiable CI. Synthetic lab only, no production access-control guarantee, durable DB or cryptographic evidence custody.
