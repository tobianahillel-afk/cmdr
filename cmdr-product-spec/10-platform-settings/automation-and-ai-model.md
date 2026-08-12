---
id: platform-settings-automation-and-ai-model
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Platform Settings Automation and AI Model — Tenant and Environment Foundations

AI is optional and cannot become an administrative authority. Every first-lot capability has a deterministic/manual path.

## Allowed assistance

- explain Tenant/Environment state and sourced transition constraints;
- summarize non-secret administrative provenance;
- explain deterministic validation failures;
- highlight incompatible Tenant/Environment context;
- suggest a next step already permitted by current sources and permissions.

## Forbidden authority

AI must not create or mutate Tenant/Environment autonomously, grant permission, invent configuration values or relationships, override Security, bypass tenant isolation, fabricate provenance, create Govern authority, resolve OPEN-013, or hide a deterministic failure.

## Product boundary

Platform Settings administers model providers/policies; Studio owns agentic automation semantics. The absence of an AI provider must not block CAP-SET-001..004.

## Hard scope gate

This lot creates no Permission ID and no Screen ID. If safe execution of any AI-assisted or deterministic path would require one, the lot is BLOCKED and that prerequisite is handled in a separate run.
