---
id: studio-palette-proposals
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: proposal
decision-status: proposed
requirements:
  - REQ-BRAND-006
  - REQ-PROD-050
  - REQ-BRAND-008
---
# Propositions de palette CMDR Studio

## Statut

`OPEN-003` reste `open`. Les trois directions sont `proposed`.

## Direction A — Workshop Sienna

### Intention

Exprimer construction, matière, itération et versioning par une chaleur d'atelier sourde, sans orange de warning.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#F0EBE6` |
| Surface claire | `#FAF7F3` |
| Canvas sombre | `#29221E` |
| Surface sombre | `#342A25` |
| Texte profond | `#29221E` |
| Accent principal | `#745344` |
| Accent secondaire | `#4E6A65` |
| Node / sélection douce | `#C8B4A7` |

- Accent sur canvas clair : environ 5.79:1.
- Texte profond sur canvas clair : environ 13.21:1.
- Différenciation non chromatique : grille, version brackets, ports et traces.
- Risque : confusion avec Govern Bronze ou une esthétique artisanale.
- Rejet si : la chaleur devient orange, décorative ou non technique.

## Direction B — Petrol Assembly

### Intention

Mettre en avant assemblage, exécution et contrôle par un petrol sobre, distinct d'un bleu IA.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#EAF0EF` |
| Surface claire | `#F6F9F7` |
| Canvas sombre | `#202829` |
| Surface sombre | `#293335` |
| Texte profond | `#202829` |
| Accent principal | `#35605F` |
| Accent secondaire | `#80664E` |
| Node / sélection douce | `#B6C9C5` |

- Accent sur canvas clair : environ 6.09:1.
- Texte profond sur canvas clair : environ 13.03:1.
- Différenciation non chromatique : connecteurs orthogonaux, étapes et run ledger.
- Risque : proximité avec Investigate Verdigris ou Command Juniper.
- Rejet si : tous les produits convergent vers le teal.

## Direction C — Iris Graphite

### Intention

Utiliser une iris grisée pour distinguer composition et simulation sans adopter le violet IA attendu.

| Rôle proposé | Valeur |
|---|---|
| Canvas clair | `#EDECEF` |
| Surface claire | `#F8F7F8` |
| Canvas sombre | `#25232A` |
| Surface sombre | `#302D36` |
| Texte profond | `#25232A` |
| Accent principal | `#5B566D` |
| Accent secondaire | `#7A6C53` |
| Node / sélection douce | `#C2BEC9` |

- Accent sur canvas clair : environ 5.95:1.
- Texte profond sur canvas clair : environ 13.19:1.
- Différenciation non chromatique : grille, versions, evaluations et Human Gates.
- Risque : lecture immédiate « produit IA violet ».
- Rejet si : saturation, gradient, étoiles ou glow sont nécessaires pour la rendre distinctive.

## Comparaison

| Critère | Workshop Sienna | Petrol Assembly | Iris Graphite |
|---|---|---|---|
| Construction | forte | forte | moyenne |
| Exécution | moyenne | forte | moyenne |
| Assurance | moyenne | forte | forte |
| Différence avec Investigate | forte | faible | forte |
| Risque cliché IA | faible | faible | élevé |
| Chaleur | forte | moyenne | faible |

## Preuves nécessaires

- Library, Builder, Tool Call, Human Gate, Automation Run, evaluation et simulation ;
- clair/sombre ;
- comparaison avec Investigate et Govern ;
- tests de sémantique erreur/succès/warning ;
- accessibilité ;
- validation sponsor.

## Interdiction de consommation

Aucune direction n'entre dans les tokens avant résolution d'OPEN-003.
