---
id: source-of-truth-policy
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-012
---
# Politique de source de vérité

## Finalité

La politique garantit qu'une décision explicite est appliquée par un seul document propriétaire, puis référencée par ses consommateurs sans créer de définition concurrente.

## Hiérarchie normative

Du niveau le plus prioritaire au moins prioritaire :

1. décisions explicites et exigences sources actives enregistrées dans `source-material/` ;
2. ADR approuvés ;
3. documents canoniques `validated` appartenant au domaine concerné ;
4. registres canoniques ;
5. documents `in-review` ;
6. documents `draft` ;
7. propositions et hypothèses ;
8. documents archivés ou `deprecated`, jamais normatifs.

À niveau égal, le document propriétaire du concept prévaut sur un consommateur. Si deux sources de même priorité revendiquent la propriété, la contradiction reste bloquante jusqu'à une décision explicite.

## Rôle de la source matérielle

- Elle conserve l'intention, la contrainte, la préférence ou la question fournie.
- Elle conserve les Requirement IDs et l'historique des remplacements.
- Elle ne devient pas une seconde spécification fonctionnelle.
- Elle n'est pas copiée intégralement dans les fichiers consommateurs.

## Rôle du document propriétaire

Le document propriétaire :

- définit l'application normative de l'exigence ;
- précise le périmètre et les exclusions ;
- établit les comportements et critères d'acceptation ;
- référence les Requirement IDs ;
- identifie les dépendants et questions ouvertes.

## Règles de consommation

Un fichier consommateur :

- référence la source canonique et les Requirement IDs applicables ;
- décrit uniquement l'usage, la projection ou l'interaction locale ;
- ne redéfinit ni objet, ni permission, ni composant, ni palette, ni capability partagée ;
- ne transforme pas une proposition ou une question ouverte en décision acquise.

## ADR et remplacement

Une ADR ne peut pas contredire une décision source active sans :

1. citer la décision et les Requirement IDs remplacés ;
2. présenter les options et la justification ;
3. recevoir l'approbation des décideurs requis ;
4. enregistrer `supersedes` ;
5. mettre à jour la matrice et tous les consommateurs.

Une ADR Draft documente une option ; elle ne remplace rien.

## Propositions et questions ouvertes

- Une proposition reste explicitement `proposed`.
- Une question ouverte conserve `status: open`.
- L'absence d'une réponse ne peut pas être masquée par un texte général.
- Une capacité `planned` ne peut pas être décrite comme livrée ou implémentée.

## Archives

- Aucun document sous `99-archive/` n'est une source normative.
- Une archive peut expliquer une migration, mais les consommateurs actifs ne la citent pas comme règle.
- Un document déprécié indique son remplaçant et la date de retrait.

## Processus de résolution des contradictions

1. **Identifier** les deux formulations, leurs chemins et propriétaires.
2. **Classer** chaque source dans la hiérarchie normative.
3. **Relever** les Requirement IDs et dépendants concernés.
4. **Appliquer la priorité** lorsqu'elle suffit ; sinon créer ou utiliser une question ouverte et demander une décision.
5. **Enregistrer le remplacement** dans l'ADR, le journal des décisions ou la source matérielle.
6. **Corriger les consommateurs** et déprécier ou supprimer la règle concurrente.
7. **Mettre à jour la matrice** avec l'état avant, après, la preuve et les dépendants.
8. **Conserver l'historique** sans laisser deux définitions normatives actives.

## Registre de résolution

Chaque résolution importante enregistre :

- identifiant de contradiction ;
- sources comparées ;
- niveau de priorité ;
- Requirement IDs ;
- décision appliquée ;
- document remplacé ;
- dépendants mis à jour ;
- date, auteur et approbateur ;
- éventuelle question ouverte restante.

## Critères d'acceptation

**Given** deux documents actifs portant des règles incompatibles,  
**When** la contradiction est examinée,  
**Then** leur priorité est démontrée, les Requirement IDs sont identifiés, une seule règle normative reste active, le remplacement est traçable, les consommateurs sont mis à jour et aucune archive n'est utilisée comme autorité.
