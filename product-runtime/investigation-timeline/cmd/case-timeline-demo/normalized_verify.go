package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

// Offline verification is deliberately NOT an HTTP path. It reuses the exact
// snapshot-byte digest check, owner validation, Case scoping, and imported
// source-identity checks used by the authenticated resumed viewer. A successful
// receipt describes only the selected principal's authorized imported events.
func verifySavedNormalizedCase(ctx context.Context, principal, snapshotPath, caseID, expectedDigest string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || snapshotPath == "" || caseID == "" ||
		expectedDigest == "" {
		return errors.New("offline Case verification requires saved snapshot, Case, and operator-held SHA-256")
	}
	data, ids, err := loadAuthorizedImportedCase(ctx, principal, snapshotPath, caseID, expectedDigest)
	if err != nil {
		return err
	}
	// Provenance slots overlap between independent Wazuh and generic
	// importers, so an operator's inventory must distinguish source types
	// using the already-authorized typed-source owner, not filename or text.
	sources := make(map[string]bool, len(data.Sources))
	for _, src := range data.Sources {
		sources[src.Identity.Reference] = src.Wazuh != nil
	}
	counts := make(map[string]int)
	wazuhAlerts, genericEvents := 0, 0
	for _, entry := range data.Entries {
		if _, ok := ids[entry.Value.ID]; !ok {
			continue
		}
		isWazuh, found := sources[entry.Value.SourceRef]
		if !found {
			return errors.New("authorized imported event lost its typed source")
		}
		counts[entry.Value.Provenance]++
		if isWazuh {
			wazuhAlerts++
		} else {
			genericEvents++
		}
	}
	if wazuhAlerts+genericEvents != len(ids) {
		return errors.New("authorized imported source inventory is incomplete")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		CaseID string `json:"case_id"`
		ImportedEvents int `json:"authorized_imported_events"`
		Feeds map[string]int `json:"authorized_events_by_import_slot"`
		WazuhAlerts int `json:"authorized_wazuh_alerts"`
		GenericEvents int `json:"authorized_generic_events"`
	}{
		Format: "cmdr-offline-authorized-case-verification-v1",
		Synthetic: true, CaseID: caseID,
		ImportedEvents: len(ids), Feeds: counts,
		WazuhAlerts: wazuhAlerts, GenericEvents: genericEvents,
	})
}

