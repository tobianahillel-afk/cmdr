# E10-CMD-105B-EXPLAIN — source-owned calendar and paused SLA explanation

## Canonical authority and four ownership identities

- **Product:** complete canonical `CAP-CMD-105` (Command SLA Tracking) §§1–27 explicitly requires explaining calendar, pause and policy source (§§3/11/15/20); `CAP-CMD-101` consumes selected Work Queue state. Validated ADR-0008 retains external contractual/calendar authority in source; `OPEN-013` makes Class-2 pause/resume/Task mutation unavailable. This unit **does not implement them**.
- **Engineering:** existing epic `E10-CMD-105`, bounded task `E10-CMD-105B-EXPLAIN`, dependency `E10-CMD-105A-READONLY` IMPLEMENTED on source branch. Initial graph had **no READY nodes**, but current canonical subset of Class-0 disclosure and exact-HEAD success of previous prerequisite permit a separately bounded next Class-0 increment; no false VERIFIED conversion.
- **Code:** first-party Go `product-runtime/work-queue`, existing `SLAOwner.ReadAuthorizedSLA` host adapter and existing selected-row `GET /command/work-queue` API/native HTML. No new route, IAM provider, parser, SLA calculation engine, calendar policy or third-party dependency.
- **Git:** new child branch from frozen parent PR #174 HEAD `b1238eb29e3e51dade0166e93de9819f5e3ba334`; all previous L2 source SHA untouched. Git parent #174 is **not** the product/engineering parent.

## Requirement → runtime before → demonstrated addition → remaining

| Requirement | Previously at PR #174 | New bounded behavior | Still absent |
|---|---|---|---|
| Selected SLA state + due/time and policy source | Host `AuthorizedSLA` selected-only JSON/HTML | Unchanged | Production policy/calendar source and live deployment |
| Explain source calendar | No calendar source/explanation exposed | Optional host-provided `calendar_source`, shown only after current selected-row and SLA source authorization; unknown if absent | Actual versioned calendar rules and due calculations |
| Explain authorized active pause | `paused` state only, no reason/time | Optional `paused_at` and `pause_reason`; source-owned UTC-offset RFC3339, only if state paused, safe HTML escaping, missing fields explicitly unknown | Pause history/audit trail, policy authorization, actual pause/restart mutations blocked by OPEN-013 |
| Denied/foreign/invalid source | Generic response on invalid SLA | Incoherent/untrusted paused metadata (healthy with pause, reason with wrong state, timestamp later than observed, controls/oversized strings) returns existing generic HTTP 503 with no selected item | Real host integrations, production IAM/permissions |

**Trust boundary:** The selected object owner and the independent SLA policy/engagement host adapter are both responsible for source-side permission. The Go validator rejects malformed source claims, but is **not** an IAM engine; source-only fields cannot be requested from client query, and 0 additional owner calls are made.

**Performance/security:** The existing one selected-row source lookup is reused, so no per-row policy lookups in a 5000-item queue, no new owner traversal and no saved snapshot. `html/template` handles source text escaping. Source fields are optional for backward compatibility; a pause with no source metadata is rendered **unknown**, not invented. CalendarSource is a descriptive source reference, NOT a policy engine or calendar validation. A source returning false claims is outside this synthetic host contract.

## Tests and merge restrictions

Behavioral tests cover real selected JSON and HTML, zero additional lookups for unselected/denied/foreign/off-page cases, source pause/time preservation, missing pause explanations, XSS-shaped source strings escaped, invalid calendar/reason or contradictions denied without source/row leakage. Existing six-view, pagination, session/ACL and first-party SLA tests remain required.

This child is DRAFT and **not accepted, merged or production-ready** until exact-HEAD `command-work-queue-runtime` on Linux/macOS/Windows, mandatory `engineering-kernel` incl changed-critical security coverage floor **90% unchanged**, and native approval/parent dependencies all pass. One independent **read-only** L2 on material authorized source-metadata boundary after stable GitHub evidence. Do not duplicate previously scheduled PR #174 L2; it is a distinct frozen SHA review.
