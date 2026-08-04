---
id: CAP-INV-326
title: Dynamic Analysis Safety and Containment
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-014
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-326 — Dynamic Analysis Safety and Containment

## 1. Définition
Définir et superviser les attentes fonctionnelles de sécurité de l’analyse dynamique — isolation, réseau, durée, ressources, interactions, fichiers, interruption, arrêt d’urgence, nettoyage et reset — sans spécifier les mécanismes techniques ni prétendre à l’invulnérabilité.

## 2. Problème utilisateur
Une sandbox peut être mal configurée, compromise ou sortir de ses limites. Sans états et procédures fonctionnelles explicites, un nouveau Run peut être lancé sur un environnement unsafe ou une propagation peut passer inaperçue.

## 3. Objectifs
- Rendre visibles isolation déclarée, limites réseau, durée, ressources, interactions et restrictions de fichiers.
- Permettre interruption, arrêt d’urgence, nettoyage, reset et blocage d’un nouveau Run.
- Signaler un environnement suspect à Platform Settings avec trace complète.
- Distinguer sécurité de l’analyse et containment d’un Endpoint réel.

## 4. Non-objectifs
Ne pas définir hyperviseur, isolation technique, instrumentation, procédure offensive, commande, API, protocole, infrastructure ou containment d’un Endpoint réel.

## 5. Propriétaire
Investigate possède l’évaluation fonctionnelle du Run; Platform Settings possède l’environnement, son health et son administration; Govern conserve l’autorité sur les cibles réelles.

## 6. Utilisateurs
Dynamic Analysis Operator; Malware Analyst; Investigation Lead; Platform Administrator en réception de handoff.

## 7. Conditions d’entrée
Sandbox Run et Environment identifiables; policies, limites et health visibles; permission de stop/blocage/signalement réévaluée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Isolation et politiques | Platform Settings | attentes déclarées | oui | version du Run | policy-unknown/block |
| Limites réseau/durée/ressources | Intake/profile/environment | bornes fonctionnelles | oui | snapshot Run | Run non autorisé |
| Signaux de sécurité | CAP-INV-317/321/Tools | erreur, sortie de limite ou suspicion | non avant Run | temps réel ou retard visible | aucun signal déclaré |
| État cleanup/reset | Platform Settings/Run | récupération | oui après Run | état courant | nouveau Run bloqué |
| Permission d’arrêt/blocage | Security | autorité fonctionnelle | oui | réévaluée | escalation requise |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Run | Investigate concept | phase, erreurs, partials et stop status | consulter/interrompre selon permission |
| Sandbox Environment | Platform Settings | health, isolation déclarée, reset et restrictions | consulter |
| Network Observation | Investigate concept | sortie réseau inattendue | consulter |
| Policy / Audit | Settings/Govern/Shared | limites, exceptions et trace | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Safety assessment | créer/versionner/contester | Investigate concept | limites et incertitude visibles |
| Run safety event | émettre | Investigate semantics/Shared Trace | Run/environment/source obligatoire |
| Environment alert | préparer/transmettre | Platform Settings | aucune mutation locale |
| Environment use block | demander/appliquer selon owner | Settings | Investigate ne réactive pas seul |

## 11. Fonctionnalités
Afficher limites et statut de sécurité; empêcher propagation et accès non autorisé au niveau fonctionnel; interrompre; arrêt d’urgence; nettoyage; reset; marquer unsafe; bloquer nouveau Run; alerter Settings; conserver trace.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire safety status | Operator | Safety assessment | 0 | Run/environment accessibles | statut visible | non |
| Interrompre/arrêt d’urgence | Operator | Sandbox Run | 2 | signal/permission | stopping/unsafe | OPEN-013 |
| Bloquer nouveau Run | Operator/Policy | Environment use | 2 | unsafe/reset incomplete | blocked | OPEN-013 |
| Signaler à Settings | Operator | Environment alert | 2 | contexte/provenance | handoff attribué | OPEN-013 |
| Demander reprise après reset | Authorized reviewer | Environment selection | 2 | health réévalué | nouvelle sélection possible | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Vérifier limites | oui | policies/seuils explicables | oui | explication | statuts et règles |
| Détecter mismatch | oui | comparaison policy/observation | oui | résumé | alertes déterministes |
| Proposer stop | oui | règles | oui | suggestion | contrôle manuel |
| Réactiver environnement | owner autorisé | health/policy | workflow de revue | jamais autonome | Settings review |

