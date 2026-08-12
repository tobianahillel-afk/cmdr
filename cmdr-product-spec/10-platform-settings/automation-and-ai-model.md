---
id: platform-settings-automation-and-ai-model
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
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

## Forbidden authority

AI must not create or mutate Tenant/Environment autonomously; create, suspend or revoke Principal autonomously; create or mutate Role autonomously; create a Principal/Role relation autonomously; grant or revoke Permission; execute an Access Review disposition autonomously; invent Group membership or effective access; override authorization; bypass tenant isolation, step-up or SoD; fabricate provenance; create Decision Authority; resolve OPEN decisions; or hide deterministic failures.

For CAP-SET-007, AI may summarize evidence or suggest a review disposition, but a `revoke` conclusion remains a human/authorized administrative disposition. Because the current source corpus does not define generic assignment-removal mechanics, AI cannot turn that disposition into a direct relation mutation.

## Product boundary

Platform Settings administers model providers/policies; Studio owns agentic automation semantics; Security owns authorization constraints; Govern owns Decision/Approval/authority semantics. AI assistance transfers none of those ownerships.

## Deterministic/manual alternative

The absence of an AI provider must not block CAP-SET-001..007. State inspection, transition validation, permission denial, Role expiry-condition evaluation, Access Review evidence inspection and disposition/handoff all retain a deterministic/manual path.

## Hard scope gate

Settings capability specification creates no Permission ID and no Screen ID. If safe execution of any AI-assisted or deterministic path would require one, the relevant lot is BLOCKED and that prerequisite is handled in a separate source-owned run.
