---
id: CAP-EPT-063
title: Live Response Operator Control, Session Closure and Audit Context
product: endpoint-agent
module: live-response
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-063 — Live Response Operator Control, Session Closure and Audit Context

## 1. Définition
Définir présence/contrôle opérateur, takeover conceptuel, cancel/stop requests, disconnect/reconnect, session timeout, close/forced-closure conceptuelle, sensitive-command visibility et audit context, sans assimiler fermeture à restoration/rollback.

## 2. Problème utilisateur
Plusieurs opérateurs, session inactivity, disconnect ou emergency closure peuvent rendre ambigu qui contrôlait l’exécution et si les opérations ont réellement cessé.

## 3. Objectifs
Conserver operator identity/presence ; contrôler ownership/takeover selon permissions ; enregistrer controls et acknowledgements ; gérer timeout/disconnect/close ; préserver transcript/audit refs ; indiquer outstanding operations et unknown termination à la fermeture.

## 4. Non-objectifs
Aucun collaboration protocol, final concurrency policy, forced kill primitive, rollback, state restoration, containment, Govern verification ou detailed terminal UX.

## 5. Propriétaire
Endpoint owns target-side operator/session-control facts. Investigate owns business participants/session context; Govern authority/Response Run/rollback; Shared generic collaboration/trace.

## 6. Utilisateurs
Response Operator, Session Supervisor, Endpoint Operator, Investigation Lead, Govern Reviewer, Auditor.

## 7. Conditions d’entrée
Technical Session ref, authenticated/authorized operator context, current session state, control capability, outstanding attempts and policy/authority.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| session state | CAP-EPT-057 | control subject | oui | current | no control |
| operator/participant context | Investigate/Identity | actor refs | oui | current | denied/unknown |
| outstanding attempts | CAP-EPT-059..062 | technical operations | non | current/last confirmed | closure warning |
| policy/authority | Settings/Govern | control bounds | oui where applicable | current | blocked |
| local audit refs | Endpoint local audit | event sink/context | oui conceptually | append time | gap marker |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Session | Endpoint | state/target | read/manage controls |
| business Live Session | Investigate | participants/purpose | read ref |
| Technical Execution Attempts | Endpoint | outstanding/status | read/control according permission |
| Response Run | Govern | authority/correlation | read only |
| local-audit-event/Trace | Endpoint/Shared | audit refs | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Operator Control Context | create/version | Endpoint | actor/role/presence/control explicit |
| Session Control Event | append | Endpoint | request/confirmation separate |
| Session Closure Context | create | Endpoint | outstanding states/limitations preserved |

## 11. Fonctionnalités
Track operator presence/control, request takeover where policy permits, record disconnect/reconnect, issue cancel/stop/close controls, apply timeout, record forced closure concept only when sourced/authorized, retain outstanding operation states and audit references.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect control/audit | Operator/Auditor | context | 0 | read | actor/state visible | non |
| request takeover/reconnect/close | authorized operator | session | 2 | policy/permission | control request | OPEN-013 |
| cancel/stop technical attempt | authorized operator | attempt | 2 | interruptible + authority | request only | according original class |
| force effectful target restoration | aucun EPT-4 | target | 3+ | Govern/EPT-5 | not executed | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| track presence/control | oui | oui | oui | summarize | event list |
| enforce timeout | oui | timer/rules | oui | explain | deterministic timer |
| summarize outstanding states | oui | oui | oui | oui | raw status table |
| close/stop autonomously beyond policy | non | authority contract | non autonome | interdit | explicit operator/Govern |

## 14. États fonctionnels
`operator-present`, `observer-only`, `control-held`, `takeover-requested`, `disconnected`, `reconnecting`, `close-requested`, `closing`, `closed`, `expired`, `forced-closure-requested`, `closure-with-unknown-operations`.

## 15. États d’interface
No detailed terminal UI. Operator identity/control, timeout and sensitive-command visibility are permission-aware; closure never hides outstanding unknown states.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Operator Control Context | Endpoint | Investigate/Audit | actor/control/time explicit |
| Session Closure Context | Endpoint | Investigate/Govern | closure != state restored |
| control/audit events | refs | CAP-EPT-064/Shared Trace | request/confirmation preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-057 | operator/control event | CAP-EPT-063 | session/actor/state | session |
| CAP-EPT-059..062 | outstanding control | CAP-EPT-063 | attempt/control/status | attempt |
| CAP-EPT-063 | close/audit | CAP-EPT-064/Investigate | closure/outstanding/provenance | session origin |

## 18. Dépendances
CAP-EPT-057..064, Endpoint session-management/local-audit, Investigate CAP-INV-209/214, Govern, Shared Trace, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT target-side operator/session controls; Investigate SOT business participants/context; Govern SOT authority/rollback; Shared SOT generic trace mechanism.

## 20. Provenance et audit
Session, target, operator identity/role, control ownership, takeover/reconnect/timeout/close/cancel/stop requests and confirmations, sensitive visibility decisions, outstanding attempts, timestamps and correlation IDs.

## 21. Permissions fonctionnelles
Session/operator control read, takeover/reconnect/close, attempt cancel/stop, sensitive transcript/output, audit/provenance, cross-tenant deny.

## 22. Limites et erreurs
Stop requested ≠ stopped; session closed ≠ target state restored; cancellation ≠ rollback; operator ≠ Govern approver automatically; forced closure does not prove attempt termination.

## 23. Métriques
Takeovers, disconnect/reconnect, timeout/closures, closure-with-unknown-operations, control denials, audit completeness.

## 24. Classification de livraison
`draft / defined / planned`; no collaboration/session-control protocol or rollback implementation.

## 25. Critères d’acceptation
**Given** stop is requested but target acknowledgement is missing, **When** session closes, **Then** closure records `unknown` attempt termination rather than stopped.

**Given** a session closes normally, **When** closure context is produced, **Then** it does not claim target state restoration or Govern rollback.

**Given** an unauthorized operator attempts takeover, **When** control is evaluated, **Then** takeover is denied and current control context remains auditable.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final concurrency/control bridge unresolved.

## 27. Consommateurs documentaires
Investigate Live Session, Govern, Shared Trace/Audit, CAP-EPT-064, Security/Quality, EPT-5 boundary.
