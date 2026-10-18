# E10-INV-400A — Search one authorized Case by actual source family

## Material usable integration

The already implemented operator-local workflow imports Wazuh 4.x and generic normalized JSONL events together into ONE Case owner, saves a scoped private snapshot, reopens an authorized Event Search Job and verifies offline inventory with separate Wazuh/generic counts. The search UI still treated both as one generic set, requiring manual reference inspection to distinguish an alert from a generic event.

This increment **reuses the same first-party `POST /demo-local-search` SearchJob and Case/Timeline/typed-source owner**. The analyst selects `Toutes les sources`, `Alertes Wazuh uniquement`, or `Événements JSONL génériques uniquement` in the existing analyst interface. The `source_family` body field is optional; omission preserves the legacy `all` behavior.

- The source family derives exclusively from canonical **typed-source owner metadata** `MemoryTimelineSource.Wazuh != nil` already attached to an authorized imported entry, never a client-supplied source type, file label or overlapping `operator-local-jsonl[-feed-N]` provenance.
- Import membership and Case/Timeline ownership are independently authorized exactly as before. Every in-window candidate is reopened using the existing authenticated `/demo-source-open` typed-source owner **before** evaluating whether it belongs to Wazuh or generic and before matching a literal, Wazuh numeric agent/rule pivot, or authorized source reference.
- The existing first-party Event Search Orchestrator receives the same explicit tenant, environment, source permission, time bounds, audit/provenance, and returns **only completed authorized refs**. Wazuh bursts derive solely from returned authorized Wazuh members, while generic events never masquerade as alerts.
- A source family can be `all`, `wazuh`, or `generic`; unrecognized values are rejected at the HTTP trust boundary. Wazuh-only `rule_id` and `agent_id` pivots cannot target generic events. The form resets stale result refs on family changes, and checks that the family has not changed before displaying asynchronous search responses.
- Retain POST JSON body ≤4KiB, no bearer or query literal in URL, no raw preview in results, 500 authorized candidates/200 hits bound and **all-or-nothing** failure on owner revocation, traversal anomaly or overflow. No new query DSL, database, permissions subsystem, untrusted source metadata or paid APIs.

### Integration acceptance

The end-to-end test creates a **real mixed private saved Case** with two Wazuh alerts and one generic event in the same UTC time window, computes the exact persisted SHA256, reopens the same scoped owner and executes real authenticated SearchJobs using `field=source_ref`:

- `source_family=all` gives 3 source-reauthorizable references; `wazuh` gives exactly 2; `generic` gives exactly 1. Numeric `rule_id=5710` on Wazuh finds the two Wazuh alerts. Every reference opens with the existing first-party source owner. Unauthorized bearer/foreign Case/unknown family/incompatible generic+rule pivot fail closed.
- Raw Wazuh descriptions, generic text and bearer are absent from SearchJob response.
- Analyst VM test verifies selector POST privacy and clears formerly authorized results when the family changes; an incompatible generic+rule pivot is rejected without a request.
- Existing Go/Node Linux, Go Windows/macOS Actions and unchanged mandatory `engineering-kernel` **90% security-critical changed coverage floor** must pass. The parent PR #79 is blocked with 88.00%/90% and earlier ancestors also remain draft. This lot DOES NOT waive the gate or imply merge-readiness.

### Scope

Local synthetic operator-approved offline Case only. Wazuh family classification is based on the imported typed source contract, **not authenticated Wazuh origin**, custody, production IAM, continuous SIEM collection or a confirmed incident. The single giant PR #36 engineering foundation and stacked draft PRs still require separate reconciliation.

**Base** immutable PR #79 SHA `753632ae5ce049e4b1a7ad8e8bfc60cfa67fc963`; no L2-reviewed parent SHA modified. One independent read-only risk-based L2 review on a verified and frozen material delivered diff, not per commit.
