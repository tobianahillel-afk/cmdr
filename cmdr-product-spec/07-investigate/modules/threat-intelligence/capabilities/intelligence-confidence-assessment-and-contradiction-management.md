---
id: CAP-INV-515
title: Intelligence Confidence, Assessment and Contradiction Management
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
# CAP-INV-515 — Intelligence Confidence, Assessment and Contradiction Management

## 1. Définition
Exprimer et réviser une confiance analytique avec justification, hypothèses, informations manquantes, reliability, credibility, contradictions, alternatives et évolution temporelle, sans l’assimiler à une probabilité calibrée ou à un fait.

## 2. Problème utilisateur
Un score unique ou le nombre de sources peut masquer leurs dépendances, contradictions et limitations et être présenté comme vérité.

## 3. Objectifs
- créer des assessments sourcés, justifiés et temporels.
- distinguer confidence, reliability, credibility et corroboration.
- voir hypothèses, missing information, contradictions et alternatives.
- comparer, contester, réviser, retirer et superseder en préservant l’historique.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Confidence, Assessment and Contradiction Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Intelligence Reviewer**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Candidate knowledge | CAP-INV-507..514 | candidate version and claims | oui | current version | no assessment |
| Reliability/credibility | CAP-INV-505 | source- and claim-specific assessments | oui | assessment versions | confidence limited |
| Supporting/contradicting sources | Materials/Sightings/Relationships | evidence and alternatives | oui | resolvable periods | inconclusive |
| Assessment method/reviewer | human process | justification, assumptions and scope | oui | assessment timestamp | invalid assessment |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Candidate knowledge / Relationship | Investigate | claim and version | lecture |
| Source Reliability / Information Credibility | Investigate | distinct assessment inputs | lecture |
| Material / Sighting / Artifact / Event | source owners | support and contradiction | lecture/lien |
| Prior Confidence Assessment | Investigate / Shared Versioning | history and changes | lecture/comparaison |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Confidence Assessment | créer, réviser, contester, retirer, superseder | Investigate concept | confidence ≠ calibrated probability/fact |
| Contradiction Record | créer, résoudre contextually, superseder | Investigate concept | contradiction ≠ automatic invalidation |
| Assessment alternative | conserver/versionner | Investigate | multiple interpretations allowed |

## 11. Fonctionnalités
- créer des assessments sourcés, justifiés et temporels.
- distinguer confidence, reliability, credibility et corroboration.
- voir hypothèses, missing information, contradictions et alternatives.
- comparer, contester, réviser, retirer et superseder en préservant l’historique.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Intelligence Reviewer | sources et projections de Intelligence Confidence, Assessment and Contradiction Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Intelligence Reviewer | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Intelligence Reviewer | concept fonctionnel Intelligence Confidence, Assessment and Contradiction Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Confidence, Assessment and Contradiction Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`not-assessed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Confidence Assessment | assessment | CAP-INV-507..518 | rationale, assumptions and limits visible |
| Contradiction Record | knowledge record | review/project/handoff | both sides and sources retained |
| Assessment change event | business event | consumers | before/after and reviewer visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-505/507..514 | assessment needed | CAP-INV-515 | candidate, sources, reliability, credibility and alternatives | source capability |
| CAP-INV-515 | confidence revised | candidate owner | assessment version, rationale and contradictions | Assessment |
| CAP-INV-515 | contradiction affects lifecycle | CAP-INV-516/517 | candidate/version, consumers and unresolved issues | Assessment |
| CAP-INV-515 | handoff ready | CAP-INV-518 | confidence, rationale, sources, contradictions and gaps | Assessment |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-505/507..518; Shared Versioning/Comparison/Trace; human review; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Confidence assessment create/update | decision/reputation influence | 2 | sensitive sources masked | OPEN-013 | assessor/reviewer | Investigate | Permissions |
| Contradiction create/update | source conflict exposure | 2 | source identity scoped | possible | analyst/reviewer | Investigate | Permissions |
| Assessment export | misuse outside context | 1 | limitations/markings mandatory | step-up possible | reviewer/exporter | Shared/Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Confidence ≠ mathematically calibrated probability, source reliability or absolute fact.
- High confidence ≠ Approval.
- Contradiction ≠ automatic invalidation; number of sources ≠ independence.
- No opaque score or model output is presented as truth.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- assessments with rationale/assumptions/missing info.
- contradictions unresolved/resolved by supersession.
- confidence changes by version.
- opaque automatic truth scores — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Contradiction
**Given** une relation soutenue et contredite  
**When** l’assessment est révisé  
**Then** sources restent visibles, confiance réévaluable, relation non supprimée et aucun score opaque ne remplace la revue.

### 2. Sources dépendantes
**Given** plusieurs sources recopiant la même origine  
**When** la confiance est évaluée  
**Then** nombre de sources n’est pas traité comme indépendance.

### 3. Sans IA
**Given** aucun modèle  
**When** la confiance est exprimée  
**Then** matrices explicables, justification et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
