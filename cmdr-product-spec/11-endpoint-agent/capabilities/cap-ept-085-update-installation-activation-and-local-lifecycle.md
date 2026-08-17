---
id: CAP-EPT-085
title: Update Installation, Activation and Local Lifecycle
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-085 — Update Installation, Activation and Local Lifecycle

## 1. Définition
Définir les transitions techniques locales d’installation et d’activation d’un update Endpoint, séparées de l’assignment Settings et de la santé post-update.

## 2. Problème utilisateur
Une installation terminée peut être confondue avec activation ou santé, tandis qu’un restart/interruption peut laisser un état partiel.

## 3. Objectifs
Représenter install pending/installing/installed, activation pending/activated, restart-required si sourcé, interruption, local transition, status et provenance.

## 4. Non-objectifs
Aucun installateur, service OS, package format, scheduler, update authority globale, support claim, health verdict final ou implementation.

## 5. Propriétaire
Endpoint owns local technical install/activation state; Settings retains assignment and Govern/Security retain applicable authority/policy boundaries.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Release Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-084 ready state, exact target/package refs, applicable policy/authority, Agent/tenant/environment and stable attempt identity.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| staged package/readiness | CAP-EPT-084 | precondition | oui | current | block |
| target assignment | CAP-EPT-082 | desired version | oui | current | block/mismatch |
| package context | CAP-EPT-083 | exact release ref | oui | pinned | block |
| policy/authority context | Settings/Security/Govern as applicable | control context | oui if required | current | denied/blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Update Readiness Assessment | Endpoint | readiness/blockers | read |
| Local Package Context | Endpoint | release/version | read |
| Endpoint Policy | Settings | execution restrictions | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Update Operation | create/transition | Endpoint | exact attempt/target |
| Installation State | transition | Endpoint | installed != activated |
| Activation State | transition | Endpoint | activated != healthy |

## 11. Fonctionnalités
Start bounded local install when authorized; track progress handoff; record installed; represent activation pending/activated; preserve restart-required/interrupted/unknown states if sourced; never infer health.

## 12. Actions utilisateur
Inspect Class 0; readiness/eligibility Class 1; install/activate request effectful Class 3 by default where applicable; class may remain source/policy dependent under OPEN-013.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate preconditions | oui | oui | oui | explain | deterministic checks |
| track lifecycle | oui | oui | oui | summarize | state events |
| authorize/activate arbitrarily | non | policy/authority only | non autonome | interdit | accountable control path |

## 14. États fonctionnels
`install-pending`, `installing`, `installed`, `install-interrupted`, `install-failed`, `activation-pending`, `activated`, `restart-required`, `activation-failed`, `state-unknown`, `superseded`.

## 15. États d’interface
No Screen ID. Installed and activated are visually/conceptually distinct; offline/unknown never implies success.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Local Update Operation | Endpoint fact | CAP-EPT-086/087/088 | exact attempt/provenance |
| install/activation state | Endpoint fact | Settings/health consumers | no health inference |
| execution events | local audit context | CAP-EPT-097 | sensitive fields masked |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-084 | ready + allowed | CAP-EPT-085 | target/package/readiness | readiness ref |
| CAP-EPT-085 | progress/failure | CAP-EPT-086 | operation/attempt/state | lifecycle ref |
| CAP-EPT-085 | activated | CAP-EPT-087 | expected/observed version + operation | health verification ref |

## 18. Dépendances
CAP-EPT-082..084/086/087/097; Settings; Security permission model; OPEN-008/013.

## 19. Source de vérité
Endpoint SOT for local install/activation facts only; assignment/authority/support remain external.

## 20. Provenance et audit
Agent/tenant, target/package, operation/attempt, requester/authority refs, start/end times, state transitions, interruption/restart hints and errors.

## 21. Permissions fonctionnelles
Update state read, execution request, activation, sensitive provenance read, cross-tenant deny; step-up/SoD and final RBAC remain canonical Security concerns.

## 22. Limites et erreurs
Ready != installed; installed != activated; activated != healthy; execution success != compatibility/support; interruption remains explicit.

## 23. Métriques
Install/activation outcomes, interruptions, restart-required, state-unknown, denied/blocked execution counts.

## 24. Classification de livraison
`draft / defined / planned`; no updater/install engine or platform implementation.

## 25. Critères d’acceptation
**Given** installation succeeds but activation is pending, **When** state is read, **Then** installed and activation-pending remain distinct.

**Given** authority/policy is missing, **When** activation is requested, **Then** execution is blocked and no state is fabricated.

**Given** AI is unavailable, **When** lifecycle events arrive, **Then** deterministic state tracking still works.

## 26. Questions ouvertes
OPEN-008/013 remain open; platform-specific restart and final action-class policy are not selected.

## 27. Consommateurs documentaires
EPT-6 Update, Settings Fleet, Security, CAP-EPT-086..088/097, Quality, registers and Roadmap.