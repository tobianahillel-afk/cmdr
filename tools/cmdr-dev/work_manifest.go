package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ProductRefsV2 struct {
	Capabilities            []string `json:"capabilities"`
	Requirements            []string `json:"requirements"`
	Screens                 []string `json:"screens"`
	Permissions             []string `json:"permissions"`
	OpenDecisions           []string `json:"open_decisions"`
	ImplementationContracts []string `json:"implementation_contracts"`
	CanonicalObjects        []string `json:"canonical_objects"`
}

type PerformanceBudgetV2 struct {
	Applicable bool     `json:"applicable"`
	Rationale  string   `json:"rationale"`
	Metrics    []string `json:"metrics"`
}

type MigrationPlanV2 struct {
	Required  bool     `json:"required"`
	Rationale string   `json:"rationale"`
	Steps     []string `json:"steps"`
}

type RollbackPlanV2 struct {
	Strategy     string   `json:"strategy"`
	Verification []string `json:"verification"`
}

type ComplexityEstimateV2 struct {
	PrimaryComponents           int `json:"primary_components"`
	EstimatedFilesChanged       int `json:"estimated_files_changed"`
	EstimatedNetLOC             int `json:"estimated_net_loc"`
	ProductionFiles             int `json:"production_files"`
	NetProductionLOC            int `json:"net_production_loc"`
	BoundedContexts             int `json:"bounded_contexts"`
	IndependentMigrations       int `json:"independent_migrations"`
	DistinctSecurityModels      int `json:"distinct_security_models"`
	SeparatelyTestableBehaviors int `json:"separately_testable_behaviors"`
}

type WorkManifestV2 struct {
	SchemaVersion      int                  `json:"schema_version"`
	ID                 string               `json:"id"`
	Title              string               `json:"title"`
	Type               string               `json:"type"`
	Parent             string               `json:"parent"`
	Goal               string               `json:"goal"`
	ProductRefs        ProductRefsV2        `json:"product_refs"`
	DependsOn          []string             `json:"depends_on"`
	AllowedPaths       []string             `json:"allowed_paths"`
	ForbiddenPaths     []string             `json:"forbidden_paths"`
	Inputs             []string             `json:"inputs"`
	Outputs            []string             `json:"outputs"`
	Invariants         []string             `json:"invariants"`
	SecurityProperties []string             `json:"security_properties"`
	AcceptanceTests    []string             `json:"acceptance_tests"`
	PerformanceBudget  PerformanceBudgetV2  `json:"performance_budget"`
	Migration          MigrationPlanV2      `json:"migration"`
	Rollback           RollbackPlanV2       `json:"rollback"`
	RequiredChecks     []string             `json:"required_checks"`
	DefinitionOfReady  []string             `json:"definition_of_ready"`
	DefinitionOfDone   []string             `json:"definition_of_done"`
	Complexity         ComplexityEstimateV2 `json:"complexity"`
}

type ManifestValidationSummary struct {
	Manifests int    `json:"manifests"`
	LegacyV1  int    `json:"legacy_v1"`
	StrictV2  int    `json:"strict_v2"`
	Mode      string `json:"mode"`
}

type manifestHeader struct {
	SchemaVersion int    `json:"schema_version"`
	ID            string `json:"id"`
}

var workUnitIDPattern = regexp.MustCompile(`^E[0-9]+-[A-Z0-9]+-[0-9]{3}[A-Z]?(?:-[A-Z0-9]+)?$`)

