# E10-INV-110F-SELECTOR — Case picker / Timeline HTTP compatibility

## Canonical and engineering ownership

- **Product:** `CAP-INV-110` (Investigation Timeline, read-only Case/source navigation) is the code/behavior owner; `CAP-INV-101` (Case Queue) consumes the picker-to-workspace transition. Shared-owned Timeline Entries and Investigate-owned Cases stay unchanged.
- **Engineering:** independent bounded task `E10-INV-110F-SELECTOR`, child of verified foundation epic `E10-INV-110`; product spec untouched, registered as **IMPLEMENTED candidate** pending exact GitHub CI and security checks. Do not mark `VERIFIED` or full CAP acceptance automatically.
- **Code:** preexisting `product-runtime/investigation-timeline` owner-neutral HTTP handlers; no new service, new permission model, new data storage, paid API or mutation.
- **Git:** intentionally forked from frozen engineering foundation PR #36 SHA `ea1770b48a4bc1b2e850a7ffab15d5839f9e9e49`. This task is independent of the blocked PR #132 UI provenance run comparison; it MUST NOT mutate their audited HEADs. Relation to #36 is **Git ancestry**, not canonical business ownership.

## Reproducible defect and change

`NewCaseDiscoveryHTTPHandler` accepted a candidate owner-provider Case ID with a newline, carriage return or NUL, subject to an otherwise successful `CheckAccess`, and could advertise that Case to an analyst. `CaseTimelineHTTPHandler` correctly rejected the same ID as an invalid HTTP Case selector. A read-only analyst journey `Case Queue → select authorized Case → open Timeline` thus advertised an unusable/unsafe Case ID. The case selection trust boundaries disagreed. Invalid owner-provider candidate IDs should fail **before permission enumeration and before any partial Case list**, not leak a partial result.

Use **one simple `validHTTPCaseSelector` predicate** already equivalent to the Timeline's checks (non-empty/non-wildcard, <=1024 UTF-8 bytes, valid UTF-8, no NUL/CR/LF). Reuse it on both handlers without weakening the existing tenant/environment scope and independent Case and source authorization gates. A malformed candidate anywhere in a list produces the same generic 404 as an unavailable owner catalog; a direct malformed Timeline HTTP selector remains 400.

## Acceptance

- For a catalog `[case-one, invalid]` with invalid owner IDs carrying NUL/CR/LF, the Case picker returns generic 404, does not invoke ACL checks, never emits `case-one` or the malformed ID.
- Direct Timeline GET of each invalid ID fails with 400.
- A normal `case-one` still opens with status 200 via the existing owner authorization gate.
- Existing Case discovery anti-enumeration and read-only Timeline tests remain green. No other product paths change.
- Verify **exact HEAD** GitHub Actions and mandatory `CHK-SECURITY-TESTS` unchanged floor >=90%; report any failure as a real merge blocker.

## Integration state

This is one small real R2 read-only interoperability/security correction, **not a new functional product capability** and not cause for repeated L2. Parent foundation #36 and all other stacks remain unmerged, so the task is not on `main`; no production claims. If only CI remains blocked, leave this child PR DRAFT, do not bypass requirements or create another layer of microtests. For a fresh agent, inspect the PR/branch exact SHA, successful/failing runs, and compare the three changed runtime/test files plus graph/manifest to frozen parent #36. The official L2 independent reviewer is for MATERIAL features, not every minor regression correction.

## Independent integration carrier (not a replacement for PR #133)

This manifest and exact selector source files were incorporated **byte-for-byte** from frozen independent draft PR #133 HEAD `868abc62da1e6a55ee1cfcb6c18db64474d6aa8e`, into a separate child of draft PR #132 HEAD `57c7cb78c8db14a58dfb8d19371100c2a229f0a1`. The three source/test files were verified byte-identical to the frozen PR #36 foundation on the #132 branch *before* this operation; this is a conflict-free functional convergence of CAP-INV-008 Class-0 Search Job provenance and CAP-INV-110/101 read-only Case Queue→Timeline navigation. Previous PRs and their frozen SHAs remain unchanged. The cross-integration child must independently prove both workflows' exact SHA CI, mandatory >=90% changed-security coverage and review/merge rules. The original #133 is still DRAFT; don't close, retarget or mark it merged until a valid reviewed integration makes that safe. This paragraph is provenance/handoff, **not** a claim of acceptance or that this child fixes the inherited #132 89.30% security coverage failure.
