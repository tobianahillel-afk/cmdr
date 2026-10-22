---
id: settings-users-and-roles-user-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# User Management

## Objectif

Définir user management dans Platform Settings.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés. Le terme UX `User` est une projection du `Principal` canonique; aucun objet `User` ou `ServiceIdentity` parallèle n'est créé.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- `Principal`
- Tenant par référence

## Fonctionnalités

- Lifecycle.
- MFA/SSO mapping.
- Suspension.
- Audit.

Le contrat fonctionnel détaillé est `../capabilities/cap-set-005-principal-administrative-lifecycle-and-identity-state.md`. Les états restent exactement ceux de `OBJ-PRINCIPAL`: `pending`, `active`, `suspended`, `revoked`. L'action d'interface `Inviter` ne crée pas un état canonique `invited`.

Human Principal et Service Principal partagent le contrat `Principal`; cette distinction ne transforme ni un Principal en Person, ni un mapping d'authentification en autorisation.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.
- Réutiliser `SET-IAM-001`; aucun nouvel Screen ID n'est créé par le lot Identity Administration.

## Permissions

Les modifications suivent le modèle défini dans `../../14-security-permissions-and-trust/permission-model.md`. `perm.settings.identity.*` et `perm.platform-settings.principal.*` restent des identifiants existants; ce document ne crée ni ne normalise en masse de Permission ID.

## États

Les états métier viennent de `../../05-domain-model/objects/principal.md`; aucune transition non sourcée n'est inventée. Les états documentaires suivent `00-governance/document-status-model.md`.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../../05-domain-model/objects/principal.md`
- `../../14-security-permissions-and-trust/permission-model.md`
- `../capabilities/cap-set-005-principal-administrative-lifecycle-and-identity-state.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Aucun objet User/ServiceIdentity parallèle n'est créé.
- `pending`, `active`, `suspended`, `revoked` restent les seuls états canoniques Principal.
- MFA/SSO mapping reste un contexte administratif et ne vaut pas autorisation.
- Aucun nouvel Screen ID ou Permission ID n'est introduit.

## Questions ouvertes

- OPEN-013 reste la décision existante pertinente pour la politique par défaut des mutations Class 2; aucune nouvelle décision n'est créée ici.
