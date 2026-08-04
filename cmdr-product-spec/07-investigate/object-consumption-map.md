---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-014
  - REQ-PROD-061
  - REQ-PROD-062
open_decisions:
  - OPEN-014
  - OPEN-015
---

# Object consumption map — Investigate Phase 4B.1

| Objet | Owner actuel | Usage Investigate | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | workspace durable | create, update, link, close/reopen conceptuels | machine finale/cardinalités | Objets |
| Hypothesis | Investigate | question testable | create, review, link, supersede | états/confiance finaux | Objets |
| Artifact | Investigate | élément analysable | register, version, derive, link | stockage/hash/dedup | Objets/Technique |
| Evidence | Investigate | élément qualifié | create, qualify, version, dispute, cite | intégrité/custody/review | Objets/Trust |
| Finding | Investigate | conclusion revue | draft, review, confirm, dispute, supersede | niveaux validation/confiance | Objets/Permissions |
| Telemetry Event | Shared | recherche/inspection | read, filter, link | moteur/index | Technique |
| Detection | Command | contexte source | read/link | engineering distinct en 4B.3 | 4B.3 |
| Signal | Command | triage/pivot | read, qualifier contribution, link | dispositions finales | Command/Objets |
| Alert | Command | contexte | read/link | chaîne Signal/Alert | Command/Objets |
| Incident | Command | contexte parent | read/link | cardinalités Case | Objets |
| Task | Command | travail persistant | request/create/link sous workflow Command | permissions | Command/Phase 7 |
| Query | Shared | authoring/exécution | draft/reference | dialecte | Technique |
| Search Job | Shared | run/results | execute/cancel/read | moteur/retenue résultats | Technique |
| Entity | Shared | relations et pivots | annotate/link, propose merge | identity resolution | Shared/Technique |
| Action Request | Govern lifecycle | prepare/submit/follow | producer operations | workflow/policy/authority | 4C/7 |
| Decision | Govern | contexte aval | read/link | 4C | Govern |
| Response Run | Govern | contexte aval | read/link/verify | bridge Run | 4C/7 |
| Result | Govern | retour au Case | read/link | verification lifecycle | 4C/7 |
| Endpoint | shared, administré par Settings | contexte/source future | read/select | objet absent et protocoles collecte | 4B.2/4D |
| Endpoint Agent Fleet | Platform Settings | projection | read only | fleet/admin | 4D |
| Workflow | CMDR Studio | automation | invoke authorized version | deployment/contracts | Studio |
| Automation Run | CMDR Studio | provenance proposals | read/link | `OPEN-015` | 4C/7 |
| Hunt | ouvert | workspace Search/Hypothesis | organize/link | objet ou composition ? | Objets |
| Attachment | ouvert/absent | fichier joint documentaire | attach, qualify proposal | `OPEN-014` | Objets |
| Report | Reporting Engine Shared | draft/content | contribute/cite | pas d’objet report actuel | Shared |
| Note / Comment | Shared concepts, fichiers absents | collaboration Case | contribute/link | objets/rétention | Shared/Objets |
| Timeline Entry | Shared | chronologie | emit content/link | modèle final | Shared/Objets |
| Audit / Provenance record | mécanismes Shared, fichiers absents | trace | emit business events | objets/retention | Shared/Trust |

Aucun schéma, JSON Schema, cardinalité ou machine d’état finale n’est défini ici.
