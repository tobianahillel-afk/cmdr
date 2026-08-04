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
  - OPEN-008
  - OPEN-014
  - OPEN-015
---
# Object consumption map — Investigate through Phase 4B.2A

| Objet | Owner actuel | Usage Investigate | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | workspace et contexte de collecte | create/update/link, return origin | machine finale/cardinalités | Objets |
| Hypothesis | Investigate | question testable | create/review/link/supersede | états/confiance | Objets |
| Artifact | Investigate | résultat matériel analysable | register/version/link/derive | stockage/hash/dedup | Objets/Technique |
| Evidence | Investigate | qualification explicite | create/qualify/dispute/cite | custody/trust final | Objets/Trust |
| Finding | Investigate | conclusion et base de demande | draft/review/supersede | validation/confiance | Objets/Permissions |
| Endpoint | concept partagé à formaliser | cible d’investigation | read/select/link | objet canonique absent | Objets |
| Endpoint Agent | Endpoint Agent | état, capacités, exécution et résultat locaux | read/invoke authorized projection | contrat détaillé | Endpoint/Technique |
| Endpoint Agent Fleet | Platform Settings | posture et affectation | read only | aucune administration locale | Settings |
| Endpoint Policy | Platform Settings | restrictions et gate | read only | évaluation détaillée future | Settings/Govern |
| Collection Request | Investigate | demande préparée et soumise | create/update/supersede/cancel/submit | états actuels mêlent demande/exécution | Objets |
| Collection Job | concept métier Investigate | suivi de l’exécution | observe/cancel/retry/link | objet absent | Objets |
| Background Job | Shared Capabilities | mécanisme de queue/progression | consume/read/cancel via contrat | contrat technique futur | Shared/Technique |
| Live Session | concept métier Investigate | session visible Case-scoped | request/open/join/extend/close | objet absent | Objets |
| Endpoint Operation | concept métier Investigate | opération autorisée de session | prepare/execute/interrupt | relation à Agent Command | Objets |
| Agent Command | Endpoint Agent | commande locale propriétaire | read status/correlation | contrat technique futur | Endpoint/Technique |
| Operation Result | concept métier Investigate | résultat local inspectable | receive/review/verify/dispute/link | objet absent | Objets |
| Custody Record | concept Investigate/Trust | chaîne de custody | emit/review/dispute | objet absent | Objets/Trust |
| Provenance Record | mécanismes Shared, sémantique Investigate | reconstruction | emit/link/review | objet absent | Shared/Trust |
| Incident / Task | Command | contexte et coordination | read/link/request Task | cardinalités/permissions | Command/Objets |
| Action Request / Decision / Response Run / Result | Govern | autorité et retour | prepare/read/link | contracts futurs | Govern/Objets |
| Workflow / Automation Run / Tool Call | Studio | orchestration et provenance | invoke/read/link | OPEN-015 | Studio/Objets |
| Memory Image | Investigate Draft | Artifact d’acquisition mémoire | create/link | format/moteur/support | Objets/4B.2B |
| Attachment | ouvert | fichier documentaire | attach/reference | OPEN-014 | Objets |
| Query / Search Job / Entity / Timeline Entry | Shared | recherche, relation et chronologie | consume/link/emit semantics | contrats techniques | Shared/Technique |

Aucun schéma, JSON Schema, cardinalité, machine d’état finale ou permission atomique n’est défini ici.