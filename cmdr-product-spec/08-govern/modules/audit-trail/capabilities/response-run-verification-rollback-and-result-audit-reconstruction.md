---
id: CAP-GOV-036
title: Response Run, Verification, Rollback and Result Audit Reconstruction
product: govern
module: audit-trail
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-036 — Response Run, Verification, Rollback and Result Audit Reconstruction

## 1. Définition
Reconstruire historiquement `Decision → Execution Handoff → Playbook/version → Execution Plan → Readiness → Response Run/Steps → technical executor refs → Verification → Rollback/Recovery → Result` en conservant exact scope/targets, erreurs, retries, partial success, residual risk et provenance sans rejouer l’exécution.

## 2. Problème utilisateur
Une exécution gouvernée peut traverser plusieurs executors, retries et états partiels. Sans reconstruction séparée, un raw technical success peut être pris pour un Result, une compensation pour un rollback ou une verification absente pour un succès prouvé.

## 3. Objectifs
- pinner Decision, Playbook, Plan, Run et Result versions ;
- reconstruire targets/scope et runtime owner refs ;
- exposer per-step/per-target progress, errors/retries/partial outcomes ;
- relier Verification Assessment, rollback/recovery et residual risk ;
- distinguer raw outputs du canonical Result ;
- rendre gaps/contradictions reviewables.

## 4. Non-objectifs
Ne pas redémarrer un Run, réexécuter un Tool/Workflow, recalculer un Result comme vérité automatique, lancer rollback, modifier Decision/Evidence/Finding, inventer un technical event ou choisir un executor/protocol.

## 5. Propriétaire
Govern / Audit Trail owns reconstruction semantics. GOV-2 remains owner of Response Run, verification, rollback/recovery governance and Result. Studio/Endpoint/providers own technical runs/outputs; Shared owns Trace/Activity.

## 6. Utilisateurs
Principal : Govern Auditor. Secondaires : Response Operator, Verification Reviewer, Rollback Reviewer, Incident Commander, Investigator, Studio/Endpoint Operator, Security/Compliance Reviewer.

## 7. Conditions d’entrée
A Decision, Response Run or Result anchor exists; source refs are readable or marked inaccessible; target/scope and historical versions are available enough to establish a chain or explicit gap.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision/Handoff | GOV-1 | authority lineage | oui | pinned historical | chain partial |
| Playbook/Plan/Readiness | CAP-GOV-017..021 | execution preparation | oui where Run exists | pinned versions | preparation gap |
| Run/Step history | CAP-GOV-022..027 | governed execution | oui | historical sequence | run gap |
| technical executor refs | Studio/Endpoint/provider | raw provenance | when execution occurred | correlated records | source-unavailable |
| Verification/Residual Risk | CAP-GOV-028/029 | outcome verification | when required | historical version | verification gap |
| rollback/recovery | CAP-GOV-030/031 | recovery lineage | when applicable | historical | rollback gap/NA explicit |
| Result | CAP-GOV-032 | canonical outcome | expected for closed run | exact version | Result gap |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Execution Handoff | Govern | authorized intent | audit read |
| Playbook / Execution Plan | Govern | versions/steps/bounds | audit read |
| Response Run / Step / Rollback / Result | Govern | lifecycle/outcomes | audit read |
| Workflow/Automation Run/Tool Call | Studio | execution refs | restricted link/read |
| technical response/local audit | Endpoint/provider | raw refs/status | restricted link/read |
| Govern Audit Event | Govern | sequence | correlate/read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run Audit Reconstruction | create/update/version | Govern local concept | no source mutation/re-execution |
| execution gap/contradiction candidate | create/review | Govern local concept | candidate only |
| source Run/Result/technical records | no mutation | source owner | read/link only |

