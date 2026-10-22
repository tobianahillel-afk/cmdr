---
id: CAP-INV-501
title: Threat Intelligence Intake and Preconditions
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
# CAP-INV-501 — Threat Intelligence Intake and Preconditions

## 1. Définition
Qualifier un besoin Threat Intelligence sourcé depuis un Case, Finding, Hypothesis, Hunt, Incident, Detection Engineering ou une analyse technique, définir son scope et préserver son return origin sans confirmer silencieusement une connaissance.

## 2. Problème utilisateur
Un besoin mal cadré peut transformer des candidates techniques en conclusions Intelligence, perdre les restrictions de source ou rompre le retour vers l’investigation d’origine.

## 3. Objectifs
- ouvrir un intake depuis les objets autorisés.
- rendre visibles sources, Evidence, Artifacts, candidates, contradictions et restrictions.
- définir besoin, scope, futurs consommateurs et critères de complétude.
- créer ou reprendre un Knowledge Project sans produire d’attribution.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Threat Intelligence Intake and Preconditions** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Origin context | Case/Finding/Hypothesis/Hunt/Incident/CAP-INV-435/technical handoff | need, sources, candidates and return origin | oui | current immutable refs | intake incomplete |
| Authorized evidence and artifacts | Investigate | Evidence, Findings, Artifacts and analysis outputs | non | linked versions | sources missing |
| Restrictions and permissions | source owner / Settings / Security | classification, markings, tenant and access | oui | current policy snapshot | permission-blocked |
| Existing Knowledge Project | CAP-INV-503 | project context and linked requirements | non | current version | create new draft project |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Hypothesis / Hunt / Artifact / Evidence | Investigate | need, sources, contradictions and return origin | lecture/lien |
| Incident / Detection / Signal / Alert | Command | operational source projection | lecture seulement |
| Continuous Improvement Package | Investigate Detection Engineering | candidate Intelligence context | lecture/lien |
| Tool/Run provenance | Studio / Shared | producer, parameters, errors and lineage | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Threat Intelligence Intake | créer, compléter, bloquer, superseder | Investigate concept | aucune connaissance confirmée silencieusement |
| Knowledge Project relation | lier ou reprendre | Investigate concept | Project distinct du Case et de l’Automation Run |
| Return-origin relation | conserver | Shared linking / source owner | ownership et permissions inchangés |

## 11. Fonctionnalités
- ouvrir un intake depuis les objets autorisés.
- rendre visibles sources, Evidence, Artifacts, candidates, contradictions et restrictions.
- définir besoin, scope, futurs consommateurs et critères de complétude.
- créer ou reprendre un Knowledge Project sans produire d’attribution.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Threat Intelligence Intake and Preconditions | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Threat Intelligence Intake and Preconditions | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Threat Intelligence Intake and Preconditions | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `blocked`, `insufficient-context`, `source-required`, `permission-blocked`, `out-of-scope`, `superseded`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Qualified intake | Threat Intelligence Intake | CAP-INV-502/503 | need, scope, missing context and restrictions visible |
| Out-of-scope disposition | business event | origin owner | reason and return origin retained |
| Source/access gap | request context | CAP-INV-504 / Settings | no administrative source change |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case/Finding/Hunt/Incident | intelligence need identified | CAP-INV-501 | question, sources, candidates, restrictions and return origin | source workspace |
| Analysis Workbench / CAP-INV-435 | candidate context available | CAP-INV-501 | technical observations, limits and provenance | technical source |
| CAP-INV-501 | intake ready | CAP-INV-502/503 | need, scope, consumers and missing information | Intake |
| CAP-INV-501 | source/access missing | CAP-INV-504 / Settings request | source need, restriction and purpose | Intake |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-001..114, CAP-INV-313/328/346/362/379/397, CAP-INV-435; Command projections; Shared Linking/Trace; OPEN-013/014/015/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Intake read/create/update | scope and sensitive origin | 0/2 | restricted sources masked | OPEN-013 | requester/reviewer | Investigate | Permissions |
| Cross-product source read | tenant and investigation sensitivity | 0 | projection scoped | possible | viewer/source owner | source owner | Permissions |
| Knowledge Project relation create | wrong-context propagation | 2 | no raw-content expansion | OPEN-013 | author/owner | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Intelligence Requirement ≠ Case Hypothesis or Detection Hypothesis.
- Intelligence Material and candidates are not Evidence or facts.
- No collection execution, attribution, Indicator confirmation, watchlist or external sharing.
- Missing permission must not reveal protected content.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- intakes by origin and disposition.
- intakes blocked by missing context/source/permission.
- return-origin links resolved.
- automatic confirmations or attributions — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Intake incomplet
**Given** un Finding sans question précise et plusieurs candidates  
**When** l’analyste ouvre l’intake  
**Then** l’état reste incomplete, les lacunes sont visibles, aucun Requirement/Indicator/attribution n’est créé silencieusement.

### 2. Permission refusée
**Given** une source restreinte et un utilisateur sans lecture brute  
**When** l’intake est consulté  
**Then** l’existence et les métadonnées autorisées peuvent rester visibles, le contenu est masqué et le refus audité.

### 3. Sans IA
**Given** aucun fournisseur de modèle  
**When** l’intake est qualifié  
**Then** formulaires, liens, checklists et revue humaine couvrent le workflow essentiel.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
