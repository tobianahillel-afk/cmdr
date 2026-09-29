# cmdr-dev

Minimal CMDR development kernel.

Properties:
- Go standard library only;
- no network access required;
- cross-platform;
- strict JSON decoding;
- repository state is validated before it is trusted.

## Commands

From any directory inside the repository:

```text
go run ./tools/cmdr-dev doctor
go run ./tools/cmdr-dev status
go run ./tools/cmdr-dev next
```

Use `--json` for machine-readable output and `--root <path>` to override repository discovery.

This tool is not the CMDR product runtime. It is engineering infrastructure.
