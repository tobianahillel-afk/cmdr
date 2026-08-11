---
id: endpoint-ept4-object-consumption-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-4 Object Consumption Map

| Concept | Canonical owner | Endpoint usage | Local operations | Boundary/gap | Future lot |
|---|---|---|---|---|---|
| Collection Request | Investigate | consume id/version/scope | technical intake only | Request != Operation | — |
| Collection Scope | Investigate / request context | consume and technically bound | plan projection | no scope expansion | — |
| Collection Operation | Endpoint technical concept | own local attempt/state | create/reconcile | != Collection Job/Background Job | — |
| Collection Item | Endpoint technical concept | own neutral output | create/partial/fail | != Artifact/Evidence/Finding | Investigate qualifies |
| Collection Package | Endpoint technical concept | assemble refs/metadata | create/version | != Evidence package | Investigate/Trust |
| Collection Progress | Endpoint | local state | reconcile | Shared Job remains Shared | — |
| Collection Output | Endpoint | raw/local output refs | create/read | OPEN-014 | Investigate qualifies |
| Collection Transfer | Endpoint | target-side transfer state | create/reconcile | Shared mechanism may carry | — |
| Live Response Session business record | Investigate | consume reference | none | != Endpoint Technical Session | — |
| Endpoint Technical Session | Endpoint | local connectivity/runtime semantics | create/reconcile/close | != Response Run/Automation Run | — |
| Technical Execution Request | Endpoint | local intake | create/validate | != Decision/Tool Call | — |
| Technical Execution Attempt | Endpoint | target-side execution | create/reconcile | != Response Run | — |
| Technical Execution Output | Endpoint | raw/local output | append/reference | != Govern Result/Evidence | external qualification |
| Command Request | Endpoint technical concept | exact bounded intent | create/accept/reject | no command catalog | — |
| Script Reference | source owner / OPEN-014-aware | consume exact ref/version | no identity mutation | no universal runtime | — |
| Operator Context | Investigate/Identity + Endpoint projection | bind actor/control | local projection | operator != approver | — |
| Endpoint Agent / Capability | Endpoint | source facts | read/local state | availability != authority | — |
| Endpoint Investigation Summary | Endpoint EPT-3 | upstream need/context | read | pivot != Collection | — |
| Artifact | Investigate | destination qualification only | no auto-create | OPEN-014 | — |
| Attachment | unresolved | reference only | none | OPEN-014 | future object decision |
| Evidence / Case / Finding | Investigate | source/destination refs | none | no auto qualification | — |
| Action Request / Decision / Response Run / Result | Govern | authority/reconciliation refs | none | Endpoint != Govern objects | EPT-5 consumes |
| Tool / Tool Call / Automation Run | Studio | caller/provenance refs | none | OPEN-015 | — |
| Secret Reference | Settings | reference-only input | use at authorized executor boundary | raw secret never stored | — |
| Job / Trace / Activity / Export | Shared | consume generic mechanisms | no ownership transfer | Endpoint records local semantics only | — |

No physical schema is created.