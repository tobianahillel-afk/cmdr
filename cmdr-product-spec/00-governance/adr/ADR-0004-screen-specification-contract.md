---
id: ADR-0004-screen-specification-contract
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-010
---
# ADR-0004 — Contrat qualitatif des écrans

## 1. Identifiant

`ADR-0004`

## 2. Titre

Contrat qualitatif des écrans

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-UX-010`

## 8. Contexte

Les 61 fichiers actuels possèdent 27 sections mais réutilisent largement le même texte.

## 9. Problème

La présence d'un titre de section ne garantit ni une UI spécifique ni un comportement testable.

## 10. Forces en présence

Comparabilité des specs, liberté de conception, accessibilité, testabilité et prévention du remplissage automatique.

## 11. Options étudiées

1. Aucun contrat.
2. Un template rigide rempli automatiquement.
3. Un contrat de couverture avec contenu spécifique et critères observables.

## 12. Décision

Conserver les 27 catégories obligatoires, mais exiger pour chacune une application locale substantive ou une mention explicite non applicable avec justification.

## 13. Justification

Le contrat maintient une couverture commune sans confondre structure et qualité.

## 14. Conséquences positives

- Revue plus fiable.
- Écrans pilotes comparables.
- États et actions testables.

## 15. Conséquences négatives

- Temps de rédaction supérieur.
- Certains écrans simples auront des sections non applicables.

## 16. Risques

- Retour au texte générique.
- Création de pages pour satisfaire le template.

## 17. Effets sur la navigation

Une page n'est créée que pour un objectif distinct ; vues, modes, filtres et panneaux restent des sous-structures.

## 18. Effets sur les objets

Les objets sont référencés et projetés, jamais redéfinis.

## 19. Effets sur les permissions

La matrice action × état × permission doit être explicite pour les actions sensibles.

## 20. Effets sur les parcours

Les transitions conservent contexte, sélection et retour réel.

## 21. Effets sur les autres documents

`04-experience-architecture/`, templates, screen register, Phase 6.

## 22. Migration

Les six écrans pilotes seront réécrits avant propagation aux autres écrans.

## 23. Critères de réévaluation

Réévaluer après revue des six pilotes et mesures de duplication.

## 24. Questions encore ouvertes

- Le niveau exact de densité par rôle reste OPEN-010.
