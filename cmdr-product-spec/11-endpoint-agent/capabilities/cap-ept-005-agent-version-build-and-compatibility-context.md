---
id: CAP-EPT-005
title: Agent Version, Build and Compatibility Context
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
# CAP-EPT-005 — Agent Version, Build and Compatibility Context

## 1. Définition
Define observation of Endpoint Agent version/build and a bounded compatibility assessment against explicit source-owned expectations without treating an observed build as approved, trusted, supported or upgrade-eligible automatically.

## 2. Problème utilisateur
Operators need to distinguish what is installed from what Settings expects and what a release decision supports; collapsing those states can create false support or upgrade claims.

## 3. Objectifs
Expose observed version/build, installation-version state, expected/admin-target projection when available, compatibility state, outdated candidate, mismatch, source, freshness and provenance.

## 4. Non-objectifs
No upgrade lifecycle, package format, signing/PKI design, download channel, rollback implementation, support declaration, physical schema or final RBAC.

## 5. Propriétaire
Endpoint Agent owns observed local version/build facts and compatibility context. Platform Settings owns upgrade waves/admin targets; release/support decisions remain external.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, compatibility/release reviewer, SOC/Investigate analyst and Auditor.

## 7. Conditions d’entrée
Agent identity, observed local version/build, tenant/environment scope and any explicit expected-version or compatibility source reference.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| observed Agent version | Endpoint local state | version fact | oui | current observation | version unknown |
| observed Agent build | Endpoint local state | build fact | selon source | current observation | build unknown |
| expected/admin target | Platform Settings projection | target context | non | current admin version | no mismatch conclusion |
| compatibility evidence | release/support source | assessment context | non | source-defined | compatibility unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | version/build/current state | existing read |
| endpoint-agent-fleet | Platform Settings | expected/admin target only | Settings read projection |
| environment | Platform Settings | target environment context | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| endpoint-agent | record observed version/build/compatibility state | Endpoint Agent | no upgrade-wave mutation |
| local-audit-event | append observation/assessment provenance | Endpoint Agent | factual trace, no package/secret |

## 11. Fonctionnalités
Observed version/build; expected target projection; compatibility known/unknown/supported-context/unsupported-context; outdated candidate; version/build mismatch; installation presence context; freshness/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect version/build | authorized user | endpoint-agent | 0 | read | observed facts | non |
| compare observed vs expected | Endpoint Operator | endpoint-agent | 1 | both refs available | match/mismatch/outdated candidate | non |
| assess compatibility | authorized reviewer | endpoint-agent | 1 | explicit compatibility source | supported-context/unsupported/unknown | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect version/build | oui | oui | oui | sourced summary only | direct fact review |
| compare versions | oui | oui | oui | explanation only | deterministic comparison |
| approve/trust/support a build | accountable source owner | source rules only | non autonome | interdit | explicit release/support decision |

## 14. États fonctionnels
`version-known`, `version-unknown`, `match`, `mismatch`, `outdated-candidate`, `compatibility-unknown`, `compatible-context`, `incompatible-context`.

## 15. États d’interface
Loading retains Agent/scope; Partial names missing build/target; Error preserves last valid observation; Offline permits last-known display with freshness; Permission denied masks restricted target data; Stale forces re-evaluation. No Endpoint Screen ID.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Version/Build Projection | endpoint-agent projection | Settings Fleet, inventory, health, consumers | observed vs expected remain distinct |
| Compatibility Assessment | derived fact | Settings/Endpoint consumers | source and unknown state explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local Agent observation | version/build read | version projection | values/source/time | Agent context |
| expected target + observed version | compare | match/mismatch candidate | both refs | no upgrade action |
| compatibility source | assessment | compatible/incompatible/unknown | source/version/platform refs | no support promotion |

## 18. Dépendances
CAP-EPT-001/003/004, Settings Fleet/upgrade-management projections, platform-support, Security and OPEN-008.

## 19. Source de vérité
Endpoint is source for observed local version/build; Settings owns admin targets/upgrade waves; support/approval is not inferred from observation.

## 20. Provenance et audit
Record Agent id, observed version/build, source/time, expected target ref, comparison result, compatibility source and any stale/unknown reason; keep history.

## 21. Permissions fonctionnelles
Reuse Endpoint Agent read for observed data; Settings target access remains Settings-owned. Sensitive/restricted version metadata needs are recorded without creating arbitrary atomic permissions.

## 22. Limites et erreurs
Observed version != approved version; build observed != trusted automatically; compatible != supported; missing target cannot become mismatch; EPT-1 initiates no upgrade.

## 23. Métriques
Known/unknown versions; build mismatches; outdated candidates; compatibility unknown/incompatible counts; freshness and source-conflict counts.

## 24. Classification de livraison
`draft / defined / planned`. Documentary context only; no upgrade engine, support declaration, package delivery, API/protocol or implementation.

## 25. Critères d’acceptation
**Given** an Agent version differs from a Settings target, **When** compared, **Then** mismatch/outdated candidate is shown without starting an upgrade.

**Given** a build is observed but no compatibility source exists, **When** assessed, **Then** compatibility remains unknown rather than supported.

**Given** AI is unavailable, **When** versions are compared, **Then** the comparison remains deterministic/manual.

## 26. Questions ouvertes
OPEN-008 remains open for supported platforms/versions and release evidence; upgrade lifecycle remains future EPT-6/Settings work.

## 27. Consommateurs documentaires
Endpoint inventory/health/capability availability, Settings Fleet/upgrade management, registers, Quality, Roadmap and future implementation contracts.