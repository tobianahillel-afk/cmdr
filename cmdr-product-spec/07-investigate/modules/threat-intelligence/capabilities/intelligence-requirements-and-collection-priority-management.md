---
id: CAP-INV-502
title: Intelligence Requirements and Collection Priority Management
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
  - OPEN-018
source-of-truth: canonical
---
# CAP-INV-502 — Intelligence Requirements and Collection Priority Management

## 1. Définition
Créer, prioriser, versionner et suivre des Intelligence Requirements et objectifs de collecte fonctionnels, sans exécuter de collecte ni les confondre avec des Hypotheses d’investigation ou de Detection Engineering.

## 2. Problème utilisateur
Sans Requirement explicite, les équipes accumulent des données sans question, critère de satisfaction, scope, consommateur ni contrainte identifiable.

## 3. Objectifs
- définir question, objectif, consommateurs, priorité et scope.
- définir périodes, environnements, catégories attendues, exclusions et contraintes.
- identifier sources candidates, gaps et dépendances.
- lier Case, Hunt ou Detection Engineering puis suspendre, rouvrir, clôturer ou superseder.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Requirements and Collection Priority Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Qualified intake | CAP-INV-501 | need, scope and return origin | oui | current intake version | Requirement draft blocked |
| Consumer questions | Case/Hunt/Detection Engineering/manager | questions and satisfaction criteria | oui | declared review cycle | incomplete |
| Source catalog projections | CAP-INV-504 / Settings | candidate availability and restrictions | non | current snapshot | source gaps explicit |
| Existing requirements | CAP-INV-502 | versions, priority and dependencies | non | latest active version | new requirement |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Threat Intelligence Intake | Investigate | need and scope | lecture/lien |
| Case / Hunt / Detection Project | Investigate | consumer context | lecture/lien |
| Intelligence Source projection | Settings / CAP-INV-504 | availability and restrictions | lecture |
| Requirement version/history | Investigate / Shared Versioning | priority and lifecycle | lecture/comparaison |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Intelligence Requirement | créer, modifier, suspendre, clôturer, rouvrir, superseder | Investigate concept | objective ≠ collection execution |
| Collection Priority | définir et réviser | Investigate concept | priority is contextual, not operational Command priority |
| Source/gap relation | lier/versionner | Investigate | no source configuration |

## 11. Fonctionnalités
- définir question, objectif, consommateurs, priorité et scope.
- définir périodes, environnements, catégories attendues, exclusions et contraintes.
- identifier sources candidates, gaps et dépendances.
- lier Case, Hunt ou Detection Engineering puis suspendre, rouvrir, clôturer ou superseder.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Intelligence Manager | sources et projections de Intelligence Requirements and Collection Priority Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Intelligence Manager | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Intelligence Manager | concept fonctionnel Intelligence Requirements and Collection Priority Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Requirements and Collection Priority Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`draft`, `active`, `paused`, `partially-satisfied`, `satisfied`, `blocked`, `obsolete`, `cancelled`, `superseded`, `archived`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Versioned Intelligence Requirement | requirement concept | CAP-INV-503..505/518 | question, scope, priority and satisfaction criteria |
| Collection gap | gap relation | Settings request / future 4B.3B.2 | no active collection |
| Satisfied/obsolete disposition | business event | origin consumers | evidence and reason retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-501 | intake qualified | CAP-INV-502 | question, scope, constraints and consumers | Intake |
| CAP-INV-502 | requirement active | CAP-INV-503 | priority, expected knowledge and dependencies | Requirement |
| CAP-INV-502 | source needed | CAP-INV-504 / Settings request | purpose, scope, restrictions and gap | Requirement |
| CAP-INV-518 | new unresolved question | CAP-INV-502 | handoff gaps and provenance | Handoff |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-501/503/504/518; Cases/Hunts/Detection Engineering; Settings sources; Shared Versioning/Collaboration; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Requirement read/create/update | priority and scope influence | 0/2 | consumer-sensitive context scoped | OPEN-013 | author/manager | Investigate | Permissions |
| Close/reopen/supersede | loss of active need | 2 | history always visible | OPEN-013 | owner/reviewer | Investigate | Permissions |
| Cross-tenant requirement prepare | tenant leakage | 2 | tenant-specific facts masked | step-up possible | requester/manager | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Intelligence Requirement ≠ Case Hypothesis ≠ Detection Hypothesis.
- Collection objective ≠ collection execution.
- Priority does not override source restrictions, permissions or Govern.
- No provider, feed, protocol or standard is selected.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- requirements by state/priority/consumer.
- time to satisfaction disposition.
- requirements with explicit criteria and gaps.
- collection executions caused locally — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Requirement incomplet
**Given** un intake sans critère de satisfaction  
**When** un Requirement est préparé  
**Then** il reste draft/incomplete, aucune collecte n’est exécutée et les lacunes sont visibles.

### 2. Requirement satisfait partiellement
**Given** certaines questions ont des réponses sourcées et d’autres non  
**When** la revue est conduite  
**Then** l’état peut être partially-satisfied, les questions non résolues restent visibles.

### 3. Sans IA
**Given** aucun modèle  
**When** le Requirement est défini  
**Then** formulaires, matrices de priorité et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
