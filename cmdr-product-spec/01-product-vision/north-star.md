---
id: north-star
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-005
  - REQ-PROD-008
---
# North Star

## Résultat d'expérience

Un utilisateur autorisé peut passer d'un signal à un résultat vérifié sans ressaisie, sans perdre la provenance et sans changer mentalement de système à chaque produit.

## Conditions

- le contexte tenant, environnement et objets actifs est conservé ;
- chaque transition indique l'owner courant et la prochaine action ;
- chaque conclusion ouvre ses Evidence et sa trace ;
- chaque action élevée montre autorité, impact, scope et rollback ;
- chaque Result revient dans Command et reste relié à la Decision ;
- l'expérience essentielle fonctionne sans IA.

## Mesure conceptuelle

La North Star sera évaluée par la combinaison de transitions sans reconstruction de contexte, ownership explicite, décisions avec preuve suffisante, actions vérifiées, handovers complets et absence de duplication d'objet.

Aucune cible chiffrée n'est décidée.

## Critère d'acceptation

**Given** un Incident, deux Cases et une Action Request,  
**When** l'Incident Commander ouvre Govern puis revient,  
**Then** l'Incident, le Case sélectionné, l'Action Request, la vue et le point de retour sont restaurés.
