---
id: CMP-INSPECTOR
domain: 03-design-system
status: draft
owner: Design Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Inspector

## Objectif

Fournir la source unique du panneau d’inspection rapide utilisé dans tous les produits.

## Anatomie

1. En-tête: type, identifiant, titre, état et fermeture.
2. Résumé décisionnel.
3. Propriétés essentielles.
4. Relations et objets liés.
5. Provenance et fraîcheur.
6. Actions autorisées.
7. Lien `Ouvrir la page complète`.

## Comportement

- S’ouvre sans réinitialiser la liste ou le canvas source.
- Conserve et restaure le focus.
- Peut être épinglé si le layout le permet.
- N’édite pas un objet si l’écran n’a pas la permission et le contrat d’édition.
- Ne redéfinit pas les champs des objets; il les consomme depuis `05-domain-model/objects/`.
- Supporte Loading, Empty, Partial, Error, Offline et Permission denied.
- En responsive, devient un drawer ou une page temporaire.

## Critères d’acceptation

- Même anatomie dans les six produits.
- Lien profond et retour source fonctionnels.
- Aucune donnée hors permission.
- Navigation clavier complète.
