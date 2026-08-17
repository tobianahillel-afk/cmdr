---
id: success-metrics
domain: 01-product-vision
status: draft
owner: Product Operations
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-010
---
# Familles de métriques de succès

Aucune cible numérique définitive n'est décidée en Phase 1.

| Métrique | Définition | Raison | Source future | Propriétaire | Risque d'interprétation | Statut |
|---|---|---|---|---|---|---|
| Signal → triage | durée entre signal disponible et première disposition | capacité opérationnelle | events futurs | Command | optimiser la vitesse au détriment de la qualité | proposed |
| Incident → owner | durée avant ownership explicite | coordination | Incident audit | Command | owner automatique non accepté | proposed |
| Incident → Case | temps jusqu'à ouverture d'une investigation lorsque nécessaire | handoff | object links | Command/Investigate | tous les incidents ne nécessitent pas de Case | proposed |
| Case → Finding | durée entre intake et Finding revu | efficacité investigation | Case timeline | Investigate | complexité des cas variable | proposed |
| Action Request → Decision | durée d'attente et motifs | gouvernance | Govern audit | Govern | réduire le délai ne doit pas supprimer les contrôles | proposed |
| Décisions avec Evidence suffisante | part des Decisions dont le package satisfait les règles | qualité décision | Decision package | Govern | définition de suffisance à décider | proposed |
| Actions vérifiées | part des Response Runs avec vérification explicite | boucle résultat | Run/Result | Govern | vérification peut être différée | proposed |
| Rollback réussi | part des rollbacks atteignant l'état vérifié | sécurité réponse | Run/rollback | Govern | ne s'applique pas à toutes actions | proposed |
| Perte de contexte | transitions nécessitant ressaisie ou reconstruction | qualité UX | navigation analytics | Product/UX | télémétrie doit respecter confidentialité | proposed |
| Duplication d'objets | objets ou définitions actives concurrentes | cohérence | register audit | Product Architecture | faux positifs sur projections légitimes | decided family |
| Handovers complets | handoffs avec owner, risque, prochaine action et contexte | continuité | handover objects | Command | checklist remplie sans qualité réelle | proposed |
| Usage essentiel sans IA | workflows critiques exécutables sans modèle | résilience et contrôle | availability tests | Product/QA | un workflow disponible mais inutilisable n'est pas un succès | decided family |
| Automations interruptibles | runs exposant pause/stop et trace selon contrat | accountable automation | Studio telemetry | Studio | certaines étapes atomiques ne sont pas stoppables | proposed |
| Faux positifs | dispositions et causes par détection | qualité détection | Detection/Alert | Detection Engineering | labels humains imparfaits | proposed |
| Couverture de détection | services/comportements couverts et qualité de données | readiness | coverage model | Detection Engineering | quantité de règles ≠ couverture | proposed |
| Satisfaction par rôle | capacité à accomplir objectifs et comprendre responsabilités | adoption | research future | Product | scores agrégés masquent les rôles | proposed |
| Adoption des workspaces | usage répété par objectif et rôle | utilité | product analytics | Product | usage forcé ou absence d'alternative | proposed |
| Temps d'investigation | durée et temps actif par catégorie de Case | efficience | Case analytics | Investigate | complexité et qualité doivent être contrôlées | proposed |
| Temps de reporting | durée sélection → review → delivery | customer value | Reporting Engine | Shared/Command | rapidité ne doit pas réduire la revue | proposed |

## Règles

- Chaque métrique doit être segmentable par rôle, tenant, environnement et type de travail lorsque la confidentialité le permet.
- Une métrique de vitesse est associée à une métrique de qualité ou sécurité.
- Les targets numériques exigent données baseline, owner et décision explicite.
- Les analytics produit n'enregistrent pas le contenu sensible des Evidence.

## Critère d'acceptation

**Given** une proposition de réduire le temps Action Request → Decision,  
**When** une cible est définie,  
**Then** elle est accompagnée de la qualité du package, du taux de refus/rework et des exigences d'autorité, afin de ne pas mesurer une gouvernance contournée.
