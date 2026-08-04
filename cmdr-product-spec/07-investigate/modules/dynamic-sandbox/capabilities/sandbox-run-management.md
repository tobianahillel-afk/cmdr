---
id: CAP-INV-317
title: Sandbox Run Management
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
  - REQ-PROD-020
  - REQ-SEC-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-317 — Sandbox Run Management

## 1. Définition
Préparer, valider, mettre en queue, lancer explicitement, suivre, interrompre, annuler, collecter, nettoyer, réinitialiser, relancer et comparer un Sandbox Run sans confondre ses étapes ni son résultat avec un succès analytique.

## 2. Problème utilisateur
Queued, started, Artifact launched, completed et environment reset sont des événements différents. Les fusionner masque échecs, timeouts, résultats partiels et risques de nettoyage.

## 3. Objectifs
- Conserver chaque phase du Run distincte et visible.
- Exiger un lancement explicite.
- Exposer progression, durée, interruptions, erreurs, résultats partiels, nettoyage et reset.
- Relancer dans un nouveau Run lié plutôt que réécrire l’historique.

## 4. Non-objectifs
Ne pas définir infrastructure, virtualisation, instrumentation, API, protocole, commande ou moteur; ne pas agir sur un Endpoint réel; ne pas présenter completion comme comportement significatif.

## 5. Propriétaire
Investigate possède le contexte du Sandbox Run. Settings possède l’environnement; Studio possède Tools/Tool Calls/Automation Runs; Shared possède Background Jobs; Govern reste hors du Run isolé sauf action réelle.

## 6. Utilisateurs
Dynamic Analysis Operator; Malware Analyst; Case Analyst; Reviewer.

## 7. Conditions d’entrée
Session, Artifact, environnement, profil, limites, initiateur, permissions et policies identifiables; environnement health compatible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Session et Artifact | CAP-INV-316 | contexte/source | oui | versions référencées | validating bloqué |
| Environnement | CAP-INV-315/Settings | cible isolée | oui | health réévalué | unavailable |
| Profil et limites | CAP-INV-323/intake | scénario, durée, réseau, ressources | oui | version du Run | rester draft |
| Tools et versions | Studio | exécution/observation | oui selon Run | status visible | incompatible |
| Confirmation de lancement | utilisateur/policy | autorisation explicite | oui | au démarrage | ne pas lancer |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Dynamic Analysis Session | Investigate | objectif et Runs liés | consulter/lier |
| Artifact | Investigate | version et restrictions | consulter/exécuter en sandbox selon permission |
| Sandbox Environment | Settings | health, version, limits, reset | consulter |
| Tool/Tool Call/Automation Run | Studio | producer/version/status | consulter/invoquer selon permission |
| Background Job | Shared | queue/progress/cancel | consulter/demander cancel |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Sandbox Run | préparer/start/stop/cancel/supersede | Investigate concept | distinct de Automation/Response Run |
| Run phase event | émettre | Investigate semantics/Shared trace | phases distinctes |
| Result/observation relations | créer | Investigate | source Run obligatoire |
| Environment safety notice | émettre | Settings handoff | pas d’administration locale |

## 11. Fonctionnalités
Préparer; valider; queue; preparing; starting; running; idle; stopping; collecting-results; cleaning; resetting; partial/completed/failed/cancelled/timed-out/unsafe; interrupt, retry, compare.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Préparer Run | Operator | Sandbox Run | 2 | préconditions valides | draft validé | OPEN-013 |
| Lancer explicitement | Operator | Sandbox Run | 1 | confirmation/policy | queued puis preparing | non |
| Interrompre/arrêter | Operator | Sandbox Run | 2 | Run actif | stopping/partials | OPEN-013 |
| Annuler avant lancement | Operator | Sandbox Run | 2 | Artifact non lancé | cancelled | OPEN-013 |
| Relancer | Operator | New Sandbox Run | 2 | justification | nouveau Run lié | OPEN-013 |
| Comparer | Reviewer | Sandbox Runs | 0 | Runs accessibles | CAP-INV-325 | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Valider préconditions | oui | règles/policies | oui | explication | checklist |
| Afficher progression | non | Background Job | oui | résumé | statuts bruts |
| Proposer stop/retry | oui | seuils explicables | oui | suggestion | contrôles manuels |
| Lancer le Run | humain explicite | jamais silencieux | workflow avec gate | jamais autonome | action humaine |

## 14. États fonctionnels
`validating`, `queued`, `preparing`, `starting`, `running`, `idle`, `stopping`, `collecting-results`, `cleaning`, `resetting`, `partial`, `completed`, `failed`, `cancelled`, `timed-out`, `unsafe-environment`. Machines finales reportées.

## 15. États d’interface
Loading conserve Run/session; Partial distingue éléments collectés; Error garde les résultats valides; Offline interdit lancement; Permission denied ne fuit rien; unsafe bloque un nouveau Run.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Sandbox Run record | Run concept | session/Case | phases, statuts et timestamps distincts |
| Run progress | Background Job projection | operator | progression et partials |
| Run observations | Analysis Results | CAP-INV-318..324 | source Run obligatoire |
| Cleanup/reset status | environment event | Settings/operator | fin du Run distincte du reset |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-316 | préparer | CAP-INV-317 | session, Artifact, environnement, profil, limites, Tools | session |
| CAP-INV-317 | Artifact lancé | CAP-INV-318..322 | Run, timestamps, sources | Run |
| CAP-INV-317 | output produit | CAP-INV-324 | Run, process/event, capture, restrictions | Run |
| CAP-INV-317 | échec/timeout | Recovery ou nouveau Run | partials, errors, environment status, reason | Run |
| CAP-INV-317 | unsafe | CAP-INV-326/Settings | Run, environment, symptômes, stop status | Run |

## 18. Dépendances
CAP-INV-315/316/318..327; Shared Jobs/Notifications/Trace; Settings Environments/Health; Studio Tools; OPEN-013/015.

## 19. Source de vérité
Investigate est source des phases métier du Run; Shared est source du mécanisme Job; Settings est source de l’environnement; Studio est source des Tool Calls.

## 20. Provenance et audit
Session, Artifact/version, environnement/version, profil, limites, initiateur, confirmation, Tools, phases, timestamps, outputs, erreurs, stop, nettoyage, reset et dispositions.

## 21. Permissions fonctionnelles
Run prepare/start/stop/cancel/retry; sensitive Artifact execute-in-sandbox; results read; Background Job cancel; restricted environment use; cross-tenant denied.

## 22. Limites et erreurs
Queue bloquée; environment reset/unavailable; Tool incompatible; Artifact launch failure; timeout/crash; résultats tronqués; stop incomplet; cleanup/reset failure; unsafe environment.

## 23. Métriques
Runs par phase/état; temps queue/start/run/cleanup/reset; interruptions/timeouts; partials; unsafe environments; retries.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Run prêt et explicitement confirmé
**When** il est lancé
**Then** queued, preparing, starting, Artifact launched, running, completed, cleaning et reset restent des étapes distinctes

**Given** un timeout
**When** le Run s’arrête
**Then** les résultats partiels sont conservés et le timeout n’est pas présenté comme absence de comportement

**Given** aucun modèle IA
**When** le Run est géré
**Then** statuts, contrôles, règles et actions humaines couvrent le workflow essentiel

## 26. Questions ouvertes
Sandbox Run reste un concept; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Session, Behavioral Timeline, Process/File/Network analysis, Runtime Artifacts, Safety, Comparison, Settings et Studio.
