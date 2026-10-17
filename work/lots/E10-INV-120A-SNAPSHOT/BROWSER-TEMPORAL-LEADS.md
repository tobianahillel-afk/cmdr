# E10-INV-120A-SNAPSHOT — browser temporal proximity triage

## Functional addition and reuse

The existing analyst browser comparison already loads **two independent, source/entry/Case-authorized paginated** Timeline snapshots (up to 5,000 entries each). The Go CLI `-compare-proximity` already calculates nearest known-time events from two Cases. This increment makes that capability usable in the local browser without a new endpoint, provider, schema, persisted state, privileged permission or external dependency.

After a successful two-Case comparison, the analyst may select an **integer 1–300-second** window and explicitly display up to 25 closest lead pairs. A sorted candidate array plus binary search matches the existing first-party Go algorithm's nearest-event selection and tie behavior, rather than introducing a database/index/scheduler. Entries with uncertain timestamps are excluded from proximity inference; their absence cannot be interpreted as absence of activity.

Each displayed pair provides a Case-specific **open entry** action. The browser never uses the pair to grant permissions: the existing `/case-source-navigation` and `/demo-source-open` requests reauthorize Case, Timeline Entry and typed source. Changing the token, either Case or filters clears the sensitive local result. Read-only and ephemeral loopback/demo restrictions apply, and all Case pages remain nontransactional.

## Acceptance and limitations

The existing Node browser harness now checks both authorized Case HTTP traversals, expected known-time gaps, max-window rejection, explicit user action, token invalidation, and fail-closed behavior when Case 2 is denied. This is **synthetic local triage**, not correlation/causality proof, production IAM or a custody-certified forensic report.

No claim that the Product Spec's broader CAP-INV-110 or OPEN-013 is complete; no new export backend, index, new privilege, ingestion or storage. Review CI for exact commit; do not reuse prior SHA test attestations. L2 should focus on truly material security, logic or browser integration defects instead of redundant micro-validation.
