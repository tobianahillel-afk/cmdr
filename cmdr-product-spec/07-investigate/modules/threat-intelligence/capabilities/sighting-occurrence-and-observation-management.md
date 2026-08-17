---
id: CAP-INV-513
title: Sighting, Occurrence and Observation Management
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
# CAP-INV-513 — Sighting, Occurrence and Observation Management

## 1. Définition
Enregistrer et revoir des Sightings, occurrences et observations liés à des candidates, événements ou Artifacts avec tenant, environnement, temps, contexte, confiance, contradictions et provenance, sans les confondre avec Signal, Incident ou compromission.

## 2. Problème utilisateur
Une occurrence importée ou inférée peut être traitée comme alerte ou compromission confirmée, tandis que l’absence de Sighting peut être prise pour absence de menace.

## 3. Objectifs
- enregistrer source, contexte, tenant, environnement, Case, timestamp et fenêtre.
- lier Observable/Indicator candidates, événements et Artifacts sources.
- conserver limitations, confiance, contradictions et Sightings multiples.
- comparer, contester, retirer et superseder sans produire d’objet Command.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Sighting, Occurrence and Observation Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **SOC Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Source observation | Telemetry Event / Artifact / Case / material | occurrence and context | oui | source timestamp/window | no Sighting |
| Candidate reference | CAP-INV-507..511 | Observable/Indicator/entity/knowledge candidate | oui | current version | unlinked observation |
| Tenant/environment/time | Settings/source owner | scope and temporal context | oui | event/source time | partial |
| Source/confidence restrictions | CAP-INV-504/505/515 | access, support and limitations | oui | assessment versions | confidence limited |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Telemetry Event / Search result | Shared | source event and query context | lecture/lien |
| Artifact / Case / Hunt | Investigate | source and investigation context | lecture/lien |
| Observable/Indicator/knowledge candidates | Investigate | candidate identity and version | lecture |
| Signal / Incident | Command | possible related operational projection | lecture/lien only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Sighting | créer, annoter, contester, retirer, superseder | Investigate concept | Sighting ≠ Signal/Incident/compromise |
| Occurrence relation | lier/versionner | Investigate / Shared Linking | source and time retained |
| Sighting comparison | produire | Shared Comparison | no same-activity inference |

## 11. Fonctionnalités
- enregistrer source, contexte, tenant, environnement, Case, timestamp et fenêtre.
- lier Observable/Indicator candidates, événements et Artifacts sources.
- conserver limitations, confiance, contradictions et Sightings multiples.
- comparer, contester, retirer et superseder sans produire d’objet Command.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | SOC Analyst | sources et projections de Sighting, Occurrence and Observation Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | SOC Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | SOC Analyst | concept fonctionnel Sighting, Occurrence and Observation Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Sighting, Occurrence and Observation Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`observed`, `imported`, `inferred-candidate`, `partial`, `contradicted`, `disputed`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Sighting | observation concept | CAP-INV-507..517 | source, time, scope and limits visible |
| Occurrence timeline | timeline relations | Knowledge Project/Case/Hunt | no compromise assertion |
| Dispute/withdrawal event | business event | candidate owners | history retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Event Search/Case/Artifact | observation identified | CAP-INV-513 | event/artifact, time, tenant, environment and limits | source workspace |
| CAP-INV-507..511 | candidate selected | CAP-INV-513 | candidate/version and expected relation | Candidate |
| CAP-INV-513 | Sighting recorded | CAP-INV-507/510/511/514/515 | source, context, time, confidence and contradictions | Sighting |
| CAP-INV-513 | case/hunt relevance | Case/Hunt handoff | observation, event source and limitations | Sighting |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Event Search/Search Jobs; Cases/Hunts/Artifacts; Command projections; CAP-INV-504/505/507..515/517/518; Shared Timeline/Linking; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Sighting create/update/review | cross-tenant and false-compromise risk | 2 | source values/context masked | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Source event read/link | raw telemetry sensitivity | 0 | field-level masking | possible | viewer/source owner | Shared/Investigate | Permissions |
| Cross-tenant analysis | data leakage | 1/2 | tenant isolation mandatory | step-up possible | authorized analyst/reviewer | Security | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Sighting ≠ Signal, Incident or confirmed compromise.
- Repeated Sightings ≠ same activity certainty.
- Absence of Sighting ≠ absence of threat.
- No Alert, response action, watchlist or monitoring loop is created.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- Sightings by source/candidate/state.
- Sightings with source event and complete time scope.
- disputed/withdrawn Sightings.
- automatic compromise or Signal creation — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Sighting candidat
**Given** un event lié à un Indicator Candidate mais contexte incomplet  
**When** un Sighting est enregistré  
**Then** event source et limites restent visibles, aucune compromission/Signal n’est affirmé.

### 2. Absence d’observation
**Given** aucun Sighting récent  
**When** le candidate est revu  
**Then** absence de Sighting ne prouve ni bénignité ni absence de menace.

### 3. Sans IA
**Given** aucun modèle  
**When** les Sightings sont gérés  
**Then** Event Search, formulaires, timelines et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
