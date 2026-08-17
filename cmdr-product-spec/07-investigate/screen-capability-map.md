---
id: investigate-screen-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-UX-001, REQ-UX-010, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-011, OPEN-012]
---
# Screen capability map — Threat Intelligence, Cloud Analysis and Mobile Forensics functional scope

| Existing screen or surface | Consumption | Change through 4B.4B |
|---|---|---|
| Entity Graph / Graph / Timeline | Intelligence candidates plus Cloud and Mobile scope/identity/resource/device/correlation candidates | links/map only; owners unchanged |
| Event Search / Hunt Workspace / Case Workspace | origins, telemetry pivots, Mobile/Cloud handoffs, feedback and return | no detailed rewrite |
| Evidence Board / Hypotheses and Findings | Evidence/Finding candidates, contradictions, Cloud/Mobile Hypotheses | read/link/prepare only; destination owners unchanged |
| Detection Engineering / Technical Workbench | TI operationalization plus Cloud/Mobile Detection Engineering handoffs | no rule/runtime UI defined |
| Analysis Workbench / Static Analysis / Reverse Engineering | technical inputs and Derived Artifact handoffs | no rewrite; destination analysis remains owner-controlled |
| Disk and Artifact Forensics | source Artifact/disk/filesystem context and deeper handoff from Mobile | Mobile analysis does not replace Disk/Filesystem Forensics |
| Network Forensics | authorized capture/flow context plus Mobile stored connectivity observations | Mobile network state does not replace Network Forensics |
| Cloud Analysis | cloud-backed/synchronized Mobile artifact handoff | Mobile does not absorb Cloud scope or provider administration |
| Endpoint Agent Fleet | optional declared device/agent capability, health and management projection | no assumption a mobile device has an Agent; no Fleet rewrite |
| Platform Settings Sources/Integrations/Tenants/Environments/Health/Secrets | provider/source/MDM-like future administration, connector, Fleet, policy and secret-reference projections | read/request only; no connector/credential/device-management UI changed |
| Studio Control Room / Builder / Library | Tool/Run/Workflow provenance and bounded analysis/recovery requests | ownership unchanged; no silent Tool execution |
| Govern Decisions / Approvals / Action Center | acquisition authority and any future real-device action request | no unlock, profile, isolation, lock, wipe, revocation or response UI defined by Mobile |
| Shared Reporting / Export / Notifications | rendering, authorized minimized export and delivery mechanisms | consumed, not redefined |
| Inspector / Context Bar | selected source/device/package/record, restrictions, sensitive-access level and return origin | consumption only |
| Threat Intelligence Workspace / Intelligence Library / Product Builder / Indicator Explorer | future detailed TI surfaces | capability links recorded; no Screen ID |
| Cloud Analysis Workspace / Cloud Inventory / Identity Graph | future detailed Cloud surfaces for CAP-INV-601..618 | conceptual capability links only; no Screen ID |
| Mobile Forensics Workspace / Device Evidence Library / Mobile Timeline | future detailed Mobile surfaces for CAP-INV-701..719 | conceptual capability links only; **no Screen ID** |

Mobile-related surfaces actually read for this phase include Disk and Artifact Forensics, Cloud Analysis, Network Forensics, Static Analysis, Reverse Engineering, Case Workspace, Evidence Board, Event Search, Endpoint Agent Fleet, Platform Settings Sources, Studio Control Room and Govern Action Center, plus shared Inspector/Context patterns from canonical UX sources. No active Mobile Screen ID existed at the starting SHA.

Screen specs modified: **0**. Detailed rewrites: **0**. New Screen IDs: **0**. Wireframes: **0**. Final buttons, columns, filters, animations and shortcuts: **0**. Mobile capability links are conceptual only until the dedicated screen phase.
