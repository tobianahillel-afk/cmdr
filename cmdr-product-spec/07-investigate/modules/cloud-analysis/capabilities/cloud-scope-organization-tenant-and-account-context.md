---
id: CAP-INV-603
title: Cloud Scope, Organization, Tenant and Account Context
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
# CAP-INV-603 — Cloud Scope, Organization, Tenant and Account Context

## 1. Définition
Représenter et comparer le périmètre Cloud sans confondre organization, tenant, account, subscription, project, folder, region, environment, management hierarchy et relations cross-account.

## 2. Problème utilisateur
Des hiérarchies de providers différentes peuvent être normalisées trop tôt, faisant passer une relation administrative, un alias ou un périmètre incomplet pour une frontière de sécurité certaine.

## 3. Objectifs
- décrire les unités de portée avec leur terminologie source et une projection fonctionnelle provider-neutral.
- définir inclusions, exclusions, périodes, régions, environnements, hiérarchies et couverture.
- conserver cross-account et cross-tenant relations candidates avec source, confiance, contradictions et limites.

## 4. Non-objectifs
Aucun modèle provider imposé, découverte active, création de compte, changement d’organisation, permission cross-account, schéma canonique final ou migration de tenant.

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
|Intake and Session scope candidates|CAP-INV-601/602|candidate hierarchy and objectives|oui|current session version|scope incomplete|
|Configured scope projections|Platform Settings|organizations, tenants, accounts, subscriptions, projects and regions|selon accès|current projection|unknown|
|Source hierarchy observations|Cloud materials / audit / inventory sources|observed parent-child and trust context|non|observed timestamp|unverified|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Tenant / Environment|Platform Settings / canonical owner|administrative boundary projection|lecture limitée|
|Cloud scope units|Platform Settings|provider terminology, identifiers, hierarchy and access state|lecture limitée|
|Entity / Graph links|Shared|candidate relations and provenance|lecture/lien|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Scope Assessment|créer, comparer, contester, versionner, superseder|Investigate|assessment ≠ configured scope|
|Scope inclusion/exclusion set|créer, modifier, retirer|Investigate|explicit and time-bounded|
|Cross-account relation observation|créer, qualifier, contester|Investigate using Shared links|relation ≠ compromise|

## 11. Fonctionnalités
- représenter organization, tenant, account, subscription, project, folder, region et environment séparément.
- conserver la terminologie source, les aliases et les identifiers déclarés.
- définir scope inclus/exclus, périodes et couverture.
- cartographier management hierarchy et relations cross-account/cross-tenant candidates.
- comparer des snapshots sans conclure à un état courant certain.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Scope, Organization, Tenant and Account Context|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-603|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Scope, Organization, Tenant and Account Context|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `incomplete`, `observed`, `partially-mapped`, `under-review`, `disputed`, `accepted-for-session`, `stale`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Scope Assessment|versioned assessment|CAP-INV-604..618|units, hierarchy, periods, exclusions and uncertainty explicit|
|Scope graph projection|relation set|Shared Graph / analyst|candidate relations retain sources|
|Scope gap|gap record|Settings / Collection owner|gap ≠ account creation or access grant|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-601/602|scope qualification requested|CAP-INV-603|candidate units, periods, exclusions and sources|session|
|CAP-INV-603|scope accepted for analysis|CAP-INV-604/605/607|units, hierarchy, regions, coverage and limits|scope assessment|
|Settings|configured scope changed|CAP-INV-603|new projection, timestamp and diff|settings|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-601/602/604..618, Platform Settings tenants/environments/sources, Shared Entity/Graph/Versioning, OPEN-008/012/013. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Cloud scope read/select|cross-tenant exposure|0/2|tenant and unit scoped|possible|analyst/access reviewer|Settings/Investigate|Permissions|
|Cross-tenant relation review|sensitive trust disclosure|0/2|relationship masking|OPEN-013|author/reviewer|Investigate/Security|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- Cloud Account ≠ tenant ≠ organization ≠ subscription ≠ project.
- configured scope ≠ accessible scope.
- management relation ≠ trust or compromise.
- scope snapshot ≠ current state certain.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- scope units by type and source.
- scopes with explicit inclusions, exclusions and coverage.
- collapsed provider-specific units — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Unités homonymes
**Given** un provider utilise account et un autre project pour des périmètres différents  
**When** les scopes sont comparés  
**Then** les termes source et projections restent distincts sans fusion sémantique.

### 2. Cross-tenant
**Given** une relation candidate traverse deux tenants  
**When** elle est ajoutée au scope  
**Then** la relation, l’autorité, les limites et les permissions sont explicites sans extension de droit.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** le scope est cartographié  
**Then** arbres, tables, graphes avec alternative tabulaire et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
