---
id: CAP-INV-508
title: Threat Entity and Identity Candidate Management
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
# CAP-INV-508 — Threat Entity and Identity Candidate Management

## 1. Définition
Gérer des Threat Entity et identity candidates avec noms, aliases, sources, comportements, infrastructures, campagnes candidates, contradictions et confiance, sans attribution certaine ni conversion automatique en Entity canonique.

## 2. Problème utilisateur
Des aliases ou co-occurrences peuvent être fusionnés en identité ou actor certain, masquant les alternatives et contradictions.

## 3. Objectifs
- gérer adversary, actor, organization, persona, account, operator, group et unknown-cluster candidates.
- voir aliases, sources, relations, behaviors, infrastructure et campaign candidates.
- proposer/refuser une fusion, séparer, annoter, contester, versionner et superseder.
- préparer une relation vers Entity sans transférer l’ownership Shared.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Threat Entity and Identity Candidate Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Candidate names/aliases | CAP-INV-506/507/509..511 | names and identity hints | oui | source versions | candidate incomplete |
| Relations and Sightings | CAP-INV-513/514 | context, time and linked knowledge | non | current versions | identity weak |
| Source/confidence assessments | CAP-INV-505/515 | support, contradiction and limits | oui | latest review | under-review |
| Shared Entity projections | Shared Entity Resolution | possible canonical entity links | non | current projection | standalone candidate |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Entity | Shared Capabilities | canonical identity projection | lecture/proposer lien |
| Observable/Indicator candidates | Investigate | identity hints and context | lecture |
| Campaign/Infrastructure/Malware knowledge | Investigate | candidate associations | lecture |
| Source/Sighting/Relationship assessments | Investigate | support and contradictions | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Threat Entity Candidate | créer, annoter, contester, versionner, superseder | Investigate concept | candidate ≠ canonical Entity/attributed actor |
| Merge/separation proposal | créer, accepter/rejeter par revue | Investigate concept | no silent merge |
| Entity-link proposal | préparer | Shared Entity owner | proposal ≠ canonical link |

## 11. Fonctionnalités
- gérer adversary, actor, organization, persona, account, operator, group et unknown-cluster candidates.
- voir aliases, sources, relations, behaviors, infrastructure et campaign candidates.
- proposer/refuser une fusion, séparer, annoter, contester, versionner et superseder.
- préparer une relation vers Entity sans transférer l’ownership Shared.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Threat Entity and Identity Candidate Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Threat Entity and Identity Candidate Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Threat Entity and Identity Candidate Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
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
| Threat Entity Candidate | candidate knowledge | CAP-INV-509..515/518 | aliases, sources and uncertainty visible |
| Merge/separation disposition | review event | CAP-INV-516 / Shared Entity | all candidates and provenance preserved |
| Entity-link proposal | relation candidate | Shared Entity owner | no ownership transfer |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-506/507 | identity hint extracted | CAP-INV-508 | names, aliases, source and context | source candidate |
| CAP-INV-513/514 | new observation/relation | CAP-INV-508 | time, context, supporting/contradicting links | Sighting/Relationship |
| CAP-INV-508 | possible duplicate | CAP-INV-516 | candidate versions, aliases, restrictions and evidence | Candidate |
| CAP-INV-508 | canonical entity relation proposed | Shared Entity review | candidate, sources, confidence and contradictions | Candidate |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Shared Entity/Entity Resolution/Graph; CAP-INV-505..507/509..516/518; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Threat Entity Candidate create/update | reputational/attribution risk | 2 | sensitive identities masked | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Merge/separate review | identity collapse risk | 2 | sources/restrictions preserved | step-up possible | reviewer/owner | Investigate | Permissions |
| Entity-link proposal | canonical identity impact | 2 | proposal only | step-up possible | requester/Shared owner | Shared/Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Threat Entity Candidate ≠ Entity canonique ≠ attributed Threat Actor.
- Aliases, shared TTP or infrastructure do not prove identity.
- Merge suggestion ≠ merge; duplicate candidate ≠ same real-world entity.
- Advanced actor assessment and attribution remain 4B.3B.2.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- entity candidates by category/state.
- merge proposals accepted/rejected.
- candidates with contradictions/alternatives.
- automatic attributions or silent merges — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Entité ambiguë
**Given** plusieurs aliases, sources et relations contradictoires  
**When** une fusion est proposée  
**Then** candidates restent distinctes, différences visibles, proposition rejetable et historique conservé.

### 2. Entity canonique absente
**Given** une candidate sans Entity correspondante  
**When** un lien est examiné  
**Then** aucune Entity n’est créée automatiquement et la lacune reste visible.

### 3. Sans IA
**Given** aucun modèle  
**When** les candidates sont comparées  
**Then** tables, graphes avec alternative tabulaire, diff et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
