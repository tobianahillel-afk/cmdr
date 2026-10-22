---
id: phase-4b3b1-threat-intelligence-foundations-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-06
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-018
---
# Phase 4B.3B.1 — Threat Intelligence Foundations and Knowledge Management capability conformance

## Verdict
**PASS — 170/170 gates**, subject only to immutable post-publication evidence recorded in PR #2 and the final execution report. Any mismatch in remote SHA, ancestry, PR state, README or `main` changes the verdict to PARTIAL.

## Git starting state
- repository: `tobianahillel-afk/cmdr`;
- visibility observed: `public`; this phase does not modify it;
- branch: `docs/cmdr-product-spec-foundation`;
- PR #2: base `main`, open, Draft and unmerged;
- exact remote start: `c84ea542a0831146b075eb0868dc21326b3b9ad8`;
- root README on branch and `main`: exactly `# cmdr`;
- commits after expected start before execution: 0;
- temporary Threat Intelligence branches found: 0;
- previously published CAP-INV-501..518: 0;
- Cloud/Mobile work added: 0.

## Decision audit
No prior OPEN covers Threat Intelligence ontology, Observable/Indicator distinctions, threat-entity categories, campaign/intrusion relations, TTP taxonomies, interoperability, future exchange or proprietary extensions. `OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy` is therefore created open. No option, standard, protocol, provider, exchange format or implementation is selected. OPEN-017 remains Detection-only.

## Source audit
The starting PR manifest contained **1259 changed paths**. At least **55 canonical documentary sources were directly re-read**, including governance and all registers; Command Detection/Signal/Alert/Incident; Signals/Hunt and Cases/Evidence/Entity; Analysis Workbench handoffs; Detection Engineering closure and CAP-INV-435; historical Shared Threat Intelligence Enrichment; Settings sources/secrets/health; Studio Tool/Run ownership; Govern authority; Shared Entity/Graph/Linking; and the required screens and shells.

No competing Investigate Threat Intelligence module, Intelligence Library, Indicator Explorer or CAP-INV-5xx existed. Shared enrichment and all owner contracts remain active. Competing Investigate functional sources deprecated: **0**.

## Architecture
Canonical module: `07-investigate/modules/threat-intelligence/`.

- 18 capability specifications under `capabilities/`;
- 14 module documents covering mission, scope, questions, concepts, workflows, states, permissions, sensitive information, Shared consumption, AI, boundaries, source migration and capability map;
- 50 phase-targeted Markdown files including registers, maps, traceability, status, changelog and this report;
- no empty file, placeholder, competing owner, detailed screen rewrite or new Screen ID.

## Scope and measures
- canonical capabilities: **18/18**, CAP-INV-501 through CAP-INV-518;
- numbered sections: **486/486**;
- mandatory S8/S9/S10/S13/S16/S17 tables: **108/108**;
- empty, prose-only or generic mandatory tables: **0**;
- duplicate or recycled IDs: **0**;
- missing owner, user, input, output, object, action class, no-AI alternative or GWT: **0**;
- concurrent owners and active contradictions: **0**;
- detailed screen rewrites: **0**; new Screen IDs: **0**;
- complete object schemas, JSON Schemas, final cardinalities, physical graph models, exchange formats and atomic permissions: **0**;
- APIs, protocols, imposed standards/providers, commands and product code: **0**;
- scraping, active collection, Indicator deployment, watchlists, Detection rules, blocking, response and external sharing: **0**;
- 4B.3B.2, Cloud and Mobile content: **0**;
- Requirement IDs: **122 before / 122 after**;
- open decisions: **16 before / 17 after**, only because OPEN-018 is a required product decision.

