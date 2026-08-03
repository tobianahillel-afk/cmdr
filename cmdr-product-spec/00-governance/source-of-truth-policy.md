---
id: source-of-truth-policy
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Politique de source de vérité

## Objectif

Garantir qu’un concept n’est défini qu’une seule fois et que tous les autres documents le référencent.

## Ordre de priorité

1. ADR approuvés et politique de gouvernance.
2. Documents canoniques explicitement enregistrés dans les registres.
3. Spécifications produit et module pour leur comportement local.
4. Spécifications d’écran pour la composition et les interactions locales.
5. Parcours, exemples et assets, qui ne sont jamais normatifs sur les objets.

## Règles

- Une page ne redéfinit ni objet, ni permission, ni composant partagé.
- Les palettes vivent uniquement sous `02-brand/`.
- Les composants vivent uniquement sous `03-design-system/components/`.
- Les objets vivent uniquement sous `05-domain-model/objects/`, sauf Permission et Reporting Engine dont les sources sont explicitement externes.
- Les questions ouvertes restent dans le document propriétaire.
- Les risques produit vivent dans `01-product-vision/product-risks.md`.
- Les dépendances vivent dans `dependency-register.md` et `18-roadmap-and-releases/dependency-roadmap.md`.
- Un document remplacé est supprimé de l’espace actif; `99-archive/` contient seulement une notice non normative ou un contenu réellement requis pour l’historique.

## Critères d’acceptation

- Aucun concept du registre n’a deux sources actives.
- Chaque référence locale résout.
- Les chemins et fichiers respectent le kebab-case, à l’exception du préfixe ADR imposé.
- Le manifeste attendu et l’arborescence réelle sont identiques.
