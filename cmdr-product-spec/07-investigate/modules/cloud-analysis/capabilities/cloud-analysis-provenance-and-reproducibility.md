---
id: CAP-INV-618
title: Cloud Analysis Provenance and Reproducibility
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
# CAP-INV-618 — Cloud Analysis Provenance and Reproducibility

## 1. Définition
Retracer de bout en bout origin, provider-neutral scope, sources, connector projections, collection context, permissions, sessions, Tools, Tool Calls, Automation Runs, functional queries, observations, sensitive accesses, timelines, handoffs, errors, limitations, human decisions and versions.

## 2. Problème utilisateur
Une conclusion Cloud peut sembler reproductible alors que ses sources, permissions, versions, transformations, Tool Calls ou accès sensibles sont manquants ou ont changé.

## 3. Objectifs
- conserver la chaîne de provenance depuis Case/Incident jusqu’aux handoffs.
- évaluer la reproductibilité selon disponibilité, permissions, versions, sources, parameters and environment.
- permettre comparaison, replay conceptuel, correction and supersession sans exécuter de requête provider ni supprimer l’historique.

## 4. Non-objectifs
Aucun replay actif, capture automatique, content retention policy, immutable ledger implementation, API, command, secret reuse or reproducibility guarantee.

Aucun choix AWS, Azure, GCP, SaaS, Kubernetes, provider, format, query language ou implementation n’est imposé.

## 5. Propriétaire
Investigate possède uniquement les concepts analytiques Cloud locaux décrits ici. Platform Settings conserve providers, connectors, credentials, secrets, configured organizations/tenants/accounts/subscriptions/projects, ingestion, schemas, parsers, health, retention, storage, policies et configuration administrative. Command conserve runtime Detection, Signal, Alert, Incident, priorité et dispositions. Govern conserve Decision, Approval, Action Request, Response Run, Result et toute mutation de cible. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Shared conserve Entity, Graph, Timeline, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration et Recovery. Aucun owner concurrent.

## 6. Utilisateurs
Principal : **Evidence Reviewer**. Secondaires : Investigation Lead, SOC Analyst, DFIR Analyst, Evidence Reviewer, Detection Engineer, Cloud Security Analyst, Platform Administrator, Security Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant et environnement CMDR, origin, scope Cloud, période, owner, sources, versions, permissions, restrictions, classification, health, erreurs et return origin sont explicites. Une information inaccessible, absente, stale ou non supportée reste `unknown`, `partial`, `restricted` ou `blocked`; elle n’est jamais inventée. Provider configured ne signifie ni accessible ni couverture complète.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
|Cloud session and analysis lineage|CAP-INV-601..617|origins, scopes, observations, hypotheses, timelines and packages|oui|selected versions|assessment incomplete|
|Tool and automation lineage|Studio|Tools, Tool Calls, Runs, parameters, errors and producer versions|selon usage|versioned|manual path only|
|Source/access/version context|Settings / source owners / Security|connector projection, collection context, permissions and restrictions|oui|relevant time|not reproducible|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Case / Incident / Artifact / Evidence projections|respective owners|origin and custody lineage|lecture/lien|
|Cloud analysis concepts|Investigate|all selected versions and dispositions|lecture|
|Tool Call / Automation Run / Trace|Studio / Shared|execution and audit lineage|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Reproducibility Assessment|créer, comparer, contester, versionner, superseder|Investigate|assessment ≠ guaranteed replay|
|Cloud Analysis Provenance Package|créer, exporter selon permission, withdraw|Investigate using Shared Trace/Export|export ≠ source access extension|
|Correction or gap proposal|préparer|source/destination owner|proposal ≠ mutation|

## 11. Fonctionnalités
- retracer origin, provider, tenant/account/project and source/connector projections.
- lier collection context, permissions, sessions, Tools, Tool Calls, Automation Runs and functional queries.
- conserver observations, Hypotheses, Artifacts, sensitive accesses, timelines and correlations.
- enregistrer handoffs, errors, limitations, human decisions and versions.
- évaluer reproducibility and prepare correction/gap without active replay.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Analysis Provenance and Reproducibility|Evidence Reviewer|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Evidence Reviewer|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Evidence Reviewer|concepts CAP-INV-618|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Analysis Provenance and Reproducibility|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`draft`, `lineage-partial`, `reproducible-with-current-access`, `conditionally-reproducible`, `not-reproducible`, `under-review`, `disputed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Analysis Provenance Package|versioned package|auditor / Evidence / destination owner|lineage and restrictions preserved|
|Reproducibility Assessment|assessment|Investigation Lead / QA|conditions and missing elements visible|
|Correction/gap proposal|handoff proposal|source owner / Settings / Studio|no source or runtime mutation|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-601..617|provenance assessment requested|CAP-INV-618|origins, selected versions, tools, decisions and handoffs|analysis capability|
|CAP-INV-618|gap or correction identified|source/Settings/Studio owner|missing lineage, impact, versions and return origin|provenance|
|CAP-INV-618|audit/evidence export authorized|Shared Export / Evidence owner|minimized package, permissions, restrictions and audit|provenance|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-601..617, Shared Trace/Activity/Versioning/Export/Recovery, Studio Tool Calls/Runs, Settings sources/connectors, Evidence, OPEN-008/012/013/014/015. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Provenance read/compare|cross-source and sensitive lineage disclosure|0|lineage and tenant masking|possible|viewer/auditor|source owners/Shared|Permissions|
|Provenance export/correction proposal|data leakage or history mutation|1/2|minimization/redaction|OPEN-013|requester/reviewer|Shared/Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- reproducibility assessment ≠ guaranteed replay.
- Tool result ≠ analyst conclusion.
- export permission ≠ source permission.
- supersession ≠ deletion.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- analyses with complete source/tool/decision lineage.
- reproducibility assessments by outcome.
- missing sensitive access audit or deleted provenance — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Tool version absente
**Given** une observation provient d’un Tool Call sans version  
**When** la reproductibilité est évaluée  
**Then** l’état devient conditional/not-reproducible et la lacune reste visible.

### 2. Accès expiré
**Given** une source n’est plus accessible au reviewer  
**When** le package est consulté  
**Then** les metadata et restrictions autorisées restent visibles sans réaccorder l’accès ni prétendre au replay.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** la provenance est auditée  
**Then** trace tables, version diff, checklists and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
