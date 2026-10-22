# Secret scanning

CMDR uses an in-repository deterministic secret scanner implemented by `cmdr-dev`. No source content is sent to a third-party service.

## Modes

- PR/push: scan only normalized changed-path evidence from `git-changes`.
- Full audit: `cmdr-dev secret-scan --root <repo> --full-scan --json` scans Git-tracked files.

Deleted paths are ignored safely. Symlinks are never followed. Binary assets, common generated/build directories, generated-source markers and known binary/media extensions are excluded. A regular text candidate larger than 16 MiB fails closed instead of being silently skipped.

## High-confidence rules

Built-in rules cover private-key PEM blocks and high-confidence provider formats for GitHub, AWS access-key IDs, Slack, Stripe live secret keys and Google API keys. A broad generic password-assignment rule is intentionally excluded because it would create low-confidence noise in a blocking PR gate.

## Redacted evidence

A finding contains only rule ID, repository-relative path, line number and a SHA-256 fingerprint over the rule ID plus matched bytes. Matched bytes are never returned in the report or error text.

## Allowlist

`secret-allowlist.json` contains no raw secret. Every exception binds an exact rule ID, canonical repository path, SHA-256 fingerprint and non-empty reason. Moving the file or changing the value invalidates the exception.

During the first implementation commit, `SEC-SECRETS-001` remains `specified`. It becomes `active` only after push and pull-request CI prove the executable scanner.
