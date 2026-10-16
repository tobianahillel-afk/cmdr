# E10-INV-110D-CLOSURE — bounded verification at exact source SHA

Source SHA: `06d1d47572b9fbea1325ad71309a895d960e1cbd`; PR #36 (draft); four successful GitHub Actions: engineering push #38963830261, engineering PR #38963841717, Timeline push #38963830462, Timeline PR #38963841785.

The checked-in two-tenant demo now has one **actual HTTP end-to-end** acceptance scenario for each synthetic analyst: authenticated Case discovery, timeline kind selection, canonical owner reauthorization on navigation, independent source preview and safe return, plus anonymous/tenant/source denials. All tests ran on Linux, and Timeline portability jobs on Windows/macOS succeeded. Go race and vet run in the Timeline workflow.

Go module total statement coverage: **91.4%**; projection package 93.1%; demo command package 79.6%; performance probe 91.3%. Distinguish global runtime threshold from isolated demo package figure. No code was added solely to hit a numeric percentage.

Benchmarks on the exact code (Timeline push #38963830462): 5000-entry synthetic projection 3.055592–3.317349 ms/op; authorized **navigation HTTP handler** 9370, 11251, 11383 ns/op, 3088 B/op, 32 allocations/op on this hosted runner. These are per-iteration averages, not p95 request latency, production authorization, storage/network latency or a release SLO. Earlier on-demand projection batch-mean p95 belongs to another SHA. The policy registry navigation on-demand handler remains **preimplementation**, and this check does not fraudulently replace it.

**Bounded D verification** means reference runtime, journey and negative/security/performance evidence passed the exact-SHA local CI gates. It does NOT mean general CAP-INV-110 product completeness, real IAM/tenanted durable owner repository, live EDR/SIEM source ingestion, host routing/TLS/audit, Shared Timeline Engine, correction/export/replay, or OPEN-013 resolution. The historical A-CONTRACT `state=preimplementation` is retained for historical evidence.

A single independent read-only L2 task has been scheduled on immutable SHA `06d1d47572b9fbea1325ad71309a895d960e1cbd`. Its report is not currently available in GitHub PR reviews/comments. Do not describe that task as approved; if a material R0/R1/R2 finding is proven, reopen this bounded work unit before any next dependent acceptance.

Reuse-first: Go stdlib `net/http/httptest`, existing first-party owner and HTTP interfaces, current GitHub CI. No vendor SDK, new index or database is justified within D. E is an evidence/handoff lot only, not production implementation. Future production-host implementation must be decomposed and separately authorized.
