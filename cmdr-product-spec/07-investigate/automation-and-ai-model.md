---
id: investigate-automation-and-ai-model
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-006
  - REQ-AI-010
  - REQ-AI-011
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-015
---

# Automation and AI model — Investigate Phase 4B.1

## Autorité par producteur

| Producteur | Observe | Propose | Modifie Investigate | Confirme Finding | Crée Decision |
|---|---:|---:|---:|---:|---:|
| Humain autorisé | oui | oui | selon permission et classe | selon rôle et revue | non, Govern |
| Règle | oui | oui | seulement contrat classe 2 explicite | non | non |
| Moteur déterministe | oui | oui | validation, liens ou transformation tracée | non | non |
| Workflow | oui | oui | étapes déployées et permissionnées | non sans humain autorisé | non |
| Automation Agent | oui | oui | proposal par défaut | non | non |
| Govern | consomme | retourne statut | lifecycle Action Request | ne confirme pas Finding | oui selon autorité |

## Matrice fonctionnelle

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern |
|---|---:|---:|---:|---:|---:|---:|
| Écrire une Query | oui | templates | validation/autocomplete | préparation répétée | proposition lisible | non |
| Exécuter une Query | déclenche/arrête | contraintes de scope | parse, plan et run | répétition autorisée | jamais silencieusement si sensible | non normalement |
| Proposer un pivot | oui | mappings de champs | occurrence/relations | enchaînement autorisé | suggestion expliquée | non |
| Organiser des résultats | oui | groupements explicites | tri/comparaison | package de promotion | regroupement proposé | non |
| Créer une Hypothesis | oui | proposition de règle | associations sourcées | collecte d’éléments | `proposed` uniquement | non |
| Créer une Evidence candidate | oui | conditions explicites | provenance et versioning | prépare un candidat | suggestion uniquement | selon collecte sensible, pas qualification |
| Qualifier une Evidence | humain autorisé | contrôles de complétude | intégrité/relations | revue structurée | jamais automatique | non par défaut |
| Préparer un Finding | oui | vérification de références | dépendances/contradictions | brouillon et revue | `draft`/`proposed` | consomme plus tard, ne confirme pas |
| Confirmer ou contester Finding | reviewer humain | non | contrôles seulement | exige étape humaine | interdit | non ; Govern consomme |
| Préparer Action Request | oui | checklist/politique | complétude et refs | package et soumission autorisée | brouillon | reçoit et traite |
| Résumer Case/Timeline | oui | agrégation de faits | sélection versionnée | génération contrôlée | brouillon attribué | non |
| Préparer Report | oui | templates/redaction rules | citations et versions | revue et export request | brouillon | selon audience/action, pas auteur |

## Usages permis

- aide à l’écriture et explication de Query ;
- proposition de pivot ou regroupement ;
- Hypothesis proposée ;
- Evidence candidate ;
- Finding draft ;
- brouillon d’Action Request ;
- résumé de Case ou Timeline ;
- préparation de rapport et suggestions de redaction.

## Alternatives sans IA

Chaque capability essentielle conserve :

- saisie et revue manuelles ;
- autocomplete et validation déterministes ;
- builders, templates et snippets ;
- règles explicables ;
- recherches sauvegardées et Query Assets ;
- filtres, comparaisons et groupements déterministes ;
- workflows non agentiques ;
- checklists de qualification, Finding et Action Request.

L’absence de fournisseur de modèle ne change que la disponibilité des suggestions et résumés.

## Provenance obligatoire des sorties automatisées

| Champ | Exigence |
|---|---|
| Initiateur | utilisateur, règle, workflow ou événement autorisé |
| Producteur | type, ID et version du moteur ou agent |
| Run | Automation Run ou identifiant d’exécution lorsqu’il existe |
| Tool Calls | outils, paramètres pertinents, résultats et erreurs inspectables |
| Sources | objets, versions, Query, période et données utilisées |
| Temps | timestamp, timezone et période couverte |
| Statut | draft, proposal, accepted, modified, rejected ou failed |
| Incertitude | confiance ou limites lorsque la méthode le permet |
| Owner humain | personne ou rôle responsable de la disposition |
| Trace | correlation ID, audit event et lien vers le workspace source |

## Interdictions

- chatbot central ou obligatoire ;
- source ou résultat brut caché ;
- Query sensible ou coûteuse exécutée silencieusement ;
- Event, Artifact ou Attachment converti automatiquement en Evidence ;
- Evidence qualifiée ou modifiée sans action autorisée et nouvelle version ;
- Finding confirmé automatiquement ;
- Entity fusionnée silencieusement ;
- Action Request soumise sans contrat et permission ;
- Decision, Approval ou exécution classe 3/4 créée dans Investigate ;
- suppression, réécriture ou masquage de provenance ;
- auto-attribution de permission ;
- dépendance essentielle à un modèle ou Automation Agent.

## Govern et Human Gate

- Govern reçoit et possède le lifecycle de l’Action Request ;
- les classes 3 et 4 nécessitent Govern ;
- OPEN-007 reste ouverte pour la relation exacte Human Gate/Govern ;
- OPEN-013 reste ouverte pour les actions classe 2 ;
- OPEN-015 reste ouverte pour le bridge Automation Run/Response Run ;
- aucune de ces décisions n’est fermée par Phase 4B.1.
