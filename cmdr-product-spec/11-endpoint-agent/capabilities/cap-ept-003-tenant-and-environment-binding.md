---
id: CAP-EPT-003
title: Tenant and Environment Binding
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
# CAP-EPT-003 — Tenant and Environment Binding

## 1. Définition
Define Endpoint-local observed/effective tenant and environment binding as references administered by Platform Settings, including mismatch, unknown, stale and cross-tenant-denied states.

## 2. Problème utilisateur
Ambiguous binding can project Agent facts into the wrong tenant/environment or create an implicit cross-tenant fallback.

## 3. Objectifs
Explicit tenant/environment refs; current binding source; mismatch/unknown/stale states; cross-tenant denial; provenance and freshness.

## 4. Non-objectifs
No tenant/environment administration, federation, authorization-policy design, routing or physical identity schema.

## 5. Propriétaire
Settings administers tenant/environment; Endpoint owns only local effective binding observation/state.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC/Investigate analyst, Govern reviewer and Auditor.

## 7. Conditions d’entrée
Endpoint Agent identity plus source-owned tenant/environment references and permission to read them.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| tenant reference | Platform Settings | scope ref | oui | current projection | binding unknown |
| environment reference | Platform Settings | scope ref | oui | current projection | binding unknown |
| local Agent identity | Endpoint Agent | technical identity | oui | current | cannot bind |
| binding observation | Endpoint Agent | local state | oui | current | unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | identity/current state | read/manage |
| tenant | Platform Settings | tenant ref | read projection |
| environment | Platform Settings | environment ref | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| endpoint-agent | local binding state | Endpoint Agent | no tenant/env administration |
| local-audit-event | binding/mismatch event | Endpoint Agent | tenant-scoped provenance |

## 11. Fonctionnalités
Current binding, binding source, mismatch, unknown/stale binding, unauthorized cross-tenant condition, provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect binding | authorized user | endpoint-agent | 0 | read | scoped projection | non |
| validate consistency | Endpoint Operator | endpoint-agent | 1 | refs present | match/mismatch/unknown | non |
| request bounded refresh | authorized operator | endpoint-agent | 2 | source refs + manage | refresh request only | no admin transfer |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| binding inspection | oui | oui | oui | summary only | direct refs review |
| mismatch detection | oui | oui | oui | explanation only | deterministic comparison |
| cross-tenant fallback | non | deny rule | non | interdit | explicit authorized scope only |

## 14. États fonctionnels
`bound`, `binding-mismatch`, `binding-unknown`, `binding-stale`, `cross-tenant-denied`.

## 15. États d’interface
Partial/Stale/Permission denied name the exact missing or restricted scope; Offline does not change tenant ownership. No Screen ID.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Binding State | endpoint-agent projection | Settings/Investigate/Command/Govern | refs/source/freshness explicit |
| Binding Audit | local-audit-event | Audit | mismatch/denial retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings refs | local refs match | bound | refs/source/time | ownership retained |
| bound | refs diverge | mismatch | both refs + freshness | no reassignment |
| any state | cross-tenant attempt | denied | origin/requested tenant | no fallback |

## 18. Dépendances
Settings tenant/environment, CAP-EPT-001/002, Security tenant isolation, OPEN-008.

## 19. Source de vérité
Settings owns tenant/environment administration; Endpoint owns local binding state only.

## 20. Provenance et audit
Record Agent id, tenant/environment refs, source/version when known, observed state, timestamp, mismatch/denial reason and correlation id.

## 21. Permissions fonctionnelles
Endpoint read/manage for local state; Settings read remains external; cross-tenant read is never implied.

## 22. Limites et erreurs
Unknown/stale bindings remain non-authoritative. Tenant mismatch never falls back; environment mismatch cannot be normalized away.

## 23. Métriques
Bound/mismatch/unknown/stale counts; cross-tenant denials; freshness; descriptive reconciliation duration.

## 24. Classification de livraison
`draft / defined / planned`; provider/platform-neutral and no implementation or final RBAC.

## 25. Critères d’acceptation
**Given** Agent tenant differs from Settings, **When** binding is validated, **Then** mismatch is explicit and no fallback occurs.

**Given** environment reference is stale, **When** read, **Then** stale is distinct from mismatch.

**Given** cross-tenant read is attempted, **When** scope is checked, **Then** access is denied without changing either tenant.

## 26. Questions ouvertes
OPEN-008 remains open; no support/source decision is inferred from binding.

## 27. Consommateurs documentaires
Endpoint foundations, Settings tenant/environment/Fleet, Security, product consumers, registers, Quality and Roadmap.