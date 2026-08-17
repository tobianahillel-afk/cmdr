---
id: CAP-INV-503
title: Threat Intelligence Workspace and Knowledge Project Management
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-PROD-055
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-018
source-of-truth: canonical
---
# CAP-INV-503 — Threat Intelligence Workspace and Knowledge Project Management

## 1. Définition
Fournir un contexte durable de travail pour relier Requirements, sources, matériaux, candidates, relations, Sightings, assessments, contradictions, versions et handoffs, distinct d’un Case et d’un Automation Run.

## 2. Problème utilisateur
Des connaissances réparties entre documents, runs et investigations perdent leur contexte, leur owner et leur historique.

## 3. Objectifs
- créer et gouverner un Knowledge Project.
- lier Requirements, Cases, Hunts, Findings, Incidents et Detection Engineering.
- organiser sources, matériaux, candidates, relations, Sightings et assessments.
- suspendre, reprendre, clôturer, rouvrir, archiver et superseder sans supprimer l’historique.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Threat Intelligence Workspace and Knowledge Project Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Intelligence Requirements | CAP-INV-502 | questions, priorities and criteria | oui | active versions | project incomplete |
| Origin relations | CAP-INV-501 | Case/Hunt/Incident/Detection and return origin | oui | resolvable links | context partial |
| Sources and materials | CAP-INV-504/506 | references, restrictions and provenance | non | current versions | empty project |
| Contributors and permissions | Settings/Security | owner, roles and tenant scope | oui | current policy | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Requirement / Intake | Investigate | scope and objectives | lecture/lien |
| Case/Hunt/Finding/Incident | Investigate / Command | origin and consumer context | lecture/lien |
| Source/Material/Candidates | source owners / Investigate | knowledge context | lecture |
| Automation Run / Tool Call | Studio | contributor provenance | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Knowledge Project | créer, modifier, pause, reprendre, archiver, superseder | Investigate concept | Project ≠ Case ≠ Automation Run |
| Project membership/relation | lier, annoter, retirer | Investigate / Shared Linking | source ownership retained |
| Project activity | émettre | Shared Activity | no duplicate Activity mechanism |

## 11. Fonctionnalités
- créer et gouverner un Knowledge Project.
- lier Requirements, Cases, Hunts, Findings, Incidents et Detection Engineering.
- organiser sources, matériaux, candidates, relations, Sightings et assessments.
- suspendre, reprendre, clôturer, rouvrir, archiver et superseder sans supprimer l’historique.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Threat Intelligence Workspace and Knowledge Project Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Threat Intelligence Workspace and Knowledge Project Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Threat Intelligence Workspace and Knowledge Project Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed-foundation`, `failed`, `archived`, `superseded`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Knowledge Project context | project concept | CAP-INV-504..518 | owners, scope, versions and restrictions retained |
| Project status event | business event | origin/assignees | partial, blocked and archived explicit |
| Foundation handoff | project relation | CAP-INV-518 | complete declared scope, gaps and provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-502 | requirement activated | CAP-INV-503 | question, priority, scope and consumers | Requirement |
| CAP-INV-503 | source work needed | CAP-INV-504/506 | project, requirement and restrictions | Project |
| CAP-INV-504..517 | knowledge updated | CAP-INV-503 | candidate/version/relation/assessment references | source capability |
| CAP-INV-503 | foundation complete | CAP-INV-518 | scope, sources, candidates, contradictions and gaps | Project |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-501/502/504..518; Shared Linking/Activity/Versioning/Collaboration/Search; Studio Runs; OPEN-013/014/015/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Knowledge Project read/create/update | cross-source sensitive context | 0/2 | restricted items masked | OPEN-013 | owner/contributor | Investigate | Permissions |
| Archive/reopen/supersede | historical continuity | 2 | history immutable | OPEN-013 | owner/reviewer | Investigate | Permissions |
| Contributor assignment | access expansion | 2 | no implicit source permission | step-up possible | owner/admin | Investigate/Settings | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Knowledge Project ≠ Case and ≠ Automation Run.
- Project membership does not grant source permissions.
- Completed-foundation is not an Intelligence Report or publication.
- No advanced fusion, attribution, dissemination or watchlist operation.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- projects by state and owner.
- projects with unresolved links/restrictions.
- foundation-complete projects with handoff.
- silent permission expansion — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Projet partiel
**Given** un Requirement actif et une source indisponible  
**When** le projet est examiné  
**Then** il reste partial/blocked, la source manquante et les éléments valides restent visibles.

### 2. Projet archivé
**Given** un projet foundation-complete  
**When** l’owner l’archive  
**Then** les versions, sources, contradictions et handoffs restent accessibles.

### 3. Sans IA
**Given** aucun modèle  
**When** le projet est géré  
**Then** catalogues, tables, graphes avec alternative tabulaire et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
