---
id: CAP-EPT-062
title: Technical Execution Output, Error, Timeout and Cancellation Semantics
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-062 — Technical Execution Output, Error, Timeout and Cancellation Semantics

## 1. Définition
Normaliser fonctionnellement les faits target-side d’une exécution technique : output stream references conceptuelles, structured/unstructured technical output, warnings, exit/completion status, errors, timeout, partial output, cancel/stop requests, lost status et unknown target-side termination.

## 2. Problème utilisateur
Un output vide, une sortie partielle, un timeout ou un cancel request peut être confondu avec succès, échec confirmé ou arrêt réel de la cible.

## 3. Objectifs
Séparer state et output ; conserver stdout/stderr-like refs sans physical schema ; distinguer request/confirmation ; préserver partial/truncation/redaction ; fournir input fiable à Investigate et Govern reconciliation sans créer Result/Evidence.

## 4. Non-objectifs
Aucun output schema physique, storage, log protocol, Govern Result, Evidence qualification, success business verdict, rollback ou containment verification.

## 5. Propriétaire
Endpoint owns raw/local technical output and execution-state facts. Investigate owns review/qualification; Govern owns canonical Result; Studio owns Tool Call/Automation Run outcomes.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Case Analyst, Govern Reviewer, Auditor, Security Reviewer.

## 7. Conditions d’entrée
Technical execution/file/script attempt identifiable, target/session refs, source status/output events or explicit gaps, output permissions and masking rules.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| execution attempt | CAP-EPT-059/060/061 | source operation | oui | exact version | orphaned output |
| status events | Endpoint runtime | target-side state | conditionnel | event time | status-unknown |
| output/warning/error refs | Endpoint runtime | technical content | conditionnel | event time | no-output/unknown |
| cancel/stop intent | CAP-EPT-057/059..061/063 | control request | non | request time | no control state |
| masking/classification | Security/Policy | visibility | oui for sensitive output | current | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Execution Attempt | Endpoint | origin/state | read |
| Technical Session | Endpoint | target/operator | read |
| Operation Result business record | Investigate | consumer correlation | read ref |
| Response Run/Result | Govern | reconciliation destination | read ref only |
| Tool Call/Automation Run | Studio | provenance | read ref only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Execution Output | create/append/version | Endpoint | raw/local technical semantics |
| Execution Error/Warning | append | Endpoint | source/status attributed |
| Target-Termination Knowledge | reconcile | Endpoint | confirmed/unknown explicit |

## 11. Fonctionnalités
Associate outputs/errors to exact attempt; preserve ordered/unordered stream refs conceptually; record warning/exit/completion facts; mark partial/truncated/redacted; reconcile timeout/cancel/stop with acknowledgements; expose termination-unknown and late status.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect output/errors | authorized user | technical output | 0 | read | sourced output | non |
| reconcile state deterministically | service/reviewer | attempt | 1 | events | state/reasons | non |
| request cancel/stop | authorized operator | attempt | 2 | interruptible | request recorded | according original action |
| mark canonical Result | aucun Endpoint | output | — | Govern only | no mutation | Govern owner |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| group output/errors | oui | oui | oui | summarize | raw streams/events |
| reconcile timeout/cancel | oui | oui | oui | explain | status rules |
| summarize technical output | oui | oui | oui | oui, attributed | structured/raw output |
| invent missing output/termination/result | non | interdit | non | interdit | explicit unknown/gap |

## 14. États fonctionnels
`output-pending`, `output-present`, `partial`, `truncated`, `redacted`, `completed`, `failed`, `timed-out`, `cancel-requested`, `cancelled-confirmed`, `stop-requested`, `termination-confirmed`, `termination-unknown`, `status-lost`, `late-status`.

## 15. États d’interface
No Screen ID. `timeout != confirmed termination`, `cancel-requested != cancelled`, `stop-requested != stopped`; permission/redaction states do not fabricate empty output.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Execution Output | Endpoint technical output | Investigate/Govern/Studio caller | not Govern Result/Evidence |
| execution state/error | Endpoint fact | CAP-EPT-063/064 | request/confirmation distinction |
| termination knowledge | diagnostic | operator/Govern | unknown stays unknown |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-059..061 | runtime events | CAP-EPT-062 | attempt/status/output/error | source attempt |
| CAP-EPT-062 | review/handoff | Investigate | neutral technical output + limitations | output ref |
| CAP-EPT-062 | governed execution reconciliation | Govern | technical status/output refs | Response Run remains Govern |

## 18. Dépendances
CAP-EPT-057..061/063/064, Investigate CAP-INV-212, Govern execution boundaries, Studio Tool Call lifecycle, Security masking, OPEN-013/015.

## 19. Source de vérité
Endpoint SOT of target-side technical output/status. Govern SOT canonical Result; Investigate SOT analytical review; Studio SOT Tool Call/Automation Run lifecycle.

## 20. Provenance et audit
Attempt/session/request refs, target/operator, source event times, output/error refs, masking, timeout/cancel/stop intents and confirmations, late/lost state, authority/correlation IDs.

## 21. Permissions fonctionnelles
Technical output read, sensitive output read, error/provenance read, cancel/stop request, cross-tenant deny; no Result/Evidence permission inherited.

## 22. Limites et erreurs
Command output ≠ Govern Result/Evidence; command success ≠ Response Run success; timeout ≠ target termination; cancellation ≠ rollback; lost status remains unknown.

## 23. Métriques
Output partial/truncated/redacted, timeout/cancel/termination-unknown, late-status, handoff completeness and false business-result promotion target zero.

## 24. Classification de livraison
`draft / defined / planned`; no output transport/storage/schema or Result engine.

## 25. Critères d’acceptation
**Given** command accepted but execution fails, **When** output is reviewed, **Then** error and any partial output remain linked to the attempt without success inference.

**Given** a command times out, **When** target termination is unconfirmed, **Then** timeout and `termination-unknown` are both retained.

**Given** technical output is passed to Govern, **When** reconciliation occurs, **Then** Endpoint output remains source data and Govern creates/updates its Result separately.

## 26. Questions ouvertes
OPEN-013/015 remain open; physical output schema and cross-product bridge are not finalized.

## 27. Consommateurs documentaires
Investigate Result Handling, Govern Runs/Result, Studio provenance, CAP-EPT-063/064, Security/Quality.
