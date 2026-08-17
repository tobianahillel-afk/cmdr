---
id: CAP-EPT-057
title: Live Response Session Lifecycle and Runtime State
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-057 — Live Response Session Lifecycle and Runtime State

## 1. Définition
Définir le lifecycle technique local d’une Live Response Session : requested/validating/connecting/active/idle/waiting/degraded/disconnected/reconnecting/closing/closed/failed/expired/cancelled, avec distinction stricte entre intention et confirmation target-side.

## 2. Problème utilisateur
Une session `requested` ou `connecting` peut être représentée à tort comme active, tandis qu’une déconnexion peut être confondue avec fermeture ou succès des opérations en cours.

## 3. Objectifs
Conserver last confirmed state/time ; gérer inactivity/expiry/concurrency/reconnect ; distinguer disconnect/close/failure ; exposer unknown state ; préserver session target/authority refs et operation independence.

## 4. Non-objectifs
Aucune state machine technique finale, transport/reconnect protocol, terminal UX, command authority, Response Run lifecycle, rollback ou containment.

## 5. Propriétaire
Endpoint possède l’état technique local de session. Investigate conserve le record métier Live Session ; Govern son Response Run ; Studio son Automation Run.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Case Analyst, Session Supervisor et Auditor.

## 7. Conditions d’entrée
CAP-EPT-056 technical session ref, target binding valide, local/runtime state observations, expiration/inactivity constraints et controls autorisés.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Technical Session | CAP-EPT-056 | lifecycle subject | oui | current version | no state |
| local connectivity/runtime state | Endpoint | confirmation | oui when available | event time | status-unknown |
| expiry/inactivity constraints | policy/session | temporal bounds | oui | current | blocked |
| open/close/reconnect/cancel intent | authorized caller | control request | non | request time | no transition request |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Technical Session | Endpoint | target/constraints/state | read/manage |
| business Live Session | Investigate | correlation | read ref |
| Response Run | Govern | authority correlation | read only |
| Automation Run/Tool Call | Studio | caller correlation | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Session State | reconcile/version | Endpoint | request != confirmation |
| Session Lifecycle Event | append | Endpoint | actor/source/time/reason |
| Reconnect/Closure Attempt | create/link | Endpoint | exact attempt/status |

## 11. Fonctionnalités
Reconcile open request and target confirmation, track active/idle/waiting/degraded, mark disconnect, accept reconnect intent and confirmation, enforce expiry, request close/cancel and preserve final/unknown status without affecting operation results.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect lifecycle | Operator | session state | 0 | read | last confirmed/requested states | non |
| validate transition eligibility | service | session | 1 | current state/constraints | allow/block reason | non |
| request reconnect/close/cancel | authorized operator | session | 2 | compatible state | control request | OPEN-013 selon policy |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| state reconciliation | oui | oui | oui | explain | source events |
| expiry/inactivity evaluation | oui | oui | oui | explain | timers/rules |
| reconnect recommendation | oui | rules | oui | suggestion | manual control |
| fabricate active/closed state | non | interdit | non | interdit | status-unknown |

## 14. États fonctionnels
`requested`, `validating`, `connecting`, `active`, `idle`, `waiting`, `degraded`, `disconnected`, `reconnecting`, `closing`, `closed`, `failed`, `expired`, `cancel-requested`, `cancelled`, `status-unknown`.

## 15. États d’interface
Aucun Screen ID. `connecting != active`, `disconnected != closed`, `close-requested != closed`; stale/unknown always show last source time.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Session State | Endpoint projection | CAP-EPT-058..063/Investigate | requested/confirmed distinction |
| lifecycle event | local audit/provenance ref | CAP-EPT-063/064 | source/time/actor |
| constraint/expiry reason | diagnostic | caller | no authority expansion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-056 | open requested | CAP-EPT-057 | session/target/constraints | request origin |
| CAP-EPT-057 active | command/script/file action selected | CAP-EPT-058/060/061 | session/actor/authority refs | same session |
| CAP-EPT-057 close | CAP-EPT-063 | closure/audit context | operations remain independent |

## 18. Dépendances
CAP-EPT-056/058..064, Investigate CAP-INV-209, Endpoint session-management source, Shared Recovery boundary, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT for local session runtime state; Investigate SOT business Live Session; Govern/Studio retain their Run states.

## 20. Provenance et audit
Session id, target, state requests/confirmations, source timestamps, inactivity/expiry, reconnect attempts, close/cancel intents, actor and correlation IDs.

## 21. Permissions fonctionnelles
Session read, reconnect/close/cancel request, sensitive transcript/state refs, operator control, cross-tenant deny.

## 22. Limites et erreurs
Session active ≠ unrestricted authority; disconnect ≠ operation failure/success automatically; closure ≠ target state restored; cancellation ≠ rollback.

## 23. Métriques
Open/reconnect failures, active/disconnected durations, expiry/cancellation, unknown-state duration, false-active/false-closed mappings target zero.

## 24. Classification de livraison
`draft / defined / planned`; no session protocol/runtime implementation.

## 25. Critères d’acceptation
**Given** a session disconnects during an operation, **When** lifecycle is reconciled, **Then** session becomes disconnected/reconnecting while operation status remains independent.

**Given** reconnect succeeds, **When** target confirms, **Then** state returns active with a new confirmation timestamp and prior gap preserved.

**Given** close is requested but not confirmed, **When** status is displayed, **Then** it remains closing/status-unknown rather than closed.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final protocol and concurrency policy not selected.

## 27. Consommateurs documentaires
EPT-4 Live Response, Investigate Live Session, CAP-EPT-063/064, Govern/Studio, Quality.
