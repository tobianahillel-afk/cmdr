---
id: CAP-INV-365
title: Disk Image Integrity and Acquisition Context Review
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
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-SEC-001
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-365 — Disk Image Integrity and Acquisition Context Review

## 1. Définition
Revoir fonctionnellement source, autorisation, acquisition, erreurs, transformations, transferts, copies, vérifications et custody d’une Disk Image sans définir algorithme, format, stockage ou méthode d’acquisition.

## 2. Problème utilisateur
Une image techniquement lisible peut être incomplète, incohérente ou dépourvue de contexte. Sans revue, ses lacunes peuvent disparaître des analyses aval.

## 3. Objectifs
Afficher source/cible, initiateur/autorisation, méthode déclarée, timestamps/durée, Job status, erreurs/plages manquantes, transformations/transferts/copies, vérifications, custody et incohérences; permettre annotation, contestation, acceptation avec limites ou non-exploitabilité; préserver original.

## 4. Non-objectifs
Ne pas définir algorithme d’intégrité, format d’image, acquisition, stockage, protocole, moteur, commande, réparation ou modification de la source.

## 5. Propriétaire
Investigate possède la revue et sa disposition. Collection/Endpoint Agent restent sources du résultat déclaré; Settings des policies; Shared des traces; Evidence qualification reste CAP-INV-107/108.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : DFIR Analyst, Investigation Lead, Custody Reviewer et Audit Analyst.

## 7. Conditions d’entrée
Disk Image liée au Case; Collection Request/Job ou import source identifiable; custody et restrictions accessibles ou explicitement manquantes; permission de revue.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Disk Image/version | Investigate | source | oui | version sélectionnée | blocked |
| Request/Job/acquisition | CAP-INV-203/205/212 | contexte d’acquisition | oui si collectée | timestamps visibles | partial/disputed |
| Integrity/custody events | CAP-INV-213/214 | vérification/provenance | oui | ordre/version visibles | unverified |
| Transformations/transfers/copies | owners respectifs | lineage | selon parcours | sourcés | incomplete |
| Policies/permissions | Settings/Security | restrictions | oui | courant | restricted |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Disk Image/Artifact | Investigate | source/version/restrictions | lire |
| Collection Request/Job | Investigate | scope/status/errors | lire |
| Endpoint/Agent | Endpoint Agent | source/méthode déclarée | projection |
| Custody/Trace events | Investigate/Shared | chaîne et lacunes | lire |
| Decision/Approval | Govern | autorité éventuelle | lire |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Image Integrity Review | créer/update/dispute/supersede | Investigate concept | distincte de pertinence/Evidence |
| Verification status | enregistrer | Investigate | source et reviewer obligatoires |
| Limitation relation | créer/propager | Investigate | partialité héritée par outputs |
| Original image | aucune mutation | Investigate | immuable |

## 11. Fonctionnalités
Voir toutes les informations de collecte et custody, les plages manquantes, incohérences et copies; annoter/contester; accepter pour analyse, marquer partial ou non-exploitable; naviguer vers les sources; préserver original et propager les limitations.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter chaîne | Reviewer | Review | 0 | read | sources/lacunes visibles | non |
| Vérifier/contester | Reviewer | status | 2 | permission/justification | disposition versionnée | OPEN-013 |
| Accepter avec limites | Reviewer | image relation | 2 | limites explicites | analyses aval autorisées | OPEN-013 |
| Marquer non-exploitable | Reviewer | review | 2 | raison | block attribué | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assembler custody | oui | oui | oui | résumé | timeline/relations |
| signaler lacune | oui | règles | oui | explication | checklist |
| comparer copies | oui | metadata checks | oui | résumé | comparaison manuelle |
| accepter/contester | humain | non | workflow seulement | assistance | revue humaine |

## 14. États fonctionnels
`unverified`, `verifying`, `verified`, `partially-verified`, `inconsistent`, `incomplete`, `corrupted`, `disputed`, `restricted`.

## 15. États d’interface
Loading conserve image/source; Empty nomme les événements absents; Partial marque chaque maillon; Error garde les éléments valides; Offline lecture limitée; Permission denied masque données sensibles; Stale montre versions superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Integrity assessment | review concept | CAP-INV-363/364/378 | complétude/cohérence/limites distinctes |
| Acquisition context package | relation package | Case/Evidence Review | source, autorisation et custody navigables |
| Inherited limitation | business event | CAP-INV-366..379 | partialité jamais effacée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Disk Image | ouvrir revue | CAP-INV-365 | source, Job, acquisition, custody, restrictions | Intake |
| Review accepted/limited | créer session | CAP-INV-364 | image, status, missing ranges, owner | Review |
| Review disputed | signaler | Case/Trace | contradictions, sources, reason | Review |

## 18. Dépendances
CAP-INV-105/203/205/212/213/214/363/364/378, Shared Trace, Security, OPEN-008/013/014.

## 19. Source de vérité
Investigate possède review/status; Disk Image/Artifact et événements propriétaires restent sources. Une revue ne remplace pas l’original ni sa custody.

## 20. Provenance et audit
Case, source/Endpoint, Request/Job, image/version, initiator/authority, declared method, timestamps, transformations, transfers, copies, checks, reviewer, disputes, disposition et correlation IDs.

## 21. Permissions fonctionnelles
Disk Image read/raw/restricted, custody read/review/dispute, integrity accept/reject, provenance export. Masquage, step-up et SoD finaux sont reportés.

## 22. Limites et erreurs
Algorithme, format, stockage et acquisition hors phase. Source absente, timestamp incohérent, copie non liée, range missing, permission révoquée ou policy conflict donnent partial/disputed sans effacer le reste.

## 23. Métriques
Images verified/partial/disputed/non-exploitable, missing ranges, custody gaps, corrections et limitations propagées.

## 24. Classification de livraison
`defined` / `planned`; aucune implémentation ou méthode technique choisie.

## 25. Critères d’acceptation
**Given** une image partielle et des plages manquantes **When** la revue est ouverte **Then** limites, erreurs et provenance sont visibles et héritées par les résultats.

**Given** une copie sans transfert documenté **When** elle est revue **Then** la lacune produit `disputed` sans supprimer les autres événements.

**Given** aucun modèle IA **When** la revue est réalisée **Then** relations, timeline et checklist déterministes suffisent.

## 26. Questions ouvertes
OPEN-008/013/014 restent ouvertes; formats, checks et modèles finaux sont reportés.

## 27. Consommateurs documentaires
INV-DSK-001, CAP-INV-363/364/366..379, Artifact/Evidence/Custody/Replay et phases Objects/Permissions/Technique.
