---
id: CAP-INV-614
title: Cloud Secrets, Keys and Sensitive Material Assessment
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
# CAP-INV-614 — Cloud Secrets, Keys and Sensitive Material Assessment

## 1. Définition
Évaluer la présence candidate et les metadata de keys, tokens, credentials, certificates and secret references avec niveaux d’accès distincts, sans utiliser un secret ni révéler sa valeur sans permission explicite.

## 2. Problème utilisateur
La simple détection d’une référence ou d’un nom sensible peut entraîner une exposition inutile, une validation du credential ou son utilisation pendant l’investigation.

## 3. Objectifs
- distinguer présence candidate, metadata, aperçu masqué, reveal, copie, export and prohibited use.
- lier source, workload, access events, expiration, restrictions and audit.
- préparer Sensitive Material Candidate et Evidence Candidate context sans tester la validité.

## 4. Non-objectifs
Aucune validation de credential, authentification, rotation, révocation, secret use, brute force, copie/export implicite ou modification de vault/provider.

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
|Secret/key/reference candidates|CAP-INV-608..613 / Artifact sources|candidate presence and metadata|oui|observed time|no assessment|
|Access policy and classification|Security / source owner|reveal, copy, export and handling authority|oui|current decision|masked only|
|Related workload and events|CAP-INV-609..613/607|source, workload, identity and access context|non|session versions|context partial|

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
|Sensitive Material Candidate source projection|source owner|presence and permitted metadata|lecture limitée|
|Workload / Storage / Activity Observation|Investigate|location and access context|lecture/lien|
|Secret Reference / policy projection|Platform Settings / Security|reference and authority, never value by default|lecture limitée|

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
|Sensitive Material Candidate|créer, qualifier, contester, versionner, superseder|Investigate|candidate ≠ valid credential|
|Sensitive access record|créer/lier|Shared Trace / source owner|every reveal/copy/export attempt audited|
|Evidence Candidate context|préparer|Evidence owner|no value use or qualification|

## 11. Fonctionnalités
- enregistrer candidate presence and type: key, token, credential, certificate or reference.
- afficher metadata and masked preview according to policy.
- séparer reveal, copy and export permissions.
- lier source, workload, access events, expiration and restrictions.
- préparer Evidence Candidate or revocation Action Request context without use.
- préserver source, observation time, versions, restrictions, uncertainty et return origin.
- fonctionner intégralement sans IA, provider ou intégration obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
|consulter, filtrer, rechercher et comparer Cloud Secrets, Keys and Sensitive Material Assessment|Evidence Reviewer|sources et projections autorisées|0|scope, tenant et lecture autorisés|vue sourcée, permission-aware et temporelle|non|
|exécuter une analyse, corrélation, reconstruction ou génération bornée|Evidence Reviewer|résultat analytique / Tool Call|1|sources, paramètres et restrictions visibles|résultat attribué avec erreurs, limites et incertitude|selon politique|
|créer, annoter, contester, versionner, retirer ou préparer un handoff|Evidence Reviewer|concepts CAP-INV-614|2|owner, provenance, permissions et séparation explicites|mutation analytique réversible ou proposition|OPEN-013|
|modifier une permission, une ressource, un credential, un réseau ou un runtime Cloud|aucun rôle local|cible réelle|3|hors périmètre Investigate|Action Request ou handoff seulement|obligatoire|
|détruire une ressource, une preuve, une trace ou divulguer irréversiblement|aucun rôle local|cible/provenance|4|interdit localement|refus audité|strict|

Investigate exécute localement uniquement les classes 0 à 2. Toute classe 3 ou 4 est refusée ou préparée comme Action Request/handoff vers Govern et le propriétaire de la cible.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
|préparer ou structurer Cloud Secrets, Keys and Sensitive Material Assessment|oui|formulaires, catalogues et règles explicables|oui|suggestion ou draft sourcé|formulaire structuré, checklist et revue humaine|
|comparer sources, versions, relations ou alternatives|oui|diff, tables, graphes et agrégations|oui|regroupement avec incertitude|comparateur, filtres et revue humaine|
|résumer observations, gaps, contradictions et limites|oui|vues sourcées et templates|oui|résumé attribué|timeline, matrice et Inspector|
|confirmer une identité, permission, compromission, Evidence, Finding ou action|humain/destination owner|contrôles seulement|jamais autonome|jamais décisionnaire|revue humaine et Govern|

Toute sortie automatisée expose initiateur, Tool/agent/model et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et accept/modify/reject. Aucun secret n’est envoyé ou révélé sans permission explicite.

## 14. États fonctionnels
`candidate`, `metadata-only`, `masked`, `reveal-requested`, `revealed-authorized`, `copy-blocked`, `export-blocked`, `disputed`, `expired-observed`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas des machines d’état physiques finales.