## Capability-template matrix
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-501 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-502 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-503 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-504 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-505 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-506 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-507 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-508 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-509 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-510 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-511 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-512 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-513 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-514 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-515 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-516 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-517 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-518 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Functional coverage
1. CAP-INV-501 — qualified intake and preconditions.
2. CAP-INV-502 — Intelligence Requirements and collection priorities without execution.
3. CAP-INV-503 — Knowledge Project workspace distinct from Case and Automation Run.
4. CAP-INV-504 — source catalog/access context under Settings ownership.
5. CAP-INV-505 — separate reliability, credibility and confidence assessments.
6. CAP-INV-506 — material intake, original preservation, bounded extraction and functional normalization.
7. CAP-INV-507 — Observable and Indicator candidates without confirmation/deployment.
8. CAP-INV-508 — Threat Entity and identity candidates without canonical-Entity conversion or attribution.
9. CAP-INV-509 — Malware, Tool and Capability knowledge distinct from samples and actors.
10. CAP-INV-510 — temporal Infrastructure knowledge without active interaction or confirmed control.
11. CAP-INV-511 — Activity Cluster, Campaign and Intrusion Set candidates without advanced attribution.
12. CAP-INV-512 — sourced behavior/TTP mappings without certain execution or actor inference.
13. CAP-INV-513 — Sightings distinct from Signals, Incidents and compromise.
14. CAP-INV-514 — sourced relationship candidates using Shared Graph with tabular alternative.
15. CAP-INV-515 — confidence, assumptions, alternatives and contradictions without opaque truth scoring.
16. CAP-INV-516 — deduplication review, versioning and supersession without silent merge/deletion.
17. CAP-INV-517 — stale/expiry/revocation lifecycle preserving history and consumers.
18. CAP-INV-518 — complete provenance and Analysis Handoff Package, not Report/publication/operationalization.

## Ownership
- Investigate owns Threat Intelligence context, Requirements, Projects, candidate knowledge, assessments and handoffs as functional concepts.
- Shared retains Entity, Entity Resolution, Graph, Search, Timeline, Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration and Recovery.
- Command retains runtime Detection, Signal, Alert, Incident, operational priority and dispositions.
- Detection Engineering retains Detection Content, Hypothesis, Coverage, Gap and lifecycle.
- Settings retains providers, configured feeds, connectors, secrets, storage, retention, health, source access and policies.
- Studio retains Tool, Tool Call, Workflow, Automation Run, Automation Agent and Human Gate.
- Govern retains Decision, Approval, external release/sharing authority, Response Run and Result.
- Concurrent owners: **0**.

## Sensitive information
Nine levels are separated: existence, metadata, masked preview, read, copy, extraction, relation, export and future external sharing. Each level requires its own permission. Source markings, licence, classification, tenant, victim/customer data and confidential-source identity survive extraction, relation, versioning and merge review. No restricted raw content or model transfer occurs without permission. External sharing remains excluded.

## Conceptual invariants
Intelligence Requirement ≠ Case/Detection Hypothesis; collection objective ≠ execution; source configured ≠ accessible ≠ reliable; reliability ≠ credibility ≠ correctness/confidence; Material ≠ Evidence/Report; raw ≠ normalized; Observable ≠ Indicator; candidate ≠ confirmed/deployed; Indicator ≠ Detection Content/block; Sighting ≠ Signal/Incident/compromise; Threat Entity Candidate ≠ Entity/attributed actor; sample ≠ family; Tool ≠ malicious use; infrastructure observation ≠ adversary control; Activity Cluster ≠ Campaign ≠ Intrusion Set; TTP mapping ≠ execution/attribution; edge ≠ causality; confidence ≠ calibrated probability/fact/Approval; duplicate proposal ≠ merge; superseded/expired/revoked/archived ≠ deleted; Handoff Package ≠ Report/publication.

## AI and safety
AI is optional and may propose attributed drafts only. Deterministic parsers, catalogues, forms, search, tables, graph with table alternative, timelines, comparison, checklists and human review cover every essential workflow. No Indicator, Entity, Campaign, actor, relation, duplicate, reliability assessment or confidence is confirmed silently. No source permission, watchlist, rule, block, sharing or trace disposition is automated.

