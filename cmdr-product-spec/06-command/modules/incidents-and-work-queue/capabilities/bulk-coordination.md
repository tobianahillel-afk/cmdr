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
Répéter une mutation est lent, mais une sélection hétérogène peut produire des effets inattendus. Sans preview et résultat par item, l’opérateur suppose une atomicité inexistante.

## 3. Objectifs
Prévisualiser scope/exclusions, limiter aux classes 0/2 réversibles, réévaluer permission/version par item et rendre partial success visible/auditable.

## 4. Non-objectifs
Aucun containment/destruction groupé, aucune permission globale substituée, aucune atomicité fictive et aucune définition finale des colonnes.

## 5. Propriétaire
Command possède l’orchestration batch des mutations Command ; Background Jobs et Export Engine restent Shared.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Team Lead.

## 7. Conditions d’entrée
Sélection non vide, action compatible, permissions par item et versions disponibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Selection | Unified Work Queue | IDs et versions | oui | vérifiée au preview | opération refusée |
| Parameters | utilisateur | assignment, priority, label, état ou export | oui | validés avant exécution | erreurs de paramètres visibles |
| Permission evaluation | Permission Model | allowed et excluded par item | oui | réévaluée à l’exécution | item exclu sans fuite de données |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | version, état, permission et champs ciblés | consulter et prévisualiser |
| Saved View / filters | Shared/Command | scope de sélection | consulter et reproduire le scope |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | assignment, priority, label ou état limité | Command | classe 2 par item, pas de transaction fictive |
| Export job | créer | Shared Export Engine | classe 0, aucune mutation métier |
| Bulk result | créer un événement de lot | Shared Audit / Command | résultats détaillés par item |

## 11. Fonctionnalités
Preview count/scope/changes/exclusions, mutations limitées, export, partial success, retry ciblé et audit lot + item.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Prévisualiser | coordinateur | sélection | 0 | versions disponibles | preview | non |
| Affecter, prioriser ou étiqueter | coordinateur | items | 2 | items autorisés | résultats par item | OPEN-013 |
| Changer un état limité | coordinateur | items | 2 | transitions valides | résultats par item | OPEN-013 |
| Exporter | rôle export | projection | 0 | redaction et classification | job | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Construire le preview | déclenché par humain | oui | oui | résumé facultatif | moteur de validation déterministe |
| Exclure les items interdits | revue possible | oui | oui | non nécessaire | Permission Model par item |
| Exécuter le lot | confirmation humaine | oui | workflow possible | jamais autonome | action humaine après preview |
| Préparer un retry | oui | oui | oui | suggestion possible | sélection déterministe des échecs sûrs |

## 14. États fonctionnels
`prepared`, `validated`, `executing`, `partial-success`, `completed`, `cancelled`, `conflicted`.

## 15. États d’interface
Loading affiche progression ; Partial liste succès/échecs/exclusions ; Error garde les résultats acquis ; Offline interdit lancement ; Permission denied exclut sans fuite ; Stale impose un nouveau preview.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Bulk result | événement de lot | utilisateur et Audit | succès, échec et exclusion par item |
| Retry selection | état de sélection | Unified Work Queue | uniquement les items non réussis et sûrs |
| Export job | job Shared | utilisateur autorisé | permission-aware, redacted et non destructif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Bulk Coordination | export utilisateur | Background Jobs / Export Engine | sélection, classification et vue | notification puis même Work Queue |
| Bulk result | retry utilisateur | Unified Work Queue | IDs échoués, paramètres et versions | nouveau preview obligatoire |
| Bulk preview | action owner-specific | capability propriétaire | sélection, action et préconditions | résultats par item dans la file |

## 18. Dépendances
CAP-CMD-101, CAP-CMD-102, CAP-CMD-104, Background Jobs, Export Engine et audit hooks.

## 19. Source de vérité
Chaque item Command reste source de sa version ; le batch n’est qu’un agrégat de résultats par item.

## 20. Provenance et audit
Lot et item enregistrent acteur, action, paramètres, version, résultat, exclusion, tenant et correlation ID.

## 21. Permissions fonctionnelles
Permissions de chaque mutation, future permission bulk et permission export distincte ; `OPEN-013` reste ouverte.

## 22. Limites et erreurs
Sélection stale/hétérogène, permission réduite, conflit, partial failure, cancellation ou dépendance indisponible ne sont jamais présentés comme transaction atomique.

## 23. Métriques
Partial success, exclusions avant exécution et items modifiés sans preview ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** dix items dont deux interdits, **When** preview est exécuté, **Then** scope, exclusions et changements sont visibles avant confirmation.

**Given** trois conflits, **When** le lot finit, **Then** succès/échecs sont détaillés et retry contient seulement les échecs sûrs.

**Given** aucun modèle IA, **When** bulk est utilisé, **Then** preview, validation et exécution déterministes suffisent.

## 26. Questions ouvertes
Permission bulk distincte et transitions suffisamment réversibles ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Jobs/Export, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.