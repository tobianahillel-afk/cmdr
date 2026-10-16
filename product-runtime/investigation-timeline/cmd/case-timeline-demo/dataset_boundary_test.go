package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportedOwnerFileMustExistAndContainRecords(t *testing.T) {
	if _, err := loadDemoDataset(""); err == nil {
		t.Fatal("empty dataset path accepted")
	}
	if _, err := loadDemoDataset(filepath.Join(t.TempDir(), "absent.json")); err == nil {
		t.Fatal("missing dataset file accepted")
	}
	empty := filepath.Join(t.TempDir(), "empty.json")
	if err := os.WriteFile(empty, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadDemoDataset(empty); err == nil {
		t.Fatal("empty local dataset accepted")
	}
	for _, tc := range []struct {
		name, source string
	}{
		{"empty-object", "{}"},
		{"all-empty-arrays", `{"Principals":[],"Cases":[],"Entries":[],"Sources":[]}`},
		{"only-principals", `{"Principals":[{"ID":"owner","Tenant":"tenant","Environment":"lab"}]}`},
		{"trailing-whitespace-with-second-object", "{}  {}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "dataset.json")
			if err := os.WriteFile(path, []byte(tc.source), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := loadDemoDataset(path); err == nil {
				t.Fatal("incomplete or ambiguous Case dataset accepted")
			}
		})
	}
}

func TestImportedDatasetNeverOverridesOwnerAuthorization(t *testing.T) {
	data := importedCaseDataset()
	data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, "event-b")
	path := writeDemoDataset(t, data)
	loaded, err := loadDemoDataset(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := newDemoHandlerWithData(testBearer, "investigator-42", loaded); err == nil {
		t.Fatal("Case catalogue referencing foreign tenant entry was accepted")
	}
	var output bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := runDemoCLI(ctx, []string{
		"-principal=investigator-42", "-dataset=" + path,
	}, &output); err == nil || output.Len() != 0 {
		t.Fatalf("invalid owner data started a listener or disclosed a token: %v %q", err, output.String())
	}
}

func TestDatasetCLIRejectsInvalidPathWithoutToken(t *testing.T) {
	var output bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := runDemoCLI(ctx, []string{
		"-principal=investigator-42",
		"-dataset=" + filepath.Join(t.TempDir(), "not-there.json"),
	}, &output); err == nil || output.Len() != 0 {
		t.Fatalf("failed dataset read reached authenticated listener: %v %q", err, output.String())
	}
}

func TestImportedSnapshotHasExplicitAuthorizedCase(t *testing.T) {
	data := importedCaseDataset()
	data.Cases[0].Readers = nil
	path := writeDemoDataset(t, data)
	loaded, err := loadDemoDataset(path)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := serveDemoWithData(ctx, "investigator-42", &output, loaded); err == nil ||
		!strings.Contains(err.Error(), "no selectable Case") || output.Len() != 0 {
		t.Fatalf("unauthorized principal Case was served: %v %q", err, output.String())
	}
}
