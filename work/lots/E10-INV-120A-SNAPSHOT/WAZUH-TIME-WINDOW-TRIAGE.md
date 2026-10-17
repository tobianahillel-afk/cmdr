# E10-INV-120A — authorized Wazuh alert-volume time windows

## Usable investigative increment

Existing offline Wazuh 4.x `alerts.json` -> import into operator-authorized Case -> authenticated Case/Timeline/typed-source traversal -> rule-level triage now additionally returns descriptive five-minute UTC windows. Each bucket gives the **number of visible imported Wazuh alerts**, highest **known** source rule level (absent if none), and count with unknown severity. Return the ten most populated windows, sorting ties chronologically; `distinct_time_windows` discloses the total number of distinct windows and `unplaced_alerts` tracks time-uncertain imported alerts.

The source is the SAME raw Wazuh bytes used by the importer, mapped to the Case-scoped derived Timeline Entry identities. Windows are counted only inside the already-existing `collectDemoCaseTriageWithObserver` callback **after** independent Case, Timeline Entry and typed-source ACLs have passed on the local in-memory owner. No second parser, HTTP route, new database, search index, stream service, external API or paid dependency. The histogram and previous JSON fields remain compatible, so existing offline consumers need not change. Test inputs live on disk and exercise observed window counts, timestamp offset normalization, unknown severities, denied principal and absence of source descriptions/IDs in the aggregate.

## Risk and scope boundaries

This is **a descriptive analyst pivot**, not a security detection rule, inference of incident boundaries, p95, chain of custody, or causal correlation. Five-minute windows are calendar UTC buckets, not a sliding-window detector. All returned counts are for **imported, authorized Wazuh alerts only**, not all Case events. Partial or denied source results fail closed without serializing data; `OPEN-013` remains unresolved. Operator-provided local Case ACLs are synthetic, not deployed IAM. Persistent Case storage, Wazuh live API and auditable evidence custody remain future integration work.

Technical selection: reuse Go `time.Truncate` and `sort.Slice` + existing normalized-source-to-Case mapping. Avoid new event-sourcing, generalized alert scoring or vendor runtime until real production authority and retention requirements exist.