## 11. Fonctionnalités
Traverse Decision-to-Run refs; verify exact Playbook/Plan/target/scope lineage; assemble Run controls and Steps; correlate Automation Run/Tool Call/Endpoint refs; preserve runtime errors/retries/partial states; align Verification and residual risk; reconstruct rollback/recovery; link Result and classification; compare Runs; flag missing/contradictory nodes; return to source.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect Run reconstruction | auditor | reconstruction | 0 | read | sourced chain | non |
| compare Run/Result versions | auditor | reconstruction | 1 | comparable refs | deterministic diff | non |
| flag missing/contradictory runtime link | auditor | audit candidate | 2 | rationale | versioned candidate | OPEN-013 |
| annotate review | auditor | reconstruction | 2 | review permission | attributed note | OPEN-013 |
| retry/start/rollback | none | Run/target | 3/4 | prohibited in audit | no effect | GOV-2 only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| traverse Run/Step IDs | oui | oui | oui | no authority | relation graph/table |
| reconcile per-target outcomes | oui | structured aggregation | oui | summary | target matrix |
| detect Result-without-verification | oui | expected-link rules | oui | candidate explanation | completeness rules |
| summarize errors/retries/rollback | oui | source aggregation | oui | sourced summary | event table |
| infer unrecorded execution | no | no | no | prohibited | explicit unknown |

## 14. États fonctionnels
`building`, `complete-enough-for-review`, `partial`, `preparation-gap`, `run-gap`, `technical-source-unavailable`, `verification-gap`, `rollback-gap`, `result-gap`, `contradictory`, `retention-limited`, `disputed`, `reviewed`, `superseded`.

## 15. États d’interface
Loading preserves anchor ; Empty means no visible reconstruction ; Partial names missing stages ; Error keeps latest valid chain ; Offline read-only ; Permission denied masks technical details ; Stale exposes later-arriving source records.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run Audit Reconstruction | audit review object | CAP-GOV-037/038/042..044 | exact lineage and limitations |
| technical-vs-canonical comparison | review artifact | auditor | raw output remains distinct from Result |
| gap/contradiction candidates | audit assessment inputs | CAP-GOV-037 | no automatic falsity claim |
| reconstruction summary | review output | CAP-GOV-038 | source-linked |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-034 | run-chain events available | CAP-GOV-036 | refs/versions/sources | Audit Trail |
| CAP-GOV-036 | gaps/contradictions | CAP-GOV-037 | expected/observed nodes | reconstruction |
| CAP-GOV-036 | review/package request | CAP-GOV-038 | scoped chain/restrictions | reconstruction |
| CAP-GOV-036 | metric derivation | CAP-GOV-042..044 | sourced observations | audit chain |

## 18. Dépendances
CAP-GOV-015..034/037/038, Studio Workflow/Automation Run/Tool, Endpoint technical records, Settings runtime/retention sources, Shared Trace/Activity/Versioning, OPEN-008/013/015.

## 19. Source de vérité
GOV-2 objects remain source for Govern execution semantics; technical owners remain source for raw technical records. Reconstruction is derivative audit context and never rewrites Decision, Run, Result, Evidence or Finding.

## 20. Provenance et audit
Record reconstruction version, all object/run/tool refs, exact versions, target/scope, timestamps, executor owners, technical source availability, errors/retries/partiality, verification, residual risk, rollback/recovery, Result and reviewer/automation provenance.

## 21. Permissions fonctionnelles
Audit read; restricted technical-output read; reconstruction create/read; compare Runs/Results; contradiction review; provenance export preparation. Executor permissions are never inherited from audit access.

## 22. Limites et erreurs
Lost executor status, retention gaps, late events, unknown target effect, inconsistent raw outputs, missing verification or Result and restricted data remain explicit. Technical success is not verified success; reconstruction does not make causal claims from timestamp order alone.

## 23. Métriques
Run reconstructions complete-enough; missing Result/verification candidates; technical-source gaps; retry/partial/rollback chain completeness; technical-success vs Result mismatch counts. No quality conclusion from count alone.

## 24. Classification de livraison
`defined` / `planned`; no replay engine, audit store, runtime adapter or verification implementation is delivered.

## 25. Critères d’acceptation
**Given** a Run exists with no linked Decision, **When** reconstruction opens it, **Then** the Run remains visible and a Decision-link gap is reported without inventing authority.

**Given** technical output exists but no canonical Result, **When** reconstruction is reviewed, **Then** raw output remains a source reference and Result absence is explicit.

**Given** Result exists but Verification is required and absent, **When** the chain is reconstructed, **Then** verification gap is shown and Result is not silently reclassified.

**Given** no AI, **When** reconstruction runs, **Then** ID traversal, matrices and human review remain sufficient.

## 26. Questions ouvertes
OPEN-008/013/015 remain open. Final event-retention/runtime bridge implementation stays future; no new OPEN.

## 27. Consommateurs documentaires
Audit Trail, CAP-GOV-037/038/042..047, Govern closure, future screens/objects/permissions/technical audit implementation and quality reports.