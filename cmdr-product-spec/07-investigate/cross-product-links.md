---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Investigate through Phase 4B.2 closure

| Transition | Context | Ownership / return |
|---|---|---|
| Case → Network Intake | Case, Hypothesis, objective, restrictions and return origin | Investigate |
| Network Capture Request → Collection Job | target projection, scope, authorization, duration and limits | Collection/Endpoint execution; request context remains Investigate |
| Collection Job → Capture Artifact | result, window, losses, errors, custody and lineage | Artifact remains Investigate; producer attribution retained |
| Capture Artifact → Integrity/Coverage/Session | source, version, acquisition, sensor, scope, restrictions and return origin | Investigate analysis; Settings/producer ownership unchanged |
| Capture → Packet → Flow → Session → Conversation | selected source, timebase, coverage, uncertainties and Tool context | Investigate observations; no certainty promotion |
| Conversation → DNS / Transaction / Encrypted Metadata / Transfer | messages, direction candidates, gaps, permissions and return origin | Investigate sub-analysis |
| Extracted Object → Static / Reverse / Sandbox | Derived Artifact, parent, restrictions, Case and provenance | Investigate; source immutable; destination owner preserved |
| Network Observation → Entity / Graph | addresses, names, certificates, sources, timestamps and confidence | Shared owns Entity/Graph; no silent merge |
| Network Observation → Memory / Disk / Endpoint / Event Search / Sandbox | selected observations, time quality, source and limitations | each destination retains ownership; return origin preserved |
| Network observations → Network Timeline | timestamp type, timezone, sensor, coverage and conflicts | Shared Timeline mechanism; Network projection remains Investigate |
| Network result → Evidence / Finding | selected records, contradictions, lineage, restrictions and uncertainty | CAP-INV-107/108/109 retain qualification |
| Network knowledge → future Detection Engineering | behavior, conditions, sources, limits and reproducibility | future Phase 4B.3; no rule created, tested or deployed |
| Network knowledge → future Intelligence | candidate domain, certificate, endpoint, relations and provenance | future Phase 4B.3; no Indicator, Campaign or Threat Actor created |
| Govern Result → Case | returned outcome, authority, scope, status and evidence | Govern owns Result; Investigate consumes projection |
| Any real-target or active network request | target, requested outcome, risk and Case | block or route to Collection/Live Response and Govern |

Tenant, environment, selection, permissions, masking, minimization, errors and return origin are preserved. Network Forensics is not Event Search, SIEM monitoring, Detection Engineering or Intelligence.
