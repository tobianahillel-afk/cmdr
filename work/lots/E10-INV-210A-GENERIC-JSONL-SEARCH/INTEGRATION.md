# E10-INV-210A — First-party Event Search over generic normalized JSONL

## Deliverable: extend investigation beyond Wazuh

Existing CMDR generic JSONL import (`-snapshot ... -import-jsonl /absolute/events.jsonl -into-case ... -serve-import`) already normalized three fields (`id`, RFC3339 timestamp, bounded `message`) and imported them into a synthetic Case owner. Until now, users could browse its Timeline but execute real `Event Search Job` only against **Wazuh-specific** metadata.

This incremental delivery makes the generic JSONL import genuinely searchable in the **existing** analyst screen using the existing first-party `event-search` orchestrator, with a simple selected Case, literal rendered-message query, inclusive start/exclusive end and explicit Execute. It is a useful multi-source foundation because the import contract is not tied to Wazuh.

### Operator path

1. Operator selects a validated Case snapshot and normalized JSONL event file via the existing CLI. The importer's unchanged strict file path, regular-file, size, UTF-8, event identity, schema and timestamp rules construct an in-memory Case/Timeline/typed-source owner.
2. Only for this **generic JSONL import mode**, attach a new `POST /demo-local-search` to the preexisting ephemeral IPv4 loopback, same bearer auth and per-Case owner. The server derives eligible imported event IDs from the importer's own appended, already-validated owner entries, never the query or unrelated pre-existing Case sources.
3. Browser POSTs Case, literal text (3–128 UTF-8 bytes, no line breaks) and RFC3339 UTC bounds in a bounded JSON body; NO bearer or sensitive literal text in URL. No CORS, external request, third-party service or new listener.
4. Backend builds an actual first-party Event Search Job with tenant/environment scope from the server's known synthetic principal, source `normalized-jsonl-authorized-import`, existing explicit synthetic search permissions, query provenance and audit trail. It traverses all authorized Case Timeline pages, skips other Case entries and events not appended by this import, handles certain timestamps in `[start,end)`, and calls the **existing typed source-open owner BEFORE any literal comparison**. A denied/changed source or incomplete page fails the entire Job.
5. The engine compares the owner's bounded rendered message case-insensitively (literal substring, no regex or DSL). Max 500 authorized time-filtered candidate previews and max 200 matching refs; overflow fails the entire Job, no silent partial results. The response contains only the completed Search Job with refs and existing provenance/audit, NOT messages/preview/raw logs. Each click goes through the existing **independent** source navigation and source-open reauthorization.

### Acceptance and security boundaries

- Real disk JSONL -> first-party normalizer -> in-memory owner -> POST authenticated search -> completed Event Search Job -> independently owner-openable telemetry result.
- Positive case-insensitive literal within time, negative unrelated/late existing source exclusion, completed zero-match result; foreign Case, missing/wrong bearer, malformed literal and timestamps, GET rejected, no raw preview in returned Job, 501 candidate overflow fails wholly. Browser tests prove POST body but no sensitive URL, source ref rendering, invalid receipt denied and identity changes clear results.
- Generic import and Wazuh modes are isolated. Normal snapshot-only viewer does NOT register `/demo-local-search`. No direct raw JSONL read at query time, no external API, no new search/indexing framework, and no new authorization model.
- This is **offline synthetic lab IAM**, not real tenant production IAM, continuous/live collection, persistence/index, transactionally consistent custody, durable Search Job, signed forensic evidence or broad normalization pipeline.
- Code deliberately reuses the first-party Event Search orchestrator, `collectDemoCaseTriageWithObserver`, `queryDemoOwner`, the existing bounded source-open handler and the existing literal validator. No new generalized abstraction or repeated safety wrapper.
- GitHub exact-head Go+Node tests Linux/macOS/Windows, mandatory engineering-kernel 90% security coverage gate, and ONE targeted independent read-only L2 required for new generic data access; no merge while blocking CI or parent reviews fail. Parent PR #59 frozen `af987750ef9218ea441b1ae5a4e6f525b4cf30a9`. Reconcile this stack with existing independent PR #52's Case Search Workbench rather than inventing a second persistence model; OPEN-013 remains open.
