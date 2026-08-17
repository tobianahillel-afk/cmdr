---
id: investigate-disk-filesystem-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Concepts and mandatory distinctions

| Concept A | Distinct from | Product rule |
|---|---|---|
| Disk Acquisition Request | Disk Image | request/execution context is not the acquired source |
| Disk Image | Collection Job / Disk Session / Memory Image / live filesystem / Evidence | no synonym or lifecycle merge |
| Partition | Volume | physical/logical interpretations remain explicit |
| Volume | Filesystem | a volume may be unsupported, empty or ambiguously interpreted |
| Filesystem | directory tree | tree is a projection, not filesystem identity |
| Filesystem Entry | canonical Artifact | promotion or relation is explicit |
| active file | deleted entry / recovered content | states and sources remain visible |
| unallocated data | identified deleted file | attribution requires evidence |
| residual space | attributed content | uncertainty is mandatory |
| carving result | original file | name, path, completeness and attribution are not guaranteed |
| metadata timestamp | certain event | type, source, timezone and uncertainty are shown |
| journal record | user action / Disk Timeline / Case Timeline | no automatic authorship or timeline replacement |
| user/application artifact | intention or author | privacy and uncertainty are mandatory |
| startup artifact | confirmed persistence | candidate until corroborated |
| encrypted/inaccessible content | malicious or absent content | no bypass or adverse inference |
| file path | stable identity | versions, content indicators and provenance remain separate |
| Disk Timeline | Memory Timeline / Case Timeline | correlation without replacement |
| Disk Forensics | Memory or Network Forensics | owners and analytical sources remain distinct |
| Tool/AI output | conclusion, Evidence, Finding or rule | human review and owner handoff required |
