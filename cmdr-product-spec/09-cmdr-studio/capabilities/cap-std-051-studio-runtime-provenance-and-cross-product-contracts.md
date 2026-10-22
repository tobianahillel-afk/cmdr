---
id: CAP-STD-051
title: Studio Runtime Provenance and Cross-Product Contracts
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-051 — Studio Runtime Provenance and Cross-Product Contracts

## 1. Définition
Définit la provenance end-to-end Studio `Workflow/version → Agent/version → Automation Run → runtime context → steps/attempts → Tool Calls → Human Gates → interventions → retries/errors → outputs → Studio Runtime Outcome → downstream handoff`, et les contrats de non-equivalence avec Govern/Investigate/Command/Endpoint/Settings/Shared.

## 2. Problème utilisateur
Sans lineage complet, une automation n’est pas reconstructible et les consumers peuvent fusionner Automation Run, Tool Call, Response Run, Result, Evidence ou Shared Job.

## 3. Objectifs
Préserver identifiers/versions/correlation/source timestamps/actors/AI use/uncertainty; reconstruire runtime history; maintenir ownership boundaries; produire handoff traceable without secret leakage.

## 4. Non-objectifs
Ne crée pas Trace engine, audit storage, canonical Result/Evidence/Finding, final Automation Run physical schema, protocol/API, Endpoint capability or STD-4 assurance/deployment.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Auditor and Studio Operator principal; Automation Designer, source-product operators, Govern Reviewer and Security Reviewer secondary.

## 7. Conditions d’entrée
Automation Run lineage objects/refs available to extent produced; permission to inspect each source; missing/restricted links represented explicitly.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Workflow/Agent/Run lineage | CMDR Studio | exact ids/versions/parent refs | oui | historical pinned | provenance incomplete |
| steps/attempts/Tool Calls/Gates | CMDR Studio | runtime event/object refs | conditionnel | source timestamps | gap explicit |
| interventions/errors/outcome/handoff | STD-3 capabilities | attributed events/refs | conditionnel | historical | gap explicit |
| cross-product refs | Govern/Command/Investigate/Endpoint/Settings/Shared | source-owned relation refs | conditionnel | source-history | restricted/missing |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow / Agent / Agent Team / Human Gate / Tool Call | CMDR Studio | identity/version/provenance refs | read |
| Automation Run | CMDR Studio | runtime lineage/projection | read |
| Response Run / Result | Govern | distinct correlation/handoff refs | read |
| Evidence / Finding / Incident | Investigate/Command | consumer refs only | read |
| Job / Trace / Activity | Shared | generic provenance/transport refs | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Studio Runtime Provenance projection | derive/reconstruct | CMDR Studio | does not replace Shared Trace/source audit |
| Cross-product relation/handoff refs | link/supersede | respective owners + Studio projection | ownership/permission retained |

## 11. Fonctionnalités
Full lineage reconstruction; raw/source refs; version/correlation; gaps/contradictions; AI/model-provider reference and uncertainty; cross-product handoff; source-owner preservation; export preparation without raw secrets.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect runtime lineage | Auditor | Automation Run provenance | 0 | read | source-backed graph/timeline refs | non |
| Validate provenance completeness | Studio Reviewer | provenance set | 1 | source refs | gaps/contradictions | non |
| Prepare provenance export/handoff | authorized user | provenance package | 2 | source/export permissions | bounded package/ref | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
complete, partial, restricted, gap-detected, contradictory, source-unavailable, stale-source, reconstructed, handoff-linked, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Studio Runtime Provenance | provenance projection | Audit/Control Room/source consumer | exact versions/refs/gaps preserved |
| cross-product handoff lineage | relation set | Govern/Investigate/Command | ownership and return-origin preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow/Agent | instantiated | Automation Run | source versions/correlation | source |
| Run events | append/source update | CAP-STD-051 | event/object refs/timestamps | Run |
| Studio outcome | handoff | destination product | outcome + lineage + source refs | Studio |

## 18. Dépendances
CAP-STD-001..050; Shared Trace/Activity/Jobs; Security provenance/integrity/export; Settings provider/secret refs; Govern CAP-GOV-025/026 and OPEN-015; Command/Investigate/Endpoint boundaries.

## 19. Source de vérité
Each product remains source for its canonical objects and raw facts. Studio owns the runtime lineage relating its objects; Shared owns generic Trace/Activity mechanics; Govern owns Response Run/Result; Investigate owns Evidence/Finding.

## 20. Provenance et audit
Record all available ids/versions, tenant/env, callers/actors, Agent/Workflow versions, Run parent/child, step/attempt/Tool Call/Gate/intervention/retry/error/output/outcome refs, source timestamps, AI/model-provider references, uncertainty, handoff/return-origin and gaps/contradictions. Raw secrets excluded.

## 21. Permissions fonctionnelles
Runtime provenance read/export preparation; source-object reads under source permission; restricted Run context; cross-tenant provenance separate; secure export remains Security/Shared governed.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing source, retention gap, restricted ref, conflicting status, orphan Tool Call/attempt, cross-tenant access denial or unavailable Shared Trace remains a documented gap; no invented event repairs history.

## 23. Métriques
Provenance completeness; orphan refs; gaps/contradictions; cross-product correlation; source-version coverage; un-attributed AI outputs; raw-secret leakage—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Studio Automation Run is correlated with a Govern Response Run, **When** lineage is reconstructed, **Then** both identities and owners remain distinct and linked through explicit correlation.

**Given** an output is consumed by Investigate, **When** provenance is inspected, **Then** Workflow/Agent/Run/step/Tool Call lineage remains accessible subject to permission and Evidence qualification remains Investigate-owned.

**Given** a provenance source is unavailable, **When** reconstruction runs, **Then** a gap is recorded and no event/status is invented to complete the chain.

## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
All Studio modules, Command, Investigate, Govern, Endpoint boundary, Settings, Shared Trace/Activity/Jobs, Security, Quality, Roadmap, future Objects/Technique.
