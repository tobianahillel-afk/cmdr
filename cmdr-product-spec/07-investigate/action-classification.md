---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-011, OPEN-012, OPEN-013, OPEN-019]
---
# Action classification — Investigate through Phase 4B.4B

| Class | Meaning | Threat Intelligence Analysis/Products examples | Cloud Analysis examples | Mobile Forensics examples | Local execution |
|---:|---|---|---|---|---|
| 0 | observation | read/filter/search/compare Sessions, hypotheses, assessments, products, publications, watchlist projections, feedback and provenance | read/navigate/search/filter/compare/inspect scope, inventory, identities, IAM, events, configurations, workloads, network, storage, timeline and provenance | read/navigate/search/filter/compare device/package metadata, filesystem/app records, communications, media metadata, location/sensor metadata, connectivity, backups, timeline and provenance; masked-presence views | allowed by permission |
| 1 | bounded analytical execution | structured analysis, source fusion, comparison, deterministic quality checks, monitoring request/run, package generation and authorized export | correlation, reconstruction, deterministic parsing/comparison, derived-result/package generation and authorized minimized export | deterministic parsing/comparison, bounded authorized extraction/recovery from represented evidence, time normalization, correlation, Derived Artifact generation and authorized minimized export | explicit, bounded, attributed and reversible |
| 2 | reversible analytical mutation | create/update/review/version Session, hypothesis, Product Draft, Release Recommendation, Dissemination Plan, Internal Publication Record, Watchlist Definition, handoff, feedback, correction/retraction and improvement package | create/update/close/reopen Session, select scope, annotate/dispute anomaly or Hypothesis, prepare handoff, request sensitive access, version/supersede local concepts | create/update/close/reopen Mobile Session, select device/platform candidate and scope, annotate/bookmark/dispute observations, manage Hypotheses, request sensitive reveal/export under policy, prepare handoffs and supersede local concepts | Investigate-owned concepts; OPEN-013/019 |
| 3 | governed production/external action | external sharing, client/public publication, watchlist activation, Indicator deployment, runtime change, global revocation/recall | permission or policy change, credential revocation, resource isolation, network blocking, workload shutdown, collection execution or response | real remote acquisition, device unlock/bypass, profile installation, configuration change, account/credential revocation, device isolation/lock, MDM action or response | not executed by Investigate; Collection/Govern/runtime owner |
| 4 | destructive/irreversible | provenance/history destruction or non-revocable sensitive disclosure | resource/data deletion, destructive remediation, trace destruction or irreversible secret/data disclosure | device wipe, destructive remediation, source deletion, trace destruction or irreversible sensitive disclosure | denied by default or strictly governed |

CAP-INV-001..719 expose only classes 0–2 locally. Internal publication in CAP-INV-529 remains policy-bound, reversible and tenant-scoped; it is not class-3 external sharing. Cloud and Mobile sensitive-material presence, metadata, masked preview, reveal, copy and export remain separate permissions; **secret use is prohibited**. Mobile class-1 extraction/recovery is bounded derivation from already authorized represented evidence, never acquisition from or mutation of the original device. Any class-3/4 action is blocked or prepared as a Collection/Govern/owner request. No Signal, Alert, Incident, rule, active watchlist, deployed Indicator, Cloud permission, mobile device configuration, credential, source or provenance is silently mutated.