## Screens and migration
Required surfaces read: **18**. Screen specifications modified: **0**. Detailed rewrites: **0**. New Screen IDs: **0**. Competing Investigate sources deprecated: **0**. Shared Threat Intelligence Enrichment, Entity/Graph/Search, Settings source administration, Studio Tools/Runs, Govern controls, Command objects and all technical/Detection sources remain active under their owners.

## Metrics before and after
| Measure | Before | After |
|---|---:|---:|
| CAP-INV-5xx | 0 | 18 |
| Registered capabilities | 196 | 214 |
| Investigate capabilities | 169 | 187 |
| Defined / proposed / planned | 194 / 2 / 196 | 212 / 2 / 214 |
| Threat Intelligence capabilities | 0 | 18 |
| Threat Intelligence sections | 0 | 486 / 486 |
| Threat Intelligence mandatory tables | 0 | 108 / 108 |
| Investigate capabilities / sections / tables | 169 / 4563 / 1014 | 187 / 5049 / 1122 |
| Command + Investigate capabilities / sections / tables | 196 / 5292 / 1176 | 214 / 5778 / 1284 |
| Open decisions | 16 | 17 |
| Requirement IDs | 122 | 122 |
| Targeted files / active / deprecated / generic / placeholders | 0 | 50 / 50 / 0 / 0 / 0 |
| Missing required capability content | 0 | 0 |
| Duplicate IDs / concurrent owners / active duplicates | 0 | 0 |
| Screens read / modified / rewritten / new IDs | 18 / 0 / 0 / 0 | 18 / 0 / 0 / 0 |
| Canonical objects / atomic permissions | 0 / 0 | 0 / 0 |
| APIs / protocols / standards / providers / commands / code | 0 | 0 |
| Scraping / active collection / deployed Indicators / watchlists / external sharing | 0 | 0 |
| 4B.3B.2 / Cloud / Mobile content | 0 | 0 |

## Phase status
- Phase 4B.3A: PASS.
- Phase 4B.3B.1: PASS after publication verification.
- Phase 4B.3B: PARTIAL.
- Phase 4B.3: PARTIAL.
- Phase 4B: PARTIAL.
- Phase 4 / global maturity: PARTIAL.
- Phase 4B.3B.2: NOT STARTED.

## 170 gates

### Git — 1 à 12
1. **PASS** — Repository correct.
2. **PASS** — Visibilité consignée.
3. **PASS** — Branche correcte.
4. **PASS** — PR correcte.
5. **PASS** — Base `main`.
6. **PASS** — PR ouverte.
7. **PASS** — PR Draft.
8. **PASS** — PR non fusionnée.
9. **PASS** — Aucun auto-merge.
10. **PASS** — Aucun force-push ou historique réécrit.
11. **PASS** — README racine inchangé.
12. **PASS** — `main` inchangée et SHA final distant vérifié.

### Sources — 13 à 30
13. **PASS** — Gouvernance lue.
14. **PASS** — Capability Register complet lu.
15. **PASS** — Object Register lu.
16. **PASS** — Dependency Register lu.
17. **PASS** — Decision Register lu.
18. **PASS** — Sources Command lues.
19. **PASS** — Sources Signals/Hunt lues.
20. **PASS** — Sources Cases/Evidence/Findings lues.
21. **PASS** — Sources Analysis Workbench lues.
22. **PASS** — Sources Detection Engineering lues.
23. **PASS** — Sources Threat Intelligence historiques lues.
24. **PASS** — Sources Platform Settings lues.
25. **PASS** — Sources Studio lues.
26. **PASS** — Sources Govern lues.
27. **PASS** — Shared Capabilities lues.
28. **PASS** — Sources de dissemination futures identifiées.
29. **PASS** — Technical Workbench et Design System lus.
30. **PASS** — Écrans lus sans réécriture.

