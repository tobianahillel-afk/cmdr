---
id: CAP-EPT-050
title: Process and System-State Collection
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001]
open_decisions: [OPEN-008, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-050 — Process and System-State Collection

## 1. Définition
Produire sur demande un snapshot technique nouveau de processus et état système autorisés — process metadata/relations, modules, sessions, services/composants, connexions et system facts selon support — distinct du contexte EPT-3 déjà disponible.

## 2. Problème utilisateur
EPT-3 peut contenir un contexte stale/incomplet. Un snapshot demandé doit être identifié comme nouvelle acquisition read-only, avec timestamp, catégories manquantes et limites, sans devenir une mutation.

## 3. Objectifs
Distinguer cached context vs fresh requested snapshot ; borner categories ; conserver snapshot time/freshness ; exposer partial/restricted/unsupported ; produire neutral Collection Items ; préparer collection complémentaire sans effectuer containment.

## 4. Non-objectifs
Aucun terminate/suspend/resume process, start/stop service, isolation, remediation, commande technique finale, Artifact/Evidence automatique.

## 5. Propriétaire
Endpoint owns fresh local snapshot acquisition. Investigate owns analytical question/qualification; Govern owns any effectful modification.

## 6. Utilisateurs
Case Analyst, DFIR Analyst, Response Operator, Endpoint Operator, Security Reviewer.

## 7. Conditions d’entrée
Bounded plan categories, target available, read-only inspection capability declared, policy/read permission and authority where applicable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| planned snapshot item | CAP-EPT-048 | categories/bounds | oui | plan version | no snapshot |
| prior EPT-3 context | CAP-EPT-037..042 | comparison baseline | non | visible timestamp | no comparison |
| local capability/health | Endpoint | support | oui | current | unsupported/unavailable |
| policy/permission | Settings/Security | read constraints | oui | current | restricted/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Plan | Endpoint | categories/target | read |
| EPT-3 contexts | Endpoint | old refs/time | read |
| Collection Request/Case | Investigate | reason | read |
| Endpoint Policy | Settings | allowed categories | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| System Snapshot Attempt | create/update | Endpoint | read-only |
| Process/System Collection Item | create | Endpoint | source/timestamp/categories explicit |
| Snapshot Gap Marker | derive | Endpoint | missing/restricted category explicit |

## 11. Fonctionnalités
Request source-backed snapshots of supported process/system categories, record exact snapshot time, compare only when previous context exists, preserve category-level partial/error states, stop before any mutation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect snapshot plan | Analyst | snapshot item | 0 | read | categories visible | non |
| collect fresh read-only snapshot | authorized operator/service | snapshot attempt | 2 | bounds + permission | snapshot item | generally no effectful authority; OPEN-013 if impact |
| prepare process/service mutation | Analyst | action context | 3+ | separate Govern path | no mutation here | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| select declared categories | oui | oui | oui | suggest | checklist |
| compare two snapshots | oui | oui | oui | summarize | deterministic diff |
| mark stale/gaps | oui | oui | oui | explain | timestamp/rules |
| mutate target | non | interdit | non | interdit | Govern/EPT-5 boundary |

## 14. États fonctionnels
`preparing`, `collecting`, `completed`, `partial`, `restricted`, `unsupported`, `stale-result`, `failed`, `timed-out`, `cancelled`.

## 15. États d’interface
No Screen ID; current snapshot, prior context and missing categories are visually/logically distinct.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Process/System Collection Item | Endpoint output | CAP-EPT-054/Investigate | fresh snapshot refs/time |
| category gap/error | diagnostic | CAP-EPT-053 | no hidden category |
| optional diff refs | derived context | Investigate | correlation != maliciousness |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-3 collection-required | request | CAP-EPT-047/048 | missing category/context | EPT-3 origin |
| CAP-EPT-048 | snapshot item | CAP-EPT-050 | categories/target | plan |
| CAP-EPT-050 | output | CAP-EPT-054/055 | item/gaps/time | attempt |

## 18. Dépendances
CAP-EPT-037..048/053..055, Investigate CAP-INV-206, Endpoint investigation sources, OPEN-008/013/014.

## 19. Source de vérité
Endpoint SOT of requested local snapshot facts; EPT-3 remains SOT of its prior context projection; Investigate owns analytical interpretation.

## 20. Provenance et audit
Request/plan, categories, Agent/platform, snapshot time, prior context refs, policy/permission, gaps/errors, actor and correlation.

## 21. Permissions fonctionnelles
Process/system/session/service/connection snapshot read/collect, sensitive output read, cancel, cross-tenant deny; no mutation permission inherited.

## 22. Limites et erreurs
Fresh snapshot acquisition ≠ EPT-3 cached context; snapshot ≠ durable target state; collection ≠ mutation/containment; output ≠ Evidence/Finding/Artifact automatically.

## 23. Métriques
Snapshot category coverage, partial/stale/error rates, EPT-3→collection pivots, mutation requests correctly redirected.

## 24. Classification de livraison
`draft / defined / planned`; no inspection engine/command/protocol.

## 25. Critères d’acceptation
**Given** EPT-3 context is stale, **When** a fresh process snapshot is authorized, **Then** new snapshot time/source is recorded rather than silently replacing history.

**Given** a service category is restricted, **When** snapshot runs, **Then** it is marked restricted and other valid categories may remain partial-success.

**Given** a terminate action is requested, **When** EPT-4 handles it, **Then** it stops at Govern/EPT-5 boundary and performs no termination.

## 26. Questions ouvertes
OPEN-008/013/014 remain open; exact supported categories vary by declared capability.

## 27. Consommateurs documentaires
EPT-4 packaging/handoff, Investigate, Govern boundary, EPT-5 future, Quality.
