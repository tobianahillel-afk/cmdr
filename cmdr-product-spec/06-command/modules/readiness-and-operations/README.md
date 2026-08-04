---
id: command-module-readiness-and-operations
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-005
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-PROD-057
open_decisions:
  - OPEN-010
  - OPEN-013
---
# Readiness and Operations

## Mission
Suivre la préparation opérationnelle globale, les exercices, plans, lacunes et actions d’amélioration sans dupliquer Studio Assurance, Sandbox ou un moteur d’exécution.

## Utilisateurs
Readiness Coordinator, Incident Commander, Team Lead, Business Owner, Platform/Product owners en projection.

## Capabilities
| ID | Capability | Status | Mode |
|---|---|---|---|
| CAP-CMD-301 | Readiness Overview | defined | planned |
| CAP-CMD-302 | Exercise Coordination | defined | planned |
| CAP-CMD-303 | Improvement Actions | defined | planned |
| CAP-CMD-304 | Operational Plans | defined | planned |
| CAP-CMD-305 | Capability Readiness | defined | planned |

## Frontières
Command coordonne records, Tasks, plans et assessments. Studio assure techniquement workflows/agents ; Settings fournit health/configuration ; Govern possède Decisions et Response Runs. Un exercice opérationnel n’est ni une sandbox, ni un test de règle, ni un Automation Run.

## Shared Capabilities consommées
Collaboration, Notifications, Reporting, Metrics, Object Linking, Versioning, Inspector, Context Bar, Trace et Audit Hooks sont consommés depuis leurs sources canoniques ; aucun moteur n’est dupliqué.

## Critère
**Given** une capability non testée, **When** Readiness est consultée, **Then** source, date, scope, owner, lacune et improvement Task sont visibles sans changer delivery mode ni statut documentaire.
