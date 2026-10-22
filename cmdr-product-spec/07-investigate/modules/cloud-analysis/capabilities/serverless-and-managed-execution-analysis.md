---
id: CAP-INV-611
title: Serverless and Managed Execution Analysis
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
# CAP-INV-611 — Serverless and Managed Execution Analysis

## 1. Définition
Analyser functions, jobs, managed services, versions, triggers, identities, permissions, configurations, environment metadata, network/storage relations, invocation observations, errors, deployments and changes.

## 2. Problème utilisateur
Un invocation record ou un déploiement observé peut être interprété comme exécution malveillante certaine, alors que le code, l’environnement et la couverture peuvent être absents.

## 3. Objectifs
- représenter functions/jobs/services et versions sans imposer un provider.
- lier triggers, identities, permissions, configuration, environment metadata, network and storage relations.
- analyser invocation observations, errors, deployments and changes avec limites et provenance.

## 4. Non-objectifs
Aucune invocation, code download, environment secret reveal, redeploy, disable, permission change ou confirmation automatique de malicious execution.

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
|Serverless/managed resource observations|CAP-INV-604/608|service, function/job, version and configuration|oui|selected versions|no analysis|
|Invocation and deployment events|CAP-INV-607 / authorized source|observed invocations, results, errors and deployments|non|event periods|execution unknown|
|Identity/permission/network/storage context|CAP-INV-605/606/612/613|related principals and resources|non|session versions|context partial|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Resource / Configuration Observation|Investigate|function/service and version context|lecture|
|Cloud Activity Observation|Investigate|invocation/deployment candidates|lecture/lien|
|Identity / Permission / dependency observations|Investigate|authorization and relations|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Serverless Observation|créer, comparer, contester, versionner|Investigate|invocation record ≠ malicious execution|
|Managed Execution Change Assessment|créer, réviser, superseder|Investigate|deployment observed ≠ code content known|
|Artifact/Detection handoff proposal|préparer|destination owner|proposal ≠ collection or rule|

## 11. Fonctionnalités
- analyser functions, jobs, services and versions.
- lier triggers, identities, permissions, configuration and environment metadata.
- cartographier network/storage dependencies.
- examiner invocation observations, errors and deployments.
- qualifier changes, gaps and handoffs without executing anything.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Serverless and Managed Execution Analysis|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-611|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Serverless and Managed Execution Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`observed`, `partial`, `invocation-observed`, `deployment-change`, `under-review`, `disputed`, `handoff-ready`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Serverless Observation|versioned observation|CAP-INV-615/616|event and configuration evidence distinguished|
|Managed Execution Change Assessment|assessment|session / reviewer|deployment ≠ malicious code|
|Artifact/Detection handoff|proposal|Artifact/Detection Engineering owner|no retrieval or rule|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-604/608|serverless resource identified|CAP-INV-611|resource, versions, triggers and config|inventory/configuration|
|CAP-INV-607|invocation or deployment event observed|CAP-INV-611|event, principal candidate, result and time|audit analysis|
|CAP-INV-611|anomaly/timeline/handoff requested|CAP-INV-615..617|observations, changes, gaps and dependencies|serverless analysis|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-604..608/612..618, Artifact Management, Detection Engineering, Shared Timeline/Linking, OPEN-008/012/014. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Serverless metadata/event read|code and environment disclosure|0|environment/sensitive fields masked|possible|viewer/source owner|Settings/Security|Permissions|
|Handoff prepare|code or sensitive material scope|2|secret values excluded|OPEN-013|analyst/reviewer|Investigate/destination|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- serverless invocation record ≠ malicious execution.
- deployment observation ≠ code content.
- environment metadata ≠ secret value.
- trigger relation ≠ confirmed causality.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- serverless observations by version and event coverage.
- unresolved triggers/dependencies.
- invocations executed by Investigate — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Invocation isolée
**Given** une invocation échoue sans autre evidence  
**When** elle est examinée  
**Then** le record reste observation avec result/error, pas malicious execution.

### 2. Environment metadata
**Given** une clé d’environnement semble sensible  
**When** elle est affichée  
**Then** seul le nom/metadata autorisé est visible et CAP-INV-614 gère le candidat.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** serverless est analysé  
**Then** tables, event timelines, dependency graphs et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
