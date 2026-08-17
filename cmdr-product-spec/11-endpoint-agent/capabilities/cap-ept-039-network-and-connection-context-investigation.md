---
id: CAP-EPT-039
title: Network and Connection Context Investigation
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-039 — Network and Connection Context Investigation

## 1. Définition
Définir l’investigation locale des connexions et listeners déjà observés : local/remote refs, direction, process/user relation, timing, DNS/context, related detections et chronology, sans packet capture ni action réseau.

## 2. Problème utilisateur
Une connexion associée à un signal nécessite du contexte, mais une adresse distante, un DNS hit ou une reputation ne prouve pas une compromission et ne justifie pas automatiquement une réponse.

## 3. Objectifs
Contextualiser connection/listener/flow facts ; relier process/user/DNS observations ; conserver timing et historique limité ; montrer gaps/rate/sampling ; permettre pivots locaux permission-aware.

## 4. Non-objectifs
Aucun packet capture, active scan, firewall/network mutation, isolation, block, network forensics acquisition, compromise verdict, SIEM/network detection engine.

## 5. Propriétaire
Endpoint Agent possède Network Context local dérivé. Investigate conserve Network Forensics/Evidence/Finding ; Govern conserve les réponses ; Shared conserve Linking/Timeline/Search.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Security Reviewer ; Auditor ; Command consumer.

## 7. Conditions d’entrée
Network observations EPT-2 ; process attribution/DNS refs si disponibles ; tenant/permission ; source limitations, sampling et historical coverage explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| network/connection observations | CAP-EPT-019 | faits locaux | oui | source freshness | context unavailable |
| process/user refs | CAP-EPT-017/020/037 | relations | non | source/snapshot time | attribution unknown |
| DNS/context refs | CAP-EPT-019 | observations | non | source-owned | no DNS inference |
| local detections | CAP-EPT-032..035 | signal context | non | candidate freshness | network context only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Network/Connection Observation | Endpoint Agent | endpoints/direction/time | read |
| Process/User Context | Endpoint Agent | relation refs | read |
| Local Detection Candidate | Endpoint Agent | related refs | read |
| Evidence/Finding | Investigate | destination refs only | no local mutation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network Investigation Context | dériver/rafraîchir | Endpoint Agent | observations existantes seulement |
| Connection-Process Relation | dériver | Endpoint Agent | attribution sourcée |
| Network Chronology Projection | dériver | Endpoint Agent | historical limitations visibles |

## 11. Fonctionnalités
Afficher connection/listener details, direction, local/remote refs, process/user attribution et DNS facts existants ; lier related candidates ; construire chronology locale tout en conservant sampling/gap/history limitations.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect connection context | Analyst | Network Context | 0 | read | facts/relations visibles | non |
| correlate existing network refs | deterministic service | context | 1 | observations disponibles | liens sourcés | non |
| request bounded context refresh | Operator | local connection state | 2 | sans packet capture/acquisition | freshness update | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| lier process/connection | oui | oui | oui | suggestion possible | source keys/time |
| construire chronology | oui | oui | oui | non nécessaire | timestamps |
| résumer réseau | oui | oui | oui | oui | structured context |
| conclure compromise | non | non | non | interdit | Investigate qualification |

## 14. États fonctionnels
`available`, `partial`, `attribution-unknown`, `sampled`, `rate-limited`, `stale`, `unsupported`, `unavailable`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Attribution inconnue, historical limitations et sampling restent visibles ; remote endpoint/reputation ne doit pas être affiché comme compromise.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network Investigation Context | concept Endpoint | CAP-EPT-042..045 | no packet capture |
| process/DNS relations | projections Endpoint | Analyst | source refs explicites |
| network context gap | diagnostic | Analyst/EPT-4 future | no response/acquisition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| network observation/candidate | pivot | Network Context | connection/process/DNS refs | origin retained |
| Network Context | process/user/timeline pivot | CAP-EPT-037/040/042 | stable refs/time | no collection |
| packet/content needed | unavailable local data | EPT-3 stop | requirement + limitation | futur Collection seulement |

## 18. Dépendances
CAP-EPT-017/019/020/024..026/032..037 ; `investigation/network-connections.md`; Shared Timeline/Linking ; Investigate Network Forensics ; `OPEN-008`.

## 19. Source de vérité
Endpoint est SOT du contexte local dérivé des network observations. Investigate reste owner de toute Evidence/Finding/forensic product ; Govern de toute network response.

## 20. Provenance et audit
Connection refs, local/remote/direction, attribution source, DNS refs, timing, sampling/rate state, related candidates, masking, tenant, actor/service et snapshot/freshness.

## 21. Permissions fonctionnelles
Network context read, sensitive identifiers/DNS/user attribution read selon classification, provenance, cross-tenant deny. Aucun droit de capture/block/isolation.

## 22. Limites et erreurs
Network context ≠ compromise ; reputation hit ≠ compromise ; missing DNS/process attribution ≠ benign ; chronology peut être partielle ; packet capture est hors EPT-3.

## 23. Métriques
Attribution coverage, DNS-link coverage, partial/stale/sampled states, context pivot success, restricted identifiers.

## 24. Classification de livraison
`draft / defined / planned`; aucune capture, scanner, firewall ou response implementation.

## 25. Critères d’acceptation
**Given** une connexion est liée à un local signal candidate, **When** l’analyste pivote, **Then** process/direction/timing/DNS refs disponibles sont montrés sans verdict compromise.

**Given** attribution process manque, **When** context est construit, **Then** `attribution-unknown` est explicite et aucune attribution n’est inventée.

**Given** un packet capture serait nécessaire, **When** le pivot atteint cette limite, **Then** EPT-3 s’arrête et ne lance pas Collection.

## 26. Questions ouvertes
`OPEN-008` conserve support des sources/protocol facts ; aucune technique de capture ou réponse n’est décidée.

## 27. Consommateurs documentaires
EPT-3 timeline/pivots/summary ; Investigate Network Forensics/Case ; Command ; Security ; Quality ; future EPT-4 boundary.