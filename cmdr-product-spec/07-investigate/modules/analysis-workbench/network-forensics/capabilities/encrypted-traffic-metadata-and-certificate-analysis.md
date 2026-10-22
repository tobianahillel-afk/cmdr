---
id: CAP-INV-389
title: Encrypted Traffic Metadata and Certificate Analysis
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-389 — Encrypted Traffic Metadata and Certificate Analysis

## 1. Définition
Examiner les conversations chiffrées candidates, métadonnées de négociation, versions et paramètres fonctionnels disponibles, identités déclarées, certificats, chaînes, validités, erreurs, incohérences et reprises candidates sans déchiffrement non autorisé.

## 2. Problème utilisateur
La validité technique d’un certificat ou la présence de chiffrement peut être confondue avec identité, confiance métier ou malveillance, tandis qu’un payload inaccessible peut être pris pour absent.

## 3. Objectifs
- identifier conversations chiffrées candidates et métadonnées disponibles
- voir versions, paramètres, noms ou identités déclarés
- voir certificats, chaînes, périodes de validité, émetteurs et sujets déclarés
- voir erreurs, incohérences, reprises ou échecs candidats
- comparer conversations, annoter, relier à Entity et préparer Hypothesis

## 4. Non-objectifs
Aucune acquisition active, administration de capteur/Fleet, packet generation/injection/crafting, scanning/interception active, replay, interaction cible, déchiffrement non autorisé, extraction ou usage de secret, exploit, évasion, règle Detection, objet Intelligence canonique, API, protocole interne, moteur, commande, code ou écran détaillé.

## 5. Propriétaire
Investigate possède Encrypted Traffic Observation and Certificate Observation, l’interprétation et les packages candidats. Command conserve Detection/Signal/Alert/Incident. Collection/Endpoint Agent produit la capture et ses limites. Settings administre capteurs/Fleet/Policies/storage/retention/health/timebase/secrets. Studio possède Tool/Tool Call/Workflow/Automation Run. Govern possède l’autorité réelle. Shared possède Entity/Graph/Timeline et les mécanismes génériques.

## 6. Utilisateurs
Principal : **Encrypted Traffic Analyst**. Secondaires : Certificate Reviewer, Network Forensics Analyst, Entity Analyst.

## 7. Conditions d’entrée
Case, source et session accessibles ; provenance, coverage, timebase, restrictions et permissions visibles ; scope borné ; aucune interaction active avec une cible.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Conversation chiffrée candidate et métadonnées/certificats disponibles | Conversation Candidate | source et contexte spécifiques | oui | versions de session | `partial` ou bloqué |
| Case, Hypothesis et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Tool/version et paramètres | Studio | traitement déterministe | oui pour classe 1 | réévalués au lancement | `tool-unavailable` |
| Permissions, payload policy et restrictions | Security/Settings | accès et minimisation | oui | décision courante | données masquées ou action refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, restrictions, return origin | consulter/lier |
| Conversation Candidate | owner source / Investigate | contenu autorisé, version, limitations | consulter uniquement |
| Tool / Tool Call / Automation Run | Studio | version, paramètres, statut, résultats | sélectionner/invoquer/lire |
| Entity / Graph / Timeline / Trace | Shared | relations, ordre, provenance | consommer sans redéfinir |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Encrypted Traffic Observation and Certificate Observation | créer, annoter, contester, superseder | Investigate | source, incertitude, auteur et version requis |
| Relation/sélection analytique | créer ou modifier réversiblement | Investigate | aucune fusion ou qualification silencieuse |
| Trace/Activity event | émettre | Shared | append-only et correlation ID |

## 11. Fonctionnalités
- identifier conversations chiffrées candidates et métadonnées disponibles
- voir versions, paramètres, noms ou identités déclarés
- voir certificats, chaînes, périodes de validité, émetteurs et sujets déclarés
- voir erreurs, incohérences, reprises ou échecs candidats
- comparer conversations, annoter, relier à Entity et préparer Hypothesis
- préserver source brute autorisée, partialité, restrictions, erreurs et return origin ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter, filtrer ou comparer | Encrypted Traffic Analyst | Encrypted Traffic Observation and Certificate Observation | 0 | lecture autorisée | vue sourcée | non |
| Comparer métadonnées et certificats | Encrypted Traffic Analyst | Tool Call / résultat | 1 | scope, Tool/version, permission | résultat attribué, éventuellement partial | selon policy |
| Annoter, contester ou relier | Encrypted Traffic Analyst | Encrypted Traffic Observation and Certificate Observation | 2 | permission réversible | nouvelle disposition versionnée | OPEN-013 |
| Préparer un handoff | Encrypted Traffic Analyst | package candidat | 2 | sources et limites présentes | package non qualifié | owner destination |

Classes 3/4 indisponibles ; toute action réelle est bloquée ou redirigée vers Collection/Live Response et Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Identifier/grouper Encrypted Traffic Observation and Certificate Observation | oui | parsers/règles/Tools | oui | suggestion attribuée | viewer, tables et filtres |
| Comparer/reconstruire | oui | oui si préconditions | oui | proposition incertaine | comparateur/reconstructeur |
| Expliquer erreur/contradiction | oui | catalogue/contrôles | oui | résumé sourcé | erreurs brutes/checklist |
| Confirmer IOC/Evidence/Finding/règle/Intelligence | owner humain | contrôles seulement | revue | jamais autonome | capabilities propriétaires |

