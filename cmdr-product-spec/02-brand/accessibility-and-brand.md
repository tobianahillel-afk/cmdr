---
id: accessibility-and-brand
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-007
  - REQ-BRAND-008
  - REQ-PROD-007
---
# Accessibilité de marque

## Principe

Une identité CMDR n'est valide que si elle reste perceptible, lisible et opérable sans dépendre de couleur, motion, thème sombre ou résolution élevée.

## Couleur

- viser WCAG 2.2 AA pour texte et composants ;
- mesurer toutes les paires réelles, pas seulement les swatches ;
- ne pas utiliser Moss, Ember, Fern ou les accents proposés comme petit texte lorsque le ratio est insuffisant ;
- fournir forme, label, position ou motif en plus de la couleur ;
- tester les visualisations avec plusieurs déficiences de perception des couleurs.

## Typographie

- taille et interligne adaptés à l'usage prolongé ;
- chiffres tabulaires pour tables et mesures ;
- distinction claire entre `0/O`, `1/l/I` dans les contenus techniques ;
- pas de texte important en capitales longues ;
- largeur de ligne contrôlée ;
- zoom à 200 % sans perte de contenu ni recouvrement.

## Thèmes

- clair et sombre sont des expressions égales, pas un thème principal et une dégradation ;
- le sombre évite noir absolu et contrastes extrêmes ;
- les surfaces restent hiérarchisées ;
- les graphes et statuts gardent leur signification.

## Mouvement

- reduced motion complet ;
- pas de clignotement ou pulsation ;
- feedback persistant sous forme d'état ou texte ;
- aucune donnée révélée uniquement par animation.

## Imagerie et illustration

- texte alternatif selon la fonction ;
- légende pour diagrammes complexes ;
- aucune information intégrée uniquement dans une image ;
- personnes représentées sans stéréotype de rôle ou de compétence.

## Iconographie

- hit target et focus relèvent du Design System ;
- labels accessibles pour actions non universelles ;
- icône et couleur ne suffisent pas à indiquer un statut ;
- géométrie lisible aux tailles prévues.

## Validation

Chaque identité et proposition de palette doit fournir :

- ratios de contraste ;
- cas clair et sombre ;
- exemple monochrome ;
- test sans couleur ;
- test zoom et densité ;
- test reduced motion ;
- revue par personnes utilisant des technologies d'assistance avant validation.

## Critère d'acceptation

**Given** une interface CMDR en monochrome, zoomée et reduced motion,  
**When** un utilisateur réalise l'activité principale,  
**Then** produit, hiérarchie, statut, sélection, provenance et action restent compréhensibles.
