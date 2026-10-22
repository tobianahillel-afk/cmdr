---
id: CAP-GOV-040
title: Approval, Authority and Separation-of-Duties Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-040 — Approval, Authority and Separation-of-Duties Metrics

## 1. Définition
Définir les métriques Govern de demandes d’approbation, dispositions, latence, expiry, availability, delegation/escalation, authority context et SoD afin d’observer le processus d’autorité sans réduire la qualité d’une Approval à sa vitesse ni déduire l’autorité d’un rôle.

## 2. Problème utilisateur
La gouvernance peut ralentir pour de bonnes raisons, ou être rapide parce qu’un contexte manque. Mesurer uniquement la latence ou le taux d’approval peut inciter à contourner SoD, ignorer l’indisponibilité des approvers ou masquer les authority conflicts.

## 3. Objectifs
- compter Approval Requests et dispositions approve/reject/abstain ;
- mesurer latency et expirations avec contexte ;
- observer unavailable approvers, delegation et escalation ;
- mesurer SoD conflict candidates/ineligible candidates ;
- observer missing/conflicting authority context ;
- supporter N-of-M completion lorsque applicable sans fixer un algorithme final.

## 4. Non-objectifs
Ne pas évaluer la « qualité » d’un approver, créer un ranking individuel opaque, transformer Role en authority, définir quorum/RBAC final, fixer un SLO universel ou modifier Approval/Authority source.

## 5. Propriétaire
Govern / Response Metrics owns approval/authority/SoD metric semantics. GOV-1 owns source assessments/Approvals; Settings owns identity/admin configuration; Security owns authority/SoD policy; Shared owns Metrics/Reporting mechanisms.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Authority Reviewer, Security Reviewer, Approval Coordinator, Compliance Reviewer, Govern Product Lead and authorized workforce planners with privacy-limited dimensions.

## 7. Conditions d’entrée
Approval/Authority/Eligibility records exist for a scoped snapshot; identity dimensions are masked or permissioned; metric definitions and denominator rules are pinned.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Approval Requests/Approvals | CAP-GOV-011 | disposition observations | oui | snapshot | partial |
| Authority Requirements/Contexts | CAP-GOV-009 | authority observations | oui | historical records | authority metric partial |
| eligibility/SoD outcomes | CAP-GOV-010 | control observations | selon request | snapshot | unknown/NA distinguished |
| delegation/escalation | CAP-GOV-012 | routing observations | when used | snapshot | zero vs unavailable distinguished |
| emergency approvals | CAP-GOV-013 | exceptional observations | when used | snapshot | separate dimension |
| identity availability | Settings | privacy/availability dimension | when actor metrics exposed | current policy | aggregate/masked only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Approval Request / Approval | Govern | disposition/timestamps/expiry | metric read |
| Authority Context | Govern | required/satisfied/conflicted state | metric read |
| Eligibility/SoD Assessment | Govern | candidate/eligible/ineligible reason | restricted metric read |
| Principal/Role | Settings | limited identity dimensions | masked/restricted read |
| Metric definition/snapshot | Shared | version/privacy/freshness | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Approval Metric Observation | calculate/version/supersede | Govern local concept | no approval-quality claim from latency alone |
| Authority/SoD Metric Observation | calculate/version/supersede | Govern local concept | role ≠ authority preserved |
| source Approval/Authority records | no mutation | GOV-1 | aggregate only |

