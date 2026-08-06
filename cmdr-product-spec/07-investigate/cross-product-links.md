---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018, OPEN-019]
---
# Cross-product links — Detection Engineering, Threat Intelligence and Cloud Analysis

| Transition | Context | Ownership / return |
|---|---|---|
| Case/Finding/Hunt/Incident/technical analysis → CAP-INV-501 | need, sources, Evidence/Artifacts, restrictions and return origin | source owner retained |
| CAP-INV-435 → CAP-INV-501/531/537 | Detection improvement context and Indicator/TTP relevance | Detection Engineering keeps its objects |
| CAP-INV-518 → CAP-INV-519 | Analysis Handoff Package, candidates, assessments, gaps, provenance | package ≠ Session/Report |
| CAP-INV-519..524 → CAP-INV-525..527 | structured assessments and alternatives | Investigate product concepts |
| CAP-INV-527 → CAP-INV-528/529 | quality/release recommendation, audience and restrictions | recommendation ≠ publication/Approval |
| CAP-INV-528 → CAP-INV-533 → Govern | candidate recipient, risks, redaction, releasability and Action Request | no transmission; Govern/destination owns release |
| CAP-INV-530 → Settings/Govern/runtime owner | Watchlist Definition and activation request | Investigate never activates; active-observed is projection |
| CAP-INV-531 → Detection Engineering | Indicator package, sources, confidence, limits, expiry | no Detection Content/rule automatic |
| CAP-INV-531 → Command | context package for prioritization/triage | no Signal/Alert/Incident automatic |
| CAP-INV-531 → Settings | destination/target/source request | Settings keeps configuration |
| CAP-INV-532 → consumers | Sighting/change assessment and deterministic notification | change notification ≠ operational Alert |
| CAP-INV-534/535 → Command/Cases/Hunts/Detection/Settings/Studio | feedback, gaps and collection/improvement package | destination reevaluates permissions and authority |
| CAP-INV-536 → consumers/Govern | correction/retraction/supersession and future recall request | history preserved; external recall not executed |
| CAP-INV-537 → Intake/Requirement/Case/Hunt/Detection | continuous-improvement package | no active mutation |
| Case/Incident/Finding/Hunt/Signal/Cloud source/Artifact → CAP-INV-601 | objective, Cloud scope candidates, periods, sources, restrictions and return origin | source owner retained; no automatic analysis |
| CAP-INV-601/602 → Platform Settings | source/access/health/schema/retention gap proposal | Settings owns providers, connectors, credentials and configuration |
| CAP-INV-603..607 ↔ Shared Entity/Graph/Timeline | candidate scope, identity, permission and event relations | Shared owns generic mechanisms; Cloud concepts remain local |
| CAP-INV-609/610 → Endpoint/Memory/Disk/Network owners | workload references, gaps, permissions and acquisition reasons | package only; destination owns collection and analysis |
| CAP-INV-612 ↔ Network Forensics | configured topology/exposure and authorized flow/capture projections | Cloud network analysis ≠ complete Network Forensics |
| CAP-INV-614 → Evidence/Govern | masked Sensitive Material Candidate or credential-revocation Action Request context | no secret use, reveal, copy, export or revocation automatic |
| CAP-INV-615/616 → Cases/Hunts | Cloud Hypothesis, anomalies, timeline and correlation candidates | candidate ≠ Finding, Case Hypothesis acceptance or causality |
| CAP-INV-617 → Evidence/Findings | Evidence Candidate Package and Finding Draft | destination owner qualifies/creates canonical objects |
| CAP-INV-617 → Detection Engineering | Cloud Detection Gap and Engineering Package | no Detection Content, rule, promotion or deployment automatic |
| CAP-INV-617 → Threat Intelligence | sourced Cloud observations and relations | no automatic attribution or canonical TI promotion |
| CAP-INV-617 → Collection/Govern | collection/future-response preparation | no collection, permission change or response execution |
| CAP-INV-618 → Settings/Studio/Shared/source owners | provenance/reproducibility gaps and correction proposals | no active replay, connector repair or history mutation |

All transitions preserve tenant, environment, Cloud scope units, versions, markings, licence, permissions, masking, errors, partiality, timestamps, source owner, authority and return origin. Cross-tenant links never extend access, and no transition silently creates a Signal, Alert, Incident, Evidence, Finding, Detection rule, Cloud permission, secret use or target mutation.