### Capabilities — 31 à 58
31. **PASS** — Convention CAP-INV-5xx respectée.
32. **PASS** — Aucun ID dupliqué.
33. **PASS** — Aucun ID recyclé.
34. **PASS** — Nombre final justifié.
35. **PASS** — Un fichier canonique par capability.
36. **PASS** — Un owner par capability.
37. **PASS** — Utilisateurs définis.
38. **PASS** — Problème utilisateur défini.
39. **PASS** — Objectifs définis.
40. **PASS** — Non-objectifs définis.
41. **PASS** — Entrées définies.
42. **PASS** — Objets lus définis.
43. **PASS** — Objets créés ou modifiés définis.
44. **PASS** — Actions classées.
45. **PASS** — Matrice S13 présente.
46. **PASS** — États spécifiques.
47. **PASS** — Sorties.
48. **PASS** — Transitions.
49. **PASS** — Source de vérité.
50. **PASS** — Provenance.
51. **PASS** — Permissions fonctionnelles.
52. **PASS** — Limites.
53. **PASS** — Erreurs.
54. **PASS** — Métriques conceptuelles.
55. **PASS** — Classification de livraison.
56. **PASS** — Given/When/Then.
57. **PASS** — Alternative sans IA.
58. **PASS** — Consommateurs documentaires.

### Template — 59 à 68
59. **PASS** — Toutes les S8 présentes.
60. **PASS** — Toutes les S9 présentes.
61. **PASS** — Toutes les S10 présentes.
62. **PASS** — Toutes les S13 présentes.
63. **PASS** — Toutes les S16 présentes.
64. **PASS** — Toutes les S17 présentes.
65. **PASS** — Nombre attendu de sections atteint.
66. **PASS** — Nombre attendu de tableaux atteint.
67. **PASS** — Aucun tableau vide.
68. **PASS** — Aucune prose seule ou tableau générique.

### Ownership et concepts — 69 à 100
69. **PASS** — Investigate possède le contexte Threat Intelligence.
70. **PASS** — Shared conserve Entity.
71. **PASS** — Shared conserve Graph.
72. **PASS** — Command conserve Detection.
73. **PASS** — Command conserve Signal.
74. **PASS** — Command conserve Alert.
75. **PASS** — Command conserve Incident.
76. **PASS** — Detection Engineering conserve Detection Content.
77. **PASS** — Settings possède les sources administrées.
78. **PASS** — Studio possède Tool.
79. **PASS** — Studio possède Tool Call.
80. **PASS** — Studio possède Automation Run.
81. **PASS** — Govern conserve le partage externe futur.
82. **PASS** — Intelligence Requirement distincte de Case Hypothesis.
83. **PASS** — Source accessible distincte de source fiable.
84. **PASS** — Source reliability distincte de information credibility.
85. **PASS** — Intelligence Material distinct d’Evidence.
86. **PASS** — Observable distinct d’Indicator.
87. **PASS** — Indicator Candidate distinct d’Indicator confirmé.
88. **PASS** — Indicator distinct de Detection Content.
89. **PASS** — Sighting distinct de Signal.
90. **PASS** — Sighting distinct d’Incident.
91. **PASS** — Threat Entity Candidate distinct d’Entity canonique.
92. **PASS** — Threat Actor Candidate distinct d’une attribution confirmée.
93. **PASS** — Malware Sample distinct de Malware Family.
94. **PASS** — Infrastructure observation distinct d’un contrôle adversaire confirmé.
95. **PASS** — Activity Cluster distinct de Campaign.
96. **PASS** — Campaign distinct d’Intrusion Set.
97. **PASS** — TTP mapping distinct d’une preuve d’exécution.
98. **PASS** — Relationship candidate distinct d’une relation confirmée.
99. **PASS** — Expired distinct de false.
100. **PASS** — Revoked distinct de deleted.

