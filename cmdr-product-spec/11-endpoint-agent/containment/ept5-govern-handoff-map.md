---
id: endpoint-ept5-govern-handoff-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Govern Handoff Map

`Finding/context → Action Request → Policy/Authority → Approval where required → Decision → Response Run/Step → CAP-EPT-065/066 → effect primitive → CAP-EPT-073 → CAP-EPT-074/075 → optional CAP-EPT-076/077/078 → CAP-EPT-079 → Govern Verification/Reconciliation → Result`.

Endpoint never creates Approval, Decision, Response Run or Result. Endpoint Technical Verification supplies target-side observations only. Govern decides whether response objectives are met and owns rollback/recovery governance.