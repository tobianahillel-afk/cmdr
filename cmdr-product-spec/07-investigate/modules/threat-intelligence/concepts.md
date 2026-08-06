---
id: investigate-threat-intelligence-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
open_decisions:
  - OPEN-018
---
# Concepts and mandatory distinctions

| Concept | Distinction |
|---|---|
| Intelligence Requirement | ≠ Case Hypothesis, Detection Hypothesis or collection execution |
| Source context | configured ≠ accessible ≠ reliable; reliability ≠ credibility ≠ correctness |
| Intelligence Material | ≠ Evidence and ≠ canonical Report |
| Normalized record | raw material ≠ normalized knowledge; extraction ≠ human validation |
| Observable Candidate | ≠ Indicator Candidate and ≠ confirmed Indicator |
| Indicator | ≠ Detection Content, watchlist, block action or confirmed malicious activity |
| Sighting | ≠ Signal, Incident or confirmed compromise |
| Threat Entity Candidate | ≠ Shared Entity and ≠ attributed Threat Actor |
| Malware Sample | ≠ Malware Family; similar ≠ same family |
| Tool | legitimate software ≠ benign use; observation ≠ attribution |
| Infrastructure knowledge | IP/domain/certificate ≠ stable identity or confirmed adversary control |
| Activity Cluster | ≠ Campaign ≠ Intrusion Set ≠ actor |
| TTP Mapping | ≠ certain execution, intent or attribution |
| Relationship | edge/co-occurrence ≠ causality/control/attribution |
| Confidence | ≠ calibrated probability, source reliability, fact or Approval |
| Lifecycle | duplicate candidate ≠ same object; superseded/expired/revoked/archived ≠ deleted |
| Handoff Package | ≠ Intelligence Report, publication, attribution or operationalization |

`OPEN-018` retains ontology and interoperability choices without selecting any standard, protocol or provider.
