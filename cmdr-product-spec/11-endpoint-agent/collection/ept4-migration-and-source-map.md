---
id: endpoint-ept4-migration-source-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-4 Migration and Source Map

Strategy is additive. The seven historical Collection files and nine Live Response files remain preserved as module/reference sources. `CAP-EPT-047..064` become the normative capability layer.

Historical source terms that may imply stronger semantics are bounded rather than deleted:
- `Artifact Collection` / historical Artifact outputs → neutral Endpoint `Collection Item/Collected Technical Output`; Investigate qualification remains under OPEN-014.
- `file-actions` delete/quarantine/restore → EPT-5 boundary; only bounded read/download/transfer and explicitly sourced temporary upload enter EPT-4.
- `network-actions` block rules → EPT-5 boundary.
- `process-actions` terminate/suspend/resume → EPT-5 boundary.
- `service-actions` start/stop/restart/disable → EPT-5 boundary.
- command/script/terminal sources → EPT-4 technical execution semantics only, no final command/runtime/protocol.

No source file is mass-deleted, no competing physical object schema is introduced and no historical owner is silently transferred.