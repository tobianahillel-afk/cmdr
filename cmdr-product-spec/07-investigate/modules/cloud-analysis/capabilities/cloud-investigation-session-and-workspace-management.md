---
id: CAP-INV-602
title: Cloud Investigation Session and Workspace Management
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
# CAP-INV-602 — Cloud Investigation Session and Workspace Management

## 1. Définition
Créer et gérer une Cloud Investigation Session versionnée qui conserve scope, owners, contributeurs, sources, requêtes fonctionnelles, Tools, observations, Hypotheses, Artifacts, erreurs, handoffs et supersession.

## 2. Problème utilisateur
Les analyses Cloud distribuées entre comptes, services et outils peuvent perdre leurs versions, leurs erreurs ou leurs décisions humaines et devenir impossibles à reproduire.

## 3. Objectifs
- créer, assigner, suspendre, reprendre, clôturer et rouvrir une session sans modifier une cible Cloud.
- conserver scope, accounts, services, sources, requêtes fonctionnelles, Tools, Tool Calls, Automation Runs, vues, observations, Hypotheses et Artifacts.
- versionner les changements, conflits, handoffs et supersession avec retour vers l’origine.

## 4. Non-objectifs
Aucun workspace technique final, terminal Cloud, console provider, commande, connecteur, scheduler, runtime d’agent, réponse ou écran détaillé.

Aucun choix AWS, Azure, GCP, SaaS, Kubernetes, provider, format, query language ou implementation n’est imposé.

## 5. Propriétaire
Investigate possède uniquement les concepts analytiques Cloud locaux décrits ici. Platform Settings conserve providers, connectors, credentials, secrets, configured organizations/tenants/accounts/subscriptions/projects, ingestion, schemas, parsers, health, retention, storage, policies et configuration administrative. Command conserve runtime Detection, Signal, Alert, Incident, priorité et dispositions. Govern conserve Decision, Approval, Action Request, Response Run, Result et toute mutation de cible. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Shared conserve Entity, Graph, Timeline, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration et Recovery. Aucun owner concurrent.

## 6. Utilisateurs
Principal : **Investigation Lead**. Secondaires : Investigation Lead, SOC Analyst, DFIR Analyst, Evidence Reviewer, Detection Engineer, Cloud Security Analyst, Platform Administrator, Security Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant et environnement CMDR, origin, scope Cloud, période, owner, sources, versions, permissions, restrictions, classification, health, erreurs et return origin sont explicites. Une information inaccessible, absente, stale ou non supportée reste `unknown`, `partial`, `restricted` ou `blocked`; elle n’est jamais inventée. Provider configured ne signifie ni accessible ni couverture complète.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
|Qualified Cloud Analysis Intake|CAP-INV-601|scope, objectives, sources, restrictions and gaps|oui|selected version|no session|
|Participants and authority|Investigate / Security|owner, contributors, reviewers and permissions|oui|current|blocked|
|Workspace materials|Artifacts / sources / Tools|observations, queries, Tool Calls and prior runs|non|versioned|empty session|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Analysis Intake|Investigate|qualified scope and return origin|lecture/lien|
|Artifacts, sources and observations|respective owner|authorized versions and restrictions|lecture limitée|
|Tool / Tool Call / Automation Run|Studio|execution metadata and results|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Investigation Session|créer, assigner, suspendre, reprendre, clôturer, rouvrir, versionner, superseder|Investigate|session ≠ provider workspace|
|Session membership and view state|créer, modifier, retirer|Investigate using Shared collaboration|membership ≠ source permission|
|Session lifecycle event|émettre|Shared Trace/Activity|history preserved|

## 11. Fonctionnalités
- conserver owner, contributeurs, reviewers et return origin.
- organiser accounts, services, sources, requêtes fonctionnelles, Tools et vues.
- lier Tool Calls, Automation Runs, observations, Hypotheses, Artifacts, erreurs et handoffs.
- gérer pause, reprise, clôture, réouverture, version, conflit et supersession.
- préserver les droits source lors du partage de session.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Investigation Session and Workspace Management|Investigation Lead|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Investigation Lead|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Investigation Lead|concepts CAP-INV-602|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Investigation Session and Workspace Management|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `active`, `paused`, `blocked`, `partial`, `under-review`, `closed`, `reopened`, `archived`, `superseded`, `conflict`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Investigation Session|versioned session|CAP-INV-603..618|scope, membership, sources and versions preserved|
|Session lifecycle event|trace event|Shared Activity / auditor|actor, reason, time and prior state preserved|
|Session handoff context|context package|Case / Detection / Evidence consumer|return origin and restrictions retained|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-601|preconditions ready|CAP-INV-602|intake, owner, scope, sources and gaps|intake|
|CAP-INV-602|scope requires qualification|CAP-INV-603|session, candidate hierarchy, periods and exclusions|session|
|CAP-INV-602|analysis closed or reopened|origin owner|disposition, results, gaps, versions and return origin|session|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-601/603..618, Shared Collaboration/Versioning/Trace/Recovery, Studio Tool Calls/Runs, Cases, OPEN-012/013/014/015. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Session create/update/close/reopen|scope and evidence integrity|2|restricted content masked|OPEN-013|author/reviewer|Investigate|Permissions|
|Session membership read/update|cross-tenant disclosure|0/2|tenant scoped|possible|owner/access reviewer|Investigate/Security|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- Cloud Investigation Session ≠ Cloud account or provider console.
- session membership ≠ permission on sources.
- closed ≠ immutable or confirmed Finding.
- Tool Call result ≠ analyst conclusion.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- sessions by state, owner and origin.
- sessions with explicit versions and return origin.
- orphan Tool Calls or silent history deletion — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Conflit de version
**Given** deux analystes modifient le scope en parallèle  
**When** la session est enregistrée  
**Then** un état conflict, un diff et les deux versions sont conservés sans écrasement silencieux.

### 2. Source restreinte
**Given** un contributeur peut voir la session mais pas une source liée  
**When** il ouvre le workspace  
**Then** la source reste masquée et la membership ne lui accorde aucun droit.

### 3. Sans IA
**Given** aucun agent n’est disponible  
**When** la session est gérée  
**Then** formulaires, tabs, tables, timeline, diff et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
