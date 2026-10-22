---
id: CAP-INV-004
title: Event Inspection and Pivot
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-002
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-014
  - REQ-PROD-045
open_decisions:
  - OPEN-014
source-of-truth: canonical
---

# CAP-INV-004 — Event Inspection and Pivot

## 1. Définition

Permet d’inspecter un Telemetry Event en formes raw et rendered selon permission, d’examiner source, parser, champs, enrichissements et provenance, puis de pivoter ou de proposer un lien au Case, Artifact ou Evidence candidate.

## 2. Problème utilisateur

Une ligne de résultat ne suffit pas à comprendre un Event. L’analyste doit distinguer donnée source, rendu, enrichissement et interprétation sans promouvoir automatiquement l’Event en Evidence.

## 3. Objectifs

- rendre raw/rendered, parser, valeurs absentes et enrichissements inspectables
- conserver provenance et transformations
- permettre copie, occurrences et pivots reproductibles
- qualifier explicitement toute utilisation comme Artifact ou Evidence candidate

## 4. Non-objectifs

- ne pas modifier l’Event source
- ne pas présenter un enrichissement comme fait source
- ne pas créer automatiquement Evidence ou Finding
- ne pas résoudre OPEN-014

## 5. Propriétaire

Investigate / Signals And Hunt / Investigate Product Lead. Investigate possède uniquement les objets et mutations explicitement listés en section 10 ; les objets consommés restent chez leur propriétaire canonique.

## 6. Utilisateurs

Principal : SOC Analyst L2 et Threat Hunter. Secondaires : Forensic Analyst, Detection Engineer et Case reviewer.

## 7. Conditions d’entrée

Event résolvable depuis une recherche, un Signal ou un Case ; permissions raw/rendered ; source, parser ou lacune identifiables.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Event reference | Event Search, Signal ou Case | identifiant et contexte | oui | résolution à l’ouverture | tombstone ou erreur sans perdre la source |
| Raw/rendered payload | Telemetry source / parser | donnée et représentation | non selon permission | timestamp/version parser | afficher la variante autorisée et la lacune |
| Enrichments and relations | Shared enrichment / Entity resolution | données dérivées | non | fraîcheur propre | ne pas les inférer |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Telemetry Event | Shared Capabilities | raw, normalized, timestamps et source | inspecter et copier selon permission |
| Parser / Data Source | Platform Settings | version, mapping et health | consulter en projection |
| Entity | Shared Capabilities | relations et confiance | consulter et pivoter |
| Case / Artifact / Evidence | Investigate | liens existants et qualification | consulter et préparer une relation |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Case annotation/link | lier l’Event avec commentaire ou rôle | Investigate | classe 2 ; l’Event reste Shared |
| Artifact proposal | proposer un enregistrement ou une référence analysable | Investigate | qualification explicite ; aucun stockage défini |
| Evidence candidate | préparer une candidature | Investigate | pas une Evidence tant que CAP-INV-107 n’est pas exécutée |

## 11. Fonctionnalités

- raw et rendered côte à côte ou par mode
- source, parser et timestamps
- champs présents/manquants
- enrichissements distingués
- occurrences et pivots
- copie avec classification
- liaison au Case
- proposition Artifact/Evidence candidate

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter | analyste | Telemetry Event | 0 | permission | détails sourcés | non |
| Copier une valeur | analyste | champ | 0 | classification autorisée | valeur avec contexte | non |
| Pivoter | analyste | Event/field | 0 | champ autorisé | Query draft | non |
| Lier au Case | analyste | Event relation | 2 | Case accessible | relation et annotation | OPEN-013 |
| Proposer comme candidate | analyste | Artifact/Evidence candidate | 2 | source et raison | brouillon qualifiable | OPEN-013 |

Une action dont l’effet cible relève des classes 3 ou 4 reste une préparation ou une demande dans Investigate. L’autorité et l’exécution demeurent chez Govern ou le mécanisme propriétaire.

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Rendre et décoder | oui | parser versionné | oui | explication facultative | raw, mappings et diagnostics |
| Résoudre les relations | oui | entity matching sourcé | oui | suggestions | recherche manuelle par valeur |
| Proposer des pivots | oui | mapping de champs | oui | suggestion | sélection manuelle |
| Proposer une candidate Evidence | oui | checklist de provenance | oui | proposition seulement | qualification humaine via CAP-INV-107 |

