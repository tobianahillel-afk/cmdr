---
id: CAP-STD-001
title: Studio Library and Asset Catalog
product: cmdr-studio
module: library
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-001 — Studio Library and Asset Catalog
## 1. Définition
Catalogue Studio permission-aware des assets et références réutilisables. Library n'est ni un package repository technique ni le moteur Search partagé.
## 2. Problème utilisateur
Un utilisateur doit découvrir un asset Studio sans confondre visibilité, disponibilité, compatibilité, autorisation et exécutabilité.
## 3. Objectifs
Lister, rechercher, filtrer, grouper et ouvrir les assets autorisés avec owner, type, version, lifecycle, dépendances, consumers et provenance.
## 4. Non-objectifs
Ne pas définir Workflow/Agent au-delà de leurs références, ne pas implémenter Search, ne pas créer de runtime ni de schéma physique.
## 5. Propriétaire
CMDR Studio / Library possède la sémantique du catalogue. Shared conserve Search; chaque asset conserve son owner canonique.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer, SOC Analyst et consommateurs autorisés secondaires.
## 7. Conditions d’entrée
Tenant/environnement résolus, Principal authentifié, contexte d'autorisation courant et projections Studio disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| requête de catalogue | utilisateur | recherche/filtre | oui | requête courante | catalogue autorisé non filtré |
| métadonnées assets | sources Studio | projections typées | oui | freshness visible | Partial |
| autorisation | Security | scope permission | oui | courant | champs/actions restreints |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | id, owner, version, status | read/link |
| Version | CMDR Studio | version, compatibility, supersession | read |
| Workflow/Agent reference | CMDR Studio | référence seulement | read/link |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Library catalog projection | composer | CMDR Studio | aucun nouveau lifecycle métier |
## 11. Fonctionnalités
Search/filter via Shared, grouping par type/owner/status, inspection des dépendances/consumers/provenance et ouverture de la surface propriétaire.
## 12. Actions utilisateur
Class 0: browse, search, filter, inspect, open. Aucune action d'exécution implicite.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| discovery | oui | oui | oui | oui | filtres + Shared Search |
| expliquer compatibilité | oui | oui | oui | oui | champs déterministes |
## 14. États fonctionnels
ready, partial, stale, permission-filtered; lifecycle de l'asset source inchangé.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale exposent la cause sans inventer de données.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| catalogue filtré | projection | utilisateur | metadata autorisée seulement |
| lien propriétaire | référence | module Studio | owner et return-origin préservés |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library | open asset | owning surface | asset/version/return-origin | même état Library |
| produit consommateur | open Studio ref | Library | tenant/env/ref | produit source |
## 18. Dépendances
Shared Search/Linking, Security, sources Studio et Context preservation. Aucune dépendance n'accorde ownership.
## 19. Source de vérité
Les objets restent dans leurs sources Studio; le catalogue n'est qu'une projection composée.
## 20. Provenance et audit
Conserver requête, filtres, actor, tenant/env, source, version, freshness et navigation sans secret brut.
## 21. Permissions fonctionnelles
Library read, restricted metadata read et owning-surface access. `perm.studio.*` et `perm.cmdr-studio.*` restent toutes deux historiques/non finalisées.
## 22. Limites et erreurs
Source absente, asset inaccessible, version inconnue, Search indisponible ou permission denied donnent Partial/refus explicite.
## 23. Métriques
Catalog completeness, stale sources, permission-filtered entries, broken owner links; aucune cible KPI/SLO.
## 24. Classification de livraison
`defined / planned`; aucune disponibilité runtime n'est revendiquée.
## 25. Critères d’acceptation
**Given** asset visible mais invoke interdit, **When** Library l'affiche, **Then** metadata permise visible et invoke bloqué.  
**Given** Shared Search indisponible, **When** Library s'ouvre, **Then** browse/filters sûrs restent possibles et Partial est visible.  
**Given** IA désactivée, **When** discovery est utilisée, **Then** recherche, filtres, inspection et liens restent fonctionnels.
## 26. Questions ouvertes
Aucune nouvelle OPEN. OPEN-003 reste présentation seulement; l'anomalie de namespace permission n'est pas tranchée ici.
## 27. Consommateurs documentaires
Library, Skills, futurs STD-2/3/4, Command/Investigate autorisés, Security, Quality et Roadmap.
