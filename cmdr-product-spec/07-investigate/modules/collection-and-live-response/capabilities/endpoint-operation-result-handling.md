---
id: CAP-INV-212
title: Endpoint Operation Result Handling
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
  - REQ-OBJ-007
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-212 — Endpoint Operation Result Handling

## 1. Définition
Recevoir, inspecter, vérifier, contester et relier le résultat local d’une opération endpoint sans le confondre avec Artifact, Evidence, Finding ou Govern Result.

## 2. Problème utilisateur
Une sortie de commande ou un statut local peut être incomplet, non vérifié ou contradictoire. Le Case doit conserver la différence entre output, Operation Result, Artifact, Evidence, Finding et Result d’un Response Run Govern.

## 3. Objectifs
- afficher opération source, Endpoint, Case, initiateur, autorité et timestamps
- séparer statut, output, errors, éléments partiels, fichiers et Artifacts
- permettre vérification, contestation, supersession et nouvelle opération
- lier à Evidence/Finding/Govern/Command sans promotion automatique

## 4. Non-objectifs
Ne pas créer ou posséder Govern Result, confirmer Finding, qualifier Evidence automatiquement, définir format d’output ou protocole ; ne pas transformer stdout/stderr en résultat vérifié.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Case Analyst ou Response Operator. Secondaires : Evidence Reviewer, Investigation Lead, Incident Commander et reviewer Govern.

## 7. Conditions d’entrée
Operation/Agent Command identifiable, Endpoint/Case liés, output/status/errors/timestamps disponibles ou lacunes explicites, permissions de raw/sensitive result et provenance accessibles.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Endpoint Operation et source command | Investigate / Endpoint Agent | origine du résultat | oui | version/statut courants | result orphaned/incomplete |
| Case, Endpoint et initiateur | owners respectifs | contexte et responsabilité | oui | relations courantes | review limitée |
| Output, errors, files et partial items | Endpoint Agent | contenu local | selon opération | timestamps visibles | incomplete |
| Decision/Response Run éventuels | Govern | autorité gouvernée | non | version/statut visibles | Operation Result non Govern |
| Verification context | reviewer / telemetry | contrôle indépendant | non | source/time visibles | received non verified |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Operation / Agent Command | Investigate concept / Endpoint Agent | scope, class, initiator et status | consulter |
| Case / Endpoint | Investigate / partagé | contexte et relations | consulter/lier |
| Artifact / Evidence / Finding | Investigate | sorties et qualifications | consulter/lier selon permission |
| Decision / Response Run / Result | Govern | autorité et résultat canonique éventuels | consulter uniquement |
| Automation Run / Tool Call | Studio | provenance éventuelle | consulter uniquement |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Operation Result record conceptuel | créer, vérifier, contester ou supersede | Investigate, modèle futur ; source Agent | distinct de Govern Result |
| Artifact relation | créer/lier pour fichier ou capture matérialisée | Investigate | source operation/result obligatoire |
| Evidence candidate relation | proposer | Investigate | qualification CAP-INV-107/108 obligatoire |
| Finding/Govern/Command relation | créer/supersede | owners respectifs / Linking | aucune mutation automatique des objets sources |

## 11. Fonctionnalités
- afficher source operation, target, initiator, authority, timestamps et duration
- présenter output, errors, partial items et files séparément
- marquer received, incomplete, under-review, verified, disputed ou superseded
- créer/lier un Artifact et proposer une Evidence candidate
- demander une nouvelle opération ou transmettre le contexte à Govern/Command

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Lire output/errors | utilisateur autorisé | Operation Result | 0 | raw/sensitive read | vue sourcée | non |
| Vérifier | reviewer | verification status | 2 | contrôle et raison | verified ou incomplete | OPEN-013 |
| Contester/supersede | reviewer | Operation Result | 2 | grounds/nouvelle source | historique conservé | OPEN-013 |
| Créer/lier Artifact | Case Analyst | Artifact | 1/2 | sortie matérielle et provenance | Artifact lié | non |
| Proposer Evidence/nouvelle opération | Evidence Reviewer | candidate/request | 1/2 | raison et permissions | handoff CAP-INV-107/202 | selon impact |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Regrouper output/errors | oui | par statut/type | oui | résumé | filtres bruts |
| Expliquer une erreur | oui | catalogue | oui | oui | message déterministe |
| Proposer Artifact/Evidence candidate | oui | règles explicites | oui | suggestion | sélection humaine |
| Comparer avec vérification | oui | diff/checks | oui | résumé | inspection manuelle |
| Marquer verified/Evidence/Finding | humain autorisé | contrôles seulement | workflow de revue | jamais autonome | revue humaine |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`received`, `incomplete`, `under-review`, `verified`, `disputed`, `superseded`, `linked-to-artifact`, `linked-to-evidence`. Machine finale reportée.

