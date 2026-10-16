# E10-INV-120A — multi-Case offline end-to-end investigation

A synthetic investigator-42 has two readable Cases in tenant A; investigator-99 has one isolated Case in tenant B. Run the existing CLI with the absolute `-snapshot` path and `-principal=investigator-42`. The startup URL opens the first Case; the existing authenticated analyst picker/HTTP Case route opens the second `case-incident-43`, its ordered Timeline `alert-3`, navigation `event-3` and bounded source preview.

Real loopback tests prove Case selection across two Cases, source drill-down, foreign tenant denial, missing bearer rejection and clean shutdown; negative startup tests remove ALL readable Cases from a principal to verify fail closed. Operator-supplied JSON ACLs are synthetic, not authoritative IAM or durable persistence. OPEN-013 remains unresolved. No production-level readiness/CI/L2 claim until exact evidence.
