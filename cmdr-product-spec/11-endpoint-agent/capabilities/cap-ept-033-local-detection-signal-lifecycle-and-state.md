---
id: CAP-EPT-033
title: Local Detection Signal Lifecycle and State
product: endpoint-agent
module: detection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-033 — Local Detection Signal Lifecycle and State

## 1. Définition
Définir le concept technique Endpoint de **Local Detection Signal Candidate** dérivé d’un ou plusieurs Detection Match, ses états documentaires et son handoff, sans redéfinir l’objet canonical `Signal` qui reste Command-owned.

## 2. Problème utilisateur
Les matches locaux peuvent se répéter, évoluer ou être transmis. Sans contrat explicite, un état technique local pourrait être confondu avec le canonical Signal, un Alert, un Incident ou une conclusion Investigate.

## 3. Objectifs
Représenter création, activité, répétition, mise à jour, groupement, suppression-state projection, expiration, invalidation et forwarding avec références exactes et provenance.

## 4. Non-objectifs
Aucun canonical Command Signal/Alert/Incident, aucune suppression destructive, aucune réponse, aucune machine d’état technique finale, aucun Case/Finding/Evidence.

## 5. Propriétaire
Endpoint Agent possède uniquement le Local Detection Signal Candidate technique. Command conserve le canonical `signal` et sa promotion opérationnelle ; Investigate conserve la qualification analytique.

## 6. Utilisateurs
Endpoint Operator ; SOC/Investigate Analyst ; Command consumer ; Detection Engineer ; Security Reviewer ; Auditor.

## 7. Conditions d’entrée
Au moins un Local Detection Evaluation/Match sourcé, tenant/Agent connus, content/version et observation refs résolubles, permissions de lecture, états de suppression éventuels sourcés.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Detection Match | CAP-EPT-032 | résultat technique | oui | evaluation time | pas de signal candidate |
| content/version | CAP-EPT-031/032 | référence | oui | version liée | état invalid/unknown |
| context refs | CAP-EPT-017..024 | observations | non | source freshness | contexte partiel |
| suppression-state projection | Endpoint/Settings/Investigate source si applicable | état externe | non | source version | aucune suppression inférée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Detection Evaluation/Match | Endpoint Agent | outcome + refs | read |
| Command Signal | Command | référence destination éventuelle | aucune mutation locale |
| Detection Content | Investigate | id/version | read |
| endpoint-agent/local-audit-event | Endpoint Agent | source/audit context | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Detection Signal Candidate | créer/mettre à jour état local | Endpoint Agent | non-canonical Command Signal |
| Recurrence/Grouping Reference | dériver | Endpoint Agent | préserve tous les matches sources |
| Forwarding Reference | ajouter | Endpoint Agent | destination owner inchangé |

## 11. Fonctionnalités
Créer une représentation locale depuis un match, lier les répétitions et mises à jour, afficher une suppression-state seulement si sourcée, invalider/expirer conceptuellement selon contenu/version et transmettre une projection permission-aware.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect local candidate | Operator/Analyst | local candidate | 0 | read | état + refs visibles | non |
| correlate repeat/update | deterministic service | candidate | 1 | matching/grouping inputs | recurrence liée | non |
| acknowledge local workflow state | authorized operator | local candidate | 2 | source autorise état local | état workflow seulement | non |
| request suppression mutation | user | foreign suppression policy | 2 request only | owner/source applicable | aucune mutation locale implicite | selon owner |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| créer depuis match | oui | oui | oui | non nécessaire | deterministic mapping |
| grouper répétitions | oui | oui | oui | suggestion possible | keys/time rules |
| expliquer candidate | oui | oui | oui | oui | source/rationale fields |
| masquer/supprimer provenance | non | non | non | interdit | provenance retained |

## 14. États fonctionnels
`created`, `active`, `repeated`, `updated`, `grouped`, `suppressed-projection`, `expired`, `invalidated`, `forwarded`. États documentaires non exhaustifs ; aucune machine runtime finale n’est sélectionnée.

## 15. États d’interface
Aucun Screen ID. Suppressed doit rester visible comme état autorisé, jamais supprimé ; Partial/Stale/Permission denied sont explicites.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| local signal candidate | concept Endpoint | CAP-EPT-034..046 | sources/matches conservés |
| recurrence/grouping refs | relation technique | Operator/Investigate | dedup ≠ erase |
| forwarded projection | handoff | Command/Investigate | destination qualifie selon son owner |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-032 Match | match créé | local candidate | evaluation/content/obs refs | provenance intacte |
| candidate | repeat/context change | updated/grouped state | old/new refs/reasons | historique conservé |
| candidate | authorized handoff | Command/Investigate consumer | stable refs + uncertainty | aucune création auto Finding/Signal |

## 18. Dépendances
CAP-EPT-031/032 ; detection-suppression ; local-correlation ; Command canonical Signal ; Investigate Detection Engineering ; `OPEN-008` ; `OPEN-017`.

## 19. Source de vérité
Endpoint est source du candidate technique local et de ses refs. Command est source de vérité du canonical Signal. Investigate est source des Finding/Evidence/Case et du Detection Content lifecycle.

## 20. Provenance et audit
Conserver chaque Match source, grouping key conceptuel, recurrence, suppression projection source/owner/reason/expiry, state transitions, forwarding destination et correlation id. Suppressed ≠ deleted.

## 21. Permissions fonctionnelles
Lecture local candidate, rationale/context restreints, provenance et handoff ; interaction suppression seulement si source owner l’autorise. Cross-tenant denied. Pas de permission de canonical Signal/Finding/Response implicite.

## 22. Limites et erreurs
Candidate peut rester partial/stale/invalidated ; grouped ≠ identical ; deduplicated ≠ erased ; multiple matches ≠ incident ; signal candidate ≠ canonical Signal, Finding, Evidence, Case, Result ou response authority.

## 23. Métriques
Candidates créés, repeats/groupings, suppression projections, invalidations, handoff completeness, stale/partial reasons. Aucune cible opérationnelle imposée.

## 24. Classification de livraison
`draft / defined / planned` ; documentation fonctionnelle uniquement, pas de runtime/queue/engine implémenté.

## 25. Critères d’acceptation
**Given** deux matches équivalents se répètent, **When** le grouping local s’applique, **Then** leurs références restent toutes conservées et grouped ne signifie pas identical.

**Given** une suppression valide est projetée, **When** le candidate est consulté, **Then** il reste visible avec owner/reason/expiry/provenance et n’est pas supprimé.

**Given** un candidate est transmis à Investigate, **When** le handoff est consommé, **Then** aucun Finding/Evidence/Case n’est créé automatiquement.

## 26. Questions ouvertes
`OPEN-017` runtime/language reste ouverte ; `OPEN-008` platform/source reste ouverte. Le bridge exact entre local candidate et canonical Command Signal est une projection/handoff, pas un transfert d’ownership.

## 27. Consommateurs documentaires
EPT-3 context/grouping/investigation ; Command ; Investigate ; Detection Engineering ; Security/Trust ; Quality ; user journeys futurs ; EPT-4+ non démarrés.