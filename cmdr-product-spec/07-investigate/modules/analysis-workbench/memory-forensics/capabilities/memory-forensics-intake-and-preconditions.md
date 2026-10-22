---
id: CAP-INV-347
title: Memory Forensics Intake and Preconditions
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
  - REQ-PROD-052
  - REQ-PROD-055
  - REQ-SEC-001
  - REQ-UX-002
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-347 — Memory Forensics Intake and Preconditions

## 1. Définition
Permet d’ouvrir une Memory Image depuis un Case, Collection Job ou Artifact Management, vérifier les préconditions et créer ou reprendre une Memory Forensics Session sans lancer d’analyse automatiquement, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans intake explicite, une image partielle, restreinte ou incompatible peut être ouverte comme si elle était exploitable, un Tool peut être lancé hors contexte et le retour au Case peut être perdu.

## 3. Objectifs
- ouvrir une Memory Image depuis un Case, Collection Job ou Artifact Management, vérifier les préconditions et créer ou reprendre une Memory Forensics Session sans lancer d’analyse automatiquement.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **DFIR Analyst**. Secondaires : Investigation Lead, Evidence Reviewer, Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles; Tool/profil sélectionné explicitement. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case, Endpoint et objectif | Investigate / Platform Settings | contexte et scope | oui | version courante | intake incomplete |
| Memory Image et relation Artifact | Artifact Management / Memory Image | source d’analyse | oui | version sélectionnée | aucune session créée |
| Acquisition context, custody et intégrité | CAP-INV-207/213/214 | provenance et restrictions | oui | dernier état attribué | partial ou disputed |
| Tools et plateformes candidates | CMDR Studio / Platform Settings | compatibilités déclarées | non pour consultation | état courant | tool-unavailable ou profile-required |
| Permissions et return origin | Security / Experience | accès et navigation | oui | à l’ouverture | denied sans fuite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, contexte et return origin | lecture/lien |
| Memory Image / Artifact | Investigate | source, taille, version, restrictions | lecture |
| Collection Request / Job | Investigate | acquisition et résultat | lecture |
| Endpoint / Agent / Policy | Endpoint Agent / Settings | source, support et restrictions | projection lecture |
| Tool / Tool Call history | CMDR Studio | compatibilité et résultats précédents | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Memory Forensics Intake | créer/modifier/superseder | Investigate | aucune analyse automatique |
| Memory Forensics Session relation | préparer ou reprendre | Investigate concept | image, Case et owner requis |
| Trace / Activity event | émettre | Shared | append-only |
| Source objects | aucune mutation | owners respectifs | projections seulement |

## 11. Fonctionnalités
- ouvrir depuis Case, Collection Job ou Artifact Management.
- afficher Endpoint source, acquisition, provenance, taille, intégrité et restrictions.
- distinguer plateforme déclarée, détectée et profils candidats.
- afficher Tools compatibles ou indisponibles et permissions.
- définir objectif et scope, créer ou reprendre une session, préserver le return origin.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter préconditions | DFIR Analyst | Intake | 0 | lecture autorisée | vue sourcée | non |
| Définir objectif et scope | DFIR Analyst | Intake | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Créer ou reprendre une session | DFIR Analyst | Memory Forensics Session | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Ouvrir explicitement un Tool autorisé | DFIR Analyst | Tool Call request | 1 | scope/Tool/policy explicites | job/extraction borné et tracé | non |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| ouvrir depuis Case, Collection Job ou Artifact Management | oui | oui | oui | proposition | checklist d’intake et projections déterministes |
| afficher Endpoint source, acquisition, provenance, taille, intégrité et restrictions | oui | oui | oui | proposition | table de provenance/custody et viewer de métadonnées |
| distinguer plateforme déclarée, détectée et profils candidats | oui | oui | oui | proposition | matrice déclarée de compatibilité Tool/plateforme |
| Conclusion finale | oui | non | non | assistance | revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `partial-image`, `corrupted`, `restricted`, `unsupported`, `profile-required`, `tool-unavailable`, `policy-blocked`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Intake assessment | Memory Forensics Intake | CAP-INV-348 | complétude et limites visibles |
| Blocking or partial reasons | Trace event | Case Workspace | aucune donnée inventée |
| Session context | Memory Forensics Session concept | Memory Workbench | Case, image et return origin conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case / Artifact | ouvrir Memory Forensics | CAP-INV-347 | tenant, Case, Endpoint, Memory Image, objectif, Hypothesis | source |
| Collection Job | ouvrir résultat | CAP-INV-347 | Job, acquisition context, errors, custody, Artifact relation | Collection |
| CAP-INV-347 | préconditions acceptées | CAP-INV-348 | image, restrictions, owner, objectif, return origin | Intake |
| CAP-INV-347 | image à revoir | CAP-INV-349 | source, intégrité, transformations et custody | Intake |

Tenant, Case, image, permissions, restrictions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-105, CAP-INV-203, CAP-INV-207, CAP-INV-213, CAP-INV-214, CMDR Studio Tool catalog, Platform Settings, Shared Trace, OPEN-005, OPEN-008, OPEN-013, OPEN-014, OPEN-015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline/Export/Versioning : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, session, profil, Tool/version, Calls/Run, paramètres, filtres, acteur, timestamps, erreurs, interruptions, données sensibles, Derived Artifacts et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Inspecter préconditions | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Définir objectif et scope | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Créer ou reprendre une session | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Ouvrir explicitement un Tool autorisé | traitement borné | 1 | selon policy | SoD si requise | Investigate/Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- aucun Tool lancé à l’ouverture.
- image partielle ou inexploitable visible.
- plateforme déclarée distincte de détectée.
- aucune acquisition locale.
- Memory Forensics ≠ Debugger Memory View ≠ Disk Forensics ≠ full Network Forensics.
- Une erreur ou incohérence ne devient pas une conclusion.

## 23. Métriques
- intakes ready/partial/unsupported.
- sessions créées avec provenance complète.
- Tools indisponibles.
- retours au contexte source réussis.
- complétude de provenance et retour au contexte sans perte.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin, plateforme ou intégration native.

## 25. Critères d’acceptation
### 1. Image mémoire partielle
**Given** une Memory Image liée à un Job partial **When** l’analyste ouvre l’intake **Then** les limites sont visibles, les analyses compatibles restent bornées et aucun Tool ne démarre.
### 2. Architecture ou plateforme unsupported
**Given** aucun Tool compatible actif **When** l’intake est consulté **Then** unsupported est visible et les résultats précédents restent accessibles.
### 3. Sans IA
**Given** aucun fournisseur de modèle **When** l’analyste crée une session **Then** les checklists et Tools déterministes restent disponibles.

## 26. Questions ouvertes
OPEN-005, OPEN-008, OPEN-013, OPEN-014, OPEN-015 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, Static/Reverse, futures phases Objects/Permissions/Journeys/Screens/Contracts, handoffs seulement vers 4B.2B.3B/4B.3.
