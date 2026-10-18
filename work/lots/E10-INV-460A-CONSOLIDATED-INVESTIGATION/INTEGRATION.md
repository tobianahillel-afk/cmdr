# E10-INV-460A — Consolidated first-party mixed-source investigation release

## Why this integration exists

GitHub's immutable compare check verified `c173c34dbe224b6427cfb37ed880da9266504e1c` (PR #81 security-coverage diagnostics foundation) to `d960cf003c187d2e92b442da99b5425725e71cae` (PR #86 private reports) is a **linear 20-commit / 16-file delta, 0 commits behind**. PR #84 consolidated the actual functional code already present in separately frozen PRs #82/#83; PR #85 fixed silent truncated receipt acknowledgements; PR #86 added a new-only private file publisher for all three authorized analyst journeys. Parent #84 and #85 exact-head coverage gates were red although PR #86's exact HEAD passed both mandatory workflows; each parent is still an independent DRAFT with its own L2 and approval requirements.

This PR collects that **one coherent first-party analyst flow** into a single reviewable feature delivery against PR #81, rather than requiring release engineering to merge #82/#83/#84/#85/#86 separately and duplicate their components:

1. Offline operator selects Wazuh 4.x and generic normalized JSONL feeds; saves a minimal authorized Case as a private no-overwrite owner snapshot, including synthetic original source previews. The principal keeps a separate SHA-256 checksum.
2. From that exact archived Case, **offline source verification** reports the authorized Wazuh/generic event counts, with full Case/principal/typed-source scoping and independently checked file digest.
3. An analyst executes the **same first-party Event Search Job** as the interactive loopback UI, but headlessly, with no TCP listener and no bearer on stdout. Results contain only authorized source refs plus Wazuh 60s bursts and mixed Wazuh↔generic ±60s temporal context derived from a completed, bounded Job.
4. The analyst compares that **same Case** with a second independently fingerprint-verified, authorized Case in the same tenant/environment using the existing Timeline/Case comparison, without merging owner graphs, changing permissions, or copying raw source previews.
5. With optional `-private-report /absolute/NEW.json`, each of the three receipts is atomically saved to its OWN new private file (0600 on POSIX) rather than streaming sensitive refs to stdout. No existing snapshot/source/report can be overwritten or silently truncated.

**New integrated end-to-end behavior test** uses one real operator-local mixed Wazuh+generic saved Case as the SAME input to all three workflows, checks the exact Case identity, 2 Wazuh+1 generic sources, an actual SearchJob with 3 authorized refs and 2 mixed context leads, plus a separate independently authorized Case comparison. It verifies private on-disk JSON reports, no raw previews/bearer/stdout and no archived Case overwritten. This is a cross-workflow acceptance journey, rather than separate isolated unit tests or governance-only changes.

## Boundaries and alternatives

- **Chosen:** rebase the single already functional 20-commit linear ancestry by opening a consolidated PR directly against the frozen #81 base, with one meaningful cross-workflow acceptance test and this decision handoff. This avoids reimplementing the engines or cherry-picking trusted code.
- **Rejected:** a new search/correlation/comparison framework, new persisted database or indexing vendor, relaxing security coverage, force-pushing/rebasing any L2-frozen SHA, or merging a red parent.
- **Not production:** operator-controlled synthetic local Case owner, no production IAM, durable server database, trusted tenant identity, live Wazuh connector, signed forensic custody, proven time truth, or guaranteed source authenticity. Temporal proximity is NOT causality; a SHA-256 held by the same untrusted operator is not a signature.
- **GitHub reality:** `main` remains at `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`, the product specification/engineering PR #2/#36 remain DRAFT (PR #36 is 763 commits and 871 changed files). The execution-state cache is not authoritative. This consolidated PR does NOT automatically mark #82/#83/#84/#85/#86 delivered or close prior L2.
- All mandatory GitHub Actions, the existing **90% security-critical changed coverage requirement**, native repository approval and the relevant independent read-only L2 are required on the exact consolidation SHA. A green workflow on the descendant SHA never retroactively turns parent red runs green.

## Integration handoff / next priorities

Reconcile superseded parallel PRs by GitHub ancestry and reviewer approval before closing/retargeting; do not alter their immutable HEADs. Then break the huge product-spec/foundation integration (#2/#36) into reviewable deployable components while retaining canonical contracts, and implement real IAM/tenant isolation, persistence and ingestion only after the OPEN product decisions have been resolved. No API costs or local machine required for tests: rely on GitHub Actions.
