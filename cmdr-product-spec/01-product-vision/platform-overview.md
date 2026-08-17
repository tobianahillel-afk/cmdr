---
id: platform-overview
domain: 01-product-vision
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-AI-002
---
# Vue d'ensemble de la plateforme

## Couche produits

Command, Investigate, Govern, CMDR Studio et Platform Settings fournissent les workspaces. Endpoint Agent est un composant produit distinct. Shared Capabilities fournit les services communs.

## Couche objets

Les produits coopèrent via les objets canoniques et leurs projections. La chaîne nominale relie Telemetry Event à Result, sans imposer un pipeline obligatoire.

## Couche de contexte

Tenant, environnement, objets actifs, filtres, sélection et retour sont propagés selon les règles d'Experience Architecture. Le Context Bar rend ce contexte visible.

## Couche capabilities partagées

Recherche, linking, timeline, graph, notification, Saved Views, Reporting, export et audit peuvent être partagés. Leur présence documentaire ne prouve pas leur implémentation.

## Endpoint Agent

L'agent fournit la cible des capacités locales. Settings administre la flotte ; Investigate utilise inspection et collecte ; Govern contrôle le risque ; Command consomme situation et résultats.

## Automatisation

Studio possède la conception, l'assurance et l'observabilité agentiques. Les produits opérationnels restent propriétaires de leur workflow et de leurs objets.

## Intégrations

Chaque capability reçoit une classification : native, integrated, temporary-integration, planned ou out-of-scope. L'expérience reste orientée activité ; le moteur et sa provenance restent visibles dans la trace.

## Audit et confiance

Toute conclusion ou action importante doit pouvoir être reliée à des données, acteurs, policies, outils, versions et résultats. Les détails du Permission Model et des contrats techniques sont ultérieurs.

## État documentaire et implémentation

| Élément | État Phase 1 |
|---|---|
| Vision et frontières | définies en Draft |
| UX détaillée | Phase 3/5 |
| Fonctionnalités modules | Phase 4 |
| Écrans pilotes | Phase 6 |
| Objets et permissions détaillés | Phase 7 |
| Architecture technique | Phase 8 |
| Implémentation | non évaluée |

## Non-décisions

Ce document ne choisit ni microservices, base de données, bus, protocole, PKI, framework, moteur de recherche ou moteur forensic.

## Critère d'acceptation

**Given** un lecteur technique,  
**When** il consulte cette vue,  
**Then** il comprend les responsabilités et flux sans en déduire une architecture de services ou un état d'implémentation non prouvé.
