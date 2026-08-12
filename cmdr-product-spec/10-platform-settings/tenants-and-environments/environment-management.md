---
id: settings-tenants-and-environments-environment-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Environment Management

## Objectif

Définir environment management dans Platform Settings.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- Environment canonique (`../../05-domain-model/objects/environment.md`).
- Tenant comme scope obligatoire, sans redéfinition.

## Fonctionnalités

- Production/Preproduction/Test/Development/Lab.
- Isolation.
- Promotion context.
- Labels.

`CAP-SET-002` formalise lifecycle/type/Tenant scope. Promotion context reste du contexte et ne signifie pas promotion automatique. Cardinalité Tenant–Environment, nesting, Environment partagé cross-tenant et migration inter-tenant ne sont pas inventés sans source canonique.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../../14-security-permissions-and-trust/permission-model.md`. Le lot réutilise `perm.platform-settings.environment.read/manage`; aucun nouvel ID n'est autorisé.

## États

Les états métier viennent de `../../05-domain-model/objects/environment.md`; aucun lifecycle concurrent n'est défini ici.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../capabilities/cap-set-002-environment-administrative-lifecycle-and-tenant-scope.md`
- `../capabilities/cap-set-004-cross-product-tenant-environment-context-preservation.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- Aucun nesting, cardinalité ou mécanisme de promotion non sourcé n'est ajouté.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.
