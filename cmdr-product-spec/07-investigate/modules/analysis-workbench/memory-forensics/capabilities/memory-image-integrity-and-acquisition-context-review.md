---
id: CAP-INV-349
title: Memory Image Integrity and Acquisition Context Review
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
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-SEC-001
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-349 — Memory Image Integrity and Acquisition Context Review

## 1. Définition
Permet d’examiner provenance, acquisition, transformations, transferts, vérifications, erreurs et custody afin de qualifier l’image pour l’analyse sans définir l’algorithme d’intégrité, sans moteur, plugin ni implémentation imposés.

## 2. Problème utilisateur
Sans revue séparée de l’intégrité et de l’acquisition, une vérification technique peut être confondue avec la pertinence analytique et une image incomplète ou transformée peut produire des conclusions trompeuses.

## 3. Objectifs
- examiner provenance, acquisition, transformations, transferts, vérifications, erreurs et custody afin de qualifier l’image pour l’analyse sans définir l’algorithme d’intégrité.
- Préserver Case, image, limites, incertitude, provenance et return origin.
- Séparer observation, candidate, Evidence, Finding et décision humaine.

## 4. Non-objectifs
Aucune acquisition, commande, méthode offensive, moteur, plugin, offset, algorithme, API, protocole, format final, action Endpoint, Disk/Filesystem/full Network Forensics, Cloud, Mobile ou règle Detection Engineering.

## 5. Propriétaire
Investigate possède contexte et interprétation; Endpoint Agent l’acquisition; Settings Fleet/Policies/stockage/rétention/santé; Studio Tools/Runs; Govern les cibles réelles; Shared les mécanismes transversaux.

## 6. Utilisateurs
Principal : **Evidence Reviewer**. Secondaires : Investigation Lead et Audit Analyst autorisés.

## 7. Conditions d’entrée
Case et Memory Image lisibles; versions, restrictions, permissions et sources visibles. Toute absence devient `partial`, `restricted`, `unsupported` ou `blocked`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Memory Image | Investigate | source et version | oui | version sélectionnée | review impossible |
| Acquisition context | CAP-INV-207/203 | initiateur, méthode déclarée, timestamps et erreurs | oui | état final du Job | incomplete |
| Custody and provenance | CAP-INV-213/214 | transformations, transferts, copies et vérifications | oui | chaîne courante | disputed |
| Permissions et restrictions | Settings / Security | lecture, export et conservation | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Memory Image / Artifact | Investigate | source, taille, statut et version | lecture |
| Collection Request / Job | Investigate | scope, progression et erreurs | lecture |
| Endpoint / Agent | Endpoint Agent / Settings | source et support déclaré | projection lecture |
| Custody / provenance records | Investigate + Shared | événements et lacunes | lecture/review |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Image Integrity Review | créer/modifier/contester | Investigate concept | distinct de pertinence |
| Memory Image relation | annoter statut fonctionnel | Investigate | original préservé |
| Custody correction | ajouter événement | Investigate / Trust | ancienne entrée conservée |
| Trace event | émettre | Shared | append-only |

## 11. Fonctionnalités
- voir source, Endpoint, initiateur, autorisation, méthode déclarée, timestamps et durée.
- voir statut de collecte, erreurs et parties manquantes.
- voir transformations, transferts, copies, vérifications et incohérences.
- contester, annoter, accepter pour analyse ou marquer partielle/non exploitable.
- préserver l’original et la chaîne de custody.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Consulter la chaîne | Evidence Reviewer | Memory Image | 0 | lecture autorisée | vue sourcée | non |
| Vérifier ou contester | Evidence Reviewer | Integrity Review | 2 | permission et justification | version attribuée | OPEN-013 |
| Marquer partial/non exploitable | Evidence Reviewer | Integrity Review | 2 | permission et justification | statut versionné | OPEN-013 |
| Accepter pour analyse | Evidence Reviewer | Session precondition | 2 | sources et limites visibles | décision attribuée | OPEN-013 |

