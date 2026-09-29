# Validation check catalog

The CMDR impact engine selects checks by **identifier**, never by arbitrary shell text.

`engineering/testing/check-catalog.json` is the strict registry of currently known validation checks. Each check declares:

- a stable `CHK-*` identity;
- a built-in `executor_key`;
- one or more closed risk domains;
- a cost tier: `fast`, `standard`, or `heavy`;
- whether it is mandatory when selected;
- whether it is always part of the PR safety floor;
- trigger path patterns;
- prerequisite check IDs;
- the expected evidence kind.

The catalog is metadata only. It cannot introduce executable commands. New executor keys require code changes and review in `cmdr-dev`.

## Safety model

- Unknown check IDs are denied.
- Unknown executor keys are denied.
- Unknown risk domains and cost tiers are denied.
- Check prerequisite cycles are denied.
- Mandatory checks require explicit risk domains and trigger paths.
- Prerequisites are closed over known check IDs.

E3-B will compute affected domains/paths. E3-C will turn those impacts into a minimum safe, risk-adaptive plan while preserving mandatory checks and prerequisite closure.
