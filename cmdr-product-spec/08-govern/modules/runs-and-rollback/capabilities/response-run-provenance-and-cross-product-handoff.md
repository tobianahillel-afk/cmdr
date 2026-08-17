---
id: CAP-GOV-033
title: Response Run Provenance and Cross-Product Handoff
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-033 — Response Run Provenance and Cross-Product Handoff

## 1. Définition
Retracer de bout en bout la chaîne GOV-2 depuis Decision/Execution Handoff jusqu’au Result, y compris Playbook/version, Execution Plan, targets/readiness, authorization reconciliation, Response Run/Steps, Studio/Tool/Endpoint projections, runtime, errors/retries, verification, rollback/recovery et human decisions, puis préparer des handoffs explicites vers les produits propriétaires sans mutation silencieuse.

## 2. Problème utilisateur
Après une réponse complexe, les informations sont réparties entre Govern, Studio, Endpoint, Settings, Command et Investigate. Sans provenance corrélée, il devient difficile d’expliquer quelle Decision/version autorisait quel effet, quel executor a réellement agi, ce qui a été vérifié et pourquoi un Result a été classifié ainsi.

## 3. Objectifs
- conserver une chaîne navigable Decision→Result ;
- pinner versions et identities de tous les objets/références ;
- relier raw technical runs/calls/commands sans transfert d’ownership ;
- conserver erreurs/retries/partial/compensation/rollback/recovery and human decisions ;
- préparer handoffs permission-aware vers Command, Investigate, Detection Engineering, TI, Settings, Studio and future GOV-3 ;
- préserver return origin et destination ownership ;
- rendre la chaîne exploitable pour future Audit Trail/Response Metrics sans les spécifier ici.

## 4. Non-objectifs
Ne pas créer un Audit Trail GOV-3, calculer Response Metrics GOV-3, modifier Incident/Case/Finding/Evidence/Workflow/Settings objects, effectuer un partage externe sous OPEN-019 sans gouvernance, fusionner technical run et Response Run, supprimer/compacter la trace ou finaliser un modèle physique de provenance.

## 5. Propriétaire
Govern owns GOV-2 provenance semantics and outgoing response handoff packages. Each linked product retains its canonical objects. Shared Trace/Activity/Linking provide generic mechanisms; GOV-3 later defines Govern Audit Trail and Response Metrics capabilities.

## 6. Utilisateurs
Principal : Response Operator / Govern Reviewer. Secondaires : Incident Commander, Investigator, Detection Engineer, Threat Intelligence Analyst, Studio Operator, Platform/Endpoint Operator, Decision Maker, Auditor.

## 7. Conditions d’entrée
A Response Run and associated GOV-2 records exist; a Result may be draft/finalized or the Run may require an earlier handoff for investigation/platform/workflow improvement; caller has permission to view referenced context and prepare the target handoff.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision + GOV-1 Execution Handoff | CAP-GOV-015/016 | authority/provenance origin | oui | historical exact versions | chain incomplete |
| Playbook/Plan/Readiness/Reconciliation | CAP-GOV-017..021 | preparation lineage | oui for executed Run | exact versions | gap explicit |
| Response Run/Steps/controls | CAP-GOV-022..024 | governed execution identity | oui | Run history | no complete run chain |
| technical executor refs | CAP-GOV-025/026 | raw execution provenance | when technical execution occurred | exact source refs | source gap explicit |
| errors/retries/verification/rollback/recovery | CAP-GOV-027..031 | execution/outcome lineage | when applicable | exact record versions | omission not inferred |
| Result/follow-up needs | CAP-GOV-032 | canonical outcome | when available | exact Result version | in-progress handoff marked |
| destination context/access | destination owner + Security | handoff permission/return path | oui for handoff | current | handoff blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision/Response Run/Result and GOV-2 records | Govern | complete governed lineage | read/link/export under permission |
| Incident | Command | source/return/follow-up target | read/link only |
| Case/Finding/Evidence | Investigate | source/follow-up target | read/link, no mutation/requalification |
| Workflow/Tool/Tool Call/Automation Run/Human Gate | CMDR Studio | executor/automation provenance | restricted read/link |
| Endpoint/Agent Command/technical result | Endpoint Agent | technical provenance | restricted read/link |
| Settings Integration/Environment/Secret Reference metadata | Platform Settings | runtime/config context | restricted read; never raw secret |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| GOV-2 Provenance Chain | build/version/append/reference/export | Govern local concept + Shared mechanisms | source records immutable/referenced |
| Cross-Product Handoff Package | create/update/withdraw/supersede | Govern local concept | destination/permission/return origin explicit |
| follow-up relation | create/link/close | Govern/Shared Linking | no destination-object mutation |
| destination objects | no silent create/update | Command/Investigate/Studio/Settings/etc. | destination owner accepts/acts explicitly |

