# E10-INV-110B — analyst kind/time read-only query

A usable in-process Case-to-Timeline HTTP flow now supports optional `kind`, `from` and `to` query selectors after complete server-side Case, Timeline Entry and typed source checks. A half-open [from,to) UTC-equivalent interval allows deterministic time triage. Unknown or clock-uncertain entries remain visible and explicitly uncertain rather than silently excluded. This is not evidence of production session/auth providers, a persistent repository, complete C-NAV or UI.

Reuse-first selection: existing Go net/http, net/url, time and immutable Result; no extra database, index, third-party framework, ACL engine or new owner dependency. O(n) over a bounded 5000-entry already authorized projection. Unknown selectors, duplicates, invalid times and forged permission/tenant parameters fail before any authority I/O. Read-only status and source/provenance remain inherited from Project; supersession of a permitted predecessor may point outside a narrowed time window.

Remaining delivery blockers: real server authentication, persistent Case and Shared Timeline owner catalog, source-owner authorization providers, host TLS/routing/audit/rate limiting, full navigation UI and source-destination routing, and independent L2. Product Spec, OPEN-013 and graph statuses unchanged. Exact-SHA CI must be checked before claiming verified operation.
