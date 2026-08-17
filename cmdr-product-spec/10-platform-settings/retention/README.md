---
id: settings-retention
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Retention

## Objectif

Définir rétention, legal hold, export et disposition.

## Périmètre

Module du produit 10-platform-settings. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- policy
- tenant

## Fonctionnalités

- Object matrix.
- Legal hold.
- Deletion evidence.
- Residency.

## UX et interactions

- Conserver le contexte de liste, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Préserver navigation clavier et liens profonds.

## Permissions

Voir `../../14-security-permissions-and-trust/permission-model.md` et le registre des permissions.

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

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
