---
id: CAP-GOV-045
title: Govern Queue, Ageing and Lifecycle Flow Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-UX-008, REQ-SEC-001]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-045 — Govern Queue, Ageing and Lifecycle Flow Metrics

## 1. Définition
Définir les métriques Govern de volume, ageing et temps de passage du Response Inbox et du lifecycle Action Request→Decision→Run→Verification afin d’identifier les étapes d’attente/blocage sans recréer les métriques de Command Work Queue ni fixer un SLO universel.

## 2. Problème utilisateur
Un backlog Govern peut augmenter parce que des requests sont incomplètes, attendent une authority review, une Approval ou une fenêtre d’exécution. Une seule moyenne d’âge masque où le flux bloque et peut être confondue avec la charge Command.

## 3. Objectifs
- mesurer Inbox volume/ageing ;
- mesurer durée `incomplete`, `information-requested`, `policy-review`, `authority-review`, `approval-pending`, `decision-ready` ;
- mesurer Decision-to-run et run-to-verification delay ;
- compter open/blocked/superseded/expired ;
- segmenter backlog par action class/tenant/environment ;
- distinguer Govern flow de Command Work Queue.

## 4. Non-objectifs
Ne pas redéfinir Work Queue Command, créer un SLO universel, modifier priorité/assignment, juger la qualité par vitesse, finaliser la state machine Action Request ou définir un scheduling engine.

## 5. Propriétaire
Govern / Response Metrics owns Govern lifecycle-flow metric semantics. Response Inbox/GOV-1 and Runs/GOV-2 own source states; Command owns its own Work Queue metrics; Shared owns generic Metrics/Reporting mechanisms.

## 6. Utilisateurs
Principal : Govern Coordinator / Control Reviewer. Secondaires : Govern Reviewer, Approval Coordinator, Decision Reviewer, Response Operator, Govern Product Lead and authorized service-management consumers.

## 7. Conditions d’entrée
Lifecycle transitions carry enough timestamps/state refs for the selected scope; metric definitions identify included states and clock rules; missing transitions remain explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Response Inbox/lifecycle events | CAP-GOV-001..003 | queue/flow observations | oui | snapshot | partial flow |
| policy/authority/approval states | CAP-GOV-007..013 | stage timing | when applicable | historical | stage unknown/NA |
| Decision timing | CAP-GOV-014..016 | decision stage | oui for decided requests | snapshot | decision gap |
| Run/start/verification timing | CAP-GOV-022/023/028/029 | downstream flow | when executed | snapshot | no downstream duration |
| tenant/environment/action class | source context | segmentation | oui when dimension used | snapshot | aggregate only |
| metric definition/snapshot | Shared Metrics | calculation contract | oui | versioned | no metric output |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request / Queue Item | Govern | lifecycle/assignment/timestamps | metric read |
| Policy/Authority/Approval records | Govern | wait-stage refs | metric read |
| Decision | Govern | decision-ready/final times | metric read |
| Response Run/Verification | Govern | downstream times | metric read |
| Command Work Queue | Command | boundary/comparison only where authorized | no metric ownership |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Govern Flow Metric Observation | calculate/version/supersede | Govern local concept | stage semantics explicit |
| backlog/ageing observation | calculate/version | Govern local concept | no Command Work Queue redefinition |
| source lifecycle objects | no mutation | Govern source capabilities | aggregate only |

