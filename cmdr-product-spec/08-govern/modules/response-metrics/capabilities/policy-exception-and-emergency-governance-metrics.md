---
id: CAP-GOV-039
title: Policy, Exception and Emergency Governance Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-039 — Policy, Exception and Emergency Governance Metrics

## 1. Définition
Définir les observations métriques Govern sur Policy Evaluation, conflicts, Exception Candidates/authorized exceptions et emergency paths afin de suivre volumes, distributions, missing inputs, expirations et re-reviews sans transformer un taux en jugement automatique sur la qualité de la gouvernance.

## 2. Problème utilisateur
Un faible taux d’exception peut refléter une bonne gouvernance ou un processus inutilisable ; un taux élevé peut refléter un contexte légitime. Sans définitions sourcées et dimensions explicites, les mêmes nombres deviennent des conclusions opposées.

## 3. Objectifs
- compter Policy evaluations par outcome `pass/warn/block/unknown/not-applicable` ;
- mesurer conflicts et missing-input frequency ;
- observer Exception Candidate, approval/refusal/expiry/revocation ;
- observer emergency-path usage et retrospective-review references ;
- comparer Policy versions/scopes/périodes ;
- rendre définition, snapshot, fraîcheur et limitations visibles.

## 4. Non-objectifs
Ne pas fixer un seuil universel, créer un KPI/SLO automatique, conclure qu’une exception est bonne/mauvaise, compter `block` comme incident évité, modifier Policy/Exception, créer un Metrics Engine ou choisir un warehouse/dashboard.

## 5. Propriétaire
Govern / Response Metrics owns metric semantics and Govern dimensions. Shared Metrics Engine owns generic definition/reconciliation/privacy/freshness mechanics; Reporting owns rendering; GOV-1 owns Policy/Exception/Emergency source records.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Policy Reviewer, Security/Compliance Reviewer, Decision Maker, Govern Product Lead, Auditor et authorized operational leaders.

## 7. Conditions d’entrée
Policy/Exception/Emergency source records are available for a permissioned tenant/environment/time window with stable definitions and known snapshot/freshness metadata.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Policy Evaluation outcomes | CAP-GOV-007 | metric observations | oui | selected snapshot | metric partial |
| conflicts/Exception Candidates | CAP-GOV-008 | governance observations | when applicable | event snapshot | explicit zero vs unavailable distinguished |
| emergency path records | CAP-GOV-013 | exceptional governance observations | when present | event snapshot | unavailable/none distinguished |
| Policy version/scope | Govern/Security source | dimension | oui for comparisons | historical version | dimension unknown |
| tenant/environment/action class | source context | segmentation | oui | snapshot | aggregate blocked |
| metric definition/snapshot | Shared Metrics | calculation context | oui | versioned | no metric published |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Policy / Policy Evaluation | Govern/Security source | version/outcome/scope | metric read |
| Exception Candidate/authorized exception context | Govern | disposition/expiry | metric read |
| Emergency Authorization Context | Govern | usage/scope/review refs | restricted metric read |
| Metric definition/snapshot | Shared | version/freshness/privacy | consume only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Policy Metric Observation | calculate/version/supersede | Govern local concept | source set + definition version required |
| Exception/Emergency Metric Observation | calculate/version/supersede | Govern local concept | no causal/quality conclusion embedded |
| source Policy/Exception/Emergency records | no mutation | GOV-1 owners | read/aggregate only |

## 11. Fonctionnalités
Define metric formulas conceptually; aggregate outcomes; segment by Policy version/action class/tenant/environment/time; distinguish zero vs unavailable; calculate conflict/missing-input frequency; count candidate/approved/refused/expired/revoked exceptions; count emergency uses and review references; compare periods; expose denominators/freshness/privacy; hand observations to trends/reporting.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter metric | reviewer | Metric Observation | 0 | metric permission | scoped observation | non |
| calculate/reconcile metric | reviewer/system | Metric Observation | 1 | definition+snapshot | deterministic value/unknown | non |
| annotate interpretation | reviewer | observation | 2 | rationale | attributed note | OPEN-013 |
| prepare comparison/report | reviewer | observation set | 2 | scope/privacy valid | derived comparison | OPEN-013 |
| change Policy/Exception automatically | none | source | 3/4 | prohibited | no action | owner workflow only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate outcomes | oui | oui | oui | explain | pivot/table |
| compare versions/periods | oui | oui | oui | summary | comparison table |
| detect missing-input trend | oui | threshold-free rules | oui | candidate explanation | trend table |
| draft interpretation | oui | templates | oui | sourced draft | reviewer note |
| decide Policy quality/change | human owner process | no automatic | no | prohibited | explicit review |

