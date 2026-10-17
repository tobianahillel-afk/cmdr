# Change impact analysis

`cmdr-dev impact` converts repository change evidence plus the selected strict-v2 work manifest into affected validation risk domains.

It does not select or execute checks; E3-C performs that step.

## Inputs

- a newline-delimited `--changes-file` containing repository-relative changed paths;
- a strict-v2 work unit (defaults to the active unit);
- the validated check catalog;
- the architecture boundary registry.

Deleted files are valid inputs because impact analysis operates on repository-relative path identity, not file existence.

## Conservative mapping

Each changed path contributes domains from matching check-catalog trigger patterns and from architecture ownership.

Additional escalation rules include:

- unowned repository path: `unknown + security`, high risk;
- path owned but matched by no check trigger: `unknown`, high risk;
- product-runtime boundary: architecture + runtime dependencies + security, high risk;
- CI control plane: work governance + repository health + security, high risk;
- `engineering/security/**`: security, high risk;
- `engineering/architecture/**`: architecture + security, high risk.

Product references in the active work manifest add product/coverage domains. Permission or OPEN-decision references also add security and high-risk status.

The result records sorted changed paths, domains, unknown paths, high-risk reasons and source evidence for every domain/reason pair.
