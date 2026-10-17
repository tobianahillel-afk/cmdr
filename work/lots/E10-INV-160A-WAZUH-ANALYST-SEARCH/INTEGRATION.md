# E10-INV-160A — Executed Wazuh Search Jobs in the analyst viewer

## Functional operator journey

Extend the existing first-party **local synthetic** Wazuh lab UI with a visible agent-or-rule search form, bounded local start/end times, and an explicit **Execute Search Job** action.

1. Read one-to-four operator-local Wazuh 4.x JSONL feeds using the existing trusted snapshot importer and launch the existing loopback analyst page with `-serve-import`.
2. Select a permitted Case and enter the ephemeral bearer. Choose `rule_id` or `agent_id`, a numeric ID, and the half-open `[from,to)` time window.
3. GET `/demo-wazuh-search?case_id=...&field=rule_id&value=5710&from=...&to=...` with the existing ephemeral bearer. The route exists **only in Wazuh import mode** and is not installed on other Case contexts.
4. A single existing first-party `Event Search` backend/orchestrator performs a new owner-authorized Case/Timeline traversal, independently re-opens every returned telemetry-event source, and returns the true completed Search Job containing stable references, scope/provenance and audit events. Nothing is released if any authorization fails. Results limit **200**, no silent truncation.
5. The analyst page shows only the completed Search Job's event references with buttons. Clicking opens the existing independent `openSource` Case/Entry/typed-source route with a **new authorization**, never an unverified preview. Changing Case, token, dates or parameters clears the results and cancels in-flight requests. UI uses `textContent`, no unsafe HTML or token in URL.

This is the same backend used by PR #54's existing CLI `-execute-wazuh-pivot` mode; shared `executeDemoWazuhMetadataSearch`, **no duplicate query engine**. The interactive UI initiates a direct metadata search rather than claiming that every query originated from Event Inspection.

## Explicit exclusions and security

- Offline synthetic principal permissions only; browser-provided rule/agent IDs are filters, **not tenant/Case/source permissions**.
- No live Wazuh subscription, multitenant production IAM, database/index, durable Search Jobs, chain-of-custody guarantees, signed evidence or real endpoint protection. No paid API, third-party library, external process, new listener, or persistent credential.
- Run HTTP server on existing ephemeral IPv4 loopback; enforce original constant-time bearer check and no-store response. Search provider uses **server-selected** tenant/environment and Case ID. Reject missing/wrong bearer, mismatched Case, duplicate/unknown parameters, invalid/missing field/value/time, unsafe HTTP method, cross-tenant sources, partial results and changing ACL.
- Tests import two real Wazuh feeds with a colliding original ID, execute two-result rule query and narrower agent query, ensure out-of-window and non-matching alerts excluded, check no-match completed state, verify same source reopens via existing owner, require no sensitive previews or bearer in Search Job receipts, and refuse invalid inputs without data.

## Integration and review

Stack on immutable PR #54 head `6166f816b22b396d888aa8ee00047d6371730920` (PR #53/#50/#36/#2 below). Do **not** modify reviewed parent branches or `main`.

Acceptance requires exact-head GitHub Actions `investigation-timeline-runtime` (Linux/macOS/Windows and Node browser tests), `engineering-kernel`, and one risk-based independent read-only L2. Parent PR #54 already fails the mandatory security critical coverage **87.50% < 90.00%**; this gate remains binding regardless of functional CI. Do not merge while red or required L2 outstanding.
