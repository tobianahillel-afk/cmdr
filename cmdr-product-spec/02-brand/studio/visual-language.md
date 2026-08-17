---
id: studio-visual-language
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-008
  - REQ-PROD-016
  - REQ-AI-002
  - REQ-AI-007
  - REQ-AI-010
---
# Langage visuel CMDR Studio

## Nodes

- type, nom, version et owner visibles ;
- ports explicites et connecteurs orthogonaux ;
- état par texte, forme et sémantique ;
- groupement par responsabilité ou phase ;
- aucun glow, particule ou flux continu ;
- node compact, non carte marketing.

## Workflow

Le canvas montre la logique. La liste ou outline fournit une alternative. Les branches, erreurs et Human Gates restent navigables au clavier et lisibles sans couleur.

## Tool Calls

Un Tool Call apparaît comme une entrée de trace : tool, version, initiateur, inputs référencés, permission, durée, résultat, erreur et provenance. Il ne se réduit pas à une animation de node.

## Human Gates

- interruption visible ;
- question, owner, délai et conséquence ;
- relation à Govern indiquée lorsque nécessaire ;
- état awaiting, approved, rejected, expired ou bypass impossible selon contrat futur ;
- aucune célébration ou animation ludique.

## Automation Runs

- timeline ou ledger ;
- état courant, étapes terminées, attente et erreur ;
- stop, pause ou reprise visibles selon permission ;
- sources et outputs ;
- lien au workflow/version ;
- distinction avec Response Run.

## Evaluations et simulations

- datasets, versions, mesures, seuils et régressions ;
- comparaison, intervalles et limites ;
- aucune note unique décorative ;
- résultat proposé ou validé clairement distinct.

## Coûts et usage

Les coûts sont contextualisés par provider, run, tenant, période et résultat. Ils ne deviennent pas une gamification.

## Erreurs

Erreur locale, cause, étape, retry et impact visibles. Les erreurs ne recolorent pas tout le canvas.

## Critère d'acceptation

Un Automation Run peut être compris et interrompu à partir de sa trace même si toutes les animations et couleurs de produit sont désactivées.
