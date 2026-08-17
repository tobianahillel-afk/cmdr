---
id: CAP-EPT-083
title: Update Package, Release Reference and Compatibility Context
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-083 — Update Package, Release Reference and Compatibility Context

## 1. Définition
Définir le contexte fonctionnel d’un package/release d’update Endpoint et son assessment de compatibilité sans choisir format, signature, transport ou support platform final.

## 2. Problème utilisateur
Un package disponible ou une version connue peut être pris à tort pour un artefact trusted, compatible ou officiellement supporté.

## 3. Objectifs
Conserver package/release ref, version/build, platform requirements, prerequisites, compatibility source, available/unsupported/unknown states et provenance.

## 4. Non-objectifs
Pas de package format, signature algorithm, certificate, CDN, downloader, support decision, API/protocol ou implementation.

## 5. Propriétaire
Endpoint owns local package/reference and compatibility consumption facts; release/package administration and support evidence remain source-owned externally.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Release Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Assignment CAP-EPT-082 or explicit release ref, observed platform/version, compatibility source and tenant/environment context.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| package/release ref | Settings/release source | immutable/versioned ref | oui for update | pinned | no package context |
| target version/build | release source | version context | oui | pinned | unknown |
| platform requirements | release/support source | compatibility input | non | source-defined | compatibility unknown |
| local platform/version | CAP-EPT-004/005 | observed facts | oui | current | assessment unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent | Endpoint | platform/version | local read |
| Environment/Fleet | Settings | target context | read projection |
| Update Package Reference | release/admin source | metadata ref only | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Package Context | create/refresh | Endpoint | metadata/ref, not contents |
| Compatibility Assessment | derive | Endpoint | source-backed; unknown explicit |
| external package/release | aucune mutation | external owner | reference only |

## 11. Fonctionnalités
Resolve package/release refs; compare platform/version prerequisites; expose compatible/incompatible/unknown/unavailable; retain integrity/authenticity requirement conceptually without defining crypto.

## 12. Actions utilisateur
Inspect package metadata Class 0; assess compatibility Class 1; request metadata refresh Class 2 if allowed. No trust declaration or package execution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare prerequisites | oui | oui | oui | explain | deterministic rules |
| summarize release context | oui | structured | oui | oui | metadata table |
| declare package trusted/support | external owner | explicit evidence | non | interdit | source decision |

## 14. États fonctionnels
`available`, `unavailable`, `metadata-partial`, `compatible-context`, `incompatible-context`, `compatibility-unknown`, `unsupported-context`, `stale`, `superseded`.

## 15. États d’interface
No Screen ID. Unknown compatibility never renders as supported; sensitive refs may be masked.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| package context | Endpoint projection | CAP-EPT-084..088 | package ref != contents |
| compatibility assessment | derived fact | readiness/health consumers | source and limitations explicit |
| provenance | audit context | Quality/Audit | version/source retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-082 | target assigned | CAP-EPT-083 | target/release refs | assignment retained |
| CAP-EPT-083 | compatible context | CAP-EPT-084 | package ref + prerequisites | compatibility ref retained |
| compatibility conflict | reassessment | Settings/release consumer | observed facts + reason | no support mutation |

## 18. Dépendances
CAP-EPT-004/005/082; Settings upgrade management; Endpoint update contract; Security; OPEN-008.

## 19. Source de vérité
Package/release source owns release metadata; Endpoint owns local assessment only.

## 20. Provenance et audit
Package/release ref, target/current version, platform facts, prerequisites, compatibility source/version, assessment result, time and masking.

## 21. Permissions fonctionnelles
Package metadata read, compatibility read/assess, sensitive ref masking, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Package available != trusted; package ref != contents; integrity metadata != crypto proof; compatible != officially supported.

## 23. Métriques
Unknown/incompatible/unavailable package contexts, stale metadata, prerequisite failures and source conflicts.

## 24. Classification de livraison
`draft / defined / planned`; no package or trust implementation.

## 25. Critères d’acceptation
**Given** a package ref exists but platform requirements are missing, **When** compatibility is assessed, **Then** state remains unknown.

**Given** compatibility rules reject the observed architecture, **When** readiness is requested, **Then** incompatible context is surfaced without download.

**Given** AI is unavailable, **When** prerequisites are compared, **Then** deterministic assessment remains available.

## 26. Questions ouvertes
OPEN-008 remains open; no supported platform/version release is declared.

## 27. Consommateurs documentaires
EPT-6 Update, Settings, Security, implementation contracts, registers, Quality and Roadmap.