---
id: CAP-EPT-040
title: User, Session and Authentication Context Investigation
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-040 — User, Session and Authentication Context Investigation

## 1. Définition
Définir le contexte local user/session/authentication à partir des observations disponibles : session refs, process activity, authentication facts, privilege context, related detections, chronology et privacy restrictions.

## 2. Problème utilisateur
Une anomalie d’authentification ou un process lié à une session nécessite du contexte ; sans contrat, elle peut être interprétée abusivement comme compromission de compte ou exposer des secrets/identifiants sensibles.

## 3. Objectifs
Relier sessions interactive/service, logon context, process/detection refs et privilege context observé ; protéger identity metadata ; interdire secret exposure ; exposer gaps/platform limitations.

## 4. Non-objectifs
Aucun account compromise verdict, credential capture/reveal, session termination, identity-provider administration, password/token use, Evidence/Finding/Case automatique.

## 5. Propriétaire
Endpoint Agent possède le contexte user/session local. Platform Settings conserve identities/admin config ; Investigate conserve qualification analytique ; Govern conserve toute account-response authority.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Security/Privacy Reviewer ; Auditor ; Detection Engineer ; Command consumer autorisé.

## 7. Conditions d’entrée
Authentication/session observations supportées ; user/session refs permissionnées ; tenant scope ; raw secret absent ; related process/detection refs optionnelles ; freshness/limitations connues.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| user/session/auth observations | CAP-EPT-020 | faits locaux | oui | source freshness | context unavailable |
| process context | CAP-EPT-017/037 | relations | non | snapshot time | process relation unknown |
| local detection refs | CAP-EPT-032..035 | context | non | candidate freshness | session context only |
| tenant/identity restrictions | Platform Settings/Security | access context | oui | current policy | deny/restrict |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| User Session/Auth Observation | Endpoint Agent | session/logon/privilege refs | read |
| Process Investigation Context | Endpoint Agent | process relation | read |
| principal/tenant projections | Platform Settings | identity/scope refs | read according permission |
| Finding/Case | Investigate | destination refs only | no creation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| User Session Investigation Context | dériver/rafraîchir | Endpoint Agent | no raw secret |
| Session-Process Relation | dériver | Endpoint Agent | source-backed |
| Authentication Chronology Projection | dériver | Endpoint Agent | anomaly ≠ compromise |

## 11. Fonctionnalités
Construire session/auth context, relier process et related detections, afficher privilege facts s’ils existent, ordonner observations, masquer sensitive identity fields, et conserver unsupported/platform-specific limitations.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect session/auth context | Analyst | User Session Context | 0 | read | facts permissionnés | non |
| correlate existing session refs | deterministic service | context | 1 | observations existantes | liens sourcés | non |
| request bounded context refresh | Operator | local session state | 2 | sans credential acquisition/termination | freshness update | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| lier session/process | oui | oui | oui | suggestion possible | stable refs/time |
| construire chronology | oui | oui | oui | non nécessaire | timestamps |
| résumer auth context | oui | oui | oui | oui | structured facts |
| conclure account compromise | non | non | non | interdit | Investigate qualification |

## 14. États fonctionnels
`available`, `partial`, `restricted`, `identity-ambiguous`, `stale`, `unsupported`, `unavailable`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Restricted identity/privilege fields et source limitations sont visibles comme tels ; aucune UI ne doit exposer raw credential material.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| User Session Investigation Context | concept Endpoint | CAP-EPT-042..045 | privacy + source refs |
| process/auth chronology | projection Endpoint | Analyst | anomaly ≠ compromise |
| restricted/missing context | diagnostic | Analyst | no secret inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| auth/session observation/candidate | pivot | User Session Context | user/session/process refs | origin retained |
| User Session Context | process/timeline pivot | CAP-EPT-037/042 | stable refs/time | no response |
| unavailable sensitive detail | expansion request | EPT-3 boundary | restriction reason | no credential collection |

## 18. Dépendances
CAP-EPT-017/020/026/032..037 ; `investigation/user-sessions.md`; authentication telemetry ; Platform Settings identity/tenant boundaries ; Security privacy ; `OPEN-008`.

## 19. Source de vérité
Endpoint est SOT du contexte local dérivé. Settings reste SOT des identities/admin scopes ; Investigate reste SOT des analyst conclusions ; aucun credential secret n’est créé/possédé ici.

## 20. Provenance et audit
Session/auth observation ids, user/principal refs, process refs, privilege context source, masking state, related candidates, chronology time, tenant et actor/service sont conservés.

## 21. Permissions fonctionnelles
User/session context read, sensitive identity/privilege/command metadata read, provenance, cross-tenant deny. Aucun droit de credential reveal/use, session termination ou account containment.

## 22. Limites et erreurs
Authentication anomaly ≠ account compromise ; identity association peut être ambiguous ; no raw secrets ; source/platform gaps restent explicites ; missing session does not imply no activity.

## 23. Métriques
Context completeness, identity ambiguity, restricted fields, process-link coverage, stale/unsupported reasons, pivot success.

## 24. Classification de livraison
`draft / defined / planned`; aucune identity-response, credential engine ou auth provider integration implémentée.

## 25. Critères d’acceptation
**Given** une authentication anomaly est liée à un signal, **When** le contexte est construit, **Then** les facts/session/process refs sont montrés sans déclarer compte compromis.

**Given** l’utilisateur manque le droit sur un command/identity field sensible, **When** le contexte est consulté, **Then** le champ est masqué et la restriction est explicite.

**Given** la plateforme ne fournit pas cette télémétrie, **When** le contexte est demandé, **Then** état = unsupported/unavailable sans support universel inféré.

## 26. Questions ouvertes
`OPEN-008` reste ouverte. Les actions de compte/session appartiennent à de futurs flux Govern/EPT-5 et ne sont pas définies ici.

## 27. Consommateurs documentaires
EPT-3 timeline/pivots/summary ; Investigate Case/Finding workflows ; Settings/Security ; Command ; Quality ; future EPT-5 boundary.