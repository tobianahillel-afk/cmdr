---
id: market-positioning
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-012
---
# Positionnement produit

## Formulation

CMDR est une plateforme opérationnelle et probatoire de cybersécurité reliant coordination, investigation, décision gouvernée, réponse et amélioration sur des objets communs.

## Différences de catégorie

| Catégorie | Ce qu'elle apporte généralement | Ce que CMDR ajoute comme destination |
|---|---|---|
| SIEM | collecte, recherche, détection | continuité vers Incident, Case, Evidence, Decision et Result |
| EDR | visibilité et action endpoint | intégration à investigation, gouvernance et preuve ; Endpoint natif comme cible |
| SOAR | orchestration et playbooks | responsabilité, Evidence, Decision, Human Gates et résultat vérifié |
| Ticketing | suivi du travail | objets cyber canoniques, contexte et transitions métier |
| Chatbot / copilote | interaction assistée | IA optionnelle dans des workspaces complets |
| Agent orchestrator | coordination d'agents | Studio gouverné, products opérationnels propriétaires |
| Vendor aggregator | accès à plusieurs outils | UX orientée activité et source de vérité commune |

## Garde-fous de positionnement

- Ne pas dire que CMDR remplace immédiatement tous les outils.
- Ne pas dire que toutes les capabilities sont natives ou livrées.
- Ne pas positionner l'IA comme condition d'usage.
- Ne pas réduire la valeur à une console unique.
- Ne pas masquer les moteurs intégrés ou leur provenance.

## Question de validation marché

Les messages, segments d'achat et comparaisons commerciales nécessitent des preuves et recherches futures. Ils ne sont pas inventés dans cette phase.

## Critère d'acceptation

**Given** une comparaison à un SOAR,  
**When** le positionnement est présenté,  
**Then** il explique Evidence, autorité et Result vérifié sans affirmer que l'orchestration native complète est déjà livrée.
