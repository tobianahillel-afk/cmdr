---
id: CAP-EPT-094
title: Endpoint Self-Protection and Anti-Tamper Boundary
product: endpoint-agent
module: security
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-094 — Endpoint Self-Protection and Anti-Tamper Boundary

## 1. Définition
Définir fonctionnellement l’état attendu/observé de self-protection et les tamper candidates sur composants Endpoint sans spécifier mécanisme kernel, crypto ou verdict de compromission.

## 2. Problème utilisateur
Un service protégé arrêté ou une configuration modifiée peut indiquer une anomalie sans prouver sabotage, compromission ou absence d’intégrité globale.

## 3. Objectifs
Represent protected component, expected/observed protection state, tamper/unauthorized-change candidate, disabled/degraded protection, sourced reaction boundary, config/authority owner and provenance.

## 4. Non-objectifs
Aucun driver/kernel mechanism, secure boot, cryptographic integrity proof, malware verdict, automatic containment, final anti-tamper implementation or support claim.

## 5. Propriétaire
Endpoint owns local protection/tamper observations. Security owns global trust policy; Settings owns relevant configuration; Govern owns response authority.

## 6. Utilisateurs
Endpoint Operator, Security Reviewer, Platform Administrator, SOC/Investigate consumer, Govern Reviewer, Auditor.

## 7. Conditions d’entrée
Known Agent/component, expected protection source/config, local observed state, tenant/environment, permission and freshness.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| protected component ref | Endpoint architecture/config | technical identity | oui | versioned | unknown component |
| expected protection state | Settings/Security policy projection | expectation | oui if defined | current | expectation unknown |
| observed protection state | Endpoint local monitor | technical fact | oui | fresh | observed unknown |
| maintenance/authority context | Settings/Govern/Security | control context | non | current | no authorized-change inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent/component | Endpoint | local runtime identity/state | read |
| Endpoint Policy/config | Settings | expected protection | read projection |
| permission/authority context | Security/Govern | allowed change context | read ref |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Self-Protection State | derive/refresh | Endpoint | technical fact only |
| Tamper Observation Candidate | create/version | Endpoint | candidate != compromise verdict |
| Protection Degradation State | derive | Endpoint | reason/source explicit |

## 11. Fonctionnalités
Compare expected/observed protection; identify candidate unauthorized changes; distinguish maintenance-authorized context where sourced; expose disabled/degraded/unknown protection; record local reaction only as sourced boundary; preserve evidence/provenance.

## 12. Actions utilisateur
Inspect Class 0; normalize/assess Class 1; bounded protection-state refresh Class 2; any protection mutation potentially Class 3 and remains policy/authority dependent under OPEN-013.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare expected/observed | oui | oui | oui | explain | deterministic comparison |
| flag tamper candidate | oui | source rules | oui | summarize | explicit indicators |
| declare compromise/authorize action | non | interdit locally | non | interdit | Investigate/Govern workflow |

## 14. États fonctionnels
`protected-observed`, `protection-degraded`, `protection-disabled`, `expected-state-unknown`, `observed-state-unknown`, `tamper-candidate`, `authorized-maintenance-context`, `unauthorized-change-candidate`, `recovery-pending`, `stale`.

## 15. États d’interface
No Screen ID. Tamper candidate and compromise verdict remain distinct; protected/running never implies uncompromised.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Self-Protection State | Endpoint fact | CAP-EPT-098/Settings/Security | local scope explicit |
| Tamper Observation Candidate | Endpoint fact | Investigate/Command consumers | no maliciousness verdict |
| protection provenance | audit context | CAP-EPT-097 | expected/observed/source retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local monitor/config | state change | CAP-EPT-094 | component/expected/observed/context | local state |
| CAP-EPT-094 | suspicious candidate | CAP-EPT-098/Investigate | candidate/reason/provenance | no Finding creation |
| CAP-EPT-094 | recovery needed | CAP-EPT-093/manual governed path | protection/dependency facts | recovery context |

## 18. Dépendances
Endpoint security anti-tamper/self-protection-recovery; CAP-EPT-008/093/095/097/098; Settings; Security; Govern; OPEN-008/013.

## 19. Source de vérité
Endpoint SOT for local protection/tamper observations; Security/Settings/Govern retain policy/config/authority.

## 20. Provenance et audit
Agent/component/version, expected/observed state, config/policy refs, maintenance/authority refs, candidate reason, time, prior/new state and recovery refs.

## 21. Permissions fonctionnelles
Self-protection state read, sensitive tamper read, bounded refresh, protection mutation only if separately sourced/authorized, provenance read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Self-protection != full endpoint security; anti-tamper state != crypto proof; integrity PASS != uncompromised; tamper candidate != malicious verdict.

## 23. Métriques
Protection degraded/disabled, tamper/unauthorized-change candidates, stale/unknown observations, recovery-pending and false auto-verdict target zero.

## 24. Classification de livraison
`draft / defined / planned`; no anti-tamper/kernel/crypto implementation.

## 25. Critères d’acceptation
**Given** a protected service is unexpectedly stopped, **When** state is assessed, **Then** protection is degraded/tamper-candidate without compromise verdict.

**Given** a sourced maintenance window explains a change, **When** reviewed, **Then** maintenance context is preserved rather than silently labeling maliciousness.

**Given** AI is unavailable, **When** expected/observed states are compared, **Then** deterministic/manual assessment remains available.

## 26. Questions ouvertes
OPEN-008/013 remain open; final platform anti-tamper implementation and governance class are not selected.

## 27. Consommateurs documentaires
EPT-6 Security/Resilience, Settings, Security Architecture, Investigate/Command/Govern consumers, Quality, registers and Roadmap.