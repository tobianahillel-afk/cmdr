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
  - REQ-BRAND-008
---
# Unresolved Decisions

Open decisions preserve uncertainty without promoting a proposal to an approved fact. A decision closes only through an explicit sponsor decision or an approved ADR, followed by updates to every affected canonical document and the traceability matrix.

## Phase 2 brand decisions

### OPEN-001 — Investigate palette direction

```yaml
decision_id: OPEN-001
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
depends_on: []
question: Which Investigate palette direction should become canonical?
options:
  - A — Verdigris Ledger
  - B — Lichen Archive
  - C — Mineral Graph
required_evidence:
  - light and dark workbench prototypes
  - contrast and color-vision-deficiency review
  - dense table, timeline, graph and code comparisons
  - family resemblance review against CMDR and Command
affected_files:
  - 02-brand/investigate/palette-proposals.md
  - 03-design-system/tokens/product-theme-tokens.md
requirement_id: REQ-PROD-048
```

### OPEN-002 — Govern palette direction

```yaml
decision_id: OPEN-002
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
depends_on: []
question: Which Govern palette direction should become canonical?
options:
  - A — Bronze Ledger
  - B — Civic Aubergine
  - C — Olive Charter
required_evidence:
  - decision-package and audit-ledger prototypes
  - non-punitive authority review
  - light and dark contrast review
  - distinction between recommendation, decision and execution
affected_files:
  - 02-brand/govern/palette-proposals.md
  - 03-design-system/tokens/product-theme-tokens.md
requirement_id: REQ-PROD-049
```

### OPEN-003 — CMDR Studio palette direction

```yaml
decision_id: OPEN-003
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
depends_on: []
question: Which CMDR Studio palette direction should become canonical?
options:
  - A — Workshop Sienna
  - B — Petrol Assembly
  - C — Iris Graphite
required_evidence:
  - builder, Control Room and Assurance prototypes
  - node-state and run-state differentiation
  - light and dark contrast review
  - review against generic purple AI and no-code conventions
affected_files:
  - 02-brand/studio/palette-proposals.md
  - 03-design-system/tokens/product-theme-tokens.md
requirement_id: REQ-PROD-050
```

### OPEN-004 — Final typography stack and licensing

```yaml
decision_id: OPEN-004
status: open
owner: Design System Lead
target_phase: Phase 2 review
blocking: false
depends_on: []
question: Which editorial sans and monospace stack should become canonical?
options:
  - Inter plus Inter Tight with a reviewed monospace
  - IBM Plex Sans plus IBM Plex Mono
  - a hybrid or equivalent editorial sans system supported by evidence
required_evidence:
  - license and legal review
  - French and English language coverage
  - tabular-number and dense-table evaluation
  - code, query, hash and identifier readability
  - Windows, macOS and Linux rendering
  - loading performance and fallback behavior
affected_files:
  - 02-brand/cmdr/typography.md
  - 03-design-system/foundations/typography.md
  - 03-design-system/tokens/*
requirement_id: REQ-PROD-051
```

### OPEN-016 — Final wordmark construction and optional symbol

```yaml
decision_id: OPEN-016
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
depends_on:
  - OPEN-004
question: Should CMDR remain wordmark-only or use an approved compact symbol, and what final optical construction should be adopted?
options:
  - wordmark-only system
  - wordmark plus distinct compact symbol
  - optically constructed CMDR monogram for compact contexts only
required_evidence:
  - small-size and favicon legibility
  - monochrome, Bone and Ink applications
  - trademark and legal review
  - optical-spacing study
  - application-icon and product-lockup tests
affected_files:
  - 02-brand/logo-and-wordmark.md
  - 02-brand/cmdr/logo.md
  - assets/brand/*
requirement_id: REQ-BRAND-008
```

`OPEN-016` blocks production of final logo assets, not the use of the textual name `CMDR`.

## Later product and technical decisions

### OPEN-005 — Initial forensic engines

```yaml
decision_id: OPEN-005
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: true
depends_on: []
question: How should initial forensic engines be classified?
options:
  - native
  - integrated
  - temporary-integration
required_evidence:
  - user outcomes
  - licensing and engine limits
  - evidence-integrity review
  - replacement strategy
affected_files:
  - 01-product-vision/capability-map.md
  - 07-investigate/modules/*
requirement_id: REQ-PROD-052
```

### OPEN-006 — Customers and Delivery applicability

```yaml
decision_id: OPEN-006
status: open
owner: Command Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Is Customers and Delivery an optional deployment capability or a generally available Command module?
options:
  - optional managed-service capability
  - generally available module
required_evidence:
  - managed-service and internal-deployment operating evidence
affected_files:
  - 06-command/modules/customer-and-reports/*
requirement_id: REQ-PROD-053
```

### OPEN-007 — Human Gate and Govern relationship

