---
id: CAP-EPT-037
title: Process Tree and Execution Context Investigation
product: endpoint-agent
module: investigation
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-037 — Process Tree and Execution Context Investigation

## 1. Définition
Définir l’investigation locale read-only du contexte process/exécution à partir des observations déjà disponibles : parent, enfants, ancestry, executable refs, user/session, modules, file/network relations, chronology et related local detection candidates.

## 2. Problème utilisateur
Un match process isolé n’explique pas son contexte d’exécution. L’analyste doit pouvoir reconstruire les relations observées sans lancer une collecte supplémentaire ni transformer une ancestry en preuve causale d’attaque.

## 3. Objectifs
Construire un Process Context permission-aware ; préserver snapshot time et partial ancestry ; lier EPT-2 observations et EPT-3 candidates ; exposer hashes/signatures uniquement s’ils sont déjà observés ; rendre les limites visibles.

## 4. Non-objectifs
Aucune acquisition process/memory/file, aucun remote shell, aucune command execution, aucun process kill, aucun causal attack-chain verdict, aucun Finding/Evidence/Case automatique.

## 5. Propriétaire
Endpoint Agent possède le contexte process local dérivé. Investigate conserve Case/Evidence/Finding et le raisonnement central ; Shared conserve generic Linking/Timeline/Search.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Security/Privacy Reviewer ; Auditor ; Command consumer autorisé.

## 7. Conditions d’entrée
Process observations EPT-2 disponibles ; Endpoint/tenant connus ; permissions pour metadata/command context ; snapshot/freshness explicites ; aucun fetch qui deviendrait Collection.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| process/execution observations | CAP-EPT-017 | faits locaux | oui | observation freshness | context unavailable/partial |
| related detection candidates | CAP-EPT-032..034 | refs de détection | non | candidate freshness | contexte sans détection |
| file/network/user/module refs | CAP-EPT-018..021 | relations observées | non | source freshness | relation omitted/unknown |
| process-tree snapshot semantics | `investigation/process-tree.md` | source fonctionnelle | oui | snapshot time | ancestry partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Process Observation | Endpoint Agent | pid/ref, parent/child, executable, time | read |
| Local Detection Signal Candidate | Endpoint Agent | related refs | read |
| telemetry-event | Shared | event reference si projetée | read |
| Case/Evidence/Finding | Investigate | destination refs seulement | aucune création locale |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Process Investigation Context | dériver/rafraîchir | Endpoint Agent | données déjà disponibles uniquement |
| Execution Ancestry Relation | dériver | Endpoint Agent | relation ≠ causalité |
| Process Context Snapshot | créer projection | Endpoint Agent | snapshot time/freshness obligatoire |

## 11. Fonctionnalités
Résoudre parent/children/ancestry connus, lier executable/hash/signature refs observées, user/session, modules, files/network et related detections ; conserver chronology, missing ancestors et provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect process context | Analyst/Operator | Process Context | 0 | read | relations visibles | non |
| correlate existing refs | deterministic service | context | 1 | observations locales | context enrichi | non |
| request bounded contextual refresh | authorized operator | existing local state | 2 | refresh sans acquisition | freshness update seulement | non |

Si des données non disponibles exigent acquisition, l’action s’arrête à la frontière EPT-3 et pointe vers futur EPT-4.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| construire ancestry | oui | oui | oui | non nécessaire | parent/child refs |
| lier observations | oui | oui | oui | suggestion possible | stable keys/time refs |
| résumer process context | oui | oui | oui | oui | structured context |
| affirmer causalité/malicious | non | non | non | interdit | analyst qualification externe |

## 14. États fonctionnels
`available`, `partial`, `stale`, `restricted`, `ancestry-incomplete`, `unsupported`, `unavailable`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Une future surface montre snapshot time, partial ancestry, masked command metadata, unavailable relations et Permission denied.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Process Investigation Context | concept Endpoint | CAP-EPT-042..045 | observations sources préservées |
| ancestry/descendant relations | relations Endpoint | Analyst | relation ≠ causal proof |
| missing-context list | diagnostic | Analyst/EPT-4 future | aucune acquisition déclenchée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| process observation/candidate | inspect/pivot | Process Context | process + detection refs | origin retained |
| Process Context | pivot file/network/session | CAP-EPT-038/039/040 | stable refs/time | no collection |
| missing local data | context expansion request | EPT-3 boundary | missing requirement | futur EPT-4 seulement |

## 18. Dépendances
CAP-EPT-017..021/032..034 ; `investigation/process-tree.md`; `host-inspection.md`; Shared Linking/Timeline ; Investigate Case boundary ; `OPEN-008`, `OPEN-017`.

## 19. Source de vérité
Endpoint est SOT du contexte process local dérivé des observations disponibles. Les objets Case/Evidence/Finding restent Investigate ; les generic links restent Shared.

## 20. Provenance et audit
Snapshot time, process ids/refs, parent-child sources, source timestamps, hashes/signatures si observés, command metadata masking, related candidates, missing ancestors et actor/service sont conservés.

## 21. Permissions fonctionnelles
Process context read, sensitive command metadata read, user/session relation read, provenance read, cross-tenant deny. Aucun droit d’acquisition, execution ou mutation.

## 22. Limites et erreurs
Process ancestry ≠ causal attack chain proof ; process relation ≠ malicious relation ; données absentes ne sont pas collectées ; PID reuse/time skew/stale snapshot doivent être explicités.

## 23. Métriques
Context completeness, ancestry depth known/unknown, stale/restricted fields, pivot resolution rate, missing-data reasons ; aucune cible SLO finale.

## 24. Classification de livraison
`draft / defined / planned`; aucun inspector runtime, collector ou remote execution implémenté.

## 25. Critères d’acceptation
**Given** un process signalé possède parent/enfants observés, **When** l’analyste pivote, **Then** la tree est reconstruite avec snapshot time et refs sans déclaration causale.

**Given** ancestry est incomplète, **When** le contexte est affiché, **Then** le gap est explicite et aucune ancestry inventée n’est ajoutée.

**Given** command metadata est restreinte, **When** un utilisateur non autorisé consulte le contexte, **Then** le champ est masqué sans cacher l’existence du process.

## 26. Questions ouvertes
`OPEN-008` conserve disponibilité par plateforme/source ; `OPEN-017` reste Detection runtime only. Les besoins nécessitant acquisition passent à EPT-4 futur.

## 27. Consommateurs documentaires
EPT-3 timeline/pivots/summary ; Investigate Case workflows ; Command ; Security ; Quality ; future EPT-4 comme boundary seulement.