# E10-INV-430A — Headless Event Search on one SHA-verified saved Case

## What changed functionally

The existing first-party mixed Wazuh 4.x + generic normalized JSONL import can save a **single authorized private Case** and later execute bounded first-party Event Search Jobs with Wazuh burst and cross-family 60-second leads in the authenticated analyst browser. Until this integration there was **no repeatable batch/CLI hunt** over an already saved Case. An analyst working offline can now replay a selected bounded search without opening a browser, exporting raw source bytes, starting a network listener, or writing credentials to stdout. This reuses the exact first-party Case owner, typed-source access, Event Search orchestrator and correlation implemented in parent PR #82—not another search engine.

### CLI example

```shell
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/previously-saved-authorized-case.json \
  -hunt-saved-case EXISTING_CASE \
  -expect-snapshot-sha256 OPERATOR_RETAINED_64_HEX \
  -hunt-field source_ref -hunt-family all \
  -hunt-value local-jsonl-source- \
  -hunt-from 2026-10-17T09:59:00Z -hunt-to 2026-10-17T10:02:00Z
```

The new CLI flags are **exclusive** of save, resume, import, export and other investigative modes. The SHA-256 fingerprint is **mandatory** and checked on the same bounded regular-file bytes as JSON parsing; it must have been saved by the operator independently. The server derives tenant/environment/principal/Case identity from the already validated source owner. Headless execution creates a random bearer in memory, opens exactly the existing Case/Timeline/typed-source handler **in-process** and POSTs the search criteria to its pre-existing `/demo-local-search` handler, also **in-process** (Go `httptest` request dispatch; no TCP bind).

The handler reopens each typed source before evaluating the search predicate. The existing 500 inspected candidates/200 results and complete pagination/authorization rules apply unchanged. Existing Wazuh-only hits, 60-second bursts and Wazuh↔generic context are derived **only from one completed SearchJob's authorized result references**.

### Output and privacy

A successful invocation emits ONE JSON receipt: `cmdr-local-headless-authorized-hunt-v1`, synthetic marker, Case ID, SearchJob ID, authorized Event Search reference IDs, pre-existing Wazuh hit/burst summaries, and existing mixed temporal context. It does NOT echo the arbitrary literal query value, raw preview, source log message, bearer token, tenant records, audit trail or partial failed job response. An authorized Entry ID is sensitive metadata and should be treated as such; source reopening still requires the first-party owner on subsequent uses.

Wrong/missing digest, unknown/foreign Case, foreign principal, bad query, disallowed family/field, malformed/out-of-order timestamps, exceeding candidate caps, denied/revoked source or incomplete SearchJob ⇒ **no receipt**. Any failure occurs before response serialization, and no listener, third-party service or paid API is involved. Successful temporal context is an *analyst lead*, **not causal attribution or proof of custody**.

### Evidence and constraints

- Integration test imports two Wazuh alerts and one generic JSONL event from real independent operator files into one saved private Case; actual CLI hunt returns three authorized references, two Wazuh hits and two contextual leads (including a generic observation), with no raw descriptions/bearer/listener. Single-family query must not expose mixed context.
- Reject wrong SHA, altered still-parseable JSON, foreign Case/principal, invalid search family, invalid time range and forbidden CLI combinations without a partial receipt.
- Existing integration tests and mandatory `engineering-kernel` **90% changed security-critical code coverage** remain required and unchanged.
- The mode is a local, synthetic, operator-controlled laboratory, NOT production SIEM, live telemetry ingestion, immutable evidence custody, signed origin, IAM service or durable Event Search index. SHA-256 is unauthenticated equality against an independently retained fingerprint.

The **frozen parent** PR #82 is HEAD `06ff4539632eedbf8c98d80dc88f15d5905ec0a4`. The parent runtime CI was green but engineering-kernel coverage **86.52% < 90.00%** and the upstream 30+ stacked draft PRs were unmerged. Do not waive coverage or merge this child before all ancestors, mandatory checks and native approvals are satisfied. Schedule only ONE independent read-only risk-based L2 on a stable exact head if code and CI support it.
