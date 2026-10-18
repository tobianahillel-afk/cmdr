// Command case-timeline-demo serves a loopback-only, synthetic Case -> Timeline
// integration. It is deliberately not a production identity or owner service.
package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runDemoCLI(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CLI parsing is separate from process exit so real requests and rejected
// arguments can be exercised in tests without a shell or subprocess runner.
const maxWazuhAlertFiles = 4

type wazuhAlertFiles []string

func (files *wazuhAlertFiles) String() string {
	return fmt.Sprintf("%d local Wazuh alert files", len(*files))
}

func (files *wazuhAlertFiles) Set(path string) error {
	if path == "" || len(*files) >= maxWazuhAlertFiles {
		return errors.New("select 1 to 4 Wazuh alerts JSONL files")
	}
	*files = append(*files, path)
	return nil
}

// Extra normalized sources are enabled only in the existing import-and-serve
// mode; other snapshot export/triage workflows keep their single-file contract.
type normalizedExtraFiles []string

func (files *normalizedExtraFiles) String() string {
	return fmt.Sprintf("%d additional normalized event files", len(*files))
}

func (files *normalizedExtraFiles) Set(path string) error {
	if path == "" || len(*files) >= 3 {
		return errors.New("select up to three additional normalized JSONL feeds")
	}
	*files = append(*files, path)
	return nil
}

func runDemoCLI(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("case-timeline-demo", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	principal := flags.String("principal", "analyst-a", "synthetic analyst ID from the selected local owner snapshot")
	snapshot := flags.String("snapshot", "", "absolute path to operator-owned read-only Case snapshot JSON (loopback demo only)")
	importJSONL := flags.String("import-jsonl", "", "absolute local normalized event JSONL for offline owner snapshot preparation (no listener)")
	var extraJSONL normalizedExtraFiles
	flags.Var(&extraJSONL, "import-jsonl-extra", "additional absolute normalized JSONL file for -serve-import; repeat at most 3 times")
	var wazuhFiles wazuhAlertFiles
	flags.Var(&wazuhFiles, "wazuh-alerts", "absolute local Wazuh 4.x alerts.json JSONL; repeat up to four times")
	intoCase := flags.String("into-case", "", "existing Case for imported local events (requires -import-jsonl)")
	triageImport := flags.Bool("triage-import", false, "analyze locally imported events; use -with-case for a two-Case comparison, optionally -compare-proximity for bounded time leads")
	serveImport := flags.Bool("serve-import", false, "import local JSONL into an in-memory Case and open the authenticated loopback analyst viewer; no merged owner file")
	saveImportSnapshot := flags.String("save-import-snapshot", "", "absolute new path for private normalized Case snapshot; explicit import only, no listener or overwrite")
	saveWazuhCaseSnapshot := flags.String("save-wazuh-case-snapshot", "", "explicit new private scoped Case snapshot from 1..4 validated Wazuh 4.x local alert files; no listener")
	saveMixedCaseSnapshot := flags.String("save-mixed-case-snapshot", "", "create one atomic private Case archive with both Wazuh and generic JSONL sources; no intermediary archive")
	resumeImportSearchCase := flags.String("resume-import-search-case", "", "Case ID in a saved normalized owner snapshot to reopen with an authenticated Search Job viewer")
	verifyImportSnapshotCase := flags.String("verify-import-snapshot-case", "", "Case ID of a saved authorized owner snapshot to verify WITHOUT a listener; requires -expect-snapshot-sha256")
	privateReport := flags.String("private-report", "", "optional absolute NEW private file for a completed authorized verification, headless hunt or cross-snapshot comparison receipt")
	huntSavedCase := flags.String("hunt-saved-case", "", "run one authorized Event Search Job on a saved Case without listener or browser")
	huntField := flags.String("hunt-field", "preview", "first-party search field (preview, rule_id, agent_id, entry_id, source_ref, provenance)")
	huntFamily := flags.String("hunt-family", "all", "source family for offline search (all, wazuh, generic)")
	huntValue := flags.String("hunt-value", "", "literal Event Search Job value, passed in-memory only")
	huntFrom := flags.String("hunt-from", "", "inclusive RFC3339 start for offline authorized Event Search Job")
	huntTo := flags.String("hunt-to", "", "exclusive RFC3339 end for offline authorized Event Search Job")
	inspectSavedCase := flags.String("inspect-saved-case", "", "explicit saved Case to privately inspect a correlated Entry (requires snapshot, SHA-256 and private-report)")
	inspectSavedEntry := flags.String("inspect-saved-entry", "", "exact imported Entry ID from a cross-Case correlation receipt; not an authorization grant")
	inspectSavedPair := flags.Int("inspect-saved-pair", 0, "1-based index (1..10) of an authorized Wazuh cross-Case pair; requires two verified Case snapshots and private-report")
	includePairGenericEvidence := flags.Bool("include-pair-generic-evidence", false, "explicitly include up to ten independently source-authorized generic previews per Case in the NEW private paired report only")
	compareSavedCase := flags.String("compare-saved-case", "", "first saved authorized Case ID for independent offline cross-snapshot comparison")
	compareSavedWithCase := flags.String("compare-saved-with-case", "", "second distinct saved Case ID for independent cross-snapshot comparison")
	compareSavedWithSnapshot := flags.String("compare-saved-with-snapshot", "", "absolute path to second independent saved Case snapshot")
	expectedOtherSHA256 := flags.String("expect-other-snapshot-sha256", "", "independently retained SHA-256 of second saved Case snapshot")
	expectedSnapshotSHA256 := flags.String("expect-snapshot-sha256", "", "optional 64-digit SHA-256 recorded at explicit Case save; verified against exact bytes on resume only")
	reportCase := flags.String("report-case", "", "produce one authorized offline Case triage JSON receipt and exit (no listener)")
	exportCase := flags.String("export-case", "", "export authorized Case Timeline events as sensitive local JSONL, without listener")
	compareCase := flags.String("compare-case", "", "first authorized Case for an offline two-Case comparison receipt")
	withCase := flags.String("with-case", "", "second authorized Case for the offline comparison (requires -compare-case)")
	compareEvents := flags.Bool("compare-events", false, "explicit sensitive chronological event-level cross-Case JSON receipt")
	compareProximity := flags.Bool("compare-proximity", false, "explicit nearest certain-time event leads across two authorized Cases")
	proximitySeconds := flags.Int("proximity-seconds", 60, "bounded nearest-event window in seconds, 1..300 (requires -compare-proximity)")
	sourceCase := flags.String("source-case", "", "authorized Case for offline source-specific event timeline")
	sourceRef := flags.String("source-ref", "", "exact authorized source reference (requires -source-case)")
	inspectCase := flags.String("inspect-case", "", "Case containing an explicitly selected source Timeline Entry")
	inspectEntry := flags.String("inspect-entry", "", "Timeline Entry ID for an authorized offline typed-source preview (requires -inspect-case)")
	inspectWazuh := flags.String("inspect-wazuh-entry", "", "selected imported Wazuh Timeline Entry for full first-party Event Inspection")
	pivotWazuh := flags.String("pivot-wazuh-field", "", "optional wazuh.rule.id or wazuh.agent.id search pivot draft")
	pivotFrom := flags.String("pivot-start", "", "inclusive RFC3339 Event Search pivot start")
	pivotTo := flags.String("pivot-end", "", "exclusive RFC3339 Event Search pivot end")
	executePivot := flags.Bool("execute-wazuh-pivot", false, "execute the selected first-party search pivot on authorized offline Wazuh alerts; bounded, no live provider")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional arguments")
	}
	// Distinguish omitted CLI modes from options explicitly set to an
	// empty string (commonly via shell variable expansion). In particular,
	// an empty Case export/reopen/inspection selector MUST NOT silently
	// fall through to the default loopback viewer or legacy stdout export.
	// This check belongs at the one CLI entry boundary, before any
	// snapshot read, listener, bearer, source export or private file write.
	emptyExplicitMode := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "private-report", "save-import-snapshot", "save-wazuh-case-snapshot",
			"save-mixed-case-snapshot", "expect-snapshot-sha256",
			"expect-other-snapshot-sha256", "resume-import-search-case",
			"verify-import-snapshot-case", "hunt-saved-case",
			"compare-saved-case", "compare-saved-with-case",
			"compare-saved-with-snapshot", "export-case", "report-case",
			"compare-case", "with-case", "source-case", "source-ref",
			"inspect-case", "inspect-entry", "inspect-wazuh-entry",
			"inspect-saved-case", "inspect-saved-entry", "snapshot":
			if f.Value.String() == "" { emptyExplicitMode = true }
		}
	})
	if emptyExplicitMode {
		return errors.New("explicit investigation, private output or snapshot integrity option cannot be empty")
	}
	if *privateReport != "" {
		modes := 0
		for _, selected := range []string{*huntSavedCase, *compareSavedCase, *verifyImportSnapshotCase, *inspectSavedCase} {
			if selected != "" { modes++ }
		}
		if modes != 1 {
			return errors.New("private-report requires exactly one authorized offline hunt, saved-Case comparison, verification or private source inspection")
		}
	}
	// Reuse the SAME authenticated owner and Event Search backend used by
	// the analyst UI, but produce an explicit minimal offline JSON receipt.
	// No local listener, bearer in stdout, independent search implementation
	// or trust of the operator's source/tenant permission claims.
	headlessSelected := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "hunt-saved-case", "hunt-field", "hunt-family", "hunt-value",
			"hunt-from", "hunt-to":
			headlessSelected = true
		}
	})
	if headlessSelected {
		allowed := true
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot", "expect-snapshot-sha256",
				"hunt-saved-case", "hunt-field", "hunt-family",
				"hunt-value", "hunt-from", "hunt-to", "private-report":
			default:
				allowed = false
			}
		})
		if !allowed || *huntSavedCase == "" || *snapshot == "" ||
			*expectedSnapshotSHA256 == "" || *huntValue == "" ||
			*huntFrom == "" || *huntTo == "" {
			return errors.New("headless Case hunt needs an independent snapshot fingerprint, authorized Case, literal and time window")
		}
		if *privateReport != "" {
			return writeAuthorizedReceiptPrivate(ctx, *privateReport, func(w io.Writer) error {
				return runSavedCaseHeadlessHunt(ctx, *principal, *snapshot, *huntSavedCase,
					*expectedSnapshotSHA256, *huntField, *huntFamily,
					*huntValue, *huntFrom, *huntTo, w)
			})
		}
		return runSavedCaseHeadlessHunt(ctx, *principal, *snapshot, *huntSavedCase,
			*expectedSnapshotSHA256, *huntField, *huntFamily,
			*huntValue, *huntFrom, *huntTo, output)
	}
	// Build a mixed Wazuh + generic JSONL Case in ONE owner-validated,
	// no-overwrite filesystem transaction, avoiding an intermediate
	// sensitive Wazuh-only archive. No online listener is started.
	if *saveMixedCaseSnapshot != "" {
		allowed := true
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot", "into-case", "wazuh-alerts",
				"import-jsonl", "import-jsonl-extra", "save-mixed-case-snapshot":
			default:
				allowed = false
			}
		})
		if !allowed || *snapshot == "" || *intoCase == "" || *importJSONL == "" ||
			len(wazuhFiles) == 0 {
			return errors.New("mixed private Case save requires one Case, Wazuh and normalized JSONL inputs only")
		}
		normalized := append([]string{*importJSONL}, extraJSONL...)
		return writePersistedMixedCase(ctx, *principal, *snapshot, *intoCase,
			[]string(wazuhFiles), normalized, *saveMixedCaseSnapshot, output)
	}
	// Bridge a real Wazuh import into the existing private, source-authorized
	// Case persistence. This is a separate explicit mode; no other local
	// report, online viewer, or saved-Case command can be combined with it.
	if *saveWazuhCaseSnapshot != "" {
		allowed := true
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot", "into-case", "wazuh-alerts",
				"save-wazuh-case-snapshot":
			default:
				allowed = false
			}
		})
		if !allowed || *snapshot == "" || *intoCase == "" || len(wazuhFiles) == 0 {
			return errors.New("Wazuh private Case save requires owner snapshot, selected Case and 1..4 Wazuh files only")
		}
		return writePersistedWazuhCase(ctx, *principal, *snapshot, *intoCase,
			[]string(wazuhFiles), *saveWazuhCaseSnapshot, output)
	}
	// An Entry ID in a saved cross-Case Wazuh correlation is a navigational
	// hint only. To inspect it, select exactly ONE independently authorized
	// saved owner, supply its operator-held fingerprint and write the
	// sensitive source preview only to a new private destination.
	savedEntrySelected := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "inspect-saved-case" || f.Name == "inspect-saved-entry" {
			savedEntrySelected = true
		}
	})
	if savedEntrySelected {
		allowed := true
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot", "inspect-saved-case",
				"inspect-saved-entry", "expect-snapshot-sha256",
				"private-report":
			default:
				allowed = false
			}
		})
		if !allowed || *snapshot == "" || *inspectSavedCase == "" ||
			*inspectSavedEntry == "" || *expectedSnapshotSHA256 == "" ||
			*privateReport == "" {
			return errors.New("private saved-entry inspection requires authorized Case, Entry, SHA-256 and private-report only")
		}
		return writeVerifiedSavedCaseEntryPrivate(ctx, *principal, *snapshot,
			*inspectSavedCase, *expectedSnapshotSHA256, *inspectSavedEntry,
			*privateReport)
	}
	// Separate, strictly scoped read-only comparison of two independently
	// persisted Cases. Never join their owner graphs or allow cross-tenant
	// source access through the comparison output.
	savedComparisonSelected := false
	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "compare-saved-case", "compare-saved-with-case",
			"compare-saved-with-snapshot", "expect-other-snapshot-sha256",
			"inspect-saved-pair", "include-pair-generic-evidence":
			savedComparisonSelected = true
		}
	})
	if savedComparisonSelected {
		allowed := true
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot", "compare-saved-case",
				"compare-saved-with-case", "compare-saved-with-snapshot",
				"expect-snapshot-sha256", "expect-other-snapshot-sha256",
				"private-report", "inspect-saved-pair", "include-pair-generic-evidence":
			default:
				allowed = false
			}
		})
		if !allowed || *snapshot == "" || *compareSavedCase == "" ||
			*compareSavedWithCase == "" || *compareSavedWithSnapshot == "" ||
			*expectedSnapshotSHA256 == "" || *expectedOtherSHA256 == "" {
			return errors.New("offline saved-Case comparison requires two Cases, snapshots and independently held SHA-256 receipts")
		}
		pairRequested := false
		flags.Visit(func(f *flag.Flag) {
			if f.Name == "inspect-saved-pair" { pairRequested = true }
		})
		if pairRequested {
			if *inspectSavedPair < 1 || *inspectSavedPair > 10 || *privateReport == "" {
				return errors.New("private pair inspection requires a 1..10 pair index and explicit new private-report path")
			}
			return writeVerifiedCorrelatedPairPrivateWithGenericEvidence(ctx, *principal, *snapshot,
				*compareSavedCase, *expectedSnapshotSHA256,
				*compareSavedWithSnapshot, *compareSavedWithCase,
				*expectedOtherSHA256, *inspectSavedPair, *privateReport, *includePairGenericEvidence)
		}
		// Detailed generic source previews require a concrete correlated
		// pair and a NEW private destination. Never enable preview export
		// for a metadata-only comparison or allow a stdout fallback.
		if *includePairGenericEvidence {
			return errors.New("generic evidence inclusion requires an explicit inspected Wazuh pair")
		}
		if *privateReport != "" {
			return writeAuthorizedReceiptPrivate(ctx, *privateReport, func(w io.Writer) error {
				return compareVerifiedSavedCases(ctx, *principal, *snapshot,
					*compareSavedCase, *expectedSnapshotSHA256,
					*compareSavedWithSnapshot, *compareSavedWithCase,
					*expectedOtherSHA256, w)
			})
		}
		return compareVerifiedSavedCases(ctx, *principal, *snapshot,
			*compareSavedCase, *expectedSnapshotSHA256,
			*compareSavedWithSnapshot, *compareSavedWithCase,
			*expectedOtherSHA256, output)
	}
	if *saveImportSnapshot != "" || *resumeImportSearchCase != "" || *verifyImportSnapshotCase != "" {
		modes := 0
		for _, chosen := range []string{*saveImportSnapshot, *resumeImportSearchCase, *verifyImportSnapshotCase} {
			if chosen != "" { modes++ }
		}
		if modes != 1 {
			return errors.New("select exactly one Case save, resume, or offline verify mode")
		}
		// Explicit modes are intentionally narrow. Never accidentally
		// combine sensitive disk writes with a Wazuh, Case export, viewer,
		// or any unrelated command mode.
		invalid := false
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "principal", "snapshot":
			case "save-import-snapshot", "into-case", "import-jsonl", "import-jsonl-extra":
				if *saveImportSnapshot == "" { invalid = true }
			case "resume-import-search-case":
				if *resumeImportSearchCase == "" { invalid = true }
			case "verify-import-snapshot-case":
				if *verifyImportSnapshotCase == "" { invalid = true }
			case "private-report":
				if *verifyImportSnapshotCase == "" { invalid = true }
			case "expect-snapshot-sha256":
				if *resumeImportSearchCase == "" && *verifyImportSnapshotCase == "" { invalid = true }
			default:
				invalid = true
			}
		})
		if invalid || *snapshot == "" {
			return errors.New("private snapshot save/resume cannot combine with other investigation modes")
		}
		if *saveImportSnapshot != "" {
			if *intoCase == "" || *importJSONL == "" {
				return errors.New("snapshot save needs existing Case and first normalized source")
			}
			paths := append([]string{*importJSONL}, extraJSONL...)
			return writeDemoPersistedNormalizedCase(ctx, *principal, *snapshot, *intoCase, paths, *saveImportSnapshot, output)
		}
		if *verifyImportSnapshotCase != "" {
			if *privateReport != "" {
				return writeAuthorizedReceiptPrivate(ctx, *privateReport, func(w io.Writer) error {
					return verifySavedNormalizedCase(ctx, *principal, *snapshot, *verifyImportSnapshotCase, *expectedSnapshotSHA256, w)
				})
			}
			return verifySavedNormalizedCase(ctx, *principal, *snapshot, *verifyImportSnapshotCase, *expectedSnapshotSHA256, output)
		}
		return serveDemoReopenedNormalizedCase(ctx, *principal, *snapshot, *resumeImportSearchCase, *expectedSnapshotSHA256, output)
	}

	if *expectedSnapshotSHA256 != "" {
		return errors.New("expected snapshot SHA-256 requires an explicit resume-import-search-case")
	}
	windowSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "proximity-seconds" {
			windowSet = true
		}
	})
	if windowSet && !*compareProximity {
		return errors.New("proximity window requires comparison proximity mode")
	}
	if *compareProximity && *compareEvents {
		return errors.New("choose exactly one cross-Case comparison mode")
	}
	if len(extraJSONL) > 0 &&
		(!*serveImport || *importJSONL == "" || *snapshot == "" || *intoCase == "" ||
			*triageImport || len(wazuhFiles) > 0) {
		return errors.New("additional normalized feeds require only -serve-import with existing Case, snapshot and first normalized file")
	}
	if *executePivot && (*inspectWazuh == "" || *pivotWazuh == "") {
		return errors.New("executed Wazuh search requires an explicit Event Inspection pivot")
	}
	if *inspectWazuh == "" && (*pivotWazuh != "" || *pivotFrom != "" || *pivotTo != "") {
		return errors.New("Wazuh pivot requires an imported Event Inspection selection")
	}
	if *inspectWazuh != "" && len(wazuhFiles) == 0 {
		return errors.New("Wazuh Event Inspection requires Wazuh alert files")
	}
	// Wazuh 4.x is an offline source, not an ACL provider. Require exactly
	// one of the existing viewer or authorized two-Case temporal paths.
	if len(wazuhFiles) > 0 {
		if *inspectWazuh != "" {
			if *snapshot == "" || *intoCase == "" || *importJSONL != "" ||
				*reportCase != "" || *exportCase != "" || *compareCase != "" ||
				*compareEvents || *compareProximity || windowSet ||
				*sourceCase != "" || *sourceRef != "" || *inspectCase != "" ||
				*inspectEntry != "" || *serveImport || *triageImport || *withCase != "" {
				return errors.New("Wazuh Event Inspection cannot combine with other investigation modes")
			}
			return writeDemoWazuhEventInspectionWithSearch(ctx, *principal, *snapshot, *intoCase,
				[]string(wazuhFiles), *inspectWazuh, *pivotWazuh, *pivotFrom, *pivotTo, *executePivot, output)
		}
		if *snapshot == "" || *intoCase == "" || *importJSONL != "" ||
			*reportCase != "" || *exportCase != "" || *compareCase != "" ||
			*compareEvents || *sourceCase != "" || *sourceRef != "" ||
			*inspectCase != "" || *inspectEntry != "" || *serveImport == *triageImport {
			return errors.New("select Wazuh viewer or explicit two-Case temporal triage")
		}
		if *serveImport {
			if *withCase != "" || *compareProximity || windowSet {
				return errors.New("Wazuh viewer cannot combine with comparison options")
			}
			return serveDemoWithWazuhAlertsFiles(ctx, *principal, *snapshot, *intoCase, []string(wazuhFiles), output)
		}
		if len(wazuhFiles) == 1 {
			if *withCase == "" {
				if *compareProximity || windowSet {
					return errors.New("Wazuh temporal leads require a second authorized Case")
				}
				return writeDemoWazuh4Triage(ctx, *principal, *snapshot, *intoCase, wazuhFiles[0], output)
			}
			if *compareProximity {
				return writeDemoWazuh4Proximity(ctx, *principal, *snapshot, *intoCase, *withCase,
					wazuhFiles[0], *proximitySeconds, output)
			}
			if windowSet {
				return errors.New("Wazuh proximity window requires temporal lead selection")
			}
			return writeDemoWazuh4Comparison(ctx, *principal, *snapshot, *intoCase, *withCase,
				wazuhFiles[0], output)
		}
		if *withCase == "" {
			if *compareProximity || windowSet {
				return errors.New("multi-feed temporal leads require a second authorized Case")
			}
			return writeDemoMultiWazuhTriage(ctx, *principal, *snapshot, *intoCase,
				[]string(wazuhFiles), output)
		}
		if *compareProximity {
			return writeDemoMultiWazuhProximity(ctx, *principal, *snapshot, *intoCase, *withCase,
				[]string(wazuhFiles), *proximitySeconds, output)
		}
		if windowSet {
			return errors.New("Wazuh proximity window requires temporal lead selection")
		}
		return writeDemoMultiWazuhComparison(ctx, *principal, *snapshot, *intoCase, *withCase,
			[]string(wazuhFiles), output)
	}
	if *serveImport && (*snapshot == "" || *importJSONL == "" || *intoCase == "") {
		return errors.New("serve-import requires an owner snapshot, normalized JSONL file and destination Case")
	}
	if *importJSONL != "" || *intoCase != "" || *triageImport || *serveImport {
		if *importJSONL == "" || *intoCase == "" || *reportCase != "" || *exportCase != "" ||
			*compareCase != "" || (*withCase != "" && !*triageImport) || *compareEvents ||
			(*compareProximity && (!*triageImport || *withCase == "")) ||
			(*serveImport && (*triageImport || *withCase != "" || *compareProximity)) ||
			*sourceCase != "" || *sourceRef != "" || *inspectCase != "" || *inspectEntry != "" {
			return errors.New("local JSONL import requires a Case and cannot be combined with investigation modes")
		}
		if *serveImport {
			if len(extraJSONL) > 0 {
				all := append([]string{*importJSONL}, extraJSONL...)
				return serveDemoWithMultipleNormalizedJSONL(ctx, *principal, *snapshot, *intoCase, all, output)
			}
			return serveDemoWithImportedJSONL(ctx, *principal, *snapshot, *intoCase, *importJSONL, output)
		}
		if *triageImport {
			if *withCase != "" {
				if *compareProximity {
					return writeDemoImportedJSONLProximity(ctx, *principal, *snapshot, *intoCase, *withCase, *importJSONL, *proximitySeconds, output)
				}
				return writeDemoImportedJSONLComparison(ctx, *principal, *snapshot, *intoCase, *withCase, *importJSONL, output)
			}
			return writeDemoImportedJSONLTriage(ctx, *principal, *snapshot, *intoCase, *importJSONL, output)
		}
		return writeDemoImportedJSONL(ctx, *principal, *snapshot, *intoCase, *importJSONL, output)
	}
	if *exportCase != "" {
		if *reportCase != "" || *compareCase != "" || *withCase != "" ||
			*compareEvents || *compareProximity || *sourceCase != "" ||
			*sourceRef != "" || *inspectCase != "" || *inspectEntry != "" {
			return errors.New("JSONL export cannot be combined with another investigation mode")
		}
		return writeDemoCaseJSONL(ctx, *principal, *snapshot, *exportCase, output)
	}
	if *inspectCase != "" || *inspectEntry != "" {
		if *inspectCase == "" || *inspectEntry == "" || *reportCase != "" ||
			*compareCase != "" || *withCase != "" || *compareEvents || *compareProximity ||
			*sourceCase != "" || *sourceRef != "" {
			return errors.New("select one complete offline inspection mode")
		}
		return writeDemoEntryInspection(ctx, *principal, *snapshot, *inspectCase, *inspectEntry, output)
	}
	if *sourceCase != "" || *sourceRef != "" {
		if *sourceCase == "" || *sourceRef == "" || *reportCase != "" || *compareCase != "" || *withCase != "" || *compareEvents || *compareProximity {
			return errors.New("select one complete offline report mode")
		}
		return writeDemoSourceTimeline(ctx, *principal, *snapshot, *sourceCase, *sourceRef, output)
	}
	if *reportCase != "" && (*compareCase != "" || *withCase != "") {
		return errors.New("select one offline report mode")
	}
	if *compareCase != "" || *withCase != "" || *compareEvents || *compareProximity {
		if *compareProximity {
			return writeDemoCaseProximity(ctx, *principal, *snapshot, *compareCase, *withCase, *proximitySeconds, output)
		}
		if *compareEvents {
			return writeDemoCaseEventSequence(ctx, *principal, *snapshot, *compareCase, *withCase, output)
		}
		return writeDemoCaseComparison(ctx, *principal, *snapshot, *compareCase, *withCase, output)
	}
	if *reportCase != "" {
		return writeDemoCaseTriage(ctx, *principal, *snapshot, *reportCase, output)
	}
	return serveDemoWithSnapshot(ctx, *principal, *snapshot, output)
}

