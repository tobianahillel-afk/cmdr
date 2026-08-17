---
id: CAP-INV-002
title: Event Search
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-010
  - REQ-PROD-011
  - REQ-PROD-014
  - REQ-PROD-045
open_decisions:
  - none
source-of-truth: canonical
---

# CAP-INV-002 — Event Search

## 1. Définition

Permet de construire, valider, exécuter, arrêter, relancer et inspecter une recherche tenant-scoped et time-scoped sur la télémétrie autorisée, sans fixer le langage, l’index ou le moteur technique.

## 2. Problème utilisateur

Le SOC Analyst ou Threat Hunter doit retrouver des Events pertinents dans plusieurs sources tout en comprenant fraîcheur, limites, coût et erreurs partielles. Une recherche opaque ou dépendante de l’IA rend le résultat non reproductible.

## 3. Objectifs

- fournir une voie manuelle et déterministe complète
- rendre période, sources, champs et erreurs visibles
- préserver requête, paramètres et provenance
- lier des résultats sélectionnés à un Case sans les transformer en Evidence

## 4. Non-objectifs

- ne pas décider la syntaxe finale du langage
- ne pas choisir moteur, index ou stockage
- ne pas exécuter automatiquement une requête sensible proposée par IA
- ne pas faire de chaque résultat un objet propriétaire

## 5. Propriétaire

Investigate / Signals And Hunt / Investigate Product Lead. Investigate possède uniquement les objets et mutations explicitement listés en section 10 ; les objets consommés restent chez leur propriétaire canonique.

## 6. Utilisateurs

Principal : SOC Analyst L2 et Threat Hunter. Secondaires : Senior Analyst, Forensic Analyst et Detection Engineer en consommation future.

## 7. Conditions d’entrée

Tenant, environnement, période et sources autorisées ; schéma ou champs disponibles ; permission de recherche et éventuel accès raw distinct.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Query draft | utilisateur, Query Asset ou proposition | expression conceptuelle et paramètres | oui | version du brouillon | éditeur vide avec aide déterministe |
| Search scope | Context Bar | tenant, environnement, période et sources | oui | confirmé avant exécution | exécution bloquée |
| Field and source metadata | Query Engine / data sources | schéma, disponibilité et fraîcheur | oui | réévaluée à l’exécution | validation partielle et source nommée |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Telemetry Event | Shared Capabilities | raw ou normalized selon permission | rechercher, filtrer, agréger et inspecter |
| Query | Shared Capabilities | texte, paramètres et version | consulter et exécuter |
| Search Job | Shared Capabilities | état, erreurs, métriques et résultats | démarrer, arrêter et relancer selon permission |
| Case | Investigate | scope et liens | relier une sélection |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Search Job | créer, arrêter ou relancer | Shared Query Engine | classe 0 pour lecture ; exécution explicitement déclenchée et auditée |
| Query execution record | enregistrer paramètres, période, sources et résultat | Shared Trace / Investigate business events | append-only conceptuel |
| Case result links | lier des résultats sélectionnés | Investigate | classe 2 ; candidats uniquement, aucune Evidence automatique |

## 11. Fonctionnalités

- éditeur ou builder conceptuel
- validation avant exécution
- exécution, annulation et relance
- résultats tabulaires et agrégés
- raw/rendered selon permission
- erreurs partielles et limites
- pivots manuels
- liaison au Case
- enregistrement via Saved Searches

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Valider | analyste | Query | 0 | scope et champs connus | diagnostics déterministes | non |
| Exécuter | analyste | Search Job | 0 | requête valide et permission | job traçable | non |
| Arrêter | analyste | Search Job | 0 | job actif | annulation demandée | non |
| Relancer | analyste | Search Job | 0 | version et paramètres visibles | nouveau run | non |
| Lier une sélection | analyste | Case relation | 2 | Case accessible | liens sourcés | OPEN-013 |

Une action dont l’effet cible relève des classes 3 ou 4 reste une préparation ou une demande dans Investigate. L’autorité et l’exécution demeurent chez Govern ou le mécanisme propriétaire.

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Valider la requête | oui | parse/validation | oui | explication facultative | diagnostics et documentation |
| Exécuter la recherche | oui | oui | workflow répété autorisé | jamais sans déclencheur/contrat | action humaine ou workflow déterministe |
| Résumer les résultats | oui | agrégations | oui | résumé attribué | table, agrégations et filtres |
| Proposer des pivots | oui | mapping de champs | oui | suggestion | pivots manuels |

