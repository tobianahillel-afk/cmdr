---
id: CAP-INV-518
title: Threat Intelligence Provenance and Analysis Handoff
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
# CAP-INV-518 — Threat Intelligence Provenance and Analysis Handoff

## 1. Définition
Retracer de bout en bout Requirements, Projects, sources, materials, extractions, Tools/Runs, candidates, knowledge, Sightings, relations, assessments, versions et lifecycle, puis préparer un Analysis Handoff Package non publié vers la future 4B.3B.2 ou vers Case, Hunt, Detection Engineering et Analysis Workbench.

## 2. Problème utilisateur
Un handoff sans lineage, restrictions, contradictions ou éléments non résolus peut être pris pour rapport, attribution ou Indicator opérationnalisé.

## 3. Objectifs
- reconstituer toutes les sources, versions, acteurs, Tools/Runs, erreurs et décisions humaines.
- préparer un package avec question, scope, reliability, credibility, candidates, relations, Sightings, confidence, gaps et limitations.
- préserver return origin et destinations autorisées.
- garantir que le package n’est ni Report, publication, attribution, Indicator déployé ni watchlist.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Threat Intelligence Provenance and Analysis Handoff** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Requirement/Project lineage | CAP-INV-501..503 | question, scope, owners and return origin | oui | versioned refs | handoff incomplete |
| Source/material/extraction lineage | CAP-INV-504..506 / Studio | access, restrictions, Tool Calls and errors | oui | immutable/versioned | provenance partial |
| Knowledge lineage | CAP-INV-507..517 | candidates, relations, assessments, versions and lifecycle | oui | current selected versions | package partial |
| Cross-product sources | Case/Hunt/Incident/Detection/Analysis Workbench | origin and consumer context | selon scope | resolvable links | source gap visible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| All Threat Intelligence concepts | Investigate | complete foundation lineage | lecture/lien |
| Case/Hunt/Finding/Evidence/Artifact | Investigate | origin and technical context | lecture/lien |
| Incident/Detection/Signal/Alert | Command | operational source projection | lecture/lien |
| Tool/Tool Call/Automation Run | Studio | automated contributor provenance | lecture/lien |
| Trace/Activity/Versioning | Shared | generic lineage and history | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Threat Intelligence Provenance Assessment | créer, compléter, contester, superseder | Investigate concept | missing links explicit |
| Analysis Handoff Package | créer, versionner, retirer, superseder | Investigate concept | package ≠ Report/publication/attribution |
| Cross-product handoff relation | préparer/lier | destination owner / Shared Linking | no ownership transfer or active change |

## 11. Fonctionnalités
- reconstituer toutes les sources, versions, acteurs, Tools/Runs, erreurs et décisions humaines.
- préparer un package avec question, scope, reliability, credibility, candidates, relations, Sightings, confidence, gaps et limitations.
- préserver return origin et destinations autorisées.
- garantir que le package n’est ni Report, publication, attribution, Indicator déployé ni watchlist.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Threat Intelligence Provenance and Analysis Handoff | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Threat Intelligence Provenance and Analysis Handoff | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Threat Intelligence Provenance and Analysis Handoff | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`draft`, `collecting`, `partial`, `complete-for-declared-scope`, `reviewed`, `handoff-ready`, `returned`, `withdrawn`, `superseded`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Analysis Handoff Package | versioned package | future 4B.3B.2 | candidate-only, contradictions/gaps visible |
| Case/Hunt handoff | context package | Case/Hunt owners | no compromise/ground truth invented |
| Detection Engineering handoff | candidate context | CAP-INV-401..435 | no Detection Content created automatically |
| Analysis Workbench handoff | technical question package | technical module | no active collection/execution |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-501..517 | foundation lifecycle event | CAP-INV-518 | IDs, versions, sources, restrictions, errors and dispositions | source concept |
| CAP-INV-518 | future analysis requested | future 4B.3B.2 | question, sources, candidates, relations, confidence, contradictions, gaps and provenance | Handoff |
| CAP-INV-518 | investigation action needed | Case/Hunt | observations, source events, limits and return origin | Handoff |
| CAP-INV-518 | detection need identified | Detection Engineering | candidate, Sightings, TTP, sources, contradictions and limitations | Handoff |
| CAP-INV-518 | technical gap identified | Analysis Workbench | Artifact/material refs, question, restrictions and unresolved items | Handoff |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-501..517; Cases/Hunts/Analysis Workbench/Detection Engineering; Command projections; Studio Tools/Runs; Shared Trace/Activity/Linking/Versioning/Export/Reporting; OPEN-013/014/015/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Provenance read/export | cross-product sensitive lineage | 0/1 | restricted sources/material masked | step-up possible | viewer/auditor | owners/Shared | Permissions |
| Analysis Handoff Package create/update | scope/consumer influence | 2 | data minimized and markings retained | OPEN-013 | author/reviewer | Investigate | Permissions |
| Cross-product handoff prepare | context propagation | 2 | destination permission re-evaluated | possible | requester/destination owner | product owners | Permissions |
| External publication/share | external disclosure | 3/4 | not available | Govern mandatory | prohibited now | Govern | 4B.3B.2 |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Analysis Handoff Package ≠ Intelligence Report, publication or attribution.
- Knowledge record/package ≠ external publication.
- No Indicator deployed, watchlist created, rule generated or response action.
- 4B.3B.2 receives a future package only and is not executed.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- handoffs by destination/disposition.
- packages with complete source/restriction/version lineage.
- unresolved contradictions/gaps preserved.
- external publications/automatic rules/watchlists — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Future analysis
**Given** un Project avec candidates et contradictions non résolues  
**When** un handoff package est préparé  
**Then** sources, gaps et contradictions visibles, aucun Report/attribution inventé et destination reste 4B.3B.2.

### 2. Detection handoff
**Given** un Indicator Candidate avec Sightings et aucune règle  
**When** un handoff est préparé  
**Then** candidate reste non déployé, sources/limits transmis, aucun Detection Content automatique et return origin préservé.

### 3. Sans IA
**Given** aucun modèle  
**When** provenance et handoff sont préparés  
**Then** trace links, tables, timelines, diff et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
