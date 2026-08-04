---
id: investigate-action-classification
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-007
  - OPEN-013
---
# Action classification — Investigate through Phase 4B.2B.1

| Classe | Sens | Autorité |
|---:|---|---|
| 0 | observation/inspection/comparison | permission |
| 1 | bounded extraction/export/static processing | capability/policy |
| 2 | reversible session/relation mutation | OPEN-013 |
| 3 | real-target containment | preparation only |
| 4 | destructive real-target action | Govern only |

| Capability | Actions représentatives | Classes | Acteur | Cible | Risque | Govern | OPEN |
|---|---|---|---|---|---|---|---|
| CAP-INV-301 | Ouvrir l’intake; choisir une famille; créer/reprendre session | 0, 2 | Analyst | analytical objects | routage prématuré | policy/OPEN-013 | OPEN-005, OPEN-014 |
| CAP-INV-302 | Créer session; modifier objectif; suspendre/reprendre | 0, 2 | Analyst | Analysis Session | contexte erroné | policy/OPEN-013 | OPEN-013, OPEN-015 |
| CAP-INV-303 | Sélectionner Tool; modifier paramètres; lancer analyse | 0, 1, 2 | Analyst | Tool Call/context | mauvais Tool/version | policy/OPEN-013 | OPEN-005, OPEN-015 |
| CAP-INV-304 | Preview; copier valeur; annoter incohérence | 0, 1, 2 | Analyst | Artifact/metadata | contenu sensible | policy/OPEN-013 | OPEN-014 |
| CAP-INV-305 | Inspecter structure; filtrer; extraire élément | 0, 1, 2 | Analyst | Artifact/Derived Artifact | extraction non bornée | policy/OPEN-013 | OPEN-005 |
| CAP-INV-306 | Extraire chaînes; filtrer; exclure candidat | 0, 1, 2 | Analyst | extracted content | candidat pris pour IOC | policy/OPEN-013 | OPEN-013, OPEN-015 |
| CAP-INV-307 | Inspecter binaire; comparer; extraire ressource | 0, 1, 2 | Analyst | binary result | verdict automatique | policy/OPEN-013 | OPEN-005 |
| CAP-INV-308 | Inspecter document; décoder; extraire objet | 0, 1, 2 | Analyst | document/result | activation de contenu | policy/OPEN-013 | OPEN-005, OPEN-014 |
| CAP-INV-309 | Inspecter archive; définir limites; extraire | 0, 1, 2 | Analyst | archive/job | expansion excessive | policy/OPEN-013 | OPEN-005, OPEN-013, OPEN-014 |
| CAP-INV-310 | Sélectionner Artifacts; configurer; comparer | 0, 2 | Analyst | comparison | différence prise pour conclusion | policy/OPEN-013 | OPEN-013, OPEN-015 |
| CAP-INV-311 | Créer dérivé; annoter; supersede/withdraw | 0, 1, 2 | Analyst | Derived Artifact | source remplacée | policy/OPEN-013 | OPEN-013, OPEN-014, OPEN-015 |
| CAP-INV-312 | Ouvrir provenance; reproduire; comparer | 0, 1, 2 | Analyst | replay/assessment | faux résultat reproduit | policy/OPEN-013 | OPEN-005, OPEN-015 |
| CAP-INV-313 | Sélectionner résultats; préparer candidate/draft | 0, 2 | Analyst | Evidence/Finding drafts | qualification prématurée | policy/OPEN-013 | OPEN-013, OPEN-015 |

No class 3/4 execution occurs in Analysis Workbench; no Artifact is executed.
