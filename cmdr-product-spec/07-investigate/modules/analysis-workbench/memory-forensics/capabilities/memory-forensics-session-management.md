---
id: CAP-INV-348
title: Memory Forensics Session Management
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
  - REQ-OBJ-009
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-348 — Memory Forensics Session Management

## 1. Définition
Permet de gérer un contexte analytique durable pour une ou plusieurs Memory Images, leurs paramètres, observations, erreurs, Derived Artifacts et handoffs, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans session durable, les profils, filtres, observations, erreurs et Derived Artifacts sont dispersés; plusieurs analystes ne peuvent ni reprendre ni comparer le travail de façon fiable.

## 3. Objectifs
- gérer un contexte analytique durable pour une ou plusieurs Memory Images, leurs paramètres, observations, erreurs, Derived Artifacts et handoffs.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**. Secondaires : Investigation Lead, Evidence Reviewer, Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles; Tool/profil sélectionné explicitement. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Intake accepté | CAP-INV-347 | image, objectif et restrictions | oui | version courante | session draft |
| Case et Memory Images | Investigate | contexte et sources | oui | versions sélectionnées | blocked |
| Profil, Tools et paramètres | CAP-INV-350 / Studio | configuration analytique | non à la création | versions enregistrées | profile-required |
| Contributeurs et permissions | Security | collaboration | oui | courant | denied |
| Résultats et erreurs existants | sessions antérieures | reprise | non | snapshot de session | nouvelle session vide |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Memory Image / Artifact | Investigate | sources et relations | lecture/lien |
| Tool / Tool Call / Automation Run | CMDR Studio | versions, paramètres et statut | lecture |
| Platform/Profile selection | Investigate/Settings projections | contexte d’analyse | lecture |
| Derived Artifact | Investigate concept | sorties liées | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Forensics Session | créer/modifier/pause/close/reopen/superseder | Investigate concept | aucun schéma final |
| Session membership | ajouter ou retirer un contributeur | Investigate concept | permission requise |
| Observation / Derived Artifact relation | lier ou retirer de l’usage actif | Investigate | historique conservé |
| Trace / Activity event | émettre | Shared | append-only |

## 11. Fonctionnalités
- créer, reprendre, suspendre, clôturer et rouvrir une session.
- lier Case et plusieurs images.
- conserver plateforme, profil, Tools, Tool Calls, paramètres, vues, requêtes, filtres et observations.
- comparer ou superseder des sessions.
- transmettre les résultats sans fusionner les objets.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Créer session | Memory Forensics Analyst | Memory Forensics Session | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Reprendre ou suspendre | Memory Forensics Analyst | Memory Forensics Session | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Clôturer ou rouvrir | Memory Forensics Analyst | Memory Forensics Session | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Comparer sessions | Memory Forensics Analyst | Session set | 0 | lecture autorisée | vue sourcée | non |
| Lancer une analyse bornée | Memory Forensics Analyst | Tool Call | 1 | scope/Tool/policy explicites | job/extraction borné et tracé | non |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| créer, reprendre, suspendre, clôturer et rouvrir une session | oui | oui | oui | proposition | formulaire de session et versioning déterministe |
| lier Case et plusieurs images | oui | oui | oui | proposition | historique de vues, filtres et sélections |
| conserver plateforme, profil, Tools, Tool Calls, paramètres, vues, requêtes, filtres et observations | oui | oui | oui | proposition | comparateur de sessions |
| Conclusion finale | oui | non | non | assistance | revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `archived`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Memory Forensics Session | concept fonctionnel | Memory Workbench | sources, owner, versions et état conservés |
| Session summary | projection | Case Replay / reviewers | résultats et erreurs séparés |
| Supersession relation | relation versionnée | future Objects phase | ancienne session accessible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-347 | créer ou reprendre | CAP-INV-348 | Case, images, objectif, restrictions | Intake |
| Session | sélectionner profil | CAP-INV-350 | images, candidats, Tools, résultats antérieurs | Session |
| Session | ouvrir domaine analytique | CAP-INV-351..360 | session, profil, sélection et filtres | Session |
| Session | préparer handoff | CAP-INV-362 | observations, contradictions et provenance | Session |

Tenant, Case, image, permissions, restrictions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-347, CAP-INV-350, CAP-INV-361, CMDR Studio, Shared Versioning, Shared Recovery, OPEN-005, OPEN-008, OPEN-013, OPEN-015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline/Export/Versioning : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, session, profil, Tool/version, Calls/Run, paramètres, filtres, acteur, timestamps, erreurs, interruptions, données sensibles, Derived Artifacts et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Créer session | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Reprendre ou suspendre | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Clôturer ou rouvrir | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Comparer sessions | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Lancer une analyse bornée | traitement borné | 1 | selon policy | SoD si requise | Investigate/Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- session distincte de Collection Job, Analysis Session, Reverse Session, Debugger Session et Automation Run.
- aucun Tool choisi.
- états fonctionnels seulement.
- aucune fusion silencieuse de sessions.
- Memory Forensics ≠ Debugger Memory View ≠ Disk Forensics ≠ full Network Forensics.
- Une erreur ou incohérence ne devient pas une conclusion.

## 23. Métriques
- sessions par état.
- reprises avec contexte restauré.
- erreurs conservées.
- sessions superseded avec relation complète.
- complétude de provenance et retour au contexte sans perte.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin, plateforme ou intégration native.

## 25. Critères d’acceptation
### 1. Reprise
**Given** une session paused **When** un contributeur autorisé reprend **Then** vues, filtres, observations et erreurs sont restaurés.
### 2. Échec Tool
**Given** un Tool échoue **When** la session continue **Then** les résultats antérieurs restent visibles et l’erreur est attribuée.
### 3. Sans IA
**Given** aucun modèle **When** la session est créée **Then** Tools déterministes, filtres et annotations fonctionnent.

## 26. Questions ouvertes
OPEN-005, OPEN-008, OPEN-013, OPEN-015 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, Static/Reverse, futures phases Objects/Permissions/Journeys/Screens/Contracts, handoffs seulement vers 4B.2B.3B/4B.3.
