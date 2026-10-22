---
id: CAP-INV-514
title: Intelligence Relationship Graph and Evidence-Backed Linking
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
# CAP-INV-514 — Intelligence Relationship Graph and Evidence-Backed Linking

## 1. Définition
Créer, comparer, revoir, contester, retirer et superseder des relations Intelligence candidates sourcées, directionnelles, temporelles et explicables, en consommant Shared Graph sans confondre edge, causalité, contrôle ou attribution.

## 2. Problème utilisateur
Un lien visuel ou une co-occurrence peut être interprété comme relation confirmée, causalité ou attribution, et une fusion peut perdre les sources.

## 3. Objectifs
- définir source, destination, type fonctionnel, direction, période et contexte.
- conserver sources, confiance, contradictions et relations alternatives.
- naviguer vers les objets et fournir une alternative tabulaire accessible.
- confirmer seulement par revue et préserver toutes les versions.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Relationship Graph and Evidence-Backed Linking** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Candidate nodes | CAP-INV-507..513 | versioned candidate knowledge | oui | current versions | relation invalid |
| Supporting sources | Materials/Artifacts/Events/Sightings | evidence-bearing references | oui | resolvable versions | unsupported relation |
| Contradictions/alternatives | CAP-INV-515 | opposing or alternative relations | non | current assessments | uncertainty incomplete |
| Graph/linking capability | Shared Graph/Object Linking | generic visualization and navigation | oui | available projection | table-only fallback |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Entity / candidate knowledge | Shared / Investigate | source and destination projections | lecture |
| Artifact / Event / Sighting / Material | source owners | supporting context | lecture/lien |
| Graph / Object Linking | Shared | generic relationships and navigation | lecture |
| Confidence/Contradiction assessments | Investigate | support and alternatives | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Intelligence Relationship | proposer, revoir, contester, retirer, superseder | Investigate concept | candidate relation ≠ causal/confirmed relation |
| Relationship source link | lier/versionner | Shared Linking / source owner | all sources preserved |
| Alternative relationship set | conserver/comparer | Investigate | no silent collapse |

## 11. Fonctionnalités
- définir source, destination, type fonctionnel, direction, période et contexte.
- conserver sources, confiance, contradictions et relations alternatives.
- naviguer vers les objets et fournir une alternative tabulaire accessible.
- confirmer seulement par revue et préserver toutes les versions.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Intelligence Relationship Graph and Evidence-Backed Linking | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Intelligence Relationship Graph and Evidence-Backed Linking | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Relationship Graph and Evidence-Backed Linking | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`proposed`, `under-review`, `supported`, `contradicted`, `ambiguous`, `disputed`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Intelligence Relationship candidate | relation concept | Knowledge Project/CAP-INV-515/518 | direction, sources and uncertainty visible |
| Graph/table projection | Shared projection | analysts/reviewers | same information available non-visually |
| Review/dispute event | business event | relation owners | history and reasons retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-507..513 | association proposed | CAP-INV-514 | source/destination versions, context and source refs | Candidate |
| Shared Graph/Object Linking | navigation requested | CAP-INV-514 | authorized projections and relation filters | Graph |
| CAP-INV-515 | confidence/contradiction changes | CAP-INV-514 | assessment, alternatives and rationale | Assessment |
| CAP-INV-514 | foundation handoff | CAP-INV-518 | relationships, sources, contradictions and unresolved edges | Relationship |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Shared Entity/Graph/Object Linking/Comparison/Inspector; CAP-INV-505/507..518; Artifact/Event sources; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Relationship propose/update/review | causal/attribution misstatement | 2 | restricted nodes/sources masked | OPEN-013 | author/reviewer | Investigate | Permissions |
| Graph navigate | relationship disclosure | 0 | permission-aware edges | possible | viewer/source owner | Shared | Permissions |
| Relation export | source/restriction leakage | 1 | markings and redaction | step-up possible | authorized exporter/reviewer | Shared/Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Graph edge ≠ causal relationship, control or attribution.
- Relationship candidate ≠ confirmed relationship; co-occurrence ≠ control.
- Shared Graph remains owner of generic visualization.
- No physical graph model, cardinality, query language or external exchange.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- relations by state/type/direction.
- relations with sources/period/contradictions.
- graph/table consistency.
- silent merges or unsupported confirmed edges — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Relation contradictoire
**Given** une relation soutenue et une source contraire  
**When** la contradiction est enregistrée  
**Then** deux sources visibles, relation non supprimée automatiquement et confiance réévaluable.

### 2. Graphe indisponible
**Given** Shared Graph indisponible  
**When** la relation est consultée  
**Then** alternative tabulaire conserve sources, direction et contexte.

### 3. Sans IA
**Given** aucun modèle  
**When** les relations sont revues  
**Then** tables, linking, comparaison et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
