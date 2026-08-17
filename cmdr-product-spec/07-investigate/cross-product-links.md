---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018, OPEN-019]
---
# Cross-product links — Detection Engineering, Threat Intelligence, Cloud Analysis and Mobile Forensics

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
| Case/Incident/Finding/Hunt/Signal/Collection Job/Artifact/Mobile Evidence Package → CAP-INV-701 | objective, device/source candidates, period, restrictions, acquisition context and return origin | source owner retained; no automatic analysis/acquisition |
| CAP-INV-701/704/705 → Collection | missing/partial source, acquisition/custody gap and complementary Collection Request context | Collection owns acquisition execution/results/custody |
| CAP-INV-703 → Platform Settings / Endpoint Agent | configured Fleet/management/source and declared capability projections | Settings/Endpoint ownership retained; Mobile does not assume agent presence |
| CAP-INV-706/715 → Disk and Filesystem Forensics | represented storage/recovery observations requiring deeper disk analysis | Mobile package only; Disk owner retains specialized analysis |
| CAP-INV-713 ↔ Network Forensics | stored mobile network/Wi-Fi/Bluetooth/SIM/endpoint observations and authorized capture refs | no active scan; Mobile network state ≠ full Network Forensics |
| CAP-INV-714 → Cloud Analysis | cloud-backed/synchronized artifact and question | cloud-backed record ≠ full Cloud Analysis; Cloud owns its Session |
| CAP-INV-718 → Static/Reverse/Dynamic | Derived Artifact, source lineage, question and restrictions | destination owns technical analysis; Derived Artifact ≠ Evidence |
| CAP-INV-718 → Evidence/Findings | Evidence Candidate Package and Finding Draft | destination owner qualifies/confirms; nothing automatic |
| CAP-INV-718 → Detection Engineering | sourced Mobile Detection Engineering Package | no Detection Content/rule/promotion/deployment automatic |
| CAP-INV-718 → Threat Intelligence | sourced Mobile observations and Intelligence Handoff Package | no attribution/knowledge promotion automatic |
| CAP-INV-712/718 → Govern | masked sensitive/response context or real-device action request | no secret use, revocation, isolation, lock, wipe or other device action by Investigate |
| CAP-INV-701..719 ↔ Studio | attributed Tools/Tool Calls/Automation Runs for bounded analysis | Studio owns execution lifecycle; no silent Tool run |
| CAP-INV-716/719 ↔ Shared Timeline/Graph/Trace/Versioning/Export | mobile chronology, candidate relations, provenance and minimized exports | Shared owns generic mechanisms; Mobile semantics remain local |

All transitions preserve tenant, environment, device/package/source scope, versions, markings/classification, permissions, masking, errors, partiality, timestamps, source owner, authority and return origin. Cross-tenant/cross-device links never extend access. No transition silently creates a Signal, Alert, Incident, Evidence, Finding, Detection rule, active Indicator, Cloud permission, secret use, Collection Job, target/device mutation or response.
