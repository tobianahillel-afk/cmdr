package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// The local file is an explicitly trusted, immutable owner snapshot for the
// loopback demo. This is not an ingestion API or a production repository.
const maxDemoDatasetBytes int64 = 16 << 20

func loadDemoDataset(path string) (timeline.MemoryTimelineData, error) {
	var data timeline.MemoryTimelineData
	if path == "" {
		return data, errors.New("local dataset path required")
	}
	file, err := os.Open(path) // #nosec G304 -- explicit local CLI path, never HTTP input.
	if err != nil {
		return data, fmt.Errorf("open local dataset: %w", err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		return data, errors.New("dataset must be a regular file")
	}
	if stat.Size() == 0 || stat.Size() > maxDemoDatasetBytes {
		return data, errors.New("dataset size outside local 16 MiB limit")
	}
	decoder := json.NewDecoder(io.LimitReader(file, maxDemoDatasetBytes+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&data); err != nil {
		return timeline.MemoryTimelineData{}, fmt.Errorf("decode local dataset: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return timeline.MemoryTimelineData{}, errors.New("dataset must contain exactly one JSON object")
	}
	if len(data.Principals) == 0 || len(data.Cases) == 0 ||
		len(data.Entries) == 0 || len(data.Sources) == 0 {
		return timeline.MemoryTimelineData{}, errors.New("dataset requires principals, Cases, Timeline Entries and sources")
	}
	return data, nil
}

// newReloadingDemoHandler is an opt-in offline host-owned source for the
// loopback demo. It rebuilds all Case, Entry and typed-source ACL state on
// EVERY request. Invalid, missing or revoked snapshots fail closed; nothing
// falls back to a previously authorized in-memory owner.
func newReloadingDemoHandler(token, principal, path string) (http.Handler, error) {
	data, err := loadDemoDataset(path)
	if err != nil {
		return nil, err
	}
	if _, err := newDemoHandlerWithData(token, principal, data); err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		current, err := loadDemoDataset(path)
		if err == nil {
			var handler http.Handler
			handler, err = newDemoHandlerWithData(token, principal, current)
			if err == nil {
				handler.ServeHTTP(w, r)
				return
			}
		}
		// Do not expose local paths, snapshots, provider errors or stale
		// authorized data to a loopback caller, even if the source is removed.
		w.Header().Set("X-Content-Type-Options", "nosniff")
		http.Error(w, "request-unavailable", http.StatusServiceUnavailable)
	}), nil
}
