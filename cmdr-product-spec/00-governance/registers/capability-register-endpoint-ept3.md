---
id: capability-register-endpoint-ept3
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Capability Register — Endpoint EPT-3

Execution lot: **EPT-3 — Local Detection and Endpoint Investigation** under Delivery Roadmap Phase 5 — Studio and Endpoint. All rows are `draft / defined / planned`. Documentary definition does not prove implementation.

| ID | Capability | Owner | Primary roles | Primary concepts | Key consumers | Requirements | OPEN | Delivery |
|---|---|---|---|---|---|---|---|---|
| CAP-EPT-031 | Detection Content Consumption and Eligibility Boundary | Endpoint Agent | Endpoint Operator, SOC Analyst, Detection Engineer | content reference, eligibility | EPT-3, Investigate, Settings | REQ-PROD-006/012/018/019, REQ-INV-006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-032 | Local Detection Evaluation and Match Semantics | Endpoint Agent | Endpoint Operator, SOC Analyst, Detection Engineer | local evaluation, match | EPT-3, Investigate, Command | REQ-PROD-012/018/019, REQ-INV-006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-033 | Local Detection Signal Lifecycle and State | Endpoint Agent | Endpoint Operator, SOC Analyst, Command consumer | local signal candidate | EPT-3, Command, Investigate | REQ-PROD-012/018/019, REQ-INV-006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-034 | Detection Context, Severity, Confidence and Rationale | Endpoint Agent | SOC Analyst, Detection Engineer | detection context, severity, confidence, rationale | EPT-3, Investigate, Command | REQ-PROD-012/018/019, REQ-INV-006, REQ-AI-002, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-035 | Detection Grouping, Deduplication and Suppression-State Boundary | Endpoint Agent | SOC Analyst, Detection Engineer, Platform Admin | grouping, duplicate candidate, suppression projection | EPT-3, Investigate, Settings | REQ-PROD-012/018/019, REQ-INV-006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-036 | Detection Coverage, Health and Gap Assessment | Endpoint Agent | Endpoint Operator, Detection Engineer, SOC Analyst | coverage, detection health, gap | Investigate, Settings, Quality | REQ-PROD-012/018/019, REQ-INV-006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-037 | Process Tree and Execution Context Investigation | Endpoint Agent | SOC Analyst, Endpoint Operator | process context, ancestry | Investigate, EPT-3 | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-038 | File and Filesystem Context Investigation | Endpoint Agent | SOC Analyst, Endpoint Operator | file context | Investigate, EPT-3 | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008 | defined / planned |
| CAP-EPT-039 | Network and Connection Context Investigation | Endpoint Agent | SOC Analyst, Endpoint Operator | network context | Investigate, EPT-3 | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008 | defined / planned |
| CAP-EPT-040 | User, Session and Authentication Context Investigation | Endpoint Agent | SOC Analyst, Security Reviewer | user/session/auth context | Investigate, Settings | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008 | defined / planned |
| CAP-EPT-041 | Service, Module, Driver and System Context Investigation | Endpoint Agent | SOC Analyst, Endpoint Operator | system context | Investigate, EPT-3 | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008 | defined / planned |
| CAP-EPT-042 | Endpoint Timeline and Activity Correlation Investigation | Endpoint Agent | SOC Analyst, Endpoint Operator | local timeline, correlation | Investigate, Shared | REQ-PROD-012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-043 | Endpoint Investigation Context Retrieval and Pivot Semantics | Endpoint Agent | SOC Analyst, Endpoint Operator | investigation pivot | Investigate, Shared | REQ-PROD-006/012/018/019, REQ-INV-001, REQ-SEC-001/002 | OPEN-008 | defined / planned |
| CAP-EPT-044 | Detection-to-Investigation Pivot and Context Expansion | Endpoint Agent | SOC Analyst, Endpoint Operator | detection-investigation pivot, expansion | Investigate | REQ-PROD-006/012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008, OPEN-017 | defined / planned |
| CAP-EPT-045 | Endpoint Investigation Summary and Consumer Handoff | Endpoint Agent | SOC Analyst, Command consumer | local summary, handoff | Investigate, Command, Govern | REQ-PROD-006/012/018/019, REQ-INV-001/006, REQ-SEC-001/002 | OPEN-008, OPEN-015, OPEN-017 | defined / planned |
| CAP-EPT-046 | Detection and Investigation Provenance and Cross-Product Contracts | Endpoint Agent | SOC Analyst, Auditor | provenance chain, cross-product handoff | Investigate, Command, Govern, Shared | REQ-PROD-006/012/018/019, REQ-OBJ-008, REQ-INV-006, REQ-SEC-001/002/004 | OPEN-008, OPEN-015, OPEN-017 | defined / planned |

## Totals
**16 capabilities / 432 numbered sections / 96 mandatory tables / at least 48 GWT.** Duplicate IDs 0; recycled IDs 0; owner conflicts 0; empty mandatory tables 0. Endpoint cumulative after EPT-3: **46 / 1242 / 276**.

`CAP-EPT-001..030` remain unchanged. `CAP-EPT-047+` is not allocated or reserved by EPT-3. EPT-4..EPT-6 remain NOT STARTED.