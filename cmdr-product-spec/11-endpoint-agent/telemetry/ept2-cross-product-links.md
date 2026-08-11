# Endpoint EPT-2 — Cross-Product Observation Boundaries

- **Shared:** owns `telemetry-event`, Trace/Activity/Jobs/Search/Reporting/Export and generic normalization mechanisms.
- **Investigate:** consumes telemetry/observations but owns Evidence/Finding/Case interpretation.
- **Command:** may consume health/availability projections; it does not own telemetry sources.
- **Studio:** may consume capability declarations; Endpoint capability is not a Tool.
- **Govern:** may consume technical observations; observation is not Result and EPT-2 creates no Response Run.
- **Platform Settings:** owns source/Fleet/Policy/provider/secret/tenant administration; Endpoint owns local facts.
- **Endpoint later lots:** EPT-3..EPT-6 consume EPT-2 contracts but remain NOT STARTED.
