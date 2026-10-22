---
id: reporting-engine
domain: 12-shared-capabilities
status: draft
owner: Shared Capabilities Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Reporting Engine

## Objectif

Être la source unique du concept **Report**, de sa composition, son snapshot, sa revue, sa publication, sa distribution et son export.

## Modèle canonique

Un Report référence un template versionné, un tenant, un auteur, une audience, un snapshot de données, des citations d’objets, une politique de redaction, un statut et une version publiée immuable.

## Fonctionnalités

- modèles opérationnels, investigation, gouvernance, exécutifs et client ;
- sections et citations stables ;
- aperçu identique à l’export ;
- snapshots reproductibles ;
- revue, approbation, publication et supersession ;
- redaction et classification ;
- exports PDF/HTML/JSON à décider ;
- planification et distribution ;
- historique et audit.

## Permissions

`perm.shared.report.read`, `perm.shared.report.create`, `perm.shared.report.review`, `perm.shared.report.publish`, `perm.shared.report.export`, avec contrôle objet et tenant.

## États

`draft → review → approved → published → superseded | archived`.

## Invariants

- Un Report publié est immuable.
- Toute valeur chiffrée référence une définition de métrique et un snapshot.
- L’export ne peut élargir la visibilité.
- Les produits Command, Investigate et Govern consomment ce moteur; ils ne redéfinissent pas Report.

## Critères d’acceptation

- Les citations restent résolubles.
- Preview et export concordent.
- Redaction et permissions sont testées.
- Version, template et snapshot sont audités.
