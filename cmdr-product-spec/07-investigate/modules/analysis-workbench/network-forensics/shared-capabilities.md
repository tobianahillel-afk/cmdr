---
id: investigate-network-forensics-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Shared capability consumption

| Shared capability | Usage Network Forensics | Non-redéfinition |
|---|---|---|
| Background Jobs | progression des reconstructions/extractions | moteur de Job reste Shared |
| Notifications | fin, erreur, permission ou policy block | centre de notification reste Shared |
| Trace / Activity / Audit Hooks | provenance et actions | records propriétaires restent Shared |
| Entity / Entity Resolution | candidates et relations | aucune fusion locale silencieuse |
| Graph Engine | visualisation relationnelle | alternative tabulaire obligatoire |
| Timeline Engine | ordering et corrélation | Network Timeline reste projection |
| Search | packets, flows, conversations et observations | pas une seconde Event Search |
| Object Linking | Case, Artifact, Entity, Evidence | ownership inchangé |
| Versioning / Recovery | sessions, assessments et disposition | contrats futurs |
| Export / Reporting | sorties autorisées | Report canonique reste owner Reporting |
| Collaboration / Inspector / Context Bar | revue et contexte | rendu Design System/Shared |

Network Forensics consomme ces capacités sans créer de doublon.
