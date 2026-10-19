package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

const maxDemoSnapshotBytes = 8 << 20

// loadDemoSnapshot imports operator-selected data at startup, never from HTTP.
// The full dataset is copied and cross-checked by the existing memory provider
// before the loopback listener starts. This is not a durable live repository.
func loadDemoSnapshot(path string) (timeline.MemoryTimelineData, error) {
	return loadDemoSnapshotWithExpectedSHA256(path, "")
}

// The operator-held digest must match the EXACT bytes decoded below. Hashing
// in a separate pass would leave a second-open time-of-check/use window.
// This is byte-change detection, not authenticated provenance or custody.
func loadDemoSnapshotWithExpectedSHA256(path, expectedHex string) (timeline.MemoryTimelineData, error) {
	var expected []byte
	if expectedHex != "" {
		var err error
		expected, err = hex.DecodeString(expectedHex)
		if err != nil || len(expected) != sha256.Size || path == "" {
			return timeline.MemoryTimelineData{}, errors.New("expected owner snapshot SHA-256 must be 64 hexadecimal digits")
		}
	}
	if path == "" {
		return demoData(), nil
	}
	if !filepath.IsAbs(path) {
		return timeline.MemoryTimelineData{}, errors.New("snapshot must be an absolute local path")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxDemoSnapshotBytes {
		return timeline.MemoryTimelineData{}, errors.New("snapshot must be an existing bounded regular file (not a symlink)")
	}
	file, err := os.Open(path) // #nosec G304 -- explicit trusted local operator CLI only, never HTTP input.
	if err != nil {
		return timeline.MemoryTimelineData{}, fmt.Errorf("cannot open operator snapshot: %w", err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return timeline.MemoryTimelineData{}, errors.New("snapshot file changed during opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxDemoSnapshotBytes+1))
	if err != nil || len(data) > maxDemoSnapshotBytes {
		return timeline.MemoryTimelineData{}, errors.New("snapshot exceeds size limit or cannot be read")
	}
	if expected != nil {
		sum := sha256.Sum256(data)
		if !bytes.Equal(sum[:], expected) {
			return timeline.MemoryTimelineData{}, errors.New("owner snapshot bytes differ from operator-held SHA-256")
		}
	}
	var snapshot timeline.MemoryTimelineData
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snapshot); err != nil {
		return timeline.MemoryTimelineData{}, fmt.Errorf("invalid snapshot JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return timeline.MemoryTimelineData{}, errors.New("snapshot must contain one JSON object")
	}
	return snapshot, nil
}

// Choose a Case already visible to the selected synthetic principal, so the
// CLI's advertised URL never reveals or points to another tenant's Case.
func firstDemoAuthorizedCase(data timeline.MemoryTimelineData, principal string) (string, error) {
	var tenant, environment string
	for _, p := range data.Principals {
		if p.ID == principal {
			tenant, environment = p.Tenant, p.Environment
			break
		}
	}
	if tenant == "" || environment == "" {
		return "", errors.New("principal absent from operator snapshot")
	}
	var caseID string
	for _, c := range data.Cases {
		if c.Tenant != tenant || c.Environment != environment {
			continue
		}
		for _, reader := range c.Readers {
			if reader == principal && (caseID == "" || c.ID < caseID) {
				caseID = c.ID
			}
		}
	}
	if caseID == "" {
		return "", errors.New("no readable Case in operator snapshot")
	}
	return caseID, nil
}