Toute automatisation expose initiateur, moteur/agent et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude et acceptation/modification/rejet.

## 14. États fonctionnels
`candidate`, `metadata-available`, `certificate-available`, `partial`, `payload-inaccessible`, `invalid-or-expired-candidate`, `conflicting`, `disputed`, `superseded`. États fonctionnels, pas machine objet définitive.

## 15. États d’interface
Loading conserve contexte/sélection ; Empty distingue absence observée et donnée absente ; Partial expose gaps/troncatures ; Error conserve les résultats valides ; Offline bloque les traitements ; Permission denied masque brut/payload ; Stale expose source/timebase/Tool obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Encrypted Traffic Observation | Encrypted Traffic Observation | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Certificate Observation | Certificate Observation | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |
| Entity/Hypothesis candidate context | Entity/Hypothesis candidate context | Network Workbench, Case ou owner destination | source, partialité et provenance visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-386 | conversation chiffrée candidate | CAP-INV-389 | conversation, negotiation metadata, certificates et restrictions | Workbench avec return origin |
| CAP-INV-389 | relier identité candidate | CAP-INV-391 | names, certificates, addresses, sources et confiance | Workbench avec return origin |
| CAP-INV-389 | handoff analytique | CAP-INV-397 | selected metadata, contradictions, payload state et provenance | Workbench avec return origin |

Les transitions conservent tenant, Case, ownership, restrictions, erreurs, permissions et return origin.

## 18. Dépendances
Case/CAP-INV-105/107/108/109, CAP-INV-208/212/213/214, related CAP-INV-380..397, Static CAP-INV-301..313, Reverse CAP-INV-329..346, Memory CAP-INV-347..362, Disk CAP-INV-363..379, Studio Tool/Tool Call/Automation Run, Settings sensor/Fleet/Policy/health, Shared Entity/Graph/Timeline/Trace/Jobs/Export/Recovery, Govern and OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source de Encrypted Traffic Observation and Certificate Observation et de la disposition humaine. Les captures/résultats bruts restent chez leurs producteurs ; Studio, Settings, Shared, Command et Govern conservent leurs sources propriétaires.

## 20. Provenance et audit
Case, capture/version, Collection Request/Job, custody, session, capteur/interface, scope/coverage/timebase, Tool/version/Calls/Run, paramètres/filtres, entrées/sorties, erreurs, partialité, accès sensibles, annotations, décisions et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Encrypted Traffic Observation and Certificate Observation metadata/read | corrélation réseau | 0 | classification | possible | viewer/reviewer | Investigate | Permissions |
| Encrypted Traffic Observation and Certificate Observation processing/read | volume, brut ou payload | 1 | payload masqué par défaut | possible | initiateur/reviewer | Investigate/Studio | Permissions/Technique |
| Encrypted Traffic Observation and Certificate Observation annotate/link/handoff | mutation analytique | 2 | valeurs protégées | OPEN-013 | auteur/reviewer | Investigate | Permissions |

Matrice atomique, namespaces, RBAC/ABAC et step-up final reportés.

## 22. Limites et erreurs
- source partial/corrupted/truncated/restricted/stale, gaps, timebase incertaine, Tool indisponible, timeout, permission ou tenant mismatch ;
- Certificat ≠ identité certaine ; validité technique ≠ confiance métier.
- Trafic chiffré ≠ trafic malveillant.
- Payload inaccessible ≠ payload absent ; aucun contournement, extraction de clé ou attaque cryptographique.
- aucun résultat Tool, score ou sortie IA ne vaut conclusion.

## 23. Métriques
Usages, états partial/failed/disputed, gaps et limites héritées, provenance complète, accès sensibles refusés, dispositions humaines et reproductibilité ; aucune cible chiffrée définitive.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire seulement. Aucun moteur, produit, format, protocole final, API, commande, modèle ML, écran détaillé ou code.

## 25. Critères d’acceptation
### 1. Scénario
**Given** une conversation chiffrée sans contenu déchiffré autorisé
**When** l’analyste l’inspecte
**Then** seules les métadonnées autorisées sont visibles et aucun déchiffrement n’est lancé

### 2. Scénario
**Given** un certificat valide techniquement mais contradictoire avec le contexte
**When** il est analysé
**Then** la validité et les contradictions sont séparées, sans confiance métier automatique

### 3. Scénario
**Given** aucun modèle IA
**When** l’analyse est réalisée
**Then** parsers déterministes, viewers et revue humaine suffisent

## 26. Questions ouvertes
OPEN-005 moteurs ; OPEN-008 support capteurs/plateformes ; OPEN-013 classe 2 ; OPEN-014 Artifact/Attachment lorsque pertinent ; OPEN-015 bridge des Runs. OPEN-011/012 restent ouvertes et hors périmètre. Objets, permissions atomiques, écrans et contrats techniques sont futurs.

## 27. Consommateurs documentaires
Network Workbench, Case/Evidence/Finding, Event Search, Dynamic Sandbox, Memory/Disk, Static/Reverse, Entity Graph, future Detection Engineering/Intelligence et phases Objets/Permissions/Écrans/Journeys/Technique/Validation.
