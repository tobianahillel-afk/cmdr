# Git change evidence

`cmdr-dev git-changes` derives the change set used by the impact engine from exact Git commit identities.

## Safety properties

- Only full 40- or 64-hex commit IDs are accepted.
- The all-zero GitHub sentinel is rejected.
- Commit IDs are passed to `git` as process arguments, never interpolated into a shell command.
- Missing commits are fetched by exact validated ID from `origin`.
- The diff uses `--no-renames` intentionally so a rename is represented as the deletion of the old path plus the addition of the new path. Both identities therefore reach impact analysis.
- Deleted files remain in the change set because impact analysis is path-based and does not require file existence.
- Paths are repository-relative, normalized, deduplicated and sorted.
- Empty ranges fail rather than silently producing a zero-check plan.

During E3-IMPACT-002A the GitHub workflow computes a real change set and validation plan in **shadow mode**. The existing full kernel remains the authoritative gate until the selected-check executor is independently verified.