func runManifestValidation(root string) (ManifestValidationSummary, error) {
	paths, err := manifestPaths(root)
	if err != nil {
		return ManifestValidationSummary{}, err
	}
	headers := map[string]string{}
	versions := map[string]int{}
	for _, path := range paths {
		header, err := decodeManifestHeader(path)
		if err != nil {
			return ManifestValidationSummary{}, err
		}
		if header.ID == "" {
			return ManifestValidationSummary{}, fmt.Errorf("%s: manifest id is empty", path)
		}
		if previous, exists := headers[header.ID]; exists {
			return ManifestValidationSummary{}, fmt.Errorf("duplicate manifest id %s in %s and %s", header.ID, previous, path)
		}
		headers[header.ID] = path
		versions[header.ID] = header.SchemaVersion
	}

	summary := ManifestValidationSummary{Manifests: len(paths), Mode: "check"}
	for id, path := range headers {
		switch versions[id] {
		case 1:
			if !strings.HasPrefix(id, "E0-") {
				return ManifestValidationSummary{}, fmt.Errorf("%s: schema-v1 manifests are legacy E0 only", path)
			}
			summary.LegacyV1++
		case 2:
			manifest, err := decodeWorkManifestV2(path)
			if err != nil {
				return ManifestValidationSummary{}, err
			}
			if err := validateWorkManifestV2(root, path, manifest, headers); err != nil {
				return ManifestValidationSummary{}, err
			}
			summary.StrictV2++
		default:
			return ManifestValidationSummary{}, fmt.Errorf("%s: unsupported manifest schema_version %d", path, versions[id])
		}
	}
	return summary, nil
}

