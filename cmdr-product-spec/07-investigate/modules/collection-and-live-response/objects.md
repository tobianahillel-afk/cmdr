# Object and concept map

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Endpoint | concept à formaliser | cible | read/select/link | objet absent | Objets |
| Endpoint Agent | Endpoint Agent | état/capacités/exécution | read/invoke authorized | contrat détaillé | Endpoint/Technique |
| Endpoint Agent Fleet | Platform Settings | posture projection | read only | aucune | Settings |
| Endpoint Policy | Platform Settings | restrictions projection | read only | policy evaluation detail | Settings/Govern |
| Case | Investigate | contexte | read/link | cardinalités | Objets |
| Collection Request | Investigate | demande | create/update/submit/cancel | états mêlent exécution | Objets |
| Collection Job | concept métier Investigate | suivi exécution | observe/cancel/retry | objet absent | Objets |
| Background Job | Shared | mécanisme générique | projection | contrat futur | Shared/Technique |
| Live Session | concept métier Investigate | session visible | request/open/join/close | objet absent | Objets |
| Endpoint Operation | concept métier | opération session | prepare/execute/interrupt | relation Agent Command | Objets |
| Operation Result | concept métier | résultat local | receive/review/link | objet absent | Objets |
| Artifact / Evidence / Finding | Investigate | résultats/analyse | create/link/qualify | schemas Draft | Objets |
| Action Request / Decision / Response Run / Result | Govern | autorité et retour | prepare/read | contrats futurs | Govern/Objets |
| Automation Run / Tool Call | Studio | provenance automation | read/link | OPEN-015 | Studio/Objets |
| Custody / Provenance Record | concepts | confiance | emit/review | objets absents | Objets/Trust |

Aucun schéma ou état final n’est défini.
