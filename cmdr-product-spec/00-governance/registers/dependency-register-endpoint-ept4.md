---
id: dependency-register-endpoint-ept4
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-4

| Family | CAP-EPT range | Required sources/owners | Boundary |
|---|---|---|---|
| Request/authority | 047 | Investigate Collection Request; Govern authority; Settings Policy | request/eligible/authorized/start remain distinct |
| Scope/plan | 048 | CAP-EPT-047; EPT-3 refs; Settings limits | plan cannot expand request |
| File collection | 049 | Endpoint file source; CAP-INV-205 | neutral Collection Item; OPEN-014 |
| Process/system snapshot | 050 | CAP-INV-206; EPT-3 contexts | fresh read-only acquisition; no mutation |
| Memory | 051 | Endpoint memory source; CAP-INV-207 | no tool/format; acquired != analyzed |
| Network capture | 052 | CAP-INV-208; Endpoint network capability | bounded capture; no containment |
| Collection state | 053 | collection-queue; CAP-INV-203; Shared Jobs | local state != business/generic Job |
| Package/integrity | 054 | triage-package; chain-of-custody; CAP-INV-213 | metadata != crypto proof; no Evidence package |
| Transfer/handoff | 055 | CAP-INV-211/213/214; Shared Jobs/Export | transfer != qualification |
| LR session | 056–057 | CAP-INV-209; Endpoint session source; Govern | technical session != business session/Response Run |
| Command/shell | 058–059 | command/terminal sources; CAP-INV-210 | request != execution; no shell/runtime selected |
| Script | 060 | script-execution; Settings Secrets; Govern/Studio | no script/runtime; effectful authority |
| Session file operations | 061 | file-actions boundary; CAP-INV-211 | temp upload Class 3; destructive actions deferred |
| Output/control | 062–063 | CAP-INV-212/214; local audit; Govern | technical output != Result; closure != restoration |
| Cross-product provenance | 064 | Investigate/Govern/Studio/Settings/Shared | OPEN-014/015 preserved |

All dependencies are documentary. OPEN-008/013/014/015/017 remain open; no provider, protocol, engine or final object identity is selected.