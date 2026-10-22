---
id: endpoint-ept1-functional-permissions
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-1 Functional Permission Needs

EPT-1 does not finalize RBAC/ABAC and creates no arbitrary atomic permission. It reuses the existing families `perm.endpoint-agent.endpoint-agent.read/manage`, `perm.endpoint-agent.agent-command.read/manage` and `perm.endpoint-agent.local-audit-event.read/manage`; Agent Command execution itself is out of EPT-1 scope.

Functional needs identified for future permission design include ordinary/restricted Agent read, identity/registration/enrollment-state read, bounded local registration or refresh action need, platform/version/inventory/health/heartbeat/capability/effective-policy-state read, sensitive inventory read, local-audit/provenance read, provenance-export preparation and explicit cross-tenant denial.

Platform Settings retains `endpoint-agent-fleet.read/manage` and `endpoint-policy.read/manage`. A projection or advertised capability never grants source permission or execution authority.