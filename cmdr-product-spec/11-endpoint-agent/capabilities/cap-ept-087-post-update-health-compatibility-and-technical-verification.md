---
id: CAP-EPT-087
title: Post-Update Health, Compatibility and Technical Verification
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-087 — Post-Update Health, Compatibility and Technical Verification

## 1. Définition
Vérifier techniquement après update la version observée, la santé Agent, les capabilities et dépendances, et la compatibilité, sans conclure support officiel ni succès métier.

## 2. Problème utilisateur
Une installation/activation techniquement réussie peut laisser un Agent dégradé, incompatible ou sur une version inattendue.

## 3. Objectifs
Compare expected/observed version, Agent health, capability/dependency health, compatibility, mismatch/degraded/unknown and verification limits with provenance.

## 4. Non-objectifs
Aucun support certification, security posture proof, business success, Govern Result, automatic rollback, SLO final ou implementation.

## 5. Propriétaire
Endpoint owns technical post-update observations/verification; release/support owner retains official support and Govern retains response-level verification/Result.

## 6. Utilisateurs
Endpoint Operator, Release Reviewer, Platform Administrator, Security Reviewer, Auditor, Govern consumer where relevant.

## 7. Conditions d’entrée
Completed/activated or failed update attempt, expected target/version, current observations, health/capability/dependency facts and compatibility source.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| update operation/outcome | CAP-EPT-085/086 | technical history | oui | exact | unverifiable |
| expected version | CAP-EPT-082/083 | target | oui | pinned/current | mismatch unknown |
| observed version/health | CAP-EPT-005/008 | local facts | oui | post-update fresh | unknown |
| capability/dependency state | CAP-EPT-011/093 | technical state | oui where available | fresh | partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Update Operation | Endpoint | target/outcome | read |
| Agent Health Assessment | Endpoint | post-update health | read |
| Technical Capability Availability | Endpoint | capability state | read |
| release/support source | external owner | compatibility context | read ref |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Post-Update Verification | create/version | Endpoint | technical scope only |
| Version Match State | derive | Endpoint | expected vs observed |
| Post-Update Compatibility State | derive | Endpoint | compatible != supported |

## 11. Fonctionnalités
Compare versions; assess Agent health; recheck capability/dependency availability; assess compatibility; preserve degraded/mismatch/unknown; record limits and recommended recovery/review path.

## 12. Actions utilisateur
Inspect Class 0; deterministic verification Class 1; request fresh self-check/reassessment Class 2 if supported; no autonomous reversion or support declaration.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare expected/observed | oui | oui | oui | explain | deterministic comparison |
| assess sourced health/capability | oui | oui | oui | summarize | structured checks |
| declare supported/business success | non | external/Govern only | non | interdit | explicit source owner workflow |

## 14. États fonctionnels
`verification-pending`, `verified-technical`, `version-mismatch`, `degraded`, `incompatible-context`, `compatibility-unknown`, `dependency-unhealthy`, `capability-degraded`, `unverifiable`, `stale`, `superseded`.

## 15. États d’interface
No Screen ID. Verified technical never renders as supported/secure/business-success automatically.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Post-Update Verification | Endpoint fact | Settings/Quality/CAP-EPT-088 | scope/limits explicit |
| health/capability delta | technical observation | operators/consumers | before/after source refs |
| compatibility state | derived fact | release reviewer | support not inferred |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-085/086 | terminal update state | CAP-EPT-087 | target/attempt/outcome | operation ref |
| CAP-EPT-087 | healthy/verified | Settings projection | observed version/health | no policy mutation |
| CAP-EPT-087 | mismatch/degraded | CAP-EPT-088/manual review | verification facts/limits | recovery context |

## 18. Dépendances
CAP-EPT-005/008/011/082..086/088/093; Settings; Govern verification boundary; OPEN-008.

## 19. Source de vérité
Endpoint SOT for technical verification observations; external release/support and Govern remain owners of their conclusions.

## 20. Provenance et audit
Update attempt, expected/observed version, health/capability/dependency observations, compatibility source, timestamps, limitations and reviewer/requester.

## 21. Permissions fonctionnelles
Verification/health/update-state read, self-check request if supported, sensitive diagnostics read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Activated != healthy; technical update success != compatible; compatible != supported; health PASS != uncompromised; unknown remains explicit.

## 23. Métriques
Version mismatch, degraded health, capability/dependency regressions, compatibility unknown/incompatible and unverifiable counts.

## 24. Classification de livraison
`draft / defined / planned`; no verification engine/support certification/implementation.

## 25. Critères d’acceptation
**Given** update activates but a required sensor is degraded, **When** verification runs, **Then** post-update state is degraded rather than healthy.

**Given** observed version differs from expected, **When** verified, **Then** mismatch is recorded even if installer returned success.

**Given** AI is unavailable, **When** verification executes, **Then** deterministic/manual checks remain available.

## 26. Questions ouvertes
OPEN-008 remains open; official support and platform-specific health criteria are not selected.

## 27. Consommateurs documentaires
EPT-6 Update/Resilience, Settings Fleet, Security, Govern consumer boundary, Quality, registers and Roadmap.