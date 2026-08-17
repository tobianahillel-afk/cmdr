---
id: endpoint-ept4-cross-product-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-4 Cross-Product Links

## Investigate
Owns Case, Collection Request, Evidence, Finding and Artifact qualification. Endpoint receives request refs and returns neutral technical output/status/provenance. Collected Output != Evidence automatically.

## Govern
Owns Action Request, Approval, Decision, Response Run, Result, production-response authority, verification and rollback. Functional governed path may be `Action Request → Decision → Response Run → Endpoint technical execution → technical output/status → Govern reconciliation → Result`. Endpoint creates none of the Govern objects.

## Studio
Owns Tool, Tool Call, Workflow, Automation Agent and Automation Run. Tool Call/Automation Run may become caller/provenance refs but never become Endpoint execution. OPEN-015 remains open.

## Settings
Owns Fleet, Endpoint Policy, credentials, secrets, providers and runtime/admin configuration. Endpoint consumes Secret References only; raw secret values do not enter request/output/audit.

## Shared
Owns generic Jobs, Trace, Activity, Recovery, Reporting and Export/transfer mechanisms when applicable. Local Collection Operation != Background Job; technical session != Shared mechanism.

## Future EPT-5
Effectful containment primitives such as isolation, process/service/network remediation, quarantine and response verification stop at EPT-5 boundary.