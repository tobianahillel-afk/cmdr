---
id: CAP-GOV-019
title: Response Execution Plan Preparation and Parameter Binding
product: govern
module: playbooks
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-019 — Response Execution Plan Preparation and Parameter Binding

## 1. Définition
Préparer un **Execution Plan** Govern à partir d’une Decision, d’un Execution Handoff Package et d’un Playbook/version compatible, en fixant steps fonctionnels, exact targets/scope, parameter definitions/bindings, Secret References, dependencies, ordering/concurrency bounds, time/retry/stop constraints, verification and rollback requirements sans exposer de secret ni démarrer l’exécution.

## 2. Problème utilisateur
Une procédure compatible peut encore être exécutée avec de mauvais paramètres, une portée trop large, un ordre incorrect ou une valeur sensible copiée dans le plan. Sans plan versionné, l’executor reçoit une intention ambiguë et le futur Run devient impossible à expliquer.

## 3. Objectifs
- lier exact Decision/Handoff/Playbook versions ;
- matérialiser l’intention d’exécution en steps fonctionnels et dépendances ;
- distinguer parameter definition, binding et Secret Reference ;
- fixer exact targets/scope, ordering, concurrency, time bounds, retry limits and stop conditions ;
- préserver rollback/verification requirements ;
- produire un plan sans effet prêt pour target/readiness review.

## 4. Non-objectifs
Ne pas résoudre une valeur de secret, administrer Settings, choisir un executor technique, exécuter une commande, démarrer Response Run, élargir le scope, modifier Decision/Playbook, définir un format/API final ou créer des steps techniques propriétaires de Studio/Endpoint.

## 5. Propriétaire
Govern / Playbooks owns Execution Plan semantics and bindings by reference. Platform Settings owns Secret References/values and runtime configuration. Studio/Endpoint/runtime owners own technical input schemas and execution semantics.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Govern Reviewer, Decision Maker, Playbook Owner, Studio/Endpoint Operator, Security Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-GOV-018 compatible Playbook version, current Decision/Handoff, declared target references, allowed scope/conditions, required parameter definitions and authorized visibility into metadata needed to bind references.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision + Handoff | CAP-GOV-015/016 | authority/bounds | oui | current planning versions | plan blocked |
| compatible Playbook/version | CAP-GOV-018 | procedure | oui | exact reviewed version | no plan |
| target references | Handoff/source owners | intended targets | oui | exact handoff refs | incomplete |
| parameter definitions | Playbook/Workflow/Tool metadata | expected inputs | selon step | exact referenced version | binding impossible |
| parameter values/non-secret refs | requester/source context | execution values | selon definition | version/source visible | incomplete |
| Secret References | Platform Settings | sensitive input reference | selon definition | metadata/current ref | blocked; raw value never requested |
| time/retry/concurrency/stop constraints | Decision/Playbook/Policy | safety bounds | selon action | current versions | explicit unknown/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Handoff Package | Govern | exact bounds/conditions | read |
| Playbook | Govern | version/steps/parameter definitions | read |
| Workflow/Tool | CMDR Studio | referenced input/capability metadata | read/link only |
| Secret Reference | Platform Settings | existence/id/metadata only | restricted read, never raw value |
| Target Reference | source owner/Shared | intended identity ref | read/link |
| Policy/Approval/Exception context | Govern | execution constraints | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Execution Plan | create/update/version/supersede | Govern local concept | exact Decision/Playbook versions pinned |
| Parameter Binding | create/update/remove | Govern local concept | definition/value source/ref explicit |
| Secret Reference binding | attach/detach reference | Govern + Settings source | reference only; no secret material |
| Decision/Playbook/Workflow/target | no source mutation | respective owner | plan does not rewrite sources |

## 11. Fonctionnalités
Assemble plan lineage; define functional step order/dependencies; bind target refs; bind non-secret values and Secret References; validate required definitions; expose masked metadata only; set concurrency/max-target/time/retry/stop constraints; carry prohibited scope; attach verification and rollback requirements; compare plan versions; detect missing/extra bindings; prepare target/readiness assessment.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect definitions/bindings | operator/reviewer | Execution Plan | 0 | read | plan context | non |
| validate binding completeness | reviewer | Plan | 1 | definitions known | missing/mismatch list | non |
| create/update plan/bindings | Response Operator | Plan/Binding | 2 | exact source versions | versioned no-effect plan | OPEN-013 |
| attach Secret Reference | authorized operator | binding | 2 | reference visible/allowed | reference only | OPEN-013 |
| resolve secret value/start executor | none here | secret/runtime | — | outside capability | no effect | Settings/executor later |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| map definitions to available inputs | oui | type/name rules | oui | mapping suggestion | binding form/matrix |
| validate required bindings | oui | oui | oui | explain missing | checklist |
| draft execution ordering | oui | explicit dependencies | oui | proposal | dependency graph/table |
| flag scope/constraint mismatch | oui | deterministic diff | oui | explanation | Decision↔Plan diff |
| choose secret value/authorize/start | authorized executor path | no autonomous | no | prohibited | explicit reference + governed start |

