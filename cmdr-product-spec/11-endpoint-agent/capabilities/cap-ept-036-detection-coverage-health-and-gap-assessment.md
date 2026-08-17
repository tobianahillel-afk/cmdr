---
id: CAP-EPT-036
title: Detection Coverage, Health and Gap Assessment
product: endpoint-agent
module: detection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-036 — Detection Coverage, Health and Gap Assessment

## 1. Définition
Définir une assessment locale de coverage/health/gaps de détection à partir du contenu éligible, des sources EPT-2, des capabilities/dépendances et de la santé d’évaluation, sans confondre santé runtime, endpoint health et visibilité complète.

## 2. Problème utilisateur
Un moteur ou capteur peut être sain alors que certaines détections ne disposent pas des sources requises. Les consommateurs ont besoin de raisons précises plutôt que d’un indicateur global trompeur.

## 3. Objectifs
Projeter content eligible/loaded state, required telemetry/source availability, dependency status, evaluation health, degraded coverage, gaps, stale/unsupported/platform/version limitations et provenance.

## 4. Non-objectifs
Aucune garantie de complete visibility, aucun SLO final, aucune remediation, aucun rollout/update, aucune activation de content, aucun health engine implémenté.

## 5. Propriétaire
Endpoint Agent possède l’assessment technique locale. Endpoint health EPT-1 reste distinct. Investigate Detection Engineering conserve les assessments de couverture engineering/globales ; Settings conserve runtime/source administration.

## 6. Utilisateurs
Endpoint Operator ; Detection Engineer ; SOC/Investigate Analyst ; Platform Administrator ; Security Reviewer ; Auditor ; Command consumer.

## 7. Conditions d’entrée
Content eligibility connu ; required telemetry/capabilities déclarées ; source/sensor availability et freshness ; évaluation status ; platform/version scope ; tenant/permissions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| content eligibility/requirements | CAP-EPT-031 | couverture attendue | oui | current version | coverage unknown |
| source/sensor/telemetry availability | CAP-EPT-015..028 | facts techniques | oui | source freshness | gap/degraded reason |
| evaluation outcomes/errors | CAP-EPT-032 | health technique | non | recent | no evaluation-health conclusion |
| Endpoint health | CAP-EPT-008..010 | contexte | non | current | detection health reste séparé |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Detection Eligibility Assessment | Endpoint Agent | content/dependencies | read |
| Observation Source/Sensor state | Endpoint Agent | availability/freshness | read |
| Agent Health | Endpoint Agent | contexte seulement | read |
| Detection Content/version | Investigate | expected requirements | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Detection Coverage Assessment | dériver/rafraîchir | Endpoint Agent | scope/limitations explicites |
| Detection Health Projection | dériver | Endpoint Agent | distinct d’Agent Health |
| Detection Gap | dériver | Endpoint Agent | absence de visibilité, pas absence d’attaque |

## 11. Fonctionnalités
Comparer requirements aux sources/capabilities disponibles, classer état de coverage, calculer gaps et reasons, séparer runtime/evaluation health de coverage, signaler stale/platform/version limitations et conserver la période évaluée.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect coverage/health | Analyst/Operator | assessment | 0 | read | état + reasons | non |
| calculate coverage/gaps | deterministic service | requirements/source state | 1 | facts disponibles | assessment sourcée | non |
| request bounded status refresh | Operator | source/capability status | 2 | refresh sans acquisition supporté | refresh seulement | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| comparer requirements/sources | oui | oui | oui | non nécessaire | deterministic set comparison |
| classer gap/reason | oui | oui | oui | explication possible | reason codes |
| résumer coverage | oui | oui | oui | oui | structured matrix |
| affirmer visibilité complète sans preuve | non | non | non | interdit | `unknown/partial` |

## 14. États fonctionnels
`complete-for-declared-scope`, `partial`, `degraded`, `gap`, `stale`, `unsupported`, `unknown`, `evaluation-error`. Aucun état ne signifie visibilité universelle.

## 15. États d’interface
Aucun Screen ID. Une future projection doit montrer scope/période/reasons ; healthy runtime et healthy endpoint ne peuvent être utilisés comme proxy de full coverage.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| coverage assessment | état Endpoint | Analyst/Detection Engineering | scope/requirements visibles |
| detection health projection | état Endpoint | Operator/Settings | distinct d’Endpoint health |
| gap record | concept Endpoint | Investigate/quality | gap ≠ malicious absence |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| content requirements + source state | assessment | coverage/gap | scope/reasons/time | provenance retained |
| source/version change | recompute | changed coverage | old/new dependencies | history retained |
| coverage/gap | consumer handoff | Investigate/Settings | requirements/missing refs | aucune remediation automatique |

## 18. Dépendances
CAP-EPT-008..010/015..032 ; Detection `local-detection-engine.md` ; Investigate Detection Engineering coverage/runtime health capabilities ; Settings source/runtime health ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint possède l’assessment local pour cet Agent/scope. Investigate possède l’engineering coverage assessment business ; Settings possède configuration/health admin ; EPT-1 possède Agent health facts.

## 20. Provenance et audit
Scope, content/version, required/available sources, capability states, freshness, evaluation errors, reason codes, platform/version, time window et assessor sont conservés.

## 21. Permissions fonctionnelles
Lire content requirements, source/capability health, coverage/gap details et provenance ; restricted details masqués ; cross-tenant denied. Aucun droit de remediation/rollout implicite.

## 22. Limites et erreurs
Healthy detection runtime ≠ complete coverage ; coverage ≠ complete visibility ; missing signal ≠ absence of malicious activity ; source healthy ≠ source complete ; stale state ne doit pas être traité comme current.

## 23. Métriques
Coverage-state distribution, gaps by dependency/reason, stale/unsupported counts, evaluation-error rates, time-to-reassessment conceptuel. Aucun seuil/SLO final.

## 24. Classification de livraison
`draft / defined / planned`; documentation uniquement, aucune preuve de moteur de monitoring ou remediation.

## 25. Critères d’acceptation
**Given** une source requise devient indisponible, **When** coverage est recalculée, **Then** l’assessment devient degraded/gap avec la dépendance exacte.

**Given** runtime/source health est healthy mais un content exige une source absente, **When** coverage est consultée, **Then** runtime health reste healthy tandis que coverage est incomplete.

**Given** un gap existe, **When** aucune observation malveillante n’est reçue, **Then** le gap n’est pas interprété comme preuve de benignité.

## 26. Questions ouvertes
`OPEN-008` et `OPEN-017` restent ouvertes. La définition d’un SLO de coverage global ou d’un runtime concret est hors EPT-3.

## 27. Consommateurs documentaires
EPT-3 investigation/handoff ; Investigate Detection Engineering ; Settings health/admin ; Command ; Quality ; Roadmap ; future EPT-6 monitoring comme consommateur seulement.