## 11. Fonctionnalités
Aggregate request/disposition counts; calculate approval latency and expiry counts; segment by action class/tenant/environment/authority family; observe unavailable candidates; delegation/escalation use; SoD conflict and ineligibility reasons; missing/conflicting authority; N-of-M progress when present; privacy-limit actor-level dimensions; compare periods without causal assertion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter metric | reviewer | observation | 0 | permission | scoped values | non |
| calculate/reconcile | reviewer/system | observation | 1 | pinned definition/snapshot | value/unknown | non |
| annotate interpretation | reviewer | observation | 2 | rationale | attributed note | OPEN-013 |
| prepare privacy-safe comparison | reviewer | observation set | 2 | dimension permission | comparison | OPEN-013 |
| change approver/authority | none | source | 3/4 | prohibited | no action | GOV-1/Settings/Security only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate dispositions | oui | oui | oui | explain | table/pivot |
| calculate latency/expiry | oui | timestamps/rules | oui | summary | distribution table |
| identify repeated SoD conflict patterns | oui | grouped observations | oui | candidate explanation | grouped table |
| draft interpretation | oui | templates | oui | sourced draft | analyst note |
| rank/judge approvers automatically | no | no | no | prohibited | privacy-aware process review |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `identity-masked`, `insufficient-data`, `privacy-limited`, `definition-changed`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains period/dimensions ; Empty distinguishes zero from unavailable ; Partial names missing authority/identity sources ; Error preserves prior snapshot ; Offline read-only ; Permission denied removes sensitive dimensions ; Stale shows source/definition age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Approval Metric Observation | metric | CAP-GOV-046/Reporting | denominator/context visible |
| Authority/SoD observations | metric | control review/CAP-GOV-046 | authority not inferred from Role |
| process-friction candidate | derived signal | CAP-GOV-047 | candidate, not blame/quality conclusion |
| privacy-safe comparison | comparison | reviewer | suppressed dimensions documented |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-009..013 | snapshot selected | CAP-GOV-040 | authority/eligibility/approval events | source records |
| CAP-GOV-040 | trend review | CAP-GOV-046 | metric refs/dimensions | Response Metrics |
| CAP-GOV-040 | improvement candidate | CAP-GOV-047 | sourced observation | Response Metrics |
| CAP-GOV-040 | report/export | Shared Reporting/Export | snapshot/privacy/classification | Response Metrics |

## 18. Dépendances
CAP-GOV-009..013/034/046/047, Shared Metrics/Reporting, Settings identities, Security SoD/authority/privacy, OPEN-007/013.

## 19. Source de vérité
GOV-1 Approval/Authority/Eligibility records remain source; Govern owns domain metric definition/observation. Settings identity data and Security policy remain source-owned. Metric observation never grants/revokes authority.

## 20. Provenance et audit
Record metric definition/version, source snapshot, time bounds, excluded/unknown records, privacy suppression, aggregation dimensions, reviewer notes and AI summary provenance.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; sensitive actor dimension read; authority/approval restricted metric read; cross-tenant metric read only with explicit permission; Reporting/Export preparation. No implicit access to source Approval rationale or identity details.

## 22. Limites et erreurs
Masked identities, changed authority models, incomplete eligibility records, expired/retention-limited Approvals and sparse N-of-M scenarios can limit comparability. Approval latency is not Approval quality.

## 23. Métriques
Approval request/disposition counts; latency; expired approvals; unavailable approvers; delegation/escalation use; SoD conflict candidates; ineligible candidates; missing/conflicting authority contexts; N-of-M completion where applicable. No universal target.

## 24. Classification de livraison
`defined` / `planned`; no workforce scoring, metrics engine, quorum implementation or SLO is selected.

## 25. Critères d’acceptation
**Given** average approval latency decreases, **When** the metric is reviewed, **Then** the system does not label Approval quality as improved without separate evidence.

**Given** requester/approver identity dimensions are restricted, **When** a user lacks permission, **Then** aggregate metrics remain possible while identities are masked/suppressed.

**Given** SoD conflict candidates rise, **When** trend review occurs, **Then** the observation is surfaced without automatically accusing users or changing policy.

**Given** no AI, **When** metrics are calculated, **Then** deterministic aggregation and human review remain complete.

## 26. Questions ouvertes
OPEN-007/013 remain open. Final actor privacy model, authority implementation and thresholds remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-046/047, Govern closure, Security/Settings reviews, Shared Reporting and future screens/technical implementation.