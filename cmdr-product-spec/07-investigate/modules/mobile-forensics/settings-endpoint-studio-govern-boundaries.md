---
id: investigate-mobile-forensics-owner-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-011, OPEN-013, OPEN-015]
---
# Settings, Endpoint, Studio and Govern Boundaries

## Platform Settings
Owns MDM/EMM providers, configured mobile sources, integrations/connectors, credentials, secret references, administratively supported platforms, device Fleet, storage, retention, source health, policies, tenant/environment configuration and access administration. Mobile consumes permission-aware projections only.

## Endpoint Agent
Owns declared native capabilities, health, policy/version state, telemetry and authorized local acquisition contribution when an agent exists. Mobile Forensics does not assume a mobile Endpoint Agent and does not own Agent enrollment, policy or execution.

## Studio
Owns Tool, Tool Call, Skill, Workflow, Automation Agent, Automation Run, Human Gate, generic evaluation, retry and execution provenance. Mobile selects or requests Tools under permission but does not redefine Tool/Run lifecycle.

## Govern
Owns Decision, Approval, Action Request, Response Run, Result and authority for all real-device actions. Remote acquisition requiring governed authority, profile installation, account revocation, configuration change, device isolation, lock, wipe or destructive remediation cannot execute inside Mobile Forensics.

## Command
Retains Detection, Signal, Alert, Incident, priority and operational dispositions. A Mobile anomaly or correlation never becomes an Alert/Incident automatically.

## Shared
Retains generic Entity, Graph, Timeline, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration and Recovery.

## Rule
A projection never transfers ownership or permission. Mobile can prepare a request/handoff, but the destination re-evaluates authorization and owns its resulting objects/actions.
