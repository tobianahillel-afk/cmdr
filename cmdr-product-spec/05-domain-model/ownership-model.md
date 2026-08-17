---
id: domain-ownership-model
domain: 05-domain-model
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Modèle de propriété

La propriété signifie: autorité sur le schéma, les invariants, les états et les changements de l’objet.

- Command possède Detection, Signal, Alert et Incident.
- Investigate possède Case, Evidence, Hypothesis et Finding.
- Govern possède Action Request, Decision, Response Run et Result.
- CMDR Studio possède Skill, Automation Agent, Agent Team, Workflow, Human Gate, Evaluation et Simulation.
- Platform Settings possède Endpoint Agent Fleet, Tenant, Environment, Integration, Parser, Model Provider, Secret Reference et Sandbox Environment.
- Shared Capabilities possède Telemetry Event, Report et les moteurs partagés.
- Endpoint Agent produit de la télémétrie et exécute des commandes mais ne possède pas la flotte administrative.

Le registre normatif est `../00-governance/ownership-register.md`.
