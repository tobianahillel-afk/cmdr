# E10-INV-120A — multi-Wazuh current-foundation integration handoff

## GitHub evidence (17 October 2026)

- Draft PR: https://github.com/tobianahillel-afk/cmdr/pull/49
- Branch: `feature/e10-inv-120a-multi-wazuh-current`; base `engineering/autonomous-development-foundation` at `7a4edd71365c16b41b839a4552a50eedda0cc8b9`.
- Functional snapshot under independent read-only L2: `0083121398c102e12fc3cf235ad836d8ae8f4957`; diff `7a4edd71365c16b41b839a4552a50eedda0cc8b9..0083121398c102e12fc3cf235ad836d8ae8f4957`.
- The L2 task is **scheduled**, not an approval; its report has not been accessed in the authoring conversation.
- On that exact snapshot: `investigation-timeline-runtime` PR GitHub Actions [39020728908](https://github.com/tobianahillel-afk/cmdr/actions/runs/39020728908) **SUCCESS**, including Linux, macOS and Windows jobs.
- Mandatory `engineering-kernel` PR GitHub Actions [39020728938](https://github.com/tobianahillel-afk/cmdr/actions/runs/39020728938) **FAILURE**: `CHK-SECURITY-TESTS` changed-security-critical coverage **88.36%**, below required **90.00%**. Earlier targeted iterations had 87.67% and 88.01%. This gate has NOT been waived and the PR is NOT ready for merge.
- Previous PR #48 multi-Wazuh work was based on an earlier foundation and is marked dirty/unmergeable; PR #49 ports the functionality without replacing existing current-base Wazuh rule/agent/window triage and HTTP analyst routes.
- No production IAM, persistent Case store, live Wazuh ingestion, signed forensic custody or causal correlation claims. `OPEN-013` remains unresolved. Main and PR #36 have not been merged or modified by this branch.

## Material user path implemented

2–4 operator-selected offline Wazuh 4.x `alerts.json` files → normalize Wazuh fields → namespace raw ID collisions by ordinal → import into one synthetic-owner-authorized Case → authenticated loopback Timeline and typed-source drill-down. The same owner chain supplies single-Case triage, independently authorized second-Case aggregate comparison and nearest certain-time cross-Case hints. Denied Cases, malformed late files and duplicate hardlinks must fail without partial receipt or bearer exposure. Existing one-file modes stay available.

## Next action after independent risk review

Triage the concrete L2 findings at immutable SHA and locate the **few materially untested trust-boundary paths** contributing to the 90% critical gate. Prefer a focused integration scenario or simpler production branch structure to another battery of near-identical tests. Re-run actual kernel and platform CI; do not soften `CHK-SECURITY-TESTS`, fabricate green results or open another L2 for this same functional increment. Document R3/R4 separately and defer. Keep PR draft and do not merge pending required reviews/CI.

The next significant product gap is binding the working local Case/Timeline/source journey to a trustworthy production authorization/persistence source under `OPEN-013`. That requires explicit ownership and a security design decision; it is not silently in scope for this synthetic offline lot.
