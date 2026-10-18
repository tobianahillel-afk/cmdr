# E10-INV-330A — Durable Wazuh Case search → temporal burst leads → source drilldown

## Why this materially integrates existing capabilities

Two genuine first-party components already exist:
1. The pre-existing Wazuh Case Event Search route provides completed source-authorized Wazuh hits and **60-second agent/rule burst leads**, using the proven bounded `deriveWazuhSearchBursts` sliding-window function.
2. PR #71 saves bounded Wazuh rule/agent/level/feed source metadata in the **authorized private Case owner** and restores first-party rule/agent text Search Jobs after the original alert files are deleted.

Until this lot, after restarting from a saved Case, an analyst could pivot by Wazuh rule/agent but **lost the event chronology and burst leads**. This integration reuses the **same exact sliding-window burst algorithm** with the **same existing completed Event Search Job**, not a second rule engine, general correlator, or rescanned raw source.

## End-to-end behavior

- Operator explicitly saves a bounded 1–4-feed Wazuh alerts.json dataset through the existing `-save-wazuh-case-snapshot`, retaining only the selected principal's Case and canonical typed sources. An operator can remove the raw JSONL files and restart with `-resume-import-search-case` (optionally with independently held SHA-256).
- The authenticated local viewer executes a first-party, bounded `POST /demo-local-search`, selecting an existing `rule_id`/`agent_id` pivot or literal owner-rendered preview. The backend traverses authorized Case Timeline pages and **reauthorizes every matching typed source BEFORE adding its reference to a Search Job**. The result must be completed; errors and canceled/partial jobs never return source metadata.
- **Only after successful Job completion** the server joins its **actual result references** with the same canonical saved owner Entry time (certain UTC only) and the existing bounded, owner-attached Wazuh metadata. The join rejects duplicate or foreign refs, wrong source type, missing canonical Entry, uncertain/invalid time or missing Wazuh metadata. It does not re-read raw files or trust client-supplied rule/agent metadata.
- The existing `deriveWazuhSearchBursts(hits)` yields up to 10 deterministic leads, each at least **3 matching events for the same agent_id+rule_id within an inclusive 60-second window**, with distinct authorized feed count, actual UTC start/end and member EntryIDs. No new correlation algorithm or score.
- Existing normalized SearchJob JSON receipt remains compatible. When **all** imported sources have owner-valid Wazuh metadata, add only `search_hits` (IDs, UTC, bounded rule/agent/severity/feed) and `search_bursts` (leads and member EntryIDs); never preview, raw log, bearer, hostname/IP or local file paths. For ordinary generic JSONL or incomplete Wazuh metadata, no fabricated Wazuh lead fields. Failing to enrich a fully Wazuh-backed completed job fails entire response, not a misleading partial lead.
- Browser renders a short lead description with a per-event button. It checks that lead members belong to the completed Job's returned references, that a lead has 3–200 distinct members, rule/agent IDs are bounded decimal strings, and times are a valid 60-second interval. Every click uses existing `openSource`, which performs **fresh Case/source authorization**. Changes in Case, token or filters erase stale leads.

## Exact behavior tests

Real disk Wazuh sources A+B include colliding original alert IDs, **3 same-rule/same-agent hits in 40 seconds across two source feeds**, a fourth same-rule/agent hit too late, and unrelated agent/rule. Save to private scoped Case, delete original JSONL, reopen with independently checked fingerprint, execute actual authorized `rule_id=5710` first-party SearchJob; expect 4 selected refs, 4 normalized hits and **one 3-member burst across two feeds** with exact UTC bounds. Every burst ref must belong to the completed SearchJob and be independently source-openable. No source content/bearer disclosed, denied bearer yields no refs. Browser tests show valid lead and reject a crafted lead with a foreign or repeated member, with zero partial results.

Security limits: **synthetic offline IAM only**; a burst is a temporal triage lead, **not proof of compromise**, causal attribution, or tamperproof signed evidence. Underlying Event Search's 500 candidate / 200 hit fail-closed limits and owner ACL must remain unchanged. No paid API, external dependency, database, live Wazuh ingestion or new governance layer.

## Integration and non-merge

Stacked **draft child of frozen PR #71** exact SHA `bdb554f6662c38c29f05fd0a8abb18ef30c5d195` (not main); all parent L2-frozen SHAs untouched. Parent PR #70 had **both** exact SHA GitHub workflows SUCCESS; PR #71 had functional CI success but mandatory changed security-critical coverage **88.42% / 90%** FAILED. Exact new SHA Linux Go+Node / macOS / Windows and the unchanged engineering-kernel security floor must be checked. Keep DRAFT/unmerged while ANY CI, L2, parent protection or approval fails. ONE independent read-only, risk-targeted L2 needed for this material correlation/security boundary; no reviewer writes, no duplication of parent reviews.