## 14. États fonctionnels
`safe-declared`, `degraded`, `policy-blocked`, `stopping`, `cleanup-pending`, `resetting`, `unsafe-environment`, `quarantined-by-settings`, `review-required`.

## 15. États d’interface
Loading conserve Run/environment; Partial montre signaux manquants; Error garde les données valides; Offline interdit nouveau Run; Permission denied ne fuit rien; unsafe est prioritaire et persistant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Safety assessment | Analysis Result | operator/Case | limites et incertitude |
| Safety event | business event | Trace/Activity | Run/environment/source visibles |
| Environment alert | Settings handoff | Platform Settings | aucune administration locale |
| Run stop context | Run event | CAP-INV-317 | partials et raison conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-317/321 | signal de sécurité | CAP-INV-326 | Run, environment, policy, observation, status | Run/network view |
| CAP-INV-326 | unsafe environment | Platform Settings alert | environment/version, symptômes, Runs, stop/reset status | safety view |
| Platform Settings | health/reset update | CAP-INV-315/326 | environment status, restrictions, review outcome | Settings |
| CAP-INV-326 | recovery autorisée | CAP-INV-317 | nouveau Run, environment version, policy | safety view |

## 18. Dépendances
CAP-INV-315/317/321/327; Settings Sandbox/Health/Policies/Audit; Shared Trace/Notifications; OPEN-005/013.

## 19. Source de vérité
Settings reste source du health et de l’administration de l’environnement; Investigate reste source du contexte de sécurité du Run; safety containment ≠ Endpoint containment.

## 20. Provenance et audit
Run, environment/version, policies, limites, signal source, Tool/version, initiateur, stop, partials, cleanup, reset, blocage, Settings handoff et review outcome.

## 21. Permissions fonctionnelles
Safety status read; Sandbox Run emergency stop; environment use block/request; Settings alert create; restricted environment use; cleanup/reset status read; cross-tenant denied.

## 22. Limites et erreurs
Policy stale; signal tardif; stop incomplet; cleanup/reset échoué; environnement inaccessible; health contradictoire; alerte non remise; permission révoquée.

## 23. Métriques
Runs stopped; unsafe environments; cleanup/reset failures; blocks; Settings alerts; recovery reviews; repeated unsafe use cible zéro.

## 24. Classification de livraison
`defined` / `planned`; aucune technologie d’isolation, moteur, hyperviseur, API, protocole, commande ou implémentation.

## 25. Critères d’acceptation
**Given** un Run qui dépasse une limite déclarée
**When** le signal est reçu
**Then** la limite, la source, le Run et l’environnement sont visibles, le Run peut être arrêté et un nouveau Run est bloqué jusqu’à revue

**Given** un reset incomplet
**When** un nouveau Run est demandé
**Then** il est bloqué et Settings reçoit le contexte

**Given** aucun modèle IA
**When** la sécurité est supervisée
**Then** policies, statuts, alertes et actions humaines couvrent le workflow

## 26. Questions ouvertes
Safety assessment reste un concept; OPEN-005 et OPEN-013 restent ouvertes; aucune sandbox n’est déclarée invulnérable.

## 27. Consommateurs documentaires
Dynamic Sandbox, Run Management, Network Analysis, Settings Sandbox/Health/Audit, Trace, Permissions et Journeys.
