---
id: investigate-detection-engineering-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Source migration

## Audit result
No active standalone Detection Engineering lifecycle module, canonical Rule Builder or canonical promotion/deployment functional source existed before 4B.3A. Existing needs were distributed across Command runtime contracts, Endpoint Agent detection/update/suppression, Settings sources/environments/health, Govern decisions/runs, Studio Control Room and generic deployment/evaluation documentation.

## Disposition
- Competing Investigate functional sources deprecated: **0**.
- Active screens deprecated: **0**.
- Command Detection/Signal/Alert/Incident contracts remain active.
- Settings source, parser, environment, target and health contracts remain active.
- Endpoint runtime, update, suppression and rollback contracts remain active.
- Govern Action Request, Decision, Approval, Response Run, Result and rollback contracts remain active.
- Studio Tool, Workflow, Automation Run, Control Room and generic deployment/evaluation contracts remain active.
- Event Search, Query and Hunt remain active and distinct.
- Canonical Investigate source: `07-investigate/modules/detection-engineering/`, with lifecycle synthesis under `lifecycle-and-production/`.

## Acceptance
The canonical module absorbs Investigate review, readiness, lifecycle assessment and proposal needs without transferring runtime ownership, choosing a runtime/language, performing a production change or starting Threat Intelligence.
