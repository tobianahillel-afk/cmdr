# E10-INV-540A — Cross-Case Wazuh event pair → separately authorized private source inspection

## Analyst capability

PR #93–#94 deliver first-party rule/agent-linked event-pair leads (<=60s inclusive), each with distinct `first_entry_id` and `second_entry_id` pointing to **different already-authorized Case owner snapshots**. A correlation result alone must never be accepted as authorization to read a source. Until this lot the pair identified both candidates, but the operator still lacked a direct private evidence-open path tied to separately verified saved owner files.

This lot connects that already published correlation output with the **existing** `inspectDemoEntryWithData` source inspection flow. The operator selects a side of the pair, independently presents its snapshot, selected Case, operator-held SHA-256, and a new private output path:

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/first-saved-case.json \
  -inspect-saved-case FIRST_CASE \
  -inspect-saved-entry FIRST_ENTRY_ID_FROM_PAIR \
  -expect-snapshot-sha256 FIRST_HEX_DIGEST \
  -private-report /absolute/first-source-inspection-new.json
```

Repeat for the second Case snapshot, **using the second Case's own independent digest and second Entry ID**. No implicit cross-Case permission is inherited. Both files contain potentially sensitive source previews and must remain under operator-controlled filesystem permissions.

### Mechanism and security

- Reuses `loadAuthorizedImportedCase`: bounded one-read exact-byte SHA-256 verification, first-party principal/tenant/environment/Case/typed-source scoping and strict canonical normalized imported Entry identity checks.
- The supplied Entry ID is only a selector. The **original first-party** `inspectDemoEntryWithData` regenerates its in-process bearer and performs the real owner-authorized Case, Entry and typed-source open operation. No new source handler, RBAC engine or trust shortcut.
- Reuses `writeAuthorizedReceiptPrivate`: the entire authorized JSON is staged, then privately created at an operator-given **absolute new path** through the existing atomic no-overwrite writer. Source text, bearer and resulting receipt are NEVER placed on stdout, in URLs, a public listener, or in cross-Case comparative output; failed authorization or a preexisting path never creates or truncates evidence.
- CLI mode is exclusive with all other save/verify/compare/import/server modes, and requires `-private-report`. No fallback to stdout when the path is missing, relative or empty. No third-party runtime dependency or paid API.

### Material behavioral acceptance

- Two real Wazuh 4.x JSONL alerts with same agent/rule and 20-second separation are imported **into two different Case owner files**, then privately saved, SHA-256 fingerprints recomputed, original raw JSONL files deleted, and the existing comparison produces one concrete pair.
- The test then opens each pair side **independently through the real CLI** with its proper Case/digest/Entry ID; each new private file includes ONLY that side's owner-authorized source preview. No raw source text or token on stdout.
- Switching Entry IDs across Cases, selecting a foreign Case/principal, wrong or malformed independent SHA, changing the file, combining CLI modes, missing private output or overwriting existing evidence fails without writing a report.
- Existing mixed-source Case search, first-party correlation, source click, private Case storage and other security tests continue unchanged.

### Scope and release condition

This is an **offline synthetic operator-controlled investigation lab**, not production tenant IAM, authenticated custody signing, evidence attestation, an online SIEM source viewer or causal Wazuh correlation. A SHA-256 retained by the operator verifies bytes against a separately held checksum; it is not a trusted signature and does not prevent an adversary who controls both file and expected digest.

Frozen parent PR #94 SHA `9c8ceb8353d1f5c3a8abf5b548ce83948330b66a` remains untouched. New PR requires its own exact-head Linux Go+Node, macOS, Windows and unchanged mandatory engineering-kernel 90% changed security-critical coverage, PLUS one independent read-only risk-based L2 on this material source-release boundary when verifiable. Prior #94 coverage 81.42% <90, prior #92 full integration likewise red, foundational #36 still draft. **No merge, CI waiver or false production claim.** Next engineering objective: reconcile the long draft chain and remove real blockers; avoid decorative tests and new governance.
