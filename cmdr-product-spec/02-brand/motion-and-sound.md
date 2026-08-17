---
id: motion-and-sound
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
  - REQ-AI-007
---
# Mouvement et son

## Rôle du mouvement

Le mouvement explique :

- l'origine et la destination d'un objet ;
- une insertion, suppression ou reclassification ;
- une transition interproduit ;
- la progression d'un run ;
- l'ouverture d'un niveau de détail ;
- la relation entre sélection et Inspector.

Il ne sert pas à maintenir l'attention artificiellement.

## Rythme

La Phase 3 définira les tokens. Les plages de conception à valider sont :

- immédiat : environ 80–120 ms pour retour local ;
- standard : environ 160–220 ms pour contrôle ou panel ;
- contextuel : environ 240–320 ms pour transition de workspace ;
- progression longue : état réel, non animation en boucle simulant l'activité.

Le mouvement doit pouvoir être interrompu et ne pas retarder l'action.

## Courbes

Les courbes favorisent départ lisible et arrivée stable. Les rebonds, élasticités, overshoot décoratif et accélérations agressives sont interdits dans les workflows opérationnels.

## Reduced motion

En reduced motion :

- les translations deviennent fondu ou remplacement direct ;
- aucune information n'est perdue ;
- les progressions restent exprimées par état et texte ;
- les graphes ne dépendent pas d'une animation pour être compris.

## Mouvement produit

- Command : mise à jour stable, priorité qui change sans saut de layout.
- Investigate : relation entre sélection, canvas et Inspector.
- Govern : progression explicite demande → décision → run → résultat.
- Studio : exécution de run et état des nodes, sans flux lumineux permanent.

## Son

Aucun son n'est activé par défaut pour une information courante. Un son éventuel exige :

- opt-in administrable ou utilisateur ;
- signification unique ;
- alternative visuelle ;
- fréquence limitée ;
- test de contexte partagé et d'accessibilité.

Aucune identité sonore n'est décidée dans cette phase.

## Interdictions

- pulsation continue ;
- glow animé ;
- ticker permanent ;
- son d'alarme pour chaque alerte ;
- animation destinée à faire croire qu'un moteur travaille ;
- transition qui cache provenance ou changement d'état.

## Critère d'acceptation

Une transition est supprimée si elle ne peut pas être décrite comme causalité, continuité, feedback ou orientation.
