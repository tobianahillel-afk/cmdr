---
id: CAP-INV-506
title: Intelligence Material Intake, Parsing and Normalization
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
# CAP-INV-506 — Intelligence Material Intake, Parsing and Normalization

## 1. Définition
Référencer ou importer un matériau autorisé, préserver l’original, effectuer une extraction bornée et produire une représentation normalisée fonctionnelle avec ambiguïtés, restrictions et provenance, sans le transformer automatiquement en Evidence.

## 2. Problème utilisateur
Une extraction peut perdre l’original, les markings et les ambiguïtés ou être présentée comme connaissance confirmée.

## 3. Objectifs
- référencer/importer un Artifact autorisé et préserver l’original.
- voir type déclaré/détecté, langue, timestamps, source et restrictions.
- extraire de façon bornée et conserver éléments non interprétés, erreurs et ambiguïtés.
- relier des candidates et relancer avec un Tool autorisé sans format externe imposé.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Material Intake, Parsing and Normalization** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Authorized material | Artifact/source reference | raw or referenced Intelligence Material | oui | immutable/versioned | intake blocked |
| Source/access context | CAP-INV-504 | restrictions, markings and permitted use | oui | current policy snapshot | restricted |
| Requirement/project context | CAP-INV-502/503 | purpose, scope and expected knowledge | oui | active versions | unscoped |
| Tool/Parser capability | Studio / Settings parser projection | declared extraction capability and version | non | selected run version | manual extraction only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Artifact / Intelligence Material reference | Investigate/source owner | original, version and custody | lecture |
| Parser / Data Source projection | Settings | declared type/capability/health | lecture |
| Tool / Tool Call / Automation Run | Studio | extraction producer and parameters | lecture/lien |
| Requirement / Knowledge Project | Investigate | purpose, restrictions and return origin | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Intelligence Material Reference | créer, retirer, superseder | Investigate concept | Material ≠ Evidence/Report |
| Normalized Knowledge Record | créer/versionner/retirer | Investigate concept | functional representation, no physical schema |
| Extraction result/relation | créer, annoter, contest | Investigate/Shared provenance | automated extraction ≠ human validation |

## 11. Fonctionnalités
- référencer/importer un Artifact autorisé et préserver l’original.
- voir type déclaré/détecté, langue, timestamps, source et restrictions.
- extraire de façon bornée et conserver éléments non interprétés, erreurs et ambiguïtés.
- relier des candidates et relancer avec un Tool autorisé sans format externe imposé.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Intelligence Material Intake, Parsing and Normalization | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Intelligence Material Intake, Parsing and Normalization | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Material Intake, Parsing and Normalization | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`referenced`, `intake-ready`, `parsing`, `normalized-partial`, `normalized`, `unsupported`, `restricted`, `failed`, `withdrawn`, `superseded`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Normalized material | Normalized Knowledge Record | CAP-INV-507..515 | original, extraction, ambiguity and restrictions linked |
| Uninterpreted/error set | partial result | analyst/Tool retry | valid output retained |
| Candidate extraction package | candidate set | CAP-INV-507..512 | candidate status explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-503/504 | material authorized | CAP-INV-506 | project, source, restrictions and purpose | Project/Catalog |
| Artifact/Source | material referenced/imported | CAP-INV-506 | original ref, version, markings and timestamps | source owner |
| CAP-INV-506 | bounded extraction complete | CAP-INV-507..512 | values, context, confidence limits and provenance | Material |
| CAP-INV-506 | retry requested | Studio Tool/Tool Call | material ref, parameters, restrictions and prior errors | Material |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Artifact Management; CAP-INV-502..505/507..515; Settings parsers/sources; Studio Tools/Runs; Shared Preview/Trace; OPEN-013/014/015/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Material read/import-reference | restricted/licensed content | 0/2 | masked preview before raw read | OPEN-014 | reader/source owner | Investigate/source owner | Permissions |
| Extraction run/read | content disclosure and model transfer | 1 | parameters/output scoped | step-up possible | requester/reviewer | Studio/Investigate | Permissions |
| Normalized record create/update | semantic distortion | 2 | original always linked | OPEN-013 | analyst/reviewer | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Intelligence Material ≠ Evidence and ≠ canonical Report.
- Raw material ≠ normalized knowledge; extracted value ≠ confirmed Observable.
- No exchange format, physical schema, ingestion architecture or parser implementation.
- No material is sent to a model without explicit permission.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- materials by intake/parse disposition.
- extractions with original/restrictions/provenance.
- uninterpreted and ambiguous elements.
- marking loss or unapproved model transfer — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Donnée restreinte
**Given** un matériau restreint et un utilisateur sans raw read  
**When** le projet est consulté  
**Then** existence/metadata autorisées restent visibles, contenu masqué, aucune copie/extraction/model transfer.

### 2. Parsing partiel
**Given** un matériau avec sections non interprétées  
**When** l’extraction termine  
**Then** résultats valides, erreurs et ambiguïtés restent séparés et aucune candidate n’est confirmée.

### 3. Sans IA
**Given** aucun modèle  
**When** le matériau est normalisé  
**Then** parsers déterministes, formulaires et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-015 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
