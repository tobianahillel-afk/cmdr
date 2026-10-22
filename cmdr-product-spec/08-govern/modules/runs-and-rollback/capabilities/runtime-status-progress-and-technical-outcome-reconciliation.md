---
id: CAP-GOV-026
title: Runtime Status, Progress and Technical Outcome Reconciliation
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-026 — Runtime Status, Progress and Technical Outcome Reconciliation

## 1. Définition
Réconcilier dans le Response Run les réponses techniques de Studio, Endpoint et autres executors — accepted/rejected, started, progress, per-target/per-step state, raw output reference, error, timeout, completion, partial completion, disconnect and unknown state — tout en conservant raw references/provenance et sans transformer un technical outcome en canonical Result.

## 2. Problème utilisateur
Les executors exposent des états et vocabulaires différents. Une traduction naïve peut marquer `success` alors qu’un output est partiel, contradictoire, ancien ou seulement accepté. Le Run doit offrir une vue cohérente sans effacer la source technique.

## 3. Objectifs
- normaliser fonctionnellement les statuts sans perdre le raw status/ref ;
- corréler per-step/per-target/progress au Run ;
- distinguer request, acceptance, start, running, completion et unknown ;
- conserver timestamp/freshness and executor identity ;
- détecter contradictions, late updates, disconnects and stale state ;
- produire des inputs explicites pour errors/verification, pas un Result métier.

## 4. Non-objectifs
Ne pas imposer un schéma/protocole runtime, modifier la source technique, conclure le business outcome, produire automatiquement Result, cacher les contradictions, auto-retry, auto-rollback ou réécrire la Decision.

## 5. Propriétaire
Govern owns normalized Run/Step reconciliation semantics. Studio/Endpoint/provider owners remain source of raw execution state/output. Shared may transport live updates/Jobs but does not own Response Run.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Runtime/Endpoint/Studio Operator, Govern Reviewer, Incident Commander, Investigator, Auditor.

## 7. Conditions d’entrée
Existing Run/Step and correlated executor handoff from CAP-GOV-025, permission to read returned technical metadata/output refs and stable correlation identifiers.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/Step + handoff correlation | Govern | reconciliation anchor | oui | current | output cannot be attributed |
| executor identity/version | Studio/Endpoint/provider | source context | oui | execution record | unknown source/block |
| accepted/rejected/started/status | technical owner | raw state | oui when emitted | source timestamp | status-unknown |
| progress/per-target/per-step state | technical owner | technical progress | selon executor | timestamp visible | partial/unknown |
| technical output/error refs | technical owner | raw result/error reference | terminal/progress dependent | exact source ref | incomplete |
| connection/health/freshness | source/Settings | reliability context | selon executor | current/last seen | stale/disconnected |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Step | Govern | expected current state/targets | read/update reconciliation |
| Executor Handoff | Govern | correlation/intent | read |
| Automation Run/Tool Call | CMDR Studio | raw status/output refs | restricted read/link |
| Agent Command/Endpoint state | Endpoint Agent | technical status/result refs | restricted read/link |
| Integration/health | Platform Settings | source health/freshness | restricted read |
| Job/Trace | Shared | transport/progress/correlation | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Runtime Status projection | create/update/version | Govern local projection | raw source/status retained |
| Response Run/Step | reconcile confirmed state/progress | Govern | no business success inference |
| contradiction/unknown record | create/resolve by new evidence | Govern local concept | conflicting sources preserved |
| technical source records | no mutation | source owner | read/reference only |

## 11. Fonctionnalités
Correlate accepted/rejected/start; map raw source states to bounded functional states; aggregate per-target/per-step progress; preserve raw labels and timestamps; attach technical output/error refs; detect source contradictions, duplicate/late events, target disconnect, stale executor and lost status; expose unknown rather than guessing; route terminal/partial/failure inputs to CAP-GOV-027/029.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect raw/normalized status | operator | Runtime Status | 0 | read | source-backed status |
| run deterministic reconciliation | reviewer | status set | 1 | correlated refs | normalized/contradiction result | non |
| annotate/select authoritative source only where contract defines precedence | reviewer | reconciliation | 2 | rationale/rule | attributed disposition | OPEN-013 |
| request refresh/status recovery | operator | source status | 2 | source supports safe query | new observation request | OPEN-013 |
| declare canonical Result/start new effect | none here | Result/target | — | outside capability | no effect | downstream/governed path |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate source events | oui | ids/time/version rules | oui | explanation | correlation table |
| map raw to functional status | oui | explicit mapping | oui | no hidden mapping | status mapping table |
| detect contradictions/staleness | oui | oui | oui | summary | comparison/freshness rules |
| summarize progress/errors | oui | structured aggregation | oui | sourced summary | per-target/step table |
| classify verified business success | downstream human/rules | no here | no | prohibited | CAP-GOV-029/032 |