func manifestPaths(root string) ([]string, error) {
	base := filepath.Join(root, "work", "lots")
	var paths []string
	err := filepath.WalkDir(base, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && entry.Name() == "manifest.json" {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk work manifests: %w", err)
	}
	sort.Strings(paths)
	return paths, nil
}

func decodeManifestHeader(path string) (manifestHeader, error) {
	var header manifestHeader
	data, err := os.ReadFile(path)
	if err != nil {
		return header, err
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return header, fmt.Errorf("%s: decode manifest header: %w", path, err)
	}
	return header, nil
}

func decodeWorkManifestV2(path string) (WorkManifestV2, error) {
	var manifest WorkManifestV2
	f, err := os.Open(path)
	if err != nil {
		return manifest, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&manifest); err != nil {
		return manifest, fmt.Errorf("%s: decode manifest v2: %w", path, err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return manifest, fmt.Errorf("%s: trailing JSON value", path)
		}
		return manifest, fmt.Errorf("%s: trailing content: %w", path, err)
	}
	return manifest, nil
}

func validateWorkManifestV2(root, path string, m WorkManifestV2, known map[string]string) error {
	fail := func(format string, args ...any) error {
		return fmt.Errorf("%s: %s", path, fmt.Sprintf(format, args...))
	}
	if m.SchemaVersion != 2 {
		return fail("schema_version must be 2")
	}
	if !workUnitIDPattern.MatchString(m.ID) {
		return fail("invalid work-unit id %q", m.ID)
	}
	if filepath.Base(filepath.Dir(path)) != m.ID {
		return fail("manifest id %s does not match directory %s", m.ID, filepath.Base(filepath.Dir(path)))
	}
	if strings.TrimSpace(m.Title) == "" || strings.TrimSpace(m.Goal) == "" {
		return fail("title and goal are required")
	}
	switch m.Type {
	case "epic":
		if m.Parent != "" {
			return fail("epic parent must be empty")
		}
	case "lot", "sublot", "task":
		if m.Parent == "" {
			return fail("%s parent is required", m.Type)
		}
	default:
		return fail("invalid type %q", m.Type)
	}
	if m.Parent != "" {
		if !workUnitIDPattern.MatchString(m.Parent) {
			return fail("invalid parent id %q", m.Parent)
		}
		if _, ok := known[m.Parent]; !ok {
			return fail("parent %s has no manifest", m.Parent)
		}
	}
	for _, dep := range m.DependsOn {
		if !workUnitIDPattern.MatchString(dep) {
			return fail("invalid dependency id %q", dep)
		}
		if dep == m.ID {
			return fail("manifest cannot depend on itself")
		}
		if _, ok := known[dep]; !ok {
			return fail("dependency %s has no manifest", dep)
		}
	}
	if err := validateProductRefs(root, m.ProductRefs); err != nil {
		return fail("%v", err)
	}
	if len(m.AllowedPaths) == 0 || len(m.ForbiddenPaths) == 0 {
		return fail("allowed_paths and forbidden_paths are required")
	}
	for _, allowed := range m.AllowedPaths {
		for _, forbidden := range m.ForbiddenPaths {
			if pathPatternsOverlap(allowed, forbidden) {
				return fail("allowed path %q conflicts with forbidden path %q", allowed, forbidden)
			}
		}
	}
	requiredLists := map[string][]string{
		"inputs": m.Inputs, "outputs": m.Outputs, "invariants": m.Invariants,
		"security_properties": m.SecurityProperties, "acceptance_tests": m.AcceptanceTests,
		"required_checks": m.RequiredChecks, "definition_of_ready": m.DefinitionOfReady,
		"definition_of_done": m.DefinitionOfDone,
	}
	for name, values := range requiredLists {
		if len(nonEmptyStrings(values)) == 0 {
			return fail("%s must contain at least one non-empty item", name)
		}
	}
	if strings.TrimSpace(m.PerformanceBudget.Rationale) == "" {
		return fail("performance_budget.rationale is required")
	}
	if m.PerformanceBudget.Applicable && len(nonEmptyStrings(m.PerformanceBudget.Metrics)) == 0 {
		return fail("applicable performance budget requires metrics")
	}
	if strings.TrimSpace(m.Migration.Rationale) == "" {
		return fail("migration.rationale is required")
	}
	if m.Migration.Required && len(nonEmptyStrings(m.Migration.Steps)) == 0 {
		return fail("required migration requires steps")
	}
	if strings.TrimSpace(m.Rollback.Strategy) == "" || len(nonEmptyStrings(m.Rollback.Verification)) == 0 {
		return fail("rollback strategy and verification are required")
	}
	values := []int{
		m.Complexity.PrimaryComponents, m.Complexity.EstimatedFilesChanged, m.Complexity.EstimatedNetLOC,
		m.Complexity.ProductionFiles, m.Complexity.NetProductionLOC, m.Complexity.BoundedContexts,
		m.Complexity.IndependentMigrations, m.Complexity.DistinctSecurityModels, m.Complexity.SeparatelyTestableBehaviors,
	}
	for _, value := range values {
		if value < 0 {
			return fail("complexity values cannot be negative")
		}
	}
	if m.Complexity.PrimaryComponents == 0 || m.Complexity.BoundedContexts == 0 || m.Complexity.SeparatelyTestableBehaviors == 0 {
		return fail("complexity primary_components, bounded_contexts and separately_testable_behaviors must be positive")
	}
	return nil
}

func validateProductRefs(root string, refs ProductRefsV2) error {
	checkIDs := func(name string, values []string, matcher *regexp.Regexp) error {
		for _, value := range values {
			if !matcher.MatchString(value) {
				return fmt.Errorf("product_refs.%s contains invalid id %q", name, value)
			}
		}
		return nil
	}
	if err := checkIDs("capabilities", refs.Capabilities, capabilityIDPattern); err != nil {
		return err
	}
	if err := checkIDs("requirements", refs.Requirements, requirementIDPattern); err != nil {
		return err
	}
	if err := checkIDs("screens", refs.Screens, screenIDPattern); err != nil {
		return err
	}
	for _, value := range refs.Permissions {
		if !strings.HasPrefix(value, "perm.") {
			return fmt.Errorf("product_refs.permissions contains invalid id %q", value)
		}
	}
	if err := checkIDs("open_decisions", refs.OpenDecisions, openDecisionIDPattern); err != nil {
		return err
	}
	for _, value := range append(append([]string(nil), refs.ImplementationContracts...), refs.CanonicalObjects...) {
		if !strings.HasPrefix(value, "cmdr-product-spec/") {
			return fmt.Errorf("product source path %q must be under cmdr-product-spec/", value)
		}
		info, err := statRepoPath(root, value)
		if err != nil || info.IsDir() {
			return fmt.Errorf("product source path %q does not exist as a file", value)
		}
	}
	return nil
}

func pathPatternsOverlap(a, b string) bool {
	a = filepath.ToSlash(strings.TrimSpace(a))
	b = filepath.ToSlash(strings.TrimSpace(b))
	if a == b {
		return true
	}
	prefix := func(pattern string) string { return strings.TrimSuffix(pattern, "**") }
	ap, bp := prefix(a), prefix(b)
	return strings.HasPrefix(ap, bp) || strings.HasPrefix(bp, ap)
}

func nonEmptyStrings(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	return out
}

func decodeWorkManifestV2Bytes(data []byte) (WorkManifestV2, error) {
	var manifest WorkManifestV2
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&manifest); err != nil {
		return manifest, err
	}
	return manifest, nil
}