## 14. États fonctionnels
`not-prepared`, `draft`, `bindings-incomplete`, `secret-reference-missing`, `dependency-incomplete`, `scope-mismatch`, `constraint-mismatch`, `reviewing`, `plan-ready-for-readiness`, `stale`, `blocked`, `superseded`, `cancelled`.

## 15. États d’interface
Loading preserves current plan/version ; Empty means no plan yet ; Partial names missing bindings/dependencies ; Error preserves last valid version ; Offline cannot finalize fresh bindings ; Permission denied masks restricted metadata ; Stale requires diff/review after Decision/Playbook/source change.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Execution Plan | versioned no-effect plan | CAP-GOV-020/021/022 | exact Decision/Playbook/targets/scope/constraints |
| Parameter Bindings | plan component | readiness/executor handoff | source/ref/type visible, secret values absent |
| binding/constraint blocker | review event | Response Operator/CAP-GOV-018 | exact issue and source |
| plan diff | review artifact | reviewer/Decision path | material changes explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-018 | compatible version accepted | CAP-GOV-019 | exact Playbook/version + limits | Playbooks |
| CAP-GOV-019 | plan prepared | CAP-GOV-020 | Plan, target refs, dependencies, bindings | Playbooks |
| CAP-GOV-019 | authorization-impacting mismatch | CAP-GOV-021/GOV-1 | plan diff vs Decision | same plan lineage |
| source/Playbook change | bindings stale | CAP-GOV-018/019 | old/new versions + affected bindings | same handoff |

## 18. Dépendances
CAP-GOV-015..018/020/021/022, Playbook, Settings Secret Reference/integrations/environments, Studio Workflow/Tool input metadata, target sources, Shared Versioning/Trace/Linking, Security secrets/access policy, OPEN-008/013/015.

## 19. Source de vérité
Govern owns Execution Plan and binding intent. Playbook/Decision remain source for procedure/authority. Settings remains source for Secret Reference and actual secret material. Studio/Endpoint/runtime owners remain source for technical input/execution contracts.

## 20. Provenance et audit
Record Decision/handoff/Playbook versions, plan version, each step/dependency, target refs, parameter definition source/version, binding source/ref, Secret Reference ids/metadata only, constraints, reviewer, diffs, missing values, automation/AI proposals/dispositions and timestamps. Never record raw secret values.

## 21. Permissions fonctionnelles
Execution Plan create/update/read; parameter binding; restricted parameter metadata; Secret Reference metadata binding; plan compare; automated recommendation request; provenance export. Secret reveal/copy/export/use and runtime execution require separate source-owner/executor permissions.

## 22. Limites et erreurs
Missing definition, invalid type, raw-secret attempt, inaccessible Secret Reference, extra target, scope expansion, stale Playbook, dependency conflict or concurrent plan edit must block/flag the plan. A missing value is never fabricated.

## 23. Métriques
Plans complete on first review; binding errors; raw-secret exposure attempts blocked; material plan diffs; stale-plan rate; average plan rework; plans that expand Decision scope — target zero.

## 24. Classification de livraison
`defined` / `planned`; functional planning only. No secret resolver, schema, API, executor, command format or provider/runtime selected.

## 25. Critères d’acceptation
**Given** a required credential input, **When** the operator binds it, **Then** the Execution Plan stores only an authorized Secret Reference and never the raw secret value.

**Given** a plan introduces a target absent from the Decision scope, **When** validation runs, **Then** scope mismatch is explicit and the plan cannot become readiness-ready.

**Given** no AI, **When** the plan is prepared, **Then** forms, dependency tables, deterministic binding checks and human review provide the complete function.

## 26. Questions ouvertes
OPEN-008 covers actual runtime support, OPEN-013 class-2 plan/binding governance, OPEN-015 execution bridge. Final technical parameter schemas/secret resolution are future source-owner/implementation decisions; no new OPEN.

## 27. Consommateurs documentaires
Playbooks, CAP-GOV-020/021/022/025, Runs & Rollback, Settings/Studio/Endpoint boundaries, future Objects/Permissions/Screens/Technique and GOV-2 conformance report.
