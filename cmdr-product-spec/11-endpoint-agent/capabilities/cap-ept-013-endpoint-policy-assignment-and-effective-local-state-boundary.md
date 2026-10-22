---
id: CAP-EPT-013
title: Endpoint Policy Assignment and Effective Local State Boundary
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-013 — Endpoint Policy Assignment and Effective Local State Boundary

## 1. Définition
Define Endpoint-side observation of an assigned Platform Settings-owned Endpoint Policy and its effective local state, including received/applied/mismatch/stale/refused/degraded outcomes without creating a second Policy object or policy administration function.

## 2. Problème utilisateur
A policy can be assigned administratively but not yet received or applied locally; treating assignment as effective state can hide mismatch, stale configuration or local refusal.

## 3. Objectifs
Expose assigned policy ref, effective version, received/applied/mismatch/stale/refused/degraded effect, last apply/check, reason, freshness and provenance.

## 4. Non-objectifs
No Policy creation/edit/assignment, Fleet selector, rollout/rollback administration, configuration schema, policy engine implementation, final RBAC or screen design.

## 5. Propriétaire
Platform Settings owns Endpoint Policy and assignment. Endpoint Agent owns only the effective local observation/state needed for technical execution/diagnostics.

## 6. Utilisateurs
Platform Administrator, Endpoint Operator, Security/Audit reviewer, Investigate/SOC analyst and Govern consumer.

## 7. Conditions d’entrée
Agent identity, tenant/environment scope, assigned policy reference when one exists and local receipt/application observations.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| assigned Endpoint Policy ref/version | Platform Settings | administrative projection | non | current assignment | no assignment inferred |
| local received/applied observation | Endpoint local state | effective-state fact | oui for effectiveness | check freshness | state unknown/not-received |
| tenant/environment | Platform Settings + CAP-EPT-003 | scope | oui | current | mismatch/blocked |
| health/dependency context | CAP-EPT-008 | effect context | non | health freshness | no health inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-policy | Platform Settings | assigned ref/version only | Settings read projection |
| endpoint-agent | Endpoint Agent | local effective state | Endpoint read/manage |
| endpoint-agent-fleet | Platform Settings | assignment context only | Settings read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Effective Policy State | derive/refresh | Endpoint Agent | technical projection; never second Endpoint Policy |
| endpoint-policy | none/read only | Platform Settings | no assignment/admin mutation |

## 11. Fonctionnalités
Assigned ref/version; received; applied; mismatch; stale; rejected/refused; degraded effect; last apply/check; reason; provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect effective policy state | authorized user | endpoint-agent/effective state | 0 | read | local vs assigned context | non |
| compare assigned vs effective | Endpoint Operator | policy ref + local state | 1 | both refs | match/mismatch/stale | non |
| acknowledge/refresh local config state | authorized operator | effective state | 2 | source supports + manage | local acknowledgement/check only | OPEN-013; no policy admin |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/compare states | oui | oui | oui | sourced explanation | deterministic version/state comparison |
| summarize refusal/degradation | oui | oui for source reasons | oui | oui, attributed | structured reason list |
| modify/assign policy | Settings only | external | non via Endpoint | interdit ici | Settings-owned path |

## 14. États fonctionnels
`not-assigned`, `assigned-not-received`, `received-not-applied`, `applied`, `mismatch`, `stale`, `refused`, `degraded-effect`, `unknown`.

## 15. États d’interface
Partial distinguishes assigned/received/applied; Stale shows last check/version; Permission denied masks policy content while preserving bounded state; Offline retains last-known effective state with freshness. No Endpoint Screen ID.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Effective Policy State Projection | Endpoint technical fact | Settings Fleet/Policy, Endpoint consumers | assignment vs application explicit |
| Policy-State Provenance | local-audit/provenance context | Audit | source ref/version/reason retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings assignment | observed locally | assigned/received state | policy ref/version/scope | Settings owner retained |
| received | local application observation | applied/refused/mismatch | effective version/reason/time | no admin mutation |
| effective state | consumer projection | Settings/products | state/freshness/limits | Endpoint context retained |

## 18. Dépendances
CAP-EPT-001/003/005/008, Settings Endpoint Policies/policy assignment/Fleet, canonical endpoint-policy, Security and OPEN-008/013.

## 19. Source de vérité
Settings is source of truth for Endpoint Policy/assignment; Endpoint is source for local effective state only. Local state never becomes the Settings Policy object.

## 20. Provenance et audit
Record Agent/scope, assigned ref/version, effective version/state, receive/apply/check times, reason, source, actor/engine and correlation id; no raw secret.

## 21. Permissions fonctionnelles
Endpoint read/manage for local state; Settings endpoint-policy read/manage remains external. EPT-1 creates no atomic policy permission or cross-tenant access.

## 22. Limites et erreurs
Assigned != applied; applied != healthy automatically; effective local config != Settings config object; refused/mismatch/stale stay explicit; no second Policy.

## 23. Métriques
Assigned-not-received, not-applied, mismatch, stale, refused/degraded counts; descriptive last-check age; permission denials.

## 24. Classification de livraison
`draft / defined / planned`; no policy engine, assignment service, rollout/rollback implementation, protocol or final RBAC.

## 25. Critères d’acceptation
**Given** Settings assigns a policy but the Agent has not applied it, **When** state is read, **Then** assigned-not-applied is explicit.

**Given** effective policy version differs from assignment, **When** compared, **Then** mismatch is reported without modifying Settings.

**Given** local policy state is stale, **When** projected, **Then** stale remains visible and is not converted to applied/current.

## 26. Questions ouvertes
OPEN-008/013 remain open. Policy administration/default class-2 governance and detailed policy execution belong to source owners/future work.

## 27. Consommateurs documentaires
Settings Endpoint Policies/Fleet, Endpoint health/capability context, Security/Govern/Investigate, registers, Quality and Roadmap.