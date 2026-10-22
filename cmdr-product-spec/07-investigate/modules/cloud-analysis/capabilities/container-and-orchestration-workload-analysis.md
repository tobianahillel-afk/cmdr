---
id: CAP-INV-610
title: Container and Orchestration Workload Analysis
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
# CAP-INV-610 — Container and Orchestration Workload Analysis

## 1. Définition
Analyser des observations Cloud sur clusters, namespaces/scopes, workloads, pods/tasks equivalents, images, identities, configurations, network, storage, secret references, audit events and available runtime observations.

## 2. Problème utilisateur
Des objets d’orchestration provider-specific peuvent être traités comme objets CMDR canoniques et des métadonnées control-plane comme forensics complète du runtime container.

## 3. Objectifs
- représenter cluster candidates, scopes, workloads, task/pod equivalents, images and identities sans imposer Kubernetes.
- lier configuration, network, storage, secret references and audit events.
- qualifier runtime observations disponibles et limitations, puis préparer des handoffs spécialisés.

## 4. Non-objectifs
Aucune commande kubectl/provider, exec dans un container, image pull, secret reveal, runtime acquisition, container isolation ou forensics complète du node/container.

Aucun choix AWS, Azure, GCP, SaaS, Kubernetes, provider, format, query language ou implementation n’est imposé.

## 5. Propriétaire
Investigate possède uniquement les concepts analytiques Cloud locaux décrits ici. Platform Settings conserve providers, connectors, credentials, secrets, configured organizations/tenants/accounts/subscriptions/projects, ingestion, schemas, parsers, health, retention, storage, policies et configuration administrative. Command conserve runtime Detection, Signal, Alert, Incident, priorité et dispositions. Govern conserve Decision, Approval, Action Request, Response Run, Result et toute mutation de cible. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Shared conserve Entity, Graph, Timeline, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration et Recovery. Aucun owner concurrent.

## 6. Utilisateurs
Principal : **DFIR Analyst**. Secondaires : Investigation Lead, SOC Analyst, DFIR Analyst, Evidence Reviewer, Detection Engineer, Cloud Security Analyst, Platform Administrator, Security Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant et environnement CMDR, origin, scope Cloud, période, owner, sources, versions, permissions, restrictions, classification, health, erreurs et return origin sont explicites. Une information inaccessible, absente, stale ou non supportée reste `unknown`, `partial`, `restricted` ou `blocked`; elle n’est jamais inventée. Provider configured ne signifie ni accessible ni couverture complète.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
|Resource and configuration observations|CAP-INV-604/608|cluster, workload and orchestration context|oui|selected versions|analysis partial|
|Audit/runtime observation materials|CAP-INV-607 / authorized sources|control-plane and available runtime observations|non|observed times|runtime unknown|
|Identity/network/storage relations|CAP-INV-605/612/613|principal and dependency context|non|session versions|relations incomplete|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Resource / Configuration Observation|Investigate|cluster/workload configuration|lecture|
|Identity / Network / Storage Observation|Investigate|related principals and resources|lecture/lien|
|Artifact / source material|respective owner|authorized manifests/logs/metadata|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Container Observation|créer, comparer, contester, versionner|Investigate|observation ≠ full container forensics|
|Orchestration Workload Observation|créer, qualifier, superseder|Investigate|provider object ≠ canonical CMDR object|
|Container/Endpoint handoff package|préparer|destination owner|no exec, pull or acquisition|

## 11. Fonctionnalités
- analyser cluster candidates, namespaces/scopes and workload topology.
- lier pod/task equivalents, images, identities and configuration.
- examiner network, storage, secret references and audit events.
- distinguer control-plane metadata from runtime observations.
- préparer handoff vers Endpoint, Artifact, Network or Detection owner.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Container and Orchestration Workload Analysis|DFIR Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|DFIR Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|DFIR Analyst|concepts CAP-INV-610|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Container and Orchestration Workload Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`candidate`, `observed`, `control-plane-only`, `runtime-partial`, `under-review`, `handoff-ready`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Container Observation|versioned observation|CAP-INV-615/616|observation layer and limitations explicit|
|Orchestration relation set|relationships|Shared Graph / analyst|provider concepts not promoted to canonical objects|
|Specialist handoff|package|Endpoint/Network/Artifact owner|no active command|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-604/608|orchestrated workload identified|CAP-INV-610|resources, config, images and source limits|inventory/configuration|
|CAP-INV-610|runtime or host evidence required|Endpoint/Artifact/Network owner|target references, gaps, permissions and return origin|container analysis|
|CAP-INV-610|candidate anomaly/timeline|CAP-INV-615/616|observations, events, relations and limitations|container analysis|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-604..609/612..618, Endpoint Agent, Network Forensics, Shared Graph, OPEN-008/012/013/014. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Container/orchestration metadata read|workload and topology disclosure|0|namespace/resource masking|possible|viewer/source owner|Settings/Security|Permissions|
|Specialist handoff prepare|target and sensitive reference leakage|2|target/secret references masked|OPEN-013|analyst/reviewer|Investigate/destination|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- orchestration object ≠ canonical CMDR object.
- container observation ≠ full container forensics.
- image reference ≠ image content.
- control-plane event ≠ runtime action certain.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- observations by control-plane/runtime coverage.
- workloads with explicit limitations.
- container commands or secret use — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Control-plane only
**Given** des manifests sont disponibles mais aucun runtime evidence  
**When** la workload est examinée  
**Then** l’état control-plane-only et les limitations sont visibles, sans claim de forensics complète.

### 2. Secret reference
**Given** un workload référence un secret  
**When** la relation est affichée  
**Then** seule la référence autorisée est visible, aucun reveal ou usage n’a lieu.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** l’orchestration est analysée  
**Then** trees, tables, graphs, manifest viewers and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