```yaml
decision_id: OPEN-007
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
depends_on: []
question: When must a Studio Human Gate create or reference a Govern Decision?
options:
  - risk-class-driven
  - policy-driven
  - combined risk and policy model
required_evidence:
  - action classes
  - authority model
  - separation of duties
  - audit requirements
affected_files:
  - 08-govern/*
  - 09-cmdr-studio/human-gates/*
requirement_id: REQ-PROD-054
```

### OPEN-008 — Endpoint platform support

```yaml
decision_id: OPEN-008
status: open
owner: Endpoint Agent Product Lead
target_phase: Phase 4
blocking: true
depends_on: []
question: Which operating systems and versions are visible at initial delivery?
options:
  - release-scoped Windows support
  - release-scoped Windows and Linux support
  - broader phased support including macOS
required_evidence:
  - user demand
  - sensor feasibility
  - support operating model
  - compatibility evidence
affected_files:
  - 11-endpoint-agent/platform-support.md
  - 10-platform-settings/endpoint-agent-fleet/*
requirement_id: REQ-PROD-055
```

### OPEN-010 — Density by role and activity

```yaml
decision_id: OPEN-010
status: open
owner: UX Architecture
target_phase: Phase 3
blocking: false
depends_on: []
question: Which density defaults and user controls apply by role and activity?
options:
  - activity defaults
  - role defaults
  - activity defaults with user overrides
required_evidence:
  - queue, case, workbench, decision and settings usability studies
affected_files:
  - 04-experience-architecture/information-density-model.md
  - 03-design-system/foundations/density.md
requirement_id: REQ-PROD-057
```

### OPEN-011 — Mobile forensic scope

```yaml
decision_id: OPEN-011
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Is mobile forensic capability planned or out of scope for initial delivery?
options:
  - planned
  - out-of-scope
required_evidence:
  - demand
  - acquisition legality
  - platform and engine strategy
affected_files:
  - 07-investigate/modules/*
requirement_id: REQ-PROD-058
```

### OPEN-012 — Cloud analysis scope

```yaml
decision_id: OPEN-012
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
depends_on: []
question: Which cloud analysis capabilities are in scope and how are they delivered?
options:
  - native
  - integrated
  - temporary-integration
  - planned
  - out-of-scope
required_evidence:
  - journeys
  - evidence sources
  - provider APIs
  - ownership boundaries
affected_files:
  - 07-investigate/modules/*
requirement_id: REQ-PROD-059
```

### OPEN-013 — Default governance for Action Class 2

```yaml
decision_id: OPEN-013
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
depends_on: []
question: Which Action Class 2 operations require Govern review by default?
options:
  - policy-controlled direct execution
  - mandatory Govern review
  - conditional review based on scope and reversibility
required_evidence:
  - risk and blast radius
  - rollback reliability
  - tenant policy
affected_files:
  - 14-security-permissions-and-trust/*
  - 08-govern/*
requirement_id: REQ-PROD-060
```

### OPEN-014 — Artifact versus Attachment

```yaml
decision_id: OPEN-014
status: open
owner: Investigate Product Lead
target_phase: Phase 7
blocking: true
depends_on: []
question: Should Attachment be a separate object or a typed Artifact relationship?
options:
  - separate Attachment object
  - typed Artifact relation
required_evidence:
  - evidence lifecycle
  - provenance
  - collaboration
  - retention and export
affected_files:
  - 05-domain-model/objects/artifact.md
  - 00-governance/ownership-register.md
requirement_id: REQ-PROD-061
```

### OPEN-015 — Automation Run to Response Run bridge

```yaml
decision_id: OPEN-015
status: open
owner: CMDR Studio Product Lead
target_phase: Phase 4
blocking: true
depends_on:
  - OPEN-007
  - OPEN-013
question: What explicit governed transition connects Automation Run to Response Run?
options:
  - governed object transition with referenced Decision
  - policy-authorized transition for eligible classes
required_evidence:
  - action classification
  - Decision requirements
  - Tool effects
  - audit model
affected_files:
  - 08-govern/*
  - 09-cmdr-studio/*
requirement_id: REQ-PROD-062
```

## Resolved history

### OPEN-009 — Capability delivery-classification authority

```yaml
decision_id: OPEN-009
status: resolved
owner: Product Architecture
target_phase: Phase 1
blocking: false
depends_on: []
question: Who approves a capability delivery classification?
decision: The capability owner proposes; the owning Product Lead and Product Architecture approve; Security reviews trust, action and sensitive-data capabilities; Engineering must review native or integrated delivery claims; QA and Traceability verify evidence and update the register.
affected_files:
  - 00-governance/review-and-approval-process.md
  - 01-product-vision/capability-map.md
requirement_id: REQ-PROD-056
```

## Maintenance rules

- An option is not a decision.
- A proposed palette remains `proposed` even when individual contrast pairs pass.
- A question closes only with decision evidence, owner and affected-file updates.
- Starting the target phase without sufficient evidence does not close the question.
- Generic placeholders are prohibited.
