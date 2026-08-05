---
id: CAP-INV-353
title: Loaded Modules, Images and Driver Analysis
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
# CAP-INV-353 — Loaded Modules, Images and Driver Analysis

## 1. Définition
Permet d’examiner modules, images et drivers observables, leurs relations, métadonnées, incohérences et possibilités d’extraction sans conclure à la malveillance, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans analyse corrélée des modules, images et drivers, l’absence de métadonnée ou une divergence entre vues peut être interprétée à tort comme dissimulation ou malveillance.

## 3. Objectifs
- examiner modules, images et drivers observables avec leurs relations et métadonnées.
- comparer plusieurs vues ou snapshots et conserver les incohérences.
- relier Static, Reverse et Dynamic sans conclusion automatique.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, API, action Endpoint ou qualification automatique de malveillance/rootkit.

## 5. Propriétaire
Investigate possède les observations/interprétations; Studio Tools/Runs; Settings policies/stockage; Govern cibles réelles; Shared mécanismes.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Reverse Engineer, Evidence Reviewer et Investigation Lead.

## 7. Conditions d’entrée
Image/session/profil/Tool et permissions lisibles; métadonnées manquantes et support de plateforme visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Platform/Profile | CAP-INV-350 | interprétation/limites | oui | courant | unsupported |
| Module/driver results | Tool Calls | projections | oui | Tool/version visibles | partial/failed |
| Policies | Settings/Security | lecture/extraction | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture |
| Process/Region Observation | Investigate concepts | relations | lecture/lien |
| Tool/Tool Call/Run | Studio | version/paramètres/statut | lecture |
| Static/Reverse/Dynamic result | Investigate | comparaison/handoff | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Module Observation | créer/annoter/contester | Investigate concept | source/incertitude obligatoires |
| Driver Observation | créer/lier/versionner | Investigate concept | présent ≠ rootkit |
| Extraction request | préparer | Investigate | via CAP-INV-360 |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- voir modules, images, drivers/composants équivalents et processus associés.
- voir emplacements, versions, symboles, signatures/métadonnées et chemins disponibles.
- voir incohérences et comparer vues/snapshots.
- annoter, extraire un Artifact et relier Static/Reverse/Dynamic.
- exposer les métadonnées absentes sans conclure à une dissimulation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Module/Driver Observation | 0 | lecture autorisée | vue sourcée | non |
| Comparer | Analyst | Observation set | 0/1 | sources compatibles | divergences visibles | non |
| Annoter/relier | Analyst | Observation | 2 | permission | version conservée | OPEN-013 |
| Extraire | Analyst | Derived Artifact request | 1 | policy/permission | résultat borné | non |

Classes 3/4 bloquées et routées hors module.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inventorier modules/drivers | oui | oui | oui | résumé | inventaire déterministe |
| afficher métadonnées/relations | oui | oui | oui | explication | tables sourcées |
| comparer snapshots | oui | oui | oui | résumé | comparateur déterministe |
| conclure malveillance/rootkit | oui | non | non | assistance | revue humaine |

Attribution complète de toute automatisation.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale conservent contexte et n’inventent rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Module Observation | result | Workbench/Static/Reverse | source/partialité visibles |
| Driver Observation | candidate | Kernel/Case/Hypothesis | aucun rootkit automatique |
| Extraction selection | request | CAP-INV-360 | parenté/restrictions conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Process/Session | inspecter modules | CAP-INV-353 | image, process, profil, Tool | source |
| Module/Driver | inspecter kernel/anomalie | CAP-INV-357/358 | observation, région, incohérences | Modules |
| Module | extraire/analyser | CAP-INV-360 puis Static/Reverse | source, process, Tool, restrictions | Modules |

Tenant, Case, image, permissions et return origin préservés.

## 18. Dépendances
CAP-INV-350/351/352/357/358/360/362, Static/Reverse/Dynamic, Studio, Shared, OPEN-005/008/013/014/015.

## 19. Source de vérité
Image/observations : Investigate; Tool/Run : Studio; policies : Settings; trace/export : Shared.

## 20. Provenance et audit
Image/session/profil, Tool/version, paramètres, module/driver tel qu’affiché, process/région, métadonnées disponibles/absentes, extraction, acteur, statut et erreurs.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| module/driver read | données système | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| compare | corrélation | 0/1 | policy | initiateur/reviewer | Investigate | Permissions |
| extract/export | diffusion | 1/2 | probable | extractor/reviewer | Investigate/Shared | Permissions |
| annotation | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- module chargé ≠ module malveillant.
- driver présent ≠ rootkit.
- absence de métadonnée ≠ dissimulation.
- aucun moteur, liste de plateforme ou signature imposés.

## 23. Métriques
- observations partial/failed.
- modules/drivers avec métadonnées/provenance complètes.
- divergences multi-vue.
- extractions tracées.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou moteur revendiqué.

## 25. Critères d’acceptation
### 1. Métadonnée absente
**Given** chemin/signature/version absente **When** inspectée **Then** l’absence est visible sans preuve de dissimulation ou malveillance.
### 2. Vues divergentes
**Given** module présent dans une vue seulement **When** comparé **Then** l’incohérence reste candidate et la partialité est visible.
### 3. Sans IA
**Given** aucun modèle **When** l’analyste travaille **Then** inventaires, tables, comparaison et extraction autorisée fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 restent ouvertes; schémas/plateformes/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351/352/357/358/360/362, Static/Reverse/Dynamic et futures phases Objects/Permissions/Screens.
