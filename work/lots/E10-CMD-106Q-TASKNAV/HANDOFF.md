# E10-CMD-106Q-TASKNAV — Incident linked Task → canonical Work Queue

## Product / engineering / Git owners

- Canonical Command product: CAP-CMD-106 Incident Coordination, consuming the existing Command-owned Task row and CAP-CMD-101 Unified Work Queue; Task object is owned by Command, linked Incident is no Task.read grant, Shared relation owner stays independent. The separate `screens/tasks.md` is **deprecated**; the navigation uses the one `GET /command/work-queue?view=tasks` route.
- Engineering parent: existing `E10-CMD-106` epic; task `E10-CMD-106Q-TASKNAV` depends `E10-CMD-106P-TASKLINKS` recorded IMPLEMENTED (not VERIFIED or ACCEPTED). No broader READY work unit existed on the parent; child is a selected bounded missing canonical navigation behavior.
- Code: original `product-runtime/work-queue` first-party Go Project, HTML and HTTP and existing `product-runtime/command-investigate-host/cmd/local-browser` synthetic loopback host. No new package, endpoint, owner, provider or library.
- Git parent: PR #215 immutable HEAD `a97ab1be523eb164ae1a263615e21ef9d406e3bc`, base branch `feature/e10-cmd-106p-authorized-incident-tasks`. Git parent is not engineering or product parent. Parent and all prior L2 SHA untouched.

## Material usable behavior

Incident Detail already displays **independently authorized** linked Task ID/status/owner using relation owner and original `Owner.List` + per Task `Owner.CanRead`; the new server-rendered `Ouvrir la Task` link navigates to **the existing** Command Work Queue `view=tasks&selected=Task:ID&restore=1` without embedding principal/tenant/env or other authority. Existing `restore=1` remains a fixed opt-in browser-only selection hint, now type-bound to either `view=incidents&selected=Incident:...` or `view=tasks&selected=Task:...`. The existing `Project` scans at most 5000 owner records, independently filters through `CanRead` for the actual session principal/scope, then selects the authorized page. User cannot supply offsets, digest, arbitrary view/filter, permission or source ID. The token does **not** grant Task.read. No new Task Detail product page, no Task editing, assignment, Case read, mutation or storage. Browser back naturally restores original Incident Detail; no new return hint authority introduced.

## Tests and risk

New direct Go handler regressions exercise 123 Tasks, find Task #100 only with a fresh owner-authorized `restore=1` HTML request, correct pagination and generic 404 after revocation; denial of hidden/foreign Task data, forged kind/view, untrusted query, duplicate, invalid cursor or JSON restore fails before owner. Same source relationship owner provides a link only while Task CanRead. A local true `httptest.NewServer` traversal with two ephemeral Basic synthetic principals checks Incident Detail → Work Queue Task and Command-only denied Task, using the preexisting host and no extra listener. Original Incident restore untouched; tests for it must remain passing. Security gofmt, vet, Go race, changed-critical security coverage >=90, govulncheck, gosec and engineering-kernel CI **MUST be checked on final exact HEAD**. L2 separate context for material navigation change only if CI and native review availability support, never mutate audited #215.

## Limits

Both host principals and data owners are synthetic only; no real IAM, production Task store, audit sink, Task Detail, generated Action Request, write operation, TLS internet deployment or full CMD-IWQ-006 acceptance. Source Task rows have no stable version here. A returned Work Queue may be newer than the Incident snapshot; both owners reauthorize independently. Parent PRs #2/#36/#215 and this descendant DRAFT are NOT merged or natively approved. Passing CI != GitHub APPROVED/production.
