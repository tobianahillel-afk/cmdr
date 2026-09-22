package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	statePath = "engineering/state/current.json"
	graphPath = "work/graph.json"
)

type CurrentState struct {
	SchemaVersion int    `json:"schema_version"`
	StateKind     string `json:"state_kind"`
	ProductSpec   struct {
		CanonicalPath                string `json:"canonical_path"`
		BaselineCommit               string `json:"baseline_commit"`
		DocumentaryStatus            string `json:"documentary_status"`
		RuntimeImplementationClaimed bool   `json:"runtime_implementation_claimed"`
		KnownCapabilities            int    `json:"known_capabilities"`
		KnownRequirements            int    `json:"known_requirements"`
		KnownOpenDecisions           int    `json:"known_open_decisions"`
		KnownActiveScreens           int    `json:"known_active_screens"`
	} `json:"product_spec"`
	EngineeringFoundation struct {
		Status     string `json:"status"`
		Phase      string `json:"phase"`
		PhaseTitle string `json:"phase_title"`
	} `json:"engineering_foundation"`
	Execution struct {
		ActiveWorkUnit    string `json:"active_work_unit"`
		PreferredNextUnit string `json:"preferred_next_unit"`
		SelectionRule     string `json:"selection_rule"`
	} `json:"execution"`
	Reconciliation struct {
		RequiredOnResume             bool `json:"required_on_resume"`
		RecordedStateIsAuthoritative bool `json:"recorded_state_is_authoritative"`
	} `json:"reconciliation"`
}

type WorkGraph struct {
	SchemaVersion int        `json:"schema_version"`
	GraphKind     string     `json:"graph_kind"`
	Nodes         []WorkNode `json:"nodes"`
}

type WorkNode struct {
	ID                      string   `json:"id"`
	Title                   string   `json:"title"`
	Type                    string   `json:"type"`
	Status                  string   `json:"status"`
	DependsOn               []string `json:"depends_on"`
	DependencyTerminalState string   `json:"dependency_terminal_state,omitempty"`
	Unlocks                 []string `json:"unlocks"`
}

type StatusOutput struct {
	Root             string `json:"root"`
	Phase            string `json:"phase"`
	PhaseTitle       string `json:"phase_title"`
	FoundationStatus string `json:"foundation_status"`
	ActiveWorkUnit   string `json:"active_work_unit"`
	PreferredNext    string `json:"preferred_next_unit"`
	ComputedNext     string `json:"computed_next_unit,omitempty"`
	ProductBaseline  string `json:"product_spec_baseline"`
}

type DoctorOutput struct {
	Root     string   `json:"root"`
	Healthy  bool     `json:"healthy"`
	Checks   []string `json:"checks"`
	Warnings []string `json:"warnings,omitempty"`
}

var lifecycleRank = map[string]int{
	"DRAFT":             0,
	"DECOMPOSED":        1,
	"READY":             2,
	"CLAIMED":           3,
	"IMPLEMENTING":      4,
	"IMPLEMENTED":       5,
	"TESTING":           6,
	"REVIEWING":         7,
	"SECURITY_REVIEW":   8,
	"MERGE_READY":       9,
	"MERGED":            10,
	"POST_MERGE_VERIFY": 11,
	"VERIFIED":          12,
	"BLOCKED":           -1,
	"BLOCKED_DECISION":  -1,
}

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	command := os.Args[1]
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	rootFlag := fs.String("root", "", "repository root")
	jsonFlag := fs.Bool("json", false, "machine-readable output")
	checkFlag := fs.Bool("check", false, "check generated output instead of writing it")
	outputFlag := fs.String("output", "engineering/spec-index/inventory.json", "generated output path")
	if err := fs.Parse(os.Args[2:]); err != nil {
		fail(err)
	}

	root, err := resolveRoot(*rootFlag)
	if err != nil {
		fail(err)
	}

	state, graph, err := loadRepositoryState(root)
	if err != nil {
		fail(err)
	}

	switch command {
	case "doctor":
		out, err := doctor(root, state, graph)
		if err != nil {
			fail(err)
		}
		printValue(out, *jsonFlag)
	case "status":
		next, _ := selectNext(graph)
		out := StatusOutput{
			Root:             root,
			Phase:            state.EngineeringFoundation.Phase,
			PhaseTitle:       state.EngineeringFoundation.PhaseTitle,
			FoundationStatus: state.EngineeringFoundation.Status,
			ActiveWorkUnit:   state.Execution.ActiveWorkUnit,
			PreferredNext:    state.Execution.PreferredNextUnit,
			ProductBaseline:  state.ProductSpec.BaselineCommit,
		}
		if next != nil {
			out.ComputedNext = next.ID
		}
		printValue(out, *jsonFlag)
	case "next":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		next, err := selectNext(graph)
		if err != nil {
			fail(err)
		}
		if next == nil {
			fail(errors.New("no executable READY work unit"))
		}
		printValue(next, *jsonFlag)
	case "spec-index":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		summary, err := runSpecIndex(root, state.ProductSpec.CanonicalPath, *outputFlag, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "spec-baseline":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/spec-index/baseline.json"
		}
		summary, err := runSpecBaseline(root, state.ProductSpec.CanonicalPath, state.ProductSpec.BaselineCommit, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	case "coverage-graph":
		if err := validateState(root, state, graph); err != nil {
			fail(err)
		}
		output := *outputFlag
		if output == "engineering/spec-index/inventory.json" {
			output = "engineering/coverage/product-graph.json"
		}
		summary, err := runCoverageGraph(root, state.ProductSpec.CanonicalPath, output, *checkFlag)
		if err != nil {
			fail(err)
		}
		printValue(summary, *jsonFlag)
	default:
		usage(os.Stderr)
		fail(fmt.Errorf("unknown command %q", command))
	}
}

