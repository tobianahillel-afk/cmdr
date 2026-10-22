---
id: review-and-approval-process
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-009
  - REQ-PROD-012
  - REQ-PROD-019
  - REQ-PROD-056
---
# Processus de revue et d'approbation

## Rôles

- **Auteur :** prépare le contenu et les preuves.
- **Propriétaire :** assume la décision et la maintenance.
- **Reviewer :** évalue un angle défini et formule des demandes.
- **Approbateur :** accepte le périmètre au nom de sa responsabilité.
- **QA / Traceability :** vérifie liens, statuts, Requirement IDs, dépendants et critères.
- **Documentation Governance :** arbitre la conformité du processus, pas le contenu produit.

Une même personne peut cumuler auteur et propriétaire, mais une action à risque ou un sujet de séparation des tâches peut exiger des approbateurs distincts.

## Matrice de revue

| Type | Product | UX / Design | Security | Engineering | Content |
|---|---:|---:|---:|---:|---:|
| Vision produit | requis | selon impact | selon impact | consultation | requis |
| Frontière produit | requis | requis | requis si confiance ou action | requis | requis |
| Identité de marque | requis | requis | non systématique | consultation | requis |
| Écran | requis | requis | selon données et actions | requis | requis |
| Permission | consultation | selon UX | requis | requis | consultation |
| Objet | requis | selon exposition | selon données | requis | consultation |
| Contrat technique | consultation | non systématique | requis | requis | consultation |
| Parcours interproduit | requis | requis | selon actions | requis | requis |
| ADR transversale | propriétaires concernés | selon impact | selon impact | selon impact | consultation |

## Déroulement

1. **Préparation** — auteur, propriétaire, exigences, sources et dépendants.
2. **Auto-contrôle** — contenu substantiel, liens, questions, AC et absence de P0 connu.
3. **Passage `in-review`** — seulement si les conditions du lifecycle sont réunies.
4. **Revue spécialisée** — chaque reviewer rend `approve`, `request-changes` ou `comment`.
5. **Résolution** — les demandes bloquantes sont corrigées ou transformées en décision ouverte avec propriétaire.
6. **Approbation** — les approbateurs requis signent le périmètre exact.
7. **Mise à jour** — matrice, journal, dépendants et statut.
8. **Révision future** — date ou déclencheur de réévaluation.

## Autorité de classification des capacités

La classification `native`, `integrated`, `temporary-integration`, `planned` ou `out-of-scope` suit ce processus :

1. le propriétaire de la capacité propose la classification et fournit son résultat utilisateur, le responsable du moteur, les limites et les dépendances ;
2. le Product Lead propriétaire et Product Architecture approuvent la classification ;
3. Security revoit toute capacité touchant aux actions, à la confiance, aux secrets, aux preuves ou aux données sensibles ;
4. Engineering revoit obligatoirement toute revendication `native` ou `integrated` avant qu'elle soit présentée comme livrée ;
5. QA and Traceability vérifie les preuves, met à jour le Capability Register et la matrice ;
6. une `temporary-integration` exige une stratégie de remplacement ou une date de réévaluation.

Cette décision résout `OPEN-009` et couvre `REQ-PROD-019` et `REQ-PROD-056`.

## Enregistrement d'une revue

Chaque revue enregistre :

```yaml
document:
version:
author:
owner:
reviewer:
review_scope:
result: approve|request-changes|comment
date:
evidence:
blocking_findings:
resolved_by:
next_review:
```

## Motifs de refus de validation

- contradiction active ;
- question bloquante sans décision ;
- contenu générique ou dupliqué ;
- propriétaire absent ;
- Requirement IDs inconnus ;
- critère non testable ;
- capability `planned` présentée comme livrée ;
- définition d'objet ou permission recopiée dans un consommateur ;
- liens cassés ou dépendants ignorés.

## Révision future

Une révision est déclenchée par :

- remplacement d'une décision source ;
- nouvelle ADR approuvée ;
- changement de propriété ;
- évolution de produit affectant les frontières ;
- écart d'implémentation ;
- incident ou audit révélant une hypothèse incorrecte ;
- date de révision planifiée.

## Critère d'acceptation

**Given** un document candidat à `validated`,  
**When** le processus est consulté,  
**Then** l'auteur, le propriétaire, chaque reviewer, chaque approbateur, le résultat, les preuves, les demandes résolues et la prochaine révision sont identifiables.
