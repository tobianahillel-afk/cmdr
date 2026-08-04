---
id: dynamic-sandbox-objects
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-OBJ-002
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-OBJ-009
open_decisions:
  - OPEN-014
  - OPEN-015
---
# Objects and concepts

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | contexte | read/link | cardinalités | Objets |
| Artifact / Derived Artifact | Investigate | source et inputs dérivés | read/link | modèle dérivé | Objets |
| Runtime Artifact | concept Investigate | sortie runtime | capture/read/export/link | objet absent | Objets |
| Evidence / Finding | Investigate | qualification aval | candidate/draft only | trust/review | Objets/Trust |
| Tool / Tool Call | Studio | exécution et attribution | select/invoke/read | objets finaux absents | Studio/Objets |
| Workflow / Automation Run | Studio | orchestration optionnelle | invoke/read/link | OPEN-015 | Studio/Objets |
| Dynamic Analysis Session | concept Investigate | contexte multi-Run | create/update/close | objet absent | Objets |
| Sandbox Environment | Platform Settings | environnement autorisé | read/select/request alternative | admin reste Settings | Settings/Objets |
| Sandbox Run | concept Investigate | exécution isolée | prepare/start/stop/retry/read | objet absent | Objets |
| Behavioral Observation | concept Investigate | événement observed/inferred/annotation | read/annotate/link | objet/record à décider | Objets |
| Process/Network/System Observation | concepts Investigate | projections analytiques | read/filter/compare | schémas absents | Objets |
| Interaction Profile | Studio/Settings definition, Investigate use | scénario fonctionnel | select/version/read | ownership final à préciser | Studio/Settings |
| Provenance Record | Shared mechanism, Investigate semantics | reconstruction | emit/review | objet absent | Shared/Trust |
| Analysis Result | concept Investigate | résultats consolidés | read/review/link | objet absent | Objets |
