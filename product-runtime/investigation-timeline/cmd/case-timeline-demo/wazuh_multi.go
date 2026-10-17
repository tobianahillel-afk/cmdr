package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Wazuh 4.x is parsed only by the existing validated normalizer. This small
// adapter adds a trusted operator-side feed label to the normalized event ID,
// before passing it to the unchanged owner-authorized JSONL importer.
func namespaceWazuhEvents(ctx context.Context, normalized []byte, feed int) ([]byte, error) {
	if feed < 1 || feed > maxWazuhAlertFiles {
		return nil, errors.New("invalid Wazuh feed number")
	}
	var result bytes.Buffer
	lines := bytes.Split(normalized, []byte("\n"))
	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if i == len(lines)-1 && len(line) == 0 {
			continue
		}
		var event normalizedLocalEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, errors.New("invalid normalized Wazuh event")
		}
		event.ID = fmt.Sprintf("feed-%d-%s", feed, event.ID)
		if !validLocalEventID(event.ID) {
			return nil, errors.New("Wazuh feed identity exceeds supported length")
		}
		if err := json.NewEncoder(&result).Encode(event); err != nil {
			return nil, err
		}
		if result.Len() > maxDemoSnapshotBytes {
			return nil, errors.New("multi-feed normalized Wazuh data too large")
		}
	}
	return result.Bytes(), nil
}

// All input files and access decisions are validated before issuing a bearer
// or starting the existing read-only loopback Case/Timeline/source viewer.
func serveDemoWithWazuhAlertsFiles(ctx context.Context, principal, snapshotPath, caseID string, paths []string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(paths) == 1 {
		return serveDemoWithWazuhAlerts(ctx, principal, snapshotPath, caseID, paths[0], output)
	}
	if snapshotPath == "" || len(paths) < 2 || len(paths) > maxWazuhAlertFiles {
		return errors.New("select one owner snapshot and up to four Wazuh alert files")
	}
	owner, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	info, err := os.Lstat(snapshotPath)
	if err != nil {
		return errors.New("cannot identify owner snapshot")
	}
	seen := []os.FileInfo{info}
	for i, path := range paths {
		if err := ctx.Err(); err != nil {
			return err
		}
		data, err := readLocalEventFile(path)
		if err != nil {
			return err
		}
		current, err := os.Lstat(path)
		if err != nil {
			return errors.New("cannot identify Wazuh alert file")
		}
		for _, previous := range seen {
			if os.SameFile(current, previous) {
				return errors.New("Wazuh alert inputs must be distinct files")
			}
		}
		seen = append(seen, current)
		events, err := normalizeWazuh4Alerts(ctx, data)
		if err != nil {
			return err
		}
		events, err = namespaceWazuhEvents(ctx, events, i+1)
		if err != nil {
			return err
		}
		owner, err = importLocalEvents(ctx, owner, principal, caseID, events)
		if err != nil {
			return err
		}
	}
	return serveDemoWithData(ctx, principal, owner, output)
}
