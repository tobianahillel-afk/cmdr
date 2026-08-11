---
id: CAP-EPT-095
title: Local Privilege, Service/Process Protection and Security Context
product: endpoint-agent
module: security
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-095 — Local Privilege, Service/Process Protection and Security Context

## 1. Définition
Définir le contexte technique local de privilège et de protection des processus/services Endpoint, sans transformer privilège disponible en permission ou autorité.

## 2. Problème utilisateur
Une capability peut échouer ou se dégrader si le runtime ne possède plus le contexte de privilège requis, sans que cela signifie un refus RBAC ou une compromission.

## 3. Objectifs
Represent required/observed privilege, insufficient privilege, protected process/service state, changed execution context, capability impact, consumer projection and provenance.

## 4. Non-objectifs
Aucun modèle OS privilege final, service definition, sandbox implementation, RBAC/ABAC, authority decision, exploit/tamper verdict or support claim.

## 5. Propriétaire
Endpoint owns observed local runtime privilege/protected-component facts. Security owns permission model; Settings owns configuration; Govern owns response authority.

## 6. Utilisateurs
Endpoint Operator, Security Reviewer, Platform Administrator, SOC consumer, Auditor.

## 7. Conditions d’entrée
Agent/component identity, required privilege source, observed execution context, capability relation, tenant/environment and permission to inspect sensitive state.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| required privilege context | Endpoint capability/source | prerequisite | oui when applicable | versioned | requirement unknown |
| observed runtime privilege | Endpoint local state | technical fact | oui | fresh | privilege unknown |
| process/service protection state | Endpoint local source | technical fact | non | fresh | state unknown |
| permission/policy context | Security/Settings | external control | non | current | no authorization inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent/component | Endpoint | runtime context | read |
| Technical Capability Availability | Endpoint | capability impact | read |
| Endpoint Policy | Settings | configured constraints | read projection |
| Permission Model | Security | distinction only | reference |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Runtime Privilege State | derive/refresh | Endpoint | privilege != permission |
| Protected Component State | create/refresh | Endpoint | local fact only |
| Capability Security-Context Impact | derive | Endpoint | effect/limits explicit |

## 11. Fonctionnalités
Compare required vs observed privilege; expose insufficient/changed context; observe process/service protection; link affected capabilities; preserve distinction from RBAC/authority and from compromise/tamper verdicts.

## 12. Actions utilisateur
Inspect Class 0; deterministic prerequisite assessment Class 1; bounded diagnostic refresh Class 2. Privilege mutation/elevation is not introduced by this capability.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare required/observed privilege | oui | oui | oui | explain | deterministic comparison |
| map capability impact | oui | oui | oui | summarize | capability dependency map |
| grant permission/elevate autonomously | non | interdit | non | interdit | Security/platform workflow |

## 14. États fonctionnels
`privilege-sufficient`, `privilege-insufficient`, `privilege-unknown`, `execution-context-changed`, `protected-component-observed`, `protection-degraded`, `capability-affected`, `stale`, `restricted`.

## 15. États d’interface
No Screen ID. Privilege and permission are separately labeled; sensitive state can be masked without fabricating normality.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Runtime Privilege State | Endpoint fact | CAP-EPT-093/098/Settings | local technical scope |
| Protected Component State | Endpoint fact | CAP-EPT-094/098 | no compromise verdict |
| Capability Impact | derived fact | operations/Quality | source/reason explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local runtime | context change | CAP-EPT-095 | required/observed privilege | local context |
| CAP-EPT-095 | insufficient privilege | CAP-EPT-093 | affected capability/dependency | reassessment |
| CAP-EPT-095 | security-state handoff | CAP-EPT-098 | state/reason/provenance | no permission grant |

## 18. Dépendances
Endpoint least-privilege/security; CAP-EPT-011/093/094/097/098; Settings; Security permission model; OPEN-008/013.

## 19. Source de vérité
Endpoint SOT for observed local privilege/component facts; Security SOT for permissions; Settings SOT for admin config.

## 20. Provenance et audit
Agent/component, required/observed context, source/version, process/service state, affected capabilities, timestamps, masking and prior/new state.

## 21. Permissions fonctionnelles
Privilege/security-context read, sensitive component-state read, diagnostic refresh, provenance read, cross-tenant deny; no final RBAC/privilege-elevation permission is defined.

## 22. Limites et erreurs
Privilege available != permission granted; service/process running != Agent uncompromised; insufficient privilege != RBAC denial automatically.

## 23. Métriques
Insufficient/unknown privilege, changed contexts, protected-component degradation, affected-capability and stale/restricted counts.

## 24. Classification de livraison
`draft / defined / planned`; no OS privilege/service/sandbox implementation.

## 25. Critères d’acceptation
**Given** a capability requires a privilege absent from runtime, **When** prerequisites are assessed, **Then** privilege-insufficient is reported without claiming permission denial.

**Given** a protected process is running, **When** security state is reviewed, **Then** running state is not treated as proof the Agent is uncompromised.

**Given** AI is unavailable, **When** privilege requirements are compared, **Then** deterministic/manual assessment remains available.

## 26. Questions ouvertes
OPEN-008/013 remain open; platform privilege implementation and effect classification are not selected.

## 27. Consommateurs documentaires
EPT-6 Security/Resilience, Settings, Security Architecture, Quality, registers, Roadmap and implementation contracts.