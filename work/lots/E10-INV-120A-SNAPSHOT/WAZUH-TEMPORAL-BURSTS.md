# E10-INV-120A — authorized Wazuh agent/rule temporal bursts

## Product outcome

The existing offline Wazuh 4.x triage (single-feed or 2–4 feeds) now identifies **repeated detections for the same Wazuh agent ID and rule ID within an inclusive 60-second interval**. An analyst receives up to ten ranked summaries with exact time-certain UTC endpoints, number of observed alerts and number of distinct operator-selected feeds. The authenticated loopback viewer displays these leads and offers a button opening the existing **server-side authorized, paginated agent/rule drilldown**. The drilldown deliberately shows the complete cohort, not an inaccurately claimed time-sliced subset.

## Scope and implementation

- Based on PR #50 exact SHA `e7c980949fd8b5b15bf905c614c10ff62a8f2a75` on a separate stacked development branch, preserving the frozen PR #50 L2 review target.
- Reuses Wazuh 4.x validated normalizer, operator-selected feed identifiers, synthetic Case owner ACL, authenticated Case/Timeline/typed-source traversal and existing browser route.
- Only **authorized imported Timeline entries with known timestamps** are eligible. Alerts with uncertain timestamps, absent agent/rule identifiers or missing source authorizations do not create leads. Any denied imported source fails the complete summary closed before writing.
- A burst is at least **three events for the same agent/rule within at most 60 seconds inclusive**. The strongest window per pair is selected; ranking is by alert count, distinct feeds, earliest UTC and pair identity. Limits: current owner 5,000-entry bound; top ten burst leads.
- Sliding two-pointer algorithm avoids quadratic scans and introduces no third-party dependency, database, listener, live ingest or write API. Source content and original alert IDs are never included in burst summaries.
- Temporal concentration is **not attribution, intrusion verification, causal correlation or a detection signature**. All ACLs and identities are lab-only synthetic inputs; `OPEN-013` remains open.

## Acceptance and remaining obligations

- Unit tests: 60-second inclusive edge and >60-second negative, single-/multi-feed counts, deterministic ranking/top-ten.
- End-to-end: two actual operator-file feeds imported to one permitted Case; after owner authorization, one three-alert burst is reported, a later unrelated alert excluded, no raw Wazuh source content exposed; a foreign principal obtains no report.
- Existing authorized drilldown, security isolation, CLI triage, Node UI tests, Go race/vet and platform CI must continue to pass at **exact PR head**.
- Documentation does not change the canonical Product Spec or mark production readiness. Independent risk-based L2 and parent PR #50 resolution remain prerequisites for merging into the engineering branch.
