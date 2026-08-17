---
id: contribution-workflow
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-009
  - REQ-PROD-012
---
# Workflow de contribution

## Préparation d'un lot

1. identifier la phase et son périmètre autorisé ;
2. lire sources, fichiers propriétaires, consommateurs et baseline ;
3. relever Requirement IDs, contradictions, questions et dépendances ;
4. définir les métriques avant modification ;
5. travailler uniquement sur la branche du lot.

## Modification

- conserver le contenu utile ;
- supprimer les définitions concurrentes ;
- ne pas transformer une proposition en décision ;
- maintenir les questions ouvertes structurées ;
- écrire des critères spécifiques ;
- mettre à jour registres et traçabilité lorsque le lot les affecte.

## Vérification

- fichiers non vides ;
- liens locaux ;
- noms et front matter ;
- Requirement IDs existants et uniques ;
- statut justifié ;
- absence de placeholder générique ;
- contrôle de répétition ;
- README et index du domaine ;
- état de la PR et du README racine.

## Git

- aucune modification directe de `main` ;
- commits fonctionnels, pas un commit par fichier ;
- pas de force-push ou réécriture d'historique sans accord ;
- PR conservée en brouillon jusqu'à la fin des lots P0 ;
- aucune fusion automatique.

## Rapport de lot

Le rapport enregistre fichiers lus, ajoutés, modifiés, supprimés, exigences, décisions, questions, métriques avant/après, contrôles, commits, SHA et recommandation de phase suivante.

## Critère d'acceptation

**Given** un lot limité à une phase,  
**When** le commit est créé,  
**Then** aucune phase ultérieure n'est commencée, les dépendants et la matrice sont à jour, le README racine est inchangé et la PR reste ouverte, brouillon et non fusionnée.
