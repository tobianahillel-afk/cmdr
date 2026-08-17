---
id: object-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: registry
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Object Register

## Objectif

Répertorier les objets canoniques CMDR, leur produit propriétaire et leur source unique. Une page, un écran ou une capability peut utiliser une projection, mais ne peut pas redéfinir l’objet ou son cycle de vie.

## Propriétés obligatoires

- Incident, Alert, Signal et Task opérationnelle → Command.
- Case, Artifact, Evidence et Finding → Investigate.
- Action Request, Decision, Approval, Response Run et Result → Govern.
- Endpoint Agent Fleet → Platform Settings.
- Skill et Automation Agent → CMDR Studio.
- Report reste défini par le [`Reporting Engine`](../../12-shared-capabilities/reporting-engine.md), pas par un objet de console.
- Business Service Catalog fournit le contexte Service tant qu’un objet `service.md` n’est pas formalisé.
- Une projection ne transfère jamais l’ownership.

| Objet | Propriétaire | Statut | Source canonique |
|---|---|---|---|
| `action-request` | Govern | draft | [`05-domain-model/objects/action-request.md`](../../05-domain-model/objects/action-request.md) |
| `agent-command` | Endpoint Agent | draft | [`05-domain-model/objects/agent-command.md`](../../05-domain-model/objects/agent-command.md) |
| `agent-team` | CMDR Studio | draft | [`05-domain-model/objects/agent-team.md`](../../05-domain-model/objects/agent-team.md) |
| `alert` | Command | draft | [`05-domain-model/objects/alert.md`](../../05-domain-model/objects/alert.md) |
| `approval` | Govern | draft | [`05-domain-model/objects/approval.md`](../../05-domain-model/objects/approval.md) |
| `artifact` | Investigate | draft | [`05-domain-model/objects/artifact.md`](../../05-domain-model/objects/artifact.md) |
| `automation-agent` | CMDR Studio | draft | [`05-domain-model/objects/automation-agent.md`](../../05-domain-model/objects/automation-agent.md) |
| `capability` | Product Architecture | draft | [`05-domain-model/objects/capability.md`](../../05-domain-model/objects/capability.md) |
| `case` | Investigate | draft | [`05-domain-model/objects/case.md`](../../05-domain-model/objects/case.md) |
| `collection-request` | Investigate | draft | [`05-domain-model/objects/collection-request.md`](../../05-domain-model/objects/collection-request.md) |
| `component` | Design System | draft | [`05-domain-model/objects/component.md`](../../05-domain-model/objects/component.md) |
| `data-source` | Platform Settings | draft | [`05-domain-model/objects/data-source.md`](../../05-domain-model/objects/data-source.md) |
| `decision` | Govern | draft | [`05-domain-model/objects/decision.md`](../../05-domain-model/objects/decision.md) |
| `deployment` | CMDR Studio | draft | [`05-domain-model/objects/deployment.md`](../../05-domain-model/objects/deployment.md) |
| `detection` | Command | draft | [`05-domain-model/objects/detection.md`](../../05-domain-model/objects/detection.md) |
| `disk-image` | Investigate | draft | [`05-domain-model/objects/disk-image.md`](../../05-domain-model/objects/disk-image.md) |
| `endpoint-agent` | Endpoint Agent | draft | [`05-domain-model/objects/endpoint-agent.md`](../../05-domain-model/objects/endpoint-agent.md) |
| `endpoint-agent-fleet` | Platform Settings | draft | [`05-domain-model/objects/endpoint-agent-fleet.md`](../../05-domain-model/objects/endpoint-agent-fleet.md) |
| `endpoint-policy` | Platform Settings | draft | [`05-domain-model/objects/endpoint-policy.md`](../../05-domain-model/objects/endpoint-policy.md) |
| `entity` | Shared Capabilities | draft | [`05-domain-model/objects/entity.md`](../../05-domain-model/objects/entity.md) |
| `environment` | Platform Settings | draft | [`05-domain-model/objects/environment.md`](../../05-domain-model/objects/environment.md) |
| `evaluation` | CMDR Studio | draft | [`05-domain-model/objects/evaluation.md`](../../05-domain-model/objects/evaluation.md) |
| `evidence` | Investigate | draft | [`05-domain-model/objects/evidence.md`](../../05-domain-model/objects/evidence.md) |
| `finding` | Investigate | draft | [`05-domain-model/objects/finding.md`](../../05-domain-model/objects/finding.md) |
| `human-gate` | CMDR Studio | draft | [`05-domain-model/objects/human-gate.md`](../../05-domain-model/objects/human-gate.md) |
| `hypothesis` | Investigate | draft | [`05-domain-model/objects/hypothesis.md`](../../05-domain-model/objects/hypothesis.md) |
| `incident` | Command | draft | [`05-domain-model/objects/incident.md`](../../05-domain-model/objects/incident.md) |
| `integration` | Platform Settings | draft | [`05-domain-model/objects/integration.md`](../../05-domain-model/objects/integration.md) |
| `local-audit-event` | Endpoint Agent | draft | [`05-domain-model/objects/local-audit-event.md`](../../05-domain-model/objects/local-audit-event.md) |
| `malware-sample` | Investigate | draft | [`05-domain-model/objects/malware-sample.md`](../../05-domain-model/objects/malware-sample.md) |
| `memory-image` | Investigate | draft | [`05-domain-model/objects/memory-image.md`](../../05-domain-model/objects/memory-image.md) |
| `model-provider` | Platform Settings | draft | [`05-domain-model/objects/model-provider.md`](../../05-domain-model/objects/model-provider.md) |
| `notification` | Shared Capabilities | draft | [`05-domain-model/objects/notification.md`](../../05-domain-model/objects/notification.md) |
| `parser` | Platform Settings | draft | [`05-domain-model/objects/parser.md`](../../05-domain-model/objects/parser.md) |
| `playbook` | Govern | draft | [`05-domain-model/objects/playbook.md`](../../05-domain-model/objects/playbook.md) |
| `policy` | Govern | draft | [`05-domain-model/objects/policy.md`](../../05-domain-model/objects/policy.md) |
| `principal` | Platform Settings | draft | [`05-domain-model/objects/principal.md`](../../05-domain-model/objects/principal.md) |
| `query` | Shared Capabilities | draft | [`05-domain-model/objects/query.md`](../../05-domain-model/objects/query.md) |
| `response-rollback` | Govern | draft | [`05-domain-model/objects/response-rollback.md`](../../05-domain-model/objects/response-rollback.md) |
| `response-run` | Govern | draft | [`05-domain-model/objects/response-run.md`](../../05-domain-model/objects/response-run.md) |
| `response-step` | Govern | draft | [`05-domain-model/objects/response-step.md`](../../05-domain-model/objects/response-step.md) |
| `result` | Govern | draft | [`05-domain-model/objects/result.md`](../../05-domain-model/objects/result.md) |
| `role` | Platform Settings | draft | [`05-domain-model/objects/role.md`](../../05-domain-model/objects/role.md) |
| `sandbox-environment` | Platform Settings | draft | [`05-domain-model/objects/sandbox-environment.md`](../../05-domain-model/objects/sandbox-environment.md) |
| `saved-view` | Shared Capabilities | draft | [`05-domain-model/objects/saved-view.md`](../../05-domain-model/objects/saved-view.md) |
| `screen` | Product Architecture | draft | [`05-domain-model/objects/screen.md`](../../05-domain-model/objects/screen.md) |
| `search-job` | Shared Capabilities | draft | [`05-domain-model/objects/search-job.md`](../../05-domain-model/objects/search-job.md) |
| `secret-reference` | Platform Settings | draft | [`05-domain-model/objects/secret-reference.md`](../../05-domain-model/objects/secret-reference.md) |
| `signal` | Command | draft | [`05-domain-model/objects/signal.md`](../../05-domain-model/objects/signal.md) |
| `simulation` | CMDR Studio | draft | [`05-domain-model/objects/simulation.md`](../../05-domain-model/objects/simulation.md) |
| `skill` | CMDR Studio | draft | [`05-domain-model/objects/skill.md`](../../05-domain-model/objects/skill.md) |
| `task` | Command | draft | [`05-domain-model/objects/task.md`](../../05-domain-model/objects/task.md) |
| `telemetry-event` | Shared Capabilities | draft | [`05-domain-model/objects/telemetry-event.md`](../../05-domain-model/objects/telemetry-event.md) |
| `tenant` | Platform Settings | draft | [`05-domain-model/objects/tenant.md`](../../05-domain-model/objects/tenant.md) |
| `timeline-entry` | Shared Capabilities | draft | [`05-domain-model/objects/timeline-entry.md`](../../05-domain-model/objects/timeline-entry.md) |
| `version` | CMDR Studio | draft | [`05-domain-model/objects/version.md`](../../05-domain-model/objects/version.md) |
| `workflow` | CMDR Studio | draft | [`05-domain-model/objects/workflow.md`](../../05-domain-model/objects/workflow.md) |

## Alignement Phase 4A

`Task` est alignée sur Command conformément aux décisions sources, à l’ownership register, au permission catalog et à la Work Queue. `12-shared-capabilities/task-inbox.md` reste une capability partagée de présentation/agrégation ; elle n’est pas un owner concurrent.

Les objets `service`, `exposure`, `report` et `audit-record` demandés comme sources Phase 4A n’existent pas sous ces chemins. La Phase 4A enregistre les besoins fonctionnels et utilise les sources partagées existantes sans créer prématurément de schéma.
