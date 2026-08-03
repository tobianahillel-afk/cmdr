---
id: canonical-object-chain
domain: 05-domain-model
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Chaîne canonique des objets

```text
Telemetry Event
→ Detection
→ Signal
→ Alert
→ Incident
→ Case
→ Evidence
→ Finding
→ Action Request
→ Decision
→ Response Run
→ Result
```

## Règles

- Chaque étape possède son fichier dans `objects/`.
- Une transition produit un nouvel objet ou une relation; elle ne renomme pas l’objet amont.
- La chaîne représente une progression décisionnelle, pas une obligation de créer chaque objet pour chaque événement.
- Toute dérivation conserve la provenance et le tenant.
- Un Result ne modifie pas rétroactivement une Evidence ou une Decision; il alimente l’Incident et les métriques.
