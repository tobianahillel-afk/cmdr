---
id: CAP-EPT-038
title: File and Filesystem Context Investigation
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
# CAP-EPT-038 — File and Filesystem Context Investigation

## 1. Définition
Définir une investigation locale du contexte fichier/filesystem à partir des observations déjà disponibles : path refs, metadata, hashes observés, process associations, chronology, related observations/signals et limitations.

## 2. Problème utilisateur
Une activité fichier liée à une détection doit pouvoir être contextualisée sans télécharger le fichier, inventer un verdict malware ou transformer un hash/reputation hit en preuve.

## 3. Objectifs
Relier file/path observations, process creation/modification chronology, hash/signature refs existantes, related signals et éventuellement reputation/context refs sourcés ; préserver sensitive path masking et provenance.

## 4. Non-objectifs
Aucune file acquisition/download, forensic image, content parsing, quarantine, malware verdict, Evidence creation, reputation provider selection ou response.

## 5. Propriétaire
Endpoint Agent possède le File Context local dérivé. Investigate conserve Artifact/Evidence/Finding et forensic acquisition ; Shared conserve generic linking/search.

## 6. Utilisateurs
SOC/Investigate Analyst ; Endpoint Operator ; Detection Engineer ; Security/Privacy Reviewer ; Auditor.

## 7. Conditions d’entrée
File/filesystem observations EPT-2 ; permissions sur paths/hashes ; source/freshness connues ; related process/signal refs optionnelles ; aucune acquisition supplémentaire.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| file/filesystem observations | CAP-EPT-018 | faits locaux | oui | source freshness | context unavailable |
| process associations | CAP-EPT-017/037 | relations | non | snapshot/observation time | relation unknown |
| detection refs/context | CAP-EPT-032..035 | détection locale | non | candidate freshness | context sans détection |
| reputation/context reference | source existante autorisée | enrichissement | non | source-owned | aucune conclusion inférée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| File Observation | Endpoint Agent | path/op/hash/time | read |
| Process Investigation Context | Endpoint Agent | process refs | read |
| Local Detection Candidate | Endpoint Agent | related refs | read |
| Artifact/Evidence | Investigate | destination refs seulement | no creation |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| File Investigation Context | dériver/rafraîchir | Endpoint Agent | données observées uniquement |
| File-Process Relation | dériver | Endpoint Agent | relation sourcée |
| File Chronology Projection | dériver | Endpoint Agent | timestamps/limitations conservés |

## 11. Fonctionnalités
Construire contexte path/file, associer process et operations, afficher hash/signature/reputation refs seulement si déjà disponibles, lier related local detections et exposer chronology/limitations sans acquisition.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect file context | Analyst | File Context | 0 | read | metadata/relations visibles | non |
| correlate existing file refs | deterministic service | context | 1 | observations existantes | liens sourcés | non |
| request bounded metadata refresh | Operator | existing metadata | 2 | pas d’acquisition/content read | freshness update | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| associer path/process | oui | oui | oui | suggestion possible | refs/time keys |
| ordonner operations | oui | oui | oui | non nécessaire | timestamps |
| résumer file context | oui | oui | oui | oui | structured summary |
| déclarer malware | non | non | non | interdit | analyst/Investigate qualification |

## 14. États fonctionnels
`available`, `partial`, `restricted-path`, `hash-unavailable`, `stale`, `unsupported`, `unavailable`, `unknown`.

## 15. États d’interface
Aucun Screen ID. Sensitive paths restent masked/minimized ; missing hash/reputation ne doit pas être présenté comme clean.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| File Investigation Context | concept Endpoint | CAP-EPT-042..045 | no content acquisition |
| process/file chronology | relation/context | Analyst | source/time refs explicites |
| missing-acquisition requirement | boundary diagnostic | futur EPT-4 | aucun Collection request créé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| file observation/detection | pivot | File Context | path/hash/process refs | origin retained |
| File Context | process/network/timeline pivot | CAP-EPT-037/039/042 | stable refs | no collection |
| file bytes/content required | expansion impossible localement | EPT-3 stop | missing-data reason | futur EPT-4 uniquement |

## 18. Dépendances
CAP-EPT-017/018/024/026/032..037 ; Endpoint Investigation README ; Investigate Artifact/Evidence ownership ; Shared Linking/Search ; `OPEN-008`.

## 19. Source de vérité
Endpoint est SOT des observations/context refs locales ; Investigate est SOT d’Artifact/Evidence et de toute qualification forensique. Une reputation source reste owner de ses propres facts.

## 20. Provenance et audit
Path refs/masking state, file operation ids/time, hashes/signatures source, process relations, detection refs, reputation source refs, snapshot/freshness, tenant et actor/service.

## 21. Permissions fonctionnelles
File context read, sensitive path read, hash/reputation ref read selon source, provenance, cross-tenant deny. Aucune permission de file content acquisition/quarantine.

## 22. Limites et erreurs
File context ≠ forensic acquisition ; hash match ≠ malware proof ; reputation hit ≠ compromise ; absent hash ≠ benign ; stale/deleted path peut limiter la résolution.

## 23. Métriques
Context completeness, process-link coverage, hash availability, restricted-path count, stale/unavailable reasons, pivot success.

## 24. Classification de livraison
`draft / defined / planned`; pas de downloader, collector, parser ou quarantine engine.

## 25. Critères d’acceptation
**Given** un hash existe dans une observation, **When** File Context est construit, **Then** le hash est référencé avec sa source sans verdict malware automatique.

**Given** le fichier n’est plus disponible localement, **When** le contexte est demandé, **Then** l’état est partial/unavailable et aucune acquisition implicite n’est lancée.

**Given** l’IA est indisponible, **When** l’analyste pivote depuis un signal vers le fichier, **Then** metadata/chronology/relations sont accessibles déterministement.

## 26. Questions ouvertes
`OPEN-008` reste ouverte pour support des sources. Toute acquisition future appartient à EPT-4/Investigate Collection selon ownership.

## 27. Consommateurs documentaires
EPT-3 timeline/pivots/summary ; Investigate Artifact/Evidence workflows ; Security ; Quality ; EPT-4 futur boundary.