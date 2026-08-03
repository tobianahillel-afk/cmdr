# Canonical Permission Model

## Objective

Define permission semantics once for every console and shared service.

## Model

CMDR combines:

- **RBAC** for role baselines.
- **ABAC** for tenant, environment, service ownership, sensitivity, legal hold and risk.
- **Decision authority** for approval thresholds independent from ordinary CRUD access.
- **Separation of duties** preventing a requester from self-approving where policy forbids it.

## Permission families

- `command.read`, `command.coordinate`, `incident.manage`
- `case.create`, `case.read`, `case.manage`
- `evidence.collect`, `evidence.read`, `evidence.verify`, `evidence.export`
- `analysis.execute`, `sandbox.execute`, `reverse.execute`, `debugger.execute`
- `finding.author`, `finding.approve`
- `response.request`, `response.approve`, `response.execute`, `response.rollback`
- `playbook.manage`, `policy.manage`, `authority.manage`
- `agent.manage`, `integration.manage`, `tenant.manage`
- `audit.read`, `audit.export`, `report.manage`
- `platform.admin`

## Enforcement rules

- Every operation is tenant-scoped before object lookup.
- Read permission does not imply export permission.
- Analysis execution may require isolation, sample-sensitivity and tool-specific policy.
- Response approval is evaluated against action type, target criticality, business service, impact and emergency mode.
- Privileged administrative access does not automatically grant response authority.
- Denials and attempted cross-scope access are audited.
- Frontend visibility is not an enforcement boundary; backend authorization is mandatory.

## UX contract

Unavailable actions are either hidden when irrelevant or disabled with a clear reason when users benefit from knowing the action exists. Approval screens explain required authority and missing conditions without exposing sensitive policy internals.
