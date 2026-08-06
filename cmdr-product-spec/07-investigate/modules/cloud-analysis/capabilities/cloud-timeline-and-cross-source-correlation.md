---
id: CAP-INV-616
title: Cloud Timeline and Cross-Source Correlation
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
# CAP-INV-616 — Cloud Timeline and Cross-Source Correlation

## 1. Définition
Construire une Cloud Timeline versionnée et des correlations candidates en distinguant observed, recorded, reconstructed, estimated and absent timestamps, timezones, gaps and conflicts across Cloud, Endpoint, Network, Memory, Disk, Detection and Threat Intelligence.

## 2. Problème utilisateur
Des horodatages de sources différentes peuvent être alignés comme faits certains, effaçant timezones, lateness, estimation, gaps et conflits, ou remplaçant la Case Timeline.

## 3. Objectifs
- conserver chaque type d’horodatage, timezone, source, precision and transformation.
- reconstruire une timeline avec gaps, conflicts, duplicates and estimated intervals.
- corréler avec Case Timeline, Endpoint, Network, Memory, Disk, Detection Engineering and Threat Intelligence sans remplacer leurs sources.

## 4. Non-objectifs
Aucune horloge corrigée dans les sources, causalité automatique, Case Timeline replacement, event deletion, universal time model or response trigger.

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
|Cloud observations and events|CAP-INV-604..615|timestamps, sources, versions and relations|oui|selected versions|timeline partial|
|External timeline projections|Case / Endpoint / Network / Memory / Disk / Detection / TI owners|authorized events and time context|non|source-defined|correlation limited|
|Time quality context|source owners / Settings|timezone, clock, latency and gaps|oui|observed period|precision unknown|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Observation families|Investigate|time-bearing observations|lecture|
|Case Timeline / forensic timelines|respective owner|authorized timeline projection|lecture/lien|
|Threat Intelligence / Detection projections|Investigate submodules / Command|Sightings, changes and detections|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Timeline|créer, reconstruire, comparer, contester, versionner, superseder|Investigate using Shared Timeline|Cloud Timeline ≠ Case Timeline|
|Cross-source Correlation Candidate|créer/lier, dispute, withdraw|Investigate using Shared Linking|correlation ≠ causality|
|Time quality assessment|créer, réviser|Investigate|estimated ≠ observed|

## 11. Fonctionnalités
- distinguer observed, recorded, reconstructed, estimated and absent timestamps.
- conserver timezone, source, precision, lateness and transformations.
- assembler gaps, conflicts and alternative orderings.
- corréler Cloud avec Endpoint, Network, Memory, Disk, Detection and TI.
- préserver source ownership and return navigation.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Timeline and Cross-Source Correlation|DFIR Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|DFIR Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|DFIR Analyst|concepts CAP-INV-616|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Timeline and Cross-Source Correlation|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `partial`, `reconstructed`, `conflicted`, `under-review`, `correlated`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Timeline|versioned timeline|CAP-INV-617 / Case consumer|time types and gaps visible|
|Cross-source Correlation Candidate|relationship set|analyst / Case / Detection|not causality or compromise|
|Time quality assessment|assessment|reviewer / auditor|precision and transformations explicit|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-604..615|timeline requested|CAP-INV-616|observations, timestamps, gaps and hypotheses|source capability|
|External timeline owner|authorized projection linked|CAP-INV-616|events, timezone, precision and restrictions|external source|
|CAP-INV-616|handoff prepared|CAP-INV-617|timeline version, correlations, conflicts and source links|timeline|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-604..618, Case Timeline, Endpoint/Network/Memory/Disk, Detection Engineering, Threat Intelligence, Shared Timeline/Linking/Versioning, OPEN-012/014/015. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Timeline and source projection read|cross-product/cross-tenant disclosure|0|source and tenant masking|possible|viewer/source owner|respective owners|Permissions|
|Correlation create/dispute/export|misleading linkage or sensitive export|1/2|minimization/redaction|OPEN-013|analyst/reviewer|Investigate/Shared|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- Cloud Timeline ≠ Case Timeline.
- reconstructed timestamp ≠ observed timestamp.
- correlation ≠ causality.
- no match ≠ no attack.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- timeline entries by timestamp type.
- gaps/conflicts and estimated intervals.
- silent timezone conversion or source mutation — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Horodatages conflictuels
**Given** deux sources placent un changement dans un ordre différent  
**When** la timeline est reconstruite  
**Then** les deux times, timezones, precision and conflict remain visible.

### 2. Case Timeline
**Given** une Cloud Timeline est liée à un Case  
**When** elle est consultée  
**Then** elle reste une projection distincte et ne remplace pas la Case Timeline canonique.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** la timeline est construite  
**Then** sorting, time normalization rules, tables, diff and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
