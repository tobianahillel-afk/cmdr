# Screen Register

## Règles

Une véritable spécification d'écran possède un objectif autonome, un Screen ID unique et le contrat de `04-experience-architecture/screen-section-contract.md`. Views, modes et filters n'obtiennent aucun Screen ID actif.

## Écrans actifs — 56

| ID | Produit | Module | Workspace | Statut | Source canonique |
|---|---|---|---|---|---|
| `CMD-CRP-001` | command | customer-and-reports | overview | draft | [`06-command/modules/customer-and-reports/screens/customer-overview.md`](../../06-command/modules/customer-and-reports/screens/customer-overview.md) |
| `CMD-EXC-001` | command | exposure-and-coverage | overview | draft | [`06-command/modules/exposure-and-coverage/screens/exposure-overview.md`](../../06-command/modules/exposure-and-coverage/screens/exposure-overview.md) |
| `CMD-IWQ-006` | command | incidents-and-work-queue | incident-detail | draft | [`06-command/modules/incidents-and-work-queue/screens/incident-detail.md`](../../06-command/modules/incidents-and-work-queue/screens/incident-detail.md) |
| `CMD-MC-001` | command | mission-control | now | draft | [`06-command/modules/mission-control/screens/mission-control-now.md`](../../06-command/modules/mission-control/screens/mission-control-now.md) |
| `CMD-MC-002` | command | mission-control | priorities | draft | [`06-command/modules/mission-control/screens/mission-control-priorities.md`](../../06-command/modules/mission-control/screens/mission-control-priorities.md) |
| `CMD-MC-003` | command | mission-control | decisions | draft | [`06-command/modules/mission-control/screens/mission-control-decisions.md`](../../06-command/modules/mission-control/screens/mission-control-decisions.md) |
| `CMD-MC-004` | command | mission-control | situation | draft | [`06-command/modules/mission-control/screens/mission-control-situation.md`](../../06-command/modules/mission-control/screens/mission-control-situation.md) |
| `CMD-MC-005` | command | mission-control | handover | draft | [`06-command/modules/mission-control/screens/mission-control-handover.md`](../../06-command/modules/mission-control/screens/mission-control-handover.md) |
| `CMD-RBI-001` | command | risk-and-business-impact | overview | draft | [`06-command/modules/risk-and-business-impact/screens/risk-overview.md`](../../06-command/modules/risk-and-business-impact/screens/risk-overview.md) |
| `CMD-RDO-001` | command | readiness-and-operations | overview | draft | [`06-command/modules/readiness-and-operations/screens/readiness-overview.md`](../../06-command/modules/readiness-and-operations/screens/readiness-overview.md) |
| `GOV-ACT-001` | govern | action-center | decision | draft | [`08-govern/modules/action-center/screens/action-center.md`](../../08-govern/modules/action-center/screens/action-center.md) |
| `GOV-AUD-001` | govern | audit-trail | audit | draft | [`08-govern/modules/audit-trail/screens/audit-trail.md`](../../08-govern/modules/audit-trail/screens/audit-trail.md) |
| `GOV-AUT-001` | govern | approvals-and-authorities | authorities | draft | [`08-govern/modules/approvals-and-authorities/screens/approvals-and-authorities.md`](../../08-govern/modules/approvals-and-authorities/screens/approvals-and-authorities.md) |
| `GOV-DEC-001` | govern | decision-register | register | draft | [`08-govern/modules/decision-register/screens/decision-register.md`](../../08-govern/modules/decision-register/screens/decision-register.md) |
| `GOV-INB-001` | govern | response-inbox | inbox | draft | [`08-govern/modules/response-inbox/screens/response-inbox.md`](../../08-govern/modules/response-inbox/screens/response-inbox.md) |
| `GOV-MET-001` | govern | response-metrics | metrics | draft | [`08-govern/modules/response-metrics/screens/response-metrics.md`](../../08-govern/modules/response-metrics/screens/response-metrics.md) |
| `GOV-PLB-001` | govern | playbooks | library | draft | [`08-govern/modules/playbooks/screens/playbooks.md`](../../08-govern/modules/playbooks/screens/playbooks.md) |
| `GOV-POL-001` | govern | policy-gates | policies | draft | [`08-govern/modules/policy-gates/screens/policy-gates.md`](../../08-govern/modules/policy-gates/screens/policy-gates.md) |
| `GOV-RUN-001` | govern | runs-and-rollback | runs | draft | [`08-govern/modules/runs-and-rollback/screens/runs-and-rollback.md`](../../08-govern/modules/runs-and-rollback/screens/runs-and-rollback.md) |
| `INV-CAS-001` | investigate | case-workspace | workspace | draft | [`07-investigate/modules/case-workspace/screens/case-workspace.md`](../../07-investigate/modules/case-workspace/screens/case-workspace.md) |
| `INV-DBG-001` | investigate | debugger | debug | draft | [`07-investigate/modules/debugger/screens/debugger.md`](../../07-investigate/modules/debugger/screens/debugger.md) |
| `INV-DSK-001` | investigate | disk-and-artifact-forensics | disk | draft | [`07-investigate/modules/disk-and-artifact-forensics/screens/disk-and-artifact-forensics.md`](../../07-investigate/modules/disk-and-artifact-forensics/screens/disk-and-artifact-forensics.md) |
| `INV-ENT-001` | investigate | entity-graph | graph | draft | [`07-investigate/modules/entity-graph/screens/entity-graph.md`](../../07-investigate/modules/entity-graph/screens/entity-graph.md) |
| `INV-EVD-001` | investigate | evidence | board | draft | [`07-investigate/modules/evidence/screens/evidence-board.md`](../../07-investigate/modules/evidence/screens/evidence-board.md) |
| `INV-EVS-001` | investigate | event-search | search | draft | [`07-investigate/modules/event-search/screens/event-search.md`](../../07-investigate/modules/event-search/screens/event-search.md) |
| `INV-HYP-001` | investigate | hypotheses-and-findings | reasoning | draft | [`07-investigate/modules/hypotheses-and-findings/screens/hypotheses-and-findings.md`](../../07-investigate/modules/hypotheses-and-findings/screens/hypotheses-and-findings.md) |
| `INV-MEM-001` | investigate | memory-forensics | memory | draft | [`07-investigate/modules/memory-forensics/screens/memory-forensics.md`](../../07-investigate/modules/memory-forensics/screens/memory-forensics.md) |
| `INV-REV-001` | investigate | reverse-engineering | reverse | draft | [`07-investigate/modules/reverse-engineering/screens/reverse-engineering.md`](../../07-investigate/modules/reverse-engineering/screens/reverse-engineering.md) |
| `INV-RPT-001` | investigate | investigation-report | report | draft | [`07-investigate/modules/investigation-report/screens/investigation-report.md`](../../07-investigate/modules/investigation-report/screens/investigation-report.md) |
| `INV-SBX-001` | investigate | dynamic-sandbox | sandbox | draft | [`07-investigate/modules/dynamic-sandbox/screens/dynamic-sandbox.md`](../../07-investigate/modules/dynamic-sandbox/screens/dynamic-sandbox.md) |
| `INV-STA-001` | investigate | static-analysis | analysis | draft | [`07-investigate/modules/static-analysis/screens/static-analysis.md`](../../07-investigate/modules/static-analysis/screens/static-analysis.md) |
| `INV-TIM-001` | investigate | case-timeline | timeline | draft | [`07-investigate/modules/case-timeline/screens/case-timeline.md`](../../07-investigate/modules/case-timeline/screens/case-timeline.md) |
| `INV-TRI-001` | investigate | intake-and-triage | triage-desk | draft | [`07-investigate/modules/intake-and-triage/screens/triage-desk.md`](../../07-investigate/modules/intake-and-triage/screens/triage-desk.md) |
| `SET-AUD-001` | platform-settings | administrative-audit | audit | draft | [`10-platform-settings/administrative-audit/screens/administrative-audit.md`](../../10-platform-settings/administrative-audit/screens/administrative-audit.md) |
| `SET-EAF-001` | platform-settings | endpoint-agent-fleet | fleet | draft | [`10-platform-settings/endpoint-agent-fleet/screens/endpoint-agent-fleet.md`](../../10-platform-settings/endpoint-agent-fleet/screens/endpoint-agent-fleet.md) |
| `SET-EPL-001` | platform-settings | endpoint-policies | policies | draft | [`10-platform-settings/endpoint-policies/screens/endpoint-policies.md`](../../10-platform-settings/endpoint-policies/screens/endpoint-policies.md) |
| `SET-HLT-001` | platform-settings | health | health | draft | [`10-platform-settings/health/screens/health.md`](../../10-platform-settings/health/screens/health.md) |
| `SET-IAM-001` | platform-settings | users-and-roles | users | draft | [`10-platform-settings/users-and-roles/screens/users-and-roles.md`](../../10-platform-settings/users-and-roles/screens/users-and-roles.md) |
| `SET-MDL-001` | platform-settings | models-and-providers | models | draft | [`10-platform-settings/models-and-providers/screens/models-and-providers.md`](../../10-platform-settings/models-and-providers/screens/models-and-providers.md) |
| `SET-PRF-001` | platform-settings | preferences | preferences | draft | [`10-platform-settings/preferences/screens/preferences.md`](../../10-platform-settings/preferences/screens/preferences.md) |
| `SET-RET-001` | platform-settings | retention | retention | draft | [`10-platform-settings/retention/screens/retention.md`](../../10-platform-settings/retention/screens/retention.md) |
| `SET-SBX-001` | platform-settings | sandbox-environments | sandboxes | draft | [`10-platform-settings/sandbox-environments/screens/sandbox-environments.md`](../../10-platform-settings/sandbox-environments/screens/sandbox-environments.md) |
| `SET-SEC-001` | platform-settings | secrets-and-connections | connections | draft | [`10-platform-settings/secrets-and-connections/screens/secrets-and-connections.md`](../../10-platform-settings/secrets-and-connections/screens/secrets-and-connections.md) |
| `SET-SRC-001` | platform-settings | sources-and-parsers | sources | draft | [`10-platform-settings/sources-and-parsers/screens/sources-and-parsers.md`](../../10-platform-settings/sources-and-parsers/screens/sources-and-parsers.md) |
| `SET-TEN-001` | platform-settings | tenants-and-environments | tenants | draft | [`10-platform-settings/tenants-and-environments/screens/tenants-and-environments.md`](../../10-platform-settings/tenants-and-environments/screens/tenants-and-environments.md) |
| `STD-AGT-001` | cmdr-studio | automation-agents | automation-agent | draft | [`09-cmdr-studio/automation-agents/screens/automation-agent-detail.md`](../../09-cmdr-studio/automation-agents/screens/automation-agent-detail.md) |
| `STD-ASR-001` | cmdr-studio | assurance | assurance | draft | [`09-cmdr-studio/assurance/screens/assurance.md`](../../09-cmdr-studio/assurance/screens/assurance.md) |
| `STD-ATM-001` | cmdr-studio | agent-teams | agent-team | draft | [`09-cmdr-studio/agent-teams/screens/agent-team-detail.md`](../../09-cmdr-studio/agent-teams/screens/agent-team-detail.md) |
| `STD-BLD-001` | cmdr-studio | builder | builder | draft | [`09-cmdr-studio/builder/screens/builder.md`](../../09-cmdr-studio/builder/screens/builder.md) |
| `STD-CTL-001` | cmdr-studio | control-room | control-room | draft | [`09-cmdr-studio/control-room/screens/control-room.md`](../../09-cmdr-studio/control-room/screens/control-room.md) |
| `STD-DEP-001` | cmdr-studio | versions-and-deployment | deployment | draft | [`09-cmdr-studio/versions-and-deployment/screens/deployment.md`](../../09-cmdr-studio/versions-and-deployment/screens/deployment.md) |
| `STD-EVL-001` | cmdr-studio | evaluations | evaluation | draft | [`09-cmdr-studio/evaluations/screens/evaluation-detail.md`](../../09-cmdr-studio/evaluations/screens/evaluation-detail.md) |
| `STD-LIB-001` | cmdr-studio | library | library | draft | [`09-cmdr-studio/library/screens/library.md`](../../09-cmdr-studio/library/screens/library.md) |
| `STD-SIM-001` | cmdr-studio | simulations | simulation | draft | [`09-cmdr-studio/simulations/screens/simulation-detail.md`](../../09-cmdr-studio/simulations/screens/simulation-detail.md) |
| `STD-SKL-001` | cmdr-studio | skills | skill | draft | [`09-cmdr-studio/skills/screens/skill-detail.md`](../../09-cmdr-studio/skills/screens/skill-detail.md) |
| `STD-WFL-001` | cmdr-studio | workflows | workflow | draft | [`09-cmdr-studio/workflows/screens/workflow-detail.md`](../../09-cmdr-studio/workflows/screens/workflow-detail.md) |

## Aliases de migration — non comptés comme écrans

| Ancien ID | Destination | Configuration | Statut |
|---|---|---|---|
| `CMD-IWQ-001` | Work Queue workspace | `view=incidents` | deprecated 2026-08-03 |
| `CMD-IWQ-002` | Work Queue workspace | `view=tasks` | deprecated 2026-08-03 |
| `CMD-IWQ-003` | Work Queue workspace | `view=unassigned` | deprecated 2026-08-03 |
| `CMD-IWQ-004` | Work Queue workspace | `view=sla-risk` | deprecated 2026-08-03 |
| `CMD-IWQ-005` | Work Queue workspace | `no system view; migration notice` | deprecated 2026-08-03 |

Ces aliases conservent les liens historiques ; ils ne sont pas des sources normatives et seront réévalués en Phase 6.

## Critère d’acceptation

**Given** l'ID `CMD-IWQ-003`, **When** le registre est consulté, **Then** il est identifié comme alias `view=unassigned`, pas comme écran actif, et `CMD-IWQ-006` reste le seul écran actif du module Work Queue.