## 15. États d’interface
Loading conserve operation/Case ; Empty distingue aucun output et redaction ; Partial affiche chaque item ; Error garde les sorties valides ; Offline n’altère pas le statut terminal ; Permission denied masque le contenu ; Stale expose versions/supersession.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Operation Result | record conceptuel | Case Workspace/Live Session | source, status, output, errors et verification visibles |
| Artifact relation | Artifact | CAP-INV-105/107 | source operation/result conservée |
| Evidence candidate | relation | CAP-INV-107/108 | aucune qualification automatique |
| Govern/Command handoff | context package | CAP-INV-113/Command | Operation Result reste distinct de Govern Result |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Endpoint Agent/Operation | résultat reçu | CAP-INV-212 | operation, target, output, errors, files, timestamps, status, provenance | Live Session/Case |
| CAP-INV-212 | materialiser output | CAP-INV-105 | file/capture, source, acquisition, result status | CAP-INV-212 |
| CAP-INV-212 | qualifier candidate | CAP-INV-107/108 | Artifact/result, Case, reason, verification | CAP-INV-212 |
| CAP-INV-212 | demander action ou informer | CAP-INV-202/113 ou Command | result, verification, next action, trace | Case |

## 18. Dépendances
CAP-INV-105/107/108/109/113/203/209/210/211/214, Endpoint Agent result reporting, Govern Result, Studio provenance, Shared Linking/Trace et OPEN-013/015.

## 19. Source de vérité
Endpoint Agent reste source de l’output local déclaré. Investigate possède le record de revue, les relations Artifact/Evidence/Case et la disposition analytique. Govern reste source de son Result canonique.

## 20. Provenance et audit
Case, Endpoint/Agent, operation/command, initiator, policy/authority, start/end/duration, output/errors/files, partial items, reviewer, verification/dispute, Artifact/Evidence relations, Automation/Govern refs et correlation IDs.

## 21. Permissions fonctionnelles
Raw/sensitive result read, transcript read, result verify/dispute, Artifact create/link, Evidence candidate, Finding link et cross-tenant restrictions. Separation of duties reportée.

## 22. Limites et erreurs
Output absent/truncated/redacted, result late/duplicate, operation interrupted, disconnect, source command inaccessible, verification contradictoire, Artifact corrupt ou permission retirée. Output ≠ verified result.

## 23. Métriques
Results received/incomplete/verified/disputed, output truncation, Artifact/Evidence candidate rates, time to review et Operation Results incorrectly linked to Govern Result cible zéro.

## 24. Classification de livraison
`defined` / `planned`. Aucun format de résultat, stockage, protocol ou verification engine n’est choisi.

## 25. Critères d’acceptation
**Given** une opération Live Response terminée sans Response Run **When** l’utilisateur consulte la sortie **Then** elle est présentée comme Endpoint Operation Result, pas comme Govern Result, et peut produire un Artifact.

**Given** un résultat partiel **When** le reviewer l’ouvre **Then** output, erreurs et éléments manquants restent distincts et aucun succès complet n’est affiché.

**Given** aucun modèle IA **When** le résultat est revu **Then** output brut, filtres, checks et disposition humaine permettent le workflow.

## 26. Questions ouvertes
OPEN-013 conserve la gouvernance de verification/supersession ; OPEN-015 le bridge Automation Run/Response Run. L’objet Operation Result et son lifecycle final restent à la phase Objets.

## 27. Consommateurs documentaires
Live Session, Collection Job, Artifact/Evidence/Finding, Case Timeline/Replay, Govern/Command handoffs, Endpoint Agent result reporting et phases Objets/Permissions.
