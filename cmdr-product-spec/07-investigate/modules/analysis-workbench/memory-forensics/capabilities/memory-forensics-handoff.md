---
id: CAP-INV-362
title: Memory Forensics Handoff to Evidence, Findings and Detection Engineering
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
  - REQ-OBJ-004
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-362 — Memory Forensics Handoff to Evidence, Findings and Detection Engineering

## 1. Définition
Permet de préparer des packages attribués vers Evidence, Finding et future Detection Engineering sans qualification, confirmation ou déploiement automatique, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans handoff explicite, une observation mémoire, une anomalie, une valeur sensible ou une connexion candidate peut être promue directement en Evidence, Finding ou règle déployée.

## 3. Objectifs
- préparer des packages attribués vers Evidence, Finding et future Detection Engineering sans qualification, confirmation ou déploiement automatique.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Investigation Lead**. Secondaires : Investigation Lead, Evidence Reviewer, Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles; Tool/profil sélectionné explicitement. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Observations et contradictions sélectionnées | CAP-INV-351..361 | résultats sourcés | oui | versions de session courantes | package incomplete |
| Case, Hypothesis et Evidence existantes | Investigate | contexte de qualification | oui | versions courantes | handoff non soumis |
| Memory Image, Session et provenance | Investigate | source et lineage | oui | versions liées | blocked |
| Permissions, restrictions et données sensibles masquées | Security / Settings | gates de lecture et diffusion | oui | courant | restricted ou denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis / Evidence / Finding | Investigate | contexte, qualification existante et contradictions | lecture/lien |
| Memory Image / Derived Artifact | Investigate | source, lineage, restrictions et versions | lecture/lien |
| Memory observations and anomalies | Investigate concepts | selected results, confidence and disposition | lecture/sélection |
| Tool / Tool Call / Automation Run | CMDR Studio | producteur, version, paramètres et statut | lecture |
| Sensitive Material Candidate | Investigate / Security | référence masquée et accès autorisé | lecture masquée uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence candidate package | créer/modifier/soumettre/retirer | Investigate concept | reste distinct d’Evidence qualifiée |
| Finding Draft | créer/modifier/soumettre/retirer | Investigate | reste distinct d’un Finding confirmé |
| Future Detection Engineering knowledge package | préparer/versionner | Investigate / future owner | aucune règle créée ou déployée |
| Trace / Activity event | émettre | Shared | append-only; source et disposition conservées |
| Source objects | aucune mutation | owners respectifs | références versionnées seulement |

## 11. Fonctionnalités
- sélectionner observations, processus, régions, modules, connexions, anomalies et Artifacts.
- sélectionner données sensibles masquées et contradictions.
- expliquer la relation avec Hypothesis.
- préparer Evidence candidate ou Finding Draft.
- préparer futurs handoffs Detection Engineering et Network Forensics sans créer de règle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter et sélectionner | Investigation Lead | observations et sources | 0 | lecture autorisée | sélection sourcée | non |
| Préparer Evidence candidate | Evidence Reviewer | candidate package | 2 | provenance et qualification requise visibles | package versionné | OPEN-013 |
| Préparer Finding Draft | Investigation Lead | Finding Draft | 2 | Hypothesis, Evidence existantes et contradictions visibles | draft attribué | OPEN-013 |
| Préparer futur handoff Detection Engineering | Analyst autorisé | knowledge package | 2 | limites, conditions et sources visibles | package uniquement | non; aucune règle |
| Retirer ou superseder un package | owner autorisé | handoff package | 2 | justification | ancienne version conservée | OPEN-013 |

