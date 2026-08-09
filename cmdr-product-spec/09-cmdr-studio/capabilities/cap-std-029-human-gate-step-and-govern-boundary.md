---
id: CAP-STD-029
title: Human Gate Step and Govern Boundary
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-029 — Human Gate Step and Govern Boundary
## 1. Définition
Contrat d’un Human Gate Step dans un Workflow : conceptual pause point, reviewer/context/choices/comment, completion/expiry/escalation reference and provenance. Human Gate ≠ Approval/Decision; completion ≠ production authorization.
## 2. Problème utilisateur
Un checkpoint humain d’orchestration peut sinon être interprété comme une approbation d’autorité et contourner Govern.
## 3. Objectifs
Permettre un step humain explicite et auditable, référencer le Human Gate canonique, et handoff vers Govern si Approval/Decision est requis.
## 4. Non-objectifs
Ne définit pas runtime pause/resume, final assignment/notification engine, Govern Approval/Decision semantics, Automation Run lifecycle or production authorization.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow/Human Gate step | CMDR Studio | step context/options/reviewer need | oui | Workflow version | incomplete |
| review context | source products/Workflow | read-only context refs | oui | current/pinned | gate cannot be reviewed |
| Govern requirement | Policy/Decision Authority projection | authority need/ref | conditionnel | current | completion non-authorizing |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Human Gate | CMDR Studio | identity/context/outcome/expiry | read/reference |
| Workflow | CMDR Studio | step/context | read/write |
| Approval / Decision | Govern | distinct authority refs | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Human Gate Step configuration | create/update | CMDR Studio | Workflow step references Human Gate contract |
| Govern handoff requirement | derive/reference | CMDR Studio | no Approval/Decision creation |
## 11. Fonctionnalités
Reviewer need; requested review/context; choices/comment; completion/expiry; escalation ref; Govern-required marker; provenance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect gate config | Automation Designer | Human Gate Step | 0 | read | context/authority need visible | non |
| Validate gate/Govern consistency | Studio Reviewer | Human Gate Step | 1 | policy/refs available | assessment | non |
| Configure Human Gate step | Automation Designer | Workflow draft | 2 | manage | draft changed | OPEN-007/013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest gate placement | oui | oui | oui | oui | manual insertion |
| summarize review context | oui | oui | oui | oui | raw context refs |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
not-configured, pending-reference, configured, authority-required, expired-reference, inconsistent, valid-for-review.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Human Gate Step contract | Workflow step | CAP-STD-030/future STD-3 | gate context + authority boundary |
| Govern authority handoff requirement | reference need | Govern | no Approval/Decision implicit |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow | reach gate conceptually | Human Gate context | Workflow/step/reviewer/context refs | Workflow |
| Human Gate completion future | evaluate next path | Workflow | outcome + provenance + authority status | Workflow |
| authority required | handoff | Govern | Approval/Decision context refs | return only with Govern ref |
## 18. Dépendances
Human Gate object/contract; OPEN-007; OPEN-013; CAP-GOV-025 and Govern Approval/Decision; Shared Notifications future; STD-3 runtime.
## 19. Source de vérité
Studio owns Human Gate and Workflow step semantics. Govern alone owns Approval/Decision/response authority. OPEN-007 remains unresolved except non-equivalence.
## 20. Provenance et audit
Record Workflow/version/step, Human Gate ref/version, reviewer context, choices, completion/expiry refs, Govern requirement/refs, actor and correlation.
## 21. Permissions fonctionnelles
Workflow/Human Gate read/configure; reviewer context read; future gate completion permission separate; Govern Approval/Decision permissions remain Govern-owned. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Reviewer unavailable, context denied, gate expired, authority required but absent, outcome invalid or tenant mismatch blocks continuation/readiness.
## 23. Métriques
Gates with unresolved authority, expired refs, missing reviewers, Govern-handoff completeness, non-equivalence violations.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Human Gate is completed but a Govern Approval is separately required, **When** continuation is assessed, **Then** completion alone does not authorize continuation.

**Given** a Govern Approval is required and absent, **When** Workflow readiness is checked, **Then** the authority dependency remains explicit and blocked for production effect.

**Given** AI is disabled, **When** a Human Gate is configured/reviewed, **Then** manual configuration and deterministic consistency checks remain available.
## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder, Human Gates, Govern/CAP-GOV-025, future STD-3 runtime, Security, Quality.
