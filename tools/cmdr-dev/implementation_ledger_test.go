package main

import "testing"

func TestValidateImplementationLedgerShape(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion: 1,
		LedgerKind: "verified-runtime-implementation-evidence",
		ProductSpecBaseline: state.ProductSpec.BaselineCommit,
		Entries: []ImplementationLedgerEntry{{
			Capability: "CAP-SET-004",
			WorkUnit: "E9-PILOT-001C-RUNTIME",
			ManifestPath: "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
			HandoffPath: "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary: "pilot-context-envelope-runtime",
			RuntimeRoot: "product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err != nil {
		t.Fatal(err)
	}
}

func TestImplementationLedgerRejectsDuplicateCapability(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	entry := ImplementationLedgerEntry{
		Capability: "CAP-SET-004",
		WorkUnit: "E9-PILOT-001C-RUNTIME",
		ManifestPath: "work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
		HandoffPath: "work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
		RuntimeBoundary: "pilot-context-envelope-runtime",
		RuntimeRoot: "product-runtime/context-envelope/**",
	}
	ledger := ImplementationLedger{
		SchemaVersion:1, LedgerKind:"verified-runtime-implementation-evidence",
		ProductSpecBaseline:state.ProductSpec.BaselineCommit,
		Entries:[]ImplementationLedgerEntry{entry, entry},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected duplicate capability rejection")
	}
}

func TestImplementationLedgerRejectsPathEscape(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion:1, LedgerKind:"verified-runtime-implementation-evidence",
		ProductSpecBaseline:state.ProductSpec.BaselineCommit,
		Entries:[]ImplementationLedgerEntry{{
			Capability:"CAP-SET-004", WorkUnit:"E9-PILOT-001C-RUNTIME",
			ManifestPath:"../outside.json",
			HandoffPath:"work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary:"pilot-context-envelope-runtime",
			RuntimeRoot:"product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected path escape rejection")
	}
}

func TestImplementationLedgerRejectsBaselineDrift(t *testing.T) {
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "0123456789abcdef0123456789abcdef01234567"
	ledger := ImplementationLedger{
		SchemaVersion:1, LedgerKind:"verified-runtime-implementation-evidence",
		ProductSpecBaseline:"fedcba9876543210fedcba9876543210fedcba98",
		Entries:[]ImplementationLedgerEntry{{
			Capability:"CAP-SET-004", WorkUnit:"E9-PILOT-001C-RUNTIME",
			ManifestPath:"work/lots/E9-PILOT-001C-RUNTIME/manifest.json",
			HandoffPath:"work/lots/E9-PILOT-001C-RUNTIME/HANDOFF.json",
			RuntimeBoundary:"pilot-context-envelope-runtime",
			RuntimeRoot:"product-runtime/context-envelope/**",
		}},
	}
	if err := validateImplementationLedgerShape(ledger, state); err == nil {
		t.Fatal("expected baseline drift rejection")
	}
}
