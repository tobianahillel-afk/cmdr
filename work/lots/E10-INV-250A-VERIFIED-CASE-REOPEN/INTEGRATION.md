# E10-INV-250A — Byte-verified resume of an authorized offline Case

## Product capability integrated with the existing owner

The first-party local normalized JSONL Case already supports 1–4 operator-approved feeds, owner-authorized Event Search Jobs, private explicit no-overwrite snapshot export, and the Case-only least-privilege projection in PR #63. The export CLI prints a SHA-256 fingerprint computed on the exact saved bytes, but until this lot nothing could check those bytes when the operator reopened the Case.

**New deliberate analyst journey:** an operator keeps the hex SHA-256 value printed by the original save and supplies it at resume:

```text
case-timeline-demo -principal analyst-a -snapshot /absolute/private-case.json \
  -resume-import-search-case CASE-ID \
  -expect-snapshot-sha256 0123456789abcdef...64-hex-characters...
```

The existing bounded regular-file snapshot loader now reads exactly ONCE, checks size, inode/regular-file protections, computes SHA-256 **on exactly the bytes it is about to JSON-decode**, compares it against the operator-supplied 32-byte expected digest, then parses JSON and delegates to the unchanged Case/principal/Timeline/source owner validation and least-privilege projection. Any content modification, truncated hash, invalid hex, mismatched digest, non-file or invalid JSON fails *before a loopback listener or fresh bearer is issued*. This avoids a naive second file read with a race between hash-checking and parsing. The mode is explicitly limited to `-resume-import-search-case` and cannot be combined with saving or unrelated report/ingest modes. No HTTP client may supply a digest to override source ACL.

Successful verified reopen uses the **existing** first-party Event Search Jobs over authorized normalized preview text, typed source-open ACL for every candidate, 500/200 bounds, same UI, and newly generated bearer. The historical resume CLI without the digest still works for backwards compatibility, but operators handling sensitive evidence should use the explicit verified path. This is **not a signature, cryptographic provenance from a trusted party, anti-rollback guarantee, or authenticated forensic custody**: an adversary who may change both local bytes and the operator's expected digest can bypass an unauthenticated hash. The guarantee is precise accidental-corruption / unexpected-change detection against an independently retained expected hash.

## Behavioral tests / acceptance

- Actual CLI save to new private owner snapshot and emitted SHA-256 hex: independently computed SHA-256 of saved bytes MUST match the printed receipt.
- Actual CLI resume with correct operator SHA launches loopback with a fresh bearer and no leaked source preview. Existing authentication and Case-only owner validation apply.
- Wrong hex digest and invalid syntax fail. Appending a harmless newline produces still-valid JSON but MUST be rejected against the previous digest, before any HTTP listener or bearer disclosure.
- Expected digest flag is refused for unrelated CLI modes, including save.
- Previous tests proving cross-tenant redaction, owner-only source access, reopen without raw JSONL files and atomic no-overwrite publication all remain in CI.
- No new storage engine, third-party library, paid API, IAM prod, live provider or custody claim. No duplicate validator framework.

**Integration:** child PR of frozen #63 HEAD `061dac08f15591c88d870e349824860a6b4fe05b`, with all parent SHA/PRs untouched. Exact-head GitHub Actions Linux Go+Node, macOS/Windows and engineering-kernel coverage **90% security floor unchanged**, no merge while red. One material risk-based independent READ-ONLY L2 after stable CI (the feature alters the boundary between on-disk data and authenticated Case bootstrap). The former full-graph snapshots are still sensitive files on disk; byte verification does not grant an unauthorized principal access.
