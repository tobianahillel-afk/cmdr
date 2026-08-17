---
id: CAP-INV-350
title: Platform and Analysis Profile Identification
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
# CAP-INV-350 — Platform and Analysis Profile Identification

## 1. Définition
Permet d’identifier et sélectionner une plateforme ou un profil d’analyse candidat en exposant sources, confiance, contradictions, compatibilités et limites, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans identification attribuée de plateforme et de profil, le produit risque d’appliquer une interprétation incompatible ou de masquer des contradictions entre candidats.

## 3. Objectifs
- identifier et sélectionner une plateforme ou un profil d’analyse candidat en exposant sources, confiance, contradictions, compatibilités et limites.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Memory Forensics Analyst**. Secondaires : Investigation Lead, Evidence Reviewer, Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles; Tool/profil sélectionné explicitement. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image et acquisition context | CAP-INV-349 | source, déclarations et limitations | oui | assessment courant | profile-required |
| Candidate profiles | Tools déterministes / analyste | candidats et preuves | oui pour analyse | Tool/version visibles | unsupported |
| Tool compatibility | CMDR Studio | Tools et versions disponibles | oui pour exécution | courant | tool-unavailable |
| Platform support projection | Endpoint Agent / Settings | support déclaré | non pour consultation | courant | warning OPEN-008 |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image | Investigate | métadonnées et contenu autorisé | lecture |
| Endpoint/platform projection | Endpoint Agent / Settings | plateforme déclarée et support | lecture |
| Tool / Tool Call | CMDR Studio | compatibilité, version et erreurs | lecture |
| Analysis Profile candidate | Investigate concept | source, confiance et contradictions | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Platform/Profile Identification | créer/versionner | Investigate concept | candidat distinct de confirmé |
| Selected profile relation | sélectionner/modifier/superseder | Investigate | auteur et justification |
| Alternative run request | préparer | Investigate / Studio | Tool Call distinct |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- voir plateforme déclarée, détectée et candidates.
- voir profils candidats, sources, confiance et contradictions.
- voir Tools compatibles et limitations.
- sélectionner, modifier ou marquer ambigu.
- tester une alternative et conserver l’historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Comparer candidats | Memory Forensics Analyst | Profile set | 0 | lecture autorisée | vue sourcée | non |
| Sélectionner profil | Memory Forensics Analyst | Analysis Profile | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Confirmer par analyste | Memory Forensics Analyst | Analysis Profile | 2 | permission, justification, version modifiable | nouvelle version; précédente conservée | OPEN-013; Govern si cible réelle |
| Tester une alternative | Memory Forensics Analyst | Tool Call request | 1 | scope/Tool/policy explicites | job/extraction borné et tracé | non |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| voir plateforme déclarée, détectée et candidates | oui | oui | oui | proposition | détecteurs déterministes et preuves brutes |
| voir profils candidats, sources, confiance et contradictions | oui | oui | oui | proposition | table comparative des candidats |
| voir Tools compatibles et limitations | oui | oui | oui | proposition | sélection manuelle attribuée |
| Conclusion finale | oui | non | non | assistance | revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`proposed`, `selected`, `confirmed-by-analyst`, `ambiguous`, `conflicting`, `unsupported`, `superseded`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Profile candidates | Analysis Profile candidates | CAP-INV-348/351 | sources, confiance et contradictions visibles |
| Selected profile | relation versionnée | analyses mémoire | sélection attribuée, pas certitude |
| Unsupported/tool-unavailable status | Trace event | Case / Session | aucun résultat fictif |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Session | identifier profil | CAP-INV-350 | image, acquisition, Tools et résultats antérieurs | Session |
| CAP-INV-350 | profil selected | CAP-INV-351..360 | profil, confiance, limites et Tool version | Profile view |
| CAP-INV-350 | ambiguïté | Alternative comparison | candidats et preuves | Profile view |

Tenant, Case, image, permissions, restrictions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-347, CAP-INV-348, CMDR Studio Tool catalog, Platform Settings, OPEN-005, OPEN-008, OPEN-013, OPEN-015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline/Export/Versioning : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, session, profil, Tool/version, Calls/Run, paramètres, filtres, acteur, timestamps, erreurs, interruptions, données sensibles, Derived Artifacts et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Comparer candidats | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Sélectionner profil | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Confirmer par analyste | mutation réversible | 2 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Tester une alternative | traitement borné | 1 | selon policy | SoD si requise | Investigate/Security | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- aucune liste définitive de plateformes.
- OPEN-008 reste ouverte.
- profil proposé distinct de confirmé.
- aucun moteur ou plugin imposé.
- Memory Forensics ≠ Debugger Memory View ≠ Disk Forensics ≠ full Network Forensics.
- Une erreur ou incohérence ne devient pas une conclusion.

## 23. Métriques
- profils ambiguous/unsupported.
- sélections révisées.
- analyses alternatives.
- résultats dépendants d’un profil partial.
- complétude de provenance et retour au contexte sans perte.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin, plateforme ou intégration native.

## 25. Critères d’acceptation
### 1. Profil ambigu
**Given** plusieurs candidats se contredisent **When** l’analyste sélectionne **Then** tous restent visibles, sélection attribuée et alternative testable.
### 2. Unsupported
**Given** aucun profil compatible **When** l’analyse est préparée **Then** unsupported est visible et aucune analyse fictive n’apparaît.
### 3. Sans IA
**Given** aucun modèle **When** les candidats sont examinés **Then** Tools déterministes et sélection manuelle fonctionnent.

## 26. Questions ouvertes
OPEN-005, OPEN-008, OPEN-013, OPEN-015 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, Static/Reverse, futures phases Objects/Permissions/Journeys/Screens/Contracts, handoffs seulement vers 4B.2B.3B/4B.3.
