---
id: CAP-CMD-204
title: Business Impact Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-003, REQ-PROD-005, REQ-PROD-013, REQ-PROD-021]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-204 — Business Impact Context

## 1. Définition
Capture ou projette l’impact métier d’un Incident avec type, portée, source, owner, certitude et date, puis le transmet à Govern.

## 2. Problème utilisateur
Un score ou texte sans auteur influence priorité et autorité. Sans statut de certitude, une hypothèse est prise pour un impact confirmé.

## 3. Objectifs
Distinguer assumed/confirmed/disputed/unknown, relier Services/owner/type, conserver source/date/justification et transmettre sans score opaque.

## 4. Non-objectifs
Ne calcule pas un impact financier définitif, ne modifie pas le catalogue, ne remplace pas Business Owner et ne confirme pas une proposition IA.

## 5. Propriétaire
Command possède l’impact porté par l’Incident ; Service, Finding et Result restent leurs sources propriétaires.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : Business Owner, Service Owner, SOC Analyst L2.

## 7. Conditions d’entrée
Incident accessible, Service confirmé ou absence documentée, acteur/source et permission de modification.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Impact statement | humain ou source métier | type, portée et certitude | oui ou état `unknown` | daté à chaque modification | follow-up créé, aucun impact confirmé |
| Service context | CAP-CMD-201 | Service, owner et criticité | non | fraîcheur du catalogue | impact conservé sans Service confirmé |
| Technical context | Incident, Case, Finding ou Result | faits et projections | non | source visible | certitude réduite ou `partial` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | impact, owner et contexte | consulter et modifier si autorisé |
| Service conceptuel | Shared Business Service Catalog | owner et criticité | consulter et relier |
| Finding / Result | Investigate / Govern | résumé autorisé et vérification | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | ajouter ou modifier impact, certitude, source et owner métier | Command | classe 2, before/after audité |
| Task | demander confirmation ou données manquantes | Command | classe 2, résultat attendu explicite |
| Action Request context | transmettre l’impact | Govern | aucune Decision créée par Command |

## 11. Fonctionnalités
Saisir impact, confirmer/contester hypothèse, lier Services/owner, afficher source/fraîcheur et préparer le contexte Govern.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Saisir assumed | coordinateur | Incident | 2 | source et raison | impact assumed | OPEN-013 |
| Confirmer | Business Owner autorisé | Incident | 2 | données suffisantes | impact confirmed | OPEN-013 |
| Contester | rôle autorisé | Incident | 2 | motif | impact disputed | OPEN-013 |
| Transmettre | requester | Action Request context | 2 | action proposée | package Govern | oui |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Agréger le contexte | oui | oui | oui | résumé attribué | lecture Incident/Service/Finding/Result |
| Proposer un impact | oui | règles possibles | oui | oui, statut assumed seulement | saisie humaine sourcée |
| Confirmer ou contester | oui | validation de permission | workflow possible | jamais décisionnelle | action Business Owner/coordinateur |
| Préparer Govern | oui | complétude/validation | oui | brouillon attribué | formulaire Action Request manuel |

## 14. États fonctionnels
`unknown`, `assumed`, `confirmed`, `disputed`, `outdated`, `superseded`.

## 15. États d’interface
Certitude/source toujours visibles ; Stale distinct d’unknown ; Offline bloque mutation ; Permission denied masque seulement les projections protégées.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Business impact context | projection Incident | Priority, Mission Control et Govern | source, certitude, owner et date visibles |
| Impact change | événement Incident | Audit et Work Queue | before/after, acteur et justification |
| Confirmation follow-up | Task Command | Business/Service Owner | attribuée et liée à l’Incident |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Business Impact | confirmation requise | Task Coordination / Business Owner | Incident, impact assumed et données manquantes | Incident Detail restauré |
| Business Impact | action nécessitant autorité | Govern | Incident, Services, impact, urgence, source et action | Incident Detail restauré |
| Business Impact | source technique à vérifier | Investigate | Incident, hypothèse d’impact et question | Incident Detail restauré |

## 18. Dépendances
CAP-CMD-201, CAP-CMD-205, CAP-CMD-106 et Business Service Catalog.

## 19. Source de vérité
Impact Incident : Command. Service/Finding/Result : propriétaires sources.

## 20. Provenance et audit
Auteur/producteur, certitude, source, before/after, time et correlation ID.

## 21. Permissions fonctionnelles
Incident read/manage, future permission de confirmation métier et Govern Action Request create ; atomisation reportée.

## 22. Limites et erreurs
Service inconnu, impact stale, owner absent, sources conflictuelles ou refus ne produisent jamais `confirmed`.

## 23. Métriques
Incidents par certitude, délai assumed→confirmed/disputed et Action Requests avec impact sourcé ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire.

## 25. Critères d’acceptation
**Given** un impact proposé par IA, **When** il est affiché, **Then** il reste assumed et le producteur est visible.

**Given** Business Owner confirme, **When** l’action est autorisée, **Then** before/after, auteur et date sont audités.

**Given** aucun modèle, **When** l’impact est saisi, **Then** la voie humaine complète reste disponible.

## 26. Questions ouvertes
Permission de confirmation métier et types d’impact canoniques ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Incident Detail, Mission Control, Priority, transition Govern, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.