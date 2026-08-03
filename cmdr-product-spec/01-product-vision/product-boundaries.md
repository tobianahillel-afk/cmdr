---
id: product-boundaries
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-015
  - REQ-PROD-016
  - REQ-PROD-017
  - REQ-PROD-018
  - REQ-AI-002
  - REQ-OBJ-001
  - REQ-OBJ-012
---
# Frontières des produits

## Règle transversale

Un produit possède une mission, des capabilities, des objets et des workflows. Il peut consommer les projections d'un autre produit, mais ne peut pas devenir propriétaire par simple affichage ou intégration.

## Command

### Mission
Prioriser, coordonner, distribuer, superviser et maintenir la situation opérationnelle.

### Objets possédés
Incident ; Task opérationnelle ; états de coordination et handover.

### Capabilities possédées
Mission Control, Work Queue, affectation, priorité, SLA, situation, impact métier, readiness, handover et customer delivery lorsque applicable.

### Capabilities consommées
Signals et Alerts, projections de Case/Evidence/Finding, Decision, Response Run, Result, Reporting Engine, Saved Views génériques.

### Workspaces principaux
Mission Control ; Incidents and Work Queue ; Risk and Coverage ; Readiness and Operations ; Customers and Delivery selon `OPEN-006`.

### Exclusions
Pas de forensic, reverse engineering, terminal endpoint, éditeur de règles, policy engine détaillé ou builder d'agents.

### Transitions
Entre : Signal/Alert → Incident, Result → Incident, handover.  
Sort : Incident → Case, Incident/Context → Action Request ou Govern.

### IA
Assistance optionnelle pour résumé, priorité proposée et handover ; Command reste propriétaire de la situation et de l'Incident.

## Investigate

### Mission
Rechercher, collecter, tester des hypothèses, analyser, établir des Evidence et produire des Findings.

### Objets possédés
Case, Hypothesis, Artifact, Evidence, Finding, collections spécialisées.

### Capabilities possédées
Signals and Hunt du point de vue investigation, Case Workspace, collection et Live Response analytique, workbench technique, forensic, static analysis, reverse engineering, debugger, sandbox, Detection Engineering et Intelligence.

### Capabilities consommées
Incident, Endpoint/Fleet projection, policies et decisions Govern, Reporting Engine, Studio automation.

### Workspaces principaux
Signals and Hunt ; Cases and Evidence ; Collection and Live Response ; Analysis Workbench ; Detection Engineering ; Intelligence.

### Exclusions
Pas d'autorité indépendante pour une réponse à risque élevé, pas d'administration de flotte, pas de seconde Work Queue générale.

### Transitions
Entre : Incident → Case, Endpoint → Collection.  
Sort : Evidence → Finding, Finding → Action Request, Finding → Detection improvement.

### IA
Peut proposer hypothèse, requête ou analyse, mais ne confirme pas seule un Finding et ne contourne pas Govern.

## Govern

### Mission
Recevoir les demandes, appliquer policies et autorités, décider, exécuter ou autoriser, observer, vérifier et rollback.

### Objets possédés
Decision, Approval, Response Run, Result, Policy, Exception et records d'audit de réponse. Action Request est produit par Command/Investigate puis gouverné ici.

### Capabilities possédées
Review Queue spécialisée, Decision Workspace, Response Runs, Policies and Exceptions, Audit Ledger.

### Capabilities consommées
Incident, Case, Evidence, Finding, Endpoint, Studio Human Gates et automatisations autorisées.

### Workspaces principaux
Review Queue ; Decision Workspace ; Run Workspace ; Policies and Exceptions ; Audit Ledger.

### Exclusions
Pas de Work Queue générale, pas de Case Workspace, pas de création de Evidence ni de builder d'agents.

### Transitions
Entre : Finding/Incident → Action Request.  
Sort : Decision → Response Run → Result → Command/Investigate.

### IA
Peut résumer ou vérifier la complétude ; l'autorité et la Decision restent humaines ou fondées sur une policy explicitement approuvée.

## CMDR Studio

### Mission
Concevoir, versionner, évaluer, déployer et superviser les automatisations et capacités agentiques.

### Objets possédés
Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate, Automation Run, Evaluation, Simulation, Version et Deployment.

### Capabilities possédées
Library, Builder, Assurance, Control Room, coûts/usage, observabilité et deployment des automatisations.

### Capabilities consommées
Objets et contexts des produits opérationnels sous permission ; policies de Govern ; models/providers de Settings.

