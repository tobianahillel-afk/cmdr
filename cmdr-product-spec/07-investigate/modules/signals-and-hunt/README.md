---
id: investigate-signals-and-hunt
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-005
  - REQ-PROD-014
  - REQ-PROD-025
  - REQ-PROD-045
---

# Signals and Hunt

## Mission

Recevoir ou consulter des Signals, rechercher dans la télémétrie, inspecter des Events, pivoter, organiser des Hunts, sauvegarder les Query Assets et transmettre des éléments sourcés vers un Case.

## Capabilities

`CAP-INV-001` à `CAP-INV-008`.

## Shared Capabilities

| Shared Capability | Usage Investigate | Données locales | Source canonique |
|---|---|---|---|
| Query/Search Job | validation, exécution, résultats | contexte, liens, annotations | Shared Query Engine / object register |
| Saved Views | vues et filtres | catalogue système éventuel | Shared Saved Views |
| Background Jobs | runs/exports | événements métier | Shared |
| Inspector | Event/Signal inspection | sections locales | Design System |
| Trace/Activity Stream | provenance | événements Search/Hunt | Shared |
| Object Linking | Case/Hunt/Signal relations | relation reason | Shared |
| Export | sélection permissionnée | manifest/redaction context | Shared |

## Limites

Aucun dialecte, moteur, index, stockage, Detection Rule ou capability 4B.3 n’est défini.