## 11. Fonctionnalités
Assemble full chain; validate missing links/version gaps; navigate from Result back to Decision and sources; include exact target/action/scope/authority; correlate technical runs/calls/commands; retain error/retry/verification/rollback/recovery records; mask restricted data/secret metadata appropriately; prepare destination-specific handoff context; route Result to Command; prepare Investigate follow-up for Case/Finding context; prepare Detection/TI feedback; prepare Settings issue; prepare Studio workflow-improvement context; expose future GOV-3 audit/metrics inputs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/navigate provenance | authorized user | Provenance Chain | 0 | read | linked trace |
| validate lineage completeness | reviewer | chain | 1 | source refs | gaps/contradictions list | non |
| prepare/update/withdraw handoff | Response Operator | Handoff Package | 2 | destination/permission/context | versioned package | OPEN-013 |
| export provenance | authorized user | chain/package | 1/2 by sensitivity | export permission/masking | bounded export | governed |
| mutate destination/source object silently | none | external object | — | forbidden | no effect | destination owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate exact ids/versions | oui | oui | oui | explain gaps | lineage graph/table |
| detect missing provenance links | oui | oui | oui | summary | completeness checklist |
| draft destination handoff summary | oui | templates/field mapping | oui | sourced draft | structured package |
| recommend follow-up destination | oui | explicit rules | oui | candidate only | destination matrix |
| share/mutate/authorize externally | accountable destination/governance path | no autonomous | no silent | prohibited | explicit handoff/permission |

## 14. États fonctionnels
Provenance chain: `building`, `complete-for-known-sources`, `partial`, `restricted`, `contradictory`, `source-unavailable`, `superseded`. Handoff: `draft`, `ready`, `permission-blocked`, `sent-to-destination-workflow`, `accepted`, `returned`, `withdrawn`, `superseded`, `closed`. `complete-for-known-sources` is not proof that every external fact exists.

## 15. États d’interface
Loading preserves chain position ; Empty requires a Run/source object ; Partial names missing/restricted links ; Error retains valid graph/table context ; Offline shows last known trace read-only ; Permission denied masks protected nodes/edges ; Stale identifies late source updates that may supersede handoff context.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| GOV-2 Provenance Chain | trace relation set | GOV-3 future/Auditor/Run users | Decision→Result links attributable/versioned |
| Result handoff to Command | bounded package | Command Incident/coordination | no Incident mutation without Command action |
| investigation/finding follow-up package | bounded package | Investigate | no Evidence/Finding requalification |
| Detection/TI feedback package | bounded package | Investigate Detection/TI | no rule/intel lifecycle mutation automatically |
| Settings/Studio improvement package | bounded package | Platform Settings / Studio | no config/Workflow mutation automatically |
| GOV-3 input refs | provenance/metric source set | future Audit Trail/Response Metrics | GOV-3 capabilities not created here |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-016..032 | records/events created | CAP-GOV-033 | exact ids/versions/relations/restrictions | originating capability |
| CAP-GOV-032 | Result finalized/follow-up needed | Command | Result/action/targets/outcome/residual risk/Run refs | Govern Result |
| CAP-GOV-029/032 | investigation follow-up needed | Investigate | anomaly/gap/Run/Result/source refs/return origin | Govern Run/Result |
| CAP-GOV-032/033 | detection/TI feedback | Investigate Detection/TI | verified outcome/technical observations/limits | Govern |
| CAP-GOV-025..027/033 | runtime/config/workflow issue | Settings/Studio/Endpoint owner | exact source refs/error/context | Govern Run |
| GOV-3 future | audit/metrics consumption | GOV-3 | immutable/ref-resolvable GOV-2 chain | source Run/Result |

