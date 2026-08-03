---
id: ADR-0004-screen-specification-contract
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# ADR-0004 — Contrat des écrans

## Contexte

Le socle documentaire doit éviter les architectures concurrentes et permettre la traçabilité produit.

## Décision

Chaque écran possède un front matter, un identifiant canonique, un propriétaire, des permissions référencées, 27 sections et les six états Loading, Empty, Partial, Error, Offline et Permission denied.

## Conséquences

- Les registres et manifestes appliquent cette décision.
- Toute exception exige un nouvel ADR.
- Les anciens chemins non conformes sont supprimés ou explicitement archivés comme non normatifs.

## Statut

Draft — à approuver par les propriétaires concernés.
