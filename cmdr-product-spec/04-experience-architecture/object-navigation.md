---
id: experience-object-navigation
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-006
  - REQ-PROD-006
---
# Navigation des objets


Object Chip pour référence inline ; Inspector pour détail et actions courtes ; page/workspace pour activité durable. Une projection indique produit propriétaire, statut, source et permission.

Ouvrir une relation conserve la sélection et fournit retour. Un consommateur ne redéfinit ni champ, ni lifecycle. Un objet interdit est remplacé par une référence non sensible ou omis.

**Given** une Evidence affichée dans Command, **When** elle est ouverte, **Then** l'Inspector identifie Investigate comme propriétaire et Command ne propose aucune mutation de lifecycle.
