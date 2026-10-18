# E10-INV-280A — Consolidated authorized saved Case boundary

## Why this integration is material

After PRs #62–#66, three separately delivered investigator journeys existed: (1) resume an imported Case and execute authorized Event Search; (2) verify one saved Case and report its accessible imported inventory entirely offline; and (3) compare two independently saved Case snapshots with the first-party chronology engine.

All three independently implemented the **same** trust boundary sequence (read bounded snapshot and operator-held SHA256 if supplied → validate owner graph → scope to one authorized principal, tenant, environment, Case, Timeline and sources → recognize canonical imported IDs). Duplication risks future divergence of tenant/source validation, and unnecessarily increases exceptional Go paths.

This integration extracts exactly **one narrow first-party function** `loadAuthorizedImportedCase(ctx, principal, path, caseID, expectedDigest)`. It returns the minimal scoped owner graph and canonical imported source IDs only after the existing loader/owner checks and context cancellation. Nothing else grants or changes permissions. Existing `serveDemoReopenedNormalizedCase`, `verifySavedNormalizedCase` and `compareVerifiedSavedCases` call that boundary rather than independently assembling it. The comparison invokes it TWICE, once per saved archive, and still requires same tenant AND environment before obtaining any comparative timelines. The ordinary owned-case file saver still projects the Case *before* filesystem serialization; no raw full-graph persistence returns.

This is a **code-level integration / reduction of duplicate security logic**, not a new runtime feature, IAM framework or compliance gate. Preserve exact HTTP owner ACL, source reopen, live Event Search Job behavior, operator-only snapshots, strict no-overwrite atomic write, separate SHA256 verification and aggregated Case comparison. The 90% mandatory changed security-critical coverage gate and negative tests are unchanged.

## Acceptance

The unchanged behavioral tests (including all tests added in PRs #62–#66) exercise:

- one authorized Case save and encrypted? **No**, explicitly private but not encrypted local source;
- actual restart with new bearer and first-party source authorization;
- bytes modified after save denied when operator supplies independently held SHA-256;
- offline authorization/tenant-aware inventory with source-denied entries removed;
- two independently verified Cases compared with no owner graph merging, source/tenant scope migration or partial report when the second Case fails;
- Linux Go+Node and macOS/Windows Go, full integration and unchanged `engineering-kernel` gate.

One independent risk-based read-only L2 on the frozen merged delta, not one per edit, is warranted because **all three sensitive paths share a changed trust boundary**.

**Stack:** base SHA frozen PR #66 `694b57401f2b5f61d31bf77e39656f8ed27510a9`, branch `feature/e10-inv-280a-unified-authorized-case-loader`. Other L2-frozen ancestor SHA untouched. `main` is still documentation/specs; #36 is still a very large draft engineering base (763 commits/871 files). No production IAM, real-time streaming, signed forensic custody, immutable database or approved merges are claimed. If required CI is red, keep DRAFT and document exact failure without reducing coverage or proliferating near-identical tests.