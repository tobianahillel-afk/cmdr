# E10-INV-260A — Offline authorized Case integrity and source-inventory verification

## Usable analyst journey

The existing local synthetic Case workflow can import up to four normalized JSONL evidence feeds, save a **single authorized Case projection** to a private snapshot without overwriting evidence, optionally resume an authenticated Event Search after restarting, and compare exact bytes of that file against a SHA-256 retained by the operator.

This lot adds an explicitly **offline, read-only verification step**, useful before enabling an investigation server or handing a locally retained snapshot to another authorized offline analyst. It does **not** launch HTTP, issue a bearer token or touch saved evidence bytes:

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/path/to/previously-saved-authorized-case.json \
  -verify-import-snapshot-case EXISTING_CASE \
  -expect-snapshot-sha256 <64-hex-digest-recorded-at-save>
```

The existing first-party loader opens one bounded regular snapshot file, verifies the operator-held digest **against the exact bytes it decodes**, parses strict JSON, applies the same tenant/environment/principal Case and typed-source owner scoping as the online resume, and recognizes only normalized imported Event/source IDs. The only receipt is a JSON object with `format: cmdr-offline-authorized-case-verification-v1`, `synthetic: true`, `case_id`, `authorized_imported_events` and counts per existing synthetic `operator-local-jsonl[-feed-N]` provenance. No raw log lines, snippets, reference IDs, bearer, endpoint, other Case or tenant details are printed.

The fingerprint is **required** for this new offline verifier. Saving and online resuming retain backward compatibility; `-expect-snapshot-sha256` cannot be used as a way to bypass Case ACL or to combine modes. CLI modes save/resume/verify are mutually exclusive. Fail closed without any receipt if SHA is wrong/missing/malformed, snapshot unreadable/tampered, source owner inconsistent, principal or Case not authorized, canceled context or output write fails.

### Scope and limitations

- The digest must be retained independently of the saved file by a trusted operator. This is **byte-integrity / unexpected-change verification**, **not authenticated origin, anti-rollback, digital signature, custody seal or a production IAM grant**. Attackers who can replace both the digest and file defeat it.
- A successful offline receipt proves only that **the selected Case's imported entries remain readable in this synthetic operator-owned owner graph** and that its bytes match the independently retained digest. It does not verify the underlying source feeds or historical event provenance beyond synthetic input-slot labels, nor can it assert that every source in the original Case was exported (denied entries are intentionally excluded).
- The saved file remains a sensitive local Case source archive; do not put it in GitHub, logs or public file sharing. No database, continuous ingestion, production multi-tenant IAM or signed forensic evidence chain.
- Reuses EXACT existing first-party snapshot loader, Case owner scoping, strict imported identity checks and synthetic source permissions. No new general framework, network listener, external library or paid API.

### Real integration acceptance

- Two real independent JSONL input files containing the SAME original event ID import and save as two different owner source IDs. The offline CLI verifies the independently recomputed SHA-256, reports exactly 2 authorized events partitioned by feed ordinal and returns without a listening socket, bearer, raw preview or source reference.
- Original evidence files removed after save; verification still succeeds.
- Missing/invalid/wrong expected digest, syntactically valid but changed JSON, missing snapshot, unauthorized principal, foreign Case, canceled context, failed receipt writer and invalid mode combinations reject without partial output.
- Existing Go source-read, snapshot saves, Event Search and UI tests remain unchanged and continue running Linux, macOS, Windows. Mandatory engineering-kernel changed security-critical coverage 90% is NOT lowered or waived. Never merge while any gate is red.

## GitHub integration decision

Draft child of frozen PR #64 exact SHA `b90d578a28dc43213162a0fdddec3457b19abac3`; preserve all parent L2 immutable commits. This new operator end-to-end workflow deserves one independent risk-based, read-only L2 on its exact diff after GitHub Actions stabilize. Future merge consolidation must reconcile PR #36 and the long parent chain; this lot itself does not unblock their separate CI checks.
