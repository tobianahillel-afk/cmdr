---
id: source-unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-048
  - REQ-PROD-062
---
# Unresolved Decisions

Generic `À compléter` placeholders are forbidden. This register preserves decisions that remain genuinely open without presenting an option as an approved fact.

## Open decisions

### OPEN-001 — Select the Investigate palette direction.

```yaml
decision_id: OPEN-001
status: open
owner: Brand Design Lead
target_phase: Phase 2
blocking: false
depends_on: []
question: Select the Investigate palette direction.
options: Compatible CMDR-family proposals only; no final values are decided.
required_evidence: Contrast, workbench density, differentiation and accessibility review.
affected_files: 02-brand/investigate/*; 03-design-system/tokens/product-theme-tokens.md
requirement_id: REQ-PROD-048
```

### OPEN-002 — Select the Govern palette direction.

```yaml
decision_id: OPEN-002
status: open
owner: Brand Design Lead
target_phase: Phase 2
blocking: false
depends_on: []
question: Select the Govern palette direction.
options: Compatible CMDR-family proposals only; no final values are decided.
required_evidence: Decision-authority tone, accessibility and non-punitive visual review.
affected_files: 02-brand/govern/*; product theme tokens
requirement_id: REQ-PROD-049
```

### OPEN-003 — Select the Studio palette direction.

```yaml
decision_id: OPEN-003
status: open
owner: Brand Design Lead
target_phase: Phase 2
blocking: false
depends_on: []
question: Select the Studio palette direction.
options: Compatible CMDR-family proposals only; no final values are decided.
required_evidence: Builder, simulation and execution-observability use cases.
affected_files: 02-brand/studio/*; product theme tokens
requirement_id: REQ-PROD-050
```

### OPEN-004 — Confirm the final typography stack and licensing.

```yaml
decision_id: OPEN-004
status: open
owner: Design System Lead
target_phase: Phase 2
blocking: false
depends_on: []
question: Confirm the final typography stack and licensing.
options: Inter, Inter Tight, IBM Plex Sans or equivalent editorial sans; monospace for technical content.
required_evidence: License, language coverage, performance and readability evaluation.
affected_files: 02-brand/cmdr/typography.md; 03-design-system/foundations/typography.md
requirement_id: REQ-PROD-051
```

### OPEN-005 — Classify the initial forensic engines.

```yaml
decision_id: OPEN-005
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: true
depends_on: []
question: Classify the initial forensic engines.
options: native, integrated or temporary-integration per capability.
required_evidence: User outcomes, licensing, engine limitations, evidence integrity and replacement strategy.
affected_files: 01-product-vision/capability-map.md; 07-investigate workbench modules
requirement_id: REQ-PROD-052
```

### OPEN-006 — Define where Customers and Delivery applies.

```yaml
decision_id: OPEN-006
status: open
owner: Command Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Define where Customers and Delivery applies.
options: optional deployment capability or generally available module.
required_evidence: Operating-model evidence from managed-service and internal deployments.
affected_files: 06-command/modules/customer-and-reports/*
requirement_id: REQ-PROD-053
```

### OPEN-007 — Define when a Studio Human Gate must create or reference a Govern Decision.

```yaml
decision_id: OPEN-007
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
depends_on: []
question: Define when a Studio Human Gate must create or reference a Govern Decision.
options: risk-class and policy-driven options; no universal rule is decided.
required_evidence: Action classes, authority model, separation of duties and audit needs.
affected_files: 01-product-vision/operating-model.md; 08-govern/*; 09-cmdr-studio/human-gates/*
requirement_id: REQ-PROD-054
```

### OPEN-008 — Define visible platform support at initial delivery.

```yaml
decision_id: OPEN-008
status: open
owner: Endpoint Agent Product Lead
target_phase: Phase 4
blocking: true
depends_on: []
question: Define visible platform support at initial delivery.
options: Windows, Linux and macOS scopes remain open by release.
required_evidence: User demand, sensor feasibility, operational support and compatibility matrix.
affected_files: 11-endpoint-agent/platform-support.md; 10-platform-settings/endpoint-agent-fleet/*
requirement_id: REQ-PROD-055
```

