---
id: CAP-INV-214
title: Collection and Live Response Provenance
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
open_decisions:
  - OPEN-007
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-214 — Collection and Live Response Provenance

## 1. Définition
Émettre la sémantique métier de provenance Collection/Live Response vers Trace, Activity, Timeline et audit sans les dupliquer.

## 2. Problème utilisateur
Sans Collection and Live Response Provenance, l’utilisateur perd le lien entre le Case, l’Endpoint, l’autorité applicable, l’exécution locale et les résultats. Les états partiels ou offline peuvent alors être pris pour un succès et les objets peuvent être confondus.

## 3. Objectifs
- fournir tous IDs et événements Collection/Session/Operation/Govern/Studio;
- exposer cible, scope, fraîcheur, policy, permission et classe d’action;
- conserver erreurs, résultats partiels, provenance et retour au Case;
- produire chaîne de provenance et événements métier sans transférer l’ownership.

## 4. Non-objectifs
- ne pas administrer la Fleet ni les Endpoint Policies;
- ne pas définir protocole, API, commande, moteur, format, PKI, stockage ou plateforme supportée;
- ne pas créer automatiquement Evidence, Finding, Decision, Response Run ou Govern Result;
- ne pas commencer Analysis Workbench.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède la sémantique métier et les relations au Case. Platform Settings, Endpoint Agent, Govern, Studio et Shared conservent leurs mécanismes et objets.

## 6. Utilisateurs
Principal : Case Reviewer / Auditor. Secondaires : Investigation Lead, Evidence Reviewer, Incident Commander, approbateur Govern ou Platform Administrator en consultation.

## 7. Conditions d’entrée
Case et Endpoint résolus, IDs des requests/jobs/sessions/opérations, sources productrices, policy/permission/autorité et événements horodatés disponibles ou lacunes explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case, Endpoint, Agent et Policy | owners respectifs | contexte et cible | oui | versions visibles | chaîne partial |
| Request, Job, Session, Operation et Result IDs | Investigate / Agent | événements métier | selon workflow | ordonnés et sourcés | broken-link |
| Automation Run et Tool Calls | Studio | provenance d’automation | si automation | version/run visibles | producer unknown explicite |
| Action Request, Decision, Response Run et Result | Govern | autorité et retour | si gouverné | statut/version visibles | aucune autorité inventée |
| Human disposition | utilisateur | accept/modify/reject/verify | si suggestion/revue | horodatée | unresolved |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Collection Request | Investigate | contexte et scope | consulter |
| Endpoint / Endpoint Agent / Policy | partagé / Agent / Settings | cible, exécution et restrictions | consulter |
| Collection Job / Live Session / Endpoint Operation / Operation Result | concepts métier | séquence et résultats | consulter |
| Automation Run / Tool Call / Workflow | Studio | producteur et outils | consulter |
| Action Request / Decision / Response Run / Result | Govern | autorité et résultat gouverné | consulter |
| Artifact / Evidence / Finding | Investigate | sorties et qualification | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Événement métier Collection/Live Response | émettre/corriger/supersede | Investigate sémantique, Shared mechanism | source, acteur, timestamp et correlation ID |
| Provenance link | créer/supersede | Object Linking / owner source | version et rôle relationnel obligatoires |
| Human disposition event | enregistrer | Investigate/Studio selon source | accept/modify/reject et owner humain visibles |
| Trace/Activity/Timeline/Audit stores | aucune redéfinition | Shared/Govern | consommation uniquement |

## 11. Fonctionnalités
Retracer Case, Endpoint, Agent, policy, permission, initiateur, Workflow, Automation Run, Tool Calls, Action Request, Decision, Response Run, opération, paramètres fonctionnels, résultats, erreurs, Artifacts, transformations et disposition humaine ; comparer et naviguer.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter/comparer | reviewer | provenance chain | 0 | read | chaîne sourcée | non |
| Annoter/corriger | reviewer autorisé | événement métier | 2 | raison/version | nouvel événement, ancien conservé | OPEN-013 |
| Exporter sous contrôle | auditor | export request | 1/2 | permission/redaction | package Export | selon sensibilité |
| Ouvrir source | utilisateur | objet propriétaire | 0 | lien/permission | navigation avec return origin | non |
| Contester | reviewer | provenance status | 2 | cause | disputed | non par défaut |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Corréler les IDs | oui | oui | oui | non nécessaire | relations déterministes |
| Détecter liens manquants | oui | règles | oui | explication | checklist |
| Résumer transcript/chaîne | oui | agrégation | oui | résumé attribué | filtres/timeline |
| Attribuer une suggestion | oui | metadata | oui | oui | lecture run/tool calls |
| Modifier une source propriétaire | humain/owner | contrat owner | selon owner | jamais silencieux | navigation source |

