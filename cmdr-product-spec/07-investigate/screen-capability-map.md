---
id: investigate-screen-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-UX-001, REQ-UX-010, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-012]
---
# Screen capability map — Threat Intelligence and Cloud Analysis functional scope

| Existing screen or surface | Consumption | Change through 4B.4A |
|---|---|---|
| Entity Graph / Graph / Timeline | Intelligence candidates plus Cloud scope, identity, permission, resource, network and correlation candidates | links/map only; owners unchanged |
| Event Search / Hunt Workspace / Case Workspace | origins, Cloud audit pivots, handoffs, feedback and return | no detailed rewrite |
| Evidence Board / Hypotheses and Findings | Evidence/Finding candidates, contradictions and Cloud Hypotheses | read/link/prepare only; destination owners unchanged |
| Detection Engineering / Technical Workbench | TI operationalization and Cloud Detection Gap/Engineering handoffs | no rule/runtime UI defined |
| Analysis Workbench surfaces | technical inputs and Cloud Artifact/workload handoffs | no rewrite |
| Network Forensics | authorized flow/capture projections correlated with Cloud topology | Cloud Analysis does not replace Network Forensics |
| Platform Settings Sources/Integrations/Tenants/Environments/Health/Secrets | provider, source, connector, scope, health, schema, policy and secret-reference projections | read/request only; no connector or credential UI changed |
| Studio Control Room / Builder / Library | Tool/Run/Workflow provenance and bounded analysis requests | ownership unchanged |
| Govern Decisions / Approvals / Action Center | external release, permission/credential/resource change or future-response Action Requests | no target mutation UI defined |
| Shared Reporting / Export / Notifications | rendering, authorized minimized export and delivery mechanisms | consumed, not redefined |
| Inspector / Context Bar | selected object, Cloud scope, source, restrictions, sensitive-access level and return origin | consumption only |
| Threat Intelligence Workspace / Intelligence Library / Product Builder / Indicator Explorer | future detailed TI surfaces | capability links recorded; no Screen ID |
| Cloud Analysis Workspace / Cloud Inventory / Identity Graph | future detailed Cloud surfaces for CAP-INV-601..618 | conceptual capability links only; **no Screen ID** |

Required Cloud-related surfaces read include possible Cloud Workspace/Inventory/Identity Graph, Event Search, Hunt, Case, Evidence, Network Forensics, Detection Engineering, Threat Intelligence, Settings Integrations/Sources, Studio Control Room, Govern Decision, Technical Workbench and Shared Inspector. Screen specs modified: **0**. Detailed rewrites: **0**. New Screen IDs, wireframes, final buttons/columns/filters/animations/shortcuts and provider/protocol UI: **0**.
