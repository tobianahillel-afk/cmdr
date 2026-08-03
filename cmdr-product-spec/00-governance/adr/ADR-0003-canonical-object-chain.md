---
id: ADR-0003-canonical-object-chain
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# ADR-0003 — Chaîne canonique d’objets

## Contexte

Le socle documentaire doit éviter les architectures concurrentes et permettre la traçabilité produit.

## Décision

La chaîne est exactement: Telemetry Event → Detection → Signal → Alert → Incident → Case → Evidence → Finding → Action Request → Decision → Response Run → Result.

## Conséquences

- Les registres et manifestes appliquent cette décision.
- Toute exception exige un nouvel ADR.
- Les anciens chemins non conformes sont supprimés ou explicitement archivés comme non normatifs.

## Statut

Draft — à approuver par les propriétaires concernés.