func randomBearer() (string, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return "", fmt.Errorf("generate demo credential: %w", err)
	}
	return hex.EncodeToString(secret[:]), nil
}

func newDemoHandler(token, principal string) (http.Handler, error) {
	return newDemoHandlerWithData(token, principal, demoData())
}

func newDemoHandlerWithData(token, principal string, data timeline.MemoryTimelineData) (http.Handler, error) {
	if len(token) != 64 {
		return nil, errors.New("demo bearer must have 32 random bytes in hex")
	}
	if _, err := hex.DecodeString(token); err != nil {
		return nil, errors.New("invalid demo bearer encoding")
	}
	if _, err := firstDemoAuthorizedCase(data, principal); err != nil {
		return nil, errors.New("principal must be a Case reader in the selected local owner dataset")
	}
	authenticate := func(r *http.Request) (string, error) {
		expected := "Bearer " + token
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte(expected)) != 1 {
			return "", errors.New("not authenticated")
		}
		return principal, nil
	}
	handlers, err := timeline.NewMemoryCaseTimelineHTTPHandlers(authenticate, data)
	if err != nil {
		return nil, fmt.Errorf("assemble synthetic owner chain: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/case-timeline", handlers.Timeline)
	mux.Handle("/case-source-navigation", handlers.Navigate)
	mux.Handle("/demo-source-open", handlers.OpenSource)
	mux.Handle("/demo-cases", handlers.Cases)
	installAnalystDemoRoutes(mux)
	return mux, nil
}

