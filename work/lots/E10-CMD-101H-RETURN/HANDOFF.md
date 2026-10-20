# E10-CMD-101H-RETURN — Exact previous-page restoration in Command Work Queue

## Truth and ownership

- **Canonical:** `CAP-CMD-101` Unified Work Queue §3 (single route and preserve view/filter/sort/selection/scroll), §§11, 16–17, 21, 25; Command's `Incident`/`Task` objects; Shared owns Saved Views and final Inspector, which this lot does **not** implement.
- **Engineering parent:** existing epic `E10-CMD-101`; task `E10-CMD-101H-RETURN` depends on `E10-CMD-101G-PAGE` (recorded IMPLEMENTED). **Git parent:** frozen PR #165 at `fbdfefe9a19d1fc6693e7e53b2649c59a2c18284` is only the technical base, not the business parent.
- **Code owner:** existing Go `product-runtime/work-queue`; no other boundary, dependencies or authorization model.

## Required → present → remaining

**Before:** page 1–N forward cursor and a return-to-FIRST link, with freshly reauthorized 50 rows but no direct way back from a deep page to its immediately preceding page. Going back from page 20 required restarting at page 1, losing the operator's position. This is material to the canonical exact-return navigation contract, not a new Queue feature family.

**New:** the existing `authorizedQueuePage` issues an optional `previous_cursor` for pages starting at offset ≥100 using the **same SHA-256 digest** of the freshly owner-authorized ordered row projection and same first-party cursor encoding. Page 2 returns to the original no-cursor route. Native same-route browser links now expose `Page précédente` and `Première page`, preserving `view`, exact `status` and `sort`. Existing `next_cursor`, first-page JSON shape, 50-row window, 5000-item candidate bound, per-page selection and six system views remain intact. No extra Owner/List authorization call and no second hash. The original cursor is NOT a capability token, signature, transactional snapshot or a final persistent schema.

**Security:** every next and previous page rebuilds current host principal/tenant/environment with one fresh `SessionResolver`, one bounded `Owner.List`, per in-scope object `Owner.CanRead` and the same authorization-derived revision check. A revoked, stale, denied or cross-context cursor yields generic 404 with **no partial row output**. No client-chosen principal/tenant or untrusted ID is an authorization source.

**Behavior proof:** a genuine Work Queue test goes through 123 authorized Items as 50→50→23→50→50, with denied same-tenant and foreign-tenant Items; exact expected IDs and filters/sort, every request authorization, no cross-tenant owner calls and denial of old backward token after source revocation. Native HTML previous link preserves the filter/sort/view and prevents denied row disclosure.

**Not delivered:** real host IAM, persistent queue, Shared Saved View state, final Inspector, Incident Detail, linked Case/Decision/Run navigation, transactional pagination and full CAP-CMD-101 acceptance. This task is only an incremental native read-only navigation integration.

## Git / CI and merge

- Parent PR #165 exact SHA `fbdfefe9a19d1fc6693e7e53b2649c59a2c18284` has two SUCCESS workflows (command-work-queue-runtime 39374253914, engineering-kernel 39374253912); that is NOT automatically proof for this child.
- Require exact-head Ubuntu/macOS/Windows `command-work-queue-runtime` and `engineering-kernel` including unchanged CHK-SECURITY-TESTS changed-critical 90% floor and Go SAST. Any failure remains a merge blocker, no threshold changes.
- This is a new material navigation/authorization integration and may merit one independent **read-only** L2 on a stable exact SHA, without duplicating parent's L2. Do not assume ChatGPT L2 is native GitHub approval.
- All ancestors including giant #36 (763 commits/871 files), #2 remain draft; no direct main mutation, no forced merge or production claim. Preserve exact SHA reviewed by L2, and do not create another micro-PR to remedy trivial test text.
