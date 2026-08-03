# Validation Results

## Scope

Validation of the canonical CMDR documentation tree prepared for pull request #2 on `2026-08-03`. These controls validate repository structure and documentation contracts; they do not validate a software implementation.

## Results

| Control | Result | Evidence |
|---|---|---|
| Required top-level domains | PASS | 22/22 present; no extra active top-level domain |
| Root specification files | PASS | `README.md`, `INDEX.md`, `STATUS.md`, `CHANGELOG.md` |
| Expected path manifest | PASS | 781 expected file paths generated |
| Non-empty files | PASS | 0 empty files |
| Local Markdown links | PASS | 620 local links resolved; 0 broken |
| Filename conventions | PASS | 0 forbidden or non-conforming active filenames |
| Legacy condensed architecture | PASS | 0 active paths under the superseded domains |
| Screen front matter | PASS | 61 screen specifications checked |
| Screen 27-section contract | PASS | all screen specifications contain sections 1–27 |
| Required screen states | PASS | Loading, Empty, Partial, Error, Offline and Permission denied present |
| Screen identifiers | PASS | 61 unique canonical identifiers; 0 duplicates |
| Referenced permissions | PASS | 186 permission identifiers registered |
| Canonical object chain | PASS | 12/12 object files present in required order |
| Object ownership | PASS | mandatory owners verified |
| Named single sources of truth | PASS | 10/10 required canonical paths present |
| Required palettes | PASS | all exact CMDR and Command values present |
| Templates | PASS | 14/14 Markdown templates present |
| Reporting ownership | PASS | Reporting Engine is a Shared Capability; no competing `report` object |
| Root repository README | PENDING REMOTE | Must remain exactly `# cmdr` after GitHub commit |
| Branch comparison and PR state | PENDING REMOTE | Verify after committing; PR #2 must remain draft and unmerged |

## Interpretation

The in-memory canonical manifest passes structural and link validation. Remote GitHub validation is mandatory after the tree is committed. Product topics whose exact approved values were not supplied remain marked as draft questions in their owning documents rather than being invented or placed in a detached gap list.
