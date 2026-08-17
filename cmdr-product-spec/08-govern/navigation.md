---
id: 08-govern-navigation
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-008, REQ-PROD-015, REQ-UX-001, REQ-UX-009]
open_decisions: [OPEN-010]
---
# Govern Navigation

## Local order

The nine canonical Govern modules remain visible according to the existing product navigation model. GOV-1 does not create a second navigation tree.

For the GOV-1 workflow, the normal path is:

`Response Inbox → Action Center → Policy Gates / Approvals & Authorities → Decision Register`.

Playbooks and Runs & Rollback remain future GOV-2 destinations. Audit Trail and Response Metrics remain future GOV-3 destinations.

## Deep-link context

A Govern deep link must preserve, subject to permission:
- tenant and environment;
- Action Request id and version;
- source product and return origin;
- selected Incident/Case/Finding references;
- target/scope identifiers without leaking restricted values;
- active review context when safe to restore.

## Cross-product entry

Command, Investigate, Detection Engineering, Threat Intelligence and analysis modules can enter Govern through an Action Request or request-preparation handoff. Govern never interprets a generic cross-product navigation link as an Approval or Decision.

## Return behavior

`request-more-information`, rejection, deferral and final Decision dispositions preserve a stable return origin. Returning to a source product restores the source object/workspace when still authorized; otherwise the source product opens its safe parent context without exposing restricted data.

## Permission behavior

Navigation never grants authority. A visible Action Center or Decision Register route does not imply permission to read all context, approve, decide or execute. Missing permission produces an explicit permission-denied/partial state and preserves safe return behavior.

## Screen rule

GOV-1 modifies no detailed screen specification and creates no Screen ID. Module routes and screen identities remain those already registered.