## 15. États d’interface
Loading conserve Session, scope et source ; Empty distingue absence, non-collecte et interdiction ; Partial nomme coverage gaps et champs manquants ; Error conserve les résultats valides ; Offline reste stale/read-only ; Permission denied ne révèle ni existence sensible ni valeur ; Stale expose dates et versions ; Conflict fournit diff, disposition et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
|Sensitive Material Candidate|versioned candidate|CAP-INV-615/617|validity not tested|
|Sensitive access audit record|trace event|auditor / source owner|actor, purpose, level and result retained|
|Evidence/response preparation context|package|Evidence/Govern owner|no secret use or revocation executed|

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
|CAP-INV-608..613|sensitive candidate identified|CAP-INV-614|source, type, location, metadata and restrictions|source analysis|
|CAP-INV-614|evidence context required|CAP-INV-617/Evidence owner|masked candidate, provenance, access history and limits|sensitive assessment|
|CAP-INV-614|revocation need suspected|Govern Action Request|candidate reference, risk, authority and return origin|sensitive assessment|

Chaque transition conserve tenant, environment, scope units, source owner, versions, permissions, restrictions, sensitive markings, errors, authority, provenance et return origin. Elle ne crée aucune permission ni mutation de cible.

## 18. Dépendances
CAP-INV-607..618, Settings Secrets/Connections, Security privacy/secrets/export, Evidence, Govern, Shared Trace, OPEN-008/012/013/014. Shared est consommé sans redéfinition. Aucun provider, API, protocol, query language, engine, schema ou connector n’est choisi.

## 19. Source de vérité
Investigate est source des observations, assessments, Sessions, candidates, Hypotheses et packages Cloud locaux. Platform Settings reste source de la configuration administrative et de l’état des integrations; les sources Cloud restent sources de leurs records; Command, Govern, Studio, Shared et les autres modules gardent leurs objets. Une projection, normalisation, correlation ou Tool result ne transfère jamais ownership ou permissions.

## 20. Provenance et audit
Conserver origin Case/Incident/Finding/Hunt/Signal ou source/Artifact, organization, tenant, account, subscription, project, provider-neutral service/resource references, source and connector projections, collection context, permissions, restrictions, Session, Tools, Tool Calls, Automation Runs, functional queries, observations, Hypotheses, Artifacts, sensitive access records, timelines, correlations, handoffs, errors, limitations, human decisions, versions, timestamps et return origin. Toute correction utilise annotation, nouvelle version, retrait ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
|Secret presence/metadata/masked preview|credential disclosure|0|mandatory masking|possible|viewer/source owner|Settings/Security|Permissions|
|Reveal/copy/export|high-impact sensitive disclosure|2/3|purpose-bound minimization|mandatory|requester/reviewer/approver|source owner/Govern|Future permissions|

Les permissions atomiques, namespaces, RBAC/ABAC, step-up définitif et séparation finale restent futurs. Permission sur Session ou package ne vaut jamais permission sur source, secret, contenu ou cible.

## 22. Limites et erreurs
- secret candidate ≠ valid credential.
- secret existence ≠ permission to reveal or use it.
- masked preview ≠ reveal.
- reveal permission ≠ copy/export/use permission.
- source unavailable, parser/schema unsupported, tenant mismatch, permission denied, stale data, timeout, cancellation, duplicate, late event, conflict et superseded version restent visibles.
- une absence ou un no-match ne prouve ni absence d’activité ni absence d’attaque.

## 23. Métriques conceptuelles
- candidates by access level and disposition.
- reveal/copy/export attempts with audit.
- secret use by Investigate — target zero.
- outputs automatisés avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, secret utilisé, target modifiée ou trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou SLA non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned`; preuve documentaire fonctionnelle uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Aucun code, API, provider, connector, command, runtime ou configuration réelle n’est livré.

## 25. Critères d’acceptation
### 1. Candidate non valide
**Given** un pattern ressemble à un token  
**When** il est évalué  
**Then** il reste candidate, sa valeur n’est pas testée et les limites sont visibles.

### 2. Reveal refusé
**Given** un analyste sans droit demande la valeur  
**When** l’action est tentée  
**Then** la valeur n’est pas révélée, le refus est audité et aucune copie/export n’est possible.

### 3. Sans IA
**Given** aucun modèle n’est disponible  
**When** le matériel sensible est géré  
**Then** pattern rules, metadata viewers, masking and human review suffice.

## 26. Questions ouvertes
OPEN-012 reste ouverte pour le scope Cloud : providers prioritaires, unités de compte, services, multi-cloud, audit models, inventory, cloud-native workloads, orchestration, serverless, SaaS, Cloud evidence and cross-tenant limits. OPEN-008/013/014/015 restent ouvertes pour support des sources, actions réversibles, relations d’objets et provenance d’automatisation. Aucun provider, protocole, langage de requête, schéma, connecteur ou politique finale n’est sélectionné.

## 27. Consommateurs documentaires
Cloud Analysis, Investigate, Cases/Hunts/Evidence/Findings, Collection and Live Response, Analysis Workbench, Network/Memory/Disk Forensics, Detection Engineering, Threat Intelligence, Command, Platform Settings, CMDR Studio, Govern, Shared Capabilities, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun contenu ou capability Mobile Forensics n’est créé.
