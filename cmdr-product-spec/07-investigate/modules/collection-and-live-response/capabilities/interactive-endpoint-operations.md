---
id: CAP-INV-210
title: Interactive Endpoint Operations
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
  - REQ-SEC-002
open_decisions:
  - OPEN-007
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-210 — Interactive Endpoint Operations

## 1. Définition
Sélectionner, valider, exécuter et interrompre uniquement des catégories d’opérations endpoint autorisées, classées et traçables, sans publier de commande ou script réel.

## 2. Problème utilisateur
Une action interactive mal classée peut contourner Govern ou transformer une observation en mutation. L’opérateur doit voir effets, préconditions, classe, autorité, rollback attendu et résultat avant et après l’exécution.

## 3. Objectifs
- présenter un catalogue fonctionnel d’observation, inspection, collecte, transfert, exécution autorisée, modification réversible et demandes de containment
- afficher cible, scope, effets, classe, policy, permission et besoin d’approbation
- bloquer les classes 3/4 sans Govern et les opérations interdites
- permettre interruption et inspection de la sortie sans confondre output et Result Govern

## 4. Non-objectifs
Ne fournir aucune commande système, script offensif, persistance, évasion ou contournement ; ne pas définir le command executor ; ne pas exécuter containment/destruction directement dans Investigate.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Investigation Lead, Case Analyst, approbateur Govern et participant de Live Session autorisé.

## 7. Conditions d’entrée
Live Session active ou contexte d’opération autorisé, Endpoint available, opération cataloguée, cible/scope, classe, policy, permission, effets et rollback attendus visibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Live Session et Case | Investigate | contexte interactif | oui sauf opération non session explicitement permise | état courant | opération interdite |
| Endpoint/Agent capability | Endpoint Agent | faisabilité locale | oui | dernière communication | offline/unsupported |
| Operation catalogue entry | owner catalogue futur | description, effets et classe | oui | version visible | validation-required |
| Scope, paramètres fonctionnels et rollback | opérateur / catalogue | contrat d’exécution | oui selon classe | snapshot avant run | incomplete |
| Permission, policy et autorité | Security / Settings / Govern | gate | oui | snapshot à l’exécution | denied/awaiting-approval |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Finding / Evidence | Investigate | objectif et justification | consulter/référencer |
| Live Session | Investigate, modèle futur | participants, statut et trace | consulter/opérer |
| Endpoint / Endpoint Agent / Policy | partagé / Agent / Settings | capability, disponibilité et restrictions | consulter |
| Agent Command | Endpoint Agent | opération locale et statut | consulter en projection |
| Decision / Response Run | Govern | autorité pour classes gouvernées | consulter uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Operation record conceptuel | créer, actualiser, interrompre ou supersede | Investigate, modèle futur | catalogue/version, scope, class et session obligatoires |
| Agent Command relation | créer une demande autorisée | Endpoint Agent | ne transfère pas ownership ni autorité |
| Operation Result relation | créer à réception | Investigate/Agent source | distinct de Govern Result |
| Action Request draft | préparer pour classe 3/4 | Govern lifecycle via CAP-INV-215/113 | aucune exécution directe |

## 11. Fonctionnalités
- parcourir les catégories autorisées sans exposer commandes réelles
- voir description, préconditions, effets, classe, scope, autorité et rollback
- valider la cible et exécuter uniquement si toutes les gates passent
- interrompre une opération et conserver output/errors partiels
- rediriger classes 3/4 vers CAP-INV-215/113 et Govern

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Consulter catalogue/effets | Response Operator | operation definition | 0 | read | classe et limites visibles | non |
| Exécuter observation/inspection | Response Operator | Endpoint Operation | 0/1 | session active et permission | operation running | non normalement |
| Exécuter modification réversible | Response Operator | Endpoint Operation | 2 | rollback, policy, permission/approval | operation running | OPEN-013 selon policy |
| Interrompre | Response Operator | Endpoint Operation | 2 | interruptible | interrupted/partial avec trace | selon policy |
| Préparer classe 3/4 | Investigation Lead | Action Request | 3/4 | Finding/Evidence/impact/rollback | handoff Govern | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Proposer une opération | oui | catalogue/règles | oui | suggestion | catalogue manuel |
| Valider préconditions | oui | validateurs | oui | explication | checklist |
| Résumer output/errors | oui | agrégation | oui | résumé attribué | sorties brutes |
| Proposer interruption/next step | oui | règles | oui | suggestion | décision opérateur |
| Exécuter action sensible | humain explicite | contrat autorisé | workflow possible | jamais autonome | opérateur/Govern |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`draft`, `validation-required`, `ready`, `queued`, `running`, `interrupted`, `partial`, `completed`, `failed`, `cancelled`, `denied`, `policy-blocked`. Machine finale reportée.

