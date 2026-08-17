---
id: naming-conventions
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
---
# Conventions de nommage

## Chemins

- Dossiers et fichiers Markdown : `kebab-case`.
- Préfixe ADR : `ADR-NNNN-description.md`.
- Produits : `command`, `investigate`, `govern`, `cmdr-studio`, `platform-settings`, `endpoint-agent`.
- Noms interdits : `final`, `new`, `version2`, `corrected`, `copy`, suffixes de dates servant de version.
- Un renommage préserve l'historique Git et met à jour liens, registres et matrice.

## Identifiants

| Type | Convention | Exemple |
|---|---|---|
| Requirement | `REQ-<DOMAIN>-NNN` | `REQ-PROD-001` |
| Question ouverte | `OPEN-NNN` | `OPEN-006` |
| ADR | `ADR-NNNN` | `ADR-0005` |
| Dépendance | `DEP-NNN` | `DEP-004` |
| Risque produit | `RISK-PROD-NNN` | `RISK-PROD-012` |
| Écran | `<PRODUCT>-<MODULE>-NNN` | `CMD-IWQ-001` |
| Objet | `OBJ-<NAME>` | `OBJ-INCIDENT` |
| Permission | `perm.<domain>.<resource>.<action>` | `perm.command.incident.read` |
| Capability | `CAP-<DOMAIN>-NNN` | `CAP-INV-004` |
| Composant | `CMP-<NAME>` | `CMP-INSPECTOR` |

## Statuts

Valeurs autorisées dans le front matter :

- `draft`
- `in-review`
- `validated`
- `implemented`
- `deprecated`

`review` et `approved` sont des valeurs historiques à migrer, pas de nouveaux statuts.

## Titres et noms d'interface

Le nom de fichier reste stable et technique. Le titre peut employer un nom d'interface approuvé. Un nom marketing n'introduit ni nouveau domaine ni nouveau propriétaire.

## Noms de capabilities et classifications

Le nom d'une capability décrit le résultat utilisateur, pas le fournisseur. La marque ou le moteur intégré apparaît dans les métadonnées d'intégration, jamais comme structure principale de navigation.

## Critères d'acceptation

- deux identifiants de même type ne désignent jamais deux concepts ;
- un renommage ne crée pas de seconde source active ;
- une question ouverte conserve son identifiant jusqu'à résolution ;
- un terme interdit ou un statut historique est détecté avant revue.
