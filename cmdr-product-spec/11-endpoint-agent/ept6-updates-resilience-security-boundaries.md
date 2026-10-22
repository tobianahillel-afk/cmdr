---
id: endpoint-ept6-updates-resilience-security-boundaries
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008, OPEN-013, OPEN-015, OPEN-017]
---
# EPT-6 — Updates, Resilience, Security and Endpoint Provenance Boundaries

## Purpose
Define ownership and non-overlap boundaries for the final Endpoint Capability Specification execution lot before individual `CAP-EPT-082..099` contracts are interpreted.

## Update boundary
Platform Settings owns administrative desired version, Fleet targeting, upgrade channels/waves, policy assignment and administrative package/provider configuration. Endpoint owns the local effective assignment projection and local technical lifecycle only.

Mandatory separations:
- administrative target version != current local version;
- assigned != downloaded != staged != ready != installed != activated != healthy;
- compatible != officially supported;
- update deferred != failed;
- retry != recovery;
- Endpoint Update Reversion != Govern Response Rollback;
- Endpoint Update Reversion != Studio Deployment Reversion;
- Settings upgrade wave != Endpoint local execution.

## Resilience boundary
Endpoint owns local buffer/queue state, retained-item facts, reconnect/replay/resume facts, restart/crash recovery facts, degraded-operation/resource-pressure facts and local dependency/capability reassessment. Shared retains generic Jobs, Retry, Recovery, Trace and Activity mechanisms.

Mandatory separations:
- buffered != delivered;
- retained != durable forever;
- replayed != consumed;
- replay != exactly-once;
- reconnect != recovered;
- Agent restarted != healthy;
- local state restored != remote state synchronized;
- degraded operation != complete failure;
- backpressure != durability guarantee;
- local recovery != Govern rollback.

## Security boundary
Endpoint owns local self-protection state, tamper observations, runtime privilege/security context, protected-component state, local sensitive-material handling facts, local audit events and technical security-state consumer handoff. Security Architecture owns global permission/privacy/audit-integrity policy. Platform Settings owns credentials and Secret administration.

Mandatory separations:
- self-protection != complete endpoint security;
- anti-tamper state != cryptographic integrity proof;
- integrity check PASS != absence of compromise;
- protected service running != uncompromised Agent;
- privilege != permission;
- local audit event != Shared Trace;
- append-only != cryptographic proof;
- tamper candidate != malicious verdict;
- Secret Reference != raw secret;
- security-state observation != Incident/Finding/Result.

## Delivery boundary
Every EPT-6 capability remains `draft / defined / planned`. Endpoint Capability Specification PASS, if later justified by closure gates, means documentary functional completeness only; it does not mean implemented, production-ready, deployed, platform-supported or globally complete.

## Screens and implementation
Endpoint Screen IDs remain 0. EPT-6 defines no UI layout, API, protocol, crypto scheme, package format, storage engine, physical queue, watchdog, driver, service implementation, final RBAC/ABAC or Phase 6 work.

## Acceptance
**Given** a Settings upgrade wave targets an Agent, **When** Endpoint reports update state, **Then** the administrative assignment and the local execution state remain separate owner-typed references.

**Given** an offline queue later replays data, **When** replay is reconciled, **Then** replay does not imply exactly-once consumption.

**Given** Endpoint observes tamper-related state, **When** it is handed to another product, **Then** the observation remains a technical fact and is not automatically converted into Incident, Finding or Result.