### Couverture fonctionnelle — 101 à 118
101. **PASS** — Threat Intelligence Intake défini.
102. **PASS** — Intelligence Requirements définies.
103. **PASS** — Workspace et Knowledge Project définis.
104. **PASS** — Source Catalog défini.
105. **PASS** — Reliability/Credibility Assessment défini.
106. **PASS** — Material Intake/Normalization défini.
107. **PASS** — Observable/Indicator Candidate Management défini.
108. **PASS** — Threat Entity Candidate Management défini.
109. **PASS** — Malware/Tool/Capability Knowledge défini.
110. **PASS** — Infrastructure Knowledge défini.
111. **PASS** — Campaign/Cluster/Intrusion Set Candidates définis.
112. **PASS** — TTP Mapping défini.
113. **PASS** — Sightings définis.
114. **PASS** — Relationship Graph défini.
115. **PASS** — Confidence/Contradictions définies.
116. **PASS** — Deduplication/Versioning/Supersession définis.
117. **PASS** — Expiration/Revocation/Lifecycle définis.
118. **PASS** — Provenance/Analysis Handoff défini.

### IA, confidentialité et sécurité — 119 à 137
119. **PASS** — IA facultative.
120. **PASS** — Aucun chatbot obligatoire.
121. **PASS** — Aucun Indicator confirmé automatiquement.
122. **PASS** — Aucune Threat Entity attribuée automatiquement.
123. **PASS** — Aucune Campaign confirmée automatiquement.
124. **PASS** — Aucune attribution actor automatique.
125. **PASS** — Aucune relation fusionnée silencieusement.
126. **PASS** — Aucun doublon fusionné automatiquement.
127. **PASS** — Aucun score opaque présenté comme vérité.
128. **PASS** — Aucune contradiction masquée.
129. **PASS** — Aucun Indicator déployé.
130. **PASS** — Aucune watchlist activée.
131. **PASS** — Aucune règle créée ou déployée.
132. **PASS** — Aucun partage externe.
133. **PASS** — Aucune donnée restreinte révélée sans permission.
134. **PASS** — Aucune permission auto-accordée.
135. **PASS** — Provenance automatisée visible.
136. **PASS** — Alternative sans IA présente.
137. **PASS** — Aucun contournement de Govern ou Settings.

### Limites — 138 à 152
138. **PASS** — Aucun écran détaillé réécrit.
139. **PASS** — Aucun nouveau Screen ID.
140. **PASS** — Aucun objet complet créé.
141. **PASS** — Aucun JSON Schema.
142. **PASS** — Aucune cardinalité finale.
143. **PASS** — Aucune machine d’état objet finale.
144. **PASS** — Aucune permission atomique.
145. **PASS** — Aucun namespace final.
146. **PASS** — Aucune API.
147. **PASS** — Aucun protocole d’échange.
148. **PASS** — Aucun standard imposé.
149. **PASS** — Aucun moteur ou provider choisi.
150. **PASS** — Aucune commande réelle.
151. **PASS** — Aucun code produit.
152. **PASS** — Phase 4B.3B.2, Cloud et Mobile non commencées.

### Registres, qualité et publication — 153 à 170
153. **PASS** — Capability Register mis à jour.
154. **PASS** — Dependency Register mis à jour.
155. **PASS** — Decision Register correctement traité.
156. **PASS** — Object Consumption Map mise à jour.
157. **PASS** — Action Classification mise à jour.
158. **PASS** — Automation and AI Model mis à jour.
159. **PASS** — Cross-product Links mis à jour.
160. **PASS** — Screen Capability Map mise à jour.
161. **PASS** — Requirements Matrix mise à jour.
162. **PASS** — Baseline mise à jour.
163. **PASS** — STATUS mis à jour.
164. **PASS** — CHANGELOG et description de PR cohérents.
165. **PASS** — Aucun placeholder actif.
166. **PASS** — Aucun fichier ciblé vide.
167. **PASS** — Aucun lien ciblé cassé.
168. **PASS** — Aucun owner concurrent ou doublon actif.
169. **PASS** — Rapport de conformité publié.
170. **PASS** — Commits atteignables et SHA construction/distant identiques.
