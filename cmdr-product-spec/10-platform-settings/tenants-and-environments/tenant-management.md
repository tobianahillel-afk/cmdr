---
id: settings-tenants-and-environments-tenant-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Tenant Management

## Objectif

Définir tenant management dans Platform Settings.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- Tenant canonique (`../../05-domain-model/objects/tenant.md`).

## Fonctionnalités

- Provisioning.
- Suspension.
- Offboarding.
- Residency.

Le lifecycle détaillé reste celui de l'objet Tenant canonique. `CAP-SET-001` formalise l'administration de ce lifecycle et sa frontière d'isolation; `CAP-SET-003` formalise validation/provenance. Suspension n'est pas suppression et offboarding n'est pas effacement irréversible.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../../14-security-permissions-and-trust/permission-model.md`. Le lot réutilise `perm.platform-settings.tenant.read/manage`; aucun nouvel ID de permission n'est autorisé.

## États

Les états métier viennent de `../../05-domain-model/objects/tenant.md`; aucun lifecycle concurrent n'est défini ici.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../capabilities/cap-set-001-tenant-administrative-lifecycle-and-isolation-boundary.md`
- `../capabilities/cap-set-003-tenant-environment-administrative-change-validation-and-provenance.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- Aucun Customer lifecycle, fallback tenant ou suppression physique n'est inventé.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.
