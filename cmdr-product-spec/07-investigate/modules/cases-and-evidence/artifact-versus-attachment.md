---
id: investigate-artifact-versus-attachment
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: proposal
requirements:
  - REQ-PROD-061
open_decisions:
  - OPEN-014
---

# Artifact versus Attachment — Phase 4B.1 proposal

## État

`OPEN-014` reste ouverte. Aucun objet Attachment canonique n’existe dans l’Object Register.

## Usages constatés

- **Artifact** : élément technique analysable, versionné, lié à un Case, pouvant avoir dérivés et rôle analytique.
- **Attachment** : fichier joint à un conteneur de collaboration ou de reporting, dont l’usage principal est documentaire.
- **Evidence** : qualification explicite d’une source ou d’un Artifact dans un Case.
- **Report, Note et Comment** : conteneurs Shared ou concepts encore insuffisamment formalisés.

## Chevauchements

Un même contenu peut être joint puis qualifié comme Artifact, ou un Artifact peut être cité dans un Report. Les relations et métadonnées ne prouvent pas que les objets doivent fusionner.

## Options à étudier

| Option | Avantage | Risque |
|---|---|---|
| Attachment distinct | permissions/rétention par conteneur | duplication de contenu et relations |
| contenu générique avec rôles | moins de doublons | modèle abstrait complexe |
| Attachment comme référence, Artifact comme qualification | progression explicite | migration et UX à formaliser |
| convergence complète | simplicité apparente | perte de distinction documentaire/analytique |

## Besoins de décision

Ownership, rétention, legal hold, versioning, permissions, qualification vers Artifact/Evidence, migration, déduplication et stockage.

Aucune option n’est retenue en Phase 4B.1.
