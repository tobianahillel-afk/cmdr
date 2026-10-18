# E10-INV-490A — Consolidated mixed-source investigation and fail-closed CLI mode selection

## Material consolidation (not another engine)

This branch descends directly from the frozen #89 HEAD `250c0e27a71dc92adfe2117d07e7ece7b4474d61`, which itself descends from the frozen #81 HEAD `c173c34dbe224b6427cfb37ed880da9266504e1c`. **PR target is #81**, allowing a single real integrated change set instead of trying to merge the individually red and L2-frozen #82–#89 staged PRs one by one. Their existing native reviews, previous coverage failures, and immutable SHAs are NOT waived and do not become accepted by changing PR ancestry.

This candidate retains and integrates the already implemented:
- Operator-selected Wazuh 4.x and normalized JSONL ingestion into a single source-authorized **synthetic local** Case, saved only as a principal-scoped 0600 private snapshot with no-overwrite atomic publication and operator-retained SHA256.
- The same first-party Event Search Job reachable via analyst browser and headless in-process authenticated POST; source authorization **before** text predicate, bounded 500 candidate / 200 hits and known-time temporal Wazuh↔generic leads from authorized result references.
- Authorized offline saved-Case inventory, independent same-tenant/environment Case comparison and completed private report publication (no stdout sensitive refs when explicit private output is requested). These are runtime tested but NOT production IAM, live SIEM ingestion or cryptographically certified custody.

## Root-cause security correction

A sensitive CLI mode selector can be provided *explicitly empty* due to environment-variable expansion. Previously an empty `-export-case`, `-report-case`, `-resume-import-search-case`, `-source-case`, `-inspect-case` or `-snapshot` could be indistinguishable from the option being omitted and **fall through to the default loopback viewer or default synthetic fixture**. An empty private destination had a similar issue fixed upstream in #88 and #89.

The existing **single** `FlagSet.Visit` entry-boundary check is now consistently used for sensitive destinations, source-owner snapshot selection, source/Case export and investigation mode selectors, and snapshot digest flags. Any explicitly supplied **empty string** rejects BEFORE the CLI loads a source, starts a listener, generates a bearer, opens output or sends an event JSONL export. Deliberately **omitted** options continue the existing legacy behavior; nonempty operator intent is unchanged. No new validator class, deep flags DSL, conditional state machine, external package or paid API.

## Real behavior and security acceptance

- Representative Case export, report, owner source, inspection, comparison and Case resume flags are tested in both `-flag ""` and `-flag=` syntaxes with a valid authorized owner snapshot. Immediate mode error; stdout empty (no source preview/bearer).
- Explicitly empty `-snapshot` cannot silently use the built-in demo dataset to prepare an unrelated Case report. A nonempty `-report-case` still successfully performs first-party authorization/reporting.
- Existing Linux Go + Node analyst UI and Windows/macOS Go workflows are required to pass on the **exact SHA**, together with `engineering-kernel` and the unchanged **90% mandatory changed security-critical code coverage** floor; previous #88 (88.74%) FAIL and #89 (targeted changed-scope) PASS are not substitutes for this *consolidated* diff check.
- The full diff #81..HEAD is NOT considered production-ready merely because code tests succeed: all branch reviews, Native GitHub rules, parent specs/engineering foundation #2/#36, and architectural decisions remain gating.
- One READ-ONLY independent L2 only after the full consolidated candidate stabilizes; do not duplicate existing L2 on #87–#89 and do not change their frozen SHAs.

## Decision

Prefer a **single coherent integrated Release PR** spanning previous Wazuh+generic Case/Search/temporal correlation, authorization, independent Case comparisons, and private reports to merging each subcomponent in separate, red and dependency-heavy PRs. Do not cherry-pick duplicated alert correlation implementations (notably independent #51) without reconciliation. Do not force merge any parent or reduce a threshold to declare this release ready. No paid services.
