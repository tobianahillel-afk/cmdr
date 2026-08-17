---
id: CAP-INV-366
title: Partition, Volume and Filesystem Identification
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
  - REQ-PROD-052
  - REQ-PROD-055
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-366 — Partition, Volume and Filesystem Identification

## 1. Définition
Identifier et comparer des partitions, volumes et filesystems candidats en exposant source, limites, confiance et contradictions, sans liste finale de support ni sélection silencieuse.

## 2. Problème utilisateur
Plusieurs interprétations peuvent coexister sur une image partielle, imbriquée ou incompatible. Une sélection implicite fausse toute navigation aval.

## 3. Objectifs
Afficher partitions/volumes/filesystems, limites/tailles/relations, espaces non attribués, structures imbriquées et volumes logiques candidats; sélectionner, modifier, comparer, marquer ambigu, tester une alternative et conserver l’historique.

## 4. Non-objectifs
Ne pas définir structures internes, offsets, signatures, support final, moteur, format, montage réel, protocole, API ou modification de la source.

## 5. Propriétaire
Investigate possède les candidats et la sélection analytique; Studio les Tools/Calls; Settings le support administré; Shared versioning/trace.

## 6. Utilisateurs
Principal : Disk Forensics Analyst. Secondaires : Artifact Analyst, Investigation Lead et Reviewer.

## 7. Conditions d’entrée
Disk Image et Integrity Review accessibles; session liée; Tool compatible explicitement choisi; limitations et permissions visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Disk Image/Review | Investigate | source/limites | oui | versions liées | blocked |
| Structure candidates | Tool Call | partitions/volumes/filesystems | oui | Tool/version visibles | unsupported |
| Platform/acquisition hints | Endpoint/acquisition | contexte déclaré | non | historique | ambiguous |
| Support/compatibility | Studio/Settings | Tools et support | oui | courant | tool-unavailable |
| Selection history | Investigate | choix/contestations | non | version courante | proposed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image | Investigate | taille/limites | lire |
| Tool/Tool Call | Studio | candidats/confiance/version | lire |
| Endpoint/acquisition | Endpoint Agent/Investigate | contexte plateforme | projection |
| Disk Session | Investigate concept | scope/choix précédent | lire |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Partition Candidate | créer/annoter/contester | Investigate concept | jamais certitude automatique |
| Volume Candidate | créer/annoter/contester | Investigate concept | relation partition explicite |
| Filesystem Candidate | proposer/select/confirm/supersede | Investigate concept | source/confiance/contradictions conservées |
| Selection history | ajouter | Investigate/Shared | append/versioned |

## 11. Fonctionnalités
Voir candidats, tailles, boundaries, relations, unallocated et nested structures; voir sources/confiance/contradictions; sélectionner/changer/comparer; marquer ambiguous/conflicting/unsupported; relancer profil alternatif; conserver historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter candidats | Analyst | candidates | 0 | résultat | vues sourcées | non |
| Comparer interprétations | Analyst | comparison | 1 | deux candidats | écarts visibles | non |
| Sélectionner/confirmer | Analyst | filesystem candidate | 2 | permission/revue | choix versionné | OPEN-013 |
| Tester alternative | Analyst | Tool Call request | 1 | Tool autorisé | résultat lié | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| proposer candidats | oui | profils/Tools | oui | suggestion | identification déterministe |
| expliquer contradiction | oui | comparaison | oui | oui | tableau pour/contre |
| sélectionner | humain | oui | workflow | jamais invisible | contrôle manuel |
| historiser | oui | oui | oui | non nécessaire | versioning |

## 14. États fonctionnels
`proposed`, `selected`, `confirmed-by-analyst`, `ambiguous`, `conflicting`, `unsupported`, `superseded`.

## 15. États d’interface
Loading conserve sélection; Empty n’invente aucune structure; Partial montre ranges non interprétés; Error garde candidats valides; Offline lecture; Permission denied masque données; Stale montre résultats superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Candidate set | observations | Session/Navigation | relations/confiance/contradictions visibles |
| Selected filesystem relation | relation | CAP-INV-367..377 | choix attribué, non absolu |
| Unsupported/ambiguous event | business event | Case/Trace | aucun support inventé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Disk Session | identifier | CAP-INV-366 | image, review, Tools, restrictions | Session |
| Filesystem Candidate | sélectionner | CAP-INV-367 | filesystem, volume, confidence, limitations | Identification |
| Ambiguous set | tester alternative | Studio Tool Call | image, candidate, functional profile | Identification |

## 18. Dépendances
CAP-INV-363/364/365/367, Studio Tools, Settings support, OPEN-005/008/013/015.

## 19. Source de vérité
Investigate possède candidats/sélection; Tool results et support restent chez Studio/Settings; image reste source immuable.

## 20. Provenance et audit
Image/session, Tool/version/Call/Run, candidate source, confidence, contradictions, selected profile, analyst, changes, errors et disposition.

## 21. Permissions fonctionnelles
Partition/volume/filesystem read, compare, select, confirm et automated identification request. Final RBAC/ABAC, step-up et SoD reportés.

## 22. Limites et erreurs
Aucune liste finale de filesystems. Image partial, structures overlapping/nested, Tool incompatible, support unknown, permission denied ou contradictions rendent le résultat ambiguous/unsupported sans inventer de structure.

## 23. Métriques
Candidats ambiguous/conflicting/unsupported, alternative profiles, selection changes et confirmations humaines.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation, moteur ou format choisi.

## 25. Critères d’acceptation
**Given** plusieurs filesystems candidats contradictoires **When** un candidat est sélectionné **Then** tous restent visibles, les preuves pour/contre et l’auteur du choix sont conservés, et une alternative peut être testée.

**Given** une image partielle avec zone non interprétée **When** l’identification termine **Then** la zone reste non attribuée et aucune partition fictive n’est créée.

**Given** aucun modèle IA **When** l’identification est réalisée **Then** profils, Tools, comparateur et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; support, profils, objets et permissions finales sont futurs.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-367..379, Workbench, Settings/Studio et phases Objects/Permissions/Technique.
