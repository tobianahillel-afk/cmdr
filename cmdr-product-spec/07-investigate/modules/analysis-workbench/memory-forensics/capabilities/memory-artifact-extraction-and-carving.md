---
id: CAP-INV-360
title: Memory Artifact Extraction and Carving
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
  - REQ-OBJ-003
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-360 — Memory Artifact Extraction and Carving

## 1. Définition
Extraire de manière bornée un Derived Artifact ou une représentation depuis une région/structure, avec parenté, contexte, limites et provenance, sans modifier l’image source.

## 2. Problème utilisateur
Sans extraction bornée et lineage, un contenu partiel peut être présenté comme récupération certaine, perdre son contexte ou altérer la source.

## 3. Objectifs
- sélectionner région/structure source, objectif et restrictions.
- produire module/contenu/représentation candidat avec parenté et contexte.
- conserver Tool/version/paramètres/erreurs et router Static/Reverse.

## 4. Non-objectifs
Aucune modification de l’image, récupération certaine, format/moteur imposé, commande, carving disque, méthode offensive ou déploiement.

## 5. Propriétaire
Investigate possède Derived Artifact/extraction context; Static/Reverse consomment; Studio Tools/Runs; Settings/Shared policies/stockage/export.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Static Analyst, Reverse Engineer, Evidence Reviewer.

## 7. Conditions d’entrée
Source sélectionnée, permission/policy, image immutable, Tool/version et objectif visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/region/structure | Investigate | source exacte | oui | version sélectionnée | blocked |
| Process/module context | CAP-INV-351..353 | lineage | non | même image | context-partial |
| Tool/version/parameters | Studio | extraction bornée | oui | sélection courante | tool-unavailable |
| Permissions/restrictions | Security/Settings | extraction/export | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source immutable | lecture |
| Region/Module/Process Observation | Investigate concepts | source context | lecture/lien |
| Tool/Tool Call/Run | Studio | version/parameters/status | lecture |
| Artifact policy | Settings/Security | restrictions | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Derived Artifact | créer/superseder/withdraw | Investigate concept | parent image obligatoire |
| Extraction Result | créer/annoter | Investigate concept | partialité/erreurs visibles |
| Static/Reverse handoff | préparer | Investigate | aucun Evidence automatique |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- sélectionner source/objectif et afficher limites/restrictions.
- extraire module, contenu ou représentation candidat.
- conserver parent/enfant, process/région, Tool/version/paramètres/erreurs.
- gérer résultat partial/invalid/restricted et comparer avec source.
- router Derived Artifact vers Static/Reverse et retirer de l’usage actif sans supprimer trace.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Sélectionner source | Analyst | Source region | 0/2 | lecture/permission | scope versionné | OPEN-013 |
| Extraire | Analyst | Derived Artifact | 1 | Tool/policy | résultat borné | non |
| Comparer | Analyst | Source/result | 0/1 | résultats disponibles | différences/limites | non |
| Withdraw/supersede | Analyst | Derived Artifact | 2 | permission | trace conservée | OPEN-013 |

Classes 3/4 indisponibles; aucun déploiement.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| sélectionner source/scope | oui | oui | oui | suggestion | sélection manuelle |
| extraire | oui | oui | oui | non nécessaire | extracteur déterministe |
| comparer source/résultat | oui | oui | oui | résumé | comparateur |
| qualifier récupération | oui | non | non | assistance | revue humaine |

Attribution complète de l’extraction.

## 14. États fonctionnels
`proposed`, `extracting`, `available`, `partial`, `invalid`, `restricted`, `superseded`, `withdrawn-from-use`. États objet finaux reportés.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale; source et résultat restent distincts.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Derived Artifact | Artifact relation | Static/Reverse | source/parenté/partialité visibles |
| Extraction Result | result | Workbench/Case | Tool/version/parameters/errors conservés |
| Handoff selection | package | CAP-INV-362 | aucune Evidence automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Region/Module | extraire | CAP-INV-360 | image, source location, process, Tool, restrictions | source |
| Derived Artifact | analyser | Static/Reverse | parent image, extraction, Case, provenance | Extraction |
| Extraction | handoff | CAP-INV-362 | Artifact, source, limitations | Extraction |

Tenant, Case, image et return origin préservés.

## 18. Dépendances
CAP-INV-311/329/352/353/357/362, Studio, Settings, Shared Export/Versioning, OPEN-005/008/013/014/015.

## 19. Source de vérité
Image/source : Investigate; Derived Artifact relation : Investigate; Tool/Run : Studio; export/versioning : Shared.

## 20. Provenance et audit
Parent image, region/structure/process/module, Tool/version, parameters, initiator, timestamps, errors, partiality, Derived Artifact hash/reference, route and disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| source read | contenu sensible | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| extract | capture contenu | 1 | policy | extractor/reviewer | Investigate | Permissions |
| export | diffusion | 1/2 | probable | exporter/approver | Shared/Security | Permissions |
| withdraw/link | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- carving ≠ récupération certaine.
- image source immutable.
- result partial/invalid/restricted visible.
- aucune extraction disque, commande ou format imposé.

## 23. Métriques
- extractions available/partial/invalid.
- parenté/contexte complets.
- routes Static/Reverse et exports autorisés/refusés.
- withdrawals/supersessions tracés.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou moteur revendiqué.

## 25. Critères d’acceptation
### 1. Extraction module
**Given** région/module et permission **When** extrait **Then** image inchangée, parenté/process/région/Tool visibles, partialité signalée et route Static/Reverse possible.
### 2. Résultat partiel
**Given** contenu incomplet **When** extraction termine **Then** `partial` est visible et aucune récupération certaine n’est déclarée.
### 3. Sans IA
**Given** aucun modèle **When** extraction réalisée **Then** sélection, extracteur, comparaison et routage déterministes fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 restent ouvertes; formats/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-311/329/357/362, Static/Reverse, Case/Evidence et futures phases Objects/Permissions/Screens.
