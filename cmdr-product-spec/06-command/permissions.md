---
id: 06-command-permissions
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-PROD-013
open_decisions:
  - OPEN-013
---
# Functional permission needs — Command

## Modèle référencé

Command consomme RBAC, ABAC tenant/environnement/ownership/classification, Decision Authority, Separation of Duties et step-up selon `../14-security-permissions-and-trust/permission-model.md`.

## Familles existantes

- `perm.command.read`
- `perm.command.coordinate`
- `perm.command.incident.read`
- `perm.command.incident.manage`
- `perm.command.alert.read/manage`
- `perm.command.signal.read/manage`
- `perm.command.task.manage`
- `perm.command.saved-view.share`
- permissions Investigate/Govern/Shared lorsqu’une transition ou projection l’exige.

## Besoins fonctionnels non finalisés

| Besoin | Capabilities | Lacune |
|---|---|---|
| lecture Task distincte | CAP-CMD-101/107 | `perm.command.task.read` absente ou non normalisée |
| assignment/reassignment | CAP-CMD-102/103 | granularité manage trop large |
| priorité/impact | CAP-CMD-002/104/204 | permission atomique et owner métier non décidés |
| pause SLA | CAP-CMD-105 | policy + permission + OPEN-013 |
| bulk actions | CAP-CMD-108 | permission batch et limites non définies |
| handover send/ack | CAP-CMD-004 | permissions par rôle non définies |
| readiness/exercise/plan | CAP-CMD-301..305 | namespaces spécifiques absents |
| customer/contract scope | CAP-CMD-401 | dépend de OPEN-006 |

## Règles

- lire n’implique ni exporter, ni modifier, ni lancer un workflow ;
- la permission d’un lien n’accorde pas celle de l’objet destination ;
- les actions de classe 2 restent soumises à OPEN-013 ;
- Command ne reçoit jamais `decision.approve`, `response-run.execute` ou permission de containment ;
- la matrice atomique et les namespaces sont reportés à la phase permissions.

## Critère

**Given** un utilisateur pouvant lire un Incident mais non le gérer, **When** il ouvre Priority Management, **Then** les facteurs sont visibles, la mutation est refusée avec raison, et aucune donnée d’un autre tenant n’est révélée.
