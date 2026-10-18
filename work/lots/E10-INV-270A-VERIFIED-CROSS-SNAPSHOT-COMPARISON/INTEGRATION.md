# E10-INV-270A — Compare two independently verified saved Cases

## Integrated functional journey

The original local CMDR demonstration already implements (a) strict operator-local snapshot import and minimal authorized saved owner projection; (b) independently operator-held exact SHA-256 byte verification of saved Case files; (c) a **first-party two-Case Timeline triage/comparison engine** via `collectDemoAuthorizedCaseTriage` and `buildCaseComparison` producing observed/inferred deltas, known UTC overlap and ordered shared certain 5-minute windows. Previously the comparison worked only for TWO Cases held in ONE owner snapshot. This lot integrates that existing comparison runtime with TWO DISTINCT previously saved and independently integrity-verified operator Cases. It requires NO new comparison/correlation algorithm, no copying combined raw source events and NO merging of their source-owner graphs/ACLs.

### Operator command (offline, no listener)

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/private-case-A.json \
  -compare-saved-case CASE-A \
  -expect-snapshot-sha256 A64HEX \
  -compare-saved-with-snapshot /absolute/private-case-B.json \
  -compare-saved-with-case CASE-B \
  -expect-other-snapshot-sha256 B64HEX
```

Both digests MUST be independently retained by the operator at the time of each save; there is no unverified variant of this cross-file capability. Existing first-party bounded JSON loader verifies each file's exact bytes BEFORE interpreting data. Each owner undergoes independent synthetic principal/tenant/environment/Case ACL scoping and check of canonical normalized imported event IDs, so a source or Case visible in one file does not confer access in the other. **Both Case principals must be within the same explicit tenant AND environment** and Cases must have distinct IDs. Each source owner is exercised through its own in-process, freshly authenticated first-party Case/Timeline/source HTTP handlers (no public socket/bearer output), collecting complete authorized pages and UTC aggregates independently. The existing `buildCaseComparison` then produces aggregate-only retrospective leads: per-Case source/count summaries, observed/inferred deltas, time overlaps, up to five co-occurring known five-minute windows.

No raw source preview, contents, complete event stream, bearer, target IP or hostname from the selected local inputs is copied into the comparison output. `top_sources` contains previously authorized source references, not raw source text; the operator is already authorized to view those references. No result is serialized until both Cases finish authentication, complete owner traversal, and first-party comparison building.

### Risk constraints / acceptance

- Real disk Case A/B, created as separate owned snapshots from 2 separate normalized JSONL inputs, both with independently recomputed operator-held SHA-256 values; original files contain known UTC events in the same 5-minute window. CLI produces a completed combined comparison with both distinct Case IDs and co-occurrence leads, without starting a listener.
- If either digest is missing, wrong, malformed or one still-parseable file has changed, no partial first-Case data is emitted. If either Case is not readable, the principal is foreign, Case IDs coincide, or owner tenant/environment differs, fail before output.
- Explicit operator mode is mutually exclusive with existing snapshot save/resume/verify, ordinary Case comparison, import, Wazuh mode and local listener; other commands unchanged. No third-party library or paid API.
- The algorithm reuses the existing first-party **comparative timeline engine**, not a new SIEM correlation model. Correlation is *temporal co-occurrence*, NOT causal attribution; no claims about identical source IDs between separate Cases, chain-of-custody signatures, production IAM, immutable datastore, source truth or continuous collection.
- SHA-256 is an **unauthenticated integrity fingerprint**: this only proves equality to independently held expected bytes, not provenance against an adversary changing both hash and snapshot.
- Synthetic operator-controlled data only. Explicit untrusted boundary validations remain server-owner authorization and Case/typed source/tenant+environment isolation. Case-specific source refs remain scoped, even if they resemble references from another snapshot.

### Required gates and merge status

Parent exact frozen PR #65 SHA `fdc088e01163a3e61a2656c47824404111b3fa01`; do not mutate audited parent SHA. Existing `main` remains specs, large #36 foundation and all ancestor PRs still draft. Exact-HEAD GitHub Actions Linux Go+Node and macOS/Windows, mandatory `engineering-kernel` changed-security-critical 90% floor (parent #65 only 89.69%), and one independent READ-ONLY risk-based L2 for material two-owner integration needed. No override/threshold reduction, no premature merge.
