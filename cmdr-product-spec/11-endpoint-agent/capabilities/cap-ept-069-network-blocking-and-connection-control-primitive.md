---
id: CAP-EPT-069
title: Network Blocking and Connection Control Primitive
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-069 — Network Blocking and Connection Control Primitive

## 1. Définition
Définir une primitive distincte d’isolation host pour appliquer ou retirer un contrôle réseau borné à une destination/connexion/scope autorisé, avec expiry/conflict/partial enforcement et verification.

## 2. Problème utilisateur
Un temporary block peut n’affecter qu’une partie du trafic et entrer en conflit avec d’autres controls. Le confondre avec host isolation masque son scope réel.

## 3. Objectifs
Pinner target/network reference ; block/unblock semantics ; scope/expiry ; active state ; conflict ; partial/failed/unknown enforcement ; verification ; reversal provenance.

## 4. Non-objectifs
Créer firewall syntax/rules engine, isoler automatiquement tout le host, conclure threat removed, définir protocol/API, ou choisir platform/network stack.

## 5. Propriétaire
Endpoint owns bounded network control primitive and local enforcement observations. Govern owns authority/Response Run/Result; Settings owns Policy/admin configuration.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Network Reviewer, Govern Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready, exact network target/scope, block/unblock request, expiry if required, authority and current policy constraints.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| network/connection target | Govern/EPT observations | target | oui | current enough | unresolved |
| block/unblock + scope/expiry | Govern Step | intent | oui | pinned | blocked |
| current network state | Endpoint | precheck | oui | timestamped | unknown |
| Policy constraints | Settings | local restriction | selon primitive | current | block/unknown |
| authority ref | Govern | effect authority | oui | current | no action |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Step/Decision | Govern | scope/authority | read |
| connection/network observations | Endpoint | target/current state | read |
| Endpoint Policy | Settings | restrictions | read |
| isolation state | Endpoint | avoid semantic collision | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Control Execution | create/transition | Endpoint | bounded scope |
| Network Control State | observe/update | Endpoint | active/partial/conflict |
| Technical Outcome | create | Endpoint | not Result |
| network control effect | block/unblock | target | implementation unspecified |

## 11. Fonctionnalités
Validate target/scope/expiry; apply bounded block/unblock; record conflict; track partial enforcement; observe active/removed state; detect drift/expiry; preserve exact allowed/blocked reference without widening to host isolation.

## 12. Actions utilisateur
Inspect Class 0; check enforcement Class 1; block/unblock Class 3 by default and Govern-dependent.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate scope/conflict | oui | oui | oui | explain | deterministic checks |
| observe enforcement | oui | oui | oui | summary | network observations |
| suggest bounded target | oui | catalog/context | oui | suggestion | manual target selection |
| authorize/apply autonomously | non | interdit | non | interdit | Govern/operator |

## 14. États fonctionnels
`not-blocked`, `block-requested`, `applying`, `active-observed`, `partial-enforcement`, `conflict`, `failed`, `unknown`, `expiry-pending`, `expired-observed`, `unblock-requested`, `removing`, `removed-observed`, `drifted`.

## 15. États d’interface
No Screen ID. Bounded block, full isolation, network outage and source-unavailable are never collapsed.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network Control State | technical state | CAP-EPT-073/075 | exact scope/time |
| enforcement technical outcome | raw outcome | Govern | partial/conflict retained |
| expiry/reversal observation | technical fact | CAP-EPT-077/078 | no rollback inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | network control ready | CAP-EPT-069 | target/scope/expiry/authority | Run |
| CAP-EPT-069 | outcome | CAP-EPT-073/076 | enforcement status | same Run |
| CAP-EPT-069 | verify/reverse | CAP-EPT-074/075/077 | observed control state | Run retained |

## 18. Dépendances
CAP-EPT-019/039/058/062/065/066/068/073..080, Govern, Settings, Security, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of local bounded network-control execution/state. Govern retains response authority/outcome.

## 20. Provenance et audit
Target/scope/expiry, before/after state, conflicts, authority/Run refs, Policy version, operator, Agent/version, timestamps, errors, verification/reversal refs.

## 21. Permissions fonctionnelles
Network control read, block/unblock request, verification read, sensitive network metadata, cross-tenant deny, step-up/SoD/Govern dependency.

## 22. Limites et erreurs
Network block ≠ host isolation; block ≠ threat removed; partial enforcement ≠ success; expiry ≠ confirmed removal unless observed; unblock ≠ secure connectivity.

## 23. Métriques
Blocks by state, partial/conflict/unknown, expiry drift, verification mismatch, unauthorized attempts target zero.

## 24. Classification de livraison
`draft / defined / planned`; no firewall syntax/API/protocol/provider/implementation.

## 25. Critères d’acceptation
**Given** a bounded network block applies only to some targets, **When** observed, **Then** partial enforcement remains explicit and host isolation is not claimed.

**Given** a block conflicts with another control, **When** execution runs, **Then** conflict and actual observed state are preserved without silent override.

**Given** unblock technically succeeds, **When** verification occurs, **Then** restored connectivity is not interpreted as endpoint security.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; no platform firewall implementation or Class-2 default is selected.

## 27. Consommateurs documentaires
EPT-5 verification/reversal, Govern, Settings, Security, Network investigation consumers, Quality.