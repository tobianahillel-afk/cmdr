---
id: CAP-STD-050
title: Automation Run Outcome, Consumer Handoff and Follow-up
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-050 — Automation Run Outcome, Consumer Handoff and Follow-up

## 1. Définition
Définit le Studio Runtime Outcome d’une Automation Run : completion state, produced outputs/Tool Call refs, warnings/failures/partial completion, human-review context, timestamps, owner, provenance, follow-up need and consumer handoff. Studio outcome ≠ Govern Result, Evidence or Finding automatically.

## 2. Problème utilisateur
Les consommateurs ont besoin du résultat technique d’une automation sans que Studio requalifie leurs objets métier ni transforme une output en Result/Evidence/Finding.

## 3. Objectifs
Assembler un outcome attribué; préserver partial/failures/unknown; qualifier uniquement le statut Studio; remettre outputs/provenance au caller; préparer follow-up/handoff under destination semantics.

## 4. Non-objectifs
Ne crée pas Govern Result, Investigate Evidence/Finding, Command Incident/Task, report publication, verification governance, rollback outcome or downstream mutation.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; source-product caller, Investigator/Incident operator/Response Operator and Auditor consumers.

## 7. Conditions d’entrée
Automation Run has terminal/stopped/cancelled/partial/unknown disposition sufficient to explain current technical outcome; output/error refs available or explicitly missing.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Automation Run terminal context | CAP-STD-043/047 | completion/failure/partial state | oui | latest confirmed | outcome inconclusive |
| step/Tool Call outputs | CAP-STD-044/047 | technical refs/outputs | conditionnel | source timestamp | partial/missing |
| Human Gate/intervention refs | CAP-STD-041/039/049 | review/control context | conditionnel | current/history | none invented |
| caller/return origin | CAP-STD-042 | consumer + source object ref | oui | pinned | handoff unresolved |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run / Steps | CMDR Studio | terminal state/history | read |
| Tool Call | CMDR Studio | technical output/error refs | read/link |
| Evidence / Finding | Investigate | downstream destination refs only | read/link |
| Result / Response Run | Govern | downstream/correlation refs only | read/link |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Studio Runtime Outcome | create/update/supersede functional outcome | CMDR Studio | technical/runtime outcome ≠ Govern Result |
| Consumer Handoff | prepare/send context reference | Studio local handoff | destination owns qualification/mutation |

## 11. Fonctionnalités
Outcome assembly; completion/partial/failure/warnings; output refs; human review/interventions; timestamps/owner; follow-up need; source consumer return; destination qualification boundary.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Run outcome | Studio Operator | Studio Runtime Outcome | 0 | read | outputs/partials/provenance visible | non |
| Validate outcome completeness | Studio Reviewer | Outcome | 1 | Run/step refs | missing/contradiction findings | non |
| Prepare consumer handoff/follow-up | Studio Operator | Consumer Handoff | 2 | destination/context permissions | bounded handoff | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
draft-outcome, complete-technical, partial, failed, cancelled, stopped, timed-out, status-unknown, follow-up-required, handed-off, disputed, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Studio Runtime Outcome | Studio record/concept | source caller/Control Room | Run state + raw refs + warnings/partials/provenance |
| Consumer Handoff | cross-product context | Command/Investigate/Govern | no destination requalification |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-043/047 | Run reaches disposition | CAP-STD-050 | Run/step/output/error refs | Run |
| CAP-STD-050 | return to caller | source product | outcome/ref/provenance/return-origin | Studio |
| CAP-STD-050 | handoff to Govern | Govern | Studio outcome + Automation Run refs | Studio |

## 18. Dépendances
CAP-STD-042..049/051; STD-1 Tool Call output contract; Command/Investigate/Govern object boundaries; Shared linking/reporting; OPEN-013/015.

## 19. Source de vérité
Studio owns Automation Run technical outcome semantics. Tool Call/raw sources remain source-attributed; Investigate owns Evidence/Finding qualification; Govern owns canonical Result; Command owns Incident/Task.

## 20. Provenance et audit
Record Outcome version, Run/step/Tool Call refs, completion/partial/failure/unknown, warnings, Human Gate/intervention refs, timestamps, owner/reviewer, follow-up/handoff refs and correlation.

## 21. Permissions fonctionnelles
Run outcome read; restricted outputs; consumer handoff prepare; provenance export preparation. Destination-object create/qualify permissions remain with destination product.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing/contradictory output, lost terminal status, partial branch, restricted output, destination permission denial or stale downstream context remains explicit; no automatic qualification.

## 23. Métriques
Outcomes by technical disposition; partial/unknown outcomes; handoff completion; downstream qualification rejection; outputs auto-promoted to Result/Evidence/Finding—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** agent-generated output is handed to Investigate, **When** Investigate receives it, **Then** it remains attributed automation output and does not become Evidence or Finding automatically.

**Given** a Studio outcome is handed to Govern, **When** Govern receives it, **Then** it remains a technical source/ref and does not become the canonical Result.

**Given** a Run is partially completed, **When** outcome is assembled, **Then** completed and failed steps remain distinct and the technical outcome is not simplified to success.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, source products, Command, Investigate, Govern, Shared Reporting/Linking, CAP-STD-051, Audit/Quality.
