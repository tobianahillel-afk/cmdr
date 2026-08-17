---
id: source-native-capability-strategy
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# Native Capability Strategy

## Classification

Every capability has exactly one delivery classification:

- `native`
- `integrated`
- `temporary-integration`
- `planned`
- `out-of-scope`

An integration is not described as native. A planned capability is not described as implemented. UX is designed around user activity rather than the integrated vendor.

## Required metadata

Every capability records owner, user outcome, delivery classification, rationale, dependencies, limits, roadmap phase and replacement plan when temporary.

## Long-term native targets

Search, investigation, collection, Live Response, endpoint telemetry, endpoint detection, forensic, static analysis, reverse engineering, debugger, sandbox, Detection Engineering, orchestration, governance, reporting, audit and automation are long-term native targets. Their current delivery classification remains subject to explicit product decisions.

## Phase boundary

Phase 0 records the strategy. Detailed engine, protocol, storage, PKI and framework decisions are deferred to later technical phases.