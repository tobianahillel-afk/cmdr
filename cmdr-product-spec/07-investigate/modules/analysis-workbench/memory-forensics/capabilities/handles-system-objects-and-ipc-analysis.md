---
id: CAP-INV-354
title: Handles, System Objects and IPC Analysis
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
# CAP-INV-354 — Handles, System Objects and IPC Analysis

## 1. Définition
Permet d’examiner handles, objets système reconstruits, ressources partagées et relations IPC candidates sans les confondre avec les objets canoniques CMDR, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans reconstruction des handles, objets système et IPC, les relations techniques restent opaques et peuvent être confondues avec des objets CMDR ou une coordination malveillante certaine.

## 3. Objectifs
- examiner handles, objets système, owners, permissions, ressources partagées et IPC candidates.
- préserver relations processus/fichiers/régions, incertitude et provenance.
- distinguer objets kernel/système et objets canoniques CMDR.

## 4. Non-objectifs
Aucune commande, méthode offensive, moteur, plugin, offset, API, action Endpoint ou qualification automatique d’une coordination malveillante.

## 5. Propriétaire
Investigate possède observations/interprétations; Studio Tools/Runs; Settings policies; Govern cibles réelles; Shared mécanismes.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**; secondaires : Investigation Lead, Evidence Reviewer et Audit Analyst.

## 7. Conditions d’entrée
Image/session/profil/Tool et permissions lisibles; structures absentes ou incompatibles visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image/Session | Investigate | source/scope | oui | versions liées | blocked |
| Platform/Profile | CAP-INV-350 | interprétation/limites | oui | courant | unsupported |
| Handle/Object/IPC results | Tool Calls | projections | oui | Tool/version visibles | partial/failed |
| Policies | Settings/Security | lecture/export | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | source/limites | lecture |
| Process/Region Observation | Investigate concepts | relations | lecture/lien |
| Tool/Tool Call/Run | Studio | version/paramètres/statut | lecture |
| Case/Hypothesis | Investigate | contexte | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Handle Observation | créer/annoter/contester | Investigate concept | source/incertitude obligatoires |
| System Object Observation | créer/lier/versionner | Investigate concept | ≠ objet CMDR |
| IPC Relation | créer/annoter | Investigate concept | candidate seulement |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- voir handles/références, owners et permissions disponibles.
- voir objets système reconstruits et relations entre processus.
- voir fichiers, ressources partagées et IPC candidates.
- voir objets orphelins/incohérents, filtrer, comparer et annoter.
- naviguer vers processus, fichiers et régions.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter/filtrer | Analyst | Handle/Object | 0 | lecture autorisée | vue sourcée | non |
| Comparer | Analyst | Observation set | 0/1 | sources compatibles | divergences visibles | non |
| Annoter/relier | Analyst | Observation/IPC | 2 | permission | version conservée | OPEN-013 |
| Lancer traitement borné | Analyst | Tool Call | 1 | scope/Tool/policy | job tracé | non |

Classes 3/4 bloquées et routées hors module.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inventorier handles/objets | oui | oui | oui | résumé | tables déterministes |
| représenter relations | oui | oui | oui | suggestion | graphe + table alternative |
| filtrer IPC/ressources | oui | oui | oui | groupement | filtres déterministes |
| conclure malveillance | oui | non | non | assistance | revue humaine |

Attribution complète de toute automatisation.

## 14. États fonctionnels
`queued`, `processing`, `partial`, `available`, `failed`, `incompatible`, `disputed`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading/Empty/Partial/Error/Offline/Permission denied/Stale conservent contexte et n’inventent rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Handle Observation | result | Workbench | source/partialité visibles |
| System Object/IPC Observation | candidate relation | Case/Hypothesis | aucune qualification automatique |
| Handoff selection | candidate context | CAP-INV-362 | provenance conservée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Process | inspecter handles/IPC | CAP-INV-354 | process, image, profil, Tool | Processes |
| Observation | naviguer | Process/Region/Artifact | relation, sélection, provenance | Objects |
| Observation | handoff | CAP-INV-362 | candidates, contradictions, sources | Objects |

Tenant, Case, image, permissions et return origin préservés.

## 18. Dépendances
CAP-INV-350/351/352/359/362, Studio, Shared, OPEN-005/008/013/015.

## 19. Source de vérité
Image/observations : Investigate; Tool/Run : Studio; policies : Settings; trace/linking : Shared.

## 20. Provenance et audit
Image/session/profil, Tool/version, paramètres, handle/object tels qu’affichés, process/resource relations, acteur, erreurs, partialité et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| handle/object read | données système | 0 | possible | reviewer si requis | Investigate/Security | Permissions |
| relation comparison | corrélation | 0/1 | policy | initiateur/reviewer | Investigate | Permissions |
| annotation/link | mutation | 2 | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Modèle d’accès final reporté.

## 22. Limites et erreurs
- handle ≠ action métier.
- objet kernel/système ≠ objet canonique CMDR.
- IPC ≠ coordination malveillante.
- structures manquantes et image partielle visibles.

## 23. Métriques
- observations partial/failed.
- objets/relations avec provenance complète.
- incohérences contestées.
- handoffs qualifiés par owner.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou moteur revendiqué.

## 25. Critères d’acceptation
### 1. IPC candidate
**Given** deux processus partagent une ressource **When** inspectée **Then** la relation reste candidate et n’implique pas une coordination malveillante.
### 2. Objet inconnu
**Given** un objet est partiellement reconstruit **When** affiché **Then** les champs absents et la source sont visibles.
### 3. Sans IA
**Given** aucun modèle **When** l’analyste travaille **Then** tables, graphes structurés, filtres et annotations fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; schémas/plateformes/permissions reportés.

## 27. Consommateurs documentaires
INV-MEM-001, CAP-INV-351/352/359/362, Case/Evidence et futures phases Objects/Permissions/Screens.
