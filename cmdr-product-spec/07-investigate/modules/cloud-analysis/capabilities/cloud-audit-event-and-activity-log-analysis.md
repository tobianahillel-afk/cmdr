---
id: CAP-INV-607
title: Cloud Audit Event and Activity Log Analysis
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
# CAP-INV-607 — Cloud Audit Event and Activity Log Analysis

## 1. Définition
Analyser des événements d’audit et activity logs Cloud autorisés en conservant actor/principal candidate, action, resource, result, timestamp, region, request context, source, duplicates, lateness, gaps and correlation.

## 2. Problème utilisateur
Un événement d’audit peut être attribué directement à un utilisateur humain, tandis que les périodes absentes, retards et duplications sont masqués dans une timeline apparemment complète.

## 3. Objectifs
- normaliser fonctionnellement les champs disponibles sans imposer de schéma provider.
- préserver actor/principal candidate, action, resource, result, timestamps, region, request context, source and errors.
- détecter périodes manquantes, late events, duplicates et correlations candidates sans confirmer l’intention.

## 4. Non-objectifs
Aucun langage de requête, API de logs, ingestion, parser déployé, suppression de doublons source, attribution humaine certaine, alerte runtime ou réponse.

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
|Audit/activity materials|Platform Settings / Artifact owner|authorized event records|oui|event and collection period|no event analysis|
|Cloud scope and identity context|CAP-INV-603/605|units, principals, aliases and limitations|oui|session versions|actor unresolved|
|Source health and coverage|Platform Settings|ingestion gaps, parser/schema and latency|oui|current projection|timeline misleading|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Telemetry Event / source event projection|source owner / Command projection where applicable|event fields and source lineage|lecture limitée|
|Identity / Resource Observations|Investigate|candidate actor and target context|lecture/lien|
|Data Source / Parser / health|Platform Settings|coverage, latency and errors|lecture|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Activity Observation|créer, comparer, annoter, contester, versionner|Investigate|audit event ≠ confirmed human action|
|Audit coverage and gap assessment|créer, réviser|Investigate|absence of logs ≠ absence of activity|
|Event correlation candidate|créer/lier, dispute|Investigate using Shared Linking|correlation ≠ causality|

## 11. Fonctionnalités
- inspecter audit events et activity logs avec source et versions.
- qualifier principal candidate, action, resource, result, region and request context.
- gérer event time, recorded time, collection time, lateness and duplicates.
- identifier missing periods and parser/source errors.
- corréler avec identities, resources and other evidence without automatic attribution.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Audit Event and Activity Log Analysis|DFIR Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|DFIR Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|DFIR Analyst|concepts CAP-INV-607|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Audit Event and Activity Log Analysis|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`unparsed`, `partial`, `observed`, `duplicate-candidate`, `late`, `gap-detected`, `correlated`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Activity Observation|versioned observation|CAP-INV-608..616|event/source times and uncertainty preserved|
|Audit gap assessment|assessment|session / Settings / CAP-INV-617|gap not absence of activity|
|Correlation candidate|link set|Shared Timeline/Graph|no causality or human attribution|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-603/605|audit scope ready|CAP-INV-607|units, principals, periods and sources|scope/identity|
|CAP-INV-607|resource or workload activity identified|CAP-INV-608..613|events, targets, results and time context|audit analysis|
|CAP-INV-607|timeline reconstruction requested|CAP-INV-616|events, gaps, duplicates, times and correlations|audit analysis|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-603..616, Settings sources/parsers/health, Shared Timeline/Linking, Command telemetry projections, OPEN-008/012/014. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Audit log read/search/compare|sensitive activity and cross-tenant data|0/1|payload and identity masking|possible|viewer/source owner|Settings/Security|Permissions|
|Activity observation/correlation create|false attribution or causality|2|restricted context masked|OPEN-013|analyst/reviewer|Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- audit event ≠ confirmed user action.
- principal field ≠ human identity.
- duplicate event ≠ corroboration.
- absence de logs ≠ absence d’activité.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- events by source, result and lateness.
- missing periods and duplicate candidates.
- human attribution without evidence — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Événement tardif
**Given** un événement arrive après la première reconstruction  
**When** il est ajouté  
**Then** son observed/recorded/received time et l’impact sur la timeline sont visibles sans écraser la version antérieure.

### 2. Principal ambigu
**Given** un event contient un assumed-role ou workload principal  
**When** l’actor est affiché  
**Then** il reste principal candidate et n’est pas présenté comme personne certaine.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** les logs sont analysés  
**Then** parsers déterministes, search, tables, timelines et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
