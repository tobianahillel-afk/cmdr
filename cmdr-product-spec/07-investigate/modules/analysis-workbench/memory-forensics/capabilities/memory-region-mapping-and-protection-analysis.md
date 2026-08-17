---
id: CAP-INV-352
title: Memory Region, Mapping and Protection Analysis
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
# CAP-INV-352 — Memory Region, Mapping and Protection Analysis

## 1. Définition
Permet d’examiner régions, mappings, protections, tailles, relations et incohérences sans qualifier automatiquement une injection, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans vue sourcée des régions, mappings et protections, une zone exécutable, modifiée ou non mappée peut être prise à tort pour une injection confirmée.

## 3. Objectifs
- inspecter régions, mappings, protections, tailles et relations.
- préserver sources, partialité, incertitude et return origin.
- permettre comparaison, annotation et extraction bornée.

## 4. Non-objectifs
Aucune méthode d’injection, signature, offset, commande, moteur, API, action Endpoint ou Memory/Disk/Network Forensics hors périmètre.

## 5. Propriétaire
Investigate possède observations/interprétations; Studio Tools/Runs; Settings policies/stockage; Govern cibles réelles; Shared mécanismes.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Reverse Engineer, Evidence Reviewer et Audit Analyst autorisés.

## 7. Conditions d’entrée
Image/session/profil lisibles, Tool/version et permissions visibles; partialité héritée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Platform/Profile | CAP-INV-350 | interprétation/limites | oui | courant | unsupported |
| Region/mapping results | Tool Calls | projections | oui | Tool/version visibles | partial/failed |
| Policies | Settings/Security | lecture/extraction | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture |
| Process/Module Observation | Investigate concepts | relations | lecture/lien |
| Tool/Tool Call/Run | Studio | producteur/version/statut | lecture |
| Case/Hypothesis | Investigate | contexte | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Region | créer/annoter/contester | Investigate concept | source/incertitude obligatoires |
| Mapping relation | créer/versionner | Investigate concept | pas de cardinalité finale |
| Extraction request | préparer | Investigate | via CAP-INV-360 |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- voir régions, mappings, tailles et protections déclarées.
- relier régions aux processus/modules et voir zones partagées/non mappées/ambiguës.
- filtrer, comparer et annoter.
- sélectionner une région et préparer une extraction.
- conserver les éléments pour et contre une anomalie.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Memory Region | 0 | lecture autorisée | vue sourcée | non |
| Comparer | Analyst | Region set | 0/1 | sources compatibles | différences/limites | non |
| Annoter/contester | Analyst | Observation | 2 | permission | version conservée | OPEN-013 |
| Extraire | Analyst | Derived Artifact request | 1 | policy/permission | résultat borné | non |

Classes 3/4 bloquées et routées hors module.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| afficher régions/mappings | oui | oui | oui | explication | viewer déterministe |
| relier processus/modules | oui | oui | oui | suggestion | tables de relations |
| comparer snapshots | oui | oui | oui | résumé | comparateur déterministe |
| qualifier injection | oui | non | non | assistance | revue humaine |

Attribution complète de toute automatisation.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale conservent contexte et n’inventent aucune donnée.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Region projection | Memory Region result | Workbench | source/partialité visibles |
| Mapping relation | observation | CAP-INV-351/353/357 | pas d’injection automatique |
| Extraction selection | request | CAP-INV-360 | source/restrictions conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Process | inspecter régions | CAP-INV-352 | image, process, profil, Tool | Processes |
| Region | analyser anomalie | CAP-INV-357 | protections, relations, comparaisons | Regions |
| Region | extraire | CAP-INV-360 | source location, process, restrictions | Regions |

Tenant, Case, image, permissions, sélection et return origin préservés.

## 18. Dépendances
CAP-INV-350/351/357/360, Studio, Settings, Shared, OPEN-005/008/013/014/015.

## 19. Source de vérité
Image et observations : Investigate; Tool/Run : Studio; policies : Settings; trace/export : Shared.

## 20. Provenance et audit
Image/session/profil, Tool/version, paramètres, région telle qu’affichée, process/module, protections déclarées, comparison, extraction, acteur, statut et erreurs.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| region read | contenu sensible | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| comparison | corrélation | 0/1 | policy | initiateur/reviewer | Investigate | Permissions |
| extraction | diffusion | 1 | probable | extractor/reviewer | Investigate/Shared | Permissions |
| annotation | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Matrice atomique et modèles d’accès finaux reportés.

## 22. Limites et erreurs
- région exécutable/modifiée/non mappée ≠ injection confirmée.
- aucun offset, signature ou moteur imposé.
- image source immuable.
- Debugger Memory View et Memory Forensics restent distincts.

## 23. Métriques
- résultats partial/failed.
- régions avec source/protection complète.
- comparaisons et extractions tracées.
- candidates acceptées/modifiées/rejetées.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou intégration revendiquée.

## 25. Critères d’acceptation
### 1. Région inhabituelle
**Given** région exécutable sans module confirmé **When** examinée **Then** elle reste candidate, éléments pour/contre visibles, aucune injection confirmée.
### 2. Région non mappée
**Given** mapping absent **When** ouverte **Then** `unmapped/partial` est visible et rien n’est inventé.
### 3. Sans IA
**Given** aucun modèle **When** l’analyste travaille **Then** viewers, tables, comparateurs et extraction autorisée fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 restent ouvertes; schémas/plateformes/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351/353/357/360/362, Static/Reverse, futures phases Objects/Permissions/Screens.
