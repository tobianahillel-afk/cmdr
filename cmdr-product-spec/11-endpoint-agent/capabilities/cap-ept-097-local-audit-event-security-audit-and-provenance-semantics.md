---
id: CAP-EPT-097
title: Local Audit Event, Security Audit and Provenance Semantics
product: endpoint-agent
module: security
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-097 — Local Audit Event, Security Audit and Provenance Semantics

## 1. Définition
Définir le Local Audit Event Endpoint et sa provenance technique, incluant acteur/processus/opération/outcome/temps/masking et limites d’ordre/rétention, sans le transformer en Shared Trace, immutable ledger ou preuve cryptographique.

## 2. Problème utilisateur
Les opérations locales, offline ou sensibles nécessitent une reconstruction factuelle sans sur-promettre immutabilité, ordre global ou intégrité cryptographique.

## 3. Objectifs
Represent local audit event, source, actor/process context, operation, outcome, timestamps, masking, ordering limitations, local-retention reference, offline/sync state and provenance.

## 4. Non-objectifs
Aucun ledger physique, WORM store, cryptographic proof, SIEM, organization-wide Audit Trail, retention engine, schema/API/protocol or Shared Trace replacement.

## 5. Propriétaire
Endpoint owns raw local audit events. Security owns global audit/integrity policy; Shared owns generic Trace/Activity; Govern owns Govern Audit Trail semantics.

## 6. Utilisateurs
Auditor, Endpoint Operator, Security Reviewer, Platform Administrator, Govern/Investigate/Command authorized consumers.

## 7. Conditions d’entrée
Sourced local operation/state transition, actor/process/source context, tenant/environment, time source, masking/classification and local retention/sync context.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local technical event | Endpoint capabilities | audit subject | oui | event time | no invented event |
| actor/process context | local runtime/request | attribution | si disponible | event time | attribution partial |
| masking/classification | Security/Policy | visibility | oui for sensitive fields | current | restrict output |
| queue/sync context | CAP-EPT-089/090 | offline delivery context | non | current | sync unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint technical objects | Endpoint | operation/state refs | read/link |
| Secret Reference | Settings | opaque ref only | masked/restricted |
| Security policy | Security | audit/masking requirements | reference |
| Shared Trace/Activity | Shared | correlation refs only | read if permitted |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Audit Event | create/append conceptually | Endpoint | history preserved; no crypto guarantee |
| Audit Sync/Queue State | link/derive | Endpoint | delivery != event creation |
| Shared/Govern audit objects | aucune mutation | external owner | refs only |

## 11. Fonctionnalités
Record exact local event/source/operation/outcome; preserve actor/process/time; apply masking; retain ordering/sequence limitations; support offline queue/sync refs; correlate cross-product refs without owner transfer; expose gaps/unknown state.

## 12. Actions utilisateur
Inspect local audit Class 0; reconstruct/normalize Class 1; request bounded audit sync/status Class 2 if supported. No audit deletion or provenance destruction.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate/sequence events | oui | oui within sourced refs | oui | summarize | ordered refs/limitations |
| apply masking | oui | oui | oui | non nécessaire | deterministic rules |
| invent/alter/delete history | non | interdit | non | interdit | explicit gap/append semantics |

## 14. États fonctionnels
`recorded`, `queued-for-sync`, `synced-observed`, `sync-unknown`, `ordering-limited`, `attribution-partial`, `masked`, `restricted`, `gap-present`, `stale-reference`, `superseded-context`.

## 15. États d’interface
No Screen ID. Local audit, Shared Trace and Govern Audit Trail remain distinguishable; sensitive fields may be masked without hiding event existence when policy permits.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Local Audit Event | Endpoint fact | Security/Quality/cross-product consumers | source/outcome/time explicit |
| provenance chain refs | technical context | CAP-EPT-099 | owner per hop retained |
| sync/ordering limitation | quality context | audit consumers | no immutable/exact-order claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Endpoint operation/state | auditable event | CAP-EPT-097 | actor/source/object/outcome/time | source ref |
| CAP-EPT-097 | offline retention/reconnect | CAP-EPT-089/090 | event refs/sequence/limits | sync context |
| local audit | consumer handoff | Shared/Govern/Security/Quality | refs/masking/limits | Endpoint ownership retained |

## 18. Dépendances
All Endpoint capability families; CAP-EPT-089/090/094..096/098/099; local-audit source/sync contract; Security audit/privacy; Shared Trace/Activity; OPEN-008.

## 19. Source de vérité
Endpoint SOT for its local audit event facts. Security owns policy, Shared generic Trace/Activity, Govern domain audit reconstruction.

## 20. Provenance et audit
The capability itself preserves source object/event, actor/process, operation, outcome, timestamps, Agent/version, tenant/environment, correlation refs, masking, sequence/ordering limits, retention/sync refs and gaps.

## 21. Permissions fonctionnelles
Local-audit read, sensitive-audit read, provenance read, sync/status request, export preparation according canonical controls, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Local audit event != Shared Trace; append-only semantics != immutable ledger != cryptographic proof; synced != globally retained forever; missing event is never fabricated.

## 23. Métriques
Audit coverage/gaps, attribution partial, masking/restriction, queued/sync-unknown, ordering limitations and provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no audit storage/ledger/crypto/sync implementation.

## 25. Critères d’acceptation
**Given** a sensitive update action occurs, **When** local audit is recorded, **Then** action/outcome and opaque refs are preserved while raw secrets are excluded.

**Given** events are queued offline, **When** later synchronized, **Then** local event history remains distinct from delivery/sync state.

**Given** AI is unavailable, **When** provenance is reconstructed, **Then** deterministic/manual reference traversal remains possible.

## 26. Questions ouvertes
OPEN-008 remains open; platform audit persistence, crypto proof and physical retention mechanisms are not selected.

## 27. Consommateurs documentaires
Endpoint EPT-1..6, Security, Shared, Govern, Investigate/Command, Settings, Quality, registers and Roadmap.