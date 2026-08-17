---
id: CAP-STD-039
title: Agent Oversight, Intervention and Escalation
product: cmdr-studio
module: automation-agents
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-039 — Agent Oversight, Intervention and Escalation

## 1. Définition
Définit human oversight of Agent/Team activity : observer/human owner, inspect proposals/state, reject proposal, modify allowed input/context, pause/stop candidate, escalation, unsafe/uncertain/repeated-failure/unexpected-Tool/scope-expansion handling. Intervention ≠ Decision modification or authority bypass.

## 2. Problème utilisateur
Sans supervision explicite, une autonomie bornée peut continuer malgré incertitude, répétition d’échec ou demande inattendue de Tool et l’intervention humaine peut contourner les mêmes règles qu’elle cherche à imposer.

## 3. Objectifs
Rendre oversight continu et attributable ; permettre reject/context correction/stop-escalation ; conserver source state; revalider après intervention; router les besoins d’autorité vers Govern.

## 4. Non-objectifs
Ne modifie pas Govern Decision, ne crée pas permission, ne bypass pas Human Gate, ne remplace pas Control Room run controls, ne définit pas UI détaillée or emergency authority.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Human supervisor/Studio Operator principal; Automation Designer, Security Reviewer, Govern operator and Auditor secondary.

## 7. Conditions d’entrée
Agent/Team or Run-related context visible; observer has functional permission; provenance and last confirmed state available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Agent/Team/proposal state | Studio runtime | current bounded state | oui | last confirmed | oversight Partial |
| human owner/observer | Settings + Studio | principal/role projection | oui | current | intervention denied |
| risk/authority context | Security/Govern | permission and required authority | conditionnel | current | escalation only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Agent / Agent Team | CMDR Studio | current config/proposals | read |
| Automation Run ref | CMDR Studio | current run projection when present | read/control request |
| Decision / Approval | Govern | authority status only | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Oversight Context | create/update observation/disposition | CMDR Studio | does not alter source authority |
| Runtime Intervention candidate | prepare/record | CMDR Studio | actual Run control uses CAP-STD-046/049 |

## 11. Fonctionnalités
Observe; reject proposal; adjust allowed input/context; mark unsafe/uncertain; request pause/stop; escalate to Human Gate/Govern; repeated failure and unexpected Tool detection; intervention provenance/revalidation.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Agent activity | human supervisor | Agent/Team/Run | 0 | read | current state/provenance | non |
| Reject unsafe proposal | human supervisor | Agent Step Proposal | 2 | review permission | proposal rejected | non |
| Modify allowed input/context | authorized operator | Run/Agent context | 2 | within authority | context version changes + revalidation | OPEN-013 |
| Escalate effectful authority need | operator | Action Request/Decision path | 2 | Govern required | handoff only | Govern |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
observing, attention-required, intervention-proposed, intervention-applied, revalidation-required, escalation-required, waiting-human, unsafe, repeated-failure, stopped-by-policy-reference, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Oversight disposition | Studio record | CAP-STD-038/042/046/049 | actor/reason/context explicit |
| Govern escalation package/reference | handoff context | Govern | no Approval/Decision fabricated |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Agent/Team | unsafe/uncertain condition | CAP-STD-039 | state/proposal/context | source |
| CAP-STD-039 | pause/stop requested | CAP-STD-046/049 | reason/run ref/actor | oversight |
| CAP-STD-039 | authority required | Govern | bounded context + source refs | Studio |

## 18. Dépendances
CAP-STD-034..038/042/046/049/051; Govern Action Request/Approval/Decision; Security; Settings Principal; OPEN-007/013/015.

## 19. Source de vérité
Studio owns oversight/intervention semantics for its agents/runs. Govern owns production authority and Decisions; source run/agent objects remain Studio source-of-truth.

## 20. Provenance et audit
Record observer/human owner, source state/proposal, trigger reason, intervention/rejection/context change, required revalidation, escalation destination, actor/time and correlation.

## 21. Permissions fonctionnelles
Oversight read; proposal reject; allowed-context edit; pause/stop request; escalation preparation; restricted context read. None imply Govern approver eligibility.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing status, stale context, insufficient permission, unsafe scope expansion, repeated failure or authority ambiguity leaves work paused/escalated/unknown rather than silently continuing.

## 23. Métriques
Interventions by trigger; repeated-failure escalations; unexpected Tool requests; context changes requiring revalidation; interventions bypassing Govern—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Agent makes an unexpected Tool request, **When** a supervisor reviews it, **Then** the request can be rejected/escalated without granting the Tool or changing permissions.

**Given** Control Room/operator intervention changes allowed input context, **When** continuation is considered, **Then** relevant permission, constraint and readiness checks run again.

**Given** a real production effect needs new authority, **When** oversight escalates, **Then** Studio prepares a Govern handoff and does not modify a Decision locally.

## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Automation Agents, Agent Teams, Human Gates, Control Room, CAP-STD-046/049/051, Govern, Security, Audit/Quality.
