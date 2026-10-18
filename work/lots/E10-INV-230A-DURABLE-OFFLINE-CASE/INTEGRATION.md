# E10-INV-230A — Explicit private Case snapshot persistence and authorized search reopen

## Functional outcome

The existing first-party **normalized JSONL** (1–4 sources), synthetic Case/Timeline/typed-source owner and **Event Search Job** are now usable beyond one temporary process session, without inventing a new persistence engine or changing canonical Case mutation contracts.

**Operator explicitly saves** the same fully validated in-memory owner graph that the importer already builds, including Case, Timeline and source previews, to a NEW local absolute destination. Saving neither starts a listener nor exposes records on stdout. The code marshals the entire bounded owner graph first, writes a private 0600 temporary file in the destination directory, syncs and closes it, and **publishes with an atomic no-overwrite hard link**. Existing regular files, symlinks, a duplicate path and creation failures are never replaced or truncated. Failures remove temporary files. The completion message contains only the path and SHA-256 digest for accidental corruption checks, **not** a signature or trustworthy chain-of-custody proof. The output may include ***other tenants and sources contained in the operator-provided original snapshot*** and must be handled as sensitive.

**Operator explicitly reopens** an existing validated saved snapshot in the same synthetic Case/Timeline/typed-source viewer; the selected Case is required and the server derives eligible normalized import IDs from existing `operator-local-jsonl` or `operator-local-jsonl-feed-N` provenance together with exact first-party canonical Entry/Source IDs (not a client-provided ID list). Reuse the **same** existing `POST /demo-local-search` first-party Event Search orchestrator, fully paginated Case authorization, per-candidate typed-source re-open **before** rendered text comparison, 500-candidate and 200-hit fail-closed bounds, and fresh authorization of opened hit refs. The new process **generates a new ephemeral bearer**; the snapshot stores NO token or temporary SearchJob results.

### Manual usage (operator filesystem only)

Save two example offline feeds, then stop; paths must already exist except the output destination:

```bash
go run ./cmd/case-timeline-demo -principal analyst-a \
  -snapshot /absolute/owner-original.json -into-case EXISTING_CASE \
  -import-jsonl /absolute/endpoint.jsonl \
  -import-jsonl-extra /absolute/firewall.jsonl \
  -save-import-snapshot /absolute/new-private-owner.json
```

Reopen after restart, even if the original JSONL files are not present:

```bash
go run ./cmd/case-timeline-demo -principal analyst-a \
  -snapshot /absolute/new-private-owner.json \
  -resume-import-search-case EXISTING_CASE
```

The loopback analyst interface is served by the pre-existing authenticated Case owner; browser POST hunts the saved normalized messages and result clicks reauthorize real typed-source owner references. Single normalized JSONL save is also supported; the original single-source identity and provenance remain unchanged. The save/resume flags are **mutually exclusive and cannot combine** with other report, export, Wazuh or interactive import modes.

### Security and limits

- **No automatic persistence**. Only an explicit `-save-import-snapshot` may write the sensitive owner graph. Output must be an absolute path. Never overwrite or follow a preexisting symlink; no broad logging of messages, Case source payloads or bearer.
- The original owner snapshot can include other Cases/tenants, even ones NOT visible to the operator's selected principal. This is an **operator-trusted local lab** and the saved file is sensitive. Do not use the raw saved graph as a multitenant production export. Filesystem access controls and Windows ACLs depend on the operator's directory/OS.
- Data bytes are `Sync`ed before an atomic hard-link publication, but no assertion is made about all host filesystems, sudden power-loss directory synchronization or tamperproof/custody signatures. A SHA-256 displayed alongside an editable file is **not authenticated evidence**.
- Only the same synthetic Case principal and verified typed-source owner provide read rights; provenance and hashes are **not authorization**. The owner snapshot remains local and operator-provided. No live ingestion, database transactions, persistent search index, server IAM, key management, retention policies or production SLO.
- If the underlying filesystem denies hard-link publication, the operation fails without corrupting another file; we do **not** fall back to an unsafe replace operation.

### Behavioral acceptance and CI

- Existing strict normalized importer, existing first-party owner validation, original source reader, Event Search and browser: no new dependency or paid API.
- Test actual CLI save of two files with identical original IDs, reload from disk, derive exact two imported source identities, rerun an actual authorized Event Search Job with two hits and **independently reopen** both saved source previews. Test original input files removed after save, wrong bearer denied, sensitive message absent from SearchJob/CLI response, output POSIX 0600, existing destination unmodified, dangling symlink/relative path/cancelled transaction rejected, invalid mode combinations rejected, edited snapshot source identity rejected and unrelated Case not searchable.
- Exact-HEAD GitHub Actions Go/Node Linux, Go macOS/Windows and mandatory engineering-kernel security coverage 90% (NEVER waive). Parent #61 L2-frozen HEAD `b70341836e47f3d590dc89fc4aef943f8faa6648` remains unchanged; all other ancestor PRs similarly unchanged. Only one material independent read-only L2 on the published SHA, not on every test/CI fix.

This is a **real local durable operator workflow**, not production CMDR persistence, a signed forensic chain of custody, or a final merge. Canonical `OPEN-013` mutation/persistence decisions are not claimed resolved.
