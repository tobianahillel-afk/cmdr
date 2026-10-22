---
id: settings-tenants-and-environments
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Tenants & Environments

## Objectif

Administrer frontières d’isolation et environnements.

## Périmètre

Module du produit 10-platform-settings. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- tenant
- environment

## Fonctionnalités

- Lifecycle.
- Residency.
- Environment labels.
- MSSP scope.

## Capability contracts

- `../capabilities/cap-set-001-tenant-administrative-lifecycle-and-isolation-boundary.md`
- `../capabilities/cap-set-002-environment-administrative-lifecycle-and-tenant-scope.md`
- `../capabilities/cap-set-003-tenant-environment-administrative-change-validation-and-provenance.md`
- `../capabilities/cap-set-004-cross-product-tenant-environment-context-preservation.md`

Les quatre contracts sont possédés uniquement par Platform Settings. Tenant–Environment relationship, validation, audit handoff et context projection sont couverts comme responsabilités de ces contracts et ne créent aucun objet ou owner concurrent.

## UX et interactions

- Conserver le contexte de liste, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Préserver navigation clavier et liens profonds.
- Réutiliser `SET-TEN-001`; aucun nouvel Screen ID dans ce lot.

## Permissions

Voir `../../14-security-permissions-and-trust/permission-model.md` et le registre des permissions. Réutiliser les permissions Tenant/Environment existantes; une nouvelle Permission ID nécessaire rend le lot BLOCKED.

## États

Les états métier viennent des fichiers d’objets canoniques; la page ajoute uniquement Loading, Empty, Partial, Error, Offline et Permission denied.

## Dépendances

- 03-design-system/
- 04-experience-architecture/
- 05-domain-model/
- 17-implementation-contracts/

## Critères d’acceptation

- Aucune définition d’objet ou de permission locale.
- Tous les écrans du module ont un front matter et 27 sections.
- Les transitions sont auditées et idempotentes.
- Les cardinalités, nesting et migrations inter-tenant non sourcés ne sont pas inventés.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
