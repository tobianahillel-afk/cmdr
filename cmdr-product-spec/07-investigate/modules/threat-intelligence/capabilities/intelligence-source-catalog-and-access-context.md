---
id: CAP-INV-504
title: Intelligence Source Catalog and Access Context
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
# CAP-INV-504 — Intelligence Source Catalog and Access Context

## 1. Définition
Consommer un catalogue fonctionnel des sources Intelligence et leur contexte d’accès, restrictions, classification, fraîcheur, couverture, rétention, santé et limites, sans administrer les providers, connectors ou secrets.

## 2. Problème utilisateur
Une source cataloguée peut être prise pour accessible, fiable ou librement réutilisable alors que ses droits, restrictions, santé et licence diffèrent.

## 3. Objectifs
- voir owner administratif, origine, méthode d’accès déclarée et provenance.
- voir restrictions, classification, markings, tenant, environnement et usage autorisé.
- voir fraîcheur, couverture, rétention, health et dernière récupération.
- préparer une demande d’accès sans configurer la source.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Intelligence Source Catalog and Access Context** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Configured sources | Platform Settings | provider/feed/integration projections | oui | current admin snapshot | catalog empty |
| Access policy and markings | Settings/Security/source owner | classification, licence and allowed use | oui | current policy | restricted/unknown |
| Health and retrieval context | Settings | health, freshness and last retrieval | non | current/time series | health unknown |
| Project requirement | CAP-INV-502/503 | purpose, tenant, period and expected knowledge | oui | active version | no relevance assessment |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Data Source / Integration / Provider | Platform Settings | identity, health and access projection | lecture |
| Secret Reference | Platform Settings | existence/status only | aucune lecture secret |
| Tenant / Environment | Platform Settings | scope and restrictions | lecture |
| Intelligence Requirement | Investigate | need and permitted purpose | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Source Catalog Entry projection | annoter local usage and limitations | Investigate concept using Settings source | no administrative mutation |
| Source Access Context | créer/versionner assessment | Investigate concept | access context ≠ permission grant |
| Access request context | préparer | Settings/Govern | request ≠ access |

## 11. Fonctionnalités
- voir owner administratif, origine, méthode d’accès déclarée et provenance.
- voir restrictions, classification, markings, tenant, environnement et usage autorisé.
- voir fraîcheur, couverture, rétention, health et dernière récupération.
- préparer une demande d’accès sans configurer la source.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Intelligence Source Catalog and Access Context | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Intelligence Source Catalog and Access Context | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Intelligence Source Catalog and Access Context | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`available-projection`, `restricted`, `access-request-required`, `unavailable`, `degraded`, `stale`, `retired`, `unknown`, `disputed`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Source catalog projection | catalog view | CAP-INV-505/506 | source owner, access and restrictions visible |
| Access request | request context | Settings/Govern | purpose and scope, no secret |
| Unavailable/retired disposition | source status relation | Knowledge Project | no data absence inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-503 | source needed | CAP-INV-504 | project, requirement and permitted purpose | Project |
| Settings | source/config/health changes | CAP-INV-504 | source version, access status and timestamps | Settings |
| CAP-INV-504 | source accessible | CAP-INV-505/506 | source reference, restrictions and freshness | Catalog |
| CAP-INV-504 | access missing | Settings/Govern request | purpose, tenant, classification and requested scope | Catalog |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Platform Settings sources/integrations/providers/secrets/health/retention; Security permissions; CAP-INV-502/503/505/506; OPEN-008/013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Source catalog read | provider and source sensitivity | 0 | provider/credential details masked | possible | viewer/admin | Settings | Permissions |
| Access context assess | misuse of licensed/restricted data | 2 | markings retained | OPEN-013 | analyst/source owner | Investigate/Settings | Permissions |
| Access request prepare | access escalation | 2 | no credential exposure | step-up possible | requester/approver | Settings/Govern | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Source configured ≠ accessible ≠ reliable.
- Catalog entry does not grant access or reveal secrets.
- Investigate does not create providers, feeds, connectors or retention policies.
- No commercial source or external standard is imposed.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- sources by access/health/freshness state.
- catalog entries with restrictions and owner.
- access requests and dispositions.
- secrets or restricted raw data exposed — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Source accessible, reliability unknown
**Given** une source accessible sans historique qualité  
**When** l’analyste consulte le catalogue  
**Then** accessibilité et reliability restent distinctes et aucun score élevé n’est attribué.

### 2. Source retirée
**Given** une source administrative retired  
**When** un projet la référence  
**Then** le retrait et la dernière version restent visibles, aucune donnée n’est inventée.

### 3. Sans IA
**Given** aucun modèle  
**When** une source est sélectionnée  
**Then** catalogue, filtres, politiques et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
