# Permission Catalog

## Objectif

Lister les identifiants de permissions utilisés par les produits CMDR. La sémantique d’autorisation, l’évaluation RBAC/ABAC, la séparation des tâches et les règles de refus restent définies dans [`permission-model.md`](permission-model.md).

## Convention

`perm.<domaine>.<ressource>.<action>` lorsque quatre segments sont utiles, ou `perm.<domaine>.<action>` pour une permission transversale. Une permission absente de ce catalogue ne doit pas être déployée sans mise à jour du registre et revue de sécurité.

| Permission | Propriétaire fonctionnel | Description normative minimale |
|---|---|---|
| `perm.cmdr-studio.agent-team.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.agent-team.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.automation-agent.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.automation-agent.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.deployment.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.deployment.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.evaluation.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.evaluation.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.human-gate.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.human-gate.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.simulation.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.simulation.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.skill.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.skill.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.version.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.version.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.workflow.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.cmdr-studio.workflow.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.alert.manage` | Command | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.alert.read` | Command | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.coordinate` | Command | Autorise l’action « coordinate » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.detection.manage` | Command | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.detection.read` | Command | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.incident.manage` | Command | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.incident.read` | Command | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.read` | Command | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.saved-view.share` | Command | Autorise l’action « share » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.signal.manage` | Command | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.signal.read` | Command | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.command.task.manage` | Command | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.design-system.component.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.design-system.component.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.agent-command.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.agent-command.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.endpoint-agent.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.endpoint-agent.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.local-audit-event.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.endpoint-agent.local-audit-event.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.action-request.create` | Govern | Autorise l’action « create » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.action-request.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.action-request.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.approval.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.approval.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.authority.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.authority.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.decision.approve` | Govern | Autorise l’action « approve » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.decision.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.decision.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.metrics.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.playbook.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.playbook.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.policy.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.policy.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-rollback.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-rollback.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-run.execute` | Govern | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-run.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-run.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-run.rollback` | Govern | Autorise l’action « rollback » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-step.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.response-step.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.result.manage` | Govern | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.govern.result.read` | Govern | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.analysis.execute` | Investigate | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.artifact.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.artifact.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.case.create` | Investigate | Autorise l’action « create » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.case.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.case.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.collection-request.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.collection-request.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.debugger.execute` | Investigate | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.disk-image.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.disk-image.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.entity.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.evidence.collect` | Investigate | Autorise l’action « collect » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.evidence.export` | Investigate | Autorise l’action « export » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.evidence.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.evidence.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.evidence.verify` | Investigate | Autorise l’action « verify » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.finding.approve` | Investigate | Autorise l’action « approve » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.finding.author` | Investigate | Autorise l’action « author » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.finding.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.finding.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.hypothesis.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.hypothesis.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.malware-sample.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.malware-sample.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.memory-image.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.memory-image.read` | Investigate | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.report.manage` | Investigate | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.reverse.execute` | Investigate | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.sandbox.execute` | Investigate | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.investigate.search.execute` | Investigate | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.data-source.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.data-source.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.endpoint-agent-fleet.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.endpoint-agent-fleet.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.endpoint-policy.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.endpoint-policy.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.environment.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.environment.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.integration.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.integration.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.model-provider.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.model-provider.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.parser.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.parser.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.principal.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.principal.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.role.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.role.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.sandbox-environment.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.sandbox-environment.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.secret-reference.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.secret-reference.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.tenant.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.platform-settings.tenant.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.product-architecture.capability.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.product-architecture.capability.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.product-architecture.screen.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.product-architecture.screen.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.security.audit.export` | Security Architecture | Autorise l’action « export » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.security.audit.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.endpoint-fleet.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.endpoint-fleet.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.endpoint-policy.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.endpoint-policy.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.health.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.identity.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.identity.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.model.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.model.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.preferences.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.retention.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.retention.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.sandbox.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.sandbox.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.secret.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.secret.read-metadata` | Platform Settings | Autorise l’action « read metadata » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.source.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.source.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.tenant.manage` | Platform Settings | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.settings.tenant.read` | Platform Settings | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.entity.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.entity.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.notification.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.notification.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.query.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.query.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.saved-view.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.saved-view.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.search-job.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.search-job.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.task.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.task.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.telemetry-event.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.telemetry-event.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.timeline-entry.manage` | Security Architecture | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared-capabilities.timeline-entry.read` | Security Architecture | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared.report.create` | Shared Capabilities | Autorise l’action « create » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared.report.export` | Shared Capabilities | Autorise l’action « export » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared.report.publish` | Shared Capabilities | Autorise l’action « publish » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared.report.read` | Shared Capabilities | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.shared.report.review` | Shared Capabilities | Autorise l’action « review » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.agent-team.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.agent-team.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.artifact.create` | CMDR Studio | Autorise l’action « create » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.artifact.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.assurance.approve` | CMDR Studio | Autorise l’action « approve » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.assurance.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.automation-agent.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.automation-agent.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.builder.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.control-room.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.deployment.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.deployment.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.evaluation.execute` | CMDR Studio | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.evaluation.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.library.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.simulation.execute` | CMDR Studio | Autorise l’action « execute » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.simulation.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.skill.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.skill.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.workflow.manage` | CMDR Studio | Autorise l’action « manage » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
| `perm.studio.workflow.read` | CMDR Studio | Autorise l’action « read » dans le périmètre RBAC/ABAC applicable ; aucun droit implicite d’export, d’approbation ou d’exécution n’est accordé. |