### Workspaces principaux
Library ; Builder ; Control Room ; Assurance.

### Exclusions
Pas de propriété d'Incident, Case, Finding, Decision ou Response Run ; pas d'interface opérationnelle principale.

### Transitions
Entre : trigger produit → Automation Run ; Human Gate → validation humaine.  
Sort : résultat proposé vers produit source ; action risquée → Action Request/Govern.

### IA
Studio est le propriétaire des capacités agentiques, pas des décisions métier.

## Platform Settings

### Mission
Administrer identité, tenants, environnements, secrets, sources, providers, flotte, policies endpoint, rétention, santé et audit administratif.

### Objets possédés
Tenant, Environment, Principal, Role, Integration, Data Source, Parser, Model Provider, Secret Reference, Endpoint Agent Fleet, Endpoint Policy et Sandbox Environment administratif.

### Capabilities consommées
Permission Model, audit, health, notification et reporting administratif.

### Workspaces principaux
Users and Roles ; Tenants and Environments ; Endpoint Agent Fleet ; Endpoint Policies ; Sources and Parsers ; Models and Providers ; Secrets and Connections ; Sandbox Environments ; Retention ; Health ; Administrative Audit ; Preferences.

### Exclusions
Pas d'investigation, Case, forensic, décision de réponse ou orchestration métier.

### Transitions
Fournit contexte de plateforme aux autres produits ; administre enrollment, policy et version de l'Endpoint Agent.

### IA
Administre les providers et policies, mais ne définit pas les workflows agentiques.

## Endpoint Agent

### Mission
Fournir progressivement les capacités locales de télémétrie, détection, inspection, collecte, Live Response, containment, résilience et sécurité.

### Statut de delivery
Cible native complète ; classification actuelle `planned` tant qu'aucune preuve de release et de moteur principal n'est liée.

### Objets et responsabilité
Le composant Endpoint Agent exécute localement ; Platform Settings possède la flotte et les policies ; Investigate possède le workflow analytique ; Govern possède l'autorité de réponse.

### Capabilities consommées
Command signing, permissions, policy, audit, context, actions Govern.

### Exclusions
Pas de workflow métier, pas de Decision, pas de Case, pas de propriété de la flotte.

### Transitions
Télémétrie vers détection/Signal ; collecte vers Artifact/Evidence ; action autorisée vers Result.

### IA
Aucune dépendance à un modèle pour les fonctions essentielles ; l'inférence locale future reste classifiée et gouvernée.

## Shared Capabilities

Shared Capabilities possède les services communs tels que Reporting Engine et Saved Views génériques. Il ne devient pas un produit parallèle et ne prend pas les décisions métier des produits consommateurs.

## Matrice transversale

| Capability / objet | Command | Investigate | Govern | Studio | Settings | Endpoint |
|---|---|---|---|---|---|---|
| Incident coordination | Owner | Consumer | Consumer | No | metadata | No |
| Case / Evidence / Finding | Consumer | Owner | Consumer | automation support | No | collection source |
| Decision / Approval | context provider | requester | Owner | Human Gate support | policy admin | executor context |
| Response Run / Result | Consumer | evidence/verification | Owner | workflow support | No | execution |
| Automation objects | trigger/consumer | trigger/consumer | policy/consumer | Owner | provider admin | Tool target possible |
| Endpoint Agent Fleet | Consumer | Consumer | Consumer | No | Owner | managed component |
| Reporting Engine | Consumer | Consumer | Consumer | Consumer | Consumer | No |
| Generic Saved Views | Consumer | Consumer | Consumer | Consumer | Consumer | No |
| Work Queue Saved Views | Owner | No | No | No | No | No |

## Comportements interdits

- Command ne devient pas un workbench forensic.
- Investigate ne devient pas l'autorité de décision de réponse.
- Govern ne devient pas une Work Queue générale.
- Studio ne remplace pas les produits opérationnels.
- Settings ne devient pas un espace d'investigation.
- Endpoint Agent ne possède pas le workflow métier.
- Shared Capabilities ne devient pas un septième produit.

## Critères d'acceptation

**Given** `Endpoint Agent Fleet`,  
**When** la matrice est consultée,  
**Then** Settings est Owner, Investigate est Consumer, Govern contrôle les actions risquées et Endpoint est le composant administré.

**Given** une Evidence affichée dans Command,  
**When** l'utilisateur la sélectionne,  
**Then** Command utilise une projection de la source Investigate et ne crée ni schéma ni lifecycle concurrent.
