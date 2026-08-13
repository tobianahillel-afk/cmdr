---
id: CAP-SET-011
title: Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
delivery_status: defined
delivery_mode: planned
updated: 2026-08-13
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-009, REQ-PROD-010, REQ-PROD-011, REQ-PROD-012, REQ-AI-001, REQ-AI-003, REQ-AI-004, REQ-AI-009, REQ-AI-010, REQ-AI-011]
open_decisions: [OPEN-008, OPEN-012, OPEN-013]
source-of-truth: canonical
---
# CAP-SET-011 — Model Routing Configuration, Eligibility, Fallback Constraints and Provider Switch Provenance

## 1. Définition
Configuration administrative provider-neutral de l'eligibility, du routage et des contraintes de fallback, avec provenance des changements et des sélections effectivement observées.

## 2. Problème utilisateur
L'administrateur doit configurer des choix autorisés sans confondre configuration Settings et moteur runtime.

## 3. Objectifs
Configuration tenant-scoped, versionnée, auditable, explicable sans IA et séparée de l'exécution technique.

## 4. Non-objectifs
Aucun canonical Model, Model Route, Routing Policy, Provider Switch, moteur de routage, automatic failover, API, protocole, nouvel écran, nouvelle permission ou implémentation.

## 5. Propriétaire
`Platform Settings Product Lead`. Les owners Security, Govern, Studio et runtime conservent leurs responsabilités existantes.

## 6. Utilisateurs
Platform Administrator principal; consommateurs autorisés en lecture/projection seulement.

## 7. Conditions d'entrée
Tenant, acteur, version courante, références Model Provider valides et contraintes d'usage disponibles lorsqu'elles sont requises. Environment seulement si sourcé.

## 8. Entrées fonctionnelles
| Entrée | Source | Requise | Si absente |
|---|---|---|---|
| Tenant | Settings | oui | refus |
| Model Provider refs | Settings | oui selon règle | configuration invalide |
| eligibility / allowed-use | source autorisée | selon usage | unknown/non éligible |
| priorité/conditions | source approuvée | si sourcée | rien n'est inventé |
| coût/latence | source autorisée | optionnel | aucune optimisation déduite |
| fallback constraints | configuration | selon besoin | aucun failover déduit |
| effective selection | source runtime autorisée | projection seulement | unknown |

## 9. Objets lus
| Objet | Owner | Usage |
|---|---|---|
| Model Provider | Platform Settings | refs, état, version, metadata |
| Tenant | Platform Settings | scope |
| Environment | Platform Settings | seulement si sourcé |
| contraintes Policy/data | owner canonique | consommation seulement |

Aucun canonical Model ou Routing Policy local.

## 10. Objets créés ou modifiés
| Objet | Opération | Règle |
|---|---|---|
| Model Provider configuration | metadata d'eligibility/routing/fallback | Tenant + version + provenance |
| Administrative audit event | outcome administratif | corrélation et provenance |

Aucun nouvel objet canonique.

## 11. Fonctionnalités
Eligibility, préférences/conditions sourcées, metadata coût/latence, fallback constraints, validation locale, versioning, projection de sélection observée et provenance d'un switch observé.

## 12. Actions utilisateur
Class 0 inspection; Class 1 validation locale sans effet externe; Class 2 mutation administrative réversible/versionnée lorsque sourcée. `OPEN-013` reste ouvert.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | IA | Alternative sans IA |
|---|---|---|---|---|
| expliquer configuration | oui | oui | assistance | lecture directe |
| valider localement | oui | oui | explication | règles déterministes |
| suggérer un changement | oui | oui pour contraintes | suggestion | saisie explicite |
| afficher sélection observée | oui | oui pour le fait | explication sourcée | projection brute |

L'IA ne change pas le routage et n'invente pas de sélection.

## 14. États fonctionnels
Aucun nouvel état. Les Model Providers conservent `configured`, `validating`, `active`, `degraded`, `disabled`.

## 15. États d'interface
Loading, Empty, Partial, Error, Offline, Permission denied, Stale. Une sélection non observée reste unknown.

## 16. Sorties
| Sortie | Consommateur | Garantie |
|---|---|---|
| routing configuration | Settings/Studio autorisés | Tenant + version + contraintes sourcées |
| eligibility outcome | admin/audit | résultat local + cause |
| fallback configuration | consommateurs autorisés | contrainte, pas automatic failover |
| effective selection projection | Settings/Studio/audit | source + fraîcheur |
| switch provenance | admin/audit | changement non silencieux si observé |

## 17. Transitions / handoffs
| Source | Déclencheur | Destination | Résultat |
|---|---|---|---|
| SET-MDL-001 | inspect/configure | Settings | configuration projetée |
| Settings | validation locale | admin/audit | résultat déterministe |
| Settings | config consommable | runtime owner | aucune exécution Settings implicite |
| runtime source | sélection observée | Settings/audit | projection sourcée |

## 18. Dépendances
Model Provider, Tenant, Security, Govern Policy si applicable, data-policy owners, Studio/runtime consumers, Administrative Audit.

## 19. Source de vérité
Model Provider object, Models & Providers module docs, Screen/Permission registers, Product Boundaries, Requirements et OPEN.

## 20. Provenance et audit
Toute mutation conserve acteur, Tenant, refs, versions, changements et corrélation. Toute sélection projetée conserve source et fraîcheur.

## 21. Permissions fonctionnelles
Réutilise les permissions Model Provider et les aliases UI existants. Nouveaux IDs = 0. L'administration n'attribue pas le runtime.

## 22. Limites et erreurs
Routing config ≠ runtime. Fallback ≠ automatic failover. Cost/latency metadata ≠ optimizer. Effective selection projection ≠ exécution Settings.

## 23. Métriques
Mutations acceptées/refusées, validations locales, observations de sélection, switches observés, stale/unknown et conflits de version. Aucun KPI runtime inventé.

## 24. Classification de livraison
`defined / planned`. Aucune disponibilité, intégration, moteur, support provider/model ou déploiement n'est revendiqué.

## 25. Critères d'acceptation
**Given** une configuration tenant-scoped, **When** elle change, **Then** version et provenance sont conservées sans créer Routing Policy.  
**Given** un fallback configuré, **When** il est enregistré, **Then** aucun automatic failover n'est affirmé.  
**Given** une sélection réellement observée, **When** elle est affichée, **Then** source et fraîcheur sont visibles.  
**Given** une sélection différente de la précédente, **When** elle est projetée, **Then** le changement n'est pas silencieux.  
**Given** une version stale ou un autre Tenant, **When** une mutation est demandée, **Then** elle est refusée.  
**Given** l'IA indisponible, **When** la configuration est administrée, **Then** le chemin manuel/déterministe reste disponible.

## 26. Questions ouvertes
`OPEN-008`, `OPEN-012` et `OPEN-013` restent ouverts; aucune décision runtime/support/autorité n'est résolue ici.

## 27. Consommateurs documentaires
Models & Providers, Administrative Audit, Security, Studio, Shared, Health si pertinent, Quality et Roadmap, sans transfert d'ownership.
