---
id: CAP-EPT-084
title: Update Download, Staging and Readiness
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-084 — Update Download, Staging and Readiness

## 1. Définition
Définir les faits locaux de disponibilité, téléchargement, staging et readiness d’un update Endpoint sans définir transport, stockage physique ou mécanisme d’installation.

## 2. Problème utilisateur
Un package téléchargé peut être indisponible pour installation à cause de ressources, dépendances, santé ou incompatibilités.

## 3. Objectifs
Représenter cible locale, package availability, download state, staged state, resource/dependency readiness, pre-install health, blocking reason et provenance.

## 4. Non-objectifs
Aucun CDN, protocole, package manager, format, stockage physique, installation, activation, support claim ou scheduler Fleet.

## 5. Propriétaire
Endpoint owns local technical download/staging/readiness facts; Settings/release sources retain assignment/package administration.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Release Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Assignment et package context valides, Agent/tenant/environment connus, permissions applicables et état local observable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local target assignment | CAP-EPT-082 | update target | oui | current | no download |
| package/release context | CAP-EPT-083 | package ref | oui | pinned/current | unavailable |
| resource state | Endpoint local health | readiness fact | oui | fresh | readiness unknown |
| dependency/pre-install health | CAP-EPT-008/093 | health context | oui | fresh | blocked/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Update Assignment Projection | Endpoint | target | read |
| Local Package Context | Endpoint | package ref/compatibility | read |
| Agent Health Assessment | Endpoint | pre-install health | read |
| Endpoint Policy | Settings | constraints | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Update Download State | create/transition | Endpoint | local technical fact |
| Update Staging State | create/transition | Endpoint | no physical storage schema |
| Update Readiness Assessment | derive/refresh | Endpoint | source-backed blocking reasons |

## 11. Fonctionnalités
Track package availability; pending/running/completed/failed download; staged/not-staged; resource/dependency readiness; pre-install health; blocking reasons; freshness/provenance.

## 12. Actions utilisateur
Inspect Class 0; readiness assessment Class 1; request bounded download/staging Class 2 or higher according policy/authority; no implicit install/activation.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assess readiness | oui | oui | oui | explain | deterministic checks |
| summarize blockers | oui | oui | oui | oui | structured reasons |
| invent package/readiness | non | interdit | non | interdit | explicit unknown/block |

## 14. États fonctionnels
`package-unavailable`, `download-pending`, `downloading`, `downloaded`, `download-failed`, `staged`, `not-staged`, `readiness-pending`, `ready`, `blocked`, `resource-insufficient`, `dependency-unavailable`, `health-degraded`, `unknown`.

## 15. États d’interface
No Screen ID. Stale/Partial/Offline/Denied preserve last facts and never convert downloaded to ready.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Download/Staging State | Endpoint fact | CAP-EPT-085/086 | downloaded != staged/ready |
| Readiness Assessment | Endpoint derived fact | operator/Settings | blockers explicit |
| provenance | local audit context | CAP-EPT-097 | source/time retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-083 | package available | CAP-EPT-084 | package/compatibility refs | package context retained |
| CAP-EPT-084 | ready | CAP-EPT-085 | staged package/readiness refs | readiness retained |
| CAP-EPT-084 | failure/block | CAP-EPT-086 | reason/progress refs | retry/defer decision context |

## 18. Dépendances
CAP-EPT-008/011/025/082/083/093; Settings Policy/Fleet; Security; OPEN-008/013.

## 19. Source de vérité
Endpoint SOT for local download/staging/readiness facts; external sources retain assignment/package administration.

## 20. Provenance et audit
Agent/target/package refs, attempt, timestamps, progress checkpoints, resource/dependency/health observations, blocker and requester.

## 21. Permissions fonctionnelles
Update state/readiness read, download/staging request, sensitive package metadata read, cross-tenant deny; final RBAC not selected.

## 22. Limites et erreurs
Assigned != downloaded; downloaded != staged; staged != ready; package available != trusted; readiness unknown stays unknown.

## 23. Métriques
Download/stage/readiness state distribution, failure/block reasons, stale assessments, resource/dependency blockers.

## 24. Classification de livraison
`draft / defined / planned`; no downloader, storage engine, transport or implementation.

## 25. Critères d’acceptation
**Given** download completes but disk/resource readiness fails, **When** readiness is assessed, **Then** state is downloaded but not ready.

**Given** package becomes unavailable mid-flow, **When** status refreshes, **Then** failure/unknown is retained without installation inference.

**Given** AI is unavailable, **When** readiness is evaluated, **Then** deterministic/manual checks still work.

## 26. Questions ouvertes
OPEN-008/013 remain open; platform-specific downloader and default governance classification are not selected.

## 27. Consommateurs documentaires
EPT-6 Update lifecycle, Settings, Security, Quality, registers, Roadmap and future implementation contracts.