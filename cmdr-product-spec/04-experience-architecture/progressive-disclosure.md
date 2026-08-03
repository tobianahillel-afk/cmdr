---
id: experience-progressive-disclosure
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-007
  - REQ-UX-004
  - REQ-UX-005
---

# Divulgation progressive

## Niveau 1 — Situation et action

Titre, statut, owner, priorité, impact, fraîcheur, information critique et prochaine action. Toujours visible lorsque pertinent, indépendamment du rôle.

## Niveau 2 — Contexte et relations

Relations, timeline, Evidence résumées, activité, métriques, dépendances et détails. Disponible dans le workspace ou l'Inspector sans navigation destructive.

## Niveau 3 — Trace et technique

Raw events, payloads autorisés, logs, provenance complète, Tool Calls, audit et détails techniques. Disponible à la demande, avec permissions et alternative accessible.

## Règles

- une donnée nécessaire à une décision ou à la sécurité n'est jamais cachée ;
- le rôle influence l'ordre et l'expansion initiale, pas l'existence de l'information ;
- un accordéon ne masque pas une erreur, une permission ou une conséquence ;
- les sections se souviennent localement de leur état sans changer la source canonique ;
- l'IA ne devient pas un quatrième niveau opaque.

## Critère d’acceptation

**Given** Business Owner et Reverse Engineer sur la même Decision,  
**When** ils ouvrent le workspace,  
**Then** le premier voit impact et alternatives en Niveau 1, le second accède directement aux Evidence techniques, mais statut, autorité et source restent identiques.
