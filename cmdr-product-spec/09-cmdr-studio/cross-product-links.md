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
