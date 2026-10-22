---
id: CAP-INV-215
title: Containment Request Preparation
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-016
  - REQ-SEC-002
open_decisions:
  - OPEN-007
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-215 — Containment Request Preparation

## 1. Définition
Spécialiser CAP-INV-113 pour préparer un containment Endpoint avec cible, impact, urgence, alternatives et rollback, sans lifecycle ni exécution concurrente.

## 2. Problème utilisateur
Sans Containment Request Preparation, l’utilisateur perd le lien entre le Case, l’Endpoint, l’autorité applicable, l’exécution locale et les résultats. Les états partiels ou offline peuvent alors être pris pour un succès et les objets peuvent être confondus.

## 3. Objectifs
- fournir Finding, Evidence, Endpoint, action, impact, urgence, alternatives, rollback;
- exposer cible, scope, fraîcheur, policy, permission et classe d’action;
- conserver erreurs, résultats partiels, provenance et retour au Case;
- produire package vers CAP-INV-113/Govern sans transférer l’ownership.

## 4. Non-objectifs
- ne pas administrer la Fleet ni les Endpoint Policies;
- ne pas définir protocole, API, commande, moteur, format, PKI, stockage ou plateforme supportée;
- ne pas créer automatiquement Evidence, Finding, Decision, Response Run ou Govern Result;
- ne pas commencer Analysis Workbench.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier et les relations au Case. Platform Settings administre Fleet/Policies; Endpoint Agent exécute localement; Govern possède l’autorité risquée.

## 6. Utilisateurs
Principal : Investigation Lead / Response Operator. Secondaires : Evidence Reviewer, Incident Commander, approbateur Govern ou Platform Administrator en consultation.

## 7. Conditions d’entrée
Tenant et environnement conservés, Case accessible, Endpoint résolu, fraîcheur et capacités visibles, Finding et Evidence sélectionnés, policy projetée, permissions vérifiées, impact et rollback explicités.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte métier | oui | version courante | rester draft |
| Endpoint et état Agent | Platform Settings / Endpoint Agent | cible et disponibilité | oui | dernière communication visible | offline/unknown explicite |
| Finding et Evidence | Investigate | justification analytique | oui | versions et statuts visibles | incomplete |
| Action, impact, urgence, alternatives, rollback | analyste / policy | package de décision | oui | revalidés à la soumission | ready interdit |
| Autorité et permission | Security / Govern | gate | oui | snapshot à l’action | awaiting-approval/denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte et liens | consulter |
| Endpoint / Endpoint Agent / Endpoint Policy | partagé / Agent / Settings | cible, capacité et restrictions | consulter uniquement |
| Finding / Evidence | Investigate | justification, incertitude et versions | consulter/sélectionner |
| Action Request | Govern | draft et retour | préparer via CAP-INV-113, consulter |
| Decision / Response Run / Result | Govern | statut futur et retour | consulter uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint containment package | créer, modifier, supersede ou annuler avant remise | Investigate | spécialisation de contenu, pas nouvel objet d’autorité |
| Action Request draft | préparer via CAP-INV-113 | Govern lifecycle | Investigate producteur ; Govern owner |
| Relations vers Finding/Evidence/Endpoint | créer/supersede | Linking / owners respectifs | sources et versions conservées |
| Decision/Response Run/Result | aucune création locale | Govern | projection seulement |

## 11. Fonctionnalités
Sélectionner Finding/Evidence/Endpoint, décrire l’action proposée, impact, urgence, alternatives, conditions, rollback et incertitudes ; classer classe 3/4 ; vérifier complétude ; remettre à CAP-INV-113 ; soumettre à Govern ; suivre retour et objets futurs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter cible et justification | analyste | projections | 0 | read | package sourcé | non |
| Créer/modifier package | Investigation Lead | draft spécialisé | 2 | Case/Finding/Evidence | draft versionné | OPEN-013 |
| Remettre à CAP-INV-113 | Investigation Lead | Action Request draft | 3 | complet | draft Govern préparé | obligatoire |
| Soumettre | humain autorisé | Action Request | 3/4 | permission et gate | lifecycle Govern | obligatoire |
| Suivre retour | analyste | Decision/Run/Result projections | 0 | lien valide | Case mis à jour par relation | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler Finding/Evidence/target | oui | relations | oui | résumé | sélection manuelle |
| Vérifier complétude | oui | checklist | oui | explication | validateur |
| Proposer impact/alternatives | oui | catalogues/policy | oui | brouillon | saisie et modèles |
| Proposer rollback | oui | catalogue autorisé | oui | suggestion | sélection humaine |
| Soumettre/exécuter | humain explicite | gate | workflow possible | jamais autonome | CAP-INV-113/Govern |

