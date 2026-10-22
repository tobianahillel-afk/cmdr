---
id: CAP-EPT-081
title: Local Account Session Lock and Termination Primitive
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-081 — Local Account Session Lock and Termination Primitive

## 1. Définition
Définir la primitive Endpoint explicitement sourcée pour verrouiller ou terminer une **session locale** associée à un principal/account context sur l’endpoint, sans administrer le compte d’annuaire ni revendiquer credential revocation externe.

## 2. Problème utilisateur
`account containment` peut être interprété comme désactivation directory globale. Le corpus Endpoint ne source que local session termination/lock et déclare les directory actions externes.

## 3. Objectifs
Pinner local session/principal reference ; distinguer session lock vs termination ; vérifier current session state/support ; exiger authority ; observer locked/terminated state ; conserver external-directory-action boundary et provenance.

## 4. Non-objectifs
Disable/delete/reset directory account, revoke cloud tokens globally, change password, kill arbitrary process as substitute, choose identity provider, provide OS commands or claim user compromise resolved.

## 5. Propriétaire
Endpoint owns local session-control primitive facts. Settings/identity providers own administrative identity state; Govern owns authority/response outcome; Investigate owns person/account findings.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Identity/Security Reviewer, Govern Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready, exact local session reference and endpoint, requested lock/terminate state, platform capability known, authority current.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local session/principal ref | EPT-3/EPT observations | target | oui | fresh | unresolved |
| lock/terminate request | Govern | effect intent | oui | pinned | reject |
| current session state | Endpoint | precheck | oui | timestamped | unknown |
| capability/policy | Endpoint/Settings | support/restriction | oui | current | unsupported/block |
| authority ref | Govern | authorization | oui | effective | no action |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Decision | Govern | authority/scope | read |
| local user/session observations | Endpoint | session identity/state | read |
| Principal/identity admin projection | Settings | identity ref only | restricted read |
| Endpoint Policy | Settings | local session-control restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Session Control Execution | create/transition | Endpoint | local target only |
| Local Session Control State | observe/update | Endpoint | lock/terminate distinct |
| Technical Outcome | create | Endpoint | not identity verdict/Result |
| local session | lock/terminate effect | endpoint target | directory account unchanged by contract |

## 11. Fonctionnalités
Resolve local session; distinguish local lock/termination; revalidate state; reject missing/stale session; execute only declared platform primitive; observe locked/terminated/partial/fail/unknown; expose directory action as external/not performed.

## 12. Actions utilisateur
Inspect/readiness Class 0/1. Local session lock/termination Class 3 by default and Govern-dependent. Directory actions are prohibited in this Endpoint capability.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve local session | oui | oui | oui | explain | session inventory |
| observe lock/termination | oui | oui | oui | summary | local state check |
| suggest local action | oui | rules/context | oui | suggestion | operator choice |
| authorize/disable account externally | non | no | non | interdit | Govern + external owner |

## 14. États fonctionnels
`session-unresolved`, `unsupported`, `ready`, `lock-requested`, `locked-observed`, `termination-requested`, `terminated-observed`, `already-ended`, `partial`, `failed`, `unknown`, `directory-action-external`, `drifted`.

## 15. États d’interface
No Screen ID. Local session, principal/account identity and external directory state remain visually/conceptually separate in any future UI.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Local Session Control State | technical fact | CAP-EPT-073/075 | exact session/state/time |
| technical outcome | raw outcome | Govern | no directory-state claim |
| external-action boundary | handoff context | Govern/Settings | no implicit identity mutation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | local session action ready | CAP-EPT-081 | endpoint/session/op/authority | Run |
| CAP-EPT-081 | outcome | CAP-EPT-073/076 | local session state | same Run |
| external identity action required | boundary reached | Govern/Settings/provider | identity ref + need only | Endpoint does not execute |

## 18. Dépendances
CAP-EPT-020/040/043/050/058/062/065/066/073..080, Govern, Settings identity, Security, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of local session-control technical facts. Settings/provider remains SOT of administrative/directory identity state. Govern is SOT of authority and response outcome.

## 20. Provenance et audit
Endpoint/session/principal refs, local vs external scope, operation, current/requested/observed states, authority/Run/Step, Policy, operator, Agent/version, timestamps, errors, verification, correlation.

## 21. Permissions fonctionnelles
Local session-control read/request, lock, terminate, sensitive session context, verification, cross-tenant deny, step-up/SoD/Govern dependency. No directory-admin permission is implied.

## 22. Limites et erreurs
Local session termination ≠ directory account disable; lock ≠ credential revocation; session ended ≠ user compromise resolved; principal ≠ person automatically; technical success ≠ Result.

## 23. Métriques
Session unresolved/already-ended, lock/terminate outcomes, unsupported, verification mismatch, attempted directory mutation through Endpoint target zero.

## 24. Classification de livraison
`draft / defined / planned`; no OS/session command, directory integration, provider API or implementation.

## 25. Critères d’acceptation
**Given** a local session exists and authority is valid, **When** lock is requested, **Then** only the local session primitive is targeted and directory account state is untouched.

**Given** the session ended before termination, **When** precheck refreshes, **Then** `already-ended` is recorded and no other session is substituted.

**Given** a requested response needs directory account disable, **When** this capability is evaluated, **Then** it returns `directory-action-external` and routes ownership outward rather than executing it.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; exact supported session types/platforms and cross-product identity-response bridge are undecided.

## 27. Consommateurs documentaires
EPT-5 verification/provenance, Govern, Settings identity/admin, Security, Investigate, Quality.