---
id: CAP-INV-355
title: Network State and Connection Artifact Reconstruction
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-355 — Network State and Connection Artifact Reconstruction

## 1. Définition
Permet de reconstruire les traces réseau candidates présentes en mémoire avec processus, endpoints, états et limites sans devenir une Network Forensics complète, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans reconstruction bornée, des sockets partielles ou sans timestamp peuvent être présentées à tort comme connexions complètes ou IOC confirmés.

## 3. Objectifs
- reconstruire connexions/sockets candidates présentes dans l’image.
- voir endpoints, ports, états, processus, timestamps et limites disponibles.
- préparer corrélation/indicateur candidat et futur handoff Network Forensics sans qualification automatique.

## 4. Non-objectifs
Aucune PCAP, reconstruction de session complète, analyse protocolaire, moteur, commande, API, action Endpoint ou capability Network Forensics complète.

## 5. Propriétaire
Investigate possède les observations mémoire; future Network Forensics possède son analyse complète; Studio Tools/Runs; Settings policies; Shared linking/trace.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Network Analyst futur, Evidence Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Image/session/profil/Tool et permissions lisibles; timestamps et structures manquants visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Platform/Profile | CAP-INV-350 | interprétation | oui | courant | unsupported |
| Network state results | Tool Calls | structures candidates | oui | Tool/version visibles | partial/failed |
| Process observations | CAP-INV-351 | attribution process | non | même image | unlinked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture |
| Process Observation | Investigate concept | process links | lecture/lien |
| Tool/Tool Call/Run | Studio | version/paramètres/statut | lecture |
| Entity | Shared/Investigate | relation candidate | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Network State Observation | créer/annoter/contester | Investigate concept | source/incertitude obligatoires |
| Connection candidate | créer/lier/versionner | Investigate concept | ≠ connexion complète/IOC |
| Future handoff package | préparer | Investigate | destination 4B.2B.3B seulement |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- reconstruire connexions candidates et structures équivalentes.
- voir endpoints locaux/distants, ports/protocoles fonctionnels, processus, états et timestamps disponibles.
- voir connexions terminées lorsque détectables et incohérences.
- filtrer, comparer, annoter et relier à Entity.
- préparer indicateur candidat et handoff futur.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Network Observation | 0 | lecture autorisée | vue sourcée | non |
| Comparer | Analyst | Observation set | 0/1 | sources compatibles | différences visibles | non |
| Annoter/relier | Analyst | Candidate | 2 | permission | version conservée | OPEN-013 |
| Préparer handoff | Analyst | Future package | 2 | sources/limites | package attribué | non |

Classes 3/4 bloquées et aucune action réseau réelle.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| reconstruire network state | oui | oui | oui | proposition | reconstructeur déterministe |
| relier processus/endpoints | oui | oui | oui | suggestion | table de relations |
| grouper/comparer | oui | oui | oui | résumé | filtres/comparateur |
| confirmer IOC | oui | non | non | assistance | revue humaine |

Attribution complète de toute automatisation.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale conservent contexte et exposent les manques.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Network State Observation | memory result | Workbench | source/partialité visibles |
| Connection/indicator candidate | candidate | Case/Hypothesis | aucun IOC automatique |
| Future Network package | handoff | Phase 4B.2B.3B | limites et provenance conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Process/Session | ouvrir network state | CAP-INV-355 | image, process, profil, Tool | source |
| Observation | corréler | Entity/Timeline | endpoints, process, timestamp quality | Network view |
| Observation | futur handoff | Phase 4B.2B.3B | candidates, links, timestamps, limites | Network view |

Tenant, Case, image, permissions et return origin préservés.

## 18. Dépendances
CAP-INV-350/351/359/362, Shared Entity/Timeline, future 4B.2B.3B, OPEN-005/008/013/015.

## 19. Source de vérité
Observation memory-resident : Investigate; full Network Forensics : future phase; Tool/Run : Studio; linking/trace : Shared.

## 20. Provenance et audit
Image/session/profil, Tool/version, paramètres, structure candidate, process links, endpoints, timestamp quality, acteur, erreurs, partialité et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| network state read | données sensibles | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| compare/link | corrélation | 0/2 | policy | auteur/reviewer | Investigate | Permissions |
| future handoff | impact externe | 2 | possible | owner futur valide | Investigate/future | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- reconstruction partielle et timestamps absents visibles.
- socket candidate ≠ échange complet.
- destination ≠ IOC confirmé.
- aucune PCAP, analyse protocolaire ou full Network Forensics.

## 23. Métriques
- observations avec/sans process/timestamp.
- taux partial/disputed.
- candidates confirmées/rejetées par revue humaine.
- handoffs futurs avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation, moteur ou protocole revendiqué.

## 25. Critères d’acceptation
### 1. Timestamp absent
**Given** connexion candidate sans timestamp complet **When** indicateur préparé **Then** absence visible, reconstruction partielle et aucun IOC confirmé.
### 2. Processus non lié
**Given** aucune attribution process fiable **When** affichée **Then** `unlinked/partial` est visible.
### 3. Sans IA
**Given** aucun modèle **When** analyse réalisée **Then** reconstructeur, tables, filtres et comparaison fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; schémas/plateformes/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351/359/362, Entity/Case, future Phase 4B.2B.3B uniquement par handoff.