func resolveRoot(explicit string) (string, error) {
	if explicit != "" {
		abs, err := filepath.Abs(explicit)
		if err != nil {
			return "", err
		}
		if isRepoRoot(abs) {
			return abs, nil
		}
		return "", fmt.Errorf("%s is not a CMDR repository root", abs)
	}

	if env := os.Getenv("CMDR_ROOT"); env != "" {
		return resolveRoot(env)
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if isRepoRoot(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("CMDR repository root not found")
		}
		dir = parent
	}
}

func isRepoRoot(root string) bool {
	required := []string{"AGENTS.md", "AI_START_HERE.md", statePath, graphPath}
	for _, rel := range required {
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

func loadRepositoryState(root string) (CurrentState, WorkGraph, error) {
	var state CurrentState
	var graph WorkGraph
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(statePath)), &state); err != nil {
		return state, graph, err
	}
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(graphPath)), &graph); err != nil {
		return state, graph, err
	}
	return state, graph, nil
}

func decodeStrict(path string, dst any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode %s: trailing JSON value", path)
		}
		return fmt.Errorf("decode %s: trailing content: %w", path, err)
	}
	return nil
}

func doctor(root string, state CurrentState, graph WorkGraph) (DoctorOutput, error) {
	out := DoctorOutput{Root: root, Healthy: false}
	if err := validateState(root, state, graph); err != nil {
		return out, err
	}
	out.Checks = append(out.Checks,
		"canonical bootstrap files present",
		"state JSON decodes strictly",
		"work graph JSON decodes strictly",
		"work graph node ids and dependencies are valid",
		"work graph is acyclic",
		"active/preferred work units exist",
		"canonical product-spec path exists",
	)

	next, err := selectNext(graph)
	if err != nil {
		return out, err
	}
	if next == nil {
		out.Warnings = append(out.Warnings, "no executable READY work unit")
	} else if state.Execution.PreferredNextUnit != "" && state.Execution.PreferredNextUnit != next.ID {
		out.Warnings = append(out.Warnings,
			fmt.Sprintf("preferred next unit %s differs from computed next %s", state.Execution.PreferredNextUnit, next.ID))
	}
	out.Healthy = true
	return out, nil
}

func validateState(root string, state CurrentState, graph WorkGraph) error {
	if state.SchemaVersion != 1 || graph.SchemaVersion != 1 {
		return errors.New("unsupported state or graph schema version")
	}
	if state.StateKind == "" || graph.GraphKind == "" {
		return errors.New("state_kind and graph_kind are required")
	}
	if state.ProductSpec.CanonicalPath == "" {
		return errors.New("canonical product-spec path is empty")
	}
	if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(state.ProductSpec.CanonicalPath))); err != nil || !info.IsDir() {
		return fmt.Errorf("canonical product-spec path %q does not exist", state.ProductSpec.CanonicalPath)
	}

	index := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		if strings.TrimSpace(node.ID) == "" {
			return errors.New("work graph contains empty node id")
		}
		if _, exists := index[node.ID]; exists {
			return fmt.Errorf("duplicate work node id %s", node.ID)
		}
		if _, known := lifecycleRank[node.Status]; !known {
			return fmt.Errorf("work node %s has unknown status %s", node.ID, node.Status)
		}
		index[node.ID] = node
	}

	for _, node := range graph.Nodes {
		for _, dep := range node.DependsOn {
			if _, ok := index[dep]; !ok {
				return fmt.Errorf("work node %s depends on missing node %s", node.ID, dep)
			}
		}
		for _, unlocked := range node.Unlocks {
			if _, ok := index[unlocked]; !ok {
				return fmt.Errorf("work node %s unlocks missing node %s", node.ID, unlocked)
			}
		}
	}
	if err := validateAcyclic(graph.Nodes, index); err != nil {
		return err
	}
	if state.Execution.ActiveWorkUnit != "" {
		if _, ok := index[state.Execution.ActiveWorkUnit]; !ok {
			return fmt.Errorf("active work unit %s is missing from graph", state.Execution.ActiveWorkUnit)
		}
	}
	if state.Execution.PreferredNextUnit != "" {
		if _, ok := index[state.Execution.PreferredNextUnit]; !ok {
			return fmt.Errorf("preferred next unit %s is missing from graph", state.Execution.PreferredNextUnit)
		}
	}
	return nil
}

