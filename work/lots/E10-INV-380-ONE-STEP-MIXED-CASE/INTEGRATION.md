# E10-INV-380 — One atomic Wazuh + generic JSONL Case import

## Integrated analyst workflow

PR #77 already provides an actual first-party, owner-authorized mixed Wazuh+normalized JSONL Case Search Job: completed Event Search Jobs return both families of source refs, Wazuh-only rule/agent pivots, Wazuh chronology and true 60s temporal bursts, with no generic alert misclassification. Until this lot, creating a durable mixed Case required TWO separate operator-visible disk saves and an intermediate sensitive Wazuh-only owner archive.

This integration provides ONE explicit command and ONE durable private output:

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/authorized-owner.json \
  -into-case EXISTING_CASE \
  -wazuh-alerts /absolute/wazuh-alerts.jsonl \
  -import-jsonl /absolute/generic-normalized.jsonl \
  -save-mixed-case-snapshot /absolute/new-private-mixed-case.json
```

Both source arguments may be repeated with their existing flags (`-wazuh-alerts` up to four; `-import-jsonl-extra` up to three). All operator inputs must be distinct absolute-path regular files; duplicate hardlinks/symlinks/owner aliases are denied. Combined evidence input size is limited to 8 MiB. The selected Case, principal, tenant, source ACLs and owner graph are checked through the existing Case/Timeline/typed-source owner.

### No new ingest engine or transient archive

- The existing strict Wazuh 4.x parser, normalized alert preview selection, existing per-source validated Wazuh rule/agent/level/feed metadata linker, and existing optional multiple-Wazuh importer are reused. The existing single-Wazuh save is refactored to use a shared *in-memory* prepared Case helper with no change in its input/output contract.
- Generic JSONL files use the existing normalized parser and imported source/Entry owner IDs (feed ordinals 1–4). Input ID collisions with a previously imported Wazuh source are rejected by the existing owner importer, not silently merged.
- Only AFTER every source is parsed, canonical metadata linked, owner graph validated and principal Case permissions checked does the existing `persistScopedCaseOwner` build one authorized Case projection. The same private 0600 temporary file, fsync, atomic no-overwrite hardlink and operator SHA-256 receipt are reused **once**. An invalid second feed or denied Case cannot create an intermediate or a partially committed final archive.
- The resulting file is reopenable through the existing `-resume-import-search-case CASE [-expect-snapshot-sha256 digest]` reader with a *fresh* ephemeral bearer, same `POST /demo-local-search` source-by-source reauthorization, Wazuh numeric pivots, 500-candidate/200-hit bounds and 60s burst derivation. Raw source Wazuh JSON, full_log beyond bounded preview, hostname/IP and generic unselected data are not added as independent evidence fields.

### Functional/security acceptance

An actual file-based test runs **one CLI save of three same-agent/rule Wazuh alerts and a separate generic JSONL event**, then reloads saved owner and executes the existing authenticated Event Search Job:
- a text query over both source families returns 4 individually openable authorized source refs;
- true Wazuh 60s burst derives only the 3 Wazuh events, not the generic one;
- numeric Wazuh pivot returns the 3 Wazuh alerts; generic-only search returns 1 generic ref, no invented Wazuh lead;
- raw previews and bearer are never present in Search Job output or CLI receipts;
- deletion of original alert/JSONL files does not prevent owner snapshot reload;
- directory contains only the *single final* archive, not a sensitive intermediate;
- invalid generic/Wazuh feeds, duplicate actual file, owner alias, existing destination, and incompatible CLI modes fail without any archive or partial output.

This is a **local synthetic** investigation workflow, not live Wazuh ingest, privileged production IAM, indexed Event Search DB, transactional custody, signed evidence or incident-verdict automation. Case readable projections intentionally exclude denied data and therefore are not complete forensic raw archives. Source timestamps, provenance and synthetic auth are not cryptographic attestations. A single save still needs one preexisting authorized Case snapshot supplied by the operator.

## Dependency and delivery plan

Base of this new integration is parent PR #77 immutable `2a671a54014d382ba8c090d939ce503c6f77f68b`, **not main**. Upstream source runtimes remain staged in DRAFT PRs, notably #36 foundational engineering and #73–#77 consolidation. No prior L2-frozen HEAD is changed and no early merge is attempted. Run GitHub's exact-HEAD Linux Go + Node browser, macOS, Windows, and the unchanged mandatory `engineering-kernel` changed-security-critical **90.00% floor**, with no waiver or increase in governance complexity. One independent read-only L2 will assess the **material new multi-source filesystem trust boundary**, not be duplicated for each code fix. Unresolved parent CI and PR dependencies still prevent delivery to main.
