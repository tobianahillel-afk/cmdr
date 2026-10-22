---
id: endpoint-ept3-functional-permissions
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Functional Permission Needs

EPT-3 identifies permission needs without finalizing RBAC/ABAC or inventing a new atomic namespace.

Required functional checks include: Detection Content reference read; restricted content/rationale read; local evaluation/match/candidate read; sensitive detection-context read; coverage/gap read; process context and restricted command metadata read; file/path/hash context read; network identifiers/DNS context read; user/session/privilege context read; system/module/driver context read; local timeline/correlation read; Endpoint Investigation Summary/provenance read; destination handoff permission; explicit cross-tenant denial.

Existing permission families remain canonical: `perm.endpoint-agent.*` for Endpoint-owned objects, `perm.investigate.*` for Case/Evidence/Finding and Detection Engineering consumers, `perm.command.detection.*`/`perm.command.signal.*` for canonical Command objects, `perm.platform-settings.*` for Settings-owned policy/source administration, and `perm.govern.*` for response/Decision/Result. Projection never grants the source permission.

EPT-3 grants no acquisition, Live Response, containment, response execution, Evidence/Finding/Case creation or cross-tenant privilege.