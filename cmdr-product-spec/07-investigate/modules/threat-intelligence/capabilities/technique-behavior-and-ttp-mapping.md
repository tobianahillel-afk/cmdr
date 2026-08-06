---
id: CAP-INV-512
title: Technique, Behavior and TTP Mapping
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
# CAP-INV-512 — Technique, Behavior and TTP Mapping

## 1. Définition
Relier observations et behaviors à des techniques ou taxonomies versionnées avec sources, contexte, plateformes, éléments favorables, contradictions et alternatives, sans preuve automatique d’exécution, intention ou attribution.

## 2. Problème utilisateur
Un mapping taxonomique peut être présenté comme preuve certaine d’une technique exécutée ou d’un actor.

## 3. Objectifs
- relier observations à behavior candidates puis taxonomies/techniques.
- voir taxonomie, version, source, contexte, plateformes et phases candidates.
- comparer mappings, annoter, contester, revoir et superseder.
- fournir une alternative structurée aux visualisations sans imposer une taxonomie unique.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Technique, Behavior and TTP Mapping** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Observed behaviors | Analysis Workbench / CAP-INV-509/510/513 | observations and context | oui | source versions | mapping impossible |
| Taxonomy projections | source material / future decision | technique labels and versions | non | declared version | free-form behavior only |
| Source/confidence assessments | CAP-INV-505/515 | support, contradiction and limits | oui | assessment versions | mapping weak |
| Platform/phase context | Artifacts/Events/Knowledge Project | environment and sequence context | non | same observation period | context incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical observation / Artifact / Sighting | Investigate | behavior and provenance | lecture/lien |
| Taxonomy reference | external/shared projection | label/version only | lecture |
| Malware/Infrastructure/Campaign candidates | Investigate | candidate context | lecture |
| Source/Confidence/Contradiction | Investigate | support and alternatives | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Behavior Candidate | créer/annoter/versionner | Investigate concept | behavior ≠ technique certainty |
| TTP Mapping | proposer, revoir, contester, superseder | Investigate concept | mapping ≠ execution/intent/attribution |
| Alternative mapping relation | conserver | Investigate | no single taxonomy truth |

## 11. Fonctionnalités
- relier observations à behavior candidates puis taxonomies/techniques.
- voir taxonomie, version, source, contexte, plateformes et phases candidates.
- comparer mappings, annoter, contester, revoir et superseder.
- fournir une alternative structurée aux visualisations sans imposer une taxonomie unique.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Technique, Behavior and TTP Mapping | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Technique, Behavior and TTP Mapping | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Technique, Behavior and TTP Mapping | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`proposed`, `under-review`, `supported`, `weakly-supported`, `alternative-mapping`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| TTP Mapping | mapping concept | CAP-INV-511/514/518 / Detection handoff | sources, version and alternatives visible |
| Unmapped behavior gap | gap relation | future analysis/taxonomy decision | no forced mapping |
| Mapping dispute/review event | business event | project/reviewer | history preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Technical observations/CAP-INV-509/510/513 | behavior identified | CAP-INV-512 | observation, platform, time, source and limits | source knowledge |
| CAP-INV-512 | mapping proposed | CAP-INV-515 | support, alternatives, contradiction and confidence context | Mapping |
| CAP-INV-512 | mapping relevant to detection | Detection Engineering handoff | behavior, candidate technique, sources and limits | Mapping |
| CAP-INV-512 | foundation handoff ready | CAP-INV-518 | mappings, versions, gaps and unresolved alternatives | Mapping |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Analysis Workbench; CAP-INV-505/509..515/518; taxonomy sources; Detection Engineering; Shared Mapping/Comparison; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Behavior/TTP mapping create/update | misclassification/attribution risk | 2 | restricted sources masked | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Taxonomy reference read | licence/version constraints | 0 | licensed content scoped | possible | viewer/source owner | source owner | Permissions |
| Mapping export/handoff | context loss | 1/2 | sources/limitations mandatory | step-up possible | author/reviewer | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- TTP mapping ≠ proof of technique execution, intent or attribution.
- Same/shared technique ≠ same actor.
- Taxonomy ≠ absolute source of truth; no proprietary taxonomy is imposed.
- No detection rule or operational action is created.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- mappings by state/taxonomy/version.
- alternative/unmapped behaviors.
- mappings with source/contradiction/context.
- automatic technique/actor conclusions — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Mapping ambigu
**Given** une observation compatible avec plusieurs techniques  
**When** un mapping est créé  
**Then** alternatives et contradictions restent visibles, aucune exécution certaine n’est affirmée.

### 2. Taxonomie absente
**Given** aucune taxonomie autorisée  
**When** le behavior est documenté  
**Then** behavior candidate reste utilisable sans mapping forcé.

### 3. Sans IA
**Given** aucun modèle  
**When** le mapping est revu  
**Then** catalogues, matrices, comparateurs et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
