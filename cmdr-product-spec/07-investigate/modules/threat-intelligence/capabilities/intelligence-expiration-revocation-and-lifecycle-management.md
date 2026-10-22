---
id: CAP-INV-517
title: Intelligence Expiration, Revocation and Lifecycle Management
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
# CAP-INV-517 — Intelligence Expiration, Revocation and Lifecycle Management

## 1. Définition
Évaluer et gérer stale, expiration, revocation, withdrawal, archive et supersession de connaissances candidates selon dates, conditions, Sightings, sources, confiance, dépendances et consommateurs, sans effacement ni modification de watchlist active.

## 2. Problème utilisateur
Une connaissance ancienne peut rester utilisée sans avertissement, ou être supprimée alors que l’historique et les consommateurs doivent rester traçables.

## 3. Objectifs
- définir date/condition d’expiration et voir last seen/Sightings.
- marquer stale/expired, proposer puis revoir une révocation.
- réactiver par nouvelle version, retirer de l’usage actif, archiver et superseder.
- notifier les consommateurs futurs tout en conservant provenance et historique.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Expiration, Revocation and Lifecycle Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Knowledge version | CAP-INV-507..516 | candidate state, version and restrictions | oui | current version | no lifecycle |
| Last Sighting/history | CAP-INV-513 | first/last seen and occurrence context | non | current timeline | recency unknown |
| Sources/confidence/contradictions | CAP-INV-505/515 | support, limits and disputes | oui | assessment versions | review required |
| Consumers/dependencies | Knowledge Project / Detection handoff / future analysis | usage and impact | oui | current links | revocation blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Candidate knowledge / version | Investigate | state, version and lifecycle context | lecture |
| Sighting / Timeline | Investigate / Shared | recency and context | lecture |
| Consumer links / Detection Content projection | Shared / Detection Engineering | potential impact | lecture/lien |
| Source/Confidence assessments | Investigate | support and contradictions | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Expiration Assessment | créer, revoir, contester, superseder | Investigate concept | expired ≠ false |
| Revocation Assessment | proposer, revoir, révoquer by authorized review | Investigate concept | revoked ≠ erased |
| Lifecycle disposition/version | stale/withdraw/archive/supersede/reactivate via new version | Investigate concept | no deletion/watchlist mutation |

## 11. Fonctionnalités
- définir date/condition d’expiration et voir last seen/Sightings.
- marquer stale/expired, proposer puis revoir une révocation.
- réactiver par nouvelle version, retirer de l’usage actif, archiver et superseder.
- notifier les consommateurs futurs tout en conservant provenance et historique.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Intelligence Manager | sources et projections de Intelligence Expiration, Revocation and Lifecycle Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Intelligence Manager | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Intelligence Manager | concept fonctionnel Intelligence Expiration, Revocation and Lifecycle Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Expiration, Revocation and Lifecycle Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`active-candidate`, `under-review`, `stale`, `expiring`, `expired`, `revocation-proposed`, `revoked`, `superseded`, `withdrawn`, `archived`, `disputed`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Lifecycle assessment | assessment | candidate owners/consumers | reason, source and impact visible |
| Expiration/revocation disposition | business event | CAP-INV-518/future consumers | history and prior Sightings retained |
| Consumer notification context | notification event | future consumers | no active watchlist change |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-513/516 | time/version condition changes | CAP-INV-517 | last seen, version, consumers and restrictions | Sighting/Version |
| CAP-INV-517 | revocation proposed | authorized lifecycle review | sources, impact, alternatives and history | Assessment |
| CAP-INV-517 | new version reactivates knowledge | candidate owner | new version and superseded relation | Lifecycle |
| CAP-INV-517 | handoff ready | CAP-INV-518 | expiration/revocation state, reasons, consumers and provenance | Lifecycle |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-505/507..518; Shared Timeline/Notifications/Versioning/Trace; Detection Engineering consumer projections; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Expiration propose/review | premature removal from active use | 2 | sensitive consumers masked | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Revocation propose/review | consumer and reputation impact | 2 | reason/source scoped | step-up possible | proposer/approver | Investigate/Govern future | Permissions |
| Archive/withdraw/reactivate version | history and consumer impact | 2 | no deletion | OPEN-013 | owner/reviewer | Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Expired ≠ false; revoked ≠ erased; archived ≠ deleted.
- No recent Sighting ≠ benign.
- No watchlist, Detection rule or external consumer is actively modified.
- Detailed dissemination/consumer lifecycle remains 4B.3B.2.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- knowledge by lifecycle state.
- items past expiry review.
- revocations with consumer/source impact.
- deleted history or active watchlist mutation — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Indicator candidate expiré
**Given** aucune observation récente et condition atteinte  
**When** le lifecycle est évalué  
**Then** candidate peut devenir expired, expired ≠ false, historique/Sightings conservés et aucun objet supprimé.

### 2. Révocation contestée
**Given** une revocation proposée avec contradictions  
**When** la revue intervient  
**Then** la proposition reste distincte de revoked et les alternatives restent visibles.

### 3. Sans IA
**Given** aucun modèle  
**When** le lifecycle est géré  
**Then** dates, règles explicables, timelines et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
