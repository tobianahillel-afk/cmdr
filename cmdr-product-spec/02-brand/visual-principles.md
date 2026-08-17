---
id: brand-visual-principles
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
  - REQ-PROD-003
  - REQ-PROD-005
  - REQ-PROD-007
  - REQ-PROD-008
---
# Principes visuels normatifs

| Principe | Règle | Conséquence | Anti-pattern | Exemple CMDR |
|---|---|---|---|---|
| Information before decoration | Toute surface doit soutenir lecture, comparaison, décision ou orientation. | Les éléments sans fonction sont retirés. | Carte, illustration ou gradient occupant de l'espace sans améliorer la tâche. | Une timeline montre événements, sources et écarts ; elle n'utilise pas une lueur décorative. |
| Typography carries hierarchy | Taille, graisse, rythme, alignement et longueur de ligne portent la hiérarchie. | La couleur reste disponible pour sélection, provenance et sémantique. | Cinq couleurs pour différencier cinq niveaux de titre. | Le titre, le statut, l'owner et la date forment un bloc éditorial stable. |
| Controlled density | La densité est volontaire, structurée et adaptée à l'activité. | Les vues expertes peuvent être compactes sans masquer l'essentiel. | Landing page vide ou accumulation de panneaux simultanés. | Investigate révèle détails techniques dans l'Inspector et la console. |
| Accountability is visible | Auteur, owner, source, statut, date, trace et prochaine action sont visibles au bon niveau. | Les contenus humains, déterministes et automatisés sont distingués. | Résumé sans provenance ou action système anonyme. | Une proposition IA affiche run, modèle, Evidence et reviewer attendu. |
| Calm over drama | Le shell reste stable et lisible, même sous urgence. | Les signaux forts sont locaux, temporaires et hiérarchisés. | Rouge permanent, clignotement ou animation alarmiste. | Une action critique utilise texte, classe, impact et confirmation, pas un écran entier rouge. |
| Editorial rhythm | Les pages suivent une cadence de titres, métadonnées, sections, notes, références et tableaux. | Le lecteur peut parcourir puis approfondir. | Mosaïque de widgets sans ordre de lecture. | Govern présente demande, preuve, conditions, décision et trace dans un ordre documentaire. |
| Functional differentiation | Les produits se distinguent par activité, rythme et motif fonctionnel, pas uniquement par couleur. | Un composant identique conserve sa fonction et son comportement. | Quatre thèmes sans autre différence qu'un accent. | Command met en avant priorité et temps ; Investigate provenance et comparaison. |
| Provenance before personality | Une signature visuelle ne doit jamais masquer la source ou le statut. | La marque encadre la preuve au lieu de la remplacer. | Badge « intelligent » sans méthode ni données. | Un Finding automatisé garde la même structure qu'un Finding humain avec origine explicite. |
| One family, several instruments | Tous les produits utilisent le même vocabulaire de surfaces, typographie, iconographie et motion. | Les variations restent contrôlées par les sources de marque. | Chaque produit invente ses propres cartes, icônes et statuts. | Inspector et Context Bar restent communs, leur contenu varie. |

## Arbitrage

En cas de conflit :

1. lisibilité, accessibilité et vérité précèdent l'expression ;
2. sémantique fonctionnelle précède couleur de marque ;
3. provenance et responsabilité précèdent gain de place ;
4. stabilité du système commun précède différenciation produit ;
5. reduced motion et usage prolongé précèdent effet spectaculaire.

## Critère d'acceptation

**Given** une maquette présentant une information critique,  
**When** elle est examinée contre ces principes,  
**Then** l'information reste lisible sans couleur, son owner et sa provenance sont accessibles, l'urgence est locale et la surface conserve le vocabulaire commun CMDR.
