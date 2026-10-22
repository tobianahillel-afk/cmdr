---
id: brand-architecture
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-005
  - REQ-BRAND-006
  - REQ-BRAND-008
  - REQ-PROD-013
  - REQ-PROD-016
---
# Architecture de marque

## Modèle

CMDR est la marque mère. Command, Investigate, Govern et CMDR Studio sont des identités produit reliées. Platform Settings et Endpoint Agent utilisent le système CMDR et les accents du contexte propriétaire ; ils ne reçoivent pas une identité indépendante pendant cette phase.

Les libellés `Command Center`, `Investigation Lab` et `Response & Governance` sont autorisés dans l'interface, mais les noms canoniques restent Command, Investigate et Govern.

## Matrice de relation

| Dimension | Commun à tous | Adaptable par produit | Interdit |
|---|---|---|---|
| Marque | wordmark CMDR, essence et signature Moss + Ember | lockup avec nom produit | logo entièrement différent |
| Typographie | familles, rôles et métriques communes | poids ou accent de titre dans les limites du système | famille différente par produit |
| Grille | alignements, rails, Context Bar, Inspector | proportion du canvas et densité | navigation incompatible |
| Surfaces | vocabulaire canvas/surface/inset/overlay | tonalité et fréquence des surfaces secondaires | nouveau Design System |
| Couleur | CMDR Ink/Bone et sémantique commune | accent produit et data-viz locale | couleur seule comme identité |
| Iconographie | style, taille, stroke et libellés | sous-ensemble fonctionnel | pack d'icônes différent |
| Mouvement | causal, bref, reduced motion | rythme local d'exécution ou d'exploration | animation décorative de produit |
| Provenance | source, owner, auteur/run, date, statut | ordre de mise en avant | trace masquée |
| Navigation | shell, Context Bar, Inspector, retour | signature du workspace | structure contradictoire |
| Visualisations | unités, source, légende, alternative | motifs adaptés à l'activité | conventions sémantiques différentes |

## Caractères produits

| Produit | Question dominante | Rythme | Densité | Motif fonctionnel |
|---|---|---|---|---|
| Command | Que faut-il traiter et qui le possède ? | stable, temporel, priorisé | moyenne à dense | bandes de situation, temps, ownership |
| Investigate | Qu'est-ce qui est établi et comment ? | exploratoire, comparatif | dense | provenance, annotation, relation, timeline |
| Govern | Que peut-on décider et sous quelles conditions ? | séquentiel, délibératif | moyenne | sections de décision, conditions, signatures, rollback |
| Studio | Comment construire, vérifier et superviser l'automatisation ? | constructif puis exécutable | moyenne à dense | grille de construction, versions, runs, gates |

## Platform Settings

Settings utilise le système CMDR avec une expression faible : structure administrative, haute lisibilité, faible accent décoratif. Il ne reçoit pas une palette produit autonome dans cette phase.

## Endpoint Agent

Endpoint Agent n'est pas une console de marque séparée. Les surfaces de flotte suivent Settings ; les surfaces d'investigation suivent Investigate ; les actions suivent Govern. Cette règle reflète les frontières produit.

## Règle de changement

Une variation produit doit démontrer une utilité liée à l'activité. Si elle ne peut être justifiée que par « rendre le produit différent », elle est rejetée.
