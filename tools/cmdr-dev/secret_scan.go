package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	secretAllowlistPath = "engineering/security/secret-allowlist.json"
	maxSecretScanBytes  = 16 * 1024 * 1024
)

type secretRule struct {
	ID      string
	Pattern *regexp.Regexp
}

type SecretFinding struct {
	RuleID      string `json:"rule_id"`
	Path        string `json:"path"`
	Line        int    `json:"line"`
	Fingerprint string `json:"fingerprint_sha256"`
}

type SecretAllowlistEntry struct {
	RuleID      string `json:"rule_id"`
	Path        string `json:"path"`
	Fingerprint string `json:"fingerprint_sha256"`
	Reason      string `json:"reason"`
}

type SecretAllowlist struct {
	SchemaVersion int                    `json:"schema_version"`
	Entries       []SecretAllowlistEntry `json:"entries"`
}

type SecretScanSummary struct {
	Mode             string          `json:"mode"`
	CandidatePaths   int             `json:"candidate_paths"`
	ScannedFiles     int             `json:"scanned_files"`
	SkippedDeleted   int             `json:"skipped_deleted"`
	SkippedBinary    int             `json:"skipped_binary"`
	SkippedGenerated int             `json:"skipped_generated"`
	SkippedSymlink   int             `json:"skipped_symlink"`
	Allowlisted      int             `json:"allowlisted"`
	FindingCount     int             `json:"findings"`
	ByRule           map[string]int  `json:"by_rule"`
	FindingItems     []SecretFinding `json:"finding_items,omitempty"`
}

var secretFingerprintPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

var secretRules = []secretRule{
	{ID: "SECRET-PRIVATE-KEY", Pattern: regexp.MustCompile(`(?s)-----BEGIN (?:RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----.*?-----END (?:RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`)},
	{ID: "SECRET-GITHUB-TOKEN", Pattern: regexp.MustCompile(`(?:gh[pousr]_[A-Za-z0-9]{36,255}|github_pat_[A-Za-z0-9_]{20,255})`)},
	{ID: "SECRET-AWS-ACCESS-KEY", Pattern: regexp.MustCompile(`(?:AKIA|ASIA|A3T[A-Z0-9])[A-Z0-9]{16}`)},
	{ID: "SECRET-SLACK-TOKEN", Pattern: regexp.MustCompile(`xox[baprs]-[A-Za-z0-9-]{20,255}`)},
	{ID: "SECRET-STRIPE-LIVE-KEY", Pattern: regexp.MustCompile(`sk_live_[A-Za-z0-9]{16,255}`)},
	{ID: "SECRET-GOOGLE-API-KEY", Pattern: regexp.MustCompile(`AIza[0-9A-Za-z_-]{35}`)},
}

func runSecretScan(root, changesFile string, fullScan bool) (SecretScanSummary, error) {
	allowlist, err := loadSecretAllowlist(root)
	if err != nil {
		return SecretScanSummary{}, err
	}
	allowed, err := indexSecretAllowlist(allowlist)
	if err != nil {
		return SecretScanSummary{}, err
	}

	var candidates []string
	mode := "changed"
	switch {
	case fullScan && strings.TrimSpace(changesFile) != "":
		return SecretScanSummary{}, fmt.Errorf("secret scan accepts either --full-scan or --changes-file, not both")
	case fullScan:
		mode = "full"
		candidates, err = trackedRepositoryPaths(root)
	case strings.TrimSpace(changesFile) != "":
		candidates, err = readChangedPaths(root, changesFile)
	default:
		return SecretScanSummary{}, fmt.Errorf("secret scan requires --changes-file or --full-scan")
	}
	if err != nil {
		return SecretScanSummary{}, err
	}

	summary := SecretScanSummary{Mode: mode, CandidatePaths: len(candidates), ByRule: map[string]int{}}
	for _, rel := range candidates {
		if isSecretGeneratedPath(rel) {
			summary.SkippedGenerated++
			continue
		}
		info, err := lstatRepoEntry(root, rel)
		if os.IsNotExist(err) {
			summary.SkippedDeleted++
			continue
		}
		if err != nil {
			return summary, fmt.Errorf("inspect secret-scan path %s: %w", rel, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			summary.SkippedSymlink++
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		if info.Size() > maxSecretScanBytes {
			return summary, fmt.Errorf("secret scan refuses oversized text candidate %s (%d bytes > %d)", rel, info.Size(), maxSecretScanBytes)
		}
		data, err := readRepoFile(root, rel)
		if err != nil {
			return summary, fmt.Errorf("read secret-scan path %s: %w", rel, err)
		}
		if isSecretBinary(data) {
			summary.SkippedBinary++
			continue
		}
		if isGeneratedSource(data) {
			summary.SkippedGenerated++
			continue
		}
		summary.ScannedFiles++
		for _, finding := range scanSecretContent(rel, data) {
			if allowed[secretAllowlistKey(finding.RuleID, finding.Path, finding.Fingerprint)] {
				summary.Allowlisted++
				continue
			}
			summary.ByRule[finding.RuleID]++
			summary.FindingItems = append(summary.FindingItems, finding)
		}
	}
	sort.Slice(summary.FindingItems, func(i, j int) bool {
		if summary.FindingItems[i].Path != summary.FindingItems[j].Path {
			return summary.FindingItems[i].Path < summary.FindingItems[j].Path
		}
		if summary.FindingItems[i].Line != summary.FindingItems[j].Line {
			return summary.FindingItems[i].Line < summary.FindingItems[j].Line
		}
		if summary.FindingItems[i].RuleID != summary.FindingItems[j].RuleID {
			return summary.FindingItems[i].RuleID < summary.FindingItems[j].RuleID
		}
		return summary.FindingItems[i].Fingerprint < summary.FindingItems[j].Fingerprint
	})
	summary.FindingCount = len(summary.FindingItems)
	if summary.FindingCount > 0 {
		return summary, fmt.Errorf("secret scan blocked: %d high-confidence finding(s); matched values are redacted", summary.FindingCount)
	}
	return summary, nil
}

func scanSecretContent(rel string, data []byte) []SecretFinding {
	var findings []SecretFinding
	for _, rule := range secretRules {
		for _, loc := range rule.Pattern.FindAllIndex(data, -1) {
			match := data[loc[0]:loc[1]]
			findings = append(findings, SecretFinding{
				RuleID: rule.ID, Path: rel, Line: 1 + bytes.Count(data[:loc[0]], []byte{'\n'}),
				Fingerprint: secretFingerprint(rule.ID, match),
			})
		}
	}
	return findings
}

func secretFingerprint(ruleID string, secret []byte) string {
	h := sha256.New()
	_, _ = h.Write([]byte(ruleID))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write(secret)
	return hex.EncodeToString(h.Sum(nil))
}

func loadSecretAllowlist(root string) (SecretAllowlist, error) {
	var allowlist SecretAllowlist
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(secretAllowlistPath)), &allowlist); err != nil {
		return allowlist, err
	}
	if allowlist.SchemaVersion != 1 {
		return allowlist, fmt.Errorf("unsupported secret allowlist schema_version %d", allowlist.SchemaVersion)
	}
	return allowlist, nil
}

