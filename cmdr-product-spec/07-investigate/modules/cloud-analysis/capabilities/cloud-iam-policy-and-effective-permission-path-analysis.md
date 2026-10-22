---
id: CAP-INV-606
title: Cloud IAM Policy and Effective Permission Path Analysis
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
# CAP-INV-606 — Cloud IAM Policy and Effective Permission Path Analysis

## 1. Définition
Analyser policies, bindings, grants, denies, héritages, trusts et contextes afin de produire des Effective Permission Candidates et permission paths explicables, sans fournir de procédure d’escalade ni confirmer un exploit.

## 2. Problème utilisateur
Une policy déclarée peut être confondue avec l’autorisation effective, tandis qu’un chemin théorique de permission peut être présenté comme chemin d’exploitation certain.

## 3. Objectifs
- représenter grants, denies, conditions, héritages, bindings et trusts avec source et version.
- calculer ou documenter des effective permission candidates bornés par les données disponibles.
- conserver permission paths, contradictions, cross-account trust candidates, limitations et hypotheses d’escalation sans procédure offensive.

## 4. Non-objectifs
Aucun changement IAM, simulation d’abus active, procédure d’escalade, utilisation de credential, exploitation, policy deployment ou décision de conformité.

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
|Identity and Role Observations|CAP-INV-605|principals, roles, assignments and ambiguity|oui|selected versions|paths incomplete|
|IAM policy materials|authorized cloud sources|policies, bindings, grants, denies and conditions|oui|observed timestamp|no assessment|
|Scope and resource hierarchy|CAP-INV-603/604|inheritance and resource context|oui|session version|effective state unknown|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Identity / Role Observations|Investigate|candidate principals and assignments|lecture|
|Policy / binding projections|Platform Settings / source owner|declared authorization data|lecture limitée|
|Entity / Graph relationships|Shared|candidate paths and relation provenance|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Permission Observation|créer, comparer, contester, versionner|Investigate|declared policy ≠ effective authorization|
|Effective Permission Candidate|créer, qualifier, retirer, superseder|Investigate|candidate ≠ certain permission|
|Permission Path Assessment|créer, réviser, dispute|Investigate using Shared Graph|path ≠ exploit path|

## 11. Fonctionnalités
- analyser policies, bindings, grants, denies, conditions et héritages.
- construire des effective permission candidates avec méthode explicite.
- cartographier permission paths et cross-account trust candidates.
- comparer declared et observed state.
- gérer contradictions, unsupported semantics et escalation hypotheses sans procédure d’abus.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud IAM Policy and Effective Permission Path Analysis|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-606|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud IAM Policy and Effective Permission Path Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `calculating`, `partial`, `candidate`, `supported`, `contradicted`, `under-review`, `disputed`, `withdrawn`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Permission Observation set|versioned observations|CAP-INV-615/616/617|source, policy version and conditions visible|
|Effective Permission Candidate|candidate assessment|analyst / Detection handoff|candidate not confirmed entitlement|
|Permission Path Assessment|graph/table assessment|Case / reviewer|path not exploit or compromise|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-605|identity context ready|CAP-INV-606|principals, roles, assignments and ambiguity|identity analysis|
|CAP-INV-603/604|hierarchy or resource changes|CAP-INV-606|scope, inheritance context and versions|scope/inventory|
|CAP-INV-606|anomaly candidate identified|CAP-INV-615|path, sources, supporting/contradicting elements and limits|IAM assessment|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-603..605/607/615..618, Shared Graph, Settings policies/sources, Security permissions, OPEN-008/012/013. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|IAM policy and permission path read|privilege disclosure|0/1|sensitive conditions masked|possible|viewer/security reviewer|Settings/Security|Permissions|
|Permission candidate/path create or dispute|false privilege conclusion|2|restricted principals masked|OPEN-013|analyst/reviewer|Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- declared policy ≠ effective authorization.
- role assignment ≠ effective permission certaine.
- permission path ≠ exploit path.
- privilege candidate ≠ privilege abuse.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- permission candidates by confidence and source completeness.
- paths with explicit grants/denies/conditions.
- exploit procedure generated — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Deny conditionnel
**Given** un grant est accompagné d’un deny conditionnel non évalué  
**When** le chemin est calculé  
**Then** le résultat reste candidate/partial et la condition non résolue est visible.

### 2. Chemin cross-account
**Given** une chaîne de trusts relie deux accounts  
**When** elle est affichée  
**Then** elle est décrite comme permission path candidate, jamais comme exploitation ou compromission.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** les permissions sont analysées  
**Then** règles déterministes, arbres, graphes, tables de grants/denies et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
