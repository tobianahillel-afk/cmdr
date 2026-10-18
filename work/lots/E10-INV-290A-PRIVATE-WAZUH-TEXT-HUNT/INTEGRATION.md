# E10-INV-290A — Sensitive Wazuh text hunt over the existing private Search Job transport

## Why this is a material integration fix

The first-party generic JSONL Case hunt already transports sensitive literal phrases as authenticated POST JSON. The separate Wazuh 4.x preview hunt, delivered earlier in PR #59, still placed an analyst's arbitrary `preview` literal in **GET query strings** and echoed it into the JSON Search Job receipt. Disabling its download did not remove either disclosure surface. If an analyst uses an internal hostname, user string or other private indicator as a hunt term, the GET URL is sensitive in browser and HTTP access/request traces. These are different transport behaviors for the **same kind of Case-authorized text search**, not justification for another engine.

**This work converges the existing paths:** Wazuh `field=preview` now uses the existing `/demo-wazuh-search` with **authenticated POST application/json**, no URL query, a strict 4096-byte body limit (including trailing bytes), unknown-field/trailing-document rejection, UTF-8 bounded literal (3–128 bytes), and the SAME Case/Timeline/source-open authorization and first-party Event Search Job backend. The response `value` for text searches is empty to avoid echoing a potentially sensitive phrase. The first-party GUI checks this response format and does not export local query replay for text hunts. **Numeric agent_id/rule_id pivots retain their existing GET URL, response value, Job semantics and replay compatibility**. POST numeric and GET preview are denied; all other cohort and drilldown paths remain GET-only.

No raw alert description/preview, bearer, phrase, full_log, or new string-index is written into the Search Job response. The Go backend's existing source-open authorization is still executed BEFORE comparing the rendered preview and obeys max 500 candidate previews and max 200 resulting refs with fail-closed overflow. Source click still reauthorizes via first-party owner. There is no new endpoint, search algorithm, permission provider, IAM, external package or paid API.

## Behavioral evidence

- Real disk Wazuh import across two feeds containing the same original Wazuh alert ID, Case owner and bearer, actual first-party `Search Job` returns two independently authorized matching event IDs for the case-insensitive `powershell` literal. Receipt has empty value, no raw alert preview or bearer; an unrelated alert is excluded.
- Wazuh GET `field=preview`, unauthenticated/wrong-Case/invalid literal, malformed times, POST query string, unknown JSON fields, trailing JSON, unsupported Content-Type, over-4096-byte body all produce denial without result refs or private text.
- Browser sends sensitive literal via POST body with no query parameters (same as normalized JSONL); maintains UI state/revision checks, independently source reopens and disabled local JSON replay. Existing numeric GET authorization and replay remain regression-tested.
- Existing Linux/macOS/Windows Go, Linux Node UI and mandatory `engineering-kernel` security coverage 90% floor are unchanged and must be verified at exact HEAD.
- **Trust scope:** only synthetic locally supplied Wazuh feeds / Case. POST removes URL/response disclosure; it does NOT encrypt the saved in-memory evidence, guarantee body privacy against a compromised endpoint, add live SIEM ingestion or replace IAM/forensic custody.

## Integration / merge

Child of frozen L2 PR #67 SHA `db8a5dcb48692b7f9b5c73bdf6107959e4155c40`. No ancestor L2 SHA changed. Parent #67 and the #36 foundation/long draft stack remain unmerged; this PR must remain DRAFT unless mandatory CI including 90% security-critical coverage, independent material L2, and all upstream rules pass. Parent #67 measured 86.41%/90% FAIL and must not be waved away by a green downstream workflow. No new L2 on unchanged parent or on microcorrective tests.