### OPEN-010 — Define density preferences by role and activity.

```yaml
decision_id: OPEN-010
status: open
owner: UX Architecture
target_phase: Phase 3
blocking: false
depends_on: []
question: Define density preferences by role and activity.
options: role defaults with user-adjustable density are candidates, not decisions.
required_evidence: Usability studies for queue, case, workbench, decision and settings activities.
affected_files: 01-product-vision/target-users.md; 04-experience-architecture/information-density-model.md
requirement_id: REQ-PROD-057
```

### OPEN-011 — Define mobile forensic scope.

```yaml
decision_id: OPEN-011
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Define mobile forensic scope.
options: planned or out-of-scope for the initial release.
required_evidence: User demand, acquisition legality, platform coverage and engine strategy.
affected_files: 07-investigate workbench capability documents
requirement_id: REQ-PROD-058
```

### OPEN-012 — Define cloud analysis scope.

```yaml
decision_id: OPEN-012
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Define cloud analysis scope.
options: native, integrated, temporary-integration, planned or out-of-scope per provider and artifact type.
required_evidence: User journeys, evidence sources, provider APIs and responsibility boundaries.
affected_files: 07-investigate workbench and intelligence documents
requirement_id: REQ-PROD-059
```

### OPEN-013 — Define default governance for Action Class 2.

```yaml
decision_id: OPEN-013
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
depends_on: []
question: Define default governance for Action Class 2.
options: policy-controlled direct execution or Govern review depending on scope and reversibility.
required_evidence: Risk assessment, rollback reliability, tenant policy and blast-radius evidence.
affected_files: 01-product-vision/operating-model.md; 14-security-permissions-and-trust/*
requirement_id: REQ-PROD-060
```

### OPEN-014 — Define the Artifact versus Attachment boundary.

```yaml
decision_id: OPEN-014
status: open
owner: Investigate Product Lead
target_phase: Phase 7
blocking: true
depends_on: []
question: Define the Artifact versus Attachment boundary.
options: separate Attachment object or typed Artifact relation remain possible.
required_evidence: Evidence lifecycle, collaboration, provenance, retention and export requirements.
affected_files: 00-governance/ownership-register.md; 05-domain-model/objects/artifact.md; evidence model
requirement_id: REQ-PROD-061
```

### OPEN-015 — Define the Automation Run to Response Run bridge.

```yaml
decision_id: OPEN-015
status: open
owner: CMDR Studio Product Lead
target_phase: Phase 4
blocking: true
depends_on: []
question: Define the Automation Run to Response Run bridge.
options: explicit governed transition; exact trigger and object mapping remain open.
required_evidence: Action classification, Decision requirements, tool effects and audit model.
affected_files: 01-product-vision/operating-model.md; Studio and Govern functional specifications
requirement_id: REQ-PROD-062
```

## Resolved during Phase 1

### OPEN-009 — Capability-status review authority

```yaml
decision_id: OPEN-009
status: resolved
owner: Product Architecture
target_phase: Phase 1
blocking: false
depends_on: []
question: Who approves a capability delivery classification?
decision: The capability owner proposes the classification. The owning Product Lead and Product Architecture approve it. Security review is mandatory for trust, action or sensitive-data capabilities. Engineering review is mandatory before a native or integrated delivery claim. QA and Traceability verify the cited evidence and update the capability register and traceability matrix.
required_evidence: owner rationale, user outcome, engine responsibility, delivery evidence, limits, dependencies and replacement plan for temporary integrations
affected_files: 00-governance/review-and-approval-process.md; 01-product-vision/capability-map.md; 00-governance/registers/capability-register.md
requirement_id: REQ-PROD-056
```

The resolved record remains here for history. It is not an open blocker.

## Maintenance rules

- A question is closed only by an explicit source decision or an approved ADR.
- Closing a question updates every affected canonical document and the traceability matrix.
- If a target phase begins without sufficient evidence, the question remains open and the dependent document cannot become Validated.
- No open decision may be replaced by a generic placeholder or an unsupported assumption.