Classes 3/4 bloquées et routées vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| assembler acquisition et custody | oui | oui | oui | explication | chaîne déterministe |
| signaler les parties manquantes | oui | oui | oui | résumé | table et checklist |
| comparer transformations et copies | oui | oui | oui | assistance | comparaison déterministe |
| Conclusion finale | oui | non | non | assistance | revue humaine |

Toute automatisation expose initiateur, producteur/version, Tool Calls, sources, paramètres, statut, erreurs et disposition humaine.

## 14. États fonctionnels
`unverified`, `verifying`, `verified`, `partially-verified`, `inconsistent`, `incomplete`, `corrupted`, `disputed`, `restricted`. Machine objet finale reportée.

## 15. États d’interface
Loading conserve le contexte; Empty n’invente rien; Partial expose les manques; Error conserve le valide; Offline limite les mutations; Permission denied masque; Stale distingue ancien et courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Integrity assessment | Image Integrity Review | CAP-INV-347/348 | statut, sources et contestations visibles |
| Acquisition context package | projection | CAP-INV-361 | initiateur, timestamps, erreurs et custody conservés |
| Partial/corrupted warning | Trace event | toutes analyses | hérité par les résultats |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Collection Job | Memory Image produite | CAP-INV-349 | request, Endpoint, acquisition, errors et custody | Collection |
| CAP-INV-349 | image acceptable | CAP-INV-348 | image, statut, restrictions et objectif | Integrity Review |
| CAP-INV-349 | contestation | Case / Evidence Review | sources, écarts et auteur | Integrity Review |

Tenant, Case, image, permissions, restrictions et return origin sont préservés.

## 18. Dépendances
CAP-INV-203, CAP-INV-207, CAP-INV-213, CAP-INV-214, CAP-INV-105, Shared Trace, OPEN-008, OPEN-013, OPEN-014. Aucune dépendance bas niveau.

## 19. Source de vérité
Image : Investigate; acquisition : Collection/Endpoint Agent; restrictions : Settings; Trace : Shared.

## 20. Provenance et audit
Case, Endpoint, Request/Job, image, custody, transformations, vérifications, acteur, timestamps, erreurs, contestations et dispositions; aucune suppression silencieuse.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up | Séparation | Owner | Phase |
|---|---|---|---|---|---|---|
| Consulter chaîne | lecture sensible | 0 | selon policy | acquisition/review distincts | Investigate/Security | Permissions |
| Vérifier/contester | mutation réversible | 2 | possible | reviewer autorisé | Investigate/Security | Permissions |
| Accepter pour analyse | mutation réversible | 2 | possible | reviewer distinct si requis | Investigate/Security | Permissions |

Matrice atomique, RBAC/ABAC et règles finales reportés.

## 22. Limites et erreurs
- aucun algorithme d’intégrité, format ou stockage choisi.
- intégrité distincte de pertinence.
- aucune modification de l’image.
- Memory Forensics distincte du Debugger, Disk et full Network Forensics.

## 23. Métriques
- images verified/partial/disputed.
- maillons de custody manquants.
- images non exploitables.
- complétude de provenance.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, plugin ou plateforme.

## 25. Critères d’acceptation
### 1. Image partielle
**Given** des plages manquent **When** la revue est ouverte **Then** la partialité est visible et aucun statut complet n’est inventé.
### 2. Custody incohérente
**Given** des événements divergent **When** le reviewer conteste **Then** la contestation est attribuée et l’original reste intact.
### 3. Sans IA
**Given** aucun modèle **When** la revue est réalisée **Then** relations, timestamps et checklist suffisent.

## 26. Questions ouvertes
OPEN-008, OPEN-013, OPEN-014 restent ouvertes. Schémas, formats, plateformes et permissions finales sont reportés.

## 27. Consommateurs documentaires
Memory module/INV-MEM-001, Case/Evidence, CAP-INV-107/108/109/311, futures phases Objects/Permissions/Journeys/Screens/Contracts.