func validateAcyclic(nodes []WorkNode, index map[string]WorkNode) error {
	const (
		gray  = 1
		black = 2
	)
	state := make(map[string]int, len(nodes))
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case gray:
			return fmt.Errorf("work graph cycle detected at %s", id)
		case black:
			return nil
		}
		state[id] = gray
		for _, dep := range index[id].DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = black
		return nil
	}
	for _, node := range nodes {
		if err := visit(node.ID); err != nil {
			return err
		}
	}
	return nil
}

func selectNext(graph WorkGraph) (*WorkNode, error) {
	index := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		index[node.ID] = node
	}
	candidates := make([]WorkNode, 0)
	for _, node := range graph.Nodes {
		if node.Status != "READY" {
			continue
		}
		ok, err := dependenciesSatisfied(node, index)
		if err != nil {
			return nil, err
		}
		if ok {
			candidates = append(candidates, node)
		}
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })
	return &candidates[0], nil
}

func dependenciesSatisfied(node WorkNode, index map[string]WorkNode) (bool, error) {
	required := node.DependencyTerminalState
	if required == "" {
		required = "VERIFIED"
	}
	requiredRank, ok := lifecycleRank[required]
	if !ok || requiredRank < 0 {
		return false, fmt.Errorf("node %s has invalid dependency terminal state %s", node.ID, required)
	}
	for _, depID := range node.DependsOn {
		dep, ok := index[depID]
		if !ok {
			return false, fmt.Errorf("node %s depends on missing node %s", node.ID, depID)
		}
		rank, known := lifecycleRank[dep.Status]
		if !known || rank < requiredRank {
			return false, nil
		}
	}
	return true, nil
}

func printValue(v any, asJSON bool) {
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(v); err != nil {
			fail(err)
		}
		return
	}

	switch x := v.(type) {
	case DoctorOutput:
		fmt.Printf("CMDR engineering doctor: healthy=%t\n", x.Healthy)
		fmt.Printf("root: %s\n", x.Root)
		for _, check := range x.Checks {
			fmt.Printf("PASS  %s\n", check)
		}
		for _, warning := range x.Warnings {
			fmt.Printf("WARN  %s\n", warning)
		}
	case StatusOutput:
		fmt.Printf("phase: %s — %s\n", x.Phase, x.PhaseTitle)
		fmt.Printf("foundation: %s\n", x.FoundationStatus)
		fmt.Printf("active: %s\n", x.ActiveWorkUnit)
		fmt.Printf("preferred next: %s\n", x.PreferredNext)
		fmt.Printf("computed next: %s\n", x.ComputedNext)
		fmt.Printf("product baseline: %s\n", x.ProductBaseline)
	case *WorkNode:
		fmt.Printf("%s — %s [%s]\n", x.ID, x.Title, x.Status)
	case SpecIndexSummary:
		fmt.Printf("spec files: %d\n", x.Files)
		fmt.Printf("active canonical documents: %d\n", x.ActiveCanonicalDocuments)
		fmt.Printf("tree digest: %s\n", x.TreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case SpecBaselineSummary:
		fmt.Printf("spec files: %d\n", x.Files)
		fmt.Printf("active canonical documents: %d\n", x.ActiveCanonicalDocuments)
		fmt.Printf("tree digest: %s\n", x.TreeDigest)
		fmt.Printf("baseline commit: %s\n", x.ProductSpecBaselineCommit)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	case CoverageGraphSummary:
		fmt.Printf("entities: %d\n", x.Entities)
		fmt.Printf("edges: %d\n", x.Edges)
		fmt.Printf("owned capabilities: %d\n", x.OwnedCapabilities)
		fmt.Printf("reference-only capabilities: %d\n", x.ReferenceOnlyCapabilities)
		fmt.Printf("distinct requirement references: %d\n", x.DistinctRequirementReferences)
		fmt.Printf("distinct open-decision references: %d\n", x.DistinctOpenDecisionReferences)
		fmt.Printf("distinct permission references: %d\n", x.DistinctPermissionReferences)
		fmt.Printf("owned screens: %d\n", x.OwnedScreens)
		fmt.Printf("registered active screens: %d\n", x.RegisteredActiveScreens)
		fmt.Printf("unregistered owned screens: %v\n", x.UnregisteredOwnedScreens)
		fmt.Printf("reference-only screens: %d\n", x.ReferenceOnlyScreens)
		fmt.Printf("tree digest: %s\n", x.SpecTreeDigest)
		fmt.Printf("output: %s\n", x.Output)
		fmt.Printf("mode: %s\n", x.Mode)
	default:
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: cmdr-dev <doctor|status|next|spec-index|spec-baseline|coverage-graph> [--root PATH] [--json] [--check] [--output PATH]")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "cmdr-dev:", err)
	os.Exit(1)
}
