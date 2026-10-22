---
id: CAP-INV-357
title: Injection, Hollowing and Memory Anomaly Analysis
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
  - REQ-PROD-052
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-357 — Injection, Hollowing and Memory Anomaly Analysis

## 1. Définition
Regrouper et examiner défensivement des anomalies mémoire candidates en exposant les éléments pour et contre, sans procédure d’injection/hollowing ni confirmation automatique.

## 2. Problème utilisateur
Sans analyse contradictoire, une protection inhabituelle ou une incohérence de mapping peut être qualifiée à tort d’injection.

## 3. Objectifs
- comparer mappings, modules, régions, protections, vues et snapshots.
- présenter éléments favorables, contradictoires et manquants.
- préparer Hypothesis, Evidence candidate ou handoff Reverse/Debugger sans confirmation automatique.

## 4. Non-objectifs
Aucune technique procédurale d’injection/hollowing, code, commande, bypass, évasion, moteur, API ou action Endpoint.

## 5. Propriétaire
Investigate possède anomaly/interprétation; Static/Reverse/Debugger restent propriétaires de leurs analyses; Studio Tools/Runs; Shared trace.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Reverse Engineer, Debug Analyst, Evidence Reviewer.

## 7. Conditions d’entrée
Image/session/profil, régions/modules/processus et sources comparables visibles; partialité explicite.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Region/process/module observations | CAP-INV-351..353 | éléments de contexte | oui | même image/session | inconclusive |
| Tool anomaly results | Tool Calls | candidates/scores/sources | non | Tool/version visibles | manual review |
| Extracted Artifacts | CAP-INV-360 | contenu candidat | non | lineage courant | no handoff |
| Hypothesis/Case | Investigate | question analytique | oui pour handoff | courant | unlinked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image/Observations | Investigate | source/relations | lecture |
| Tool/Tool Call/Run | Studio | producteur/version/paramètres | lecture |
| Derived Artifact | Investigate | extraction/lineage | lecture/lien |
| Hypothesis | Investigate | question/relations | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Anomaly | créer/annoter/contester/superseder | Investigate concept | candidate ≠ Finding |
| Hypothesis relation | lier | Investigate | source/incertitude obligatoires |
| Reverse/Debugger handoff | préparer | Investigate | aucune action réelle |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- grouper anomalies candidates et comparer mappings/modules/protections.
- voir relations process/région et différences entre vues/snapshots.
- voir Artifacts extraits et éléments pour/contre.
- classer candidate, weakly-supported, supported, contradicted ou inconclusive.
- préparer Hypothesis/Evidence/Reverse/Debugger handoff.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/comparer | Analyst | Memory Anomaly | 0/1 | sources lisibles | dossier contradictoire | non |
| Classer/annoter | Analyst | Memory Anomaly | 2 | permission/justification | version conservée | OPEN-013 |
| Extraire | Analyst | Derived Artifact request | 1 | policy | résultat borné | non |
| Préparer handoff | Analyst | Package | 2 | provenance complète | candidate/draft | non |

Classes 3/4 indisponibles.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| grouper anomalies | oui | oui | oui | proposition | règles explicables |
| comparer vues | oui | oui | oui | résumé | comparateurs/tables |
| expliquer éléments pour/contre | oui | oui | oui | assistance | checklist humaine |
| confirmer injection | oui | non | non | non automatique | revue humaine |

Attribution et disposition humaine obligatoires.

## 14. États fonctionnels
`candidate`, `weakly-supported`, `supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. États objet finaux reportés.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale; aucune anomalie fictive ou confirmation silencieuse.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Memory Anomaly | candidate | Workbench/Hypothesis | éléments pour/contre visibles |
| Extraction/handoff selection | context | CAP-INV-360/362/Reverse | lineage conservé |
| Disposition | versioned event | Case/Replay | humain et justification visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Region/module | analyser anomalie | CAP-INV-357 | protections, relations, comparisons | source |
| Anomaly | extraire | CAP-INV-360 | source location, process, restrictions | Anomaly |
| Anomaly | handoff | Reverse/Debugger/CAP-INV-362 | candidates, evidence, contradictions | Anomaly |

Tenant, Case, image, permissions et return origin préservés.

## 18. Dépendances
CAP-INV-351/352/353/360/362, CAP-INV-329..346, Studio, Shared, OPEN-005/008/013/015.

## 19. Source de vérité
Sources/observations/anomaly : Investigate; Tool/Run : Studio; Reverse/Debugger propriétaires de leurs sessions.

## 20. Provenance et audit
Image/session/profil, Tool/version, rules/sources, régions/modules/processus, comparisons, Artifact lineage, acteur, état, erreurs et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| anomaly review | lecture sensible | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| compare/extract | traitement | 1 | policy | initiateur/reviewer | Investigate | Permissions |
| classify/handoff | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- région/anomalie ≠ injection confirmée.
- score/Tool/IA ≠ conclusion.
- aucune procédure, code, bypass ou évasion.
- partialité et contradictions visibles.

## 23. Métriques
- anomalies par état/disposition.
- taux de contradictions et partialité.
- handoffs avec provenance complète.
- suggestions acceptées/modifiées/rejetées.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou technique offensive revendiquée.

## 25. Critères d’acceptation
### 1. Région inhabituelle
**Given** région inhabituelle sans module/comportement runtime **When** examinée **Then** elle reste candidate, éléments pour/contre visibles, aucune injection confirmée.
### 2. Suggestion automatisée
**Given** agent propose anomalie **When** consultée **Then** agent/sources/Calls visibles et analyste accepte/modifie/rejette.
### 3. Sans IA
**Given** aucun modèle **When** analyse réalisée **Then** règles, comparateurs, extraction et revue humaine fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; schémas/plateformes/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351..353/360/362, Reverse/Debugger, Case/Evidence et futures phases Permissions/Screens.
