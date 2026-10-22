---
id: CAP-INV-516
title: Knowledge Deduplication, Versioning and Supersession
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
# CAP-INV-516 — Knowledge Deduplication, Versioning and Supersession

## 1. Définition
Détecter et comparer des doublons candidats, proposer ou refuser des fusions, préserver les objets distincts, créer des versions et relations de supersession avec aliases, sources, restrictions, Sightings, contradictions et historique complets.

## 2. Problème utilisateur
Une déduplication silencieuse peut fusionner des réalités différentes, perdre les restrictions ou effacer les contradictions et la provenance.

## 3. Objectifs
- comparer valeurs, types, sources, périodes, relations, restrictions et Sightings.
- proposer/refuser une fusion et conserver des candidates distinctes.
- créer versions, diff, aliases et supersession.
- restaurer une version active documentaire sans supprimer l’historique.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Knowledge Deduplication, Versioning and Supersession** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Candidate records | CAP-INV-507..515 | values, types, versions and relations | oui | immutable versions | no comparison |
| Sources/restrictions | CAP-INV-504..506 | markings, licences and provenance | oui | source versions | merge blocked |
| Sightings/time context | CAP-INV-513 | periods, environments and occurrences | non | current history | identity uncertain |
| Contradictions/confidence | CAP-INV-515 | alternatives and rationale | oui | assessment versions | review incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Candidate knowledge records | Investigate | content and versions | lecture/comparaison |
| Version / Comparison mechanisms | Shared / Studio version projection | diff and history | lecture |
| Source/Restriction/Sighting/Relationship | source owners / Investigate | identity and merge constraints | lecture |
| Entity projection | Shared | possible canonical relation | lecture only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Duplicate Candidate Assessment | créer, contester, superseder | Investigate concept | duplicate candidate ≠ same real object |
| Merge Proposal/Disposition | proposer, accepter/rejeter by review | Investigate concept | no silent merge |
| Knowledge Version / Supersession relation | créer, restore context, superseder | Investigate concept using Shared Versioning | superseded ≠ deleted |

## 11. Fonctionnalités
- comparer valeurs, types, sources, périodes, relations, restrictions et Sightings.
- proposer/refuser une fusion et conserver des candidates distinctes.
- créer versions, diff, aliases et supersession.
- restaurer une version active documentaire sans supprimer l’historique.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Knowledge Deduplication, Versioning and Supersession | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Knowledge Deduplication, Versioning and Supersession | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Knowledge Deduplication, Versioning and Supersession | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`comparison-ready`, `duplicate-candidate`, `merge-proposed`, `merge-rejected`, `kept-distinct`, `version-created`, `superseded`, `restored-context`, `disputed`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Duplicate comparison | assessment | reviewers/candidate owners | differences and restrictions explicit |
| Merge/keep-distinct disposition | business event | knowledge project | all source lineage preserved |
| Knowledge version/supersession | version relation | CAP-INV-517/518 | history and aliases retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-507..515 | similarity detected | CAP-INV-516 | records, versions, sources, restrictions and relations | Candidates |
| CAP-INV-516 | merge proposal reviewed | candidate owners | accept/reject/keep-distinct rationale and lineage | Comparison |
| CAP-INV-516 | new version/supersession | CAP-INV-517 | active/stale consumers, history and lifecycle context | Version |
| CAP-INV-516 | handoff ready | CAP-INV-518 | versions, aliases, supersessions and unresolved duplicates | Version |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-504..518; Shared Versioning/Comparison/Entity Resolution/Trace; OPEN-013/014/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Duplicate comparison run/read | sensitive-value matching | 1 | values/restrictions scoped | possible | analyst/reviewer | Shared/Investigate | Permissions |
| Merge proposal/review | identity and provenance loss | 2 | all sources/markings preserved | OPEN-013 | proposer/reviewer | Investigate | Permissions |
| Version/supersession create | consumer impact | 2 | history immutable | step-up possible | owner/reviewer | Investigate/Shared | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Duplicate candidate ≠ same real-world object.
- Merge suggestion ≠ merge; no silent fusion.
- Superseded ≠ deleted; restore context creates a new disposition/version.
- No physical matching algorithm, canonical identifier format or final cardinality.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- duplicate candidates and review dispositions.
- merges rejected/kept distinct.
- versions with complete source/restriction lineage.
- silent merges or provenance loss — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Doublon candidat
**Given** deux Indicator Candidates similaires avec sources/restrictions différentes  
**When** le système propose un doublon  
**Then** aucune fusion silencieuse, différences/markings visibles, analyste peut refuser et provenance complète.

### 2. Supersession
**Given** une knowledge version corrigée  
**When** elle supersede l’ancienne  
**Then** ancienne version reste résoluble et aucune donnée n’est effacée.

### 3. Sans IA
**Given** aucun modèle  
**When** la déduplication est réalisée  
**Then** comparateurs déterministes, diff et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