Classes 3/4 bloquées et aucune qualification, confirmation ou règle n’est exécutée localement.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Sélectionner observations, contradictions et Artifacts | oui | oui | oui | proposition | sélection manuelle et filtres |
| Préparer Evidence candidate | oui | oui | oui | brouillon | checklist et package déterministe |
| Préparer Finding Draft | oui | oui | oui | brouillon | formulaire sourcé et revue humaine |
| Préparer future connaissance Detection Engineering | oui | oui | oui | brouillon | matrice comportement/conditions/limites |
| Qualifier Evidence, confirmer Finding ou déployer une règle | humain propriétaire uniquement | non | non | interdit | capabilities propriétaires et revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, statut, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `ready-for-review`, `submitted`, `returned`, `accepted-as-candidate`, `superseded`, `withdrawn`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evidence candidate package | candidate package | CAP-INV-107/108 | observation distincte d’Evidence; qualification requise |
| Finding Draft | draft | CAP-INV-109 | aucune confirmation automatique |
| Future Detection Engineering package | knowledge package | future Phase 4B.3 | aucune règle créée ou déployée |
| Future Network Forensics package | candidate context | future Phase 4B.2B.3B | uniquement traces réseau présentes en mémoire |
| Handoff disposition event | Trace / Activity event | Case Replay / Audit | auteur, sources, décision et retour conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Memory observations / Session | préparer handoff | CAP-INV-362 | image, session, selections, contradictions, uncertainty, provenance | source Workbench |
| CAP-INV-362 | soumettre Evidence candidate | CAP-INV-107/108 | candidate package, sources, restrictions and qualification required | Memory Workbench |
| CAP-INV-362 | soumettre Finding Draft | CAP-INV-109 | Hypothesis, Evidence existantes, observations, contradictions and author | Memory Workbench |
| CAP-INV-362 | préparer connaissance future | Phase 4B.3 | behavior, processes, regions, modules, conditions, limits and sources | Memory Workbench |
| CAP-INV-362 | préparer network handoff futur | Phase 4B.2B.3B | connection candidates, process links, timestamps, limits and provenance | Memory Workbench |

Tenant, Case, image, permissions, restrictions, sélection et return origin sont préservés.

## 18. Dépendances
CAP-INV-107, CAP-INV-108, CAP-INV-109, Future Phase 4B.3, OPEN-008, OPEN-013, OPEN-014, OPEN-015. Aucune dépendance bas niveau.

## 19. Source de vérité
Image/contexte : Investigate; acquisition : Collection/Endpoint Agent; administration : Settings; Tools/Runs : Studio; Trace/Timeline/Export/Versioning : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, session, profil, Tool/version, Calls/Run, paramètres, filtres, acteur, timestamps, erreurs, interruptions, données sensibles, Derived Artifacts et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Sélectionner observations | lecture sensible | 0 | selon policy | SoD si requise | Investigate/Security | Permissions |
| Evidence candidate prepare/submit | impact de qualification | 2 | possible | preparer distinct du reviewer | Investigate | Permissions |
| Finding Draft prepare/submit | impact analytique | 2 | possible | author distinct du confirmer | Investigate | Permissions |
| Detection/Network handoff prepare | diffusion future | 2 | selon policy | owner futur valide | Investigate/future owner | Permissions |
| Sensitive reference include | donnée sensible | 0/2 | masking et step-up | reviewer Security si requis | Security/Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC, step-up et SoD finaux reportés.

## 22. Limites et erreurs
- observation non Evidence.
- candidate non Finding, IOC ou credential.
- aucune règle créée ou déployée.
- Memory Forensics ≠ Debugger Memory View ≠ Disk Forensics ≠ full Network Forensics.
- Une erreur ou incohérence ne devient pas une conclusion.

## 23. Métriques
- résultats Memory Forensics Handoff to Evidence, Findings and Detection Engineering par état.
- taux partial/failed/unsupported.
- observations avec provenance complète.
- candidates accepted/modified/rejected.
- complétude de provenance et retour au contexte sans perte.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin, plateforme ou intégration native.

## 25. Critères d’acceptation
### 1. Observation non qualifiée
**Given** un résultat mémoire candidat est sélectionné **When** l’analyste prépare Evidence **Then** le résultat reste distinct, la provenance est liée et une qualification reste nécessaire.
### 2. Résultat partiel
**Given** des données ou structures sont manquantes **When** l’analyse termine **Then** partial est visible et aucune donnée n’est inventée.
### 3. Sans IA
**Given** aucun fournisseur de modèle **When** l’analyste travaille **Then** Tools déterministes, tables, filtres et revue humaine restent disponibles.

## 26. Questions ouvertes
OPEN-008, OPEN-013, OPEN-014, OPEN-015 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, Static/Reverse, futures phases Objects/Permissions/Journeys/Screens/Contracts, handoffs seulement vers 4B.2B.3B/4B.3.
