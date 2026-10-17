# E10-INV-150A — Executed first-party Wazuh pivot search

## Material functional delivery

The existing E10-INV-140A Wazuh→Event Inspection→PivotDraft flow now optionally **executes** the normalized Rule ID or Agent ID pivot using CMDR's existing `event-search.NewOrchestrator`, `Create` and `Execute` lifecycle with a bounded first-party Case-owned Wazuh backend.

From the existing Wazuh Event Inspection CLI, specify all three `-pivot-wazuh-field=wazuh.rule.id` (or `wazuh.agent.id`), `-pivot-start`, `-pivot-end` **and** `-execute-wazuh-pivot`. Without the execute flag, the existing inspection and draft behavior stays unchanged.

The output includes a real first-party SearchJob with completed state, immutable randomized job/run IDs, tenant/environment/time/source/query provenance, a create/start/finish audit trail and **stable matched event references**, not source payloads. Query interval is half-open [start,end), with only certain-time Wazuh owner-authorized alerts matching the exact selected rule/agent. The `wazuh4-authorized-local-import` source is an explicitly **synthetic in-process import scope**, never a live Wazuh backend or externally trusted datasource. Max 200 matches per execution; overflow fails the complete search rather than silently truncating it.

## Trust and simple technical choice

- Reuse existing operator-selected read-only Wazuh import; no second parser or storage; reuse the same existing Timeline 200-row owner traversal.
- The Event Search Backend reauthorizes a **new complete Case traversal** when executing, and independently reopens each matching typed telemetry source. It refuses partial source grants, changed import totals, revoked bearer, tenant mismatch, invalid timestamps, overflow and cancellation. No result is emitted until successful backend/job completion.
- Use existing `eventsearch.Validate` four permission checks and its lifecycle. Demonstration search permissions are assigned **ONLY to the locally selected synthetic operator**, not extracted from Wazuh or caller-controlled HTTP, and MUST be replaced by real IAM at production integration.
- Search is **metadata equality** on `wazuh.rule.id` or `wazuh.agent.id` only. No raw message search, no direct arbitrary DSL, no cross-tenant fallback, no changes to owner snapshots, no custody signature, no live Wazuh service, no persistent index, no vendor SDK and no paid API.
- Reuse exact first-party Go modules by local replace, with a narrow recorded architecture dependency. The parent PR #53 and its L2 SHA remain unchanged.

## Tests and acceptance

GitHub Actions exact HEAD: `investigation-timeline-runtime` Linux/macOS/Windows and `engineering-kernel`; require security coverage floor 90% and independent one-shot risk-based L2. Real test fixtures read Wazuh feeds from disk with original alert ID collision across two feeds, execute a 2-result time-bounded rule job and a narrower agent job, verify completed SearchJob audit and provenance, exclude late event at time boundary, reject foreign principal/Case, revoked owner credential, inconsistent visible import, invalid/conflicting flags and cancelled jobs, and verify original owner/feed bytes unchanged.

**Limits and merge gates:** stacked on PR #53 (parent PR #50; original docs/engineering PR #2/#36), draft only until CI, independent L2 and upstream merge gates; not production IAM, live SIEM/EDR, governed case mutation or durable job persistence. The earlier security coverage failure on PR #53 is NOT bypassed by this functional enhancement.
