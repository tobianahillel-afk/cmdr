---
id: studio-permission-namespace-audit
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: quality-boundary
---
# Studio Permission Namespace Audit

STD-1 confirms two historical families coexist in the current repository:
- object/register references including `perm.cmdr-studio.*`;
- Studio screen/catalog references including `perm.studio.*`.

## Disposition
- no bulk rename;
- no existing screen permission is broken;
- no final RBAC/ABAC namespace is selected;
- capabilities state functional permission needs independently of atomic identifier choice;
- no new OPEN is created because the ambiguity does not block functional STD-1 documentation;
- Product Architecture/Security may normalize later when final atomic permissions are designed.

The distinction between read, manage, invoke/execute, sensitive use, cross-tenant scope, export and authority remains mandatory regardless of final namespace spelling.
