---
id: ADR-0005-page-view-mode-filter-rules
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-008
  - REQ-UX-009
---
# ADR-0005 — Distinction Page, Workspace, Vue, Mode et Filtre

## 1. Identifiant

`ADR-0005`

## 2. Titre

Distinction Page, Workspace, Vue, Mode et Filtre

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-UX-001`, `REQ-UX-008`, `REQ-UX-009`

## 8. Contexte

La structure actuelle représente plusieurs variantes de Work Queue et Mission Control comme des écrans séparés.

## 9. Problème

Créer une route pour chaque sous-ensemble ou représentation multiplie navigation, états et maintenance.

## 10. Forces en présence

Clarté des objectifs, liens profonds, conservation du contexte, personnalisation et simplicité.

## 11. Options étudiées

1. Une page par variante.
2. Une page unique sans vues enregistrées.
3. Un workspace par objectif, avec vues, modes et filtres.

## 12. Décision

Un objectif durable correspond à une page ou un workspace. Une vue change le sous-ensemble, un mode la représentation, un filtre réduit temporairement les données.

## 13. Justification

La structure suit le modèle mental de l'utilisateur plutôt que l'inventaire des fichiers.

## 14. Conséquences positives

- Moins de clones.
- Contexte préservé.
- Saved Views réutilisables.

## 15. Conséquences négatives

- Migration des identifiants d'écran actuels.
- Compatibilité des liens profonds à gérer.

## 16. Risques

- Perte d'accès direct à une vue si l'URL n'est pas sérialisée.
- Workspace surchargé.

## 17. Effets sur la navigation

Work Queue et Mission Control deviennent chacun un workspace avec vues et modes.

## 18. Effets sur les objets

Les mêmes objets sont projetés dans les vues ; aucun nouvel objet n'est créé.

## 19. Effets sur les permissions

Changer de vue ne change pas l'autorisation et ne révèle pas d'objets supplémentaires.

## 20. Effets sur les parcours

Les liens entrants restaurent workspace, vue, filtre, sélection et contexte.

## 21. Effets sur les autres documents

`04-experience-architecture/page-view-mode-filter-rules.md`, screen register et écrans Command.

## 22. Migration

Conserver les fichiers actuels comme entrées de migration jusqu'à la Phase 6, puis consolider routes et registre.

## 23. Critères de réévaluation

Réévaluer après validation des écrans pilotes Command.

## 24. Questions encore ouvertes

- Les identifiants legacy à conserver comme alias seront décidés en Phase 6.
