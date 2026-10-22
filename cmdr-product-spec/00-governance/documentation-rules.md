---
id: documentation-rules
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-012
---
# Règles documentaires

## Objet

Ce document définit comment produire une spécification CMDR normative, traçable et non redondante. Il s'applique à tous les fichiers actifs sous `cmdr-product-spec/`.

## Classes de documents

| Classe | Rôle | Peut être normative ? |
|---|---|---:|
| Source matérielle | Conserver l'intention, la décision, la contrainte ou la question fournie | Oui sur l'intention, non sur l'application détaillée |
| ADR | Arbitrer une décision transversale avec options et conséquences | Oui après approbation |
| Document canonique | Définir l'application normative dans le domaine propriétaire | Oui selon son statut |
| Registre | Indexer propriétaire, identifiant, statut et source canonique | Oui comme index, pas comme définition détaillée |
| Spécification locale | Décrire l'usage d'une source canonique dans un produit, module ou écran | Oui sur le comportement local |
| Parcours | Relier étapes, objets et produits | Non pour redéfinir objets ou permissions |
| Template | Prescrire une structure minimale | Non sur le contenu métier |
| Proposition | Explorer une option non décidée | Non |
| Archive | Conserver une trace historique | Jamais |

## Requirement IDs

- Toute décision ou contrainte source possède un identifiant stable.
- Un document canonique liste les Requirement IDs qu'il applique.
- Une référence prouve la couverture uniquement si une section substantive applique l'exigence.
- Un Requirement ID ne doit pas être réutilisé pour une nouvelle intention.
- Une exigence remplacée conserve son identifiant et référence `supersedes` ou `superseded-by`.
- La matrice de traçabilité est mise à jour dans le même lot que le document canonique.

## Rôle des ADR

Une ADR est requise lorsque la décision :

- affecte plusieurs produits ou propriétaires ;
- crée ou modifie une frontière ;
- change une source de vérité ;
- impose une convention durable ;
- possède plusieurs options raisonnables et des conséquences difficiles à annuler.

Une ADR Draft ne remplace pas une décision source. Une ADR approuvée applique ou remplace explicitement une décision source et indique les Requirement IDs concernés.

## Rôle des registres

Les registres répondent à « où est la source ? » et « qui en est propriétaire ? ». Ils ne doivent pas recopier :

- le schéma complet d'un objet ;
- le comportement complet d'un composant ;
- les conditions détaillées d'une permission ;
- le contenu des écrans ;
- la logique d'un workflow.

## Décision source et spécification

La source matérielle conserve la formulation ou l'intention reçue. Le document propriétaire transforme cette intention en règle utilisable. Les consommateurs citent le document propriétaire et le Requirement ID ; ils ne recopient pas la décision source dans chaque fichier.

## Propriété et consommation

Le propriétaire :

- définit la sémantique, les états et les invariants ;
- approuve les changements ;
- maintient la source canonique ;
- coordonne les dépendants.

Un consommateur peut afficher, filtrer, lier, déclencher une transition autorisée ou créer une projection. Il ne peut pas créer une définition, machine d'état, permission ou source concurrente.

## Citation canonique

Une citation normative comprend :

1. le chemin relatif de la source canonique ;
2. le ou les Requirement IDs appliqués ;
3. la section ou le concept consommé ;
4. la description de l'usage local.

Un lien seul ne constitue pas une couverture fonctionnelle.

## Règles anti-duplication

- Une définition partagée vit dans un seul fichier propriétaire.
- Un fichier consommateur ne répète que les données nécessaires à son usage local.
- Les paragraphes standards doivent être placés dans une règle ou un composant canonique, puis référencés.
- Les critères d'acceptation sont spécifiques au résultat du document.
- Les états génériques d'interface sont configurés localement ; leur définition commune reste dans le Design System.
- Une copie temporaire doit indiquer son propriétaire, sa date d'expiration et sa migration.

## Statuts

Les valeurs autorisées sont `draft`, `in-review`, `validated`, `implemented` et `deprecated`. Les conditions sont définies dans [`status-lifecycle.md`](status-lifecycle.md). Les anciens termes `review` et `approved` ne doivent plus être introduits.

## Questions ouvertes

Une question ouverte utilise le schéma canonique :

```yaml
decision_id:
status: open
owner:
target_phase:
blocking:
depends_on:
question:
options:
required_evidence:
affected_files:
```

Un texte tel que « À compléter » sans identifiant, propriétaire et phase est interdit.

## Critères d'acceptation

Un critère doit identifier :

- le contexte initial ;
- l'acteur ou le rôle ;
- l'action ;
- le résultat observable ;
- l'objet ou le statut affecté ;
- l'erreur ou le refus attendu ;
- la trace requise ;
- une limite mesurable lorsqu'elle est décidée.

Les critères qui vérifient seulement l'existence du fichier, son propriétaire ou ses liens ne suffisent pas pour une spécification fonctionnelle.

## Mise à jour des dépendants

Avant modification, l'auteur recherche :

- références entrantes ;
- registres ;
- écrans consommateurs ;
- parcours ;
- contrats ;
- tests ;
- décisions ouvertes.

Après modification, chaque dépendant est mis à jour, déclaré non affecté avec justification ou enregistré dans le dependency register.

## Contrôles avant revue

- contenu substantiel et distinct du template ;
- Requirement IDs valides ;
- source canonique et propriétaire identifiés ;
- liens locaux résolus ;
- aucune contradiction P0 connue ;
- questions ouvertes structurées ;
- dépendants listés ;
- critères fonctionnels présents ;
- historique de décision conservé.

## Contrôles avant validation

- reviewers et approbateurs requis ont conclu ;
- demandes de modification résolues ;
- aucune décision bloquante ouverte ;
- matrice de traçabilité complète ;
- tests documentaires passés ;
- périmètre validé explicitement ;
- version et date enregistrées ;
- aucun document de priorité inférieure ne conserve une règle concurrente.

## Critère d'acceptation du document

**Given** une spécification candidate et ses sources,  
**When** elle entre en revue,  
**Then** sa classe, son statut, son propriétaire, ses Requirement IDs, ses dépendants, ses questions ouvertes et ses critères fonctionnels sont identifiables sans interprétation, et aucune définition partagée n'est recopiée comme une seconde source.
