---
id: settings-users-and-roles-role-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Role Management

## Objectif

Définir role management dans Platform Settings.

## Périmètre

Document canonique du domaine. Il définit l'administration du `Role` canonique et renvoie vers Security pour la signification des permissions et l'autorisation. Il ne définit ni Group, ni moteur d'assignation, ni Effective Access.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- `Role`
- `Principal` par référence de relation typée déjà sourcée
- Tenant par référence

## Fonctionnalités

- Role baselines.
- Custom constraints.
- Expiry.
- No implicit authority.

Le contrat fonctionnel détaillé est `../capabilities/cap-set-006-role-administrative-lifecycle-constraints-and-principal-relation-boundary.md`. Les états restent exactement `draft`, `active`, `deprecated`. `Expiry` est une condition/contrainte temporelle et n'est jamais un état canonique `expired`.

Une relation typée Principal/Role reste tenant-scoped et sourcée; elle ne vaut ni Permission definition, ni effective authorization, ni Decision Authority. Aucun `RoleAssignment`, `AccessAssignment` ou `PermissionAssignment` n'est créé.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.
- Réutiliser `SET-IAM-001`; aucun nouvel Screen ID n'est créé par ce lot.

## Permissions

Les modifications suivent `../../14-security-permissions-and-trust/permission-model.md`. Les permissions existantes `perm.settings.identity.*` et `perm.platform-settings.role.*` sont consommées sans nouvel identifiant et sans bulk normalization.

## États

`draft`, `active`, `deprecated` viennent de `../../05-domain-model/objects/role.md`. Une condition d'expiry n'ajoute aucun état. Les transitions invalides sont refusées côté serveur par l'owner canonique.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../../05-domain-model/objects/role.md`
- `../../05-domain-model/objects/principal.md`
- `../../14-security-permissions-and-trust/permission-model.md`
- `../../14-security-permissions-and-trust/decision-authority.md`
- `../capabilities/cap-set-006-role-administrative-lifecycle-constraints-and-principal-relation-boundary.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Role ≠ Permission, Group, job title ou Decision Authority.
- `draft`, `active`, `deprecated` restent les seuls états canoniques Role.
- Expiry reste condition/contrainte, jamais état.
- Aucune sémantique générique d'assignment, inheritance, precedence ou effective access n'est inventée.
- Aucun nouvel Screen ID ou Permission ID n'est introduit.

## Questions ouvertes

- OPEN-013 reste ouvert pour la politique par défaut des mutations Class 2; les sémantiques d'assignment absentes ne sont pas résolues localement.
