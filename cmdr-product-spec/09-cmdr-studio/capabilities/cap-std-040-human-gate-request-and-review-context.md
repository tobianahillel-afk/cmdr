---
id: CAP-STD-040
title: Human Gate Request and Review Context
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
# CAP-STD-040 — Human Gate Request and Review Context

## 1. Définition
Définit la request/runtime context d’un Studio Human Gate : originating Workflow/Automation Run, step, reason, reviewer candidate, allowed choices, data shown/masked, deadline, return destination and provenance. Human Gate Request ≠ Govern Approval Request.

## 2. Problème utilisateur
Un point de revue humaine sans contexte borné peut exposer trop de données, perdre la Run/step origin ou être confondu avec une demande d’approbation de production.

## 3. Objectifs
Construire un review package minimal et sourcé ; identifier reviewer candidate sans le rendre approver Govern ; définir choices/deadline/return-origin ; appliquer masking and tenant isolation.

## 4. Non-objectifs
Ne crée pas Govern Approval/Decision, n’établit pas approver eligibility, ne finalise pas UI, notification engine, SoD policy or production authority.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Human reviewer principal; Studio Operator/Automation Designer secondary; Govern/Security reviewer consumes boundary context.

## 7. Conditions d’entrée
Originating Workflow/Run/step known, reason and allowed choices defined, reviewer candidate resolvable, permission/masking context available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| originating Workflow/Run/step | CMDR Studio | exact runtime lineage | oui | current/pinned | request invalid |
| reason + allowed choices | Workflow/Run | review contract | oui | versioned | request incomplete |
| reviewer candidate | Platform Settings principal projection | candidate identity/role | oui | current | unassigned |
| review data refs | source owners | masked authorized context | conditionnel | freshness visible | Partial/restricted |
| deadline/return origin | Workflow/Run | temporal/navigation context | oui | current | expiry/return undefined |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Human Gate | CMDR Studio | gate identity/config | read |
| Automation Run / Workflow | CMDR Studio | origin/step/context refs | read |
| Principal | Platform Settings | reviewer identity projection | read |
| Approval / Decision | Govern | distinct authority refs only | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Human Gate Request context | create/update before response | CMDR Studio | not Govern Approval Request |
| Reviewer assignment candidate | assign/reassign functional | CMDR Studio + Settings identity | does not grant Govern approver authority |

## 11. Fonctionnalités
Origin/step/reason; candidate reviewer; allowed choices; authorized context package; masking; deadline; return destination; reassignment/escalation candidate; provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Gate request | human reviewer | Human Gate | 0 | read | context/choices/deadline visible | non |
| Validate review package/masking | Studio Operator | Human Gate Request | 1 | source refs/permissions | pass/block findings | non |
| Request/reassign Human Gate | Studio Operator | Human Gate | 2 | workflow/run contract + permission | awaiting review | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
prepared, requested, awaiting-review, viewed, information-requested, reassignment-needed, restricted-context, expired-candidate, cancelled, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Human Gate Request | Studio review context | human reviewer/CAP-STD-041 | origin/choices/deadline/return/provenance |
| review notification intent | event/ref | Shared Notifications | deep-link context only; Shared owns delivery |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow/Run step | requires human review | CAP-STD-040 | origin/step/reason/context | same Workflow/Run |
| CAP-STD-040 | reviewer opens | CAP-STD-041 | gate/request/authorized data | source Run |
| request needs authority | escalate/link | Govern | bounded context; not Approval object | Studio |

## 18. Dépendances
STD-2 CAP-STD-029; CAP-STD-039/041/042/043/049; Human Gate object/contract; Settings Principal; Shared Notifications; Govern OPEN-007; Security; OPEN-013/015.

## 19. Source de vérité
Studio owns Human Gate/request context semantics. Settings owns reviewer identity administration; source products own displayed data; Govern owns Approval/Decision and approver authority.

## 20. Provenance et audit
Record gate/request id, origin Workflow/Run/step/version, reason, choices, data refs/masking, reviewer candidate/reassignments, deadline, return origin, creator/time and correlation. Never raw secret values.

## 21. Permissions fonctionnelles
Human Gate read/request/reassign/escalate; reviewer context read under source permissions; sensitive-context read separately. Reviewer assignment does not imply Govern Approval permission.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing origin, invalid reviewer, denied source data, over-broad context, cross-tenant data, expired deadline, stale source or unavailable notification remains explicit; request does not auto-approve.

## 23. Métriques
Requests with complete origin; masking violations; reassignments; expired before response; cross-tenant blocks; requests mistaken for Govern Approval—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Workflow step requests a Human Gate, **When** the request is created, **Then** Run/Workflow/step, reason, choices, deadline and return destination are preserved.

**Given** review data includes a restricted source field, **When** the reviewer lacks permission, **Then** the field is masked/omitted while the Gate remains attributable.

**Given** a Human Gate request concerns an action needing production authority, **When** it is prepared, **Then** it remains a Studio review request and separately references/escalates to Govern.

## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Human Gates, Workflows, Automation Runs, Control Room, Shared Notifications, Settings, Govern, Security, Quality.
