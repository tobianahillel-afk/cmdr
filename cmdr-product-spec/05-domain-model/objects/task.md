---
id: OBJ-TASK
type: object
domain: 05-domain-model
owner: Command
status: draft
updated: 2026-08-04
source-of-truth: canonical
---

# Task

## Objectif

Représenter un travail opérationnel assigné, suivi dans Command et lié à un objet source sans modifier automatiquement cet objet source.

## Propriétaire

Command pour la Task opérationnelle. Shared Capabilities peut fournir une Task Inbox ou une projection générique d’affectation, sans posséder le cycle de vie métier Command.

## Position dans la chaîne

- Amont : Incident, Result, exercice, plan ou autre objet source référencé ;
- Aval : résultat de coordination ou relation vers l’objet source.

## Champs canoniques existants

- `id immuable`
- `tenant-id`
- `created-at`
- `updated-at`
- `version`

La Phase 4A n’ajoute aucun schéma complet. Les besoins fonctionnels owner, équipe, contributeurs, watcher, échéance, priorité, état, dépendance, blocage et résultat sont transmis à la future phase Objets.

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.
- Une Task liée à un Case, une Decision ou un Result ne devient pas propriétaire de cet objet.

## États existants

- open
- in-progress
- blocked
- done
- cancelled

Cette liste reste Draft et ne vaut pas machine d’état finale.

## Invariants

- Le tenant est obligatoire.
- Les mutations sont auditées.
- Les références utilisent des identifiants stables.
- Une Task ne modifie jamais automatiquement l’objet source.
- L’ownership opérationnel d’une Task ne change pas l’ownership canonique des objets liés.

## Permissions existantes

- `perm.command.task.manage`
- famille de lecture Command à atomiser ultérieurement.

## Audit et provenance

Toute création, affectation, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le propriétaire correspond aux décisions canoniques de Command.
- Shared Task Inbox reste une capability consommatrice et non un owner concurrent.
- Aucun champ, API, protocole ou stockage définitif n’est ajouté en Phase 4A.
- Les références restent résolubles après versionnement.

## Questions transmises

La future phase Objets doit préciser les champs, cardinalités, spécialisations éventuelles et machine d’état sans créer une Task concurrente par produit.
