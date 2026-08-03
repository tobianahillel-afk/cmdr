---
id: investigate-disk-and-artifact-forensics
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Disk & Artifact Forensics

## Objectif

Analyser Disk Image et artefacts endpoint en lecture seule.

## Périmètre

Module du produit 07-investigate. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Investigate Product Lead.

## Objets concernés

- disk-image
- artifact
- evidence
- timeline-entry

## Fonctionnalités

- Filesystem.
- Registry/log/browser.
- Carving.
- Super timeline.

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
