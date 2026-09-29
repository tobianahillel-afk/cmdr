(() => {
  "use strict";
  document.body.innerHTML = "";
  const ROWS = 250;
  const UPDATES = 160;
  const INITIAL_ROUNDS = 30;
  const MUTATED = 12;
  const data = Array.from({ length: ROWS }, (_, i) => ({
    id: `evt-${String(i).padStart(5, "0")}`,
    source: i % 3 === 0 ? "endpoint" : i % 3 === 1 ? "network" : "identity",
    severity: ["info", "low", "medium", "high"][i % 4],
    summary: `Synthetic forensic event ${i}`
  }));

  const quantile = (values, p) => {
    const sorted = [...values].sort((a, b) => a - b);
    return sorted[Math.min(sorted.length - 1, Math.max(0, Math.ceil(p * sorted.length) - 1))];
  };
  const mean = values => values.reduce((a, b) => a + b, 0) / values.length;

  function buildRow(event) {
    const row = document.createElement("div");
    row.className = "row";
    row.dataset.id = event.id;
    for (const value of [event.id, event.source, event.severity, event.summary]) {
      const span = document.createElement("span");
      span.textContent = value;
      row.append(span);
    }
    return row;
  }

  function nativeInitial() {
    const root = document.createElement("section");
    const status = document.createElement("output");
    status.textContent = `${ROWS} results`;
    const list = document.createElement("div");
    const fragment = document.createDocumentFragment();
    for (const event of data) fragment.append(buildRow(event));
    list.append(fragment);
    root.append(status, list);
    return { root, status, rows: [...list.children] };
  }

  function nativeUpdate(context, n) {
    context.status.textContent = `${ROWS - (n % 7)} results · q=${n}`;
    const base = (n * 17) % ROWS;
    for (let j = 0; j < MUTATED; j++) {
      const row = context.rows[(base + j * 7) % ROWS];
      row.toggleAttribute("data-selected", ((n + j) & 1) === 0);
      row.children[2].textContent = ["info", "low", "medium", "high"][(n + j) % 4];
    }
  }

  function serverTemplate(n) {
    const root = document.createElement("section");
    const status = document.createElement("output");
    status.textContent = `${ROWS - (n % 7)} results · q=${n}`;
    const list = document.createElement("div");
    const fragment = document.createDocumentFragment();
    const base = (n * 17) % ROWS;
    const selected = new Set(Array.from({ length: MUTATED }, (_, j) => (base + j * 7) % ROWS));
    for (let i = 0; i < ROWS; i++) {
      const row = buildRow(data[i]);
      if (selected.has(i)) row.setAttribute("data-selected", "");
      row.children[2].textContent = ["info", "low", "medium", "high"][(n + i) % 4];
      fragment.append(row);
    }
    list.append(fragment);
    root.append(status, list);
    return root;
  }

  const serverTemplates = Array.from({ length: 16 }, (_, i) => serverTemplate(i));
  const serverInitial = () => serverTemplates[0].cloneNode(true);
  function serverUpdate(host, n) {
    const next = serverTemplates[n % serverTemplates.length].cloneNode(true);
    host.replaceChildren(...next.childNodes);
  }

  function measureInitial(factory) {
    const samples = [];
    for (let i = 0; i < INITIAL_ROUNDS; i++) {
      const host = document.createElement("div");
      document.body.append(host);
      const start = performance.now();
      const candidate = factory();
      host.append(candidate.root || candidate);
      void host.offsetHeight;
      samples.push(performance.now() - start);
      host.remove();
    }
    return samples;
  }

  function measureNativeUpdates() {
    const host = document.createElement("div");
    document.body.append(host);
    const context = nativeInitial();
    host.append(context.root);
    void host.offsetHeight;
    const samples = [];
    for (let i = 0; i < UPDATES; i++) {
      const start = performance.now();
      nativeUpdate(context, i);
      void host.offsetHeight;
      samples.push(performance.now() - start);
    }
    host.remove();
    return samples;
  }

  function measureServerUpdates() {
    const host = document.createElement("div");
    document.body.append(host);
    host.append(serverInitial());
    void host.offsetHeight;
    const samples = [];
    for (let i = 0; i < UPDATES; i++) {
      const start = performance.now();
      serverUpdate(host, i);
      void host.offsetHeight;
      samples.push(performance.now() - start);
    }
    host.remove();
    return samples;
  }

  for (let i = 0; i < 4; i++) {
    const context = nativeInitial();
    nativeUpdate(context, i);
    serverInitial();
  }

  const nativeInitialSamples = measureInitial(nativeInitial);
  const serverInitialSamples = measureInitial(serverInitial);
  const nativeUpdateSamples = measureNativeUpdates();
  const serverUpdateSamples = measureServerUpdates();

  return {
    schema_version: 1,
    workload: {
      rows: ROWS,
      initial_rounds: INITIAL_ROUNDS,
      updates: UPDATES,
      mutated_rows_per_native_update: MUTATED,
      server_fragment_rows_per_update: ROWS
    },
    candidates: {
      native_web: {
        initial_p95_ms: quantile(nativeInitialSamples, 0.95),
        initial_mean_ms: mean(nativeInitialSamples),
        update_p95_ms: quantile(nativeUpdateSamples, 0.95),
        update_mean_ms: mean(nativeUpdateSamples),
        third_party_runtime_bytes: 0
      },
      server_progressive: {
        initial_p95_ms: quantile(serverInitialSamples, 0.95),
        initial_mean_ms: mean(serverInitialSamples),
        update_p95_ms: quantile(serverUpdateSamples, 0.95),
        update_mean_ms: mean(serverUpdateSamples),
        third_party_runtime_bytes: 0
      }
    }
  };
})()
