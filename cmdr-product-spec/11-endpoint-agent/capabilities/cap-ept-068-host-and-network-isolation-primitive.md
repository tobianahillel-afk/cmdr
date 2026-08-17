---
id: CAP-EPT-068
title: Host and Network Isolation Primitive
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-068 — Host and Network Isolation Primitive

## 1. Définition
Définir la primitive de containment d’un host visant un état d’isolation réseau borné, avec exceptions de control/management path seulement comme concept sourcé, état requested/observed, impact/connectivity et release eligibility.

## 2. Problème utilisateur
Isolation peut être partielle ou couper des chemins inattendus. Un endpoint inaccessible ne prouve ni que l’isolation est active ni qu’elle est saine.

## 3. Objectifs
Pinner host target ; current/requested isolation state ; exceptions conceptuelles ; effectuer un effect governed ; observer enforcement/connectivity ; distinguer active/partial/failed/unknown ; préparer release/verification.

## 4. Non-objectifs
Définir firewall rules/syntax, réseau exact, control-plane protocol, déclarer endpoint offline, conclure incident contained, autoriser isolation, ou fournir UX finale.

## 5. Propriétaire
Endpoint owns technical isolation primitive/state observations. Govern owns authority and response-level containment success. Settings owns Policy/Fleet/admin configuration.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Network/Security Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready, exact host/Agent binding, isolation capability declared, authority current, policy constraints and expected verification path available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| host/Agent target | EPT-1 | target | oui | current | unresolved |
| requested isolation state | Govern Run/Step | intent | oui | pinned | blocked |
| isolation capability/state | Endpoint | technical facts | oui | current observation | unknown/unavailable |
| policy/exception projection | Settings | restriction | selon policy | current | block/unknown |
| authority ref | Govern | authority | oui | effective | no execution |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Decision | Govern | scope/conditions/authority | read/ref |
| Endpoint Agent | Endpoint | identity/connectivity/capability | read |
| Endpoint Policy | Settings | isolation constraints/exceptions | read |
| network observation | Endpoint/EPT-2 | reachability context | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Isolation Technical Execution | create/transition | Endpoint | exact target/authority |
| Isolation State | observe/update | Endpoint | requested vs observed distinct |
| Technical Outcome | create | Endpoint | not Govern Result |
| network reachability | bounded effect | target | implementation unspecified |

## 11. Fonctionnalités
Request isolation/release only under authority; track requested/start/active/partial/fail/unknown; observe allowed management-path concept when sourced; preserve policy exceptions; verify connectivity effects; expose release eligibility without auto-release.

## 12. Actions utilisateur
Inspect = Class 0; readiness/verification = Class 1; isolate/release = Class 3 by default, governed. OPEN-013 prevents arbitrary Class 2 promotion.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare requested/observed isolation | oui | oui | oui | explain | state matrix |
| verify connectivity indicators | oui | oui | oui | summarize | source observations |
| propose isolation | oui | rules/catalog | oui | suggestion | human selection |
| authorize/execute independently | non | no | non autonome | interdit | Govern path |

## 14. États fonctionnels
`not-isolated`, `isolation-requested`, `isolating`, `isolated-observed`, `partial-isolation`, `isolation-failed`, `isolation-unknown`, `management-path-unknown`, `release-eligible`, `release-requested`, `releasing`, `released-observed`, `drifted`.

## 15. États d’interface
No Screen ID. Offline, network failure and isolation states remain distinct; inability to contact target never fabricates isolation success.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Isolation State | technical fact | CAP-EPT-073/075/078 | source/time/state |
| isolation technical outcome | execution fact | Govern | partial/unknown preserved |
| connectivity observations | verification inputs | CAP-EPT-074/075 | not business verdict |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | isolation ready | CAP-EPT-068 | host/state/policy/authority | Govern Run |
| CAP-EPT-068 | technical status | CAP-EPT-073/076 | execution + partial/unknown | same Run |
| CAP-EPT-068 | verify/release | CAP-EPT-074/075/078 | observed state/limits | Govern retained |

## 18. Dépendances
CAP-EPT-009/019/028/039/062/065/066/073..080, Govern Runs/Verification, Settings Policy/Fleet, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of local isolation technical state; Govern SOT of response containment outcome; Settings SOT of isolation policy/admin configuration.

## 20. Provenance et audit
Host/Agent, requested state, policy/exceptions, authority/Decision/Run/Step, precheck, connectivity before/after, execution/verification states, operator, Agent/version, timestamps/errors/correlation.

## 21. Permissions fonctionnelles
Isolation read/request, release request, verification read, sensitive network context, cross-tenant deny, step-up/SoD/Govern dependency; no final RBAC.

## 22. Limites et erreurs
Host isolated ≠ endpoint offline; isolation ≠ network failure; isolated technically ≠ incident contained; target unreachable ≠ isolated; release ≠ secure endpoint.

## 23. Métriques
Isolation active/partial/failure/unknown, management-path unknown, verification mismatch, release attempts, unauthorized/cross-tenant target zero.

## 24. Classification de livraison
`draft / defined / planned`; no firewall implementation, protocol, provider, native command or supported platform claim.

## 25. Critères d’acceptation
**Given** isolation applies only partially, **When** state is observed, **Then** `partial-isolation` is retained and no contained verdict is created.

**Given** host is isolated but management state is unknown, **When** verification runs, **Then** management-path uncertainty is explicit and not converted to failure or success automatically.

**Given** the target becomes unreachable, **When** state is reconciled, **Then** unreachable and isolated remain separate hypotheses/facts.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; control-plane implementation and exact platform behavior are not chosen.

## 27. Consommateurs documentaires
EPT-5 verification/release/provenance, Govern, Settings, Security, Command/Investigate, Quality.