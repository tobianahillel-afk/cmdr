---
id: CAP-INV-359
title: Memory Timeline and Cross-Source Correlation
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
  - REQ-UX-007
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-359 — Memory Timeline and Cross-Source Correlation

## 1. Définition
Organiser les événements temporels disponibles et corréler mémoire, Case, Collection, Dynamic Analysis et Endpoint telemetry en exposant qualité et incertitude, sans remplacer la Case Timeline.

## 2. Problème utilisateur
Sans qualité temporelle explicite, des timestamps observés, reconstruits ou estimés peuvent être fusionnés et produire une chronologie artificiellement certaine.

## 3. Objectifs
- distinguer timestamps observés, reconstruits, estimés et absents.
- corréler processus, threads, modules, connexions, Artifacts et acquisition.
- comparer images, exposer conflits/trous et préparer Evidence candidate.

## 4. Non-objectifs
Aucune fabrication de timestamp, remplacement de Case Timeline, timeline Disk/Network complète, moteur, format ou API.

## 5. Propriétaire
Investigate possède Memory Timeline projection; Shared possède Timeline Engine; Case Timeline reste propriétaire de la chronologie du Case.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Investigation Lead, Timeline Reviewer et Evidence Reviewer.

## 7. Conditions d’entrée
Observations sourcées et qualité temporelle disponibles; timezone et absence de timestamp visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory observations | CAP-INV-349..358 | événements/relations | oui | même session/image | empty/partial |
| Timestamp quality/timezone | sources/analyst | observed/reconstructed/estimated/absent | oui | par événement | ambiguous |
| Case/Collection/Dynamic/Endpoint events | owners respectifs | corrélation externe | non | freshness visible | local-only |
| Comparison scope | analyst | images/période/filtres | non | versionné | default scope |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory observations | Investigate | source/time quality | lecture |
| Case Timeline | Investigate | événements liés | lecture/lien |
| Collection/Dynamic/Endpoint events | owners respectifs | projections | lecture |
| Timeline mechanism | Shared | ordering/filtering | consommation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Timeline projection | créer/annoter/contester | Investigate concept | ≠ Case Timeline |
| Correlation relation | créer/lier/versionner | Investigate concept | source/quality obligatoires |
| Evidence candidate selection | préparer | Investigate | qualification future |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- organiser événements et afficher source, timezone et qualité temporelle.
- voir processus, threads, modules, connexions, Artifacts et acquisitions.
- corréler Case/Collection/Dynamic/Endpoint telemetry.
- comparer images, filtrer, annoter, afficher conflits et trous.
- préparer Evidence candidate sans inventer les timestamps absents.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Consulter/filtrer | Analyst | Memory Timeline | 0 | lecture autorisée | vue sourcée | non |
| Corréler/comparer | Analyst | Correlation | 0/1 | sources lisibles | relations/limites | non |
| Annoter/contester | Analyst | Timeline entry | 2 | permission | version conservée | OPEN-013 |
| Préparer Evidence | Analyst | Candidate | 2 | provenance | package | non |

Classes 3/4 indisponibles.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| ordonner événements | oui | oui | oui | résumé | timeline déterministe |
| qualifier timestamps | oui | oui | oui | proposition | labels/règles explicites |
| proposer corrélations | oui | oui | oui | suggestion | filtres/comparateur |
| conclure séquence | oui | non | non | assistance | revue humaine |

Attribution complète et sources visibles.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. États objet finaux reportés.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale; trous et timestamps absents restent visibles.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Memory Timeline | projection | Workbench/Case | quality/source visibles |
| Correlation relation | relation | Case/Hypothesis | pas de timestamp inventé |
| Evidence selection | candidate | CAP-INV-362 | provenance/contradictions conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Observations | construire timeline | CAP-INV-359 | timestamps, quality, sources | source |
| Timeline | corréler Case | Case Timeline | linked events, conflits, limites | Memory Timeline |
| Timeline | handoff | CAP-INV-362 | selected events, uncertainty | Memory Timeline |

Tenant, Case, image, filtres et return origin préservés.

## 18. Dépendances
CAP-INV-349..358/361/362, Shared Timeline/Linking, Case Timeline, Collection/Dynamic/Endpoint projections, OPEN-005/008/013/015.

## 19. Source de vérité
Observations/Memory Timeline : Investigate; Case Timeline : owner Case; ordering mechanism : Shared; sources externes restent propriétaires.

## 20. Provenance et audit
Chaque événement conserve image/session/source, timestamp value/quality/timezone, Tool/version, relation externe, auteur, filtre, conflit, erreur et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| timeline read | données corrélées | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| correlate/compare | agrégation sensible | 0/1 | policy | initiateur/reviewer | Investigate/Shared | Permissions |
| annotate/handoff | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- Memory Timeline ≠ Case Timeline.
- timestamp absent non fabriqué; estimé clairement signalé.
- conflit/trou temporel visible.
- aucune timeline Disk/Network complète.

## 23. Métriques
- événements par qualité temporelle.
- trous/conflits et corrélations contestées.
- images comparées.
- handoffs avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou moteur revendiqué.

## 25. Critères d’acceptation
### 1. Timestamp absent
**Given** observation sans timestamp **When** affichée **Then** `absent` est visible et aucun horaire n’est inventé.
### 2. Corrélation contradictoire
**Given** sources divergent **When** corrélées **Then** conflit et sources restent visibles.
### 3. Sans IA
**Given** aucun modèle **When** timeline construite **Then** ordering, filtres, comparaison et annotations déterministes fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; schémas/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, Case Timeline, CAP-INV-349..362, Evidence/Replay et future 4B.2B.3B.
