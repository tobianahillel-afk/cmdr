---
id: investigate-cloud-analysis-user-questions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-014, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-012]
---
# Cloud Analysis user questions

The module must let an authorized analyst answer, with explicit uncertainty:

1. Why was the analysis opened and what is the return origin?
2. Which organization, tenant, account, subscription, project, region and period are included or excluded?
3. Which providers/services are represented, configured, accessible and actually covered?
4. Which sources, schemas, periods and event families are absent, stale, late or restricted?
5. Which resources and states were observed, when and by which source?
6. Which identities and principals are involved, and where are aliases or human attribution uncertain?
7. Which grants, denies, inheritance and trust relations support an effective-permission candidate?
8. Which audit events occurred and which actor/action/resource/result fields are only candidates?
9. Which configurations changed, and is the difference current, stale, partial or unsupported?
10. What is known about compute, containers, orchestration and serverless workloads without executing commands?
11. Which network exposures and connections are configured versus observed?
12. Which storage and data-access observations exist without unauthorized content access?
13. Which sensitive-material candidates exist, and which access level is authorized?
14. Which anomalies and competing Hypotheses are supported or contradicted?
15. What is the Cloud Timeline, where are its gaps and how does it correlate with other sources?
16. Which Evidence/Finding/Detection/TI/Collection handoffs can be prepared?
17. Can the analysis be reproduced with the available versions, permissions and Tools?