## 11. Fonctionnalités
Count inbox/backlog; calculate item age and stage durations; separate waiting-on-information/policy/authority/approval/decision/run/verification; calculate Decision-to-run and Run-to-verification delays; count blocked/superseded/expired; segment by action class/tenant/environment; compare periods; expose missing transitions; distinguish active queue from historical throughput.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter flow metrics | reviewer | observation | 0 | permission | scoped flow view | non |
| calculate ageing/stage duration | reviewer/system | observation | 1 | event pairs | deterministic values | non |
| annotate bottleneck candidate | reviewer | observation | 2 | rationale | attributed candidate | OPEN-013 |
| prepare comparison/report | reviewer | observation set | 2 | comparable definitions | comparison | OPEN-013 |
| reassign/change priority from metric | none | source queue | 2/3 | outside capability | no action | source workflow only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| calculate ages/stage durations | oui | oui | oui | explain | flow table |
| group backlog by stage/class | oui | oui | oui | summary | pivot |
| flag bottleneck candidate | oui | comparison/ranking | oui | hypothesis | sorted table |
| draft flow summary | oui | templates | oui | sourced draft | analyst note |
| change workflow/SLO automatically | human owner process | no | no | prohibited | explicit improvement review |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `transition-incomplete`, `insufficient-data`, `privacy-limited`, `definition-changed`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading preserves period/stage ; Empty distinguishes zero backlog from unavailable source ; Partial lists missing transitions ; Error retains prior snapshot ; Offline read-only ; Permission denied masks tenant/actor dimensions ; Stale shows source age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Queue/Ageing Observation | metric | CAP-GOV-046/Reporting | stage/scope/denominator explicit |
| lifecycle-flow durations | metric | CAP-GOV-046/047 | missing stages not coerced to zero |
| bottleneck candidate | derived observation | CAP-GOV-046/047 | hypothesis only |
| period comparison | comparison | control reviewer | definitions/differences visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| GOV-1/GOV-2 lifecycle events | snapshot selected | CAP-GOV-045 | states/timestamps/dimensions | source objects |
| CAP-GOV-045 | trend/control review | CAP-GOV-046 | flow observations | Response Metrics |
| CAP-GOV-045 | improvement candidate | CAP-GOV-047 | stage/bottleneck evidence | Response Metrics |
| CAP-GOV-045 | report/export | Shared Reporting/Export | snapshot/definitions | Response Metrics |

## 18. Dépendances
CAP-GOV-001..003/007..016/022/023/028/029/034/046/047, Shared Metrics/Reporting, Settings tenant context, Command boundary, OPEN-010/013.

## 19. Source de vérité
Govern lifecycle objects remain source for Govern flow. Command Work Queue remains Command-owned and is not merged. Govern owns derived flow metric semantics only.

## 20. Provenance et audit
Record metric definition/version, lifecycle event pairs, clocks/exclusions, snapshot/time window, dimensions, missing transitions, privacy suppression, reviewer notes and AI summaries.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; restricted queue/tenant dimension read; cross-tenant flow comparison separately authorized; Reporting/Export preparation. No assignment/priority permission implied.

## 22. Limites et erreurs
Missing timestamps, state-machine changes, unresolved time zones/clock sources, retention gaps, superseded requests and mixed deployment models can limit comparisons. Faster flow is not automatically better governance.

## 23. Métriques
Inbox volume/ageing; incomplete/information/policy/authority/approval/decision-ready durations; Decision-to-run; Run-to-verification; open/blocked/superseded/expired; backlog by action class/tenant/environment. No universal SLO.

## 24. Classification de livraison
`defined` / `planned`; no SLO, scheduling engine, queue engine or dashboard implementation selected.

## 25. Critères d’acceptation
**Given** a backlog grows because requests wait on information, **When** flow metrics are reviewed, **Then** information wait is separated from active Govern review time.

**Given** a rejected Decision has no Run, **When** Decision-to-run metric is calculated, **Then** no fabricated zero delay is assigned.

**Given** Command Work Queue also has ageing metrics, **When** Govern metrics are viewed, **Then** the two domains retain distinct objects, owners and definitions.

**Given** no AI, **When** flow metrics are calculated, **Then** deterministic event-pair calculations and human analysis remain complete.

## 26. Questions ouvertes
OPEN-010/013 remain open. Final role-density, SLO and state-machine implementation remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-046/047, Govern closure, Command boundary docs, Shared Reporting/Metrics and future screen/technical phases.