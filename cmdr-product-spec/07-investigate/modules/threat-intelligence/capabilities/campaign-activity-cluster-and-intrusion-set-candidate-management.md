---
id: CAP-INV-511
title: Campaign, Activity Cluster and Intrusion Set Candidate Management
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
# CAP-INV-511 — Campaign, Activity Cluster and Intrusion Set Candidate Management

## 1. Définition
Créer et gérer des Activity Cluster, Campaign et Intrusion Set Candidates avec fenêtres, scope, Sightings, behaviors, TTP, malware, infrastructure, entity candidates, gaps et contradictions, sans attribution avancée.

## 2. Problème utilisateur
Des événements proches ou partageant des TTP peuvent être fusionnés en campagne, intrusion set ou actor certain.

## 3. Objectifs
- distinguer Activity Cluster, Campaign et Intrusion Set.
- lier Sightings, behaviors, TTP, malware, infrastructure et Threat Entity candidates.
- comparer, annoter, contester, versionner, séparer et fusionner seulement après revue.
- préserver gaps et alternatives en reportant l’attribution à 4B.3B.2.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Campaign, Activity Cluster and Intrusion Set Candidate Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Sightings/observations | CAP-INV-513 | time, environment and source context | oui | time-bounded | candidate unsupported |
| Behavior/TTP mappings | CAP-INV-512 | candidate techniques and alternatives | non | mapping versions | behavior gap |
| Malware/infrastructure/entity candidates | CAP-INV-508..510 | linked candidate knowledge | non | current versions | relations limited |
| Source/confidence/contradictions | CAP-INV-505/515 | support, credibility and gaps | oui | assessment versions | under-review |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sighting | Investigate concept | occurrence and context | lecture/lien |
| TTP Mapping / Behavior | Investigate | candidate behavior evidence | lecture |
| Malware/Infrastructure/Threat Entity candidates | Investigate | candidate associations | lecture |
| Case/Hunt/Incident | Investigate / Command | source and operational context | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Activity Cluster Candidate | créer, séparer, fusion proposal, versionner | Investigate concept | cluster ≠ campaign |
| Campaign Candidate | créer, contester, versionner, superseder | Investigate concept | campaign ≠ actor/intrusion set |
| Intrusion Set Candidate | créer, contester, versionner, superseder | Investigate concept | candidate ≠ attributed actor |

## 11. Fonctionnalités
- distinguer Activity Cluster, Campaign et Intrusion Set.
- lier Sightings, behaviors, TTP, malware, infrastructure et Threat Entity candidates.
- comparer, annoter, contester, versionner, séparer et fusionner seulement après revue.
- préserver gaps et alternatives en reportant l’attribution à 4B.3B.2.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Campaign, Activity Cluster and Intrusion Set Candidate Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Campaign, Activity Cluster and Intrusion Set Candidate Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Campaign, Activity Cluster and Intrusion Set Candidate Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`proposed`, `under-review`, `supported`, `weakly-supported`, `ambiguous`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Cluster/Campaign/Intrusion Set candidates | candidate knowledge | CAP-INV-514/515/518 | scope, time, sources and gaps visible |
| Merge/separation disposition | review event | CAP-INV-516 | no silent merge |
| Advanced-analysis handoff context | candidate package | future 4B.3B.2 | no attribution or Report |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-513 | multiple relevant Sightings | CAP-INV-511 | time, scope, source and candidate links | Sightings |
| CAP-INV-512/508..510 | candidate relation changes | CAP-INV-511 | behavior, TTP, malware, infrastructure and entity refs | source knowledge |
| CAP-INV-511 | duplicate/merge review | CAP-INV-516 | versions, similarities, differences and restrictions | Candidate |
| CAP-INV-511 | foundation package ready | CAP-INV-518 | candidates, support, contradictions, gaps and unresolved attribution | Candidate |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-505/508..510/512..516/518; Case/Hunt/Incident; Shared Graph/Timeline/Comparison; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Campaign/cluster candidate create/update | attribution and operational influence | 2 | sensitive sources/targets masked | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Merge/separate review | knowledge collapse | 2 | provenance/restrictions preserved | step-up possible | reviewer/owner | Investigate | Permissions |
| Future analysis handoff | advanced attribution boundary | 2 | candidate-only package | possible | author/future reviewer | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Activity Cluster ≠ Campaign ≠ Intrusion Set ≠ Threat Actor.
- Co-occurrence/shared TTP ≠ same operation/actor.
- Campaign Candidate ≠ confirmed campaign or attribution.
- Advanced fusion, actor assessment and strategic/tactical production remain 4B.3B.2.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- candidates by category/state/window.
- merge/separation dispositions.
- candidates with contradictions/gaps.
- automatic campaign/actor confirmations — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Campaign candidate
**Given** plusieurs Sightings, TTP similaires et périodes proches sans attribution fiable  
**When** une candidate est créée  
**Then** elle reste distincte d’un actor, sources/contradictions visibles et Activity Cluster reste alternative.

### 2. Fusion contestée
**Given** deux clusters avec restrictions et différences  
**When** une fusion est proposée  
**Then** aucune fusion silencieuse, les versions restent distinctes jusqu’à revue.

### 3. Sans IA
**Given** aucun modèle  
**When** les candidates sont gérées  
**Then** comparateurs, timelines, graphes/tabular views et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
