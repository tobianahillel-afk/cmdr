---
id: CAP-INV-507
title: Observable and Indicator Candidate Management
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
  - OPEN-018
source-of-truth: canonical
---
# CAP-INV-507 — Observable and Indicator Candidate Management

## 1. Définition
Créer, comparer, revoir, versionner, retirer, superseder ou révoquer des Observable Candidates et Indicator Candidates avec sources, Sightings, contexte, restrictions, confiance et contradictions, sans confirmation ou déploiement automatique.

## 2. Problème utilisateur
Une valeur extraite ou observée peut être promue en Indicator, règle ou action de blocage sans contexte, revue ni lifecycle.

## 3. Objectifs
- distinguer explicitement Observable et Indicator.
- gérer source, first/last seen, Sightings, contexte, relations, restrictions et expiration.
- proposer des doublons, annoter, contester et confirmer uniquement par revue fonctionnelle.
- préserver la qualification candidate jusqu’à une future décision objet/opérationnelle.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Observable and Indicator Candidate Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Normalized/extracted values | CAP-INV-506 | candidate values and context | oui | material version | no candidate |
| Sightings/observations | CAP-INV-513 / source events | first/last seen and environment | non | time-bounded | observation limited |
| Source assessments | CAP-INV-505 | reliability, credibility and limitations | oui | assessment version | confidence limited |
| Relations/lifecycle | CAP-INV-514/516/517 | links, duplicates, expiry and revocation | non | current versions | lifecycle incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Normalized Knowledge Record / Material | Investigate | value, extraction and provenance | lecture |
| Telemetry Event / Artifact / Sighting | Shared / Investigate | observation context | lecture/lien |
| Detection Content / Signal / Alert | Investigate Detection / Command | consumer or match projection only | lecture |
| Source/Confidence assessments | Investigate | support, contradictions and restrictions | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Observable Candidate | créer, annoter, contester, retirer, superseder | Investigate concept | Observable ≠ Indicator |
| Indicator Candidate | créer, revoir, retirer, superseder, révoquer | Investigate concept | candidate ≠ confirmed/deployed Indicator |
| Candidate relation/alias | lier/versionner | Investigate / Shared Linking | no silent merge |

## 11. Fonctionnalités
- distinguer explicitement Observable et Indicator.
- gérer source, first/last seen, Sightings, contexte, relations, restrictions et expiration.
- proposer des doublons, annoter, contester et confirmer uniquement par revue fonctionnelle.
- préserver la qualification candidate jusqu’à une future décision objet/opérationnelle.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Observable and Indicator Candidate Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Observable and Indicator Candidate Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Observable and Indicator Candidate Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `stale`, `expired`, `revoked`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Observable Candidate | candidate knowledge | CAP-INV-508..518 | source and candidate status explicit |
| Indicator Candidate | candidate knowledge | Detection Engineering handoff / future 4B.3B.2 | no rule/watchlist/action |
| Lifecycle/duplicate review context | assessment context | CAP-INV-516/517 | history and restrictions retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-506 | value extracted | CAP-INV-507 | value, type, context, source and ambiguity | Material |
| CAP-INV-513 | new Sighting | CAP-INV-507 | candidate ref, event/artifact, time and limits | Sighting |
| CAP-INV-507 | Indicator relevance reviewed | Detection Engineering handoff | candidate, sources, Sightings, contradictions and limits | Candidate |
| CAP-INV-507 | duplicate/lifecycle condition | CAP-INV-516/517 | versions, restrictions, first/last seen and consumers | Candidate |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-505/506/513..518; Artifact/Event Search; Detection Engineering; Shared Linking/Versioning; OPEN-013/014/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Candidate read/create/update/review | sensitive values and false qualification | 0/2 | value masked by class/source | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Indicator revoke/withdraw | consumer impact | 2 | history preserved | step-up possible | owner/reviewer | Investigate | Permissions |
| Detection handoff prepare | future operational influence | 2 | no automatic rule/action | OPEN-013 | requester/detection owner | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Observable ≠ Indicator; Candidate ≠ confirmed Indicator.
- Indicator ≠ Detection Content and ≠ block action.
- Indicator match ≠ confirmed malicious activity; absence ≠ threat absence.
- No watchlist, rule, deployment, network/Endpoint block or external sharing.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- candidates by type/state/source.
- candidates with Sightings/expiry/restrictions.
- candidate-to-handoff dispositions.
- automatic confirmations/deployments — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Observable non qualifié
**Given** une adresse avec un seul Sighting sans contexte malveillant suffisant  
**When** l’analyste l’enregistre  
**Then** elle reste Observable Candidate, aucune watchlist/règle n’est créée et la confiance reste limitée.

### 2. Indicator expiré
**Given** un Indicator Candidate sans observation récente et condition atteinte  
**When** le lifecycle est évalué  
**Then** il peut devenir expired, expired ≠ false, historique/Sightings restent accessibles.

### 3. Sans IA
**Given** aucun modèle  
**When** les candidates sont gérées  
**Then** formulaires, recherche, comparateurs et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-014 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
