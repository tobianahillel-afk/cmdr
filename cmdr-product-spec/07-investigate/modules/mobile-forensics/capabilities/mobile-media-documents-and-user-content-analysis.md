---
id: CAP-INV-710
title: Mobile Media, Documents and User Content Analysis
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-011, OPEN-013, OPEN-014, OPEN-015]
source-of-truth: canonical
---
# CAP-INV-710 — Mobile Media, Documents and User Content Analysis

## 1. Définition
Analyser photos, videos, audio, documents, downloads, notes, archives, thumbnails/previews, metadata et relations de création/import/export candidates, applications associées, shared records, duplicates, deleted candidates, restricted content et sélection de Derived Artifacts.

## 2. Problème utilisateur
Un fichier média présent sur un device ou backup ne prouve pas qu’il a été créé par l’utilisateur. EXIF-like metadata, thumbnails, previews et timestamps peuvent être copiés, modifiés, synchronisés ou incomplets et ne constituent pas une vérité certaine.

## 3. Objectifs
- explorer metadata et contenu selon permissions distinctes;
- préserver source/package/app/path/version, duplication/sync/deleted states et restrictions;
- comparer metadata/content relations without certain attribution;
- préparer location/timeline/recovery/Derived Artifact and Evidence handoffs.

## 4. Non-objectifs
No facial/biometric identity inference, content upload to external service, hidden metadata fabrication, automatic authorship/location truth, unsupported media repair, decryption bypass, provider integration, codec/format implementation or product code.

## 5. Propriétaire
Investigate owns Media/Document Observations and derivative selection. Raw content remains source material; Shared owns Export/Reporting; Security controls private/sensitive content access; Evidence workflow qualifies Evidence.

## 6. Utilisateurs
Principal : Mobile Forensics Analyst. Secondaires : Sensitive Data Reviewer, Evidence Reviewer, DFIR Analyst and Investigation Lead.

## 7. Conditions d’entrée
Session, authorized media/document source records or files, integrity/completeness/accessibility limits, content/metadata permissions and explicit purpose.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Media/document files or records | package/app/filesystem | source material | oui | source version | no content claim |
| Metadata/thumbnails/previews | source | descriptive records | non | source version | unknown |
| App/path/source relations | CAP-INV-706/708 | context | non | Session version | provenance gap visible |
| Integrity/completeness limits | CAP-INV-705 | trust boundary | oui | reviewed version | partial/unverified |
| Privacy/classification | Security/source | access boundary | oui | access time | metadata-only/masked/denied |
| Sync/deleted indicators | source/CAP-INV-714/715 | state candidates | non | source version | unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Session / Filesystem / App Data observations | Investigate | source/path/app context | read |
| Media/document source | source | content/metadata | read according to permission |
| Artifact | Investigate | derived relation | read/link |
| Location Observation | CAP-INV-711 concept | metadata/location comparison | read/link candidate |
| Shared Export/Comparison | Shared | mechanisms | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Media Observation | create/review/dispute/supersede | Investigate concept | source/version required; file ≠ user-created |
| Document Observation | create/review/dispute/supersede | Investigate concept | metadata ≠ certain truth |
| Duplicate/import/export relation candidate | create/dispute | Investigate | source/confidence required |
| Derived Artifact proposal | prepare | CAP-INV-718 | content restrictions inherited |

## 11. Fonctionnalités
Gallery/list/table metadata views, search/filter by type/time/app/path, authorized preview/content, metadata comparison, duplicate relation candidates, source/app links, thumbnail/preview distinction, deleted candidate marking, restricted-content masking, location/timeline pivots and bounded Derived Artifact selection.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Browse metadata | analyste | observations | 0 | metadata read | scoped list | non |
| View restricted content | sensitive reviewer | content | 0/2 | explicit permission/purpose | audited view | step-up possible |
| Compare metadata/duplicates | analyste | files/records | 1 | source access | diff/candidate relation | non |
| Annotate/dispute relation | analyste | observation | 2 | source | versioned disposition | OPEN-013 |
| Prepare Derived Artifact/export | analyste | selected file/record | 1/2 | content/export permission | CAP-INV-718 input | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Parse metadata | oui | oui where supported | oui | non | deterministic parser |
| Detect exact/metadata duplicates | oui | rules/hashes if source supplies | oui | suggestion | compare table |
| Group related media | oui | rules | oui | suggestion | tags/filters |
| Summarize authorized document | oui | deterministic metadata | oui | possible with permission | viewer/manual summary |
| Infer creator/person/location truth | non | non | non | interdit | human sourced assessment |

