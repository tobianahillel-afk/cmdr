---
id: investigate-palette-proposals
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: proposal
decision-status: proposed
requirements:
  - REQ-BRAND-006
  - REQ-PROD-048
  - REQ-BRAND-008
---
# Propositions de palette Investigate

## Statut

`OPEN-001` reste `open`. Les trois directions ci-dessous sont `proposed`. Elles servent à une revue comparative ; aucune n'est canonique.

Toutes héritent de CMDR Ink, Bone, Moss, Ember, Fog, Graphite et Stone. Les noms et valeurs ci-dessous sont des candidats de produit, pas des tokens de production.

## Direction A — Verdigris Ledger

### Intention

Évoquer analyse matérielle, patine, profondeur et traçabilité sans utiliser un bleu de cybersécurité générique.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#EDEFEA` |
| Surface claire | `#F7F7F2` |
| Canvas sombre | `#202927` |
| Surface sombre | `#293431` |
| Texte profond | `#24312F` |
| Accent principal | `#2F6A65` |
| Accent secondaire | `#7A6A52` |
| Trace / sélection douce | `#B8C8C1` |

- Accent sur canvas clair : environ 5.38:1.
- Texte profond sur canvas clair : environ 11.66:1.
- Différenciation non chromatique : crochets de provenance, annotations marginales, règles de comparaison.
- Risque : proximité avec Command si Juniper et l'accent Investigate sont employés de la même manière.
- Rejet si : le produit semble seulement « Command plus sombre ».

## Direction B — Lichen Archive

### Intention

Créer une atmosphère d'archive vivante, de matière et de patience analytique, plus organique que technologique.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#EFF0E8` |
| Surface claire | `#F8F7F1` |
| Canvas sombre | `#222720` |
| Surface sombre | `#2B3229` |
| Texte profond | `#252B25` |
| Accent principal | `#536B55` |
| Accent secondaire | `#886A4A` |
| Trace / sélection douce | `#C5CFBD` |

- Accent sur canvas clair : environ 5.08:1.
- Texte profond sur canvas clair : environ 12.61:1.
- Différenciation non chromatique : couches, marqueurs d'artefact et annotations archivistiques.
- Risque : trop proche de Moss ou d'une esthétique environnementale.
- Rejet si : la relation à l'investigation technique devient trop douce ou naturaliste.

## Direction C — Mineral Graph

### Intention

Mettre l'accent sur structure, relation, graphe et précision minérale avec un bleu-vert grisé non spectaculaire.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#ECEFEE` |
| Surface claire | `#F7F8F5` |
| Canvas sombre | `#20272A` |
| Surface sombre | `#2A3336` |
| Texte profond | `#232A2D` |
| Accent principal | `#425F65` |
| Accent secondaire | `#826A5B` |
| Trace / sélection douce | `#BAC9C9` |

- Accent sur canvas clair : environ 5.94:1.
- Texte profond sur canvas clair : environ 12.60:1.
- Différenciation non chromatique : connecteurs, grilles de comparaison et marqueurs temporels.
- Risque : glissement vers un bleu sécurité conventionnel.
- Rejet si : l'accent devient saturé ou associé à un glow.

## Comparaison

| Critère | Verdigris Ledger | Lichen Archive | Mineral Graph |
|---|---|---|---|
| Profondeur technique | forte | moyenne | forte |
| Différence avec Command | moyenne | forte | forte |
| Compatibilité dark workbench | forte | moyenne | forte |
| Risque de cliché | proximité teal | naturaliste | bleu cyber |
| Chaleur éditoriale | moyenne | forte | moyenne |
| Accent non couleur | provenance | strates | relation |

## Preuves nécessaires

- maquettes clair/sombre sur Case, code, timeline, graph et Evidence ;
- tests de contraste et déficiences colorimétriques ;
- comparaison directe avec Command ;
- test longue session ;
- validation sponsor et Brand Design Lead.

## Interdiction de consommation

Tant qu'OPEN-001 est ouverte, `03-design-system/tokens/product-theme-tokens.md` ne doit sélectionner aucune de ces directions.
