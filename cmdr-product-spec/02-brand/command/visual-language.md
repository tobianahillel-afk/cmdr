---
id: command-visual-language
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-004
  - REQ-BRAND-008
  - REQ-PROD-013
  - REQ-PROD-007
---
# Langage visuel Command

## Surfaces

- Canvas porte la situation générale.
- Surface porte work queues, tables et listes.
- White porte document, handover ou détail sélectionné.
- Mist indique regroupement ou sélection.
- Une surface Ink est rare et réservée à un contexte précis, jamais au shell complet.

## Tables et listes

- colonnes alignées sur la décision : priorité, objet, owner, temps, impact, prochaine action ;
- densité compacte et lignes clairement séparées ;
- badges limités ;
- changement en direct sans déplacement de ligne incontrôlé ;
- surbrillance locale et temporaire pour mise à jour.

## Priorité

La priorité combine ordre, libellé, raison, conséquence et temps. La couleur sémantique peut compléter, jamais remplacer. Un score reste accompagné de facteurs.

## Temps et handover

- temps relatif accompagné d'une date absolue accessible ;
- fraîcheur visible ;
- SLA et échéances placés près de l'action ;
- handover structuré : situation, changements, décisions, risques, owner et suivi.

## Impact

L'impact métier est présenté par périmètre, services, utilisateurs et conséquences. Il ne devient pas une grande carte rouge.

## Cartes et graphiques

- graphiques compacts, small multiples et tendances ;
- carte seulement si emplacement ou juridiction change la décision ;
- aucune carte décorative ;
- légende, période, source et état de données obligatoires.

## Navigation

Juniper marque destination et sélection. Les autres niveaux restent Ink/Slate. La navigation n'utilise pas une succession de capsules.

## Traitement de la mise à jour

Les mises à jour live conservent position, scroll et sélection. Un changement significatif reçoit un repère temporel ou une annotation, pas une animation alarmiste.

## Critère d'acceptation

**Given** une Work Queue avec incident critique,  
**When** la couleur est supprimée,  
**Then** l'ordre, l'owner, l'échéance, l'impact et la prochaine action restent lisibles.
