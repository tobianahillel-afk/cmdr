---
id: gov-readme
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-012
  - REQ-OBJ-001
  - REQ-OBJ-012
---
# Gouvernance documentaire CMDR

## Mission du domaine

`00-governance/` définit comment une décision entre dans le référentiel, comment elle devient une règle canonique, qui en est responsable, comment les consommateurs la référencent et quelles preuves sont nécessaires avant revue, validation, implémentation ou dépréciation.

La gouvernance protège trois propriétés du référentiel :

1. une décision explicite ne peut pas être remplacée silencieusement ;
2. un concept ne possède qu'une source normative active ;
3. un statut décrit un niveau de preuve réel, jamais la seule présence d'un fichier.

## Propriétaire et statut

- **Propriétaire :** Product Architecture.
- **Statut du domaine :** Draft, en préparation de revue produit.
- **Consommateurs :** tous les domaines CMDR, les auteurs, reviewers, approbateurs, équipes Design, Security, Engineering, QA et Content.

## Requirement IDs principaux

`REQ-PROD-006`, `REQ-PROD-009`, `REQ-PROD-012`, `REQ-OBJ-001` à `REQ-OBJ-012`, `REQ-UX-001`, `REQ-AI-002`.

## Ordre de lecture

1. [`source-material/README.md`](source-material/README.md) — décisions et contraintes fournies.
2. [`source-of-truth-policy.md`](source-of-truth-policy.md) — hiérarchie normative et résolution des contradictions.
3. [`documentation-rules.md`](documentation-rules.md) — règles d'écriture, de référence et d'acceptation.
4. [`ownership-register.md`](ownership-register.md) — propriété et consommation des concepts.
5. [`terminology-rules.md`](terminology-rules.md) — vocabulaire structurel.
6. [`status-lifecycle.md`](status-lifecycle.md) — états documentaires et preuves.
7. [`review-and-approval-process.md`](review-and-approval-process.md) — revues et approbations.
8. [`dependency-register.md`](dependency-register.md) — dépendances produit et documentaires.
9. [`decision-log.md`](decision-log.md) puis `adr/` — décisions d'architecture.
10. `registers/` — index canoniques.

## Documents normatifs

- `documentation-rules.md`
- `source-of-truth-policy.md`
- `ownership-register.md`
- `terminology-rules.md`
- `naming-conventions.md`
- `status-lifecycle.md`
- `review-and-approval-process.md`
- `dependency-register.md`
- ADR approuvés, lorsqu'ils le seront
- registres canoniques

## Documents d'appui et de migration

- `source-material/` conserve l'intention source, sans remplacer les spécifications propriétaires.
- `legacy-screen-mapping.md` décrit les migrations et n'est pas une source UX.
- `contribution-workflow.md` décrit l'exécution Git et documentaire.
- les fichiers `deprecated` ne sont jamais normatifs.

## Inventaire des enfants

### Règles et processus

- [`documentation-rules.md`](documentation-rules.md)
- [`source-of-truth-policy.md`](source-of-truth-policy.md)
- [`ownership-register.md`](ownership-register.md)
- [`terminology-rules.md`](terminology-rules.md)
- [`naming-conventions.md`](naming-conventions.md)
- [`status-lifecycle.md`](status-lifecycle.md)
- [`review-and-approval-process.md`](review-and-approval-process.md)
- [`dependency-register.md`](dependency-register.md)
- [`contribution-workflow.md`](contribution-workflow.md)
- [`legacy-screen-mapping.md`](legacy-screen-mapping.md)
- [`decision-log.md`](decision-log.md)

### Compatibilité documentaire

- `document-status-model.md` — déprécié au profit de `status-lifecycle.md`.
- `documentation-lifecycle.md` — déprécié au profit de `status-lifecycle.md`.
- `review-and-approval-policy.md` — déprécié au profit de `review-and-approval-process.md`.

### ADR

- `adr/ADR-0001-product-separation.md`
- `adr/ADR-0002-single-source-of-truth.md`
- `adr/ADR-0003-canonical-object-chain.md`
- `adr/ADR-0004-screen-specification-contract.md`
- `adr/ADR-0005-page-view-mode-filter-rules.md`
- `adr/ADR-0006-endpoint-agent-ownership.md`
- `adr/ADR-0007-agentic-studio-placement.md`

### Registres

- `registers/capability-register.md`
- `registers/component-register.md`
- `registers/object-register.md`
- `registers/permission-register.md`
- `registers/screen-register.md`

### Sources de décision

Le contenu complet de `source-material/` est inventorié dans son propre README. Les exigences sont reliées au référentiel par `source-material/requirements-traceability-matrix.md`.

## Dépendances

- `01-product-vision/` applique les décisions de mission, principes et frontières.
- `04-experience-architecture/` applique les décisions page, workspace, vue, mode, filtre et contexte.
- `05-domain-model/` applique la propriété et les distinctions d'objets.
- `14-security-permissions-and-trust/` possède le modèle de permission.
- `16-quality-and-validation/` vérifie structure, contenu, liens et traçabilité.

## Décisions ouvertes

Les questions `OPEN-001` à `OPEN-015` restent dans `source-material/unresolved-decisions.md`. Une question n'est pas recopiée comme texte libre : les documents concernés citent son identifiant et décrivent seulement son impact local.

## Règles de contribution

- travailler hors `main` ;
- limiter chaque lot à un objectif documentaire ;
- lire les sources et dépendants avant modification ;
- mettre à jour la traçabilité dans le même lot ;
- ne pas passer un document en `validated` sans approbations enregistrées ;
- ne pas fusionner une contradiction active ou un P0 connu.

## Critères d'acceptation

**Given** une décision source, un document canonique et plusieurs consommateurs,  
**When** une modification est proposée,  
**Then** la priorité des sources est identifiable, le propriétaire canonique est unique, les Requirement IDs sont reliés, les dépendants sont listés, les contradictions sont résolues ou enregistrées comme ouvertes et aucun document de priorité inférieure ne reste normatif avec une règle concurrente.
