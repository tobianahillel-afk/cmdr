package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
)

// All three operator-local investigation receipts are already independently
// authorized and rendered by their existing first-party workflows. Preserve
// those checks, then atomically publish the WHOLE JSON report as a private,
// new-only file; never stream an incomplete report to stdout or overwrite
// an existing source/snapshot. No new storage engine or sidecar is needed.
func writeAuthorizedReceiptPrivate(ctx context.Context, destination string, produce func(io.Writer) error) error {
	if ctx == nil || ctx.Err() != nil || !filepath.IsAbs(destination) || produce == nil {
		return errors.New("private investigation report requires an absolute new output path")
	}
	var receipt bytes.Buffer
	if err := produce(&receipt); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if receipt.Len() == 0 || receipt.Len() > maxDemoSnapshotBytes || !json.Valid(receipt.Bytes()) {
		return errors.New("authorized investigation report is not a complete bounded JSON document")
	}
	// Existing owner-snapshot publisher uses an OS-private temporary file,
	// sync and atomic hard-link publication without replacing a path/symlink.
	// It deliberately fails closed if that filesystem cannot make hard links.
	return writePrivateSnapshotAtomic(ctx, destination, receipt.Bytes())
}