func indexSecretAllowlist(allowlist SecretAllowlist) (map[string]bool, error) {
	index := map[string]bool{}
	for i, entry := range allowlist.Entries {
		if !secretRuleKnown(entry.RuleID) {
			return nil, fmt.Errorf("secret allowlist entry %d uses unknown rule %q", i, entry.RuleID)
		}
		normalized, err := normalizeChangedPaths([]string{entry.Path})
		if err != nil || len(normalized) != 1 || normalized[0] != entry.Path {
			return nil, fmt.Errorf("secret allowlist entry %d has non-canonical path %q", i, entry.Path)
		}
		if !secretFingerprintPattern.MatchString(entry.Fingerprint) {
			return nil, fmt.Errorf("secret allowlist entry %d has invalid SHA-256 fingerprint", i)
		}
		if strings.TrimSpace(entry.Reason) == "" {
			return nil, fmt.Errorf("secret allowlist entry %d requires a reason", i)
		}
		key := secretAllowlistKey(entry.RuleID, entry.Path, entry.Fingerprint)
		if index[key] {
			return nil, fmt.Errorf("duplicate secret allowlist entry for %s at %s", entry.RuleID, entry.Path)
		}
		index[key] = true
	}
	return index, nil
}

func secretAllowlistKey(ruleID, path, fingerprint string) string {
	return ruleID + "\x00" + path + "\x00" + fingerprint
}

func secretRuleKnown(id string) bool {
	for _, rule := range secretRules {
		if rule.ID == id {
			return true
		}
	}
	return false
}

func trackedRepositoryPaths(root string) ([]string, error) {
	cmd := exec.Command("git", "-C", root, "ls-files", "-z", "--cached") // #nosec G204,G702 -- executable and arguments are fixed; root is the resolved repository root and no shell is used.
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("list tracked files for secret scan: %w", err)
	}
	return parseNULPaths(out)
}

func isSecretBinary(data []byte) bool {
	sample := data
	if len(sample) > 8192 {
		sample = sample[:8192]
	}
	return bytes.IndexByte(sample, 0) >= 0
}

func isGeneratedSource(data []byte) bool {
	sample := data
	if len(sample) > 4096 {
		sample = sample[:4096]
	}
	lower := strings.ToLower(string(sample))
	return strings.Contains(lower, "code generated") && strings.Contains(lower, "do not edit")
}

func isSecretGeneratedPath(path string) bool {
	path = filepath.ToSlash(path)
	for _, prefix := range []string{"node_modules/", "vendor/", "dist/", "build/", "out/", "coverage/", "target/", ".next/", "bin/", "obj/"} {
		if strings.HasPrefix(path, prefix) || strings.Contains(path, "/"+prefix) {
			return true
		}
	}
	lower := strings.ToLower(path)
	for _, suffix := range []string{
		".min.js", ".min.css", ".map", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf",
		".zip", ".gz", ".tgz", ".7z", ".rar", ".jar", ".war", ".class", ".wasm", ".exe", ".dll",
		".so", ".dylib", ".woff", ".woff2", ".ttf", ".otf", ".mp3", ".mp4", ".mov", ".avi",
		".parquet", ".sqlite", ".db",
	} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}
