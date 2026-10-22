---
id: CAP-EPT-001
title: Endpoint Agent Identity and Instance Registration
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-001 — Endpoint Agent Identity and Instance Registration

## 1. Définition
Define one Endpoint Agent instance identity, endpoint association and local registration state without defining enrollment credentials, PKI, keys or transport.

## 2. Problème utilisateur
Without a stable Agent-instance identity, an installation, Endpoint, Device, Host Identity and Agent instance can be confused, and duplicate or superseded registrations can appear authoritative.

## 3. Objectifs
Stable instance reference; endpoint association; first/current state; revoked/superseded state; duplicate-registration candidate; tenant/environment and provenance.

## 4. Non-objectifs
No token/certificate/key format, protocol, Fleet administration, Device schema, identity auto-merge, Screen ID or implementation.

## 5. Propriétaire
Endpoint Agent owns local Agent identity/state. Settings retains Endpoint/Fleet administration; no foreign ownership moves.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator as consumer, SOC/Investigate analyst and Auditor.

## 7. Conditions d’entrée
Resolvable Agent observation, tenant/environment scope, source reference and read permission. Missing/restricted facts remain explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Agent instance observation | Endpoint local state | identity fact | oui | current observation | unregistered/unknown |
| Endpoint association ref | shared Endpoint/Settings projection | context ref | si disponible | source freshness | unresolved association |
| tenant/environment refs | Platform Settings | scope refs | oui | current projection | registration untrusted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | id, version, local lifecycle | existing read/manage family |
| tenant/environment | Platform Settings | scoped references | read projection |
| Endpoint shared model | Settings-administered shared model | association only | read if available |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| endpoint-agent | record identity/registration transition | Endpoint Agent | no Device/Endpoint redefinition |
| local-audit-event | append identity provenance | Endpoint Agent | append-only, no raw secret |

## 11. Fonctionnalités
Instance identifier; identity status; first seen; endpoint association; current registration; revoked/superseded state; duplicate candidate; provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect identity | authorized user | endpoint-agent | 0 | read | identity/provenance | non |
| validate identity consistency | Endpoint Operator | endpoint-agent | 1 | sources present | duplicate/mismatch assessment | non |
| request local registration handoff | authorized operator | endpoint-agent | 2 | Settings ref + manage need | bounded request/state change | OPEN-013 if applicable |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| identity inspection | oui | oui | oui | sourced summary only | structured source review |
| duplicate assessment | oui | oui | oui | explanation only | deterministic comparison |
| invent identity/authority | non | validation only | non | interdit | explicit human/source path |

## 14. États fonctionnels
`unregistered`, `registration-pending`, `registered`, `superseded`, `revoked`, `identity-unknown`.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale preserve scope, last valid source and explicit missing facts; no Endpoint screen is defined.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent Identity Projection | endpoint-agent projection | Settings, Investigate, Command, Govern | identity/status/source explicit |
| Registration Provenance | local-audit-event | Audit | history retained, no raw secret |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local instance observation | identity sufficient | registered identity | Agent/tenant/env/source refs | origin retained |
| registered identity | duplicate/supersession evidence | candidate/superseded | both refs + evidence | no silent merge |
| registered identity | revocation observed | revoked | revocation ref/time | history retained |

## 18. Dépendances
Canonical endpoint-agent, Settings tenant/environment/enrollment administration, local-audit-event, Security, OPEN-008. EPT-2..6 remain NOT STARTED.

## 19. Source de vérité
Endpoint Agent owns local Agent-instance identity/registration. Endpoint/Device and Settings administrative records remain referenced projections.

## 20. Provenance et audit
Record Agent id, association, tenant/environment, first/current timestamps, transitions, source, duplicate/supersession evidence and correlation id; never raw enrollment secrets.

## 21. Permissions fonctionnelles
Reuse `perm.endpoint-agent.endpoint-agent.read/manage` and local-audit-event read/manage. Restricted identity/cross-tenant needs are documented, not atomically finalized.

## 22. Limites et erreurs
Duplicate, stale identity, revoked state, tenant/environment mismatch, missing association and permission denial stay explicit; no fallback creates identity equivalence.

## 23. Métriques
Registered/unknown/revoked/superseded counts; duplicate candidates; unresolved associations; freshness and denied cross-tenant attempts.

## 24. Classification de livraison
`draft / defined / planned`. Documentary provider-neutral capability only; no software, platform support, API/protocol, physical schema or final RBAC.

## 25. Critères d’acceptation
**Given** two Agent instances reference the same Endpoint, **When** consistency is checked, **Then** both remain distinct and a duplicate candidate is surfaced.

**Given** an Agent is revoked while offline, **When** state is read, **Then** revoked is preserved and is not converted to deleted/offline.

**Given** AI is unavailable, **When** identity is validated, **Then** the essential path remains deterministic/manual.

## 26. Questions ouvertes
OPEN-008 remains open. No platform/source support or new OPEN decision is resolved by this capability.

## 27. Consommateurs documentaires
Endpoint foundations, Settings Fleet/enrollment, Investigate/Command/Govern projections, registers, Quality, Roadmap Phase 5 and future EPT lots.