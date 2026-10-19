# E10-INV-101R-SELECT — Reauthorized selection in the native Case Queue browser

**Canonical owner:** CAP-INV-101, Investigate / Case Queue (particularly §17), REQ-PROD-014/020; read-only Case fields remain Investigate-owned; Command Incident and Findings are independently source-authorized; Shared owns the FULL canonical Inspector UI. **Engineering parent:** existing E10-INV-101 epic, prerequisite Q-BROWSER IMPLEMENTED. **Code owner:** product-runtime/case-queue static browser consumer only; **Git parent:** immutable #148 HEAD `ab997eb3ab2cfa334cc2e2d53210120c69e53430`, NOT main or business parent.

## Material functional delta

The Q browser lists 25 currently authorized Case rows, allows literal query and cursor pagination, but the Case Queue's already implemented first-party `GET /case-queue?inspect_case=EXACT_ID` endpoint had no browser consumer. Each row now exposes one keyboard-focusable **Inspect** action. This action makes a **NEW authorized request**, passing exact original Case ID and CURRENT `q`, `limit` and `cursor`. It does not trust the old list row for access. The host server again authenticates its current session, authorizes each Case/Incident/Finding individually, applies current Shared Saved View and page snapshot, and returns only a currently visible authorized `Row`. Rejection, stale cursor, source/Case revocation or a different returned Case ID fails closed.

The browser shows only inert textContent for selected Case title, exact ID, status, owner, freshness/partial/stale, updated timestamp and next action. No raw source, Incident or Finding bodies, arbitrary HTML or bearer are displayed, and no client permission or scope is introduced. Closing via button or Escape restores keyboard focus to the trigger when it still exists, keeping the original Queue page/search. Search/page changes and pagehide abort and discard any pending selection and clear displayed sensitive details.

This is a deliberately small **Case selection context** to connect a canonical CAP-INV-101 transition. It is **NOT a locally cloned Shared Inspector component** (no Shared sections/actions, keyboard system, resizing or full-screen), nor the final Case Workspace host UI. Shared owns the eventual reusable Inspector presentation; future integration must consume its component rather than copy it here.

## Risk, tests and limits

- Independent source-authorized exact ID request is mandatory even after an allowed list response; response ID is checked byte-for-byte without trimming significant leading/trailing spaces.
- No stale or denied previous Case details survive another selection, rejected request, page change or navigation; no partial unauthorized response is rendered. Every DOM user-controlled value is assigned with `textContent`.
- Node-built-in tests exercise true browser script in a DOM harness: successful case with significant spaces, inert malicious title, current host credentials, focus restoration, denied second selection clearing first and mismatched Case response denied, pagination purging selection. Existing Go browser host security tests verify current server ownership and no-static-data rules. No real production browser screenshot or IAM claim.
- No new modules, dependencies, server APIs, payload caches, browser URL authorization, automatic fetch/poll, Case writes or paid APIs. Read authorization stays completely within existing native Go handler.
- The unchanged engineering-kernel changed-security-critical coverage floor remains **90% mandatory**. Product acceptance, production readiness and merging remain blocked by draft upstream PRs, native approval requirements and outstanding full Shared Inspector/Workspace obligations.

**Next:** exact-HEAD GitHub Actions + one independent read-only L2 of this material browser/authorization integration. Keep all parent L2 SHAs unchanged; treat this new task as IMPLEMENTED until actual CI/approval evidence, not VERIFIED.
