---
id: studio-cross-product-links
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Studio Cross-Product Links

## STD-1 boundaries — preserved
- Workflow != Playbook; Human Gate != Approval; Automation Run != Response Run; Tool Call output != Result.
- Investigate qualifies its own Evidence/Finding from technical outputs.
- Command retains Incident/Work Queue ownership.
- Endpoint owns technical primitives.
- Settings owns provider/runtime/integration/Secret Reference administration.
- Shared owns Search/Linking/Versioning/Jobs/Notifications/Trace/Activity/Reporting/Export/Collaboration/Recovery.

## STD-2 additions

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

### Shared
Workflow uses generic Trace/Activity/Jobs/Versioning/Search/Notifications/Recovery mechanisms without creating competing engines.