## 14. États fonctionnels
Observations use `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn`; source content may be `available`, `restricted`, `inaccessible`, `deleted-candidate`, `partial`.

## 15. États d’interface
Loading preserves selection; Empty distinguishes no represented media from inaccessible source; Partial lists missing types/areas; Error preserves metadata; Offline read-only; Permission denied hides protected preview/content; Stale shows source/version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Media/Document Observations | concepts | CAP-INV-711/714..718 | source, metadata, access, limits |
| Duplicate/import/export candidates | relations | analyst/timeline | no authorship inference |
| Location metadata candidates | relation | CAP-INV-711 | metadata ≠ human presence |
| Derived Artifact proposal | package input | CAP-INV-718 | bounded source + privacy restrictions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-706/708 | media/document found | CAP-INV-710 | source/path/app/metadata/restrictions | source view |
| CAP-INV-710 | location metadata | CAP-INV-711 | metadata/source/confidence | media |
| CAP-INV-710 | sync/duplicate relation | CAP-INV-714 | source/version/relation candidate | media |
| Deleted candidate | recovery review | CAP-INV-715 | source region/record/quality | media |
| Timestamps/relations | correlate | CAP-INV-716/717 | observations + uncertainty | media |
| Selected file/record | derive | CAP-INV-718 | source, purpose, restrictions | media |

## 18. Dépendances
CAP-INV-705/706/708/711/714..719, Security Privacy, Shared Export/Comparison, Artifact/Evidence, OPEN-011/013/014/015.

## 19. Source de vérité
Raw content/metadata remains source material. Mobile observations and relation candidates are Investigate. Generic export/report mechanisms are Shared. Human creation/location/identity assertions require explicit analyst assessment beyond metadata.

## 20. Provenance et audit
Source package/path/app, file/record version, metadata parser/Tool/Run, access level, preview/content accesses, duplicate/comparison inputs, metadata/time fields, deleted/sync indicators, annotations, Derived Artifact selection, errors and human disposition.

## 21. Permissions fonctionnelles
Media metadata read, media content read, document metadata/content read, restricted content reveal, search/compare, location metadata read, Derived Artifact create/read/export, annotation/dispute, cross-tenant content relation and provenance export.

## 22. Limites et erreurs
Unsupported format, corrupted/partial file, misleading/copied metadata, missing original, thumbnail-only record, sync duplication, restricted content, permission denial or parser error remain explicit; no truth/authorship/location certainty is inferred.

## 23. Métriques
Media/documents observed, metadata-only vs content access, restricted/denied views, duplicate candidates, deleted candidates, derived selections, parser errors and downstream location/timeline/evidence handoffs.

## 24. Classification de livraison
`defined` / `planned`. No codec/parser library, biometric model, external content service, format list, API or product code delivered.

## 25. Critères d’acceptation
**Given** a photo with location metadata **When** reviewed **Then** metadata is shown with source/confidence and does not establish certain user presence.

**Given** a thumbnail without original **When** displayed **Then** it remains a thumbnail/preview record and no original file is invented.

**Given** no AI **When** media/documents are analyzed **Then** viewers, metadata tables, filters, comparison and human review provide the complete workflow.

## 26. Questions ouvertes
OPEN-011 future platform/media coverage; OPEN-013 sensitive access/annotations; OPEN-014 Artifact/material relations; OPEN-015 Tool/Run contracts.

## 27. Consommateurs documentaires
Location/Sensors, Backup/Sync, Deleted/Recovery, Timeline, Anomaly/Hypothesis, Artifact/Evidence handoff, Privacy, Objects, Permissions, Screens and Quality.
