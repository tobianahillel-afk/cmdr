---
id: portfolio-map
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-018
---
# Carte du portefeuille

| Produit | Question principale | Objets propres | Entrées | Sorties | Shared Capabilities clés |
|---|---|---|---|---|---|
| Command | Que faut-il traiter, par qui et avec quel impact ? | Incident, Task | Signals, Alerts, Results | Case context, Action Request context | Saved Views, Reporting, context |
| Investigate | Que s'est-il passé et que peut-on établir ? | Case, Hypothesis, Artifact, Evidence, Finding | Incident, Endpoint, telemetry | Finding, Action Request, detection improvement | Search, graph, timeline |
| Govern | Que peut-on autoriser et comment vérifier ? | Decision, Approval, Response Run, Result | Finding, Incident, Action Request | Result, rollback, audit | policy, audit, reporting |
| CMDR Studio | Comment construire et assurer l'automatisation ? | Skill, Tool, Tool Call, Agent, Workflow, Human Gate, Automation Run | product triggers/context | proposals, outputs, Action Request | models, metrics, audit |
| Platform Settings | Comment administrer la plateforme et la flotte ? | Tenant, Role, Fleet, Policy, Provider, Integration | admin intent | configuration and health | audit, notifications |
| Endpoint Agent | Que peut observer ou exécuter le composant local ? | Endpoint Agent local records | signed/authorized command, policy | telemetry, collections, results | trust, audit, context |

## Boucle principale

```text
Command → Investigate → Govern → Command
                 ↘ Studio support ↗
Settings → Endpoint Agent → Investigate / Govern
Shared Capabilities traversent les produits sans posséder leur décision métier.
```

## Critère d'acceptation

Le portefeuille ne crée ni produit « Platform » parallèle, ni produit « AI », ni produit « Shared Capabilities » visible comme destination métier autonome.
