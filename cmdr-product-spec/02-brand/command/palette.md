---
id: command-palette
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-004
  - REQ-BRAND-005
  - REQ-BRAND-008
  - REQ-PROD-013
---
# Palette canonique Command

Ces huit valeurs sont décidées. Les couleurs sémantiques ne sont pas définies ici et appartiennent à la Phase 3.

| Token | Valeur | Rôle |
|---|---|---|
| Command Canvas | `#EDF2EF` | fond opérationnel général |
| Command Surface | `#F7F9F7` | surface de travail principale |
| Command White | `#FCFDFB` | surface la plus élevée ou document |
| Command Juniper | `#1F6658` | accent principal, navigation et sélection |
| Command Fern | `#4D897A` | accent secondaire, continuité et série de données |
| Command Mist | `#D5E4DE` | sélection douce, regroupement et contexte |
| Command Slate | `#53645F` | texte secondaire et métadonnées |
| Command Ink | `#17221F` | texte principal et surfaces sombres |

## Hiérarchie de surfaces

1. Canvas pour le shell et les zones de situation.
2. Surface pour tables, listes et blocs opérationnels.
3. White pour document, inspection ou focus fort.
4. Mist pour sélection, regroupement et état local.
5. Ink pour une surface sombre ponctuelle, jamais pour transformer Command en salle de contrôle nocturne.

## Accents

- Juniper porte navigation active, contrôle sélectionné, lien principal et série de données dominante.
- Fern porte série secondaire, tendance ou continuité.
- Mist peut montrer une sélection étendue sans réduire la lisibilité.
- Moss et Ember peuvent signer la marque, mais ne remplacent pas Juniper/Fern ni la sémantique.

## Contraste mesuré

| Paire | Ratio approximatif | Usage |
|---|---:|---|
| Command Ink sur Canvas | 14.43:1 | texte normal |
| Command Slate sur Canvas | 5.53:1 | texte normal |
| Command Juniper sur Canvas | 5.98:1 | texte et contrôle |
| Command Fern sur Canvas | 3.58:1 | accent, grand texte ou data-viz ; pas petit texte normal |
| Command Ink sur White | 16.02:1 | texte normal |
| Command Mist sur Command Ink | 12.44:1 | surface ou texte clair sur fond sombre |

## Sélection et focus

La sélection combine Mist, une bordure Juniper et un indicateur non chromatique. Le focus doit rester distinct de la sélection ; son token et sa géométrie seront définis en Phase 3.

## Navigation

- Juniper marque le produit et la destination active.
- Les éléments inactifs utilisent Ink ou Slate.
- La navigation n'utilise pas un remplissage vert sur chaque élément.
- Ember n'est pas un indicateur de navigation.

## Graphiques

- Juniper est la série prioritaire.
- Fern est une série secondaire.
- Slate, Mist et les neutres structurent comparaison et contexte.
- Les couleurs sémantiques restent réservées à leur signification.
- Une légende, une unité, une période et une source sont obligatoires.

## Badges non sémantiques

Mist, Juniper ou Slate peuvent identifier une catégorie, un owner ou une vue. Un badge de sévérité, de succès, d'échec ou de risque utilise les tokens sémantiques de Phase 3.

## Usages éditoriaux

Juniper peut souligner un titre, une règle ou une donnée principale. Fern peut marquer une continuité. Ember peut attirer l'attention sur une échéance éditoriale, à condition que le texte explique la signification.

## Interdictions

- interface entièrement verte ;
- rouge permanent pour évoquer un SOC ;
- Juniper utilisé comme succès ;
- Fern utilisé comme seul canal ;
- grands aplats sombres de type HUD ;
- radar ou carte décorative ;
- ajout de couleurs produit concurrentes.

## Critère d'acceptation

**Given** une situation Command contenant priorité, sélection, statut et série de données,  
**When** la palette est appliquée,  
**Then** l'accent produit, la sélection et la sémantique restent distincts et l'interface évoque une priorité calme plutôt qu'un centre militaire.
