---
id: CAP-CMD-109
title: Work Freshness and Staleness
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-UX-005
open_decisions:
  - none
source-of-truth: canonical
---

# CAP-CMD-109 — Work Freshness and Staleness

## 1. Définition
Expose dernière mise à jour, source, retard, qualité partielle et indisponibilité des dépendances pour chaque work item et pour la vue.

## 2. Problème utilisateur
Un item peut sembler actuel alors qu’owner, SLA, Case ou Run provient d’une source retardée. Sans fraîcheur visible, une décision repose sur une projection périmée.

## 3. Objectifs
Montrer fraîcheur item/champ/vue, distinguer aging/stale/unavailable/unknown, préserver les données valides et permettre refresh ou follow-up sûr.

## 4. Non-objectifs
Ne masque pas stale, n’invente pas un seuil universel, ne réordonne pas silencieusement et ne réduit pas l’explication à un timestamp secondaire.

## 5. Propriétaire
Command possède l’interprétation locale ; chaque source possède timestamps/health ; Data Quality et Jobs restent Shared.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Auditor, Service Delivery Manager.

## 7. Conditions d’entrée
Work item ou vue et métadonnées source/time lorsqu’elles existent.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Update metadata | objets et sources | updated-at, source et ingestion | oui | par champ/source | `unknown` |
| Dependency health | Shared ou Settings | disponibilité et dégradation | non | état courant | `source-unavailable` explicite |
| Freshness policy | tenant ou module | seuils par type | non | version applicable | `aging` ou `unknown`, sans seuil inventé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | timestamps et source | consulter et filtrer |
| Case / Decision / Response Run | produits propriétaires | freshness metadata | consulter en projection |
| Health | Settings ou Shared | dependency status | consulter et relier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Task | créer un suivi de refresh/validation | Command | classe 2, objet stale lié |
| Acknowledgement | accuser lecture stale | Command audit | ne rend pas la donnée fraîche |
| Objets sources | aucune mutation | propriétaires sources | métadonnées en lecture seule |

## 11. Fonctionnalités
Afficher source/date au premier niveau, calculer un statut explicable, conserver valeurs valides, bloquer mutation non garantie et créer une action de validation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter fraîcheur | lecteur | item ou vue | 0 | métadonnée ou unknown | explication visible | non |
| Demander refresh | lecteur autorisé | source/job | 0 | opération idempotente | job ou message | non |
| Créer follow-up | coordinateur | Task | 2 | stale affecte une décision | Task liée | non |
| Accuser lecture stale | utilisateur | view state | 2 | warning visible | acknowledgement | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Calculer la fraîcheur | interprétation humaine possible | oui | oui | résumé facultatif | timestamps et policy versionnée |
| Détecter une dépendance indisponible | oui | oui | oui | non nécessaire | health projection et règles |
| Demander un refresh | oui | job déterministe | oui | non nécessaire | action manuelle de refresh |
| Créer un suivi | oui | validation/déduplication | oui | brouillon Task | Task manuelle |

## 14. États fonctionnels
`current`, `aging`, `stale`, `source-unavailable`, `partial`, `unknown`.

## 15. États d’interface
Stale montre source/date/policy ; Partial conserve les valeurs sûres ; Error nomme la dépendance ; Offline bloque mutation ; Permission denied masque le champ source.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Freshness projection | metadata et statut | Queue, Mission Control et Reporting | source, timestamp et policy/version visibles |
| Validation follow-up | Task | data owner et coordinateur | raison et objet stale liés |
| Refresh request | job ou événement | source propriétaire | attribué, idempotent et sans mutation métier |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Staleness | ouvrir la source | produit propriétaire ou Settings | objet, source, timestamp et return origin | même vue restaurée |
| Staleness | action bloquée | Task Coordination | item, conséquence, owner et validation attendue | vue après validation |
| Staleness | demander refresh | Background Jobs | source, objet et correlation ID | notification puis vue inchangée |

## 18. Dépendances
Data Quality Service, Background Jobs, Notification Center, CAP-CMD-107 et CAP-CMD-110.

## 19. Source de vérité
Chaque source possède timestamps/health ; Command applique une policy versionnée et expose ses facteurs.

## 20. Provenance et audit
Policy/version, source, timestamps, calcul, acknowledgement, refresh et follow-up enregistrent acteur et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, permission refresh/job selon source et `perm.command.task.manage` ; granularité reportée.

## 22. Limites et erreurs
Métadonnée absente, horloge incohérente, source down, policy manquante, tenant incompatible ou refus produisent unknown/partial, jamais current inventé.

## 23. Métriques
Items avec fraîcheur explicite, âge des projections critiques et mutations bloquées ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** une projection Case retardée, **When** la file s’affiche, **Then** source, date, état stale et conséquence sont visibles.

**Given** une source indisponible, **When** le reste est valide, **Then** les données locales restent utilisables et seule la projection est Partial.

**Given** aucun modèle IA, **When** la fraîcheur est évaluée, **Then** timestamps et policies déterministes suffisent.

## 26. Questions ouvertes
Quels seuils sont définis par tenant/module/source et quelle fraîcheur est exportée ? — Requirement IDs ci-dessus. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Work Queue, Mission Control, Reporting, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.