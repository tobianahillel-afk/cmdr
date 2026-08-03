---
id: cmdr-palette
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-005
  - REQ-BRAND-006
---
# Palette canonique CMDR

Ces valeurs sont décidées. Toute copie ou modification crée une contradiction. Les tokens d'implémentation de Phase 3 devront référencer ce fichier.

| Token | Valeur | Rôle de marque | Usage fonctionnel autorisé | Usage interdit |
|---|---|---|---|---|
| CMDR Ink | `#18201F` | autorité calme et profondeur | texte principal, surfaces sombres, wordmark monochrome | noir générique pour tous les fonds |
| CMDR Bone | `#F2EFE6` | chaleur éditoriale | canvas clair, documents, surfaces marketing | texte clair sur surface claire |
| CMDR Moss | `#8EA68F` | continuité et mémoire | accent de marque, liaison, repère non sémantique | succès universel ou petit texte sur Bone |
| CMDR Ember | `#D7644B` | engagement et point d'attention | accent de marque, échéance ou transition éditoriale avec texte | critique universel, erreur ou long texte sur Bone |
| CMDR Fog | `#D7DDD7` | respiration et séparation | surface secondaire, règle large, fond de métadonnée | texte principal |
| CMDR Graphite | `#4B5451` | précision secondaire | texte secondaire et métadonnées sur Bone | texte sur Ink sans vérification |
| CMDR Stone | `#C8C6BE` | matérialité neutre | bordure, division, surface discrète | texte principal ou état |

## Associations

- **Ink + Bone** : contraste fondateur.
- **Graphite + Bone** : métadonnées et texte secondaire.
- **Fog/Stone + Ink** : surfaces ou repères sur thème sombre, après validation.
- **Moss + Ember** : signature de marque limitée, jamais paire sémantique succès/erreur.

## Contraste mesuré

Les ratios suivants servent d'information de conception ; les tests d'implémentation restent obligatoires :

| Paire | Ratio approximatif | Usage |
|---|---:|---|
| Ink sur Bone | 14.44:1 | texte normal et grands titres |
| Graphite sur Bone | 6.80:1 | texte normal |
| Moss sur Bone | 2.28:1 | pas de texte normal ; accent de surface ou grand motif |
| Ember sur Bone | 3.14:1 | pas de texte normal ; grand texte ou accent seulement selon validation |
| Bone sur Ink | 14.44:1 | texte clair sur thème sombre |
| Moss sur Ink | 6.32:1 | texte/accent possible après test |
| Ember sur Ink | 4.60:1 | texte normal possible à la limite AA, à tester selon poids et contexte |

## Thème clair

- Bone peut être le canvas de marque ou éditorial.
- Ink porte texte principal et wordmark.
- Graphite porte métadonnées.
- Fog et Stone structurent sans créer une mosaïque de cartes.
- Moss et Ember restent rares.

## Thème sombre

- Ink sert de référence, mais le Design System doit créer des surfaces sombres teintées distinctes plutôt qu'un fond uniforme ou noir absolu.
- Bone, Fog ou Stone servent de texte/surface selon ratio.
- Moss et Ember gardent leur rôle de marque ; ils ne deviennent pas automatiquement statut.
- Les valeurs dérivées de thème sont des tokens de Phase 3, pas de nouvelles couleurs de marque.

## Accessibilité

- La couleur n'est jamais le seul canal.
- Les accents dont le ratio est insuffisant sont accompagnés de texte, forme ou bordure.
- Les graphiques doivent être testés pour déficiences de perception des couleurs.
- Les états sémantiques utilisent une palette distincte définie en Phase 3.

## Critère d'acceptation

**Given** un token CMDR dans une maquette ou un fichier de tokens,  
**When** sa source est examinée,  
**Then** son nom et sa valeur correspondent exactement à ce tableau et son rôle ne remplace pas une sémantique fonctionnelle.
