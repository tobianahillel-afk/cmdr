---
id: capability-register-studio-std1
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: registry
---
# Capability Register — Studio STD-1

Parent: Delivery Roadmap Phase 5 — Studio and Endpoint. Execution lot: STD-1. Namespace: `CAP-STD-*`. All entries are `draft / defined / planned`; none is implemented/native/integrated.

| Capability ID | Name | Owner/module | Status | Delivery | Canonical file | Primary roles | Primary objects/concepts | Requirements | OPEN | Consumers |
|---|---|---|---|---|---|---|---|---|---|---|
| CAP-STD-001 | Studio Library and Asset Catalog | Studio/library | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-001-studio-library-and-asset-catalog.md` | Automation Designer; SOC Analyst | Skill; Version; asset refs | REQ-PROD-006,016; REQ-OBJ-009 | — | Studio; Command; Investigate |
| CAP-STD-002 | Studio Asset Metadata, Ownership and Classification | Studio/library | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-002-studio-asset-metadata-ownership-and-classification.md` | Automation Designer; Auditor | asset metadata; owner refs | REQ-PROD-006,016; REQ-OBJ-009 | — | Library; Quality |
| CAP-STD-003 | Tool Definition and Functional Capability Contract | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-003-tool-definition-and-functional-capability-contract.md` | Automation Designer; Studio Operator | Tool; Version refs | REQ-PROD-006,016; REQ-OBJ-009; REQ-AI-002 | — | Skills; consumers |
| CAP-STD-004 | Tool Input, Parameter and Validation Contract | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-004-tool-input-parameter-and-validation-contract.md` | Automation Designer; analyst | Tool; Secret Reference | REQ-PROD-006,016; REQ-SEC-001 | — | Tool Calls; Skills |
| CAP-STD-005 | Tool Output, Result and Error Contract | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-005-tool-output-result-and-error-contract.md` | Studio Operator; analyst | Tool output; Result ref | REQ-PROD-002,005,016; REQ-OBJ-009 | — | callers; Trace |
| CAP-STD-006 | Tool Versioning, Compatibility and Deprecation | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-006-tool-versioning-compatibility-and-deprecation.md` | Automation Designer; reviewer | Tool Version | REQ-PROD-006,016; REQ-OBJ-009 | — | Library; Skills; Calls |
| CAP-STD-007 | Tool Permission, Risk and Execution Eligibility | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-007-tool-permission-risk-and-execution-eligibility.md` | Studio Operator; Security Reviewer | Tool; Decision ref | REQ-PROD-006,016; REQ-SEC-001,002 | OPEN-013 | callers; Govern |
| CAP-STD-008 | Tool Call Request and Invocation Context | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-008-tool-call-request-and-invocation-context.md` | Studio Operator; analyst | Tool Call request; Tool | REQ-PROD-006,016; REQ-OBJ-009; REQ-SEC-001 | OPEN-015 | runtime; callers |
| CAP-STD-009 | Tool Call Lifecycle, Technical Outcome and Provenance | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-009-tool-call-lifecycle-technical-outcome-and-provenance.md` | Studio Operator; Auditor | Tool Call; Trace ref | REQ-PROD-005,016; REQ-OBJ-009 | OPEN-015 | callers; future Control Room |
| CAP-STD-010 | Skill Definition and Reusable Capability Contract | Studio/skills | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-010-skill-definition-and-reusable-capability-contract.md` | Automation Designer; reviewer | Skill; Tool refs | REQ-PROD-006,016; REQ-OBJ-009; REQ-AI-002 | — | Library; future Workflow/Agent |
| CAP-STD-011 | Skill Composition, Dependencies and Preconditions | Studio/skills | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-011-skill-composition-dependencies-and-preconditions.md` | Automation Designer; reviewer | Skill dependency refs | REQ-PROD-006,016; REQ-OBJ-009 | — | Skills; consumers |
| CAP-STD-012 | Skill Input, Output and Parameter Contract | Studio/skills | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-012-skill-input-output-and-parameter-contract.md` | Automation Designer; analyst | Skill I/O; Secret refs | REQ-PROD-006,016; REQ-SEC-001 | — | Skills; Tools; consumers |
| CAP-STD-013 | Skill Versioning, Lifecycle and Deprecation | Studio/skills | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-013-skill-versioning-lifecycle-and-deprecation.md` | Automation Designer; reviewer | Skill; Version | REQ-PROD-006,016; REQ-OBJ-009 | — | Library; consumers |
| CAP-STD-014 | Skill Discovery, Reuse and Consumer Linking | Studio/skills | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-014-skill-discovery-reuse-and-consumer-linking.md` | Automation Designer; analysts | Skill; consumer refs | REQ-PROD-006,016; REQ-OBJ-009 | — | Command; Investigate; future Studio |
| CAP-STD-015 | Provider, Runtime, Integration and Secret Reference Boundaries | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-015-provider-runtime-integration-and-secret-reference-boundaries.md` | Studio Operator; Platform Admin | provider/runtime/integration/Secret refs | REQ-PROD-006,016,017; REQ-SEC-001 | — | Tools; Skills; Settings |
| CAP-STD-016 | Studio Foundations Cross-Product Contracts and Provenance | Studio/foundations | draft | defined/planned | `09-cmdr-studio/capabilities/cap-std-016-studio-foundations-cross-product-contracts-and-provenance.md` | Studio Operator; cross-product roles | Tool/Call/Skill refs; Decision/Run/Result refs | REQ-PROD-006,009,016; REQ-OBJ-009; REQ-AI-002 | OPEN-007,015 | all products |

**Count: 16; defined: 16; proposed: 0; planned: 16; sections: 432; mandatory tables: 96.** IDs are unique, newly allocated and not recycled.
