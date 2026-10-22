---
id: endpoint-ept6-cross-product-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# EPT-6 Cross-Product Links

## Platform Settings
Consumes administrative target version, Fleet, update channel/wave, Endpoint Policy and Secret References. Settings keeps administration ownership; Endpoint returns local observed execution/health/recovery facts.

## CMDR Studio
Studio owns Tool/Skill/Workflow/Agent asset release, deployment and deployment reversion. Endpoint Agent update lifecycle is independent: Agent binary/runtime update != Studio asset deployment.

## Govern
Govern owns Decision, Response Run, response verification, Response Rollback and Result. Endpoint update reversion/local recovery may be referenced by Govern but never completes a Govern rollback or creates Result.

## Shared Capabilities
Shared owns generic Jobs, Retry, Recovery, Trace, Activity, Reporting and generic mechanisms. Endpoint owns only local queue/replay/restart/dependency/security facts.

## Security
Security owns Permission Model, privacy/minimization, global audit/integrity requirements, tenant isolation and secret policy. Endpoint exposes local security observations and local audit events under those controls.

## Investigate / Command
Endpoint may hand off telemetry, detection/investigation context and security-state/tamper candidates. These remain technical facts; Investigate decides Findings/Evidence and Command owns Incident.

## Canonical distinctions
- Settings desired version != Endpoint current version.
- Endpoint Update != Studio Deployment.
- Endpoint Update Reversion != Studio Deployment Reversion.
- Endpoint local recovery != Govern Response Rollback.
- Endpoint Local Audit Event != Shared Trace != Govern Audit Trail.
- Endpoint security-state observation != Finding != Incident != Result.