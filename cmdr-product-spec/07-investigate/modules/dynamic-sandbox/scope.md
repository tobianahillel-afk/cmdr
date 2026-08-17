---
id: dynamic-sandbox-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-INV-005
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
---
# Scope

## Inclus
Préconditions, environnement autorisé, sessions/runs, comportements observés, process tree, changements fichiers/système, réseau, persistance candidate, interactions simulées, Runtime Artifacts, comparison, safety, provenance et handoff.

## Exclu
Instrumentation, virtualisation, produits tiers, techniques d’évasion, reverse/debugger, forensics avancé, cibles réelles et action de réponse.

## Règle
Le produit définit ce que l’analyste demande, voit, interrompt, compare et transmet ; il ne définit pas comment la sandbox est implémentée.
