---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---

# Action classification — Investigate Phase 4B.1

## Classes

| Classe | Sens fonctionnel | Autorité Investigate |
|---:|---|---|
| 0 | observation, recherche, inspection, filtre, pivot ou comparaison | exécutable selon permission |
| 1 | collecte ou export contrôlé | demande ou opération non destructive ; exécution de collecte détaillée reportée à 4B.2 |
| 2 | modification réversible d’un objet, lien ou état | selon permission ; gouvernance par défaut reste ouverte sous OPEN-013 |
| 3 | containment | Investigate prépare une Action Request ; Govern autorise ou décide |
| 4 | destructive ou irréversible | Investigate fournit Evidence et contexte ; Govern décide |

## Matrice par capability

| Capability | Action | Classe | Acteur | Cible | Risque | Govern | OPEN |
|---|---|---:|---|---|---|---|---|
| CAP-INV-001 | qualifier ou lier un Signal | 2 | SOC Analyst | qualification / Case relation | confusion avec priorité Command | selon politique | OPEN-013 |
| CAP-INV-002 | exécuter ou annuler une recherche | 0 | Analyst / Hunter | Query / Search Job | coût, scope ou données sensibles | non normalement | — |
| CAP-INV-003 | accepter une Query proposée | 2 | Analyst | Query draft | exécution sensible ou coûteuse | non ; confirmation requise | OPEN-013 |
| CAP-INV-004 | pivoter ou lier un Event au Case | 0/2 | Analyst | Event context / relation | Event présenté comme Evidence | non | OPEN-013 pour liaison |
| CAP-INV-005 | créer, suspendre ou clôturer un Hunt | 2 | Threat Hunter / Hunt Lead | Hunt workspace state | objet Hunt non finalisé | non normalement | OPEN-013 |
| CAP-INV-006 | publier, partager ou déprécier un Query Asset | 2 | Query Author | asset/version | usage avec sources incompatibles | non | OPEN-013 |
| CAP-INV-007 | grouper, annoter, exclure ou exporter des résultats | 1/2 | Analyst | workspace/result set | exclusion non justifiée, fuite export | selon sensibilité | OPEN-013 |
| CAP-INV-008 | enregistrer une provenance de run | 2 | système / Analyst | Trace relation | perte ou attribution erronée | non | OPEN-015 si automatisé |
| CAP-INV-101 | appliquer une Saved View de Cases | 0 | Case Analyst | Case Queue | confusion avec Work Queue Command | non | — |
| CAP-INV-102 | créer, affecter, suspendre, clôturer ou rouvrir un Case | 2 | Investigation Lead | Case | machine d’état et séparation des tâches | selon politique | OPEN-013 |
| CAP-INV-103 | créer, modifier ou supersede une Hypothesis | 2 | Analyst / Reviewer | Hypothesis | Hypothesis prise pour Finding | non | OPEN-013,015 |
| CAP-INV-104 | proposer une fusion Entity ou ajouter une relation | 2 | Analyst | Entity / relation | fusion silencieuse ou identité fausse | non ; revue requise | OPEN-013 |
| CAP-INV-105 | importer, versionner ou dériver un Artifact | 1/2 | Analyst | Artifact | provenance, données sensibles, dérivé non attribué | selon acquisition/sensibilité | OPEN-014 |
| CAP-INV-106 | joindre ou promouvoir une Attachment | 2 | Contributor | Attachment relation | fusion Artifact/Evidence implicite | non par défaut | OPEN-014 |
| CAP-INV-107 | qualifier une source comme Evidence | 2 | Analyst autorisé | Evidence | qualification automatique ou mutation silencieuse | non par défaut | OPEN-013 |
| CAP-INV-108 | confirmer, contester ou retirer de l’usage une Evidence | 2 | Reviewer | Evidence qualification | intégrité confondue avec pertinence | séparation des tâches possible | OPEN-013 |
| CAP-INV-109 | proposer, confirmer, contester ou supersede un Finding | 2 | Analyst / Reviewer | Finding | Finding assimilé à Decision/Result | Govern consomme seulement | OPEN-013,015 |
| CAP-INV-110 | ajouter ou corriger une entrée analytique | 2 | Analyst / Reviewer | Timeline Entry | inférence prise pour observation, historique effacé | non | OPEN-013 |
| CAP-INV-111 | publier une Note, joindre un fichier ou promouvoir une Task | 2 | Contributor | Note/Comment/Attachment/Task handoff | Task concurrente ou Attachment qualifiée implicitement | non | OPEN-013,014 |
| CAP-INV-112 | créer une review et préparer un handoff d’amélioration | 2 | Investigation Lead | Case review / handoff | conclusion non sourcée ou capability future simulée | non | OPEN-013,015 |
| CAP-INV-113 | soumettre une Action Request ou préparer containment | 2/3 | Investigation Lead | Action Request Govern | auto-approbation ou exécution locale | obligatoire pour classe 3/4 | OPEN-007,013,015 |
| CAP-INV-114 | rédiger, redacter ou demander l’export d’un Report | 1/2 | Author / Reviewer | Report draft / export request | divulgation ou citation stale | selon données/politique | OPEN-013,014 |

## Invariants

- la classe décrit l’effet réel, jamais le libellé du bouton ;
- une classe 3 ou 4 n’est jamais exécutée dans Investigate ;
- une automatisation ne s’accorde aucune permission et ne transforme pas une proposition en mutation effective ;
- la création ou qualification d’Evidence est réversible dans l’usage, mais la trace historique n’est pas supprimée ;
- les bulk links restent classe 2, affichent le périmètre et les échecs partiels ;
- OPEN-013 reste ouverte et aucune matrice atomique de permissions n’est finalisée.