Toute proposition automatisée expose initiateur, producteur, version ou run, sources, facteurs, éventuels Tool Calls, incertitude et disposition. Elle ne devient pas silencieusement un état effectif.

## 14. États fonctionnels

`available`, `partial`, `raw-restricted`, `parser-unknown`, `enrichment-stale`, `linked`, `candidate-proposed`, `source-unavailable`.

Ces états décrivent le travail de la capability ; ils ne finalisent pas la machine d’état canonique des objets, reportée à la phase Objets.

## 15. États d’interface

Raw access peut être refusé sans masquer le rendered autorisé ; Partial sépare champs source et enrichis ; Stale est par enrichissement ; le retour restaure ligne, query et scroll. Le Design System conserve la propriété du rendu et des interactions génériques.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Inspected Event context | projection détaillée | analyste et Case | source, parser, timestamps et permissions visibles |
| Pivot draft | Query context | Event Search | champ, valeur, période et origine conservés |
| Artifact/Evidence candidate | proposal | Artifact/Evidence workflows | non qualifiée, attribuée et rejetable |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Event Inspector | pivot sur valeur | Event Search | Event, champ, valeur, période et return origin | même Event restauré |
| Event Inspector | lier au Case | Case Workspace | Event ref, query/run, annotation et provenance | Inspector restauré |
| Event Inspector | qualifier une candidate | Evidence Creation | source Event, contexte, raison et transformations | retour à l’Event |

Chaque transition réévalue les permissions, conserve le tenant et l’environnement et ne transfère jamais l’ownership du produit destination.

## 18. Dépendances

- CAP-INV-002
- CAP-INV-007
- CAP-INV-104
- CAP-INV-105
- CAP-INV-107
- Telemetry Normalization
- Entity Resolution

## 19. Source de vérité

Telemetry Event et parser restent Shared/Settings ; Case, Artifact, Evidence et annotations sont Investigate ; enrichissements conservent leur producteur.

## 20. Provenance et audit

Event ID, source, raw/rendered distinction, parser/version, timestamps, enrichments/version, query/run, acteur, copied field, candidate reason et Case links.

## 21. Permissions fonctionnelles

Besoins : event read, raw access distinct, sensitive field reveal/copy, Case link, Artifact proposal et Evidence create. Les namespaces atomiques, le step-up et la séparation des tâches définitive restent à la phase Permissions.

## 22. Limites et erreurs

Event supprimé/inaccessible, payload tronqué, parser absent, champ classifié, enrichissement stale ou conflit d’identité restent visibles ; aucune causalité ni Evidence n’est inventée.

## 23. Métriques

- Events inspectés avec source/parser visibles
- pivots restaurant le contexte
- candidates qualifiées/rejetées
- accès raw refusés sans fuite

Aucune cible chiffrée définitive n’est fixée en Phase 4B.1.

## 24. Classification de livraison

`defined` / `planned`; aucun parser, moteur d’enrichissement ou format de stockage n’est choisi.

## 25. Critères d’acceptation

**Given** un Event issu d’une recherche  
**When** l’analyste l’inspecte  
**Then** raw/rendered, source, parser, champs, enrichissements et provenance sont distingués

**Given** l’analyste propose l’Event comme Evidence candidate  
**When** le Case est ouvert  
**Then** aucune Evidence n’existe avant qualification et la raison reste visible

**Given** raw est interdit mais rendered autorisé  
**When** l’Event est ouvert  
**Then** le rendu reste disponible et le raw n’est ni affiché ni inféré

## 26. Questions ouvertes

Quand un Event doit-il être enregistré comme Artifact versus référencé directement ? `OPEN-014` reste ouverte et la phase Objets doit préciser les relations.

## 27. Consommateurs documentaires

Event Search, Signal Triage, Hunt, Case Workspace, Evidence workflows et futurs écrans Event Inspector.
