---
id: phase-4b4a-cloud-analysis-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-07
source-of-truth: quality-report
open_decisions: [OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Phase 4B.4A — Cloud Analysis Source Audit

## Purpose
This audit records what was actually consulted during the corrective closure. It distinguishes full reads from references, identified paths and absent sources. It does not infer that an unlisted file was read.

## Fully read — governance and source material
- `00-governance/source-material/cmdr-master-product-brief.md`
- `00-governance/source-material/README.md`
- `00-governance/source-material/product-boundaries.md`
- `00-governance/source-material/canonical-object-and-ownership-decisions.md`
- `00-governance/source-material/ai-and-automation-constraints.md`
- `00-governance/source-material/native-capability-strategy.md`
- `00-governance/source-material/explicit-non-goals.md`
- `00-governance/source-material/unresolved-decisions.md` — canonical and temporary variants compared.
- `00-governance/source-material/requirements-traceability-matrix.md`
- `00-governance/source-material/qualitative-baseline.md`
- `00-governance/source-of-truth-policy.md`
- `00-governance/documentation-rules.md`
- `00-governance/ownership-register.md`
- `00-governance/dependency-register.md` — canonical and temporary variants compared.
- `00-governance/registers/capability-register.md` — canonical and temporary variants compared.
- `00-governance/registers/object-register.md`
- `00-governance/registers/screen-register.md`
- `01-product-vision/product-boundaries.md`
- `templates/capability-specification-template.md`

## Fully read — Investigate and adjacent capability families
- `07-investigate/product-definition.md`
- `07-investigate/README.md`
- `07-investigate/modules/signals-and-hunt/README.md`
- `07-investigate/modules/cases-and-evidence/README.md`
- `07-investigate/modules/collection-and-live-response/README.md`
- `07-investigate/modules/collection-and-live-response/platform-settings-boundaries.md`
- `07-investigate/modules/analysis-workbench/README.md`
- `07-investigate/modules/analysis-workbench/platform-settings-boundaries.md`
- `07-investigate/modules/analysis-workbench/tool-and-studio-boundaries.md`
- `07-investigate/modules/analysis-workbench/memory-forensics/README.md`
- `07-investigate/modules/analysis-workbench/disk-and-filesystem-forensics/README.md`
- `07-investigate/modules/analysis-workbench/network-forensics/README.md`
- `07-investigate/modules/dynamic-sandbox/README.md`
- `07-investigate/modules/analysis-workbench/reverse-engineering-and-debugger/README.md`
- `07-investigate/modules/detection-engineering/README.md`
- `07-investigate/modules/detection-engineering/scope.md`
- `07-investigate/modules/threat-intelligence/README.md`
- `07-investigate/action-classification.md` — canonical and temporary variants compared.
- `07-investigate/automation-and-ai-model.md` — canonical and temporary variants compared.
- `07-investigate/capability-map.md` — canonical and temporary variants compared.
- `07-investigate/cross-product-links.md` — canonical and temporary variants compared.
- `07-investigate/object-consumption-map.md` — canonical and temporary variants compared.
- `07-investigate/screen-capability-map.md` — canonical and temporary variants compared.

## Fully read — Platform Settings, Command, Studio, Govern and Shared
- `10-platform-settings/sources-and-parsers/README.md`
- `10-platform-settings/secrets-and-connections/README.md`
- `10-platform-settings/tenants-and-environments/README.md`
- `10-platform-settings/health/README.md`
- `10-platform-settings/users-and-roles/README.md`
- `06-command/product-definition.md`
- `09-cmdr-studio/product-definition.md`
- `08-govern/product-definition.md`
- `12-shared-capabilities/README.md`
- `12-shared-capabilities/graph-engine.md`
- `12-shared-capabilities/timeline-engine.md`
- `12-shared-capabilities/query-engine.md`
- `12-shared-capabilities/export-engine.md`
- `12-shared-capabilities/reporting-engine.md`

## Fully read — security, permissions and trust
- `14-security-permissions-and-trust/permission-model.md`
- `14-security-permissions-and-trust/secrets-and-key-management.md`
- `14-security-permissions-and-trust/tenant-isolation.md`
- `14-security-permissions-and-trust/evidence-trust.md`

## Fully read — canonical object sources
- `05-domain-model/objects/case.md`
- `05-domain-model/objects/artifact.md`
- `05-domain-model/objects/evidence.md`
- `05-domain-model/objects/finding.md`
- `05-domain-model/objects/hypothesis.md`
- `05-domain-model/objects/principal.md`
- `05-domain-model/objects/role.md`
- `05-domain-model/objects/secret-reference.md`
- `05-domain-model/objects/data-source.md`
- `05-domain-model/objects/integration.md`
- `05-domain-model/objects/environment.md`
- `05-domain-model/objects/tenant.md`
- `05-domain-model/objects/timeline-entry.md`

## Fully read — existing screen sources used for conceptual mapping
- `07-investigate/modules/event-search/screens/event-search.md`
- `07-investigate/modules/case-workspace/screens/case-workspace.md`
- `07-investigate/modules/evidence/screens/evidence-board.md`
- `07-investigate/modules/hypotheses-and-findings/screens/hypotheses-and-findings.md`
- `10-platform-settings/sources-and-parsers/screens/sources-and-parsers.md`
- `10-platform-settings/secrets-and-connections/screens/secrets-and-connections.md`
- `10-platform-settings/tenants-and-environments/screens/tenants-and-environments.md`

## Fully read — quality, roadmap and publication evidence
- `STATUS.md`
- `CHANGELOG.md`
- `INDEX.md`
- `16-quality-and-validation/README.md`
- `16-quality-and-validation/architecture-manifest.md`
- `16-quality-and-validation/expected-path-manifest.md` — inspected as a historical minimum-path baseline, not treated as an exhaustive later-phase inventory.
- `16-quality-and-validation/validation-status.md`
- `18-roadmap-and-releases/README.md`
- `18-roadmap-and-releases/phase-4b4a-cloud-analysis.md`
- PR #2 metadata and state.
- temporary closure PR metadata, changed-file list and per-file patches.
- four published Cloud commit metadata and ancestry.

## Referenced and structurally validated
- `07-investigate/modules/cloud-analysis/README.md`
- the Cloud module’s supporting documents under `07-investigate/modules/cloud-analysis/`.
- `CAP-INV-601..618` capability files.
- the Cloud capability-register shard and conformance report.

These files are the published functional artifacts of the prior run. The corrective closure preserves them. Their existence, continuous IDs, front matter, 27-section contract and six mandatory-table contract are checked in the conformance workflow; they are not represented here as newly authored or re-scoped.

## Identified but not used as normative evidence in this correction
- additional product screens not required for the conceptual Cloud mapping;
- detailed implementation contracts under `17-implementation-contracts/`;
- provider-specific or vendor-specific Cloud documentation, because no provider is selected;
- archived sources under `99-archive/`, which are non-normative.

## Confirmed absent before Phase 4B.4A
- a pre-existing canonical Investigate Cloud Analysis module;
- a pre-existing Cloud subphase identifier conflicting with `4B.4A`;
- a prior `CAP-INV-6xx` family;
- a canonical Cloud provider schema, connector contract or query language;
- a canonical Cloud object family with final schemas and state machines;
- a dedicated Cloud Analysis Screen ID;
- a Mobile Forensics capability family or implementation.

## Findings
1. No source conflict requires modifying `CAP-INV-601..618`.
2. Platform Settings remains owner of providers, connectors, credentials, configured scopes, schemas, health, retention and administration.
3. Investigate remains owner only of provider-neutral Cloud analytical Sessions, observations, candidates, Hypotheses, correlations and packages.
4. Command, Govern, Studio and Shared ownership remains unchanged.
5. The temporary closure branch is append-oriented and preserves earlier Dependency, Requirement and baseline evidence after the corrective additions in this lot.
6. The conformance verdict must remain pending until remote publication checks complete.
