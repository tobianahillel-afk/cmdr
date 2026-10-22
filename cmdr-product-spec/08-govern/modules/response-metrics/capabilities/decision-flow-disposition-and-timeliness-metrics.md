---
id: CAP-GOV-041
title: Decision Flow, Disposition and Timeliness Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-041 — Decision Flow, Disposition and Timeliness Metrics

## 1. Définition
Définir les métriques Govern de volume, disposition et temps de passage des Decisions afin d’observer `approve`, `approve-with-conditions`, `reject`, `defer`, `request-more-information`, supersession/expiry et les temps de préparation/attente/Decision-to-run, sans assimiler rapidité à qualité ni fixer un SLO universel.

## 2. Problème utilisateur
Un processus de Decision peut être lent à cause d’informations manquantes, d’Approvals nécessaires ou d’une forte complexité. Une moyenne agrégée sans décomposition pousse à conclure que « plus rapide = meilleur » et masque les causes d’attente.

## 3. Objectifs
- compter Decisions et dispositions ;
- mesurer preparation duration, waiting-on-information et waiting-on-approval ;
- mesurer Decision-to-run delay lorsque le Run existe ;
- observer conditional/superseded/expired Decisions ;
- segmenter par action class/tenant/environment/module/période ;
- exposer distributions, denominators, snapshot/freshness et limitations.

## 4. Non-objectifs
Ne pas définir un SLO universel, scorer la qualité d’un Decision Maker, conclure qu’une Decision rapide est meilleure, modifier Decision/Run, créer un workflow engine ou imposer un dashboard.

## 5. Propriétaire
Govern / Response Metrics owns Decision-flow metric semantics. GOV-1 owns Decision source records; GOV-2 owns Run source records; Shared owns generic Metrics/Reporting infrastructure.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Decision Reviewer/Maker, Govern Product Lead, Security/Compliance Reviewer, Response Operator and authorized operational leadership.

## 7. Conditions d’entrée
Decision lifecycle records and relevant request/approval/run timestamps are available for a scoped snapshot; definitions specify start/stop events and exclusions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision records/dispositions | CAP-GOV-014..016 | decision observations | oui | historical snapshot | metric unavailable/partial |
| request lifecycle times | CAP-GOV-003 | flow context | for preparation/waiting metrics | historical | component unknown |
| Approval waiting times | CAP-GOV-011 | wait context | when required | historical | explicit NA/unknown |
| Response Run start refs | CAP-GOV-022/023 | Decision-to-run endpoint | when run exists | correlated snapshot | no delay metric for item |
| action/tenant/environment dimensions | Govern/Settings context | segmentation | oui for dimension | snapshot | aggregate only |
| metric definition/snapshot | Shared Metrics | calculation contract | oui | versioned | no published metric |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | submitted/info wait timestamps | metric read |
| Approval | Govern | approval wait timestamps | metric read |
| Decision | Govern | disposition/version/timestamps | metric read |
| Response Run | Govern | first authorized start reference | metric read |
| Metric definition/snapshot | Shared | formula/freshness/privacy | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Decision Metric Observation | calculate/version/supersede | Govern local concept | explicit start/end definitions |
| flow-duration observation | calculate/version | Govern local concept | missing phases not collapsed into zero |
| source Decision/Run | no mutation | Govern source capabilities | metric read only |

## 11. Fonctionnalités
Aggregate Decision volume/dispositions; calculate preparation and wait durations; distinguish active vs expired/superseded; calculate Decision-to-run delay; segment by action class/tenant/environment/Decision type; compare periods; expose percentile/distribution conceptually without fixed SLO; identify missing timestamp/source coverage; feed trend/control-health review.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter metrics | reviewer | observation | 0 | permission | scoped metrics | non |
| calculate durations/distributions | reviewer/system | observation | 1 | defined event pairs | deterministic observations | non |
| annotate interpretation | reviewer | observation | 2 | rationale | attributed note | OPEN-013 |
| prepare comparison | reviewer | observation set | 2 | comparable definitions | comparison | OPEN-013 |
| change Decision/SLO automatically | none | Decision/policy | 3/4 | prohibited | no action | owner process only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate dispositions | oui | oui | oui | explain | pivot table |
| calculate stage durations | oui | timestamp rules | oui | summary | duration table |
| compare periods | oui | oui | oui | trend narrative | comparison chart/table |
| identify long-wait candidates | oui | ranking without target claim | oui | candidate explanation | sorted list |
| judge Decision quality | human review | no | no | prohibited | separate review evidence |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `timestamp-incomplete`, `definition-changed`, `insufficient-data`, `privacy-limited`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains period ; Empty distinguishes zero Decisions from unavailable data ; Partial lists missing stage timestamps ; Error retains prior snapshot ; Offline read-only ; Permission denied suppresses restricted dimensions ; Stale exposes snapshot age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Decision disposition observation | metric | CAP-GOV-046/Reporting | exact definition/snapshot |
| stage-duration observations | metric | control review | stage boundaries explicit |
| Decision-to-run delay | metric | CAP-GOV-046 | only where correlation exists |
| flow-friction candidate | derived signal | CAP-GOV-047 | candidate, not quality judgment |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-003/011/014..016/022..023 | snapshot selected | CAP-GOV-041 | lifecycle timestamps/dispositions | source records |
| CAP-GOV-041 | trend review | CAP-GOV-046 | observations/definitions | Response Metrics |
| CAP-GOV-041 | improvement candidate | CAP-GOV-047 | waiting-stage evidence | Response Metrics |
| CAP-GOV-041 | report/export | Shared Reporting/Export | metric refs/snapshot | Response Metrics |

## 18. Dépendances
GOV-1 request/Approval/Decision sources, GOV-2 Run source, Shared Metrics/Reporting, Settings tenant context, Security privacy, OPEN-013.

## 19. Source de vérité
Govern lifecycle records are source. Govern owns metric semantics/observations, not a generic Metrics Engine. Duration is derived from named lifecycle events and cannot imply quality by itself.

## 20. Provenance et audit
Record definition/version, event-pair rules, source snapshots, excluded/incomplete records, dimensions, aggregation method, freshness, privacy suppression, notes and AI summaries.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; Decision/restricted context metric read; cross-tenant comparisons separately authorized; Reporting/Export preparation. Metric access never grants Decision mutation.

## 22. Limites et erreurs
Clock inconsistencies, missing event pairs, supersession, partial history, retention and definition changes can make periods incomparable. Faster Decision never automatically equals better Decision.

## 23. Métriques
Decision volume and disposition; preparation duration; waiting-on-information; waiting-on-approval; Decision-to-run delay; conditional/superseded/expired frequency. No universal SLO.

## 24. Classification de livraison
`defined` / `planned`; no SLO, dashboard, engine or implementation chosen.

## 25. Critères d’acceptation
**Given** Decision latency becomes faster, **When** metrics are reviewed, **Then** the observation is shown without a quality conclusion.

**Given** a Decision waits primarily on information, **When** flow is decomposed, **Then** information wait is separated from approval wait and active review time.

**Given** no Run exists after a rejected Decision, **When** Decision-to-run metrics calculate, **Then** the item is not assigned a fabricated zero delay.

**Given** no AI, **When** Decision metrics are produced, **Then** deterministic timestamp aggregation and human analysis are complete.

## 26. Questions ouvertes
OPEN-013 remains open. Final SLO/targets and event model remain future. No new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-046/047, Govern closure, Shared Reporting/Metrics and future screen/technical phases.