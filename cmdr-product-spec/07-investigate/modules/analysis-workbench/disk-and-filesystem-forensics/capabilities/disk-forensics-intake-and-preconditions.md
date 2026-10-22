---
id: CAP-INV-363
title: Disk Forensics Intake and Preconditions
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
# CAP-INV-363 — Disk Forensics Intake and Preconditions

## 1. Définition
Ouvrir une Disk Image ou un Artifact équivalent dans un contexte Disk Forensics explicite, vérifier les préconditions fonctionnelles et préparer ou reprendre une session sans lancer automatiquement d’analyse.

## 2. Problème utilisateur
Une image peut être partielle, corrompue, chiffrée, restreinte ou incompatible. Sans intake, l’analyste peut interpréter une source incomplète comme fiable ou démarrer un Tool inadapté.

## 3. Objectifs
- ouvrir depuis Case, Collection Job ou Artifact Management;
- montrer source, Endpoint, acquisition, custody, taille, intégrité, complétude et restrictions;
- montrer partitions, volumes, filesystems et Tools candidats, résultats antérieurs et permissions;
- définir objectif et scope, créer ou reprendre une Disk Forensics Session et préserver le return origin.

## 4. Non-objectifs
Ne pas acquérir la source, monter ou modifier un live filesystem, choisir moteur/format/outil tiers, fournir commande, offset, algorithme, API, protocole, code, bypass ou capability Network Forensics complète.

## 5. Propriétaire
Investigate / Analysis Workbench possède l’intake, l’objectif, le scope et les relations analytiques. Endpoint Agent contribue à l’acquisition; Settings administre Fleet/Policies/storage/retention; Studio possède Tools/Runs; Govern l’autorité réelle; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : DFIR Analyst. Secondaires : Investigation Lead, Evidence Reviewer, Artifact Analyst et Audit Analyst.

## 7. Conditions d’entrée
Tenant, environnement et Case conservés; Disk Image/Artifact lisible; source et acquisition connus ou explicitement incomplets; permissions et restrictions évaluées; original immuable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case, objectif, Hypothesis | Investigate | contexte | oui | version courante | draft |
| Disk Image / Artifact | CAP-INV-105 ou Collection | source/version | oui | version sélectionnée | incomplete |
| Endpoint/source et acquisition context | Endpoint Agent / CAP-INV-203/213/214 | provenance/custody | oui | timestamps visibles | partial/disputed |
| Integrity/completeness/restrictions | source et policies | préconditions | oui | dernière revue | review-required |
| Structure candidates et Tools | Tool Calls / Studio | compatibilité | non | résultats attribués | filesystem-required/tool-unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif et return origin | lire/lier |
| Disk Image / Artifact | Investigate | version, source, limites | lire |
| Collection Request / Job | Investigate | acquisition et erreurs | lire |
| Endpoint / Endpoint Agent | Endpoint Agent | source/capability déclarée | projection read-only |
| Tool / Tool Call / Automation Run | Studio | compatibilité et résultats antérieurs | lire/sélectionner |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Disk Forensics Intake | créer/mettre à jour/supersede | Investigate concept | aucune analyse automatique |
| Disk Forensics Session relation | préparer/reprendre | Investigate concept | distincte de Case et Collection Job |
| Trace/Activity event | émettre | Shared | ouverture, sélection et blocage attribués |

## 11. Fonctionnalités
Afficher l’Endpoint/source, acquisition, custody, taille, integrity, completeness, restrictions, structures candidates, résultats antérieurs, Tools compatibles/indisponibles, permissions et limitations; définir objectif/scope; créer ou reprendre une session. Une image partial/corrupted/encrypted/restricted/unsupported reste visible comme telle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter intake | Analyst | Disk Image context | 0 | read | limites visibles | non |
| Définir objectif/scope | Analyst | Intake | 2 | modifiable | version attribuée | OPEN-013 |
| Créer/reprendre session | Analyst | Disk Session | 2 | image exploitable ou limites acceptées | session liée | OPEN-013 |
| Lancer Tool autorisé | Analyst | Tool Call request | 1 | sélection explicite | job attribué | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier préconditions | oui | oui | oui | explication | checklist/validateurs |
| proposer structure/Tool | oui | oui | oui | suggestion | catalogue/profils déterministes |
| résumer limites | oui | oui | oui | oui | tableau sourcé |
| démarrer analyse | explicite | oui | workflow | non silencieuse | action manuelle |

Attribution complète; aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `partial-image`, `corrupted`, `encrypted`, `restricted`, `unsupported`, `filesystem-required`, `tool-unavailable`, `policy-blocked`.

## 15. États d’interface
Loading conserve Case/image; Empty explique la source absente; Partial détaille les limites; Error garde les données valides; Offline limite les mutations; Permission denied masque contenu et chemins; Stale distingue source persistante et état actuel.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Intake assessment | concept | CAP-INV-364/365 | état et limites explicites |
| Session context | relation | Disk Workbench | objectif, scope, owner, return origin |
| Block/limitation event | Trace event | Case/Notifications | aucune précondition inventée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case/Artifact/Collection Job | ouvrir | CAP-INV-363 | tenant, Case, source, image, objective, Hypothesis | source |
| Intake | revoir intégrité | CAP-INV-365 | image, acquisition, custody, restrictions | Intake |
| Intake ready/limited | créer/reprendre | CAP-INV-364 | image, scope, owner, limitations | Intake |

## 18. Dépendances
CAP-INV-105/203/205/212/213/214, Disk Image, Studio Tools, Settings projections, Shared Trace/Jobs/Inspector, OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate possède intake et relations analytiques; Disk Image/Artifact restent sources; Endpoint Agent reste source de l’acquisition déclarée; Studio, Settings, Shared et Govern conservent leurs objets.

## 20. Provenance et audit
Tracer Case, image/version, source/Endpoint, acquisition/custody, actor, objective/scope, permissions, Tools/versions, résultats antérieurs, limitations, décision et correlation IDs.

## 21. Permissions fonctionnelles
Disk Image read/raw/restricted read, intake create/update, session create/resume, Tool compatibility read et automated analysis request. Step-up/SoD/matrice atomique sont reportés.

## 22. Limites et erreurs
Source absente/superseded, image partial/corrupted/encrypted, structure unsupported/ambiguous, Tool unavailable, policy/permission denied, tenant mismatch. Les résultats valides restent visibles et la source n’est jamais modifiée.

## 23. Métriques
Intakes par état, images partielles/corrompues/restreintes, Tools indisponibles, sessions reprises et limitations correctement propagées.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, format, Tool tiers ou protocole.

## 25. Critères d’acceptation
**Given** une image issue d’un Job `partial` avec plages manquantes **When** l’intake est ouvert **Then** `partial-image`, limites, provenance et analyses compatibles sont visibles sans complétude fictive.

**Given** une image chiffrée sans moyen autorisé **When** elle est ouverte **Then** elle reste visible comme inaccessible et aucune attaque n’est lancée.

**Given** aucun modèle IA **When** l’intake est utilisé **Then** formulaires, validateurs, profils et Tools déterministes couvrent le workflow.

## 26. Questions ouvertes
OPEN-005/008/013/014/015 restent ouvertes; moteurs, plateformes, objets et permissions finales sont reportés.

## 27. Consommateurs documentaires
INV-DSK-001, Case Workspace, Artifact Management, CAP-INV-364..379, phases Objects/Permissions/Journeys/Screens/Technique et future Network Forensics uniquement comme handoff.