Toute proposition automatisée expose initiateur, producteur, version ou run, sources, facteurs, éventuels Tool Calls, incertitude et disposition. Elle ne devient pas silencieusement un état effectif.

## 14. États fonctionnels

`draft`, `valid`, `invalid`, `queued`, `running`, `partial`, `completed`, `failed`, `cancelled`, `expired`.

Ces états décrivent le travail de la capability ; ils ne finalisent pas la machine d’état canonique des objets, reportée à la phase Objets.

## 15. États d’interface

Loading conserve query et scope ; Partial liste les sources échouées ; Error garde les résultats valides ; Offline interdit une nouvelle exécution ; Permission denied masque champs/raw ; Stale affiche freshness par source. Le Design System conserve la propriété du rendu et des interactions génériques.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Search results | projection de Telemetry Events | analyste, Hunt ou Case | source, période, champs et erreurs visibles |
| Execution record | événement de recherche | Trace et audit | requête/version/paramètres/run conservés |
| Selected result context | relations candidates | Case Workspace | non destructif et sans qualification Evidence |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Event Search | sélection d’un Event | Event Inspection | Event ref, raw permission, query, position et return origin | résultat et scroll restaurés |
| Event Search | ajout au Case | Case Workspace | query, run, résultats sélectionnés, période et annotations | recherche restaurée |
| Event Search | organiser une investigation exploratoire | Hunt Management | question, scope, query et résultats | retour à la recherche |

Chaque transition réévalue les permissions, conserve le tenant et l’environnement et ne transfère jamais l’ownership du produit destination.

## 18. Dépendances

- CAP-INV-003
- CAP-INV-004
- CAP-INV-006
- CAP-INV-007
- CAP-INV-008
- Shared Query Engine
- Background Jobs
- Data Quality Service

## 19. Source de vérité

Telemetry Event, Query et Search Job restent propriétaires de leurs sources Shared ; Investigate possède l’usage analytique, les liens au Case et les annotations locales.

## 20. Provenance et audit

Texte ou builder normalisé, Query/version, auteur, paramètres, période, tenant/env, sources, run, timestamps, erreurs, résultats sélectionnés, assistance éventuelle et exports.

## 21. Permissions fonctionnelles

Besoins : event search, raw event access, source-specific fields, cross-tenant interdiction par défaut, saved search et Case link. Les namespaces atomiques, le step-up et la séparation des tâches définitive restent à la phase Permissions.

## 22. Limites et erreurs

Query invalide, source indisponible, résultat tronqué, timeout, coût refusé, données stale ou permission réduite sont visibles ; aucun résultat complet n’est simulé.

## 23. Métriques

- temps validation→premier résultat
- part des jobs avec provenance complète
- fréquence des résultats partiels
- taux de retours restaurant query et position

Aucune cible chiffrée définitive n’est fixée en Phase 4B.1.

## 24. Classification de livraison

`defined` / `planned`, cible native à long terme ; aucun langage, moteur ou index n’est choisi.

## 25. Critères d’acceptation

**Given** aucun fournisseur de modèle, une source autorisée et des champs connus  
**When** l’utilisateur construit et exécute une requête  
**Then** validation, résultats, erreurs partielles, fraîcheur, pivots et liaison au Case restent disponibles

**Given** une requête touche deux sources dont une échoue  
**When** le job termine  
**Then** les résultats valides sont conservés, la source échouée est nommée et l’état est partial

**Given** un résultat est lié à un Case  
**When** l’utilisateur revient à la recherche  
**Then** query, période, sélection, filtres et position sont restaurés

## 26. Questions ouvertes

Quels contrats fonctionnels de coût, rétention du résultat et compatibilité de dialectes sont nécessaires ? La syntaxe et l’architecture restent aux phases Technique/Implémentation.

## 27. Consommateurs documentaires

Signals and Hunt, Cases and Evidence, futurs parcours Hunt/Case, écrans Event Search/Inspector, phases Query Engine et Permissions.
