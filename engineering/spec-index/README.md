# Product Spec generated index policy

The canonical product truth remains under `cmdr-product-spec/**`.

`cmdr-dev spec-index` generates the detailed inventory on demand. The generated inventory contains every Markdown source path, source SHA-256, front-matter identity fields and recognized CMDR references. It is intentionally **not committed** while a complete scan remains cheap.

`baseline.json` is the small committed freshness contract. It records:
- compiler schema version;
- canonical spec root;
- documentary baseline commit;
- complete tree digest;
- Markdown file count;
- active canonical document count.

CI runs:

```text
cmdr-dev spec-baseline --check
```

Any Product Spec file addition, deletion or byte-level mutation changes the tree digest and therefore fails the check until the baseline is deliberately regenerated and reviewed.

This avoids a large redundant generated artifact while preserving deterministic drift detection. If scan cost later becomes material, the policy can evolve to a content-addressed cache without changing product truth.
