---
id: settings-users-and-roles-access-reviews
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Access Reviews

## Objectif

Définir access reviews dans Platform Settings.

## Périmètre

Document canonique du domaine. Il couvre la revue périodique de l'accès visible à partir des références Principal/Role sourcées, son owner, son evidence/provenance et sa disposition. Il ne crée ni objet `AccessReview`, ni moteur d'assignation, ni moteur de conformité.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- `Principal` par référence
- `Role` par référence
- provenance/audit administratif par référence

## Fonctionnalités

- Periodic review.
- Owner.
- Evidence.
- Revocation.

Le contrat fonctionnel détaillé est `../capabilities/cap-set-007-access-review-evidence-and-revocation-disposition.md`.

`Revocation` signifie ici une disposition keep/revoke sourcée. Le corpus canonique courant ne définit pas de mécanique générique de retrait d'un `RoleAssignment` ou `AccessAssignment`; en conséquence une conclusion `revoke` produit une **revocation disposition / revocation-required handoff** et ne supprime aucune relation implicitement.

Review evidence n'est pas automatiquement l'objet Evidence d'Investigate. Access Review n'est ni Govern Approval, ni Govern Decision, ni generic Compliance engine.

## UX et interactions

- Réutiliser `SET-IAM-001` et son action `Lancer une revue`.
- Conserver Tenant, owner, scope, versions, références et gaps visibles.
- Réutiliser `SET-AUD-001` pour la projection de provenance administrative.
- Aucun nouvel Screen ID n'est créé.

## Permissions

Les modifications suivent `../../14-security-permissions-and-trust/permission-model.md`. La surface consomme `perm.settings.identity.read/manage` et les lectures objet Principal/Role applicables. Aucun `review.execute`, `role.assign` ni autre Permission ID n'est créé.

## États

Aucune state machine canonique Access Review n'est définie par la source et ce document n'en invente aucune. Les états Principal/Role restent dans leurs objets canoniques; la review expose seulement progression/complétude documentaire, gaps et disposition sourcée.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../../05-domain-model/objects/principal.md`
- `../../05-domain-model/objects/role.md`
- `../../14-security-permissions-and-trust/permission-model.md`
- `../capabilities/cap-set-007-access-review-evidence-and-revocation-disposition.md`
- `../administrative-audit/admin-event-catalog.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Owner, scope, evidence/provenance et disposition sont auditables.
- Une disposition `revoke` ne crée pas de mécanique d'assignation absente.
- Aucun objet AccessReview/Assignment, Group, Permission ID ou Screen ID n'est créé.
- Les décisions et approvals Govern restent distincts.

## Questions ouvertes

- OPEN-007/013/014/015/019 restent ouvertes selon leur domaine; aucune n'est résolue pour compléter artificiellement une review.
