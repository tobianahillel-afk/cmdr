---
id: investigate-network-forensics-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Workflows

| Étape | Source | Résultat | Limite |
|---|---|---|---|
| Intake | Case, Collection Job ou Artifact | assessment et objectif | aucune analyse automatique |
| Trust review | capture, acquisition, custody | integrity/scope assessment | format et acquisition hors périmètre |
| Coverage | sensors, interfaces, timebases | coverage assessment | absence observée ≠ absence certaine |
| Inspection | capture | packets/frames sourcés | aucune mutation ou replay |
| Reconstruction | packets | flows, sessions, conversations | partialité visible |
| Specialized analysis | conversations | DNS, transactions, encrypted metadata, transfers | aucune qualification automatique |
| Reasoning | observations | Entities, anomalies, timeline et comparisons | candidates et contradictions visibles |
| Extraction | source bornée | Derived Artifact | aucun contenu exécuté |
| Provenance | tous les records | reproducibility assessment | Shared Trace non redéfini |
| Handoff | sélection humaine | Evidence/Finding/future package | destination owner qualifie |

Le return origin, le Case, les permissions, les restrictions et les erreurs sont préservés à chaque transition.
