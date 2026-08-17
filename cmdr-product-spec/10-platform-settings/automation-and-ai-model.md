---
id: platform-settings-automation-and-ai-model
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
---
# Platform Settings Automation and AI Model

AI is optional and cannot become an administrative authority. Every Settings capability has a deterministic/manual path.

## Tenant and Environment Foundations — allowed assistance

- explain Tenant/Environment state and sourced transition constraints;
- summarize non-secret administrative provenance;
- explain deterministic validation failures;
- highlight incompatible Tenant/Environment context;
- suggest a next step already permitted by current sources and permissions.

## Identity Administration — allowed assistance

- explain sourced Principal state;
- explain sourced Role state and Role constraints;
- summarize Access Review evidence/provenance without fabricating missing facts;
- explain deterministic denial;
- identify a potential SoD issue for human review when the Security rule is available;
- suggest a safe next administrative step already allowed by sources and current permissions.

## Sources & Parsers — allowed assistance

Bounded assistance may:
- explain Data Source or Parser configuration and sourced lifecycle constraints;
- summarize non-secret metadata, version, freshness, error/quality metadata and audit history;
- explain a fixture or deterministic validation error;
- suggest a non-authoritative mapping or next administrative step already permitted by current sources and permissions;
- explain that a projected health/test result is stale, partial or sourced elsewhere.

The deterministic/manual path remains authoritative. AI is never required to configure, inspect, validate, version or safely disable a Data Source or Parser.

## Forbidden authority

AI must not create or mutate Tenant/Environment autonomously; create, suspend or revoke Principal autonomously; create or mutate Role autonomously; create a Principal/Role relation autonomously; grant or revoke Permission; execute an Access Review disposition autonomously; invent Group membership or effective access; override authorization; bypass tenant isolation, step-up or SoD; fabricate provenance; create Decision Authority; resolve OPEN decisions; or hide deterministic failures.

For Sources & Parsers, AI must additionally not:
- enable ingestion or collection;
- execute a connector, external source test, parser engine, plugin, sandbox, normalization engine or stream processor;
- fabricate health, freshness, successful validation, parsed output, schema fields or runtime support;
- invent a Data Source→Parser assignment, compatibility, route, selection, fallback or precedence;
- reveal, copy or place Secret values in prompts, outputs, logs, fixtures or examples;
- silently activate/disable production configuration or retire/switch Parser versions;
- select ECS, OCSF, CIM, OpenTelemetry or another schema standard absent a canonical decision.

For CAP-SET-007, AI may summarize evidence or suggest a review disposition, but a `revoke` conclusion remains a human/authorized administrative disposition. Because the current source corpus does not define generic assignment-removal mechanics, AI cannot turn that disposition into a direct relation mutation.

## Product boundary

Platform Settings administers platform configuration, including Data Source/Parser administration. Studio owns agentic automation semantics and runtime Tool/Tool Call/Automation Run behavior; Security owns authorization constraints; Govern owns Decision/Approval/authority semantics. Investigate/Collection/Endpoint and other source owners retain their runtime responsibilities. AI assistance transfers none of those ownerships.

## Deterministic/manual alternative

The absence of an AI provider must not block Settings capability operation. State inspection, transition validation, permission denial, freshness inspection, configuration checks, fixture inspection, administrative version selection and handoff/result projection all retain deterministic/manual paths. Runtime execution remains outside Settings unless separately sourced.

## Hard scope gate

Settings capability specification creates no Permission ID and no Screen ID. If safe execution of any AI-assisted or deterministic path would require one, the relevant lot is BLOCKED and that prerequisite is handled in a separate source-owned run.
