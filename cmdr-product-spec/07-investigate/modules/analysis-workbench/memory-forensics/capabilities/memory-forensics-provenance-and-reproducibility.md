---
id: CAP-INV-361
title: Memory Forensics Provenance and Reproducibility
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
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-361 — Memory Forensics Provenance and Reproducibility

## 1. Définition
Permet de retracer acquisition, custody, image, session, profil, Tools, paramètres, observations, accès sensibles, erreurs et dispositions puis évaluer la reproductibilité sans redéfinir Trace ou Audit, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans provenance complète, une analyse ne peut pas être reproduite, contestée ou expliquée lorsque l’image, le profil, le Tool, la version ou les paramètres manquent.

## 3. Objectifs
- retracer acquisition, custody, image, session, profil, Tools, paramètres, observations, accès sensibles, erreurs et dispositions puis évaluer la reproductibilité sans redéfinir Trace ou Audit.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Audit Analyst**. Secondaires : Investigation Lead, Evidence Reviewer, Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles; Tool/profil sélectionné explicitement. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image et Session | Investigate | source, scope et restrictions | oui | versions liées | blocked ou partial |
| Selected platform/profile | CAP-INV-350 | interprétation et limites | oui pour analyse | sélection courante | unsupported ou profile-required |
| Memory Forensics Provenance and Reproducibility results | Tool Calls / analyst | projections sourcées | oui | Tool et version visibles | partial ou failed |
| Case, permissions et policies | Investigate / Security / Settings | contexte et autorisation | oui | courant | denied ou restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image / Artifact | Investigate | source, version, limites et lineage | lecture |
| Reproducibility Assessment | Investigate concept | projection, source et confiance | lecture |
| Tool / Tool Call / Automation Run | CMDR Studio | producteur, version, paramètres et statut | lecture |
| Case / Hypothesis | Investigate | contexte et relations | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Reproducibility Assessment | créer/annoter/contester/superseder | Investigate concept | aucun schéma final |
| Provenance relation | créer/lier/versionner | Investigate concept | source et incertitude obligatoires |
| Trace / Activity event | émettre | Shared | append-only |
| Source objects | aucune mutation | owners respectifs | projections seulement |

## 11. Fonctionnalités
- retracer Case, Endpoint, request, Job, image, custody et session.
- retracer profil, Tools, versions, Tool Calls, Automation Runs et paramètres.
- retracer observations, Derived Artifacts, accès sensibles, erreurs et décisions.
- évaluer la reproductibilité et ses préconditions manquantes.
- contester, superseder ou reproduire explicitement.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter et filtrer | Audit Analyst | Reproducibility Assessment | 0 | lecture autorisée | vue sourcée | non |
| Comparer plusieurs vues ou images | Audit Analyst | Comparison | 0 | lecture autorisée | vue sourcée | non |
| Annoter, contester ou relier | Audit Analyst | Reproducibility Assessment | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Lancer un traitement borné | Audit Analyst | Tool Call | 1 | scope/Tool/policy explicites | job/extraction borné et tracé | non |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| retracer Case, Endpoint, request, Job, image, custody et session | oui | oui | oui | proposition | assemblage déterministe des liens de provenance |
| retracer profil, Tools, versions, Tool Calls, Automation Runs et paramètres | oui | oui | oui | proposition | checklist de reproductibilité |
| retracer observations, Derived Artifacts, accès sensibles, erreurs et décisions | oui | oui | oui | proposition | comparaison Tool/version/paramètres |
| Conclusion finale | oui | non | non | assistance | revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`reproducible`, `partially-reproducible`, `not-reproducible`, `missing-image`, `incomplete-image`, `missing-profile`, `missing-tool`, `missing-version`, `unsupported-platform`, `policy-blocked`, `disputed`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Reproducibility Assessment | Memory Forensics Provenance and Reproducibility result | Memory Workbench | source, statut, partial et incertitude visibles |
| Provenance relation | relation ou candidate | Case / Hypothesis | aucune conclusion automatique |
| Handoff selection | candidate context | CAP-INV-362 | provenance et restrictions conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Memory Forensics Session | ouvrir Memory Forensics Provenance and Reproducibility | Current capability | image, profil, Tool, scope et return origin | Session |
| Current capability | naviguer vers source | Processes / Regions / Modules / Case | selection, relations et provenance | Current view |
| Current capability | préparer handoff | CAP-INV-362 | observations, contradictions, uncertainty et lineage | Current view |

Tenant, Case, image, permissions, restrictions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-213, CAP-INV-214, CAP-INV-312, CAP-INV-345, OPEN-005, OPEN-008, OPEN-013, OPEN-014, OPEN-015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline/Export/Versioning : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, session, profil, Tool/version, Calls/Run, paramètres, filtres, acteur, timestamps, erreurs, interruptions, données sensibles, Derived Artifacts et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Inspecter et filtrer | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Comparer plusieurs vues ou images | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Annoter, contester ou relier | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Lancer un traitement borné | traitement borné | 1 | selon policy | SoD si requise | Investigate/Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- ne duplique pas Trace, Activity ou Audit.
- aucune infrastructure technique définie.
- not-reproducible ne prouve pas invalidité.
- Memory Forensics ≠ Debugger Memory View ≠ Disk Forensics ≠ full Network Forensics.
- Une erreur ou incohérence ne devient pas une conclusion.

## 23. Métriques
- résultats Memory Forensics Provenance and Reproducibility par état.
- taux partial/failed/unsupported.
- observations avec provenance complète.
- candidates accepted/modified/rejected.
- complétude de provenance et retour au contexte sans perte.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin, plateforme ou intégration native.

## 25. Critères d’acceptation
### 1. Profil manquant
**Given** image et Tool existent mais le profil n’est plus disponible **When** le reviewer évalue **Then** missing-profile est visible et aucune fausse reproductibilité n’est déclarée.
### 2. Résultat partiel
**Given** des données ou structures sont manquantes **When** l’analyse termine **Then** partial est visible et aucune donnée n’est inventée.
### 3. Sans IA
**Given** aucun fournisseur de modèle **When** l’analyste travaille **Then** Tools déterministes, tables, filtres et revue humaine restent disponibles.

## 26. Questions ouvertes
OPEN-005, OPEN-008, OPEN-013, OPEN-014, OPEN-015 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, Static/Reverse, futures phases Objects/Permissions/Journeys/Screens/Contracts, handoffs seulement vers 4B.2B.3B/4B.3.