Toute sortie automatisée expose initiateur, moteur/agent/version, Automation Run, Tool Calls, sources, paramètres, timestamp, incertitude, owner humain, accept/modify/reject et trace.

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-cap-113`, `submitted`, `returned-for-information`, `superseded`, `cancelled-before-submission`. Aucun état ne signifie containment exécuté.

## 15. États d’interface
Loading conserve le draft ; Empty demande Finding/Evidence/Endpoint ; Partial nomme données/rollback manquants ; Error préserve le draft ; Offline n’indique aucune exécution ; Permission denied masque les données ; Stale exige revalidation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint containment package | contenu spécialisé | CAP-INV-113 | cible, Finding, Evidence, impact, alternatives et rollback versionnés |
| Action Request draft | projection Govern | Govern Review Queue | aucun ownership local du lifecycle |
| Return/status link | relation | Case Workspace | Decision/Run/Result restent Govern |
| Provenance event | événement | CAP-INV-214/Timeline/Trace | auteur, sources, disposition et correlation ID |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Finding/Case | préparer containment | CAP-INV-215 | Finding, Evidence, Endpoint, impact, urgence, return origin | Case |
| CAP-INV-215 | package complet | CAP-INV-113 | package Endpoint, alternatives, rollback, incertitudes | CAP-INV-215/Case |
| CAP-INV-113 | soumission | Govern | Action Request, classe, target, risques, provenance | Case |
| Govern Result | retour | Case Workspace | Decision, Response Run, Result, vérification, effet Finding | Case/Result |

## 18. Dépendances
CAP-INV-109/113/201/212/214, Govern Action Center/Policy Gates/Runs/Rollback, Endpoint Agent containment capability, Settings Policy, OPEN-007/008/013/015.

## 19. Source de vérité
Finding/Evidence/Case et package préparatoire restent Investigate ; Action Request lifecycle, Decision, Response Run, Result, rollback et autorité restent Govern ; exécution locale reste Endpoint Agent.

## 20. Provenance et audit
Auteur, Case, Endpoint/Agent/Policy, Finding/Evidence versions, action, impact, urgence, alternatives, conditions, rollback, incertitudes, validations, submission, CAP-INV-113 handoff et objets Govern.

## 21. Permissions fonctionnelles
Containment request prepare/submit, Finding/Evidence read, Endpoint read, sensitive output, cross-tenant/environment, step-up, separation of duties et Govern review. Matrice finale reportée.

## 22. Limites et erreurs
Finding contesté, Evidence superseded, Endpoint stale/unsupported, rollback absent, impact inconnu, policy blocked, permission refusée ou Govern indisponible gardent un draft ; aucune isolation ou modification n’est présentée comme exécutée.

## 23. Métriques
Packages incomplets, retours pour information, délais de soumission, liens Finding/Evidence complets, supersessions et retours Govern liés au Case.

## 24. Classification de livraison
`defined` / `planned`. Aucune exécution de containment, Decision, Run ou plateforme n’est prouvée ; promotion dépend de Govern, Endpoint Agent, Permissions et OPEN-007/008/013/015.

## 25. Critères d’acceptation
**Given** un Endpoint lié à un Finding et aucune Decision **When** l’analyste tente le containment **Then** l’exécution directe est bloquée et un package CAP-INV-113/Govern peut être préparé.

**Given** une demande retournée **When** elle est complétée **Then** les liens Finding/Evidence/Endpoint et l’historique sont conservés.

**Given** aucun modèle IA **When** le package est préparé **Then** formulaires, checklists, catalogues et revue humaine permettent la soumission.

## 26. Questions ouvertes
OPEN-007, OPEN-008, OPEN-013 et OPEN-015 restent ouvertes. CAP-INV-215 ne ferme ni ne remplace CAP-INV-113.

## 27. Consommateurs documentaires
CAP-INV-113, Govern Action Center/Policy Gates/Runs, Endpoint Agent containment, Case Workspace, parcours Incident/Endpoint containment, permissions et future phase Objets.