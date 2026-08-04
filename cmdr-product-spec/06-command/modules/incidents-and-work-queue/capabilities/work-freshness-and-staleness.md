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
Un item peut sembler actuel alors qu’owner, SLA, Case ou Run provient d’une source retardée. Principal : SOC Analyst L1/L2 ; secondaires : Incident Commander, Auditor, Service Delivery Manager. Sans capacité, une décision repose sur une projection périmée.

## 3. Objectifs
Montrer fraîcheur item/champ/vue ; distinguer aging/stale/unavailable/unknown ; préserver données valides ; permettre refresh ou follow-up sûr.

## 4. Non-objectifs
Ne pas masquer stale, inventer un seuil universel, réordonner silencieusement ou réduire l’explication à un timestamp secondaire.

## 5. Propriétaire
Command possède l’interprétation locale de fraîcheur ; chaque source possède ses timestamps/health ; Data Quality/Jobs restent Shared.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Auditor, Service Delivery Manager.

## 7. Conditions d’entrée
Work item ou vue et métadonnées source/time lorsqu’elles existent.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Update metadata | objets/sources | updated-at, source, ingestion | oui | par champ/source | unknown |
| Dependency health | Shared/Settings | availability/degradation | non | courante | unavailable explicite |
| Freshness policy | tenant/module | seuils par type | non | version applicable | aging/unknown sans arbitraire |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | timestamps/source | lecture |
| Case / Decision / Run | produits propriétaires | freshness metadata | projection |
| Health | Settings/Shared | dependency status | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Task | follow-up refresh/validation | Command | classe 2 |
| Acknowledgement | lecture stale | local/audit | ne rend pas la donnée fraîche |

## 11. Fonctionnalités
Afficher dernière mise à jour/source au premier niveau ; calculer un statut explicable ; conserver valeurs valides en Partial/Error ; bloquer mutation non garantie ; créer une action de validation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter fraîcheur | lecteur | item/view | 0 | metadata ou unknown | explication | non |
| Demander refresh | lecteur autorisé | source/job | 0 | idempotent | job/message | non |
| Créer follow-up | coordinateur | Task | 2 | stale affecte décision | Task liée | non |
| Accuser lecture stale | utilisateur | view state | 2 | warning visible | acknowledgement | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Calculer fraîcheur | interprétation/correction | policy | déterministe | possible | résumé | non | timestamps+policy |
| Créer suivi | décision humaine | possible | validation | possible | proposition | non | Task manuelle |

## 14. États fonctionnels
`current`, `aging`, `stale`, `source-unavailable`, `partial`, `unknown`.

## 15. États d’interface
Stale montre source/date/policy ; Partial conserve valeurs sûres ; Error nomme dépendance ; Offline bloque mutation ; Permission denied ne révèle pas le champ source. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Freshness projection | metadata/status | Queue, Mission Control, reporting | source, time, policy/version |
| Validation follow-up | Task | data owner/coordinateur | raison et objet stale liés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Staleness | ouvrir source owner | produit/Settings | object, source, timestamp | vue |
| Staleness bloque action | Task/Escalation | Command | item, consequence, owner | après validation |

## 18. Dépendances
Data Quality Service, Background Jobs, Notification Center, CAP-CMD-107 et CAP-CMD-110.

## 19. Source de vérité
Chaque source possède ses timestamps/health ; Command applique une policy versionnée et expose ses facteurs sans réécrire les données.

## 20. Provenance et audit
Policy/version, source, timestamps, calcul, acknowledgement, refresh et follow-up enregistrent acteur et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, permission refresh/job selon source, `perm.command.task.manage`. Granularité reportée.

## 22. Limites et erreurs
Metadata absente, horloge incohérente, source down, policy manquante, tenant incompatible ou refus produisent unknown/partial, jamais current inventé.

## 23. Métriques
Items avec fraîcheur explicite ; âge projections critiques ; mutations bloquées pour version non fiable. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** une projection Case retardée, **When** la file s’affiche, **Then** source, dernière mise à jour, état stale et conséquence sont visibles.

**Given** une source indisponible, **When** le reste de l’item est valide, **Then** les données locales restent utilisables et seule la projection concernée est Partial.

**Given** aucun modèle IA, **When** la fraîcheur est évaluée, **Then** timestamps et policies déterministes suffisent.

## 26. Questions ouvertes
Quels seuils sont définis par tenant/module/source ? Quelle fraîcheur reste visible dans exports/handovers ? — REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-005. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Work Queue, Mission Control, reporting, handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
