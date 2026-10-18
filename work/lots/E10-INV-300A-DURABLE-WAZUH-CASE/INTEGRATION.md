# E10-INV-300A — Wazuh 4.x → durable authorized Case → offline search

## Integrated functional delivery

The current offline Case runtime already has two working pieces that were isolated:
- Real operator-supplied Wazuh 4.x `alerts.json` → first-party Wazuh 4.x normalizer → Case/Timeline/typed-source owner with full **temporary** Wazuh-specific viewer (one to four Wazuh input files).
- Generic normalized JSONL → private **Case-only least-privilege snapshot**, atomic create/no-overwrite, SHA-256 digest, verified resume, searchable offline Event Search Jobs and independent saved Case comparison.

This PR joins those existing components. An analyst can explicitly preserve already-normalized, authorized Wazuh events to a durable offline Case snapshot without writing the raw multi-tenant owner dataset or developing a second persistence engine.

### Save a real Wazuh investigation

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a -snapshot /absolute/original-case-owner.json \
  -into-case EXISTING_AUTHORIZED_CASE \
  -wazuh-alerts /absolute/wazuh-collector-A.jsonl \
  -wazuh-alerts /absolute/wazuh-collector-B.jsonl \
  -save-wazuh-case-snapshot /absolute/new-private-wazuh-case.json
```

Strictly 1–4 Wazuh 4.x files, existing normalizer only. Distinct operator files, no symlink or alias to the owner, full bounded parsing and synthetic principal/Case/typed-source authorization **before any write**. Two files with the same original Wazuh alert ID receive distinct canonical Case/Entry/source IDs, reusing the already implemented Wazuh multi-feed namespacing; no input bytes are modified on disk.

The **same** `persistScopedCaseOwner` routine is reused for generic JSONL and Wazuh: minimize to the selected principal and Case and sources/entries whose respective Readers authorize that principal, validate the resulting existing first-party owner graph, serialize an 8-MiB-bounded dataset, write mode 0600 local temporary file, Sync, then publish with atomic no-replace hard-link. Output is merely path and independently reusable SHA-256 digest; no bearer/preview or raw alert text in CLI output. Saving does not open a server. Existing destination, symlink, invalid/foreign Case and conflicting CLI modes fail closed without overwriting evidence.

### Resume, investigate, compare

After saving, even if both original Wazuh `alerts.json` inputs are deleted, use the **existing**:
- `-snapshot /absolute/new-private-wazuh-case.json -verify-import-snapshot-case EXISTING_AUTHORIZED_CASE -expect-snapshot-sha256 HASH` to verify the snapshot offline with no listener;
- `-snapshot /absolute/new-private-wazuh-case.json -resume-import-search-case EXISTING_AUTHORIZED_CASE -expect-snapshot-sha256 HASH` to relaunch the same local Case/Timeline/source viewer, fresh bearer, and real first-party `POST /demo-local-search` to hunt rendered normalized alert descriptions with owner-authorized source drilldown;
- optionally compare this Case with another distinct, independently verified saved Case via the parent `-compare-saved-case` path.

**Important fidelity boundary:** the saved snapshot retains Wazuh event time, a bounded analyst-facing `rule.description` (or bounded `full_log` fallback), Case Entry IDs and typed-source previews. It DOES NOT retain original full raw Wazuh JSON, agent ID, rule ID/level metadata, signed collector identity, or the specialized live-process-only Wazuh rule/agent cohorts. After resumption, the operator receives **generic normalized Case text hunting**, not the original specialized Wazuh numeric-pivot engine. This is a declared, functional limitation and a next possible integration once Wazuh metadata has a durable, authorized schema. A local SHA-256 checksum is not a signature, tamperproof custody or anti-rollback guarantee.

No network ingestion, paid API, external library, new search engine, generalized persistence wrapper, real tenant IAM, durable database or product-wide architecture change. Canonical OPEN-013 remains unresolved.

## Behavioral acceptance

- Two genuine local Wazuh 4.x files with the SAME source event ID but different `rule.description`; explicit CLI save; exact on-disk SHA-256 recomputed from bytes equals printed receipt; authorized saved Case has one principal/Case and three imported distinct Entry/source IDs.
- First-party Event Search Job after disk reload selects case-insensitive `PowerShell` across both feeds, reports two distinct authorized refs without raw messages or token, independently re-opens both typed-source previews. Remove original Wazuh files; offline verification and real CLI resume with new bearer still work.
- Single Wazuh file save works. Other principal/unknown Case, existing evidence destination, relative destination, missing Wazuh inputs or combined Wazuh/JSONL/viewer/save modes reject without receipts or overwriting existing evidence.
- Existing Linux Go+Node, macOS/Windows, GitHub `engineering-kernel` and strict security-critical changed-code 90% gate are unchanged. **No merge while any gate fails.** One independent risk-based read-only L2 on exact SHA once functional CI stabilizes.

## Integration constraints

Parent PR #68 frozen SHA `f3870051946f83512dc5e9473ba8d953cba71b13`. PR #36 engineering foundation and all stacked PRs remain unmerged. The L2 review of this bridge must not mutate any previously reviewed SHA. The entire original owner graph is NOT written to disk: the export is a readable projection, not a complete forensic chain-of-custody archive.