Toute automation expose initiateur, agent/moteur/version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, incertitude, owner humain, accept/modify/reject et trace.

## 14. États fonctionnels
`complete`, `partial`, `broken-link`, `redacted`, `superseded`, `disputed`, `unavailable`. Ce ne sont pas des états finaux d’objet.

## 15. États d’interface
Loading conserve Case/Endpoint ; Empty explique l’absence de source ; Partial nomme les liens manquants ; Error conserve événements valides ; Offline montre le dernier snapshot ; Permission denied redacted ; Stale expose versions superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Provenance chain | relations/événements | Case Replay, Timeline, audit | producteurs, sources, statuts et lacunes visibles |
| Business events | événements métier | Shared Trace/Activity/Timeline | aucune duplication de mécanisme |
| Automation attribution | liens Studio | analyste/reviewer | Run, Tool Calls, version et disposition visibles |
| Govern attribution | liens Govern | Case/Incident | Decision/Response Run/Result distincts |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Request/Job/Session/Operation | événement | CAP-INV-214 | IDs, cible, acteur, policy, paramètres, statut, erreurs | source |
| CAP-INV-214 | revoir l’investigation | CAP-INV-112/110 | chaîne, lacunes, corrections, dispositions | provenance view |
| Artifact/Evidence/Finding | ouvrir origine | CAP-INV-214 | relation/version/return origin | objet source |
| Govern/Studio projection | ouvrir run | owner destination | IDs, source et return origin | provenance view |

## 18. Dépendances
CAP-INV-008/110/112/203/209/212/213/215, Shared Trace/Activity/Timeline/Linking/Export, Studio Workflow/Automation Run/Tool Calls, Govern Action Request/Decision/Run/Result, OPEN-007/015.

## 19. Source de vérité
Investigate émet la sémantique de ses événements ; Endpoint Agent, Settings, Studio et Govern restent sources de leurs objets ; Shared/Govern restent propriétaires des mécanismes de trace/audit.

## 20. Provenance et audit
Case, Endpoint, Agent/Policy, permissions, initiateur, automation, tools, authority, requests, jobs, sessions, operations, outputs, errors, Artifacts, transformations, human dispositions et correlation IDs.

## 21. Permissions fonctionnelles
Provenance read/export, transcript read, raw result read, sensitive output, cross-tenant/environment, automation/govern projection read et audit access. Matrice atomique reportée.

## 22. Limites et erreurs
Lien cassé, source indisponible, redaction, clock inconsistency, version superseded, producer inconnu, tenant mismatch ou permission refusée rendent la chaîne partial/disputed ; aucune source n’est recréée localement.

## 23. Métriques
Chaînes complètes, liens cassés, événements sans producteur, suggestions sans disposition, transitions Govern/Studio attribuées et navigation source réussie.

## 24. Classification de livraison
`defined` / `planned`. Aucun audit ledger, Trace store, API ou contrat technique n’est déclaré livré ; OPEN-007/015 restent ouvertes.

## 25. Critères d’acceptation
**Given** une collecte automatisée **When** la provenance est ouverte **Then** initiateur, Workflow, Automation Run, Tool Calls, target, paramètres, résultats, erreurs et disposition humaine sont visibles.

**Given** un Response Run Govern **When** il est relié au Case **Then** il reste distinct d’Automation Run et Operation Result.

**Given** aucun modèle IA **When** la chaîne est revue **Then** relations, timestamps, filtres et sources permettent la reconstruction.

## 26. Questions ouvertes
OPEN-007 et OPEN-015 restent ouvertes. Les objets Provenance/Audit, rétention et contrats techniques appartiennent aux phases Objets/Trust/Technique.

## 27. Consommateurs documentaires
Case Timeline/Replay, Evidence/Finding, Govern Runs/Audit, Studio provenance, Endpoint Agent result reporting, Shared Trace/Activity/Timeline/Export et phases Objets/Permissions/Trust.