---
id: cmdr-typography-study
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-007
  - REQ-BRAND-008
  - REQ-PROD-051
---
# Étude typographique CMDR

## Statut

Cette étude définit les besoins et les candidats. Elle ne ferme pas `OPEN-004` et ne constitue pas une décision de famille finale.

## Besoins

- lecture dense prolongée ;
- hiérarchie éditoriale forte ;
- chiffres tabulaires ;
- distinction des caractères techniques ;
- support Latin étendu et besoins multilingues à confirmer ;
- rendu cohérent Windows, macOS et Linux ;
- variable fonts ou chargement maîtrisé ;
- performances web et application ;
- familles sans et monospace compatibles ;
- licence permettant intégration produit, sous réserve de revue juridique.

## Candidats sans

| Direction | Usage proposé | Avantages | Risques / tests |
|---|---|---|---|
| Inter + Inter Tight | Inter pour UI/body, Inter Tight pour display et titres | forte lisibilité écran, x-height élevée, chiffres tabulaires, large palette de poids | apparence SaaS générique si mal dirigée ; vérifier package Inter Tight, italique, hinting et petits corps |
| IBM Plex Sans | famille principale, éventuellement Condensed pour métadonnées | personnalité éditoriale et technique, nombreuses variantes de scripts, cohérence avec Plex Mono | association corporate forte, métriques plus larges, tests de densité nécessaires |
| Inter seul | rôles UI, body et titres par poids/espacement | simplicité de chargement et cohérence | distinction de marque moins forte ; dépend davantage du rythme et du wordmark |

Les dépôts officiels d'Inter et d'IBM Plex indiquent une licence SIL Open Font License 1.1. Une vérification juridique et de version reste obligatoire avant intégration.

## Candidats monospace

| Famille | Atouts | Risques |
|---|---|---|
| IBM Plex Mono | relation éditoriale avec Plex, ton calme | peut sembler trop lié à la stack Plex si Inter est retenu |
| JetBrains Mono | excellente distinction code et développement, OFL 1.1 selon dépôt officiel | personnalité développeur marquée ; ligatures à contrôler |
| Source Code Pro ou équivalent | lecture technique classique | disponibilité, licence et cohérence à vérifier avant shortlist |

Aucun fichier de police ne doit être ajouté dans cette phase.

## Rôles typographiques proposés

| Rôle | Intention |
|---|---|
| Display | communication de marque rare, grande taille, rythme éditorial |
| Page title | nom précis de workspace ou objet |
| Section title | rupture de section et navigation de lecture |
| Body | explication et contenu continu |
| UI label | action, contrôle et navigation |
| Metadata | owner, date, source, scope et fraîcheur |
| Table | densité, chiffres tabulaires et scan vertical |
| Badge | libellé court, jamais information complète |
| Code | code source et snippets |
| Query | requêtes, expressions et filtres techniques |
| Identifier | hashes, IDs, chemins, adresses et corrélations |

## Principes de hiérarchie

- pas plus de niveaux que le contenu n'en exige ;
- poids avant changement de famille ;
- taille et rythme avant couleur ;
- capitales réservées aux labels très courts ;
- métadonnées lisibles, non miniaturisées ;
- monospace réservé au contenu qui bénéficie d'une largeur fixe.

## Fallbacks proposés

Sans : `Inter`, `"IBM Plex Sans"`, `"Segoe UI"`, `"Helvetica Neue"`, Arial, sans-serif.  
Mono : `"IBM Plex Mono"`, `"JetBrains Mono"`, `ui-monospace`, SFMono-Regular, Menlo, Consolas, `"Liberation Mono"`, monospace.

L'ordre final dépendra de la famille approuvée et de la stratégie de distribution.

## Matrice de validation

Avant fermeture d'OPEN-004 :

- licence et reserved font names ;
- Latin étendu, français, anglais et langues cibles ;
- chiffres tabulaires, slashed zero et ponctuation ;
- corps denses et grands titres ;
- tables, code, hashes et chemins ;
- Windows ClearType, macOS et Linux ;
- écran standard et haute densité ;
- variable font, subset autorisé et poids réseau ;
- zoom, dyslexie, basse vision et longue session ;
- cohérence du wordmark et des identités produit.

## Critère d'acceptation

La famille finale n'est approuvée qu'après comparaison sur les mêmes écrans et contenus, avec mesures de lisibilité, performance et licence.
