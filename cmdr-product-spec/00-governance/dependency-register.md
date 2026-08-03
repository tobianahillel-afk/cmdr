---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Registre des dépendances

| Dépendance | Propriétaire | Consommateurs | Risque actuel | Décision attendue |
|---|---|---|---|---|
| Moteur de recherche | Platform Engineering | Command, Investigate | langage non choisi | architecture et dialecte |
| Moteur de politiques | Security Architecture | Govern, Platform Settings | technologie non choisie | règles et explication |
| Orchestrateur de workflows | CMDR Studio | Govern, Studio | langage non choisi | contrat de version |
| Moteurs forensics | Investigate | Investigation Lab | natif/intégré à décider | catalogue de lancement |
| Infrastructure sandbox | Investigate / Platform | Investigate, Studio | profils non figés | images et egress |
| PKI Endpoint Agent | Platform Settings | Endpoint Agent, Govern | exigences à préciser | enrollment et rotation |
| Rétention et legal hold | Compliance | tous | durées non définies | matrice par objet |
| Signature et horodatage | Security Architecture | Govern, Audit | standard non choisi | norme et conservation |

Les dates et séquences de résolution sont maintenues dans `../18-roadmap-and-releases/dependency-roadmap.md`.
