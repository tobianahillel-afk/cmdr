# Ransomware Reference Scenario

## Purpose

Provide one end-to-end example used to validate cross-console completeness. It is illustrative, not a canonical data model.

## Scenario

Command Center prioritises `INC-2026-0712`, a critical incident affecting the Finance Platform. The incident timeline shows phishing, credential use, lateral movement, data access and encryption signals. The analyst opens `CASE-2026-0712`.

Investigation Lab registers email, process, registry, network, memory and encrypted-file evidence. Analysts correlate entities, test hypotheses and confirm a Finding: ransomware execution on `FIN-SRV-02` with likely lateral movement. The Finding includes confidence and evidence provenance.

A recommendation proposes isolating the host. Response & Governance creates request `R-0241`, previews business impact, checks policy `RG-17`, collects required approvals and records rationale. The approved action launches a run with rollback capability.

Command Center receives the approved decision and execution result. The incident becomes contained only when telemetry confirms isolation. Residual risk, affected services and next steps remain visible.

## Acceptance use

The scenario passes only when every object link, permission check, audit event, state transition and return path is represented without manual re-entry.
