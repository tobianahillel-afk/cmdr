---
id: investigate-mobile-forensics-user-questions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020]
open_decisions: [OPEN-011, OPEN-013, OPEN-014]
---
# Mobile Forensics User Questions

Authorized users must be able to answer, with sources and uncertainty visible:

1. Which device or device candidate is in scope, and what declared/candidate platform and version apply?
2. What evidence representation is available, who acquired or transferred it, under what authority, and what is missing?
3. Is integrity verified, partially verified or unverified; is the representation complete, partial, corrupted, encrypted, locked, restricted or inaccessible?
4. Which partitions, filesystems, storage areas, directories, application areas and residual regions are actually represented?
5. Which applications/packages are present, which versions and permissions are declared, and what contradictions exist without assuming usage?
6. Which application records, databases, caches, preferences, browser/download data and deleted candidates exist without inferring user intent?
7. Which messaging, call, contact, email or communication records exist, and what attribution/delivery/read-state limits remain?
8. Which media/documents and metadata exist, and which creation/import/export relations are only candidates?
9. Which location, movement, navigation, sensor or health/activity records exist, and what prevents treating device location as certain user presence?
10. Which account/token/key/secret candidates exist, what is masked, and what access level is authorized without using credentials?
11. Which Wi-Fi, Bluetooth, NFC, SIM/eSIM, carrier, interface, endpoint or paired-device records exist without active interaction?
12. Which local/synchronized backups and cross-device/cloud-backed relations exist, and which records cannot be proven local/current?
13. Which deleted/residual/recovered/carved records exist and what recovery quality or attribution limits apply?
14. What sequence can be reconstructed in the Mobile Timeline, with observed/recorded/reconstructed/estimated/synchronized/absent timestamps distinguished?
15. Which anomalies or persistence/compromise hypotheses are supported, contradicted or inconclusive without automatic malware/compromise confirmation?
16. Which bounded source records should become Derived Artifacts or candidate handoff packages, and which destination owns qualification?
17. Can another authorized reviewer reproduce the analytical reasoning from source package, Tools/Runs, parameters, permissions, versions, errors and human decisions?

These questions are functional. They do not imply a particular mobile platform, commercial extraction tool, acquisition method or implementation architecture.
