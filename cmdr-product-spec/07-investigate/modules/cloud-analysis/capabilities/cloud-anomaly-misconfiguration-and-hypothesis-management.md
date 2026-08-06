---
id: CAP-INV-615
title: Cloud Anomaly, Misconfiguration and Hypothesis Management
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
# CAP-INV-615 — Cloud Anomaly, Misconfiguration and Hypothesis Management

## 1. Définition
Créer et gérer des Cloud Anomaly candidates et Hypotheses à partir d’observations de configuration, identité, permission, ressource, activité et data access, avec supporting/contradicting elements, confidence, dispute, withdrawal and supersession.

## 2. Problème utilisateur
Un score d’anomalie ou une misconfiguration candidate peut être promu automatiquement en compromission, Finding ou vulnérabilité sans preuve suffisante.

## 3. Objectifs
- créer des anomaly candidates typées et sourcées.
- lier supporting, contradicting and missing elements, alternatives and confidence.
- gérer Hypothesis, dispute, withdrawal, review and supersession sans confirmation automatique.

## 4. Non-objectifs
Aucun moteur de scoring imposé, seuil universel, Finding confirmé, Incident, Alert, vulnerability validation, attribution ou réponse.

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
|Cloud observations|CAP-INV-604..614|resource, identity, permission, activity, config, network, storage and sensitive candidates|oui|selected versions|no hypothesis|
|Case/Hunt/Finding context|respective owner|questions, competing hypotheses and constraints|non|versioned|local analysis only|
|Source quality and limitations|sources / CAP-INV-601/603/607|coverage, gaps, errors and restrictions|oui|current session|confidence invalid|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Cloud Observation families|Investigate|supporting and contradicting projections|lecture|
|Case / Hunt / existing Hypothesis|respective owner|question and evidence context|lecture/lien|
|Graph / Timeline projections|Shared|relations and temporal context|lecture|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Cloud Anomaly|créer, qualifier, contester, retirer, versionner, superseder|Investigate|anomaly ≠ compromise|
|Cloud Hypothesis|créer, soutenir, contredire, dispute, withdraw, supersede|Investigate|Hypothesis ≠ Finding|
|Misconfiguration Candidate disposition|créer/modifier|Investigate|candidate ≠ vulnerability|

## 11. Fonctionnalités
- créer anomaly candidates par type.
- lier supporting/contradicting/missing elements and alternatives.
- documenter confidence and basis without opaque universal score.
- gérer analyst review, dispute, withdrawal and supersession.
- préparer Evidence/Finding handoff only after explicit human disposition.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Anomaly, Misconfiguration and Hypothesis Management|Cloud Security Analyst|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Cloud Security Analyst|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Cloud Security Analyst|concepts CAP-INV-615|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Anomaly, Misconfiguration and Hypothesis Management|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`candidate`, `under-analysis`, `supported`, `contradicted`, `inconclusive`, `disputed`, `withdrawn`, `superseded`, `handoff-ready`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Cloud Anomaly|versioned candidate|CAP-INV-616/617|not compromise or Finding|
|Cloud Hypothesis|versioned hypothesis|Case/Hunt / reviewer|support and contradictions visible|
|Disposition record|review event|Shared Trace / auditor|human rationale preserved|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-604..614|candidate pattern identified|CAP-INV-615|observations, sources, gaps and restrictions|source capability|
|CAP-INV-615|temporal validation needed|CAP-INV-616|hypothesis, relevant observations and time questions|anomaly management|
|CAP-INV-615|human disposition supports handoff|CAP-INV-617|candidate, evidence links, alternatives, confidence and limits|anomaly management|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-604..618, Cases/Hunts/Hypotheses, Shared Graph/Timeline/Trace, OPEN-012/013/014/015. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Anomaly/Hypothesis create/review/dispute|false accusation or escalation|2|sensitive observations masked|OPEN-013|author/reviewer|Investigate|Permissions|
|Cross-tenant correlation|scope leakage|1/2|tenant partition and minimization|mandatory by policy|analyst/access reviewer|Security/Investigate|Permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- configuration anomaly ≠ Finding.
- misconfiguration candidate ≠ confirmed vulnerability.
- cloud anomaly score ≠ compromise.
- Hypothesis ≠ analyst conclusion or Incident.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- anomalies/hypotheses by type and disposition.
- candidates with explicit contradictions and limitations.
- automatic compromise/Finding confirmation — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Score élevé
**Given** un outil retourne un score élevé sans evidence suffisante  
**When** une anomalie est créée  
**Then** le score reste Tool result, supporting/contradicting elements sont requis et aucune compromission n’est confirmée.

### 2. Hypothèse contredite
**Given** un événement ultérieur invalide une relation  
**When** la Hypothesis est revue  
**Then** elle devient contradicted/withdrawn ou superseded et l’historique reste visible.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** les anomalies sont gérées  
**Then** matrices pour/contre, rules, tables, timelines and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
