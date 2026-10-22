---
id: CAP-INV-510
title: Infrastructure and Resource Knowledge Management
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
# CAP-INV-510 — Infrastructure and Resource Knowledge Management

## 1. Définition
Documenter des adresses, domaines, URLs, certificats, services, hosting, comptes publics, repositories, communication resources et clusters d’infrastructure comme connaissances candidates sourcées et temporelles, sans interaction active ni contrôle adversaire confirmé.

## 2. Problème utilisateur
Une IP, un domaine, un certificat ou une co-occurrence d’hébergement peut être présenté comme identité stable ou infrastructure contrôlée par un adversaire.

## 3. Objectifs
- documenter resources, services, hosting et ownership candidates.
- conserver first/last seen, Sightings, changements, restrictions, expiration et révocation.
- relier campaigns, malware et entities avec confiance et contradictions.
- interdire toute requête active ou action de blocage.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Infrastructure and Resource Knowledge Management** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Network/technical observations | Network Forensics / Event Search / Analysis Workbench | addresses, domains, URLs, certificates and services | oui | time-bounded source refs | no infrastructure record |
| Material/source claims | CAP-INV-505/506 | ownership/hosting claims and limits | non | assessment versions | claim unknown |
| Sightings and changes | CAP-INV-513 | first/last seen, environments and context | non | current timeline | temporal context limited |
| Entity/campaign/malware relations | CAP-INV-508/509/511/514 | candidate associations | non | current versions | standalone infrastructure |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Network observations / Telemetry Events | Investigate / Shared | resource values, time and source | lecture/lien |
| Artifact / Certificate / public resource refs | source owners | observed metadata | lecture |
| Threat Entity/Malware/Campaign candidates | Investigate | candidate relationships | lecture |
| Source/Confidence/Lifecycle assessments | Investigate | support, limits and expiry | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Infrastructure Knowledge Record | créer, annoter, contester, versionner, superseder | Investigate concept | observation ≠ adversary control |
| Infrastructure cluster candidate | créer/séparer | Investigate concept | shared hosting ≠ shared actor |
| Ownership/hosting relation candidate | proposer/revoir | Investigate concept | domain/IP/cert ≠ owner identity |

## 11. Fonctionnalités
- documenter resources, services, hosting et ownership candidates.
- conserver first/last seen, Sightings, changements, restrictions, expiration et révocation.
- relier campaigns, malware et entities avec confiance et contradictions.
- interdire toute requête active ou action de blocage.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections de Infrastructure and Resource Knowledge Management | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Threat Intelligence Analyst | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Threat Intelligence Analyst | concept fonctionnel Infrastructure and Resource Knowledge Management | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Infrastructure and Resource Knowledge Management | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`proposed`, `observed`, `under-review`, `supported`, `shared-resource`, `ambiguous-control`, `contradicted`, `stale`, `expired`, `revoked`, `superseded`, `withdrawn`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Infrastructure knowledge | knowledge record | CAP-INV-511..518 | time, context, restrictions and uncertainty |
| Resource relation candidates | relations | CAP-INV-514 | association ≠ control/attribution |
| Expiration/revocation context | lifecycle context | CAP-INV-517 | historical Sightings preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Network Forensics/Event Search | resource observed | CAP-INV-510 | value/type, timestamp, environment, source and limitations | source workspace |
| CAP-INV-506 | infrastructure claim extracted | CAP-INV-510 | claim, source, credibility and restrictions | Material |
| CAP-INV-513/514 | new Sighting/relation | CAP-INV-510 | temporal and association context | Sighting/Relationship |
| CAP-INV-510 | lifecycle or handoff needed | CAP-INV-517/518 | record, first/last seen, consumers and contradictions | Infrastructure |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
Network Forensics CAP-INV-380..397; Event Search; CAP-INV-505/506/508/509/511/513..518; Shared Entity/Graph/Timeline; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Infrastructure knowledge read/create/update | sensitive/customer infrastructure | 0/2 | values/tenants masked by policy | OPEN-013 | analyst/reviewer | Investigate | Permissions |
| Ownership/control relation review | attribution/reputation risk | 2 | source identities scoped | step-up possible | reviewer/owner | Investigate | Permissions |
| Active query/block | production interaction | 3/4 | not available | Govern required | prohibited/local none | Govern/owner | Future |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- IP address ≠ stable identity; domain/certificate ≠ certain owner.
- Infrastructure observation ≠ adversary-controlled infrastructure confirmed.
- Shared hosting ≠ shared actor.
- No active request, scan, scraping, block or Endpoint/network action.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- resources by type/state/first-last seen.
- shared/ambiguous infrastructure records.
- relations with supporting/contradicting sources.
- active interactions or automatic blocks — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Infrastructure partagée
**Given** une adresse hébergeant plusieurs services liée à une campaign candidate  
**When** elle est documentée  
**Then** shared hosting visible, relations candidates et aucun contrôle adversaire/blocage automatique.

### 2. Certificat ambigu
**Given** un certificat partagé par plusieurs domains  
**When** une relation d’identité est proposée  
**Then** certificate ≠ organization identity and alternatives remain visible.

### 3. Sans IA
**Given** aucun modèle  
**When** l’infrastructure est gérée  
**Then** tables, timelines, graphes sourcés et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
