---
id: CAP-INV-367
title: Filesystem Navigation and Metadata Inspection
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-UX-002
  - REQ-UX-006
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-367 — Filesystem Navigation and Metadata Inspection

## 1. Définition
Naviguer en lecture seule dans une projection de filesystem et inspecter entrées, liens, chemins, tailles, types, identités, permissions, timestamps, attributs et états sans confondre chemin, identité ou événement certain.

## 2. Problème utilisateur
Une arborescence partielle ou ambiguë peut cacher des entrées supprimées, inaccessibles ou incohérentes; un chemin ou timestamp peut être interprété à tort comme identité stable ou action certaine.

## 3. Objectifs
Parcourir répertoires/entrées/liens; rechercher, filtrer, trier et sélectionner; afficher metadata et états; annoter/bookmark; préserver navigation, sélection, historique et return origin.

## 4. Non-objectifs
Ne pas monter ou modifier un live filesystem, exécuter un fichier, définir colonnes/filtres finaux, structures bas niveau, moteur, format, API, commande ou objet File complet.

## 5. Propriétaire
Investigate possède la projection, sélection, bookmarks et annotations. Disk Image reste source; Studio possède Tools; Shared possède Search/Tree/Inspector/History.

## 6. Utilisateurs
Principal : Disk Forensics Analyst. Secondaires : Artifact Analyst, Investigation Lead, Privacy Reviewer et Evidence Reviewer.

## 7. Conditions d’entrée
Session active; filesystem sélectionné avec confiance/limites visibles; Disk Image lisible; permissions metadata/content distinctes; source immuable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Selected filesystem | CAP-INV-366 | structure/confiance | oui | version sélectionnée | filesystem-required |
| Entry/metadata projection | Tool Call | tree, links, metadata, errors | oui | Tool/version visibles | partial/failed |
| Disk limitations | CAP-INV-365 | missing ranges/restrictions | oui | inherited | partial |
| Search/filter/sort context | analyst/Shared | navigation | non | session courante | default view |
| Permissions | Security | metadata/content scope | oui | courant | permission-denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/Filesystem Candidate | Investigate | source/selection | lire |
| Filesystem Entry | Investigate concept | path, metadata, state | lire |
| Tool/Tool Call | Studio | producer/version/errors | lire |
| Case/Hypothesis | Investigate | contexte | lire/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Bookmark/selection | créer/update/remove | Investigate concept | source/version requises |
| Navigation history | ajouter/restaurer | Shared/Investigate | return context conservé |
| Entry annotation | créer/update/supersede | Investigate | path ≠ stable identity |
| Source entry | aucune mutation | Disk Image | read-only |

## 11. Fonctionnalités
Naviguer, voir entrées/liens/paths/sizes/types/owners/permissions/timestamps/attributes et états active/deleted-candidate/inaccessible/partial; rechercher/filtrer/trier; annoter/bookmark/multi-select; préserver history/back et contexte.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Naviguer/inspecter | Analyst | entry/tree | 0 | read | projection visible | non |
| Rechercher/filtrer/trier | Analyst | view | 0 | scope autorisé | vue persistée | non |
| Annoter/bookmark | Analyst | local knowledge | 2 | session modifiable | version attribuée | OPEN-013 |
| Produire vue/comparison | Analyst | analysis result | 1 | source disponible | résultat lié | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| construire arbre | oui | Tool/parser | oui | non nécessaire | tree/table |
| rechercher/filter | oui | oui | oui | suggestion | query explicite |
| proposer relations | oui | règles | oui | oui | links/metadata |
| annoter | oui | oui | oui | aide | édition manuelle |

Aucun filesystem n’est choisi invisiblement; attribution complète.

## 14. États fonctionnels
`loading`, `available`, `partial`, `inaccessible`, `deleted-candidate`, `inconsistent`, `restricted`, `failed`, `superseded`.

## 15. États d’interface
Loading conserve tree/selection; Empty explique; Partial marque branches manquantes; Error garde données valides; Offline lecture; Permission denied ne révèle paths/content; Stale montre source/version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Filesystem projection | view | analyst/Inspector | source et partialité visibles |
| Entry selection/bookmark | relation | CAP-INV-368..375 | path, metadata, source version liés |
| Navigation/annotation event | Trace/Activity | Session/Replay | actor et disposition visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Filesystem Candidate | ouvrir | CAP-INV-367 | image, volume, filesystem, confidence, limits | Identification |
| Entry | inspecter contenu/identité | CAP-INV-368 | entry, metadata, source, permissions | Navigation |
| Deleted candidate | analyser | CAP-INV-369 | entry, residual metadata, source | Navigation |

## 18. Dépendances
CAP-INV-365/366/368/369/370, Shared Tree/Search/Inspector/History, Studio Tool Calls, OPEN-005/008/013/014.

## 19. Source de vérité
Disk Image et Tool result restent sources; Investigate possède projection/interprétation/annotations; path n’est jamais l’identité absolue.

## 20. Provenance et audit
Image/session/filesystem, Tool/version, entry identifiers as displayed, path, metadata, search/filter/sort, selection, annotation, permission, errors et correlation IDs.

## 21. Permissions fonctionnelles
Filesystem browse, metadata read, sensitive path read, search, bookmark, annotation et cross-tenant restrictions. Content preview est séparé et appartient à CAP-INV-368/376.

## 22. Limites et erreurs
Entry missing/inaccessible, broken link, partial tree, metadata conflict, unsupported attribute, permission denied, stale/superseded source. Aucun path/timestamp/contenu n’est inventé.

## 23. Métriques
Entries/branches partial/inaccessible, navigation restore, searches permission-aware, bookmarks et annotations avec provenance.

## 24. Classification de livraison
`defined` / `planned`; aucune UI finale, implémentation, moteur ou format choisi.

## 25. Critères d’acceptation
**Given** une entrée supprimée avec chemin incomplet **When** elle est ouverte **Then** le path non confirmé n’est pas inventé et l’état deleted-candidate reste visible.

**Given** un timestamp metadata sans corroboration **When** il est affiché **Then** sa source/type sont visibles et aucun événement certain n’est affirmé.

**Given** aucun modèle IA **When** l’analyste navigue **Then** tree, table, search, filters, history et annotations fonctionnent.

## 26. Questions ouvertes
OPEN-005/008/013/014 restent ouvertes; identité, schemas et permissions finales sont reportés.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-368..379, Artifact Explorer/Inspector/Search, Case/Evidence et phases Objects/Permissions/Screens.
