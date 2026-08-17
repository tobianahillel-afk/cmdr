---
id: investigate-network-forensics-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Concepts

| Concept | Statut | Owner | Invariant |
|---|---|---|---|
| Network Capture Request | capability context | Investigate Collection | ≠ Capture Artifact |
| Network Capture Artifact | Artifact spécialisé, schéma futur | Investigate | ≠ Collection Job ou Session |
| Network Forensics Session | concept fonctionnel | Investigate | ≠ reconstructed network session |
| Capture Integrity Review | concept fonctionnel | Investigate | validité ≠ représentativité |
| Coverage Assessment | concept fonctionnel | Investigate | gap ≠ silence certain |
| Packet / Frame Observation | concepts fonctionnels | Investigate | packet ≠ frame ≠ flow |
| Flow / Network Session / Conversation | concepts fonctionnels | Investigate | niveaux distincts |
| Protocol / DNS / Transaction / Certificate Observation | concepts fonctionnels | Investigate | candidates et observations, pas certitudes |
| Network Relationship / Entity Candidate | relation analytique | Investigate + Shared Entity | aucune fusion silencieuse |
| Network Anomaly Candidate | concept fonctionnel | Investigate | ≠ Finding |
| Network Timeline Entry | projection Investigate | Shared Timeline mechanism | ≠ autres timelines |
| Derived Artifact | concept Investigate | Investigate | ≠ Evidence |
| Reproducibility Assessment | assessment Investigate | Investigate + Shared Trace | schéma futur |

Aucun schéma, JSON Schema, cardinalité, packet format, protocole ou machine d’état objet finale n’est défini.
