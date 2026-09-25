# Event Search frontend runtime benchmark — BENCH-0001

This is engineering decision evidence for `ENG-DEC-0001`. It is **not** a product benchmark or a production Search SLO.

## Question

Compare the first-party Native Web candidate with the first-party server-progressive alternative on a bounded Event Search interaction shape without introducing any third-party UI runtime.

## Workload

- 250 rendered synthetic forensic result rows.
- 30 initial-render samples per run.
- 160 state updates per run.
- Native Web mutates 12 existing rows plus result status per update.
- Server-progressive replaces a complete 250-row pre-rendered fragment per update.
- Seven repeated runs after warm-up.

This intentionally isolates frontend rendering/update mechanics. It excludes backend query latency, network transport, final virtualization policy and production browser-matrix claims.

## Environment

- Linux x86_64.
- Chrome/Chromium 144.0.7559.96 headless.
- V8 14.4.258.22.
- GPU disabled.
- No network used by the benchmark.
- No third-party JavaScript runtime.

The exact benchmark expression is `frontend-runtime-benchmark.js`; SHA-256 at measurement time:
`53cae19d268941a29420ad91791e524e8061ab52e92fcdafd9ae31305a95bea0`.

## Results

Across-run median p95:

- Native Web initial render: **12.5 ms**.
- Server-progressive initial fragment: **9.5 ms**.
- Native Web repeated update: **0.7 ms**.
- Server-progressive full-fragment update: **9.8 ms**.

The server-progressive candidate therefore retains a real initial-render advantage in this workload. The Native Web candidate has roughly 14× lower repeated-update p95, which is material for an interactive Event Search workspace.

## Reproduction

Run the JavaScript expression from `frontend-runtime-benchmark.js` in a clean Chromium page through DevTools Protocol or the DevTools console. Repeat seven times and aggregate the per-run p95 values by median. Do not load external scripts or network resources.

The committed JSON preserves all seven raw measurements and the environment.

## Limitations

- Synthetic result rows are not the final Event Search DOM.
- No final list virtualization strategy is selected here.
- Server/network round-trip cost is intentionally excluded; including it would change the server-progressive profile.
- Headless timing is not a substitute for the future supported-browser performance matrix.
- This evidence selects an architecture family only; it does not choose query syntax, storage, indexing, transport provider or OPEN-013 behavior.
