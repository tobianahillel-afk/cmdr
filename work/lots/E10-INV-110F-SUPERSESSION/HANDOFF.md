# E10-INV-110F-SUPERSESSION — Authorized read-only supersession visibility

## Why this is product behavior, not a verification micro-cycle

The CAP-INV-110 canonical timeline requires corrections and superseded older versions to remain visible. The existing Go projection preserves **successor -> visible predecessor** via `supersedes`, but consumers inspecting an older row cannot distinguish it from a still-current row without scanning and joining the full timeline themselves. That fails especially on bounded 200-row HTTP pages. This increment produces a directly inspectable `superseded: true` on that older row, computed once over the **full authorized projection before page/filter selection**, keeping the old entry and source intact. It does **not** create a correction, decide which of concurrent successors is canonical, or show a protected successor ID. Absence of the flag means *no authorized visible successor in this read*, not proof that the source has none.

## Four owner identities

- **Product:** CAP-INV-110 Investigation Timeline read, Case belongs Investigate, Timeline Entry and source lineage remain Shared/typed source owners. No new ownership or mutation of Timeline Entry.
- **Engineering:** existing `E10-INV-110` epic, independent `E10-INV-110F-SUPERSESSION` child of VERIFIED E10-INV-110E. Other READY `E10-INV-120A-SNAPSHOT` is a **separately recorded broad, stale and already implemented-in-parts path**; do not rewrite its history, auto-verify it or conflate this work with it.
- **Code:** existing Go first-party `product-runtime/investigation-timeline/projection.go` and bounded regression test. Uses existing in-process case/catalog/source adapters, existing `/case-timeline` JSON, no extra provider, route, dependency or persisted data.
- **Git:** branch `feature/e10-inv-110f-visible-supersession` from PR #36 frozen head `ea1770b48a4bc1b2e850a7ffab15d5839f9e9e49`, independent of DRAFT PR #187 Wazuh HEAD `53544c2a65f66324d004d48427b07493f7ac2d9e`. Zero changes in PR #187 paths or SHA; the relation is git ancestry, not product parenthood.

## Access/security invariant and tests

`Project` already rejects invalid source versions and hides `Supersedes` targets missing from the fully owner-authorized entry cohort. Using the same visible ID->version mapping, the projection now also identifies the visible predecessor row and marks it `superseded: true` **only after** the successor passes all existing validation. It changes neither source permissions nor Case selection. The field is omitted when false (backward compatible, non-definitive). Tests cover authorized two-version and three-version chains, hidden/unknown/missing/foreign successor not changing visible history, HTTP bearer and tenant Case denial, and an old version still marked non-current when its authorized successor lies on another HTTP page or outside a non-security UI filter. There is no inferred causal link, external persistence, mutation Class-2 (OPEN-013), final user-interface badge, new trusted root or production IAM.

## Release evidence and handoff

This document is committed as the audit handoff before CI. The PR must record exact candidate HEAD, changed files, mandatory engineering checks, `go test -race` and `go vet` as GitHub-visible evidence. No independent L2 review and no production acceptance have occurred at authoring; report their actual later statuses separately. A native independent GitHub APPROVED review, merged upstream `#36/#2`, real IAM/Shared Timeline source, correction permission/audit workflow and final analyst UI remain outside this bounded increment.
