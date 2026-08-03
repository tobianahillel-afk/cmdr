---
id: endpoint-agent-readme
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Endpoint Agent

## Objectif

Définir l’EDR natif CMDR complet: télémétrie, détection locale, investigation, collecte, live response, containment, résilience et sécurité.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Endpoint Agent Product Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Produit technique natif distinct de Platform Settings.
- La flotte administrative appartient à Platform Settings.
- Toute commande est signée, autorisée et auditée.
- Le mode hors ligne conserve sécurité et traçabilité.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- 00-governance/source-of-truth-policy.md

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.

## Exigences EDR natives obligatoires

L’Endpoint Agent est un EDR natif complet. Il couvre explicitement la télémétrie, la détection locale, l’inspection, la collecte, le terminal et l’exécution de commandes, la quarantaine et l’isolation, le rollback et sa vérification, le mode hors ligne, la protection anti-altération, la signature des commandes, l’audit local et le stockage sécurisé.
