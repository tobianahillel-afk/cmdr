---
id: CAP-STD-048
title: Runtime Context, State and Data Propagation
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-048 — Runtime Context, State and Data Propagation

## 1. Définition
Définit Run context, step-local/shared context, source object refs, variable values, sensitive values/Secret References, outputs, transient state, parent/child Run context, cross-step propagation, masking and lifecycle. Transient context ≠ canonical object or permanent memory.

## 2. Problème utilisateur
Runtime data peut perdre source/provenance, traverser tenant boundaries ou être persistée comme mémoire/secrets sans autorisation.

## 3. Objectifs
Préserver source refs and scope; séparer definitions from runtime values; limiter shared/local context; mask sensitive data; propagate only explicit outputs; manage parent/child refs and transient lifecycle.

## 4. Non-objectifs
Ne définit pas memory architecture, storage format, raw-secret persistence, final schema, model context window, source-object mutation or permanent knowledge base.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator/Automation Designer; Agent supervisor, Security Reviewer and Auditor.

## 7. Conditions d’entrée
Automation Run and source data contract exist; authorized inputs/refs resolved; tenant/env known; masking/Secret Reference permissions available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/source context | CAP-STD-042 + source caller | input/object refs | oui | pinned/current | context incomplete |
| Workflow variables/mappings | STD-2 CAP-STD-019/023 | definition/mapping semantics | conditionnel | pinned | value unbound |
| step outputs | CAP-STD-044/047 | runtime value/output refs | conditionnel | source timestamp | missing/partial |
| Secret References | Platform Settings | opaque refs only | conditionnel | current auth | restricted/unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run / parent Run | CMDR Studio | context lineage | read |
| Workflow / Tool Call / Step | CMDR Studio | data contracts/output refs | read |
| source business object | source product | authorized reference/value projection | read only |
| Secret Reference | Platform Settings | opaque metadata/reference | restricted read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run Context | create/update transient values/refs | CMDR Studio | temporary context ≠ canonical source object |
| step-local/shared value projection | bind/propagate/mask/expire | CMDR Studio | no source mutation/no raw secret persistence |

## 11. Fonctionnalités
Run/step scopes; variable values vs definitions; object refs; sensitive masking; Secret References; outputs; parent/child context; cross-step propagation; transient lifecycle/expiry; provenance/freshness.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect permitted Run context | Studio Operator | Run Context | 0 | read | masked/source-backed values | non |
| Validate propagation/scope | Studio Reviewer | Run Context | 1 | mapping + permissions | pass/block | non |
| Modify allowed transient context | authorized operator | Run Context | 2 | within source permission/scope | versioned transient change | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
unbound, partially-bound, valid, restricted, masked, stale, source-missing, propagation-blocked, parent-linked, child-linked, expired-transient, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run Context snapshot/projection | temporary context | CAP-STD-044/049/050 | source refs/scope/masking/freshness |
| propagated step input/output refs | runtime bindings | next step/consumer | explicit mapping + provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| caller/Workflow | bind Run inputs | CAP-STD-048 | source refs/values/tenant | source |
| Run Step | emit output | CAP-STD-048 | step/output ref/provenance | Run |
| CAP-STD-048 | propagate allowed value | next Run Step/child Run | mapping/scope/source ref | Run |

## 18. Dépendances
STD-2 CAP-STD-019/023; STD-1 Tool output contracts; CAP-STD-035/036/042/044/047/049/050; Settings Secret References; Security tenant isolation; source-product permissions; OPEN-013/015.

## 19. Source de vérité
Studio owns temporary Run-context semantics. Source products own source objects/data; Settings owns Secret References/raw secret resolution; Workflow owns definitions/mappings; temporary state is not a canonical business object.

## 20. Provenance et audit
Record Run/step/source refs, binding/mapping, sensitivity/masking classification, value provenance or opaque reference (not secret), parent/child relation, freshness, changes, actor/time and correlation.

## 21. Permissions fonctionnelles
Run context read; sensitive Run context read separately; source-object permission; Secret Reference metadata; context update within bounded scope. Cross-tenant read requires explicit permission.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Cross-tenant ref, sensitive permission denial, Secret Reference unavailable, stale source, missing mapping, output partial, parent/child mismatch or expired transient value remains explicit and is not substituted.

## 23. Métriques
Masked values; propagation blocks; stale/missing sources; cross-tenant denials; transient-context retention violations; raw-secret persistence—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Run references data from another tenant without explicit permission, **When** context is resolved, **Then** the reference is blocked/masked and no cross-tenant value is propagated.

**Given** a Secret Reference becomes unavailable, **When** a dependent step needs it, **Then** the step is blocked without copying or substituting a secret.

**Given** a temporary Run context value is no longer needed, **When** its lifecycle ends, **Then** it is not automatically promoted into permanent Agent memory or a canonical source object.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Automation Runs, Workflows, Tool Calls, Control Room, Settings, Security, source products, CAP-STD-049..051, Quality.
