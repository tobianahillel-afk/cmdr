---
id: brand-readme
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-005
  - REQ-BRAND-006
  - REQ-BRAND-007
  - REQ-BRAND-008
---
# Marque CMDR

## Mission du domaine

`02-brand/` traduit la vision produit en un langage visuel reconnaissable, durable et exploitable. Il définit l'identité de la marque mère, les règles communes, la palette CMDR, la palette Command, les propositions non validées des autres produits et les contraintes que le Design System devra transformer en tokens pendant la Phase 3.

La marque doit rendre visibles la maîtrise, la responsabilité, la provenance et la continuité. Elle ne sert ni à simuler la complexité ni à transformer le produit en spectacle.

## Propriétaire et statut

- **Propriétaire :** Brand Design Lead.
- **Statut du domaine :** Draft, prêt pour revue de marque mais non validé.
- **Direction décidée :** Operational Editorial Modernism.
- **Décisions encore ouvertes :** `OPEN-001`, `OPEN-002`, `OPEN-003`, `OPEN-004` et `OPEN-016`.

## Ordre de lecture

1. [`operational-editorial-modernism.md`](operational-editorial-modernism.md)
2. [`brand-personality.md`](brand-personality.md)
3. [`brand-essence-and-signature.md`](brand-essence-and-signature.md)
4. [`visual-principles.md`](visual-principles.md)
5. [`brand-architecture.md`](brand-architecture.md)
6. [`forbidden-directions.md`](forbidden-directions.md)
7. [`cmdr/palette.md`](cmdr/palette.md)
8. [`command/palette.md`](command/palette.md)
9. [`cmdr/typography.md`](cmdr/typography.md)
10. [`logo-and-wordmark.md`](logo-and-wordmark.md)
11. [`shape-and-surfaces.md`](shape-and-surfaces.md)
12. [`imagery-and-iconography.md`](imagery-and-iconography.md)
13. [`illustration-system.md`](illustration-system.md)
14. [`data-visualization-language.md`](data-visualization-language.md)
15. [`motion-and-sound.md`](motion-and-sound.md)
16. [`editorial-expression.md`](editorial-expression.md)
17. [`ai-visual-expression.md`](ai-visual-expression.md)
18. [`accessibility-and-brand.md`](accessibility-and-brand.md)
19. [`brand-quality-gates.md`](brand-quality-gates.md)
20. les identités Command, Investigate, Govern et Studio.

## Sources normatives

- `cmdr/palette.md` est l'unique source des sept couleurs CMDR.
- `command/palette.md` est l'unique source des huit couleurs Command.
- `operational-editorial-modernism.md` porte la direction visuelle.
- `forbidden-directions.md` porte les interdictions de marque.
- `cmdr/typography.md` porte l'étude typographique, sans fermer `OPEN-004`.
- `investigate/palette-proposals.md`, `govern/palette-proposals.md` et `studio/palette-proposals.md` contiennent uniquement des propositions.
- Les valeurs sémantiques, échelles et tokens d'implémentation appartiennent à la Phase 3.

## Inventaire

### Marque mère

- `cmdr/README.md`
- `cmdr/identity.md`
- `cmdr/logo.md`
- `cmdr/palette.md`
- `cmdr/typography.md`
- `cmdr/usage.md`
- `cmdr/voice.md`

### Identités produit

- `command/` — identité décidée et palette canonique.
- `investigate/` — identité définie, palette ouverte.
- `govern/` — identité définie, palette ouverte.
- `studio/` — identité définie, palette ouverte.

### Système transversal

Les fichiers racine de ce domaine définissent architecture de marque, personnalité, principes, surfaces, imagerie, illustration, iconographie, data visualisation, mouvement, voix et contrôles.

## Frontière avec le Design System

La marque définit l'intention, les couleurs canoniques, les comportements visuels et les interdictions. Le Design System définit les tokens, composants, valeurs d'espacement, rayons, élévations et règles d'implémentation. Un fichier de Phase 3 référence les sources de marque ; il ne recopie pas une palette.

## Contribution

Toute modification doit :

- citer les Requirement IDs ;
- distinguer `decided`, `proposed` et `open` ;
- conserver les valeurs canoniques exactes ;
- démontrer la différenciation autrement que par la couleur ;
- vérifier clair, sombre, contraste, daltonisme et reduced motion ;
- éviter toute promesse visuelle de capacité non livrée ;
- mettre à jour les décisions ouvertes et la matrice.

## Critère d'acceptation

**Given** un designer qui ne connaît pas CMDR,  
**When** il lit ce domaine dans l'ordre recommandé,  
**Then** il peut reconnaître la marque, expliquer sa personnalité, appliquer ses palettes décidées, produire des propositions compatibles pour les palettes ouvertes, différencier les quatre produits sans changer de Design System et rejeter les clichés cyber ou IA.
