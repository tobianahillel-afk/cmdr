# E10-INV-720A — Paginated source-authorized Wazuh event-pair discovery

## Functional gap and choice

The first-party DFIR Case comparator already correlates the same Wazuh rule+agent from TWO independently authenticated Case owners within +/-60 seconds. Its default receipt intentionally includes **only ten** earliest, deterministic event-linked leads. PR #114 permits private drilldown by exact first/second Entry IDs beyond those ten, but does not provide a simple bounded discovery path for analysts who do not already know the IDs. Do not add another search engine or expose raw source previews just to navigate more evidence.

This lot **reuses the existing chronological nearest-peer algorithm** and offers explicit paginated metadata-only traversal. It preserves the original top-ten comparison receipt's exact format and behavior. One page contains up to ten paired Entry IDs, with a bounded offset and a `has_more` boolean. Only eligible real 60-second matches count toward the offset, so pages are deterministic and do not duplicate the same Entry on a page boundary. The authorization and expected source scopes are inherited from two independently verified saved Case snapshots.

## Operator journey

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/saved-first-case.json \
  -compare-saved-case CASE_A -expect-snapshot-sha256 A64HEX \
  -compare-saved-with-snapshot /absolute/saved-second-case.json \
  -compare-saved-with-case CASE_B -expect-other-snapshot-sha256 B64HEX \
  -list-saved-wazuh-pair-offset 10
```

The output is **only** JSON metadata `cmdr-local-authorized-wazuh-pair-page-v1`, selected Case IDs, offset, ten or fewer `sharedAuthorizedWazuhEventPair` leads, and `has_more`. Use `offset=0` to see the original first ten; choose offsets 10, 20, ... only while another page exists. The exact Entry IDs allow the already existing private `-inspect-saved-pair-first-entry ID -inspect-saved-pair-second-entry ID -private-report /absolute/NEW.json` workflow to reopen *both* owners and each selected typed source before writing sensitive descriptions. The top-ten default comparison, ordinary offline verification and inspect-by-ordinal behavior remain unchanged.

Security:
- Both operator-held SHA-256 digests are required, verified against the exact bounded file bytes used for JSON decode. Each principal's Case, tenant/environment, Timeline and typed-source owner access is independently checked using the SAME first-party loaded-owner path and full Case pagination before any pair metadata is emitted. Entry IDs are only **navigation selectors**, never authorization grants.
- No cross-tenant or cross-environment pairing, combined owner graph, continuous ingest, remote access, IAM production, new dependency or paid API. No HTTP listener/bearer and no source previews in stdout.
- Paging is constrained to `0..5000` eligible candidates (the existing max per Case), ten pair leads per page and a one-item lookahead. Invalid offsets, mismatched digests, denied Cases, incompatible `-inspect-saved-pair` or generic evidence expansion, and unrequested private output are rejected **before any metadata receipt**.
- The leads show retrospective *time co-occurrence* under an operator-supplied synthetic Wazuh snapshot, not causal attribution, globally unique event identity, production custody or proof of network activity.

Risk-based tests: persist two genuine independently owned Wazuh files with twelve first-Case events and one second-Case peer; page 0 yields ten and `has_more`, page 10 yields two and no more, page 12 empty. A selected 11th Entry from page 10 is actually opened into the existing 0600 no-overwrite two-source private dossier with no stdout evidence leak. Deny wrong second SHA, foreign principal/Case, parseable-but-tampered second file, out-of-bounds offset and incompatible sensitive options. Existing macOS/Linux/Windows Go and Linux Node UI CI must pass at exact HEAD. Unchanged engineering-kernel 90% changed-security-critical coverage is mandatory; no waiver or premature merge.

Parent PR #114 immutable SHA `9f43ceab9f96ac7648429f2361bf8b4324861786`; frozen ancestor PR #112 SHA `4d6c8057557a68657b9b64584c4f41bac17b3868`. Neither SHA is modified. New work stays in a single cohesive draft PR and should receive ONE independent risk-based READ-ONLY L2 after publication and verifiable exact-head CI. The large foundation PR #36 and other draft integration candidates remain unmerged: this feature does not pretend to solve that merge backlog.
