---
id: CAP-INV-364
title: Disk Forensics Session Management
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
  - REQ-OBJ-009
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-364 — Disk Forensics Session Management

## 1. Définition
Créer, reprendre, suspendre, clôturer, rouvrir, comparer et superseder un contexte durable Disk Forensics lié au Case et à une ou plusieurs Disk Images.

## 2. Problème utilisateur
Sans session durable, images, filesystems sélectionnés, Tools, requêtes, observations, erreurs et Derived Artifacts perdent leur contexte et leur attribution.

## 3. Objectifs
Conserver Case, images, objective/scope, owner/contributors, partitions/volumes/filesystems, Tools/Calls, paramètres, vues, requêtes/filtres, observations, Derived Artifacts, erreurs, partials, historique et handoffs.

## 4. Non-objectifs
Ne pas devenir Case, Collection Job, Automation Run ou objet canonique complet; ne pas définir schéma, moteur, stockage, API, protocole, commande ou full Network Forensics.

## 5. Propriétaire
Investigate possède le concept de session et son contenu analytique. Studio possède Tools/Runs; Settings l’administration; Shared collaboration/versioning/recovery; Govern les cibles réelles.

## 6. Utilisateurs
Principal : Disk Forensics Analyst. Secondaires : Investigation Lead, co-analystes, Evidence Reviewer et Audit Analyst.

## 7. Conditions d’entrée
Intake accessible; au moins une Disk Image liée; objective/scope définis; owner et permissions connus; limitations héritées conservées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Intake et Case | CAP-INV-363 | contexte | oui | version courante | draft |
| Disk Images | Investigate | sources/versions | oui | versions liées | blocked |
| Owner/contributors | identity/access | responsabilité | oui | courant | incomplete |
| Structures/Tools/parameters | CAP-INV-366 / Studio | contexte analytique | selon activité | versionné | partial |
| Observations/errors/Derived Artifacts | capabilities aval | travail durable | non | session courante | empty |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case/Hypothesis | Investigate | objective et contexte | lire/lier |
| Disk Image/Artifact | Investigate | sources et restrictions | lire/lier |
| Tool/Tool Call/Automation Run | Studio | versions/exécutions | lire |
| Background Job/Trace/Activity | Shared | progression/audit | lire |
| User/principal | Settings | membership projection | lire |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Disk Forensics Session | créer/update/pause/close/reopen/supersede | Investigate concept | schéma final reporté |
| Session membership | ajouter/retirer/versionner | Investigate concept | permission-aware |
| Session links | créer/supersede | Investigate | source et return origin conservés |

## 11. Fonctionnalités
Créer/reprendre; lier Case et plusieurs images; définir objective/scope/owner/contributors; conserver structures sélectionnées, Tools/Calls, paramètres, vues, requêtes/filtres, observations, Derived Artifacts, erreurs/partials; suspendre/reprendre/clôturer/rouvrir/comparer/superseder/transmettre.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter/reprendre | Analyst | Session | 0/2 | read/update | contexte restauré | OPEN-013 pour mutation |
| Créer/modifier | Analyst | Session | 2 | intake prêt | version nouvelle | OPEN-013 |
| Pause/close/reopen | Owner | Session | 2 | état compatible | état attribué | OPEN-013 |
| Comparer sessions | Analyst | Comparison | 1 | sessions lisibles | résultat lié | non |
| Superseder | Owner | Session | 2 | raison | historique conservé | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| restaurer contexte | oui | oui | oui | non nécessaire | versioning/recovery |
| résumer session | oui | oui | oui | oui | tables et activity |
| proposer prochaine vue | oui | règles/checklist | oui | oui | navigation manuelle |
| transmettre résultats | oui | workflow | oui | brouillon | sélection/validation humaine |

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `archived`, `superseded`.

## 15. États d’interface
Loading restaure la structure; Empty explique l’absence de résultat; Partial garde les outputs valides; Error conserve la session; Offline limite les mutations; Permission denied masque sources; Stale expose versions superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Session context | concept | CAP-INV-365..379 | sources, scope et limitations conservés |
| Session state event | Trace/Activity | Case/Collaboration | actor/reason/version visibles |
| Comparison/handoff context | package | CAP-INV-377/379 | sessions et provenance liées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Intake | créer/reprendre | CAP-INV-364 | Case, images, objective, owner, limits | Intake |
| Session | identifier filesystem | CAP-INV-366 | image, scope, prior results, Tools | Session |
| Session | transmettre | CAP-INV-379 | selected observations, contradictions, provenance | Session |

## 18. Dépendances
CAP-INV-363/366/377/378/379, Case/Artifact, Studio Tools/Runs, Shared Versioning/Collaboration/Recovery/Trace, OPEN-005/008/013/015.

## 19. Source de vérité
Investigate est source de la session conceptuelle et des relations; chaque objet consommé reste chez son owner.

## 20. Provenance et audit
Case, images/versions, owner/contributors, objective/scope, structures, Tools/Calls/Runs, paramètres, queries/filters, observations, errors, interruptions, states, handoffs et dispositions.

## 21. Permissions fonctionnelles
Session create/read/update/pause/close/reopen/archive/supersede, membership manage, comparison create et handoff prepare. Matrice atomique, step-up et SoD finaux reportés.

## 22. Limites et erreurs
Image inaccessible, contributor unauthorized, Tool missing, session conflict, partial output, stale selection, failed recovery ou cross-tenant mismatch. Aucune source n’est modifiée.

## 23. Métriques
Sessions créées/reprises, partial/blocked/failed, recovery success, supersessions et sessions avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; objet, états finaux, permissions et implémentation restent futurs.

## 25. Critères d’acceptation
**Given** une session active avec filtres et sélection **When** elle est reprise **Then** contexte, source, vues, erreurs et return origin sont restaurés.

**Given** une image devient inaccessible **When** la session est reprise **Then** elle est blocked/partial sans supprimer les résultats valides.

**Given** aucun modèle IA **When** la session est gérée **Then** versioning, recovery, tables et workflows déterministes suffisent.

## 26. Questions ouvertes
OPEN-005/008/013/015 restent ouvertes; schéma et permission model finaux sont reportés.

## 27. Consommateurs documentaires
INV-DSK-001, Case Workspace, CAP-INV-363..379, Collaboration/Trace/Recovery et phases Objects/Permissions/Journeys/Screens/Technique.
