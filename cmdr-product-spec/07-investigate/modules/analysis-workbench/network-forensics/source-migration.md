---
id: investigate-network-forensics-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Source migration

## Audit result
Aucun module fonctionnel ou écran autonome Network Forensics n’existait avant cette phase. Aucun document n’est donc déprécié.

| Existing source | Useful needs consumed | Canonical relationship | Status |
|---|---|---|---|
| Collection `network-capture-request.md` | bounded acquisition, losses, scope, Artifact handoff | CAP-INV-380/382 consume result; acquisition remains Collection | active owner source |
| Dynamic Sandbox `network-behavior-analysis.md` | Run-specific simulated/blocked/allowed observations | CAP-INV-392 may correlate; no duplicate | active owner source |
| Memory `network-state-and-connection-artifact-reconstruction.md` | memory-resident candidates | CAP-INV-393 correlates; no PCAP analysis duplicated | active owner source |
| Disk system/network configuration observations | persistent configuration candidates | CAP-INV-393/397 correlate; no connection claim | active owner source |
| Event Search / Signal Triage | telemetry search and operational triage | CAP-INV-393 links; Network Forensics is not SIEM | active owner source |
| Entity Graph / Shared Entity | relation display and resolution | CAP-INV-391 consumes without local merge ownership | active owner source |
| Endpoint network connections/telemetry | live or endpoint-resident projections | correlation only; no acquisition administration | active owner source |

## Screen migration
No Network Forensics Screen ID existed. No screen is created or deprecated. Detailed composition is deferred to the Screens phase; Case Workspace, Event Search, Entity Graph and Technical Workbench provide temporary entry or pivot surfaces.

## Acceptance
No competing functional architecture, active duplicate, deprecated source or screen rewrite is introduced. Existing owners retain their sources and the new module records explicit handoffs.
