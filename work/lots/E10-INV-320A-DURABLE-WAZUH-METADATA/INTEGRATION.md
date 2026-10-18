# E10-INV-320A — Retain typed Wazuh rule/agent pivots in a saved authorized Case

## Functional gap closed

The original Wazuh 4.x Case workbench provides bounded **rule.id / agent.id** pivots and a rule.level indicator only while the operator's alert files are still loaded in memory. E10-INV-300A persisted a *sanitized* authorized Case and made its rendered alert text searchable after restart, but it could not recover rule/agent numeric pivots after original input files were deleted.

This increment makes the existing first-party Wazuh normalizer's **already validated bounded numeric metadata** durable as an OPTIONAL property of the corresponding canonical typed `telemetry-event` source. It stores strictly:
- `rule_id` (0–16 decimal digits, empty means unknown);
- `agent_id` (0–16 decimal digits, empty means unknown);
- `rule_level` when known (integer 0..16);
- operator-assigned `feed` ordinal (0..4).

No original alert JSON, original alert ID, `full_log`, host name, agent name, IP address, path, username, raw payload or bearer is added to this metadata. The existing source's bounded authorized *preview* retains the alert summary exactly as before. Generic normalized JSONL sources without a `Wazuh` block remain untouched; older saved snapshots load without it.

The existing Wazuh 1–4 file loader parses **the same bytes** into the already implemented `wazuhRuleLevels`/feed-namespace mapping and Case owner. Before private snapshot publication, it associates those metadata with the exact canonical source IDs: the number of linked sources MUST equal the validated number of Wazuh alerts. It then invokes the SAME selected-Case/principal/Entry/typed-source ACL projection and private 0600 atomic no-overwrite snapshot writer. The parent #70 redaction of foreign tenants/denied sources stays intact, including for metadata.

## Actual analyst journey

1. Operator saves a 1–4 feed Wazuh Case with `-save-wazuh-case-snapshot`, receives the existing SHA-256 digest, and may delete the original alerts.json files.
2. Operator resumes via `-resume-import-search-case` with optional `-expect-snapshot-sha256`; existing first-party Case owner validates all metadata and access grants. Browser uses the existing generic `POST /demo-local-search` Event Search Job with a new bounded `field`:
   - `preview` (default when absent): exact existing case-insensitive rendered-preview literal;
   - `rule_id`: numeric Wazuh detection rule ID (1..16 decimal digits);
   - `agent_id`: numeric Wazuh agent ID (1..16 decimal digits).
3. No new provider, DSL or query language. The existing Case/Timeline owner reauthorizes each time-eligible typed source **BEFORE** comparing its retained source metadata, and the backend returns only completed authorized Entry refs. Same max 500 eligible candidates / max 200 hits, no partial results, no raw text or token in SearchJob receipt; per-hit open reauthorizes again. Numeric queries fail closed if this Case's imported sources lack a complete retained Wazuh metadata set (e.g., older generic/legacy saved files). The backend never interprets rule/agent metadata as an ACL or caller authority. Sensitive terms and numeric values are sent via the existing authenticated bounded JSON POST, NOT URLs.

## Proof and limitations

- Two actual disk Wazuh files contain the *same original alert ID* but distinct rule/agent/source info; Case saved privately, exact on-disk SHA rechecked, **original files deleted**, owner reloaded; rule 5710 returns 2 authorized distinct source refs, rule 1000 returns 1, agent 001 one, agent 002 two; case-insensitive preview search still works, missing rule returns completed zero-hit Job.
- Every returned entry reopens via the existing typed-source owner authorization; no unauthorized raw previews or bearer in Job payload. Invalid nondecimal IDs, unsupported arbitrary source selectors, and wrong bearer fail before any source result. Tampering source metadata to nondecimal, level 17 or feed 5 fails owner validation before a listener can be created.
- The exact owner metadata is **synthetic operator-controlled evidence description**, not a verified sensor provenance, production IAM, deduplicated collector, SOC alert verdict or chain of custody. A forged input snapshot with a different valid agent/rule could pass structural validation if the operator separately accepts its hash; the independently retained SHA helps only detect change, not trusted Wazuh origin. A level is a *Wazuh declared severity*, not a CMDR risk score.
- This extends the optional local reference-owner snapshot schema, NOT the canonical production source contract. Avoid interpreting optional metadata on other source types or for tenant access control.
- No external packages, paid APIs, new local service, persistent database, continuous Wazuh collector or production deployment. Single synthetic principal and local loopback only.

## Integration and review

Draft child PR of frozen #70 exact HEAD `0c1e704d23a966f8723f1561c2f429fbe9ab7402`, whose functional and mandatory security CI both passed on that SHA. All ancestor L2 SHAs untouched, but earlier ancestor PRs remain DRAFT and sometimes fail their own security gate. On the new exact HEAD require runtime Go+Node Linux/macOS/Windows and `engineering-kernel` unchanged mandatory 90% changed security coverage. Never lower/waive any gate or merge prematurely. One independent read-only L2 appropriate for this material new saved source metadata/trust boundary after CI.
