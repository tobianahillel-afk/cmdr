---
id: investigate-disk-filesystem-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Workflows

| Source | Trigger | Destination | Ownership | Context and return |
|---|---|---|---|---|
| Case / Collection Job / Artifact | open Disk Forensics | CAP-INV-363 | Investigate | tenant, Case, source, Disk Image, objective, Hypothesis, return origin |
| Collection Job | result received | Disk Image / CAP-INV-365 | acquisition owners unchanged | request, Job, Endpoint/source, errors, custody and Artifact relation |
| Disk Image | review source | CAP-INV-365 | Investigate | acquisition context, transformations, transfers, verification and restrictions |
| Integrity Review | accept with limits | CAP-INV-364 | Investigate | image, exploitability, owner, scope and inherited limitations |
| Session | identify structures | CAP-INV-366 | Investigate using Studio Tools | image, candidates, Tool/version, restrictions |
| Filesystem | navigate | CAP-INV-367 | Investigate | volume/filesystem selection, confidence, partiality and return |
| Filesystem Entry | inspect | CAP-INV-368 | Investigate | path projection, metadata, content permissions and source |
| Deleted Entry / Unallocated Source | analyze or carve | CAP-INV-369/375 | Investigate | source range, reconstruction quality, scope and restrictions |
| Journal Record | correlate | CAP-INV-370/374 | Investigate + Shared Timeline | record, timestamp quality, file relation and gaps |
| System/User/Startup Artifact | interpret | CAP-INV-371/372/373 | Investigate | source, age, identity context, contradictions and privacy |
| Disk observation | correlate Memory | CAP-INV-347..362 | Investigate | Artifact relation, time quality and source distinction |
| File / Derived Artifact | handoff | Static or Reverse | destination owner | parent source, extraction, restrictions, objective and return |
| Disk observation | candidate package | CAP-INV-379 → 107/108/109 | owner capability retains qualification | observations, contradictions, provenance and uncertainty |
| Persistent network artifact | future handoff | Phase 4B.2B.3B.2 | future owner | persistent source, candidate configuration, no connection claim |
| Any real-target request | explicit request | Collection/Live Response + Govern | owners retained | target, desired outcome, risk, Case and reason |
