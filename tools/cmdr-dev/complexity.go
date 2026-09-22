package main

import (
	"fmt"
	"sort"
)

const (
	softPrimaryComponents = 2
	softFiles             = 8
	softLOC               = 400
	hardFiles             = 12
	hardLOC               = 800
)

type ComplexityResult struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	GraphStatus        string   `json:"graph_status"`
	Decision           string   `json:"decision"`
	Warnings           []string `json:"warnings,omitempty"`
	SplitTriggers      []string `json:"split_triggers,omitempty"`
	ReadinessViolation bool     `json:"readiness_violation"`
}

type ComplexityAuditSummary struct {
	Evaluated           int                `json:"evaluated"`
	WithinBudget        int                `json:"within_budget"`
	WithWarnings        int                `json:"with_warnings"`
	SplitRequired       int                `json:"split_required"`
	ReadinessViolations int                `json:"readiness_violations"`
	Results             []ComplexityResult `json:"results"`
}

func runComplexityAudit(root string, graph WorkGraph) (ComplexityAuditSummary, error) {
	paths, err := manifestPaths(root)
	if err != nil {
		return ComplexityAuditSummary{}, err
	}
	graphIndex := make(map[string]WorkNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		graphIndex[node.ID] = node
	}

	var summary ComplexityAuditSummary
	for _, path := range paths {
		header, err := decodeManifestHeader(path)
		if err != nil {
			return ComplexityAuditSummary{}, err
		}
		if header.SchemaVersion != 2 {
			continue
		}
		manifest, err := decodeWorkManifestV2(path)
		if err != nil {
			return ComplexityAuditSummary{}, err
		}
		if manifest.Type == "epic" {
			continue
		}
		node, ok := graphIndex[manifest.ID]
		if !ok {
			return ComplexityAuditSummary{}, fmt.Errorf("%s: strict-v2 manifest has no work-graph node", manifest.ID)
		}
		result := evaluateComplexity(manifest, node.Status)
		summary.Results = append(summary.Results, result)
		summary.Evaluated++
		switch result.Decision {
		case "within-budget":
			summary.WithinBudget++
		case "warning":
			summary.WithWarnings++
		case "split-required":
			summary.SplitRequired++
		}
		if result.ReadinessViolation {
			summary.ReadinessViolations++
		}
	}
	sort.Slice(summary.Results, func(i, j int) bool {
		return summary.Results[i].ID < summary.Results[j].ID
	})
	return summary, nil
}

func evaluateComplexity(manifest WorkManifestV2, graphStatus string) ComplexityResult {
	c := manifest.Complexity
	result := ComplexityResult{
		ID:          manifest.ID,
		Type:        manifest.Type,
		GraphStatus: graphStatus,
		Decision:    "within-budget",
	}

	files := maxInt(c.EstimatedFilesChanged, c.ProductionFiles)
	loc := maxInt(c.EstimatedNetLOC, c.NetProductionLOC)

	if c.PrimaryComponents > softPrimaryComponents {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("primary components %d exceed target %d", c.PrimaryComponents, softPrimaryComponents))
	}
	if files > softFiles {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("material files %d exceed target %d", files, softFiles))
	}
	if loc > softLOC {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("net LOC %d exceed target %d", loc, softLOC))
	}

	if files > hardFiles {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("material files %d exceed hard limit %d", files, hardFiles))
	}
	if loc > hardLOC {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("net LOC %d exceed hard limit %d", loc, hardLOC))
	}
	if c.BoundedContexts > 1 {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("bounded contexts %d exceed hard limit 1", c.BoundedContexts))
	}
	if c.IndependentMigrations > 1 {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("independent migrations %d exceed hard limit 1", c.IndependentMigrations))
	}
	if c.DistinctSecurityModels > 1 {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("distinct security models %d exceed hard limit 1", c.DistinctSecurityModels))
	}
	if c.SeparatelyTestableBehaviors > 1 {
		result.SplitTriggers = append(result.SplitTriggers,
			fmt.Sprintf("separately testable behaviors %d exceed hard limit 1", c.SeparatelyTestableBehaviors))
	}

	sort.Strings(result.Warnings)
	sort.Strings(result.SplitTriggers)
	if len(result.SplitTriggers) > 0 {
		result.Decision = "split-required"
		if rank, ok := lifecycleRank[graphStatus]; ok && rank >= lifecycleRank["READY"] {
			result.ReadinessViolation = true
		}
	} else if len(result.Warnings) > 0 {
		result.Decision = "warning"
	}
	return result
}

func validateComplexityReadiness(summary ComplexityAuditSummary) error {
	if summary.ReadinessViolations == 0 {
		return nil
	}
	var ids []string
	for _, result := range summary.Results {
		if result.ReadinessViolation {
			ids = append(ids, result.ID)
		}
	}
	sort.Strings(ids)
	return fmt.Errorf("work units require mandatory split before READY: %v", ids)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
