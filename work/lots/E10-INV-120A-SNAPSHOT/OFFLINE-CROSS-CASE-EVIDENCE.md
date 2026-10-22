# E10-INV-120A — offline two-Case comparison: verified execution handoff

## Published functional code and risk-based validation

- Functional code snapshot SHA: `b079d4b85d35e4dd77e90eedc934b93a9d0daf0a`
- Diff: `7fc09d39fd648e34201765b521c3366d9dd684bf..b079d4b85d35e4dd77e90eedc934b93a9d0daf0a`
- Commits: `a9804f15c80888720ff3314073e254006039bc9e` (first offline comparison) and `b079d4b85d35e4dd77e90eedc934b93a9d0daf0a` (material disjoint/unknown-time regression tests).
- Files: `product-runtime/investigation-timeline/cmd/case-timeline-demo/{main.go,compare.go,compare_test.go}` and `work/lots/E10-INV-120A-SNAPSHOT/OFFLINE-CROSS-CASE-REPORT.md`.

All four GitHub Actions runs on the **exact functional code SHA** were observed in the GitHub API as `completed/success`:

| Workflow/event | Run | Actual conclusion |
|---|---:|---|
| engineering-kernel push | [38994350032](https://github.com/tobianahillel-afk/cmdr/actions/runs/38994350032) | success |
| engineering-kernel pull_request | [38994353207](https://github.com/tobianahillel-afk/cmdr/actions/runs/38994353207) | success |
| investigation-timeline-runtime push | [38994350034](https://github.com/tobianahillel-afk/cmdr/actions/runs/38994350034) | success |
| investigation-timeline-runtime pull_request | [38994353217](https://github.com/tobianahillel-afk/cmdr/actions/runs/38994353217) | success |

The first code commit `a9804f15...` was **not** green: engineering-kernel push run `38994186831` reported changed security-critical coverage `84.25% < 90%`. The second code commit adds materially different forensic cases (disjoint Case time ranges and unknown source time) and restores all four checks. Do not erase the failing evidence.

## Functional acceptance

- Runs entirely in existing **synthetic local demo** via `-compare-case` and `-with-case` on an operator-selected snapshot; no additional HTTP route or third-party runtime.
- Each Case is separately authorized using the existing Case picker, paginated Timeline and typed-source owner checks. Test fixture has **207 permitted entries across two pages** on first Case and one on second; a linked denied source and another tenant are present but omitted.
- Comparison output contains only two authorized aggregate receipts, visible-entry and observed/inferred deltas, and optional intersection of **certain UTC bounds**, not individual evidence events. An unknown or disjoint chronology does not yield a fabricated overlap.
- Fail closed with no report bytes on unauthorized second Case, cancelled operation, invalid snapshot or failed output writer. No bearer, raw Timeline Entry identifier or source preview is exported.
- Time-overlap of known outer bounds does **not** prove the same event, shared real-world incident or causal link. No transactionally stable multi-Case snapshot.

## Independent review and next actions

A **single** read-only L2 task titled `CMDR — Revue L2 — E10-INV-120A-OFFLINE-COMPARISON` was successfully scheduled using the immutable code SHA `b079d4b85d35e4dd77e90eedc934b93a9d0daf0a`, same diff and explicit four run IDs. As of this handoff the actual L2 findings were **not accessible** here: `L2_EN_ATTENTE`. This scheduled ChatGPT review is not a native second-user GitHub approval.

The active graph unit remains `E10-INV-120A-SNAPSHOT` READY, not VERIFIED. Continue meaningful authorized Case-to-Timeline/source integration and resolve confirmed R0/R1/R2 reports before marking the unit complete. The existing open draft PR #36 is a large integration branch; feature PR #41 and #42 also existed as drafts targeting the foundation branch and require independent merge/conflict validation. No merge to main, production ACL/storage/authentication integration or OPEN-013 closure is claimed. Do not schedule duplicate L2 for this already-reviewed code SHA or cosmetic documentation-only updates.
