---
id: CAP-CMD-205
title: Risk Prioritization Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-205 — Risk Prioritization Context

## 1. Définition
Combine explicitement impact, urgence, criticité Service, Exposure, Coverage, confiance, dépendances et SLA pour éclairer une priorité sans imposer un moteur définitif.

## 2. Problème utilisateur
Un score unique masque facteurs, inconnues et responsabilités. Sans contexte explicable, une recommandation devient de fait une priorité non contestable.

## 3. Objectifs
Afficher facteurs/sources, indiquer missing/conflicts, produire une proposition distincte de la priority et permettre accept/reject humain audité.

## 4. Non-objectifs
Ne choisit pas l’algorithme final, ne crée pas de score universel, ne modifie pas les sources et ne transforme pas une proposition en Decision.

## 5. Propriétaire
Command possède le contexte et la proposition ; les sources gardent leurs facteurs ; CAP-CMD-002 possède la mutation effective de priorité.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Business Owner, Risk stakeholder.

## 7. Conditions d’entrée
Work item, au moins un facteur, définitions/fraîcheur par facteur et permissions sources.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Impact, urgency et SLA | Command | facteurs opérationnels | au moins un facteur | version courante | contexte `partial` |
| Service, Exposure et Coverage | sources propriétaires | projections de risque | non | fraîcheur propre à chaque source | facteur marqué absent, jamais inventé |
| Recommendation | règle, moteur, workflow ou agent | proposition avec facteurs | non | version/run visible | analyse humaine disponible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | impact, urgence, priority et SLA | consulter et comparer |
| Service / Exposure / Coverage | sources propriétaires | criticité, exposition et gaps | consulter en projection |
| Signal / Alert | source ou Command | confidence et severity | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Priority recommendation | créer, recalculer ou disposer | Command | proposition non effective, facteurs visibles |
| Incident / Task | modifier priority seulement via CAP-CMD-002 | Command | classe 2 et acceptation autorisée |
| Facteurs sources | aucune mutation | propriétaires sources | lecture seule |

## 11. Fonctionnalités
Afficher facteur/unité/source/fraîcheur, comparer action/inaction, calculer une proposition déterministe ou IA attribuée et renvoyer vers Priority Management.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter | lecteur | risk context | 0 | sources autorisées | facteurs visibles | non |
| Recalculer | coordinateur | context | 0 | inputs et version | proposition | non |
| Accepter ou rejeter | coordinateur | proposition | 2 | rationale visible | disposition | OPEN-013 |
| Ouvrir une source | lecteur | source product | 0 | permission | transition | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Collecter les facteurs | oui | oui | oui | résumé facultatif | consultation des sources |
| Calculer une proposition | oui | oui, facteurs visibles | oui | proposition attribuée | moteur déterministe ou analyse humaine |
| Expliquer les inconnues | oui | oui | oui | reformulation possible | liste des facteurs absents/conflictuels |
| Accepter/rejeter | oui | validation/version | workflow possible | jamais automatique | disposition humaine |

## 14. États fonctionnels
`complete`, `partial`, `conflicting`, `proposal-available`, `proposal-rejected`, `unknown`.

## 15. États d’interface
Chaque facteur affiche source/fraîcheur ; Partial/Conflict ne sont jamais masqués ; Offline interdit acceptation stale ; Permission denied masque le facteur protégé.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Risk context | projection multi-source | Priority, Mission Control et Govern | facteurs, sources et inconnues visibles |
| Priority proposal | proposition Command | coordinateur et CAP-CMD-002 | attribuée, versionnée et non effective |
| Proposal disposition | événement d’audit | producteur et analystes | accept/reject distinct de la priority |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Risk Context | accepter une proposition | Priority Management | work item, facteurs, proposition et rationale | Risk Context ou work item restauré |
| Risk Context | action à effet élevé | Govern | Incident, facteurs, impact, urgence et action proposée | aucune Decision locale |
| Risk Context | facteur à approfondir | produit source ou Investigate | facteur, source, work item et question | même contexte restauré |

## 18. Dépendances
CAP-CMD-002, CAP-CMD-104, CAP-CMD-105, CAP-CMD-201, CAP-CMD-202, CAP-CMD-203, CAP-CMD-204 et Metrics Engine.

## 19. Source de vérité
Chaque facteur reste source-owned ; la proposition est Command et ne devient pas effective sans disposition autorisée.

## 20. Provenance et audit
Producer type, version/run, facteurs, missing data, rationale et disposition.

## 21. Permissions fonctionnelles
Command read, coordinate pour disposition et source-specific reads ; `OPEN-013` reportée.

## 22. Limites et erreurs
Facteurs incompatibles/stale/unknown, source refusée ou moteur indisponible donnent Partial/Conflict, jamais score universel.

## 23. Métriques
Propositions avec facteurs/source/fraîcheur, dispositions par producer et contexts conflicting/unknown ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; moteur final non choisi.

## 25. Critères d’acceptation
**Given** des facteurs partiels, **When** une proposition est calculée, **Then** missing data et facteurs sont visibles et priority reste inchangée.

**Given** une proposition IA, **When** elle est consultée, **Then** producer/run/rationale sont visibles et acceptation humaine requise.

**Given** aucun modèle, **When** un recalcul est demandé, **Then** moteur déterministe ou analyse manuelle fournit le contexte.

## 26. Questions ouvertes
Quel moteur ultérieur et quels facteurs sont comparables entre tenants ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Priority Management, Mission Control, package Govern, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.