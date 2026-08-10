---
id: studio-std1-cross-product-links
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
open_decisions: [OPEN-007, OPEN-015]
---
# STD-1 Cross-Product Links

## Govern
Workflow != Playbook; Human Gate != Approval; Automation Run != Response Run; Tool Call output != Result. Studio may consume Decision/Response Run references but does not own response authority.

## Investigate
Investigate may select an authorized Tool/Skill, trigger a Tool Call, consume technical output and then produce its own conclusions. Tool output != Evidence or Finding automatically.

## Command
Command may consume authorized Studio asset references and future execution projections. Studio never becomes the Command Work Queue or Incident owner.

## Endpoint Agent
A Tool/Workflow may later reference Endpoint technical capabilities. Studio never becomes owner of an Endpoint primitive or device/agent health.

## Platform Settings
Studio consumes provider/runtime/integration/Secret Reference and health/config projections. Settings retains providers, integrations, credentials, raw secrets, tenant/environment administration and connection lifecycle.

## Shared Capabilities
Studio consumes Search, Linking, Versioning, Jobs, Notifications, Trace, Activity, Reporting, Export, Collaboration and Recovery where applicable. It creates no competing generic engine.

## STD-2 — Workflow Builder & Orchestration addendum

### Govern
Workflow Definition/Version and validation remain Studio-owned. Branch != Decision, condition != Policy, Human Gate completion != Approval/Decision, compensation != Govern rollback. CAP-GOV-025 remains the future execution handoff boundary. OPEN-007/013/015 stay open.

### Investigate
Investigate may reference a Workflow/Tool/Skill and receive attributed technical outputs later, but retains Case/Evidence/Finding/Hypothesis ownership and qualification.

### Command
Command may reference authorized Workflow versions without becoming Studio or transferring Incident/Task ownership.

### Endpoint Agent
A Workflow may reference a future Endpoint capability only by typed reference; Studio never owns or defines the Endpoint primitive.

### Platform Settings
Workflow data/runtime bindings consume opaque Settings-owned provider/integration/environment/Secret References. Raw secrets are excluded.

### Shared Capabilities
Workflow uses generic Trace/Activity/Jobs/Versioning/Search/Notifications/Recovery mechanisms without creating competing engines.

## STD-3 — Agents, Human Gates & Runtime Control addendum

### Govern
Automation Run != Response Run; Human Gate != Approval/Decision; Studio Runtime Outcome != Result. Effectful start/resume/retry uses explicit Govern authority when required. OPEN-007/013/015 remain open.

### Investigate
Investigate may receive attributed automation/agent outputs, but qualification as Evidence/Finding remains Investigate-owned.

### Command
Command may call/observe an authorized Automation Run but Studio does not become Incident/Task/Work Queue owner.

### Endpoint Agent
Studio may later orchestrate a typed Endpoint technical capability; Endpoint retains its primitive, local runtime/state and device-side result. STD-3 creates no Endpoint capability.

### Platform Settings
Agent/runtime consumers use Settings-owned Principal/Role, provider/integration, environment, health and opaque Secret Reference projections; no raw secret/admin ownership transfers.

### Shared Capabilities
Automation Run != Shared Job. Generic queue/scheduling, Jobs, Trace, Activity, Notifications and Recovery remain Shared; Studio owns Run business semantics.