## 14. États fonctionnels
`not-calculated`, `calculating`, `available`, `partial`, `insufficient-data`, `privacy-limited`, `stale`, `definition-changed`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains selected dimensions ; Empty distinguishes true zero from unavailable data ; Partial names missing sources ; Error keeps last valid snapshot ; Offline read-only ; Permission denied suppresses restricted dimensions ; Stale shows snapshot age/definition version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Policy Metric Observation | metric observation | CAP-GOV-046/Reporting | definition/snapshot/dimensions visible |
| Exception/Emergency observations | metric observations | control review/CAP-GOV-046 | no health conclusion encoded |
| missing-input/control signal candidate | derived observation | CAP-GOV-047 | candidate only |
| period/version comparison | comparison artifact | control reviewer | same definitions or differences explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-007/008/013 | observation snapshot selected | CAP-GOV-039 | outcomes/versions/context | source records |
| CAP-GOV-039 | comparison requested | CAP-GOV-046 | metric definitions/snapshots/dimensions | Response Metrics |
| CAP-GOV-039 | improvement candidate | CAP-GOV-047 | sourced observation/limitation | Response Metrics |
| CAP-GOV-039 | report/export | Shared Reporting/Export | metric refs/snapshot/classification | Response Metrics |

## 18. Dépendances
CAP-GOV-007/008/013/034/046/047, Shared Metrics/Reporting/Export, Settings tenant context, Security privacy/permissions, OPEN-013/019.

## 19. Source de vérité
GOV-1 records remain source for Policy/Exception/Emergency events. Govern owns metric definitions in its domain and derived observations; Shared owns generic metric computation infrastructure. A metric observation is not a Policy or Decision.

## 20. Provenance et audit
Record metric definition/version, source query/snapshot refs, numerator/denominator semantics, dimensions, exclusions, freshness, privacy suppression, reviewer notes, AI summaries and supersession.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read` family; sensitive Policy/Exception dimension read; cross-tenant metric read only when separately authorized; report/export preparation; no permission to mutate source Policy/Exception implied.

## 22. Limites et erreurs
Missing policy inputs, changed definitions, sparse/private dimensions, retention gaps, duplicate source events or mixed Policy versions can make a metric partial/incomparable. A low/high value is not automatically healthy/unhealthy.

## 23. Métriques
This capability defines: evaluation count/distribution, conflict frequency, candidate/approved/refused/expired exception counts, emergency-use count, re-review count, missing-input frequency and Policy-version distribution. It sets no target threshold.

## 24. Classification de livraison
`defined` / `planned`; no Metrics Engine, warehouse, dashboard, SLO or implementation is selected.

## 25. Critères d’acceptation
**Given** a high exception rate, **When** metrics are reviewed, **Then** rate/denominator/context are shown and no automatic unhealthy conclusion is produced.

**Given** Policy `block` count rises, **When** a report is prepared, **Then** the count is not labelled “prevented incidents” without separate evidence.

**Given** emergency records lack retrospective-review references, **When** metrics are calculated, **Then** the missing-review dimension is visible without declaring misconduct.

**Given** no AI, **When** metrics are produced, **Then** deterministic aggregation and tables provide complete functionality.

## 26. Questions ouvertes
OPEN-013 and OPEN-019 remain open. Final metric engine/thresholds/compliance interpretation remain future. No new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-046/047, Govern closure, Shared Reporting/Metrics integration, future screens/permissions/technical implementation and quality reports.