---
id: CAP-EPT-002
title: Enrollment Handoff and Endpoint Enrollment State
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-002 — Enrollment Handoff and Endpoint Enrollment State

## 1. Définition
Define Endpoint-side consumption of a Settings-owned enrollment handoff and local enrollment state without moving credential issuance, enrollment administration or revocation authority out of Settings.

## 2. Problème utilisateur
Administrative enrollment may be approved while the local Agent is still pending, rejected, mismatched or offline; configured intent must not be mistaken for technical completion.

## 3. Objectifs
Consume enrollment reference; represent enrolling/enrolled/rejected/failed/revoked; expose tenant/environment mismatch; preserve unenrollment/revocation local effect and history.

## 4. Non-objectifs
No token/certificate format, mutual-auth protocol, Fleet membership administration, credential lifecycle, retry implementation or transport.

## 5. Propriétaire
Platform Settings owns enrollment administration. Endpoint Agent owns only local technical enrollment state and its provenance.

## 6. Utilisateurs
Platform Administrator, Endpoint Operator, Security/Audit reviewer and authorized Investigate/Govern consumer.

## 7. Conditions d’entrée
Settings enrollment reference, tenant/environment assignment and local Agent identity. Raw credential values are never required by this capability document.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| enrollment reference | Platform Settings | authorized handoff ref | oui | current admin version | pending/blocked |
| tenant/environment assignment | Platform Settings | scope | oui | current | mismatch/reject |
| local registration state | Endpoint Agent | technical state | oui | current | unknown |
| revocation/unenrollment ref | Platform Settings | admin intent | si applicable | current event | no revocation conclusion |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | registration/enrollment state | read/manage |
| endpoint-agent-fleet | Platform Settings | enrollment context only | read projection |
| tenant/environment | Platform Settings | assignments | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| endpoint-agent | local enrollment transition | Endpoint Agent | admin record remains Settings-owned |
| local-audit-event | append handoff/outcome | Endpoint Agent | no raw token/certificate |

## 11. Fonctionnalités
Handoff consumption; enrolling/enrolled/rejected/failed/revoked; mismatch detection; local unenrollment effect; retry-candidate concept; provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect local enrollment | authorized user | endpoint-agent | 0 | read | local + admin-ref state | non |
| validate handoff scope | Endpoint Operator | endpoint-agent | 1 | refs present | eligible/mismatch/blocked | non |
| request/ack handoff | authorized operator | endpoint-agent | 2 | manage + Settings authority | bounded local transition | OPEN-013; no admin transfer |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| handoff inspection | oui | oui | oui | summary only | source-table review |
| scope validation | oui | oui | oui | explanation only | deterministic refs comparison |
| enroll without authority | non | validation only | non autonome | interdit | Settings-authorized path |

## 14. États fonctionnels
`not-enrolled`, `enrollment-pending`, `enrolling`, `enrolled`, `rejected`, `failed`, `revoked`.

## 15. États d’interface
Partial distinguishes admin intent from local completion; Offline does not erase state; Permission denied masks restricted refs; Stale shows source version/freshness. No screen design.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Local Enrollment State | endpoint-agent projection | Settings Fleet + consumers | admin intent != local completion |
| Enrollment Audit | local-audit-event | Audit | source refs/outcome retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings enrollment admin | handoff consumable | pending/enrolling | enrollment ref + scope | Settings record retained |
| enrolling | local prerequisites complete | enrolled | completion evidence | admin owner retained |
| enrolled/pending | revocation observed | revoked/not-enrolled effect | source ref + time | history retained |

## 18. Dépendances
Settings enrollment, CAP-EPT-001, tenant/environment, Security, OPEN-008/013. No EPT-2+ execution.

## 19. Source de vérité
Settings owns enrollment administration and credential-bearing artifacts; Endpoint owns local enrollment state only.

## 20. Provenance et audit
Record enrollment ref, scope, local transition, reason, timestamps, mismatch evidence and correlation id; no token/key/certificate values.

## 21. Permissions fonctionnelles
Endpoint read/manage for local state; Settings Fleet/enrollment permissions remain external. No atomic enrollment permission is created.

## 22. Limites et erreurs
Admin approved != enrolled; enrolled != online; revoked != deleted; stale/mismatch/offline/permission denial remain explicit.

## 23. Métriques
Pending/enrolled/rejected/failed/revoked counts; mismatch count; descriptive handoff-to-local-completion duration; retry candidates.

## 24. Classification de livraison
`draft / defined / planned`; no PKI, protocol, credential implementation, platform delivery or final RBAC.

## 25. Critères d’acceptation
**Given** Settings approves enrollment but local completion is absent, **When** state is read, **Then** it remains pending/enrolling rather than enrolled.

**Given** tenant/environment differs, **When** scope is validated, **Then** the handoff is blocked with mismatch provenance.

**Given** revocation is received, **When** local state reconciles, **Then** revoked is recorded while history remains.

## 26. Questions ouvertes
OPEN-008 and OPEN-013 remain open; EPT-1 selects no credential, protocol or default class-2 policy.

## 27. Consommateurs documentaires
Endpoint identity/binding, Settings enrollment/Fleet, Audit, registers, Quality, Roadmap and future technical enrollment contracts.