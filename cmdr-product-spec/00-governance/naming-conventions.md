---
id: naming-conventions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Conventions de nommage

## Chemins

- Dossiers et fichiers Markdown en `kebab-case`.
- Produits canoniques : `command`, `investigate`, `govern`, `cmdr-studio`, `platform-settings`, `endpoint-agent`.
- ADR : `ADR-0001-description.md`.
- Écrans : nom fonctionnel en kebab-case et identifiant dans le front matter.
- Interdits : `final.md`, `new.md`, `version2.md`, `corrected.md`, suffixes de copie et noms en majuscules hors ADR.

## Identifiants

- Écran : `<PRODUCT>-<MODULE>-<NNN>`.
- Objet : `OBJ-<NAME>`.
- Permission : `perm.<domain>.<resource>.<action>`.
- Capacité : `CAP-<DOMAIN>-<NNN>`.
- Composant : `CMP-<NAME>`.
