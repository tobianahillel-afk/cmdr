# E10-INV-220A — Multi-feed normalized JSONL Case imports

## Functional deliverable

Extend the existing first-party **generic normalized JSONL** operator import and **Event Search Job** in PR #60 from a single local JSONL to **2–4 explicit independent JSONL files** feeding one authorized Case.

Example lab command (requires a preexisting operator-owned owner snapshot and synthetically authorized Case):

```bash
go run ./cmd/case-timeline-demo -principal analyst-a \
  -snapshot /absolute/case-snapshot.json -into-case existing-case \
  -import-jsonl /absolute/firewall-normalized.jsonl \
  -import-jsonl-extra /absolute/endpoint-normalized.jsonl \
  -serve-import
```

Input contract for **each feed** is unchanged `{"id":"...","timestamp":"RFC3339","message":"..."}` with strict UTF-8, max 4096-byte lines, 512-byte messages, bounded regular file (never symlinks, paths only via CLI), no unknown fields. The entire selected set is checked before the viewer starts. A combined **8 MiB** feed content ceiling and the existing **5000-Case-entry** owner limit apply; repeated paths and hardlinks are rejected.

- Every feed receives a deterministic ordinal (1–4) in a derived stable Case/Entry/source identity hash and owner-visible `Provenance: operator-local-jsonl-feed-N`. Two actual files containing an alert with the **same original ID** produce **distinct** authorized Case entries and typed sources, not collisions or overwrites; original source bytes/IDs are not rewritten.
- The **original single-feed** `-import-jsonl` path (including offline export, triage, comparisons, existing `-serve-import` without extra feeds) preserves its existing identity and `operator-local-jsonl` provenance exactly. No canonical spec mutation or breaking API.
- Reuse the exact Case/Timeline/typed-source owner validation after each append and final in-memory graph; no new permissions, fake metadata ACL, database, background ingestion daemon, network listener, search engine or dependency.
- All derived new Case entry IDs for the multi-feed set become the ONLY eligible set in the existing authenticated `POST /demo-local-search` handler. Execute the same first-party Event Search Job over a literal message and selected [from,to) window; within each candidate the owner reauthorizes typed source **before** comparing text. The same 500-candidate/200-hit bounds, whole-job failure on overflow and audit/tenant/Case scope remain.
- Existing analyst viewer navigates every returned result via independent source owner. Search results never include raw messages or bearer.

## Evidence and acceptance

- Real two-file import with colliding original alert IDs, two in-window case-insensitive `PowerShell` alerts returned as **two distinct** Event Search references and individually owner-openable sources; a third unrelated event and a fourth outside the time window are excluded.
- Denied bearer/foreign Case cannot import/search; repeated file path blocked; different files with identical normalized original event IDs stay separate; default legacy single-file IDs/provenance unchanged; extras without explicit `-serve-import` rejected.
- GitHub Actions on exact published HEAD: Go+Node tests and vet on Linux, Go tests Windows and macOS, existing `engineering-kernel` security checks with 90% floor not weakened, plus one independent risk-based L2 on material multi-source/identity change.
- Parent PR #60 fixed L2-reviewed SHA `69d29c8b8b1ecb8d67be54f02215f9d212a8cdaa` is IMMUTABLE, as are all earlier L2 parents. Upstream mandatory `CHK-SECURITY-TESTS` coverage at parent was **89.12% < 90%**, a genuine merge blocker; do NOT reduce requirement or merge while red.

## Explicit limitations

Offline operator-provided **synthetic lab**, not multiple live source connectors, indexed store, continuously updating Case, real IAM, multiple simultaneous analyst accounts, durable Search Job, verified event timestamps, or chain-of-custody evidence. The Case/Timeline source owners remain the sole authority for access; ordinal provenance identifies only the operator-selected import slot, not a cryptographically proven log origin.