## 18. Dépendances
CAP-GOV-001..032, Command Incident, Investigate Case/Finding/Evidence/Detection/TI, Studio Workflow/Tool/Automation Run/Human Gate, Endpoint Agent, Settings integrations/secrets/env, Shared Trace/Activity/Linking/Versioning/Reporting/Export, Security permissions, OPEN-007/008/013/015/019.

## 19. Source de vérité
Each product remains source of its canonical records. Govern is source of Response Run/Result and the cross-record GOV-2 provenance/handoff composition. Shared Trace may persist generic correlation; it does not become the owner of Govern lifecycle semantics.

## 20. Provenance et audit
The provenance record itself retains source product/object/id/version, relation type, target/action/scope, Decision/Approval/Exception context, Plan/Playbook, executor owner/runtime refs, timestamps, human/automation actions, errors/retries, verification/rollback/recovery, Result classification, restrictions/masking, destination/return origin and supersession. No raw secret material is recorded.

## 21. Permissions fonctionnelles
Provenance read/export; restricted context read; cross-product handoff prepare/send to destination workflow; Result handoff; follow-up package; masking; automated summary request. Destination-object write permissions remain destination-owned; cross-tenant/external sharing requires explicit authority and OPEN-019 remains unresolved.

## 22. Limites et erreurs
Missing source refs, deleted/inaccessible external records, permission changes, late runtime results, broken correlation, cross-tenant boundary or destination refusal leave explicit partial/restricted/returned state. A handoff never silently compensates by copying protected content or mutating destination records.

## 23. Métriques
Runs with Decision→Result navigable lineage; missing/late source links; handoff acceptance/return rates; provenance export masking incidents; destination object mutations without explicit owner transition — target zero; GOV-3 metrics are not defined here.

## 24. Classification de livraison
`defined` / `planned`; no physical audit store, graph database, event protocol, external-sharing integration, API or GOV-3 implementation selected.

## 25. Critères d’acceptation
**Given** a Response Run used a Studio Automation Run and Endpoint technical command, **When** provenance is inspected, **Then** both source identities/owners and their links to the Govern Run remain distinct and navigable.

**Given** a Result requires Investigate follow-up, **When** a handoff is prepared, **Then** the package references Case/Finding/Evidence context without changing or requalifying those objects.

**Given** a Result is handed to Command, **When** Command receives it, **Then** Govern remains owner of Result and Command decides its own Incident transition; no silent update occurs.

**Given** no AI, **When** provenance/handoffs are managed, **Then** deterministic id/version linking, completeness checks, structured packages and human workflows provide full functionality.

## 26. Questions ouvertes
OPEN-007/008/013/015 remain relevant; OPEN-019 remains open for external/cross-tenant dissemination. GOV-2 creates or closes no OPEN. Final physical provenance/audit/metrics architecture belongs later phases/GOV-3/Technique.

## 27. Consommateurs documentaires
Command, Investigate, Detection Engineering, Threat Intelligence, Platform Settings, Studio, Endpoint Agent, Shared, future GOV-3 Audit Trail/Response Metrics, Objects/Permissions/Screens/Journeys/Technique, quality reports and roadmap closure.
