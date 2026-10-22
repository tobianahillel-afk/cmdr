---
id: ai-visual-expression
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-007
  - REQ-AI-010
  - REQ-BRAND-002
---
# Expression visuelle de l'IA et de l'automatisation

## Position

CMDR est augmentable par l'IA, non dépendant de l'IA. Le langage visuel doit montrer une assistance attribuée et gouvernée, pas une présence magique ou omnisciente.

## Représentation obligatoire

Une production IA ou automatisée expose selon le contexte :

- type : règle, moteur déterministe, workflow, Automation Agent ou modèle ;
- initiateur ;
- run ou version ;
- sources et Tool Calls ;
- confiance ou incertitude ;
- statut : suggestion, draft, reviewed, accepted ou rejected ;
- owner humain ou prochaine validation ;
- accès à la trace.

## Distinction

| Production | Traitement |
|---|---|
| Humaine | auteur et rôle |
| Déterministe | règle/moteur, version et inputs |
| Automation Run | workflow, état, étapes et interruption |
| IA proposée | modèle/run, provenance, incertitude et statut de proposition |
| Decision | autorité Govern, jamais badge IA |
| Result | exécution et vérification |

## Iconographie

Il n'existe pas de couleur ou d'icône universelle « IA ». Une icône peut signaler une origine automatisée si elle est accompagnée d'un label. Les étoiles, baguettes, robots et cerveaux sont interdits comme représentation par défaut.

## Placement

L'assistance apparaît près de l'activité qu'elle soutient : recherche, résumé, hypothèse, explication, rédaction ou proposition. Elle ne crée pas un bouton `Ask AI` répété dans chaque surface ni un chat principal obligatoire.

## États indisponibles

Lorsque le modèle ou provider est indisponible :

- le workflow principal reste actif ;
- la fonction optionnelle explique son indisponibilité ;
- aucune zone essentielle ne devient vide ;
- les chemins manuels et déterministes restent visibles.

## Assurance

Les états d'évaluation, simulation, interruption et erreur utilisent le langage Studio et les tokens sémantiques communs. Ils ne sont pas décorés par un gradient ou une aura.

## Critère d'acceptation

**Given** deux conclusions identiques, l'une humaine et l'autre produite par un modèle,  
**When** elles sont affichées,  
**Then** leur structure reste comparable, mais l'origine, la trace, l'incertitude et l'autorité sont explicitement différentes.
