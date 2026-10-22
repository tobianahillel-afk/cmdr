---
id: CAP-STD-041
title: Human Gate Lifecycle, Response, Expiration and Govern Boundary
product: cmdr-studio
module: human-gates
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-041 — Human Gate Lifecycle, Response, Expiration and Govern Boundary

## 1. Définition
Définit le runtime lifecycle d’un Human Gate et sa response : requested, awaiting-review, viewed, information-requested, responded, accepted-for-workflow, rejected-for-workflow, expired, cancelled, superseded, escalated. `accepted-for-workflow` ≠ Govern Approval/Decision or production authority.

## 2. Problème utilisateur
Des états historiques `approved/rejected` peuvent être interprétés comme une autorisation Govern. Runtime orchestration a besoin d’un outcome de revue sans créer une seconde chaîne d’approbation.

## 3. Objectifs
Exprimer response/lifecycle/expiry/escalation ; distinguer workflow continuation from authority ; préserver reviewer/provenance; requérir Govern lorsque l’effet l’exige.

## 4. Non-objectifs
Ne modifie pas le physical Human Gate object schema/states, ne ferme pas OPEN-007, ne crée pas Approval/Decision, ne définit pas final approver SoD/authority or UI.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Human reviewer principal; Studio Operator/Automation Designer, Govern Reviewer and Auditor secondary.

## 7. Conditions d’entrée
Valid Human Gate Request, authorized reviewer, current deadline/context and return destination; any external authority requirements separately identifiable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Human Gate Request | CAP-STD-040 | request/context/choices/deadline | oui | current | no response |
| reviewer identity/permission | Settings + Security | eligible Studio reviewer context | oui | current | response denied |
| current source/Run state | Studio/source owner | review context freshness | oui | current enough | information request/stale |
| Govern authority requirement | Govern/Security | whether separate Approval/Decision needed | conditionnel | current | continuation blocked if required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Human Gate | CMDR Studio | current review lifecycle | read/respond |
| Automation Run / Workflow | CMDR Studio | waiting step + return path | read/link |
| Approval / Decision | Govern | separate authority state | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Human Gate response | record/supersede functional response | CMDR Studio | workflow review outcome only |
| Human Gate lifecycle projection | transition/expire/cancel/escalate | CMDR Studio | physical object schema unchanged |

## 11. Fonctionnalités
Review response; info-request; accepted/rejected-for-workflow; expiry/cancel/supersession; escalation; separate Govern authority check; return-to-Run; no-equivalence messaging/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Respond to Gate | human reviewer | Human Gate | 2 | review permission + current request | response recorded | non-authoritative |
| Expire/cancel Gate | Studio runtime/operator | Human Gate | 2 | deadline/cancel condition | terminal review state | OPEN-013 |
| Request Govern authority | Studio Operator | Govern handoff | 2 | effect requires authority | separate Govern request/ref | Govern |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
requested, awaiting-review, viewed, information-requested, responded, accepted-for-workflow, rejected-for-workflow, expired, cancelled, superseded, escalated.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Human Gate response | Studio review outcome | Automation Run/Workflow | choice/reviewer/time/provenance; no Govern authority |
| Govern escalation/ref | cross-product handoff | Govern | separate Approval/Decision identity |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-040 | review starts | CAP-STD-041 | request/reviewer/context | origin Run |
| CAP-STD-041 | accepted/rejected/expired | CAP-STD-043/044 | gate outcome + return step | Human Gate |
| CAP-STD-041 | authority required | Govern | request context + source refs | same Studio Run waits |

## 18. Dépendances
CAP-STD-029/040/042..044/049; Govern Approval/Decision and CAP-GOV boundaries; Settings identity; Security/SoD; Shared Notifications; OPEN-007/013/015.

## 19. Source de vérité
Studio is source for Human Gate review lifecycle/outcome. Govern is source for Approval/Decision and production authority. The bridge remains unresolved under OPEN-007 except non-equivalence.

## 20. Provenance et audit
Record gate/request/response versions, reviewer, choice, context viewed/masked, information request, deadline/expiry/cancel/supersession/escalation, linked Govern refs, origin/return step and correlation.

## 21. Permissions fonctionnelles
Human Gate read/respond/reassign/escalate/cancel as functional needs; sensitive context permission separate; no automatic Govern approver eligibility.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Expired Gate, stale context, reviewer permission loss, source status change, conflicting response, missing Govern authority or lost return Run remains explicit; no implicit approval/rejection translation.

## 23. Métriques
Response latency; expired gates; info requests; escalations; accepted-for-workflow paths still blocked by missing Govern authority; false Approval equivalences—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Human Gate expires without response, **When** its waiting Run is evaluated, **Then** the Gate becomes expired and no Govern Approval rejection is fabricated.

**Given** a reviewer selects `accepted-for-workflow`, **When** the next step requires Govern authority, **Then** the Run remains blocked until the separate Govern Approval/Decision contract is satisfied.

**Given** a Govern Decision is absent for an action that requires authority, **When** the Human Gate is complete, **Then** Studio cannot start that effect.

**Given** AI is unavailable, **When** a Human Gate is reviewed, **Then** explicit human choices and deterministic deadline/authority checks provide the complete path.

## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Human Gates, Automation Runs, Workflows, Control Room, Govern, Security, Settings, Audit/Quality.
