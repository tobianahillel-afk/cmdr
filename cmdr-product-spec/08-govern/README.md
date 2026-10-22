---
id: 08-govern-readme
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-OBJ-005, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Govern

## Mission
Govern is CMDR's authority, governed-response and response-review product. GOV-1 governs Action Requests through Decision/Handoff; GOV-2 governs Playbooks, execution planning, Response Runs, verification, rollback/recovery and canonical Result; GOV-3 defines Govern audit interpretation, Govern-specific metrics, continuous-improvement packages and documentary closure without duplicating Shared infrastructure.

## Canonical programme identity
- Parent: **Delivery Roadmap Phase 4 — Govern**, id `roadmap-phase-4-govern` — **PASS**.
- GOV-1 — Action Requests, Policy, Authorities and Decisions: **PASS**, 16 capabilities / 432 sections / 96 tables / historical 180 gates.
- GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback: **PASS**, 17 / 459 / 102 / historical 190 gates.
- GOV-3 — Audit Trail, Response Metrics and Govern Closure: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**, 14 / 378 / 84.
- Govern capability specification: **PASS**, 47 capabilities / 1269 sections / 282 tables.
- GOV-1/GOV-2/GOV-3 are execution lots only. `Phase 4C/4D/4E Govern` do not exist.

## Ownership
Govern owns the GOV-1/GOV-2 semantics plus Govern Audit Event interpretation/reconstruction, completeness/gap/contradiction assessment, Audit Review/Evidence Package composition, Govern-specific metrics, Trend/Control Health Assessment, Continuous Improvement Package and closure provenance.

Consumed but not owned: Shared Trace/Activity/Search/Metrics/Reporting/Export/Jobs/Versioning; Settings users/roles/tenants/environments/retention/storage/providers/integrations/secrets; Security permissions/privacy/integrity/legal-hold; Command Incident/Work Queue/KPIs; Investigate Case/Evidence/Finding/metrics; Studio Workflow/Tool/Automation Run/metrics; Endpoint technical primitives/local audit/runtime metrics.

## Module coverage
1. Response Inbox — GOV-1.
2. Action Center — GOV-1.
3. Decision Register — GOV-1.
4. Policy Gates — GOV-1.
5. Approvals & Authorities — GOV-1.
6. Playbooks — GOV-2.
7. Runs & Rollback — GOV-2.
8. Audit Trail — GOV-3 CAP-GOV-034..038.
9. Response Metrics — GOV-3 CAP-GOV-039..047.

All nine modules have provider-neutral capability coverage.

## Audit/metrics invariants
Audit Trail ≠ Trace/Activity/raw log/SIEM event. Reconstruction ≠ execution. Completeness ≠ truth. Gap/contradiction ≠ wrongdoing/falsity automatically. Integrity requirement ≠ implemented cryptographic proof. Audit Evidence Package ≠ canonical Evidence.

Metric ≠ objective/Policy/SLO/KPI automatically. Count ≠ quality; throughput ≠ effectiveness; faster Decision ≠ better Decision; Policy block count ≠ prevented incidents; runtime success ≠ verified success; rollback rate ≠ failure rate; Result success ≠ business value. Trend ≠ causal explanation; anomaly ≠ control failure; dashboard ≠ source of truth.

## AI / sensitive data
AI is optional/proposal-only and cannot invent events, alter history, declare fraud/violation as fact, mutate Decision/Result, hide contradictions, change thresholds, publish externally or apply an improvement. Raw secret values remain excluded; identity/tenant/target/exception/emergency dimensions remain permission-aware.

## Verified GOV-3 publication
- baseline: `36edacb4eb374e0b56d6c9e9c45931fdb1e0af20`;
- fifth functional SHA: `042f70d3cfd13467acc294bfff726edde9e16cb0`;
- five commits ahead / zero behind, same merge base;
- 200/200 mandatory gates PASS;
- Requirements: 122 = 99/20/3/0; OPEN: 18;
- global capabilities: 317 = 27 Command + 243 Investigate + 47 Govern; 315 defined / 2 proposed / 317 planned; 8559 sections / 1902 tables.

## Boundary / next roadmap
Documentary PASS does not mean implementation complete. No SIEM, audit/metrics engine, warehouse/storage schema, API/protocol, final RBAC/retention policy, detailed screen rewrite, product code or external compliance certification is claimed.

Next verified historical candidate: **Delivery Roadmap Phase 5 — Studio and Endpoint**, id `roadmap-phase-5-studio-and-endpoint`, title `Phase 5 Studio And Endpoint`. It is identified only and **NOT STARTED** by GOV-3.