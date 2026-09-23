package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	developmentToolRegistryPath = "engineering/security/development-tools.json"
	engineeringGoVersion        = "1.26.8"
	gosecModule                 = "github.com/securego/gosec/v2/cmd/gosec"
	gosecVersion                = "v2.28.0"
	govulncheckModule           = "golang.org/x/vuln/cmd/govulncheck"
	govulncheckVersion          = "v1.8.0"
)

type DevelopmentTool struct {
	ID            string `json:"id"`
	Module        string `json:"module"`
	Version       string `json:"version"`
	Commit        string `json:"commit"`
	License       string `json:"license"`
	Execution     string `json:"execution"`
	DataPolicy    string `json:"data_policy"`
	NetworkPolicy string `json:"network_policy"`
}

type DevelopmentToolRegistry struct {
	SchemaVersion int               `json:"schema_version"`
	Tools         []DevelopmentTool `json:"tools"`
}

type SASTFinding struct {
	RuleID      string `json:"rule_id"`
	Severity    string `json:"severity"`
	Confidence  string `json:"confidence"`
	CWE         string `json:"cwe,omitempty"`
	Path        string `json:"path"`
	Line        int    `json:"line"`
	Fingerprint string `json:"fingerprint_sha256"`
}

type SASTSummary struct {
	Tool         string         `json:"tool"`
	Version      string         `json:"version"`
	ScanRoot     string         `json:"scan_root"`
	FindingCount int            `json:"findings"`
	BySeverity   map[string]int `json:"by_severity"`
	Findings     []SASTFinding  `json:"finding_items,omitempty"`
}

