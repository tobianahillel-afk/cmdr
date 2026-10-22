---
id: personalization
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-007
  - REQ-UX-003
  - REQ-UX-009
---

# Personnalisation

## Autorisée

Thème, densité, ordre et largeur de colonnes, vues personnelles, favoris, raccourcis non conflictuels, ordre de modules locaux, taille de panneaux dans les limites, préférences de notification et sections non critiques de l'Inspector.

## Interdite

Sémantique, permissions, provenance, statuts, audit, actions obligatoires, step-up, Human Gates, classes de risque, libellés d'autorité et données sources.

## Portée

Les préférences sont personnelles et tenant-scoped lorsque nécessaire. Une vue partagée ne remplace pas les permissions du destinataire. Un reset restaure les défauts d'activité sans effacer le travail métier.

## Critère

**Given** une vue partagée contenant une colonne non autorisée,  
**When** un utilisateur l'ouvre,  
**Then** la colonne est retirée avec explication, la vue source n'est pas modifiée et aucune donnée n'est divulguée.
