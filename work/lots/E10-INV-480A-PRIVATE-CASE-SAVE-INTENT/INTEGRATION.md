# E10-INV-480A — Preserve explicit private Case output intent

## Material privacy boundary fixed

An operator runs the **real** first-party offline Case importer and deliberately selects a private output file via `-save-import-snapshot`. Before this change, an explicitly empty destination (common from an unset environment variable in a shell, e.g. `-save-import-snapshot "$OUT"`) was treated as if the option had never been supplied. A valid `-import-jsonl` and `-into-case` would instead fall through to the **legacy stdout JSONL exporter**, disclosing sensitive rendered normalized source previews to stdout/log collection. The privacy controls protecting HTTP Case searches do not protect this CLI fallback.

The same intent ambiguity affected `-save-wazuh-case-snapshot` and `-save-mixed-case-snapshot` (their unsupported fallbacks currently fail but must never silently select an unexpected mode). The prior independent PR #88 had already corrected the closely related `-private-report ""` and `-expect-snapshot-sha256 ""` ambiguity. This change **extends the SAME single CLI input trust-boundary validation to all five sensitive options**, rather than adding validators to each downstream import, storage and search function.

### Minimal implementation

Immediately after the existing Go FlagSet parse and argument verification, use the already existing `FlagSet.Visit` for **explicitly supplied** `-private-report`, `-save-import-snapshot`, `-save-wazuh-case-snapshot`, `-save-mixed-case-snapshot`, and `-expect-snapshot-sha256`. Any explicitly empty string (either `-flag ""` or `-flag=`) fails before ALL dispatch, stdout, filesystem import, HTTP listener, or bearer generation. An option genuinely omitted remains backward compatible; an explicit nonempty private path uses exactly the previous first-party private atomic no-overwrite Case writer and selected principal/tenant/Case/source owner.

### Evidence and limits

- Real, valid local normalized JSONL with a distinctive private message, complete Case snapshot and selected Case. Explicit empty private save path previously downgraded to a stdout full source export; **now both empty flag spellings must fail with exactly zero stdout bytes**.
- Both other protected save flags (Wazuh and mixed source) reject empty intent at the same trust boundary.
- A real nonempty private output still publishes one readable owner snapshot, prints only filepath/sha256 receipt, and never prints source preview.
- Old tests for explicit empty `private-report`, SHA256, owner-scoped Case, authorized Wazuh and generic mixed Event Search Job, no-overwrite private publisher and cross-tenant isolation run unchanged.
- No new IAM, new parser, source backend, data store, package, paid API or security-policy threshold change. Synthetic local lab only; output remains sensitive to operator filesystem access permissions.

## Branch / gate

Child of **immutable previously L2-reviewed PR #88** HEAD `4c46b644e95a65f09057e80508aea437a8832c16`. Parent PR #88 failed `engineering-kernel` mandatory changed security-critical 90% floor at 88.74%; no override, waive or premature merge. Full exact SHA GitHub Actions Linux Go+Node and macOS/Windows, engineering kernel and **one** targeted independent read-only L2 review of the new CLl privacy boundary are required. This issue is R1-level confidentiality leakage under a plausible shell-empty private destination; do not treat it as cosmetic.
