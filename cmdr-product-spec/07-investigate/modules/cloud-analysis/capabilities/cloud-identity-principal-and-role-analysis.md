---
id: CAP-INV-605
title: Cloud Identity, Principal and Role Analysis
product: investigate
module: cloud-analysis
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
  - REQ-INV-001
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-008
  - OPEN-012
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-605 — Cloud Identity, Principal and Role Analysis

## 1. Définition
Analyser les observations d’identités Cloud, principals, utilisateurs, service identities, rôles, groupes, workloads, applications et fédérations sans assimiler une identité technique à une personne certaine.

## 2. Problème utilisateur
Les aliases, sessions, fédérations et identités de workloads peuvent être attribués à tort à un humain ou fusionnés malgré des sources contradictoires.

## 3. Objectifs
- conserver type d’identité, identifiers déclarés, aliases, assignments, source, timestamp et confiance.
- distinguer humain déclaré, principal technique, workload, application, groupe, rôle, fédération et session candidate.
- gérer contradictions, ambiguïtés, dedup candidates et relations sans attribution automatique.

## 4. Non-objectifs
Aucune résolution d’identité certaine, connexion annuaire, création d’utilisateur, changement de rôle, révocation de session, déchiffrement de token ou attribution humaine automatique.

Aucun choix AWS, Azure, GCP, SaaS, Kubernetes, provider, format, query language ou implementation n’est imposé.

## 5. Propriétaire
Investigate possède uniquement les concepts analytiques Cloud locaux décrits ici. Platform Settings conserve providers, connectors, credentials, secrets, configured organizations/tenants/accounts/subscriptions/projects, ingestion, schemas, parsers, health, retention, storage, policies et configuration administrative. Command conserve runtime Detection, Signal, Alert, Incident, priorité et dispositions. Govern conserve Decision, Approval, Action Request, Response Run, Result et toute mutation de cible. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Shared conserve Entity, Graph, Timeline, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration et Recovery. Aucun owner concurrent.

## 6. Utilisateurs
Principal : **Cloud Security Analyst**. Secondaires : Investigation Lead, SOC Analyst, DFIR Analyst, Evidence Reviewer, Detection Engineer, Cloud Security Analyst, Platform Administrator, Security Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant et environnement CMDR, origin, scope Cloud, période, owner, sources, versions, permissions, restrictions, classification, health, erreurs et return origin sont explicites. Une information inaccessible, absente, stale ou non supportée reste `unknown`, `partial`, `restricted` ou `blocked`; elle n’est jamais inventée. Provider configured ne signifie ni accessible ni couverture complète.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
|Cloud scope and resource observations|CAP-INV-603/604|authorized units and resource context|oui|session version|analysis partial|
|Identity and assignment materials|Platform Settings / authorized source|users, service identities, roles, groups and federations|oui|observed period|no identity assessment|
|Audit/session observations|CAP-INV-607 / source owner|candidate usage and time context|non|event timestamp|activity unknown|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Principal / Role projections|Platform Settings / source owner|identifiers, type, assignments and restrictions|lecture limitée|
|Cloud Resource Observation|Investigate|workload/application/resource context|lecture/lien|
|Entity / relationship projection|Shared|candidate identity links|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Identity Observation|créer, annoter, contester, versionner, superseder|Investigate|identity observation ≠ person|
|Role Observation|créer, comparer, qualifier|Investigate|role label ≠ effective permission|
|Identity relationship candidate|créer/lier, retirer, dispute|Investigate using Shared links|alias ≠ identity proof|

## 11. Fonctionnalités
- analyser users, service identities, roles, groups, workloads, applications et federated identities.
- conserver candidate sessions, aliases, assignments, timestamps et sources.
- qualifier confiance et contradictions.
- comparer identifiers et relations sans fusion automatique.
- préparer contexte pour l’analyse IAM et la timeline.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Identity, Principal and Role Analysis|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-605|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Identity, Principal and Role Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`candidate`, `observed`, `ambiguous`, `corroborated`, `contradicted`, `under-review`, `disputed`, `withdrawn`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Identity Observation|versioned observation|CAP-INV-606/607/614/616|type, source, confidence and ambiguity visible|
|Role Observation|observation|CAP-INV-606|assignment ≠ effective permission|
|Identity relation candidate|relationship|Shared Graph / analyst|candidate, not canonical identity resolution|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-603/604|identity-bearing scope selected|CAP-INV-605|units, resources, identity materials and restrictions|scope/inventory|
|CAP-INV-605|IAM assessment requested|CAP-INV-606|identity/role observations, assignments and contradictions|identity analysis|
|CAP-INV-607|activity refers to principal|CAP-INV-605|event, principal candidate, timestamp and source|audit analysis|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-603/604/606/607/614/616, Settings identities/access, Shared Entity Resolution/Graph, OPEN-008/012/013. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Identity and role read|personal or security-sensitive metadata|0|identity fields masked|possible|viewer/source owner|Settings/Security|Permissions|
|Identity observation/relation create|misattribution|2|restricted aliases masked|OPEN-013|analyst/reviewer|Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- identity ≠ human user.
- principal ≠ person.
- alias match ≠ identity proof.
- role assignment ≠ effective permission certaine.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- identity observations by type and confidence.
- ambiguous or contradicted relations.
- automatic human attribution — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Service identity
**Given** un audit event contient un principal technique avec un nom humain-like  
**When** l’identité est analysée  
**Then** le type technique et l’ambiguïté restent visibles et aucune personne n’est attribuée.

### 2. Fédération contradictoire
**Given** deux sources associent un alias à des principals différents  
**When** la relation est revue  
**Then** les deux candidates et leurs sources restent visibles en état disputed.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** les identités sont analysées  
**Then** tables, exact matching, graphes et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
