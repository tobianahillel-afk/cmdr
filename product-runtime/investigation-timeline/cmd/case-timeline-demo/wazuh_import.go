package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

// Wazuh 4.x alerts.json is an alert stream, not a trusted Case/tenant/ACL
// catalog. Only selected observational fields are imported. Unknown Wazuh
// metadata cannot grant access; the existing Case owner validates all reads.
type wazuh4Alert struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	FullLog   string `json:"full_log"`
	Rule      struct {
		Description string `json:"description"`
	} `json:"rule"`
}

func normalizeWazuh4Alerts(ctx context.Context, raw []byte) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(raw) == 0 || len(raw) > maxDemoSnapshotBytes || !utf8.Valid(raw) {
		return nil, errors.New("Wazuh alerts must be bounded valid UTF-8 JSONL")
	}
	lines := bytes.Split(raw, []byte("\n"))
	if len(lines) > 5001 {
		return nil, errors.New("too many Wazuh alerts")
	}
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	count := 0
	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if i == len(lines)-1 && len(line) == 0 {
			continue
		}
		if len(line) == 0 || len(line) > 8192 {
			return nil, fmt.Errorf("invalid Wazuh alert line %d", i+1)
		}
		var alert wazuh4Alert
		if err := json.Unmarshal(line, &alert); err != nil {
			return nil, fmt.Errorf("invalid Wazuh alert line %d", i+1)
		}
		if alert.ID == "" {
			return nil, fmt.Errorf("missing Wazuh alert identity at line %d", i+1)
		}
		eventID := "wazuh-" + alert.ID
		if !validLocalEventID(eventID) {
			return nil, fmt.Errorf("invalid Wazuh alert identity at line %d", i+1)
		}
		when, err := time.Parse(time.RFC3339Nano, alert.Timestamp)
		if err != nil {
			when, err = time.Parse("2006-01-02T15:04:05.999999999-0700", alert.Timestamp)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid Wazuh alert time at line %d", i+1)
		}
		// A rule description is a bounded analyst-facing summary, never
		// interpreted as an instruction or as a source authorization.
		message := alert.Rule.Description
		if message == "" {
			message = alert.FullLog
		}
		if len(message) == 0 || len(message) > 512 || !utf8.ValidString(message) ||
			strings.ContainsAny(message, "\x00\r\n") {
			return nil, fmt.Errorf("invalid Wazuh preview at line %d", i+1)
		}
		if err := encoder.Encode(normalizedLocalEvent{
			ID: eventID, Timestamp: when.Format(time.RFC3339Nano), Message: message,
		}); err != nil {
			return nil, fmt.Errorf("cannot normalize Wazuh alert at line %d", i+1)
		}
		count++
		if out.Len() > maxDemoSnapshotBytes {
			return nil, errors.New("normalized Wazuh alerts exceed size limit")
		}
	}
	if count == 0 {
		return nil, errors.New("no Wazuh alerts found")
	}
	return out.Bytes(), nil
}

// This is an offline lab integration: one Wazuh alerts.json file flows
// straight to the already-authorized in-memory Case and loopback UI.
// No merged owner graph is exported and no new ACL service is introduced.
func serveDemoWithWazuhAlerts(ctx context.Context, principal, snapshotPath, caseID, alertsPath string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if snapshotPath == "" || alertsPath == "" || snapshotPath == alertsPath {
		return errors.New("select two distinct absolute local owner and Wazuh files")
	}
	owner, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	raw, err := readLocalEventFile(alertsPath)
	if err != nil {
		return err
	}
	events, err := normalizeWazuh4Alerts(ctx, raw)
	if err != nil {
		return err
	}
	merged, err := importLocalEvents(ctx, owner, principal, caseID, events)
	if err != nil {
		return err
	}
	return serveDemoWithData(ctx, principal, merged, output)
}

// writeDemoWazuh4Proximity turns a local Wazuh alert stream into time leads
// across two independently authorized Cases. It never exports owner data,
// does not interpret Wazuh agent/rule metadata as permissions, and emits no
// report until both Case/typed-source traversals have completed.
func writeDemoWazuh4Proximity(ctx context.Context, principal, snapshotPath, intoCase, withCase, alertsPath string, seconds int, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || intoCase == "" || withCase == "" || intoCase == withCase ||
		len(intoCase) > 1024 || len(withCase) > 1024 ||
		snapshotPath == "" || alertsPath == "" || snapshotPath == alertsPath ||
		seconds < 1 || seconds > 300 {
		return errors.New("select two distinct Cases, separate local files and a 1..300-second window")
	}
	owner, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	raw, err := readLocalEventFile(alertsPath)
	if err != nil {
		return err
	}
	events, err := normalizeWazuh4Alerts(ctx, raw)
	if err != nil {
		return err
	}
	merged, err := importLocalEvents(ctx, owner, principal, intoCase, events)
	if err != nil {
		return err
	}
	imported := len(merged.Entries) - len(owner.Entries)
	if imported < 1 || imported > 5000 {
		return errors.New("inconsistent Wazuh import count")
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, merged)
	if err != nil {
		return err
	}
	both, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, intoCase, withCase)
	if err != nil {
		return err
	}
	leads, err := selectCaseProximity(ctx, both, seconds)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localImportedProximity{
		Format: "cmdr-local-wazuh4-case-proximity-v1",
		Synthetic: true,
		ImportedEntries: imported,
		Proximity: leads,
	})
}
