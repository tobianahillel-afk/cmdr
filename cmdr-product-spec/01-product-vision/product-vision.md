---
id: product-vision
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-002
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-010
  - REQ-PROD-011
---
# Vision CMDR

## Destination

CMDR vise une plateforme opérationnelle de cybersécurité dans laquelle un signal peut devenir une situation compréhensible, une investigation peut produire une preuve traçable, une conclusion peut devenir une demande d'action explicite, une décision peut être prise sous autorité, et une action peut être observée, vérifiée puis réinjectée dans la situation opérationnelle.

La valeur de CMDR ne réside pas dans l'addition d'outils. Elle réside dans la continuité entre :

```text
détection → compréhension → investigation → preuve → recommandation
→ décision → exécution → résultat → amélioration
```

## État futur recherché

Dans l'état cible :

- humains, règles, moteurs déterministes, outils natifs, intégrations et Automation Agents travaillent sur les mêmes objets canoniques ;
- le contexte traverse Command, Investigate et Govern sans ressaisie ni perte de provenance ;
- chaque conclusion expose ses sources, sa confiance, ses contradictions et son auteur ou run ;
- chaque action importante expose autorité, périmètre, impact, rollback et résultat ;
- les capacités agentiques restent optionnelles, gouvernées et interrompables ;
- les workflows essentiels restent utilisables sans fournisseur de modèle ;
- les intégrations temporaires peuvent être remplacées progressivement par des capacités natives sans changer le modèle utilisateur.

## Ce que la vision n'implique pas

Cette vision ne signifie pas que :

- toutes les capabilities sont déjà natives ;
- l'Endpoint Agent complet est livré ;
- l'architecture technique est décidée ;
- l'IA est nécessaire à l'utilisation ;
- un produit unique remplace les responsabilités distinctes ;
- toute alerte crée obligatoirement chaque objet de la chaîne.

## Principes de réussite

CMDR atteint sa vision lorsque :

1. la situation opérationnelle et la preuve restent reliées ;
2. la responsabilité humaine est visible ;
3. les transitions sont explicites et attribuables ;
4. l'usage d'un moteur externe n'apparaît pas comme une rupture de produit ;
5. la vitesse n'est pas obtenue au prix de la gouvernance ou de l'intégrité ;
6. les résultats d'action alimentent détection, readiness et amélioration.

## Vision cible et état présent

| Dimension | Vision cible | État documentaire actuel |
|---|---|---|
| Continuité interproduit | contexte persistant de bout en bout | principes définis, UX détaillée en Phase 3/5 |
| Objets canoniques | sources uniques et transitions explicites | propriété définie, modèle détaillé en Phase 7 |
| Capabilities natives | moteur et expérience CMDR pour les capacités critiques | cible définie, classification et moteurs encore progressifs |
| IA | assistance optionnelle, traçable et gouvernée | principes définis, runtime ultérieur |
| Endpoint | EDR natif mature | cible `planned`, expérience en Phase 4, technique en Phase 8 |
| Implémentation | modules livrés avec preuves et tests | non évaluée par cette phase |

## Critères d'acceptation

**Given** une investigation sans modèle externe disponible,  
**When** l'analyste recherche, collecte, examine des Evidence et prépare un Finding,  
**Then** le workflow reste accessible par UI, règles, moteurs déterministes ou API appropriée ; les fonctions IA sont optionnelles et leur absence ne bloque pas l'activité essentielle.

**Given** un Result produit par Govern,  
**When** il revient dans Command,  
**Then** l'Incident, le Case, la Decision, le Response Run et les preuves restent reliés et attribuables.
