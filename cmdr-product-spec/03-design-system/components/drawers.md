---
id: component-drawers
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
---

# Drawers

## 1. Rôle

Héberger une interaction secondaire courte. Le composant fournit un contrat partagé ; le produit propriétaire fournit objets, actions et données.

## 2. Source canonique

Ce fichier possède l'anatomie et l'interaction. Les objets viennent de `05-domain-model/`, les permissions de `14-security-permissions-and-trust/`, le langage de `15-content-and-language/`.

## 3. Anatomy

title, scope, body, actions, close.

## 4. Variants

right, left, bottom responsive.

## 5. Sizes

Compact, standard et comfortable consomment les dimensions de `../foundations/density.md`. Une taille ne change ni contenu ni permission.

## 6. States

Default, hover, focus, selected lorsque applicable, disabled, read-only, loading, empty, partial, error, offline et permission-denied. Disabled explique la raison ; permission-denied n'est pas disabled.

## 7. Behavior

Un seul drawer; pas de navigation complexe ni investigation. Les mises à jour préservent focus, sélection et données valides. Toute action à effet indique portée et résultat.

## 8. Keyboard

focus moves to title/first control; trap; Esc close; restore trigger.

## 9. Accessibility

Nom accessible, focus visible, ordre logique, contraste, annonces proportionnées et alternatives aux interactions visuelles. Target minimum effectif 44×44 px lorsque nécessaire.

## 10. Responsive

Wide conserve la composition complète ; standard réduit les détails ; compact utilise drawer/full-screen ou alternative structurée. Aucune fonction critique ne disparaît sans explication.

## 11. Tokens consommés

`color.surface.*`, `color.text.*`, `color.border.*`, `color.focus.ring`, `space.*`, `radius.*`, `motion.*` et `component.drawers.*`. Aucun hex de marque direct.

## 12. Content rules

Libellés précis, verbes d'action, statut textuel, owner/source/date lorsque pertinents. Les identifiants complets restent copiables.

## 13. Adaptations produit

Les produits peuvent changer ordre, accent résolu, densité par défaut et contenu ; ils ne changent pas clavier, sémantique, statut ou structure de base.

## 14. Usages autorisés

Seulement lorsque le rôle décrit répond au besoin utilisateur et que le produit possède l'action ou la projection.

## 15. Usages interdits

Créer un clone local, masquer provenance/permission, utiliser couleur seule, transformer une suggestion IA en fait ou Decision.

## 16. Exemple

Une sélection au clavier ouvre le détail sans perdre l'origine ; une fermeture rend le focus au déclencheur.

## 17. Critère d’acceptation

**Given** un utilisateur clavier, thème sombre, densité compacte, données partielles et permissions limitées,  
**When** il parcourt et active le composant,  
**Then** l'état et la portée sont annoncés, le focus est visible et restauré, aucune donnée interdite n'apparaît, les tokens sémantiques sont utilisés et une alternative non-IA demeure disponible.
