---
id: endpoint-ept6-migration-and-source-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# EPT-6 Migration and Source Map

Strategy: **strictly additive**. Historical Update/Resilience/Security documents remain source/background material; `CAP-EPT-082..099` become the normative capability layer. No mass deletion or history rewrite.

| Historical source | Status/use | EPT-6 destination | Reconciliation |
|---|---|---|---|
| resilience/update-and-rollback.md | active source/skeleton | CAP-EPT-083..088 | signed-update/rollback language narrowed to functional requirements; no crypto/implementation claim |
| detection/detection-update.md | active specialized source | CAP-EPT-083/086 plus EPT-3 ownership | detection-content update remains distinct; generic Agent update semantics live in EPT-6 |
| resilience/offline-mode.md | active source | CAP-EPT-089/090 | deferred sync/local policy context retained |
| resilience/queueing-and-retry.md | active source | CAP-EPT-089/090 | physical/encryption wording treated as source intent, capability layer remains implementation-agnostic |
| resilience/crash-recovery.md | active source | CAP-EPT-091 | journal/idempotency intent retained without physical journal design |
| resilience/health-monitoring.md | active source | CAP-EPT-087/093 | reuses EPT-1 health; no duplicate health owner |
| resilience/resource-guardrails.md | active source | CAP-EPT-092 | budgets/kill-switch source intent bounded; no resource-control implementation |
| security/anti-tamper.md | active source | CAP-EPT-094 | candidate/protection semantics; no compromise or kernel mechanism claim |
| security/least-privilege.md | active source | CAP-EPT-095 | privilege context separated from permissions |
| security/secret-protection.md | active source | CAP-EPT-096 | raw secrets prohibited; Settings/Security ownership preserved |
| security/secure-storage.md | active source | CAP-EPT-096/089 boundary | storage requirements remain source intent; no storage engine/schema |
| security/local-audit.md | active source | CAP-EPT-097 | append-only treated as functional/history intent, not immutable/crypto proof |
| security/command-signing.md | active source | security constraint consumed across Endpoint | no signature/protocol format selected |
| security/mutual-authentication.md | active source | trust constraint consumed across Endpoint | no certificate/PKI implementation selected |
| security/privacy.md | active source | CAP-EPT-096..098 | global privacy ownership remains Security |
| security/self-protection-recovery.md | active source | CAP-EPT-094/093/097 | signed-repair/tamper-evidence intent retained without implementation |
| 17-implementation-contracts/endpoint-agent-update-contract.md | implementation-agnostic contract source | CAP-EPT-083..088 | package/wave/health/reversion boundaries aligned |
| 17-implementation-contracts/local-audit-sync-contract.md | implementation-agnostic contract source | CAP-EPT-090/097 | sequence/ack/dedup remain conceptual |

## Blocking gaps
No blocking capability-family gap was found after `CAP-EPT-082..099`. OPEN-008/013/015/017 remain explicitly unresolved but do not require fabricated implementation detail.

## Deprecated/duplicate semantics
No historical source is deleted. If old wording implies a physical solution or competing ownership, the capability contract narrows it to the canonical owner/boundary and the old file remains historical/source evidence rather than a second normative capability definition.