type SCAFinding struct {
	OSV          string `json:"osv"`
	Ecosystem    string `json:"ecosystem"`
	Module       string `json:"module"`
	Version      string `json:"version"`
	Package      string `json:"package,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	FixedVersion string `json:"fixed_version,omitempty"`
}

type SCASummary struct {
	Tool                  string       `json:"tool"`
	Version               string       `json:"version"`
	Database              string       `json:"database,omitempty"`
	GoVersion             string       `json:"go_version,omitempty"`
	Modules               int          `json:"modules"`
	InformationalFindings int          `json:"informational_findings"`
	ActionableFindings    int          `json:"actionable_findings"`
	Findings              []SCAFinding `json:"finding_items,omitempty"`
}

type gosecJSONReport struct {
	Errors map[string][]json.RawMessage `json:"Golang errors"`
	Issues []struct {
		Severity   string `json:"severity"`
		Confidence string `json:"confidence"`
		CWE        *struct {
			ID string `json:"id"`
		} `json:"cwe"`
		RuleID string `json:"rule_id"`
		File   string `json:"file"`
		Line   string `json:"line"`
		Code   string `json:"code"`
	} `json:"Issues"`
}

type govulnJSONMessage struct {
	Config *struct {
		ProtocolVersion string `json:"protocol_version"`
		ScannerName     string `json:"scanner_name,omitempty"`
		ScannerVersion  string `json:"scanner_version,omitempty"`
		DB              string `json:"db,omitempty"`
		GoVersion       string `json:"go_version,omitempty"`
	} `json:"config,omitempty"`
	SBOM *struct {
		GoVersion string `json:"go_version,omitempty"`
		Modules   []struct {
			Path    string `json:"path,omitempty"`
			Version string `json:"version,omitempty"`
		} `json:"modules,omitempty"`
	} `json:"SBOM,omitempty"`
	Finding *struct {
		OSV          string `json:"osv,omitempty"`
		FixedVersion string `json:"fixed_version,omitempty"`
		Trace        []struct {
			Module   string `json:"module"`
			Version  string `json:"version,omitempty"`
			Package  string `json:"package,omitempty"`
			Function string `json:"function,omitempty"`
			Receiver string `json:"receiver,omitempty"`
		} `json:"trace,omitempty"`
	} `json:"finding,omitempty"`
}

func loadDevelopmentTools(root string) (DevelopmentToolRegistry, error) {
	var registry DevelopmentToolRegistry
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(developmentToolRegistryPath)), &registry); err != nil {
		return registry, err
	}
	if registry.SchemaVersion != 1 {
		return registry, fmt.Errorf("unsupported development-tool schema_version %d", registry.SchemaVersion)
	}
	expected := map[string]struct {
		module  string
		version string
	}{
		"gosec":       {module: gosecModule, version: gosecVersion},
		"govulncheck": {module: govulncheckModule, version: govulncheckVersion},
	}
	if len(registry.Tools) != len(expected) {
		return registry, fmt.Errorf("development-tool registry must contain exactly %d reviewed tools", len(expected))
	}
	seen := map[string]bool{}
	for _, tool := range registry.Tools {
		want, ok := expected[tool.ID]
		if !ok {
			return registry, fmt.Errorf("unreviewed development tool %q", tool.ID)
		}
		if seen[tool.ID] {
			return registry, fmt.Errorf("duplicate development tool %q", tool.ID)
		}
		seen[tool.ID] = true
		if tool.Module != want.module || tool.Version != want.version {
			return registry, fmt.Errorf("development tool %s provenance mismatch", tool.ID)
		}
		if strings.TrimSpace(tool.Commit) == "" || strings.TrimSpace(tool.License) == "" ||
			strings.TrimSpace(tool.Execution) == "" || strings.TrimSpace(tool.DataPolicy) == "" ||
			strings.TrimSpace(tool.NetworkPolicy) == "" {
			return registry, fmt.Errorf("development tool %s has incomplete provenance", tool.ID)
		}
	}
	return registry, nil
}

func runGoSAST(root string) (SASTSummary, error) {
	if _, err := loadDevelopmentTools(root); err != nil {
		return SASTSummary{}, err
	}
	moduleRoot := filepath.Join(root, "tools", "cmdr-dev")
	tmp, err := os.MkdirTemp("", "cmdr-gosec-")
	if err != nil {
		return SASTSummary{}, err
	}
	defer os.RemoveAll(tmp)
	reportPath := filepath.Join(tmp, "gosec.json")
	spec := gosecModule + "@" + gosecVersion
	cmd := exec.Command("go", "run", spec,
		"-no-fail",
		"-fmt=json",
		"-out="+reportPath,
		"-severity=medium",
		"-confidence=medium",
		"-exclude-generated",
		"-nosec-require-rules",
		"-nosec-require-justification",
		"./...",
	)
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOSECGOVERSION=go"+engineeringGoVersion)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return SASTSummary{}, fmt.Errorf("gosec %s execution failed; diagnostic output withheld from CI evidence: %w", gosecVersion, err)
	}
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return SASTSummary{}, fmt.Errorf("read gosec report: %w", err)
	}
	return parseGosecReport(root, moduleRoot, data)
}

func parseGosecReport(root, moduleRoot string, data []byte) (SASTSummary, error) {
	var report gosecJSONReport
	if err := json.Unmarshal(data, &report); err != nil {
		return SASTSummary{}, fmt.Errorf("decode gosec report: %w", err)
	}
	errorCount := 0
	for _, values := range report.Errors {
		errorCount += len(values)
	}
	if errorCount > 0 {
		return SASTSummary{}, fmt.Errorf("gosec reported %d processing error(s); details withheld from CI evidence", errorCount)
	}
	summary := SASTSummary{
		Tool:       "gosec",
		Version:    gosecVersion,
		ScanRoot:   "tools/cmdr-dev",
		BySeverity: map[string]int{},
	}
	for _, issue := range report.Issues {
		severity := strings.ToUpper(strings.TrimSpace(issue.Severity))
		confidence := strings.ToUpper(strings.TrimSpace(issue.Confidence))
		if severity != "HIGH" && severity != "MEDIUM" {
			return summary, fmt.Errorf("gosec returned unexpected filtered severity %q", issue.Severity)
		}
		if confidence != "HIGH" && confidence != "MEDIUM" {
			return summary, fmt.Errorf("gosec returned unexpected filtered confidence %q", issue.Confidence)
		}
		if strings.TrimSpace(issue.RuleID) == "" {
			return summary, fmt.Errorf("gosec finding is missing rule_id")
		}
		path, err := normalizeSecurityToolPath(root, moduleRoot, issue.File)
		if err != nil {
			return summary, err
		}
		line, err := parseSecurityLine(issue.Line)
		if err != nil {
			return summary, fmt.Errorf("gosec %s line: %w", issue.RuleID, err)
		}
		cwe := ""
		if issue.CWE != nil {
			cwe = issue.CWE.ID
		}
		summary.Findings = append(summary.Findings, SASTFinding{
			RuleID: issue.RuleID, Severity: severity, Confidence: confidence,
			CWE: cwe, Path: path, Line: line,
			Fingerprint: sastFindingFingerprint(issue.RuleID, path, issue.Code),
		})
		summary.BySeverity[severity]++
	}
	sort.Slice(summary.Findings, func(i, j int) bool {
		if summary.Findings[i].Path != summary.Findings[j].Path {
			return summary.Findings[i].Path < summary.Findings[j].Path
		}
		if summary.Findings[i].Line != summary.Findings[j].Line {
			return summary.Findings[i].Line < summary.Findings[j].Line
		}
		return summary.Findings[i].RuleID < summary.Findings[j].RuleID
	})
	summary.FindingCount = len(summary.Findings)
	if summary.FindingCount > 0 {
		return summary, fmt.Errorf("gosec blocked: %d medium/high confidence medium/high severity finding(s); code snippets are withheld", summary.FindingCount)
	}
	return summary, nil
}

func runGoSCA(root string) (SCASummary, error) {
	if _, err := loadDevelopmentTools(root); err != nil {
		return SCASummary{}, err
	}
	moduleRoot := filepath.Join(root, "tools", "cmdr-dev")
	spec := govulncheckModule + "@" + govulncheckVersion
	cmd := exec.Command("go", "run", spec, "-json", "./...")
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return SCASummary{}, fmt.Errorf("govulncheck %s execution failed; diagnostic output withheld from CI evidence: %w", govulncheckVersion, err)
	}
	return parseGovulncheckStream(stdout.Bytes())
}

func parseGovulncheckStream(data []byte) (SCASummary, error) {
	summary := SCASummary{Tool: "govulncheck", Version: govulncheckVersion}
	moduleVersions := map[string]string{}
	type rawFinding struct {
		osv, fixed                               string
		module, version, pkg, function, receiver string
		actionable                               bool
	}
	var raw []rawFinding
	dec := json.NewDecoder(bytes.NewReader(data))
	for {
		var message govulnJSONMessage
		if err := dec.Decode(&message); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return summary, fmt.Errorf("decode govulncheck stream: %w", err)
		}
		if message.Config != nil {
			if message.Config.ProtocolVersion == "" {
				return summary, fmt.Errorf("govulncheck config is missing protocol version")
			}
			if message.Config.ScannerName != "" && message.Config.ScannerName != "govulncheck" {
				return summary, fmt.Errorf("unexpected vulnerability scanner %q", message.Config.ScannerName)
			}
			summary.Database = message.Config.DB
			if message.Config.GoVersion != "" {
				summary.GoVersion = message.Config.GoVersion
			}
		}
		if message.SBOM != nil {
			if message.SBOM.GoVersion != "" {
				summary.GoVersion = message.SBOM.GoVersion
			}
			for _, module := range message.SBOM.Modules {
				if module.Path == "" {
					continue
				}
				moduleVersions[module.Path] = module.Version
			}
		}
		if message.Finding != nil {
			f := rawFinding{osv: message.Finding.OSV, fixed: message.Finding.FixedVersion}
			if len(message.Finding.Trace) > 0 {
				frame := message.Finding.Trace[0]
				f.module, f.version, f.pkg = frame.Module, frame.Version, frame.Package
				f.function, f.receiver = frame.Function, frame.Receiver
				f.actionable = frame.Function != ""
			}
			raw = append(raw, f)
		}
	}
	summary.Modules = len(moduleVersions)
	seen := map[string]bool{}
	for _, finding := range raw {
		if !finding.actionable {
			summary.InformationalFindings++
			continue
		}
		if finding.osv == "" || finding.module == "" {
			return summary, fmt.Errorf("actionable govulncheck finding is missing OSV or module identity")
		}
		version := finding.version
		if version == "" {
			version = moduleVersions[finding.module]
		}
		if version == "" && finding.module == "stdlib" {
			version = summary.GoVersion
		}
		if version == "" {
			return summary, fmt.Errorf("actionable govulncheck finding %s for %s has no exact version", finding.osv, finding.module)
		}
		symbol := finding.function
		if finding.receiver != "" {
			symbol = finding.receiver + "." + symbol
		}
		key := strings.Join([]string{finding.osv, finding.module, version, finding.pkg}, "\x00")
		if seen[key] {
			continue
		}
		seen[key] = true
		summary.Findings = append(summary.Findings, SCAFinding{
			OSV: finding.osv, Ecosystem: "Go", Module: finding.module, Version: version,
			Package: finding.pkg, Symbol: symbol, FixedVersion: finding.fixed,
		})
	}
	sort.Slice(summary.Findings, func(i, j int) bool {
		if summary.Findings[i].OSV != summary.Findings[j].OSV {
			return summary.Findings[i].OSV < summary.Findings[j].OSV
		}
		if summary.Findings[i].Module != summary.Findings[j].Module {
			return summary.Findings[i].Module < summary.Findings[j].Module
		}
		return summary.Findings[i].Package < summary.Findings[j].Package
	})
	summary.ActionableFindings = len(summary.Findings)
	if summary.ActionableFindings > 0 {
		return summary, fmt.Errorf("govulncheck blocked: %d reachable known vulnerability finding(s)", summary.ActionableFindings)
	}
	return summary, nil
}

func sastFindingFingerprint(ruleID, path, code string) string {
	sum := sha256.Sum256([]byte(ruleID + "\x00" + path + "\x00" + code))
	return fmt.Sprintf("%x", sum[:])
}

func formatSASTSafeFindings(summary SASTSummary) string {
	if len(summary.Findings) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, finding := range summary.Findings {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, "%s/%s/%s:%d/%s/%s",
			finding.RuleID, finding.Severity, finding.Path, finding.Line, finding.CWE, finding.Fingerprint)
	}
	b.WriteByte(']')
	return b.String()
}

func normalizeSecurityToolPath(root, moduleRoot, reported string) (string, error) {
	reported = strings.TrimSpace(reported)
	if reported == "" {
		return "", fmt.Errorf("security tool reported an empty file path")
	}
	path := reported
	if !filepath.IsAbs(path) {
		path = filepath.Join(moduleRoot, filepath.FromSlash(path))
	}
	rel, err := filepath.Rel(root, filepath.Clean(path))
	if err != nil {
		return "", err
	}
	values, err := normalizeChangedPaths([]string{filepath.ToSlash(rel)})
	if err != nil || len(values) != 1 {
		return "", fmt.Errorf("security tool path escapes repository: %q", reported)
	}
	return values[0], nil
}

func parseSecurityLine(value string) (int, error) {
	value = strings.TrimSpace(value)
	if i := strings.IndexByte(value, '-'); i >= 0 {
		value = value[:i]
	}
	line, err := strconv.Atoi(value)
	if err != nil || line < 1 {
		return 0, fmt.Errorf("invalid line %q", value)
	}
	return line, nil
}
