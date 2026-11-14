package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type HealthGateStatus string

const (
	StatusMainHealthy    HealthGateStatus = "MAIN_HEALTHY"
	StatusMainUnhealthy  HealthGateStatus = "MAIN_UNHEALTHY"
	StatusMainUnverified HealthGateStatus = "MAIN_UNVERIFIED"

	ExitCodeMainHealthy    = 0
	ExitCodeMainUnhealthy  = 1
	ExitCodeMainUnverified = 2
)

type HealthDimensionResult struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Status  string   `json:"status"` // "PASS", "FAIL", "UNVERIFIED"
	Details []string `json:"details,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type MainlineHealthGateReport struct {
	Status     HealthGateStatus                 `json:"status"`
	ExitCode   int                              `json:"exit_code"`
	Summary    string                           `json:"summary"`
	Timestamp  string                           `json:"timestamp"`
	Root       string                           `json:"root"`
	HeadCommit string                           `json:"head_commit"`
	Dimensions map[string]HealthDimensionResult `json:"dimensions"`
}

func evaluateCodeDimension(root string, state CurrentState, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "code", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	doc, err := doctor(root, state, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("doctor validation failed: %v", err))
		return res
	}
	if !doc.Healthy {
		res.Errors = append(res.Errors, doc.Warnings...)
		return res
	}

	cat, err := runCheckCatalogAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("check catalog audit failed: %v", err))
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("doctor checks passed (%d checks verified)", len(doc.Checks)),
		fmt.Sprintf("check catalog verified (%d checks registered)", cat.Checks),
	)
	return res
}

func evaluateSecurityDimension(root string) HealthDimensionResult {
	res := HealthDimensionResult{Name: "security", Passed: false, Status: "FAIL"}

	secScan, err := runSecretScan(root, "", true)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("secret scan failed: %v", err))
		return res
	}
	if secScan.FindingCount > 0 {
		res.Errors = append(res.Errors, fmt.Sprintf("%d secret findings detected across %d paths", secScan.FindingCount, secScan.CandidatePaths))
		return res
	}

	gates, err := runSecurityGateAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("security gates audit failed: %v", err))
		return res
	}

	if engineeringGoVersion != "1.26.11" {
		res.Errors = append(res.Errors, fmt.Sprintf("engineering Go toolchain version mismatch: expected 1.26.11, got %s", engineeringGoVersion))
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("0 secrets detected across %d candidate paths", secScan.CandidatePaths),
		fmt.Sprintf("security gate registry valid (%d registered gates, %d blocking)", gates.Gates, gates.Blocking),
		fmt.Sprintf("Go toolchain pinned to %s", engineeringGoVersion),
	)
	return res
}

func evaluateGitDimension(root string) (HealthDimensionResult, string) {
	res := HealthDimensionResult{Name: "git", Passed: false, Status: "FAIL"}

	cmdHead := exec.Command("git", "-C", root, "rev-parse", "HEAD") // #nosec G204,G702 -- git rev-parse HEAD within confined root
	headBytes, err := cmdHead.Output()
	if err != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, fmt.Sprintf("unable to resolve HEAD commit: %v", err))
		return res, ""
	}
	headCommit := strings.TrimSpace(string(headBytes))
	if !fullCommitIDPattern.MatchString(headCommit) {
		res.Errors = append(res.Errors, fmt.Sprintf("invalid HEAD commit SHA: %q", headCommit))
		return res, headCommit
	}

	cmdStatus := exec.Command("git", "-C", root, "status", "--porcelain") // #nosec G204,G702 -- git status within confined root
	statusBytes, err := cmdStatus.Output()
	if err != nil {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, fmt.Sprintf("git status check failed: %v", err))
		return res, headCommit
	}
	statusStr := strings.TrimSpace(string(statusBytes))
	if len(statusStr) > 0 {
		res.Errors = append(res.Errors, fmt.Sprintf("working tree has uncommitted modifications: %s", statusStr))
		return res, headCommit
	}

	const baselineCommit = "82e630fce3a7f6236e03709116eca6dd0a1acd40"
	if gitCommitExists(root, baselineCommit) || ensureGitCommit(root, baselineCommit) == nil {
		cmdDiff := exec.Command("git", "-C", root, "diff", baselineCommit, "--", "cmdr-product-spec") // #nosec G204,G702 -- invariant diff check
		diffBytes, err := cmdDiff.Output()
		if err != nil {
			res.Status = "UNVERIFIED"
			res.Errors = append(res.Errors, fmt.Sprintf("git diff product spec baseline failed: %v", err))
			return res, headCommit
		}
		if len(bytes.TrimSpace(diffBytes)) > 0 {
			res.Errors = append(res.Errors, "cmdr-product-spec has non-zero diff vs baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40")
			return res, headCommit
		}
	} else {
		res.Status = "UNVERIFIED"
		res.Errors = append(res.Errors, "product spec baseline commit 82e630fce3a7f6236e03709116eca6dd0a1acd40 not available in git history")
		return res, headCommit
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("HEAD commit resolved: %s", headCommit),
		"working tree clean (0 uncommitted changes)",
		"cmdr-product-spec 0-byte invariant verified against baseline 82e630fce3a7f6236e03709116eca6dd0a1acd40",
	)
	return res, headCommit
}

func evaluateTraceabilityDimension(root string, state CurrentState, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "traceability", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	ledger, err := runImplementationLedgerAudit(root, state, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("implementation ledger audit failed: %v", err))
		return res
	}
	if ledger.Status != "valid" && ledger.Status != "PASS" {
		res.Errors = append(res.Errors, fmt.Sprintf("implementation ledger status is %s", ledger.Status))
		return res
	}

	dec, err := runDecisionRegistryAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("decision registry audit failed: %v", err))
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("implementation ledger %d entries verified against spec baseline %s", ledger.Entries, ledger.ProductSpecBaseline),
		fmt.Sprintf("decision registry %d decisions reconciled with work graph", dec.Decisions),
	)
	return res
}

func evaluateOperationalQualityDimension(root string, graph WorkGraph, stateErr error) HealthDimensionResult {
	res := HealthDimensionResult{Name: "operational_quality", Passed: false, Status: "FAIL"}
	if stateErr != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("repository state failed to load: %v", stateErr))
		return res
	}

	comp, err := runComplexityAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("complexity audit failed: %v", err))
		return res
	}
	if comp.ReadinessViolations > 0 {
		res.Errors = append(res.Errors, fmt.Sprintf("%d complexity readiness violations detected", comp.ReadinessViolations))
		return res
	}

	perf, err := runPerformanceRegistryAudit(root)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("performance registry audit failed: %v", err))
		return res
	}
	if perf.RuntimeCoverageStatus != "complete" {
		res.Errors = append(res.Errors, fmt.Sprintf("performance registry runtime coverage status is %s", perf.RuntimeCoverageStatus))
		return res
	}

	rec, err := runRecoveryJournalAudit(root, graph)
	if err != nil {
		res.Errors = append(res.Errors, fmt.Sprintf("recovery journal audit failed: %v", err))
		return res
	}
	if rec.Status != "valid" {
		res.Errors = append(res.Errors, fmt.Sprintf("recovery journal status is %s", rec.Status))
		return res
	}

	res.Passed = true
	res.Status = "PASS"
	res.Details = append(res.Details,
		fmt.Sprintf("complexity within limits across %d evaluated units (%d within budget)", comp.Evaluated, comp.WithinBudget),
		fmt.Sprintf("performance registry complete with %d targets and %d metrics", perf.Targets, perf.Metrics),
		fmt.Sprintf("recovery journal valid with %d events intact", rec.Events),
	)
	return res
}

func runMainlineHealthGate(root string, state CurrentState, graph WorkGraph, stateErr error) (MainlineHealthGateReport, int) {
	report := MainlineHealthGateReport{
		Root:       root,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Dimensions: make(map[string]HealthDimensionResult, 5),
	}

	codeDim := evaluateCodeDimension(root, state, graph, stateErr)
	report.Dimensions["code"] = codeDim

	secDim := evaluateSecurityDimension(root)
	report.Dimensions["security"] = secDim

	gitDim, headCommit := evaluateGitDimension(root)
	report.Dimensions["git"] = gitDim
	report.HeadCommit = headCommit

	traceDim := evaluateTraceabilityDimension(root, state, graph, stateErr)
	report.Dimensions["traceability"] = traceDim

	opDim := evaluateOperationalQualityDimension(root, graph, stateErr)
	report.Dimensions["operational_quality"] = opDim

	hasFail := false
	hasUnverified := false

	for _, dim := range report.Dimensions {
		if dim.Status == "FAIL" || !dim.Passed && dim.Status != "UNVERIFIED" {
			hasFail = true
		}
		if dim.Status == "UNVERIFIED" {
			hasUnverified = true
		}
	}

	if hasFail {
		report.Status = StatusMainUnhealthy
		report.ExitCode = ExitCodeMainUnhealthy
		report.Summary = "Mainline health gate FAILED; one or more required health dimensions failed verification."
	} else if hasUnverified {
		report.Status = StatusMainUnverified
		report.ExitCode = ExitCodeMainUnverified
		report.Summary = "Mainline health gate UNVERIFIED due to unresolved environmental or diagnostic conditions."
	} else {
		report.Status = StatusMainHealthy
		report.ExitCode = ExitCodeMainHealthy
		report.Summary = "Mainline health gate PASSED across all 5 dimensions (Code, Security, Git, Traceability, Operational Quality)."
	}

	return report, report.ExitCode
}

func printMainlineHealthGateReport(report MainlineHealthGateReport) {
	fmt.Printf("=== MAINLINE HEALTH GATE REPORT ===\n")
	fmt.Printf("Status:      %s (exit %d)\n", report.Status, report.ExitCode)
	fmt.Printf("Summary:     %s\n", report.Summary)
	fmt.Printf("Timestamp:   %s\n", report.Timestamp)
	fmt.Printf("Root:        %s\n", report.Root)
	fmt.Printf("Head Commit: %s\n", report.HeadCommit)
	fmt.Println("-----------------------------------")
	dimOrder := []string{"code", "security", "git", "traceability", "operational_quality"}
	for _, name := range dimOrder {
		dim, ok := report.Dimensions[name]
		if !ok {
			continue
		}
		statusTag := fmt.Sprintf("[%s]", dim.Status)
		fmt.Printf("%-12s %-20s (passed=%t)\n", statusTag, dim.Name, dim.Passed)
		for _, d := range dim.Details {
			fmt.Printf("    ✓ %s\n", d)
		}
		for _, e := range dim.Errors {
			fmt.Printf("    ✗ %s\n", e)
		}
	}
	fmt.Println("===================================")
}
