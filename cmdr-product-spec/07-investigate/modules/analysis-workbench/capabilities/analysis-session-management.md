---
id: CAP-INV-302
title: Analysis Session Management
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-002
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-302 — Analysis Session Management

## 1. Définition
Gérer le contexte fonctionnel durable d’une analyse statique, ses objectifs, participants, Artifacts, Tools, paramètres, résultats et dispositions, sans la confondre avec Case, Tool Call ou Automation Run.

## 2. Problème utilisateur
Les analyses longues utilisent plusieurs Tools, dérivés et reprises. Sans session dédiée, les paramètres, erreurs et décisions humaines se dispersent et l’analyse devient non reproductible.

## 3. Objectifs
- Créer, reprendre, suspendre, clôturer et rouvrir une session selon permission.
- Lier un ou plusieurs Artifacts au même objectif analytique.
- Conserver Tools, Tool Calls, paramètres, outputs et Derived Artifacts.
- Comparer des sessions et transmettre leurs résultats.

## 4. Non-objectifs
- Ne pas définir une machine d’état objet finale.
- Ne pas posséder Automation Run ou Tool Call.
- Ne pas remplacer le Case ou créer un workflow de réponse.

## 5. Propriétaire
Investigate possède le contexte et l’interprétation; Studio, Settings, Govern et Shared conservent leurs objets.

## 6. Utilisateurs
- Malware Analyst
- Case Analyst
- Investigation Lead
- Reviewer

## 7. Conditions d’entrée
- Case actif.
- Au moins un Artifact analysable.
- Objectif et owner définis.
- Permissions de session et d’Artifact disponibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case | Investigate | contexte durable | oui | état courant | session non liée interdite |
| Artifact(s) | Investigate | sources d’analyse | oui | versions figées/référencées | session incomplete |
| Objectif et scope | analyste | intention analytique | oui | version courante | rester draft |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope et participants | consulter/lier |
| Artifact | Investigate | versions et relations | consulter/lier |
| Tool/Tool Call | Studio | version, paramètres, statut, output | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | créer/modifier/clôturer/reouvrir/supersede | schéma final reporté |
| Session relation | Investigate | lier Case/Artifacts/Tools/results | sourcée et permission-aware |
| Session disposition | Investigate | enregistrer | completed n’implique pas Finding confirmé |

## 11. Fonctionnalités
- Définir objectif, owner, contributeurs et scope.
- Enregistrer Artifacts sources et dérivés, Tools, versions, Tool Calls, paramètres et outputs.
- Suspendre et reprendre sans perdre sélection ni contexte.
- Comparer sessions et conserver erreurs, interruptions et décisions humaines.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Créer session | Analyst | Analysis Session | 2 | Case/Artifact/objective | draft créé | OPEN-013 |
| Modifier objectif/participants | Owner | Analysis Session | 2 | session modifiable | nouvelle version | OPEN-013 |
| Suspendre/reprendre | Contributor | Analysis Session | 2 | permission et état compatible | état fonctionnel changé | OPEN-013 |
| Clôturer/réouvrir | Owner/Lead | Analysis Session | 2 | disposition/justification | historique conservé | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Résumé de session | oui | agrégation | oui | résumé attribué | timeline/filtre |
| Suggestion de prochaine analyse | oui | règles/checklists | oui | proposition | catalogue manuel |
| Capture des Tool Calls | non | oui | oui | non | Trace déterministe |
| Clôture de session | humain | contrôles de complétude | workflow de revue | jamais autonome | checklist humaine |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `cancelled`, `archived`, `superseded`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Session summary | Analysis Session projection | Case/Reviewer | scope, statuts et lacunes visibles |
| Session results | Analysis Result concept | CAP-INV-313 | sources et dispositions conservées |
| Session comparison | comparison result | Reviewer | préconditions et versions explicites |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-301 | créer/reprendre | CAP-INV-302 | Case, Artifact, objectif, owner, restrictions | intake |
| CAP-INV-302 | sélectionner Tool | CAP-INV-303 | session, Artifact actif, paramètres, permissions | session |
| CAP-INV-302 | transmettre résultats | CAP-INV-313 | observations, outputs, dérivés, sources, contradictions | session |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-301
- CAP-INV-303
- CAP-INV-311/312/313
- Studio Tool Call/Automation Run
- Shared Collaboration/Trace/Versioning
- OPEN-013/015

## 19. Source de vérité
Investigate possède le contexte et la disposition de la session ; Studio possède chaque Tool Call et Automation Run ; Shared possède les mécanismes de trace/versioning.

## 20. Provenance et audit
Enregistrer création, versions, owner, contributeurs, objectifs, Artifacts, Tools, paramètres, outputs, erreurs, interruptions, dispositions et return origin.

## 21. Permissions fonctionnelles
- Analysis Session create/update/close/reopen
- Artifact read
- Tool output read
- collaboration
- reproducibility review
- cross-tenant denied

Matrice atomique reportée.

## 22. Limites et erreurs
- Session sans Artifact courant.
- Contributeur retiré ou permission révoquée.
- Tool Call orphelin ou output expiré.
- Réouverture après Artifact superseded.

## 23. Métriques
- Sessions actives/paused/partial.
- Durée et reprises.
- Sessions avec provenance complète.
- Comparaisons et handoffs réalisés.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** une session active avec plusieurs Tool Calls
**When** l’analyste la suspend puis la reprend
**Then** objectif, Artifacts, sélection, paramètres et outputs restent liés

**Given** un utilisateur sans droit de réouverture
**When** il tente de rouvrir une session clôturée
**Then** l’action est refusée et l’historique reste inchangé

**Given** aucun modèle IA
**When** une session est gérée de bout en bout
**Then** formulaires, états, checklists et résultats déterministes restent disponibles

## 26. Questions ouvertes
- Analysis Session n’est pas encore un objet canonique.
- OPEN-013 régit les changements de classe 2.
- OPEN-015 conserve la distinction Automation Run/Response Run.

## 27. Consommateurs documentaires
- Analysis Workbench
- Case Replay et Timeline
- Static/Dynamic future workbenches
- phases Objets, Permissions, Journeys
