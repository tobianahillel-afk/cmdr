# Canonical Domain Model

## Objective and ownership

This document is the canonical definition of CMDR business objects. Owner: Product Architecture with Security Architecture review.

## Core objects

- **Tenant** — isolated customer or organisational scope.
- **Principal** — human or service identity acting in CMDR.
- **Role / Policy / Authority** — access and decision constraints.
- **Event** — immutable normalised telemetry record.
- **Alert** — detection output requiring triage.
- **Incident** — operational container representing a security situation and business impact.
- **Case** — investigation workspace, optionally linked to one or more incidents.
- **Evidence** — preserved artefact or derived analytical result with provenance and integrity metadata.
- **Entity** — host, account, IP, domain, process, file, service, cloud resource or other identifiable subject.
- **Timeline entry** — time-qualified event, assertion or action linked to sources.
- **Hypothesis** — testable analytical explanation with confidence and evidence links.
- **Finding** — supported conclusion with confidence, scope and evidence package.
- **Recommendation** — proposed response or follow-up based on findings.
- **Response request** — governed request to perform an action.
- **Decision** — approval, conditional approval, request for information or refusal with rationale.
- **Playbook** — versioned orchestration definition.
- **Run / Step / Rollback** — execution instance and reversible action history.
- **Agent** — managed execution or collection component.
- **Audit event** — immutable record of security-relevant activity.
- **Report** — versioned presentation or export of selected facts and conclusions.
- **Integration** — configured external system connection and capability contract.

## Key relationships

- An Incident aggregates Alerts, affected Entities and operational Timeline entries.
- A Case may investigate multiple related Incidents and owns Evidence, Hypotheses and Findings.
- Evidence references source Events, files, captures, snapshots or tool outputs and can relate to Entities.
- A Finding must link to supporting Evidence and may refute or confirm Hypotheses.
- A Response request references at least one Incident, Case, Finding or emergency justification.
- A Decision governs one Response request and may impose Conditions.
- A Run executes a versioned Playbook or atomic action under a Decision.
- Audit events reference the acting Principal, tenant, object, action, timestamp and correlation identifier.

## Integrity rules

Identifiers are stable and non-semantic. Human-readable numbers may be displayed but never used as the sole relational key. Evidence provenance and hashes are immutable after verification; corrections create a new version or annotation. Cross-tenant relationships are prohibited unless an explicit MSSP aggregation object contains only authorised summaries.
