---
id: ADR-0002-single-source-of-truth
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# ADR-0002 — Source de vérité unique

## Contexte

Le socle documentaire doit éviter les architectures concurrentes et permettre la traçabilité produit.

## Décision

Chaque objet, permission, composant, palette et contrat a une source unique enregistrée. Les autres documents font des liens et décrivent seulement leur usage local.

## Conséquences

- Les registres et manifestes appliquent cette décision.
- Toute exception exige un nouvel ADR.
- Les anciens chemins non conformes sont supprimés ou explicitement archivés comme non normatifs.

## Statut

Draft — à approuver par les propriétaires concernés.
