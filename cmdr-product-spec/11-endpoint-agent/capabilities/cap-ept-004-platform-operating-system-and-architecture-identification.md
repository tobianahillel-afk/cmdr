---
id: CAP-EPT-004
title: Platform, Operating System and Architecture Identification
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
# CAP-EPT-004 — Platform, Operating System and Architecture Identification

## 1. Définition
Define provider-neutral observation of platform family, OS name/version and architecture, with unknown/unsupported/not-evaluated support states and no claim that observed means supported or delivered.

## 2. Problème utilisateur
Detected OS/architecture is necessary context but can be mistaken for an officially supported delivery target if observation and support decisions are collapsed.

## 3. Objectifs
Observe platform/OS/architecture; optional kernel/build or virtualization/container facts only when sourced; explicit unknown/unsupported/not-evaluated state; source/freshness.

## 4. Non-objectifs
No Windows/Linux/macOS support declaration, distro/version selection, driver/collector design, cloud/container commitment or platform implementation.

## 5. Propriétaire
Endpoint owns observed technical facts. Official support remains a release/product decision under OPEN-008.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC/Investigate analyst, compatibility reviewer and Auditor.

## 7. Conditions d’entrée
Agent identity, authorized local platform source and tenant/environment context. Missing source yields unknown.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| platform observation | Endpoint local state | platform fact | oui | observation freshness | unknown |
| OS name/version | Endpoint local state | OS facts | si disponible | observation freshness | unknown |
| architecture | Endpoint local state | architecture fact | si détectable | observation freshness | unknown |
| support decision ref | release/OPEN-008 evidence | support context | non | decision freshness | not-evaluated |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | Agent context | read |
| environment | Platform Settings | environment context | read projection |
| Endpoint shared model | Settings-administered shared model | target context | read if available |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| endpoint-agent | record observed platform facts | Endpoint Agent | observation does not create support policy |
| local-audit-event | record fact/source change | Endpoint Agent | factual provenance only |

## 11. Fonctionnalités
Platform family; OS name/version; architecture; sourced kernel/build; sourced virtualization/container context; support-status representation; freshness.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect platform facts | authorized user | endpoint-agent | 0 | read | source-backed facts | non |
| normalize observed facts | Endpoint deterministic logic | endpoint-agent | 1 | source present | labels + unknowns | non |
| assess support context | reviewer | endpoint-agent | 1 | explicit decision if any | support/unsupported/not-evaluated | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| fact inspection | oui | oui | oui | sourced summary | direct source review |
| normalization | oui | oui | oui | explanation only | deterministic mapping |
| invent platform/support | non | validation only | non | interdit | unknown/not-evaluated state |

## 14. États fonctionnels
`observed`, `partially-observed`, `unknown`, `support-not-evaluated`, `unsupported`.

## 15. États d’interface
Partial/Stale show missing/source age; Permission denied masks restricted metadata; no final screen/layout.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Platform Facts | endpoint-agent projection | Inventory/Health/Settings/Investigate | observed/source/freshness explicit |
| Support Context | derived state | compatibility consumers | no invented support |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local platform source | facts observed | platform facts | platform/OS/arch/source/time | origin retained |
| facts | source stale/changes | partial/unknown | old/new refs | history retained |
| no support decision | support requested | not-evaluated | OPEN-008 ref | no delivery claim |

## 18. Dépendances
Endpoint architecture/platform-support, Settings environment, CAP-EPT-001/003, OPEN-008.

## 19. Source de vérité
Endpoint owns observed facts; official support/delivery is not inferred and remains unresolved under OPEN-008.

## 20. Provenance et audit
Record source, timestamp, platform/OS/architecture values, normalization context, previous ref and support-decision ref when one exists.

## 21. Permissions fonctionnelles
Endpoint read covers ordinary platform facts; restricted metadata needs remain future. No platform-specific permission created.

## 22. Limites et erreurs
Observed != supported; referenced != delivered; OS identified != compatible; compatible != supported; missing source = unknown.

## 23. Métriques
Fact completeness; unknown/not-evaluated counts; freshness; source conflicts/support-context mismatches.

## 24. Classification de livraison
`draft / defined / planned`; platform-neutral/provider-neutral, no implementation or support claim.

## 25. Critères d’acceptation
**Given** Windows is observed, **When** OPEN-008 has no release decision, **Then** Windows is recorded while support remains not-evaluated.

**Given** platform source is unavailable, **When** identification runs, **Then** facts remain unknown rather than inferred.

**Given** an OS version is identified, **When** compatibility is unassessed, **Then** no compatibility/support claim is made.

## 26. Questions ouvertes
OPEN-008 remains open for platforms, distributions, versions and source support.

## 27. Consommateurs documentaires
Endpoint inventory/health/capabilities, Settings Fleet, Investigate/Command/Govern consumers, registers, Quality, Roadmap and future EPT lots.