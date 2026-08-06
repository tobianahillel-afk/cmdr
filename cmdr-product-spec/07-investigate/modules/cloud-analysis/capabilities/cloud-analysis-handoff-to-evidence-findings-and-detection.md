---
id: CAP-INV-617
title: Cloud Analysis Handoff to Evidence, Findings and Detection
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
  - OPEN-017
source-of-truth: canonical
---
# CAP-INV-617 — Cloud Analysis Handoff to Evidence, Findings and Detection

## 1. Définition
Préparer des Evidence Candidate Packages, Finding Drafts, Cloud Detection Gaps, Detection Engineering Packages, Threat Intelligence handoffs, Collection Requests and future response preparation without automatic qualification or execution.

## 2. Problème utilisateur
Une observation ou Hypothesis Cloud peut être transformée silencieusement en Evidence, Finding, rule ou response, contournant review, ownership, permissions and provenance.

## 3. Objectifs
- assembler observations, timeline, hypotheses, sources, restrictions, confidence, gaps and alternatives.
- préparer des packages distincts pour Evidence, Finding, Detection, TI, Collection and Govern.
- conserver return origin, destination owner, acceptance/rejection and no-effect boundary.

## 4. Non-objectifs
Aucune Evidence qualifiée, Finding confirmé, Detection Content, rule deployment, collection, Signal, Incident, Action Request approval or response execution.

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
|Qualified Cloud analysis context|CAP-INV-601..616|session, observations, hypotheses, timeline and provenance|oui|selected versions|handoff blocked|
|Destination requirements and permissions|Evidence / Case / Detection / TI / Collection / Govern owners|required fields, authority and restrictions|oui|current contract projection|package incomplete|
|Review disposition|Investigation Lead / Evidence Reviewer|human rationale and unresolved gaps|oui|current review|no submission|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Investigation Session and outputs|Investigate|selected analysis versions|lecture|
|Evidence / Finding / Detection / Action Request projections|destination owners|handoff requirements and status|lecture/lien|
|Artifacts and source restrictions|respective owners|custody, licence and handling|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Evidence Candidate Package|préparer, versionner, retirer, superseder|Investigate|candidate ≠ qualified Evidence|
|Finding Draft|préparer, contester, withdraw|Investigate|draft ≠ confirmed Finding|
|Cloud Detection Gap / Engineering Package|préparer/lier|Detection Engineering owner|package ≠ Detection Content or rule|
|TI / Collection / future response handoff|préparer|destination owner|handoff ≠ execution|

## 11. Fonctionnalités
- sélectionner analysis versions and source links.
- préparer Evidence Candidate with custody/restrictions.
- préparer Finding Draft with supporting/contradicting elements.
- préparer Detection Gap/Engineering, TI and Collection packages.
- routage, return, rejection, withdrawal and supersession without effect on target.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Analysis Handoff to Evidence, Findings and Detection|Investigation Lead|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Investigation Lead|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Investigation Lead|concepts CAP-INV-617|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Analysis Handoff to Evidence, Findings and Detection|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `incomplete`, `under-review`, `ready-for-handoff`, `submitted`, `returned`, `accepted-by-destination`, `rejected`, `withdrawn`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Evidence Candidate Package|versioned package|Evidence owner|not qualified Evidence|
|Finding Draft|draft|Case/Finding owner|not confirmed Finding|
|Detection Engineering Package|handoff package|Detection Engineering|no rule or deployment|
|Collection/TI/Govern preparation|handoff context|respective owner|no collection, attribution or response|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-615/616|human review supports handoff|CAP-INV-617|hypothesis, observations, timeline, gaps and provenance|analysis|
|CAP-INV-617|evidence/finding package submitted|Evidence/Case owner|candidate/draft, restrictions, versions and return origin|handoff|
|CAP-INV-617|detection/TI/collection need submitted|destination owner|gap, context, limits, permissions and return origin|handoff|
|Destination owner|accepted/returned/rejected|CAP-INV-617|status, reasons, destination object link and authority|destination|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-601..618, Evidence/Findings/Cases, Detection Engineering, Threat Intelligence, Collection and Live Response, Govern, Shared Linking/Trace, OPEN-012/013/014/015/017. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Evidence/Finding/Detection handoff prepare|cross-owner promotion and disclosure|2|source restrictions retained|OPEN-013|author/reviewer|Investigate/destination|Permissions|
|Collection or response preparation|target impact|2/3|target minimization|mandatory for execution|requester/approver/executor|Govern/destination|Future permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- Evidence candidate ≠ qualified Evidence.
- Finding Draft ≠ confirmed Finding.
- Detection package ≠ deployed rule.
- investigation action ≠ response action.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- packages by destination and disposition.
- packages preserving source versions/restrictions/return origin.
- automatic Evidence/Finding/rule/response — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Evidence candidate
**Given** des observations soutiennent une Hypothesis  
**When** un package Evidence est préparé  
**Then** il reste candidate, les gaps et contradictions sont visibles et seul l’owner Evidence peut qualifier.

### 2. Detection gap
**Given** une activité n’est couverte par aucune détection connue  
**When** le package est soumis  
**Then** aucune rule n’est créée ou déployée et Detection Engineering conserve son lifecycle.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** les handoffs sont préparés  
**Then** templates, checklists, source links and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