## 14. États fonctionnels
Normalized observations include `not-accepted`, `accepted`, `rejected`, `start-requested`, `started`, `running`, `progressing`, `partially-completed`, `technically-completed`, `technically-failed`, `timed-out`, `disconnected`, `stale`, `contradictory`, `status-unknown`, `cancel/stop-pending`, `terminal-awaiting-verification`, `superseded`.

## 15. États d’interface
Loading retains last confirmed status ; Empty no runtime observation ; Partial shows per-target gaps ; Error keeps raw refs and last confirmed status ; Offline distinguishes platform UI offline from executor state ; Permission denied masks protected raw output ; Stale visibly shows source timestamp and prevents fresh-success claim.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| normalized Runtime Status | Run/Step projection | CAP-GOV-022/027/029 | raw source/ref/status/time preserved |
| technical output/error refs | source links | CAP-GOV-027/029/032 | not canonical Result |
| partial/contradiction/unknown condition | runtime event | CAP-GOV-027 | affected target/step/source explicit |
| terminal-awaiting-verification | Run condition | CAP-GOV-028/029 | technical completion ≠ verified success |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-025 / executor | status/output event | CAP-GOV-026 | source ref/raw status/time/correlation | source runtime |
| CAP-GOV-026 | failure/partial/timeout | CAP-GOV-027 | Run/step/targets/raw refs/known state | reconciliation |
| CAP-GOV-026 | technical terminal state | CAP-GOV-029 | intended vs raw outcome refs/current target state | Run |
| CAP-GOV-026 | contradiction/lost state | operator/source | affected refs/question | same Run |

## 18. Dépendances
CAP-GOV-022..025/027..033, Studio Automation Run/Tool Call, Endpoint Agent/Agent Command, Settings health/integration, Shared Jobs/Trace/Activity/live updates, Security access, OPEN-008/013/015.

## 19. Source de vérité
Technical owner is source of raw execution facts. Govern Runtime Status is source of their normalized relationship to the Response Run. It never overrides raw source history and is not the canonical Result.

## 20. Provenance et audit
Record Run/step/handoff ids, executor identity/version, raw object/ref/status/output/error, source timestamp/receipt time, mapping rule/version, normalized state, per-target progress, contradictions, stale/disconnect intervals, reviewer annotations, refresh requests, automation/AI provenance and correlation ids.

## 21. Permissions fonctionnelles
Run/step status read, technical output inspect under source permission, reconciliation run/read, restricted raw-output metadata, status refresh request, provenance export. No automatic retry/start/rollback or source-object write permission.

## 22. Limites et erreurs
Missing correlation, duplicate/late/out-of-order events, source outage, stale heartbeat, contradictory terminal statuses, inaccessible output or partial target results remain explicit. Unknown never maps silently to success or failure.

## 23. Métriques
Unknown/contradictory/stale runtime duration; event correlation failure; per-target gaps; terminal technical success later failing verification; raw-to-normalized mapping errors; technical output directly stored as canonical Result — target zero.

## 24. Classification de livraison
`defined` / `planned`; no provider/runtime schema, websocket/protocol, event format or implementation selected.

## 25. Critères d’acceptation
**Given** an executor reports success but another correlated source reports a target error, **When** reconciliation occurs, **Then** the contradiction remains visible and the Run cannot silently become successful.

**Given** executor status is lost after start, **When** the Run is viewed, **Then** state becomes unknown/stale rather than assuming completion; affected targets and last confirmed observations remain visible.

**Given** no AI, **When** runtime is reconciled, **Then** deterministic correlation/status mappings and raw-reference tables provide full functionality.

## 26. Questions ouvertes
OPEN-008 actual runtime support, OPEN-013 reconciliation mutations, OPEN-015 technical-run bridge remain open. Final status ontology/precedence protocol is future technical/object work.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-027..033, Studio/Endpoint/Settings source owners, Command/Investigate projections, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
