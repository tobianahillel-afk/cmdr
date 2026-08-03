# Audit structurel de la fondation condensée

## Verdict initial

**FAIL structurel** au 3 août 2026 avant migration. La PR #2 possédait 73 fichiers cohérents entre eux, mais ne respectait pas le manifeste canonique : produits, marque, design system, Studio, Platform Settings, Endpoint Agent, parcours, sécurité, contrats, roadmap, templates et archives n’étaient pas représentés selon les chemins obligatoires.

## Tableau d’écarts et de migration

| Chemin réel initial | Chemin canonique attendu | Statut initial | Défaut | Correction appliquée | Contenu conservé / scindé | Priorité |
|---|---|---|---|---|---|---|
| `00-governance/*.md` | `00-governance/`, `adr/`, `registers/` | Incomplet / mal nommé | Noms hors convention, registres et ADR absents | Gouvernance recréée en kebab-case, quatre ADR et cinq registres ajoutés | Politique de source unique, cycle documentaire et décisions conservés puis normalisés | P0 |
| `01-product/` | `01-product-vision/` | Mal placé | Domaine non canonique et risques dispersés | Vision, portée, personas, modèle opératoire et risques migrés | Contenu produit conservé et scindé | P0 |
| absent | `02-brand/` | Absent | Cinq identités et palettes manquantes | Identités CMDR, Command, Investigate, Govern et Studio créées | Valeurs décidées intégrées ; valeurs non fournies signalées comme questions ouvertes | P0 |
| `03-experience/` | `03-design-system/` et `04-experience-architecture/` | Couvrant plusieurs sources | Design et architecture UX mélangés | Foundations, layouts, components, patterns et règles d’expérience séparés | Accessibilité, navigation et interactions migrées | P0 |
| `02-domain-model/` | `05-domain-model/objects/` | Incomplet / redondant | Tous les objets réunis dans un document | Un fichier par objet et chaîne canonique exacte | Définitions utiles scindées ; permissions déplacées vers le domaine sécurité | P0 |
| `10-command-center/` | `06-command/` | Mauvais produit / chemins non canoniques | Pages générales sans structure module/écran | Modules Command et écrans avec identifiants/front matter créés | Contenu opérationnel migré ; saved views consolidées dans leur source unique | P0 |
| `20-investigation-lab/` | `07-investigate/` | Mauvais produit / chemins non canoniques | Chaîne forensics utile mais non structurée | Modules Investigate et écrans canoniques créés | Triage, case, evidence, SIEM, static, sandbox, reverse, debugger, memory et disk conservés | P0 |
| `30-response-governance/` | `08-govern/` | Mauvais produit / chemins non canoniques | Gouvernance et exécution non séparées en modules canoniques | Govern recréé avec décisions, politiques, autorités, runs et rollback | Contenu utile migré | P0 |
| absent | `09-cmdr-studio/` | Absent | Skills, Automation Agents et workflows non spécifiés | Domaine complet Library, Builder, Control Room, Assurance et lifecycle ajouté | Nouvelles sources canoniques, sans duplication | P0 |
| éléments dans `40-platform/` | `10-platform-settings/` | Incomplet / mélangé | Administration confondue avec capacité partagée | Domaine Settings séparé : tenants, fleet, policies, sources, models, secrets, sandbox, retention, health, audit, preferences | Contenu plateforme migré et scindé | P0 |
| mention d’Agent seulement | `11-endpoint-agent/` | Absent fonctionnellement | Pas d’EDR natif complet | Huit domaines EDR ajoutés avec télémétrie, détection, investigation, collecte, réponse, confinement, résilience et sécurité | Mentions antérieures remplacées par liens | P0 |
| `40-platform/` | `12-shared-capabilities/` | Mal placé / trop large | Reporting et services partagés sans propriétaires précis | Capacités partagées individualisées ; reporting-engine devient source unique | Recherche, intégrations, audit, notifications et reporting migrés | P0 |
| scénarios dispersés | `13-user-journeys/` | Absent | Parcours et transitions non inventoriés | Parcours principaux et scénario ransomware ajoutés | Scénario existant conservé et relié | P1 |
| permissions dans `02-domain-model/` et pages | `14-security-permissions-and-trust/` | Source parallèle | Modèle de permission mal placé et partiellement répété | Modèle et catalogue centralisés ; pages ne gardent que des références | Sémantique existante consolidée | P0 |
| absent | `15-content-and-language/` | Absent | Aucun contrat éditorial ou terminologique | Terminologie, voix, erreurs, états et localisation ajoutés | Nouvelles sources canoniques | P1 |
| `50-quality/` et `OPEN_GAPS.md` | `16-quality-and-validation/` + propriétaires | Non canonique | Dossier et nom interdits ; lacunes détachées | Qualité migrée ; risques, dépendances, roadmap et questions répartis | Critères utiles conservés ; ancienne liste supprimée | P0 |
| absent | `17-implementation-contracts/` | Absent | Contrats API, événements, auth, audit et exécution manquants | Contrats d’implémentation créés par domaine | Nouvelles sources canoniques | P1 |
| absent | `18-roadmap-and-releases/` | Absent | Ordonnancement et release gates manquants | Roadmap, dépendances, versions et critères de sortie ajoutés | Questions d’ordonnancement migrées | P1 |
| absent | `templates/` | Absent | Aucun modèle normatif | Quatorze templates créés | Nouvelles sources canoniques | P1 |
| absent | `assets/` | Absent | Aucun inventaire d’actifs | Sous-domaines documentés avec README non vides | Références visuelles à rattacher ultérieurement | P2 |
| ancienne structure active | `99-archive/` | Redondante | Risque de deux architectures actives | Anciennes voies supprimées ; archive contient uniquement le registre de migration non normatif | Git conserve l’historique complet | P0 |

## Règle de clôture

Ce rapport documente l’écart historique ; il n’est pas une source de vérité produit. Les exigences actives se trouvent exclusivement dans les chemins canoniques enregistrés par les registres.
