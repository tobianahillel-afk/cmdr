---
id: shared-saved-views
domain: 12-shared-capabilities
status: draft
owner: Shared Capabilities Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-011
  - REQ-UX-003
  - REQ-UX-009
---

# Saved Views génériques

## Propriété

Shared Capabilities possède le stockage, le versioning, le partage et l'application permission-aware d'une configuration de vue. Le produit consommateur possède les vues système et la signification de leurs filtres.

## Configuration autorisée

- route/page/workspace propriétaire ;
- filtres et opérateurs ;
- tri et groupement ;
- colonnes, ordre, largeur et pinning ;
- densité ;
- mode ;
- périmètre personnel ou partagé ;
- version de schéma.

Une Saved View n'enregistre jamais les objets, résultats, permissions ou secrets.

## Ouverture

1. vérifier tenant et produit ;
2. charger la version ;
3. réévaluer permissions ;
4. supprimer/masquer les champs interdits avec explication ;
5. migrer le schéma ou signaler l'incompatibilité ;
6. appliquer la configuration et marquer l'état dirty si modifiée.

## Portées

System : définie et versionnée par le produit. Personal : owner unique. Shared : permission de partage, audience et audit. Une vue partagée reste immutable pour le destinataire jusqu'à duplication.

## Relation Work Queue

Command définit exclusivement `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `My Work` dans `../06-command/modules/incidents-and-work-queue/saved-views.md`. Le composant commun est `../03-design-system/components/saved-views.md`.

## Critère d’acceptation

**Given** une vue partagée avec une colonne devenue interdite,  
**When** un destinataire l'ouvre,  
**Then** la permission est réévaluée, la colonne est retirée, la vue source reste inchangée et aucune donnée n'est divulguée.
