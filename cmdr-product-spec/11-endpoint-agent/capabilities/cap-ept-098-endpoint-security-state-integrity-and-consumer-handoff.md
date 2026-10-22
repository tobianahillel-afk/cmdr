---
id: CAP-EPT-098
title: Endpoint Security-State, Integrity and Consumer Handoff
product: endpoint-agent
module: security
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-012, REQ-PROD-014, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-098 — Endpoint Security-State, Integrity and Consumer Handoff

## 1. Définition
Consolider les observations locales de protection, intégrité et sécurité Endpoint et les transmettre aux consommateurs autorisés sans les promouvoir automatiquement en Finding, Incident, Decision ou Result.

## 2. Problème utilisateur
Des faits de self-protection, privilège, tamper ou intégrité peuvent être surinterprétés comme compromission ou conclusion métier lorsqu’ils quittent Endpoint.

## 3. Objectifs
Represent protection/integrity/security-state observations, component state, suspicious/tamper candidate, degraded/unknown state, health/capability effects, consumer handoff and provenance.

## 4. Non-objectifs
Aucun security score global, malware verdict, Finding/Incident creation, Decision/Result, SIEM event model, integrity engine, crypto proof or implementation.

## 5. Propriétaire
Endpoint owns local security-state observations. Investigate owns Findings/Evidence, Command owns Incident, Govern owns Decision/Result, Settings owns admin security configuration.

## 6. Utilisateurs
Endpoint Operator, Security Reviewer, Investigator, Incident Commander, Platform Administrator, Govern Reviewer, Auditor.

## 7. Conditions d’entrée
Source-backed CAP-EPT-094..097 facts, Agent/tenant/environment, health/capability context, consumer permissions and freshness.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| self-protection/tamper state | CAP-EPT-094 | security observation | non | fresh | partial state |
| privilege/component context | CAP-EPT-095 | security context | non | fresh | partial state |
| sensitive-material state | CAP-EPT-096 | handling context | non | fresh | no secret inference |
| local audit/provenance | CAP-EPT-097 | provenance | oui for handoff | source time | handoff limited |
| health/capability state | CAP-EPT-008/093 | effect context | non | fresh | impact unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint security states | Endpoint | sourced facts | read |
| Agent Health/Capability | Endpoint | technical impact | read |
| Endpoint Policy/config | Settings | context only | read projection |
| Case/Finding/Incident/Result | external owners | destination refs only | no mutation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Security-State Projection | derive/version | Endpoint | technical observations only |
| Security-State Handoff | create/link | Endpoint | destination owner retained |
| Finding/Incident/Result | aucune création/mutation | Investigate/Command/Govern | external workflow only |

## 11. Fonctionnalités
Aggregate source-backed local security observations without hiding uncertainty; link component/health/capability impact; distinguish suspicious/tamper candidates from verdicts; apply masking; hand off typed refs to authorized consumers; preserve source owner and provenance.

## 12. Actions utilisateur
Inspect Class 0; normalize/correlate local facts Class 1; prepare consumer handoff Class 2 no-effect. Any investigation/incident/response mutation occurs in destination owner workflow.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate source-backed state | oui | oui | oui | summarize | structured state projection |
| explain impact/uncertainty | oui | oui | oui | oui, attributed | reason list |
| declare compromise/Finding/Incident/Result | non | external owner only | non | interdit | destination workflow |

## 14. États fonctionnels
`normal-for-declared-checks`, `security-degraded`, `tamper-candidate`, `suspicious-candidate`, `integrity-unknown`, `protection-unknown`, `capability-affected`, `restricted`, `stale`, `handoff-ready`, `handoff-partial`.

## 15. États d’interface
No Screen ID. Candidate, observation and destination-owned verdict/object remain visually/conceptually distinct; sensitive details may be masked.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Security-State Projection | Endpoint fact | Command/Investigate/Settings/Govern | source/limits explicit |
| Security-State Handoff | typed refs/context | authorized destination | no object promotion |
| provenance/limitations | audit context | Quality/Security | no invented facts |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-094..097 | state available | CAP-EPT-098 | local security facts/provenance | source refs retained |
| CAP-EPT-098 | analytical need | Investigate | observations/candidates/limits | no Finding auto-create |
| CAP-EPT-098 | operational need | Command/Govern/Settings | typed technical state | destination ownership retained |

## 18. Dépendances
CAP-EPT-008/011/093..097/099; Command/Investigate/Govern/Settings; Security; OPEN-008/015/017.

## 19. Source de vérité
Endpoint SOT for local security observations only; destination products remain SOT for their canonical objects and conclusions.

## 20. Provenance et audit
Agent/component/source refs, observed states, candidate rationale, health/capability effects, masking, timestamps, consumer/destination refs and handoff outcome.

## 21. Permissions fonctionnelles
Security-state read, sensitive-security-state read, handoff preparation, provenance read, destination-object permission checked by destination, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Security-state observation != Finding != Incident != Result; integrity check PASS != absence of compromise; local fact cannot create destination authority.

## 23. Métriques
Degraded/unknown/candidate states, affected capabilities, handoff completeness/denials and false automatic object-promotion target zero.

## 24. Classification de livraison
`draft / defined / planned`; no security engine, SIEM integration or implementation.

## 25. Critères d’acceptation
**Given** a tamper candidate exists, **When** handed to Investigate, **Then** it remains an Endpoint observation and does not automatically become a Finding.

**Given** security state is restricted, **When** Command consumes it, **Then** permitted summary/limits are shown without revealing protected details.

**Given** AI is unavailable, **When** security-state handoff is prepared, **Then** deterministic/manual structured facts remain available.

## 26. Questions ouvertes
OPEN-008/015/017 remain open; no platform support, runtime bridge or detection runtime decision is closed.

## 27. Consommateurs documentaires
Endpoint closure, Command, Investigate, Govern, Settings, Security, Quality, user journeys, registers and Roadmap.