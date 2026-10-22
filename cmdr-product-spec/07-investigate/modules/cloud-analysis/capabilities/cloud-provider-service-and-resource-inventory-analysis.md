---
id: CAP-INV-604
title: Cloud Provider, Service and Resource Inventory Analysis
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
# CAP-INV-604 — Cloud Provider, Service and Resource Inventory Analysis

## 1. Définition
Analyser des inventaires Cloud autorisés afin de produire des observations versionnées sur services, ressources candidates, types, régions, ownership administratif, tags, état observé, relations et différences de snapshots.

## 2. Problème utilisateur
Un inventaire partiel ou ancien peut être présenté comme état exhaustif et courant, masquant les sources divergentes, les ressources inconnues et les zones non couvertes.

## 3. Objectifs
- agréger des snapshots et observations sans imposer un provider ou schéma.
- conserver identifiers déclarés, types, régions, tags, owners administratifs, source et timestamp.
- comparer sources et snapshots, exposer inconnus, doublons, stale data et couverture.

## 4. Non-objectifs
Aucune découverte active, scanner CSPM, appel provider, création ou suppression de ressource, validation de conformité ou inventaire temps réel garanti.

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
|Accepted Cloud Scope|CAP-INV-603|units, regions, periods and exclusions|oui|session version|no inventory analysis|
|Inventory snapshots or materials|Platform Settings / Artifact owner|provider-neutral resource observations|oui|observed timestamps|partial|
|Source health and schema projection|Platform Settings|coverage, parser/schema and limitations|oui|current projection|quality unknown|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud scope units|Investigate / Settings projection|authorized hierarchy and coverage|lecture|
|Data Source / Parser / health|Platform Settings|source lineage and quality limits|lecture limitée|
|Entity / Graph projection|Shared|candidate resource relations|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Resource Observation|créer, comparer, contester, versionner, superseder|Investigate|resource observed ≠ resource current|
|Cloud Inventory Snapshot Assessment|créer, réviser|Investigate|inventory ≠ certain current state|
|Unknown or inventory gap|créer/lier|Investigate / destination owner|gap ≠ missing resource proof|

## 11. Fonctionnalités
- cataloguer services, resource candidates, types, identifiers, régions, tags et owners administratifs.
- conserver source, timestamp, snapshot version et état observé.
- comparer plusieurs snapshots et sources.
- identifier unknown types, unsupported resources, duplicates and contradictions.
- produire relations candidates sans créer d’Entity canonique.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Provider, Service and Resource Inventory Analysis|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-604|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Provider, Service and Resource Inventory Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `loading-snapshot`, `partial`, `observed`, `diff-ready`, `under-review`, `stale`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Resource Observation set|versioned observations|CAP-INV-608..614/616|source and observed time retained|
|Inventory Snapshot Assessment|assessment|session / reviewer|coverage and differences visible|
|Inventory gap|gap record|Settings / Collection / CAP-INV-617|no active discovery|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-603|scope accepted|CAP-INV-604|scope units, regions, periods and coverage|scope|
|CAP-INV-604|resource family selected|CAP-INV-608..614|observations, relations, versions and source limits|inventory|
|New snapshot|comparison requested|CAP-INV-604|prior/current snapshots and timestamps|snapshot source|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-603/608..618, Settings sources/parsers/health, Shared Entity/Graph/Versioning, OPEN-008/012/014. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Inventory read/compare|resource and tenant disclosure|0/1|tenant/resource scoped|possible|viewer/source owner|Settings/Investigate|Permissions|
|Observation create/dispute|incorrect inventory conclusion|2|restricted metadata masked|OPEN-013|analyst/reviewer|Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- account inventory ≠ current state certain.
- resource observed ≠ resource still active.
- tag or owner label ≠ canonical ownership.
- duplicate source records ≠ corroboration.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- observations by source, type and freshness.
- unknown or unsupported resource types.
- inventory presented as current without timestamp — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Snapshots divergents
**Given** deux sources décrivent des états différents  
**When** les inventaires sont comparés  
**Then** les deux observations, timestamps, sources et limites restent visibles sans choisir silencieusement.

### 2. Ressource stale
**Given** une ressource n’apparaît que dans un ancien snapshot  
**When** elle est examinée  
**Then** elle reste observed/stale et n’est pas déclarée active ou supprimée.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** l’inventaire est analysé  
**Then** tables, trees, diff, filters et règles déterministes suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
