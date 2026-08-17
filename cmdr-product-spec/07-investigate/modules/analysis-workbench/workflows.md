---
id: analysis-workbench-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-PROD-014
  - REQ-PROD-020
---
# Workflows and transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case | ouvrir Artifact | CAP-INV-301 | tenant, environment, Case, Artifact, Hypothesis, objectif, return origin | Case |
| CAP-INV-301 | créer/reprendre | CAP-INV-302 | Artifact, type, provenance, restrictions, famille, owner | intake |
| CAP-INV-302 | sélectionner analyse | CAP-INV-303 | session, Artifact, objectif, permissions | session |
| CAP-INV-303 | lancer Tool | Tool Call Studio | Tool/version, paramètres, Artifact, limites | Workbench |
| Analysis Result | transformation | CAP-INV-311 | source, transformation, Tool/version, paramètres | résultat |
| Derived Artifact | analyser | CAP-INV-301/302 | parent, provenance, restrictions | source session |
| Analysis Result | qualifier | CAP-INV-313 | observations, sources, contradictions, provenance | Workbench |
| CAP-INV-313 | Evidence candidate | CAP-INV-107/108 | candidate, Case, sources, qualification requise | Workbench |
| CAP-INV-313 | Finding Draft | CAP-INV-109 | draft, Evidence existantes, incertitude, auteur | Workbench |
