---
id: CAP-CMD-108
title: Bulk Coordination
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-004
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-SEC-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---
# CAP-CMD-108 — Bulk Coordination

## 1. Définition
Applique à une sélection autorisée des modifications réversibles homogènes avec preview, exclusions, résultats par item et échecs partiels.

## 2. Problème utilisateur
Répéter une mutation sur plusieurs items est lent, mais une sélection hétérogène peut produire des effets inattendus. Principal : Incident Commander ; secondaires : SOC Analyst L2, Team Lead.

## 3. Objectifs
Prévisualiser scope/exclusions ; limiter aux classes 0/2 réversibles ; réévaluer permission/version par item ; rendre partial success visible/auditable.

## 4. Non-objectifs
Aucun containment/destruction groupé, aucune permission globale substituée, aucune atomicité fictive, aucune colonne finale définie.

## 5. Propriétaire
Command possède l’orchestration batch des mutations Command ; Background Jobs/Export Engine restent Shared.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Team Lead.

## 7. Conditions d’entrée
Sélection non vide ; action compatible ; permissions par item ; versions disponibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Selection | Unified Work Queue | IDs/versions | oui | preview | refuser |
| Parameters | utilisateur | assignment/priority/label/state/export | oui | validés | erreurs visibles |
| Permission evaluation | Permission Model | allowed/excluded | oui | exécution | exclure sans fuite |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | version/state/permission/champs | lecture |
| Saved View / filters | Shared/Command | scope | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | assignment, priority, label, état limité | Command | classe 2 par item |
| Export job | créer | Shared Export Engine | classe 0, aucune mutation objet |

## 11. Fonctionnalités
Preview count/scope/changes/exclusions ; mutations limitées ; export ; partial success et retry ciblé ; audit lot + item.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Prévisualiser | coordinateur | selection | 0 | stable | preview | non |
| Affecter en masse | coordinateur | items | 2 | autorisés | résultats/item | OPEN-013 |
| Priority/label | coordinateur | items | 2 | valeurs valides | résultats/item | OPEN-013 |
| État limité | coordinateur | items | 2 | transition valide | résultats/item | OPEN-013 |
| Exporter | rôle export | projection | 0 | redaction | job | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Preview/validation | confirmation | oui | oui | possible | résumé | non | preview déterministe |
| Exécution | déclencheur | policy | item/version | possible | jamais autonome | OPEN-013 | action humaine |

## 14. États fonctionnels
`prepared`, `validated`, `executing`, `partial-success`, `completed`, `cancelled`, `conflicted`.

## 15. États d’interface
Loading affiche progression ; Partial liste succès/échecs/exclusions ; Error garde résultats acquis ; Offline interdit lancement ; Permission denied exclut sans fuite ; Stale impose nouvelle preview. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Bulk result | batch event | utilisateur/audit | succès/échec/exclusion par item |
| Retry selection | state | Queue | items sûrs non réussis uniquement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Bulk | export | Jobs/Export | selection, classification, view | notification/return |
| Failure | retry | Queue | failed IDs, params | nouvelle preview |

## 18. Dépendances
CAP-CMD-101, 102, 104, Background Jobs, Export Engine, audit hooks.

## 19. Source de vérité
Chaque item Command reste source de sa version ; le batch n’est qu’un agrégat de résultats par item.

## 20. Provenance et audit
Lot et chaque item enregistrent acteur, action, params, version, résultat, exclusion, tenant et correlation ID.

## 21. Permissions fonctionnelles
Permissions de chaque mutation, future permission bulk, permission export distincte. `OPEN-013` reste ouverte.

## 22. Limites et erreurs
Sélection stale/hétérogène, permission réduite, conflit, partial failure, cancellation ou dépendance indisponible ne doivent jamais être présentés comme transaction atomique.

## 23. Métriques
Partial success ; exclusions avant exécution ; items modifiés sans preview (cible conceptuelle zéro).

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** dix items dont deux interdits, **When** preview est exécutée, **Then** scope, deux exclusions et changements sont visibles avant confirmation.

**Given** trois conflits pendant exécution, **When** le lot finit, **Then** succès et échecs sont détaillés, retry ne contient que les échecs sûrs.

**Given** aucun modèle IA, **When** bulk est utilisé, **Then** preview/validation/exécution déterministes suffisent.

## 26. Questions ouvertes
Permission bulk distincte ? Quelles transitions suffisamment réversibles ? — REQ-PROD-004, REQ-PROD-009, REQ-PROD-013, REQ-SEC-001. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Jobs/Export, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
