# E10-INV-101I-FILTERS — Native authorized Case Queue GET filtering and sorting

## Ownership / canonical obligation

- **Product:** `CAP-INV-101` Case Queue / Investigate / Cases and Evidence; `REQ-PROD-014`, `REQ-PROD-020`; `INV-CAS-001`. Class-0 accessible queue filter/search/sort behavior in `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` §§3, 8, 11, 12, 15, 21, 25. Tenant/permission contract and `pagination-filtering-sorting.md` remain canonical.
- **Engineering parent:** `E10-INV-101`, unit `E10-INV-101I-FILTERS`; code dependency `E10-INV-101H-HTTP` IMPLEMENTED with exact-head engineering-kernel success, NOT VERIFIED or accepted. The verified runtime foundation `E10-INV-101F` already provides `ProjectWithSavedViewFallback`.
- **Code owner:** `product-runtime/case-queue`; neither a Shared Saved View backend nor a Timeline demo. Product specs remain unmodified.
- **Git parent only:** frozen L2 PR #139 SHA `47a5eb698d85356881d13a506483e28f1c238638`; new child branch `feature/e10-inv-101i-native-queue-http-filters`. PR #138 is an **independent parallel child** and no business parent.

## Requirement → prior behavior → new capability → exclusions

- Existing Case Queue core **already** has first-party deterministic literal filter and sorting semantics with fixed stable tie-breakers. The native authenticated GET adapter in #139 previously exposed **only `q`**. This child extends that SAME server-authorized adapter with one optional literal each of `status`, `owner`, `freshness`, and `sort=title:asc` or any existing allowed field `updated_at|human_id|title|owner|status` with `asc|desc`. Query `q` remains optional; at most 5 total unique keys. No new search engine, rule evaluator, new entity or pagination algorithm.
- All candidates must pass the existing per-Case server authorizer, then any optional Incident and Finding are checked by the host-supplied independent typed-object authorizer BEFORE `ProjectWithSavedViewFallback` applies any filter/sort. A denied Case is removed completely, including from total counts, query matching and filter decisions.
- Filter literal valid UTF-8 <=128 bytes, no ASCII control, nonempty after trim. Search <=256 bytes, one sort field/direction from the projector's **existing** allowlist. Unknown keys (tenant/env/ACL/Saved View/permissions etc), duplicates, malformed sort, invalid UTF-8, control characters and oversized query strings are rejected with generic HTTP 400 BEFORE any host session or provider call.
- Existing no-store/no-referrer/Vary Authorization/nosniff, 5000 candidate, 4MiB serialized output and generic fail-closed response remain unchanged.
- **Out of scope:** production auth/IAM/provider, final Case Queue UI, pagination/cursor semantics (decisions remain pending), persisted Saved Views, Case assignment/status (OPEN-013), export, Task/Decision/Response Run, TLS, rate limits and durable audit. Successful Go handler tests are NOT a deployed production Case Queue.

## Behavioral acceptance evidence (exact-head CI pending until available)

- Authorized HTTP GET with Case title literal search plus status, owner, freshness and deterministic title sorting returns **only** the two permitted matching Cases, never the denied matching Case, its title, ID or denied count. Independently authorizes all candidates before filtering and returns correct native sort metadata.
- Alternate `sort=updated_at:desc` uses the existing deterministic Case Queue sort, no side effects.
- Unknown tenant/permission key, duplicate status/sort, missing filter value, control byte, oversized filter, malformed/disallowed sort or six parameters fail HTTP 400 **before** session resolver, provider and authorization handler execute.
- Existing #139 positive and negative owner security suite remains unchanged. All tests are Go standard-library `httptest` with no external dependency.
- **CI**: observe exact GitHub Actions after publication; `engineering-kernel` mandatory changed security-critical coverage floor >=90%, no waiver. `investigation-timeline-runtime` is not the Case Queue test workflow, so do not claim it ran.
- **L2**: ONE independent read-only review after published stable exact SHA for HTTP query/authorization trust boundary; no duplicate for #139. `L2_À_PLANIFIER` until that task is actually submitted; GitHub native approvals and upstream reconciliation remain separate.

## Recovery / integration

Work graph now records `E10-INV-101I-FILTERS` as IMPLEMENTED, not VERIFIED, because source changes exist but CI/L2/merge remain independent. Parent #139 remains frozen. PR #36 and specs #2 remain unmerged; no product acceptance/status promotion. New agent should start at `AGENTS.md`, `AI_START_HERE.md`, the graph, this manifest and the exact GitHub SHA/CI.
