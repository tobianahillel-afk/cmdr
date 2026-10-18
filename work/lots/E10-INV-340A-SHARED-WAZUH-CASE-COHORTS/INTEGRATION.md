# E10-INV-340A — Shared Wazuh rule/agent temporal leads across two saved Cases

## Genuine investigation integration

CMDR already has a first-party read-only chronological Case comparison, which works on two **separate, individually authorized and independently SHA-256-verified saved Case snapshots**. The integrated draft PR #73 separately restores source-owned Wazuh 4.x `rule_id`, `agent_id`, optional level and feed metadata after deleting the original JSONL alert files, and supports a local Case Event Search Job and same-Case 60-second bursts.

This increment connects these TWO existing capabilities to answer a concrete analyst question:

> Did two independent authorized Cases observe the same Wazuh rule **and** the same agent during the **same definite five-minute UTC window**, and how many observations occurred in each Case?

Reuses the existing CLI command and first-party `buildCaseComparison`:

```sh
go run ./cmd/case-timeline-demo -principal analyst-a \
  -snapshot /absolute/authorized-wazuh-case-A.json \
  -compare-saved-case CASE-A -expect-snapshot-sha256 A64HEX \
  -compare-saved-with-snapshot /absolute/authorized-wazuh-case-B.json \
  -compare-saved-with-case CASE-B -expect-other-snapshot-sha256 B64HEX
```

After **both** Case files have been individually verified byte-for-byte, scoped to the same actual principal and tenant/environment, authorized at Case/Timeline/typed-source boundaries, and traversed by the existing first-party triage, compare bounded canonical source metadata only from the **selected authorized imported events**. Group by exact pair `(rule_id, agent_id)` and five-minute UTC bucket. Output at most **10** time-ordered cohort leads, each with UTC window, rule ID, agent ID and event counts for Case A and B. It adds `shared_authorized_wazuh_five_minute_windows` to the existing comparison JSON; all existing aggregate fields and their semantics remain unchanged.

No raw Wazuh alerts, descriptions, IP/hostname, bearer, original Wazuh alert identifiers, or combined data-store/source-owner objects in the new lead. The existing older `top_sources` aggregate in Case comparison remains unchanged and may include already-authorized synthetic source references. Operator-import slot labels and Wazuh metadata are **descriptive**, never treated as ACLs. Entries with uncertain time or malformed Wazuh source time cannot contribute; generic JSONL Case owners with no Wazuh metadata produce NO invented cohorts. An agent identifier is meaningful only to the same tenant/environment and might still refer to different real-world agents: **this is a retrospective co-occurrence lead, not an incident verdict, provenance guarantee or causal attribution**.

### Technical choice

No additional full-text backend, search orchestrator, materialized index, permission validator framework, SQL engine or paid API. Uses first-party saved authorized owner projection, canonical imported identity, source-owned immutable bounded Wazuh metadata, and existing independent Case chronology comparison. Counting is O(authorized entries) on each Case, bounded by existing 5000 Case-entry owner limit; max 10 returned leads.

### Acceptance, risks and required gates

- Real Wazuh alert.json files parsed using existing Wazuh 4.x normalizer and persisted as two separate private owner Case files. Each has an independently retained SHA-256; original alert files removed. Two identical `rule_id=5710, agent_id=001` alerts in Case A and one in Case B within `2026-10-17T10:00:00Z` bucket form exactly one lead with counts `2/1`. Different agent, unrelated rule or next five-minute UTC interval is excluded.
- Confirm no source text/credentials in aggregated lead. Generic JSONL-only Cases return zero Wazuh-specific leads without failing existing Case comparison.
- Existing typed source ACL, tenant/environment separation, digest tamper detection and no-partial-output behavior must hold unchanged for two independently authorized owners.
- Exact-head `investigation-timeline-runtime` Linux Go+Node / Windows / macOS and mandatory `engineering-kernel` security-critical 90% floor. Prior latest #73 security CI was **87.96% < 90%** and is NOT mergeable; this PR must not waive/change the gate or force merges. Only one independent risk-based READ-ONLY L2 on stable exact SHA, not duplicate earlier L2 of #71/#72/#73. No production IAM, real ingest, durable DB or signed forensic evidence claim.

This is a **single coherent integration increment**, not independent microtest PRs. Parent PR #73 SHA `b375c78eb5645a2d2e686e3b970c21d2e5b10b6a` is L2-frozen and must not be altered; PR #36 and the long ancestor chain remain draft.
