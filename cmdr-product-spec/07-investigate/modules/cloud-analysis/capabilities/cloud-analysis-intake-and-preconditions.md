---
id: CAP-INV-601
title: Cloud Analysis Intake and Preconditions
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
# CAP-INV-601 — Cloud Analysis Intake and Preconditions

## 1. Définition
Qualifier l’ouverture d’une analyse Cloud depuis un Case, Incident, Finding, Hunt, Signal, source Cloud ou Artifact, en rendant explicites le scope, les sources disponibles, la couverture, les restrictions et les conditions d’accès avant toute analyse.

## 2. Problème utilisateur
Un contexte Cloud peut être ouvert avec un provider supposé accessible, des périodes manquantes ou une couverture partielle non signalée, ce qui transforme des lacunes de collecte en conclusions trompeuses.

## 3. Objectifs
- ouvrir depuis Case, Incident, Finding, Hunt, Signal, source Cloud ou Artifact sans transférer leur ownership.
- qualifier organization, tenant, account, subscription, project, période, services, sources, health, restrictions et permissions.
- exposer logs absents, couverture partielle, erreurs, résultats antérieurs, Tools compatibles, objectifs et limites sans lancer d’analyse automatique.

## 4. Non-objectifs
Aucune connexion réelle à un provider, collecte, requête provider, API, protocole, langage de requête, scanner, commande Cloud, modification de tenant, permission Cloud, réponse ou code produit.

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
|Investigation origin|Case / Incident / Finding / Hunt / Signal owner|authorized origin and return context|oui|selected version|intake blocked|
|Cloud source or Artifact context|Platform Settings / Artifact owner|provider-neutral source projection|selon origine|declared capture period|source gap visible|
|Scope and access context|Platform Settings / Security|organization, tenant, account, project, health, restrictions and permissions|oui|current projection|blocked or partial|
|Prior results and compatible Tools|Investigate / Studio|previous observations, errors and compatible capabilities|non|versioned|manual intake remains|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Case / Incident / Finding / Hunt / Signal|respective owner|origin, objective, constraints and return origin|lecture/lien|
|Cloud source / connector projection|Platform Settings|provider, configured scope, health, coverage and limitations|lecture limitée|
|Artifact / prior result / Tool projection|respective owner / Studio|available material, versions, compatibility and errors|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Analysis Intake|créer, qualifier, versionner, retirer, superseder|Investigate|intake ≠ analysis result|
|Cloud Preconditions Assessment|créer, contester, réviser|Investigate|accessible ≠ complete coverage|
|Collection or access gap proposal|préparer/lier|destination owner|proposal ≠ collection or permission grant|

## 11. Fonctionnalités
- sélectionner l’origine et conserver le return origin.
- qualifier provider-neutral scope, période, services, sources, permissions et restrictions.
- exposer health, logs absents, périodes manquantes, erreurs et couverture partielle.
- lier résultats antérieurs et Tools compatibles sans les exécuter.
- produire une Preconditions Assessment réversible.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Analysis Intake and Preconditions|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-601|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Analysis Intake and Preconditions|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `incomplete`, `access-review-required`, `source-unavailable`, `coverage-partial`, `ready`, `blocked`, `withdrawn`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Analysis Intake|versioned intake|CAP-INV-602/603|origin, scope, restrictions, gaps and return origin visible|
|Cloud Preconditions Assessment|assessment|Investigation Lead / source owner|readiness ≠ provider accessibility guarantee|
|Gap proposal|handoff proposal|Settings / Collection owner|no collection or permission change|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|Case/Incident/Finding/Hunt/Signal|cloud analysis requested|CAP-INV-601|origin, objective, tenant/environment, evidence links and return origin|origin|
|Cloud source or Artifact|analysis requested|CAP-INV-601|source version, capture context, restrictions and health|source|
|CAP-INV-601|preconditions accepted|CAP-INV-602/603|intake, scope candidates, periods, sources, gaps and permissions|intake|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
Cases, Signals/Hunt, Collection and Live Response, Platform Settings sources/health, Studio Tool catalogue, Shared Linking/Trace, OPEN-008/012/013/014/015. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Cloud source and scope read|cross-tenant or restricted metadata|0|tenant/source scoped|possible|viewer/source owner|Settings / Security|Permissions|
|Intake create/update/withdraw|scope contamination or premature readiness|2|restricted inputs masked|OPEN-013|author/reviewer|Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- configured provider ≠ accessible provider.
- accessible provider ≠ complete coverage.
- absence de logs ≠ absence d’activité.
- intake ready ≠ analysis completed.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- intakes by origin and readiness state.
- intakes with explicit coverage and missing periods.
- automatic analysis at opening — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Couverture partielle
**Given** un Case demande une analyse mais une période de logs manque  
**When** l’intake est qualifié  
**Then** la période manquante, le health et l’impact restent visibles et l’état n’est pas ready sans disposition humaine.

### 2. Provider inaccessible
**Given** un provider est configuré mais l’accès est refusé  
**When** l’analyste ouvre l’intake  
**Then** configured et accessible restent distincts, aucune donnée n’est inventée et un gap est préparé.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** l’intake est préparé  
**Then** formulaires, catalogues, health, checklists et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