func serveDemo(ctx context.Context, principal string, output io.Writer) error {
	return serveDemoWithSnapshot(ctx, principal, "", output)
}

func serveDemoWithSnapshot(ctx context.Context, principal, snapshot string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || principal == "*" {
		return errors.New("select a synthetic principal from the local owner dataset")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	return serveDemoWithData(ctx, principal, data, output)
}

// serveDemoWithData reuses exactly the same authenticated loopback handler for
// a disk snapshot and a strictly in-memory imported snapshot.
func serveDemoWithData(ctx context.Context, principal string, data timeline.MemoryTimelineData, output io.Writer) error {
	return serveDemoWithDataRoute(ctx, principal, data, output, nil)
}

func serveDemoWithDataRoute(ctx context.Context, principal string, data timeline.MemoryTimelineData, output io.Writer, augment func(token string, base http.Handler) http.Handler) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || principal == "*" {
		return errors.New("select a synthetic principal from the local owner dataset")
	}
	caseID, err := firstDemoAuthorizedCase(data, principal)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	if augment != nil {
		handler = augment(token, handler)
	}
	// Never accept a listener address from HTTP, CLI or an environment flag.
	// The fixture is accessible only on the local IPv4 loopback interface.
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen on loopback: %w", err)
	}
	defer listener.Close()
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 << 10,
	}
	fmt.Fprintf(output, "CMDR synthetic Case Timeline demo; NOT a production service\n")
	fmt.Fprintf(output, "URL: http://%s/case-timeline?case_id=%s&limit=2\n", listener.Addr(), url.QueryEscape(caseID))
	fmt.Fprintf(output, "Local analyst viewer: http://%s/analyst\n", listener.Addr())
	fmt.Fprintf(output, "Bearer token (this run only): %s\n", token)
	// Quote operator-controlled identifiers to prevent multiline terminal output.
	fmt.Fprintf(output, "Principal: %q; available synthetic Case: %q\n", principal, caseID)
	fmt.Fprintln(output, "Only authenticated loopback requests are accepted. Ctrl-C exits.")
	ready := make(chan error, 1)
	go func() {
		ready <- server.Serve(listener)
	}()
	select {
	case err := <-ready:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("demo HTTP server: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("demo HTTP shutdown: %w", err)
		}
		if err := <-ready; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("demo HTTP server: %w", err)
		}
		return nil
	}
}