## 15. États d’interface
Loading conserve opération/scope ; Empty signifie catalogue non disponible ; Partial affiche output/errors reçus ; Error garde trace et résultats ; Offline n’indique pas success ; Permission denied masque paramètres sensibles ; Stale exige revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Operation | record conceptuel | CAP-INV-212/214 | catalogue/version, cible, classe, initiateur et statut |
| Operation output/error | Operation Result source | CAP-INV-212 | partiel/terminal distingués |
| Artifact éventuel | Artifact | CAP-INV-105/107 | création explicite et provenance |
| Action Request context | package | CAP-INV-215/113/Govern | classe 3/4 non exécutée localement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Live Session | sélectionner opération | CAP-INV-210 | session, target, participant, catalogue, policy | Live Session |
| CAP-INV-210 | exécution autorisée | Endpoint Agent | operation, scope, class, authority, correlation | CAP-INV-210/Session |
| Endpoint Agent | output/error | CAP-INV-212 | operation, target, output, errors, files, timestamps, status | Live Session/Case |
| Classe 3/4 | préparer demande | CAP-INV-215/113 | Finding, Evidence, Endpoint, action, impact, rollback | Case |

## 18. Dépendances
CAP-INV-209/211/212/214/215/113, Endpoint Agent command executor target, Settings Policy, Govern authority, Studio optional Workflow, Shared Trace et OPEN-007/008/013/015.

## 19. Source de vérité
Investigate possède la sélection, le record métier et les relations au Case/Session. Endpoint Agent reste source de l’exécution locale ; Govern de l’autorité et des Runs ; Studio de toute Automation Run.

## 20. Provenance et audit
Case, session, Endpoint/Agent, catalogue/version, initiateur, scope, class, policy, permission, approval, rollback, start/interrupt/end, output/errors, Artifacts et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint operation execute/interrupt, catalogue read, sensitive output read, session access, class-2 step-up et containment request. Aucune permission n’est auto-accordée.

## 22. Limites et erreurs
Operation unsupported, target stale, session expired, policy changed, permission revoked, approval absent, output truncated, interruption impossible, disconnect ou result late. Output ≠ verified Govern Result.

## 23. Métriques
Operations par classe/statut, denied/policy-blocked, interruptions, partial/failures, classes 3/4 redirigées et opérations avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`. Aucun catalogue final, commande, script, executor, plateforme ou protocole n’est livré.

## 25. Critères d’acceptation
**Given** une opération classe 3 sans Decision **When** l’analyste tente de l’exécuter **Then** l’exécution directe est bloquée et une Action Request peut être préparée.

**Given** une opération interrompue **When** des sorties partielles existent **Then** elles restent visibles comme partial/interrupted sans être déclarées succès.

**Given** aucun modèle IA **When** l’opérateur sélectionne une opération **Then** catalogue, préconditions, classes et contrôles déterministes permettent le workflow.

## 26. Questions ouvertes
OPEN-007, OPEN-008, OPEN-013 et OPEN-015 restent ouvertes. Le catalogue final, les commandes et l’executor appartiennent aux phases Objets/Permissions/Technique.

## 27. Consommateurs documentaires
Live Session, Result Handling, Provenance, Artifact/Evidence, Containment Preparation, Govern, Endpoint Agent Live Response et phases Objets/Permissions.
