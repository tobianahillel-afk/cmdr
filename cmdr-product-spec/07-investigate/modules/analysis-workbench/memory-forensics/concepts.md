---
id: investigate-memory-forensics-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Concepts et distinctions

| Concept | Nature dans cette phase | Distinction obligatoire |
|---|---|---|
| Memory Image | objet canonique Investigate | ≠ Request, Job, Session, Derived Artifact ou Evidence |
| Memory Forensics Session | concept fonctionnel Investigate | ≠ Analysis/Reverse/Debugger Session, Automation Run ou Collection Job |
| Platform Candidate / Analysis Profile | concepts fonctionnels | proposed/selected ≠ certain |
| Process / Thread Observation | reconstructions | ≠ état runtime courant |
| Memory Region / Mapping | observations | exécutable/unmapped ≠ injection |
| Module / Driver Observation | observations | chargé/présent ≠ malveillant/rootkit |
| Handle / System Object / IPC | observations techniques | ≠ objet canonique CMDR ou coordination malveillante |
| Network State Observation | trace présente en mémoire | ≠ session réseau complète ou IOC |
| Sensitive Material Candidate | candidate masquée | ≠ credential valide ou autorisation d’usage |
| Memory Anomaly | interprétation candidate | ≠ Finding |
| Kernel State Observation | projection dépendante du support | incohérence ≠ rootkit confirmé |
| Memory Timeline | projection locale | ≠ Case Timeline |
| Extraction Result | résultat borné | carving ≠ récupération certaine |
| Reproducibility Assessment | évaluation fonctionnelle | ≠ infrastructure technique |

Aucun schéma, format mémoire, cardinalité, offset, signature ou état objet définitif n’est défini.
