---
id: dynamic-sandbox-workflows
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

| Source | Déclencheur | Destination | Contexte transmis | Return/error |
|---|---|---|---|---|
| Case | ouvrir analyse dynamique | CAP-INV-314 | tenant, environment, Case, Artifact, objective, return origin | Case; draft preserved |
| Static Analysis Session | demander comportement | CAP-INV-314 | Artifact, Derived Artifacts, static results, provenance | static Workbench |
| CAP-INV-314 | préconditions satisfaites | CAP-INV-316 | objective, scope, risks, permissions, limits | intake |
| CAP-INV-316 | choisir environnement | CAP-INV-315 | requirements, restrictions, profile | session |
| CAP-INV-316 | préparer/lancer | CAP-INV-317 | Artifact, environment, profile, limits, initiator | session |
| Sandbox Run | événement | CAP-INV-318..322 | timestamps, source, process/file/network/system refs | Run |
| Sandbox Run | sortie matérielle | CAP-INV-324 | source event/process, Tool/version, restrictions | Run |
| Sandbox Run | comparer | CAP-INV-325 | Runs, environments, profiles, versions, observations | session |
| Failed Run | recover/retry | CAP-INV-317 | errors, partials, cleanup/reset status | session |
| Unsafe Environment | signaler | Platform Settings | environment/version, symptoms, trace, block status | Run/session |
| Runtime Artifact | inspecter statiquement | CAP-INV-301 | parent Run, provenance, restrictions, return origin | dynamic Workbench |
| Runtime Artifact | future reverse | future 4B.2B.2B | Artifact, results, provenance | dynamic Workbench |
| Dynamic Result | qualifier | CAP-INV-328 → 107/108/109 | observations, Runtime Artifacts, contradictions, provenance | dynamic Workbench |
