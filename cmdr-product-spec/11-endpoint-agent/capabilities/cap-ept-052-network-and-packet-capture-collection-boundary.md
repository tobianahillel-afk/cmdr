---
id: CAP-EPT-052
title: Network and Packet Capture Collection Boundary
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-020, REQ-PROD-055, REQ-SEC-001, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-052 — Network and Packet Capture Collection Boundary

## 1. Définition
Définir une capture réseau/packet technique bornée lorsque supportée : interface/context ref, duration/size/filter conceptuels, start/stop, progress, partial/loss/error et neutral output/provenance, sans moteur ni syntaxe de filtre.

## 2. Problème utilisateur
Une capture non bornée peut collecter trop de données sensibles ou être faussement considérée complète malgré pertes/arrêt tardif.

## 3. Objectifs
Valider capability/interface/context ; borner durée/volume/filter conceptuel ; distinguer start request/start confirmed/stop request/stopped ; exposer loss/partial ; produire neutral Collection Item.

## 4. Non-objectifs
Aucun packet capture command, interface API, BPF-like syntaxe, format, detection verdict, network block/isolation or universal platform support.

## 5. Propriétaire
Endpoint owns target-side capture execution and loss/status facts. Investigate owns request/Case and later Artifact/network analysis qualification.

## 6. Utilisateurs
DFIR Analyst, Response Operator, Endpoint Operator, Evidence Reviewer, Privacy Reviewer.

## 7. Conditions d’entrée
CAP-EPT-048 network item, capability declared, interface/context available, bounded duration/volume/filter concept, policy/permission/authority current.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| capture plan item | CAP-EPT-048 | scope/bounds | oui | plan version | no capture |
| interface/context capability | Endpoint | target support | oui | current | unsupported |
| duration/volume/filter concept | request/plan | bounds | oui | start snapshot | blocked |
| policy/authority | Settings/Govern | gate | oui | current | denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request/Case | Investigate | reason/scope | read |
| Technical Plan | Endpoint | capture bounds | read |
| endpoint-agent | Endpoint | network capability/context | read |
| Endpoint Policy | Settings | sensitivity/limits | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Capture Attempt | create/start/stop/update | Endpoint | bounded only |
| Network Collection Item | create/partial | Endpoint | neutral output |
| Capture Loss Marker | derive | Endpoint | observed/declared loss explicit |

## 11. Fonctionnalités
Validate scope, request/confirm start, track duration/volume, request/confirm stop, record limit reached/loss/partial/error, produce output reference with actual period and limitations.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect capture scope | Analyst | plan | 0 | read | bounds visible | non |
| start bounded capture | authorized path | capture attempt | 2/3 by impact | support + authority | start-requested/capturing | according impact |
| stop/cancel capture | Operator | attempt | 2 | running/cancellable | stop/cancel requested | OPEN-013 |
| create block rule | aucun EPT-4 | network target | 3+ | EPT-5/Govern | not executed | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate bounds/support | oui | oui | oui | explain | rules |
| track duration/volume | oui | oui | oui | summarize | counters |
| explain loss/errors | oui | oui | oui | oui | raw status |
| extend duration automatically | non | interdit | non | interdit | explicit reauthorization |

## 14. États fonctionnels
`preparing`, `unsupported`, `awaiting-authority`, `start-requested`, `capturing`, `stop-requested`, `stopping`, `partial`, `completed`, `failed`, `timed-out`, `cancelled`.

## 15. États d’interface
No Screen ID. Start-requested != capturing; stop-requested != stopped; losses/partial stay visible.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network Collection Item | Endpoint output | CAP-EPT-054/055 | capture != compromise verdict |
| actual capture period/loss | metadata | Investigate | limitations explicit |
| status/errors | Endpoint state | CAP-EPT-053 | no hidden loss |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-048 | capture item | CAP-EPT-052 | target/bounds/authority | plan |
| CAP-EPT-052 | output | CAP-EPT-054/055 | item/period/loss | attempt |
| Investigate | analyze | Network Forensics | qualified ref | collection source retained |

## 18. Dépendances
Investigate CAP-INV-208, CAP-EPT-048/053..055, Settings Policy, OPEN-008/013/014.

## 19. Source de vérité
Endpoint SOT of local capture execution/loss facts; Investigate SOT of request and analytical qualification.

## 20. Provenance et audit
Request/plan, Agent/platform, interface/context ref, duration/volume/filter concept, policy/authority, start/stop timestamps, loss/errors/output and correlation.

## 21. Permissions fonctionnelles
Network capture request/start/stop/cancel, sensitive network output read, provenance, cross-tenant deny.

## 22. Limites et erreurs
Packet capture ≠ compromise/detection verdict ; start request ≠ start ; stop request ≠ stopped ; loss may make output partial ; no network containment.

## 23. Métriques
Capture start/stop failures, actual duration/volume, loss/partial/failure, timeouts/cancellations, provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no capture engine, filter syntax, format or protocol.

## 25. Critères d’acceptation
**Given** network capture unsupported on platform, **When** requested, **Then** capture remains unsupported and never starts.

**Given** loss is reported, **When** capture ends, **Then** output is partial/limited and no complete claim is made.

**Given** stop is requested but not confirmed, **When** status is shown, **Then** target-side termination remains unknown/stopping rather than completed.

## 26. Questions ouvertes
OPEN-008/013/014 remain open; no low-level capture technology selected.

## 27. Consommateurs documentaires
EPT-4 packaging/transfer, Investigate Network Forensics, Govern boundary, Security/Privacy, Quality.
