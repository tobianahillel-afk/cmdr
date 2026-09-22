package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type ProductEntity struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	SourcePath    string `json:"source_path,omitempty"`
	SourceSHA256  string `json:"source_sha256,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Product       string `json:"product,omitempty"`
	Module        string `json:"module,omitempty"`
	Status        string `json:"status,omitempty"`
	SourceOfTruth string `json:"source_of_truth,omitempty"`
	ReferenceOnly bool   `json:"reference_only"`
}

type ProductEdge struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Kind       string `json:"kind"`
	SourcePath string `json:"source_path"`
}

type ProductGraph struct {
	SchemaVersion  int             `json:"schema_version"`
	SpecTreeDigest string          `json:"spec_tree_digest"`
	Entities       []ProductEntity `json:"entities"`
	Edges          []ProductEdge   `json:"edges"`
}

type CoverageGraphSummary struct {
	Entities                      int    `json:"entities"`
	Edges                         int    `json:"edges"`
	OwnedCapabilities             int    `json:"owned_capabilities"`
	ReferenceOnlyCapabilities     int    `json:"reference_only_capabilities"`
	DistinctRequirementReferences int    `json:"distinct_requirement_references"`
	DistinctOpenDecisionReferences int   `json:"distinct_open_decision_references"`
	DistinctPermissionReferences  int    `json:"distinct_permission_references"`
	OwnedScreens                  int    `json:"owned_screens"`
	ReferenceOnlyScreens          int    `json:"reference_only_screens"`
	SpecTreeDigest                string `json:"spec_tree_digest"`
	Output                        string `json:"output"`
	Mode                          string `json:"mode"`
}

var capabilityIDPattern = regexp.MustCompile(`^CAP-[A-Z0-9]+-[0-9]{3}$`)
var requirementIDPattern = regexp.MustCompile(`^REQ-[A-Z0-9]+-[0-9]{3}$`)
var openDecisionIDPattern = regexp.MustCompile(`^OPEN-[0-9]{3}$`)
var adrIDPattern = regexp.MustCompile(`^ADR-[0-9]{4}$`)
var dependencyIDPattern = regexp.MustCompile(`^DEP-(?:[A-Z0-9]+-)?[0-9]{3}$`)

func runCoverageGraph(root, specRel, output string, check bool) (CoverageGraphSummary, error) {
	inventory, err := buildSpecInventory(root, specRel)
	if err != nil {
		return CoverageGraphSummary{}, err
	}
	graph, err := buildProductGraph(inventory)
	if err != nil {
		return CoverageGraphSummary{}, err
	}
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return CoverageGraphSummary{}, err
	}
	data = append(data, '\n')

	outputPath := output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(root, filepath.FromSlash(outputPath))
	}
	mode := "write"
	if check {
		mode = "check"
		existing, err := os.ReadFile(outputPath)
		if err != nil {
			return CoverageGraphSummary{}, fmt.Errorf("read product graph: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return CoverageGraphSummary{}, fmt.Errorf("product graph is stale: %s", outputPath)
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
			return CoverageGraphSummary{}, err
		}
		if err := os.WriteFile(outputPath, data, 0o644); err != nil {
			return CoverageGraphSummary{}, err
		}
	}

	summary := summarizeProductGraph(graph)
	summary.Output = outputPath
	summary.Mode = mode
	return summary, nil
}

func buildProductGraph(inventory SpecInventory) (ProductGraph, error) {
	entities := map[string]ProductEntity{}
	var edges []ProductEdge

	for _, doc := range inventory.Documents {
		if !doc.Active {
			continue
		}
		sourceID := documentEntityID(doc)
		entity := ProductEntity{
			ID:            sourceID,
			Kind:          classifyOwnedDocument(doc),
			SourcePath:    doc.Path,
			SourceSHA256:  doc.SHA256,
			Domain:        doc.Domain,
			Product:       doc.Product,
			Module:        doc.Module,
			Status:        doc.Status,
			SourceOfTruth: doc.SourceOfTruth,
			ReferenceOnly: false,
		}
		if err := addEntity(entities, entity); err != nil {
			return ProductGraph{}, err
		}

		explicit := map[string]string{}
		for _, id := range doc.Requirements {
			explicit[id] = "requirement"
		}
		for _, id := range doc.OpenDecisions {
			if id != "none" {
				explicit[id] = "open-decision"
			}
		}
		for _, id := range doc.Permissions {
			explicit[id] = "permission"
		}

		for target, kind := range explicit {
			targetKind := classifyReference(target)
			addReferenceEntity(entities, target, targetKind)
			edges = append(edges, ProductEdge{
				From:       sourceID,
				To:         target,
				Kind:       kind,
				SourcePath: doc.Path,
			})
		}

		for _, target := range doc.References {
			if target == sourceID {
				continue
			}
			if _, alreadyExplicit := explicit[target]; alreadyExplicit {
				continue
			}
			targetKind := classifyReference(target)
			addReferenceEntity(entities, target, targetKind)
			edges = append(edges, ProductEdge{
				From:       sourceID,
				To:         target,
				Kind:       "reference",
				SourcePath: doc.Path,
			})
		}
	}

	entityList := make([]ProductEntity, 0, len(entities))
	for _, entity := range entities {
		entityList = append(entityList, entity)
	}
	sort.Slice(entityList, func(i, j int) bool {
		if entityList[i].ID == entityList[j].ID {
			return entityList[i].Kind < entityList[j].Kind
		}
		return entityList[i].ID < entityList[j].ID
	})
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		if edges[i].To != edges[j].To {
			return edges[i].To < edges[j].To
		}
		if edges[i].Kind != edges[j].Kind {
			return edges[i].Kind < edges[j].Kind
		}
		return edges[i].SourcePath < edges[j].SourcePath
	})
	edges = dedupeEdges(edges)

	if err := validateProductGraph(entityList, edges); err != nil {
		return ProductGraph{}, err
	}
	return ProductGraph{
		SchemaVersion:  1,
		SpecTreeDigest: inventory.TreeDigest,
		Entities:       entityList,
		Edges:          edges,
	}, nil
}

func documentEntityID(doc SpecDocument) string {
	if doc.ID != "" {
		return doc.ID
	}
	return "doc:" + doc.Path
}

func classifyOwnedDocument(doc SpecDocument) string {
	if capabilityIDPattern.MatchString(doc.ID) {
		return "capability"
	}
	if strings.EqualFold(doc.Type, "screen") || strings.EqualFold(doc.SourceOfTruth, "screen") || strings.Contains(doc.Path, "/screens/") {
		return "screen"
	}
	if doc.Domain == "13-user-journeys" || strings.Contains(doc.Path, "/13-user-journeys/") {
		return "journey"
	}
	if strings.Contains(doc.Path, "/17-implementation-contracts/") {
		return "implementation-contract"
	}
	if strings.Contains(doc.Path, "/05-domain-model/objects/") {
		return "canonical-object"
	}
	if requirementIDPattern.MatchString(doc.ID) {
		return "requirement"
	}
	if openDecisionIDPattern.MatchString(doc.ID) {
		return "open-decision"
	}
	if strings.HasPrefix(doc.ID, "perm.") {
		return "permission"
	}
	return "document"
}

func classifyReference(id string) string {
	switch {
	case capabilityIDPattern.MatchString(id):
		return "capability"
	case requirementIDPattern.MatchString(id):
		return "requirement"
	case openDecisionIDPattern.MatchString(id):
		return "open-decision"
	case strings.HasPrefix(id, "perm."):
		return "permission"
	case adrIDPattern.MatchString(id):
		return "architecture-decision"
	case dependencyIDPattern.MatchString(id):
		return "dependency"
	default:
		return "reference"
	}
}

func addEntity(entities map[string]ProductEntity, entity ProductEntity) error {
	if existing, ok := entities[entity.ID]; ok {
		if existing.ReferenceOnly && !entity.ReferenceOnly {
			entities[entity.ID] = entity
			return nil
		}
		if !existing.ReferenceOnly && !entity.ReferenceOnly && existing.SourcePath != entity.SourcePath {
			return fmt.Errorf("product entity %q has multiple active source documents: %s and %s", entity.ID, existing.SourcePath, entity.SourcePath)
		}
		return nil
	}
	entities[entity.ID] = entity
	return nil
}

func addReferenceEntity(entities map[string]ProductEntity, id, kind string) {
	if _, ok := entities[id]; ok {
		return
	}
	entities[id] = ProductEntity{
		ID:            id,
		Kind:          kind,
		ReferenceOnly: true,
	}
}

func validateProductGraph(entities []ProductEntity, edges []ProductEdge) error {
	known := make(map[string]ProductEntity, len(entities))
	for _, entity := range entities {
		if entity.ID == "" {
			return fmt.Errorf("product graph contains empty entity id")
		}
		if _, exists := known[entity.ID]; exists {
			return fmt.Errorf("product graph contains duplicate entity id %q", entity.ID)
		}
		if !entity.ReferenceOnly && (entity.SourcePath == "" || entity.SourceSHA256 == "") {
			return fmt.Errorf("owned product entity %q lacks source provenance", entity.ID)
		}
		known[entity.ID] = entity
	}
	for _, edge := range edges {
		if _, ok := known[edge.From]; !ok {
			return fmt.Errorf("edge source %q is missing", edge.From)
		}
		if _, ok := known[edge.To]; !ok {
			return fmt.Errorf("edge target %q is missing", edge.To)
		}
		if edge.SourcePath == "" {
			return fmt.Errorf("edge %s -> %s lacks source path", edge.From, edge.To)
		}
	}
	return nil
}

func dedupeEdges(edges []ProductEdge) []ProductEdge {
	if len(edges) == 0 {
		return nil
	}
	out := edges[:0]
	var previous ProductEdge
	for i, edge := range edges {
		if i == 0 || edge != previous {
			out = append(out, edge)
			previous = edge
		}
	}
	return out
}

func summarizeProductGraph(graph ProductGraph) CoverageGraphSummary {
	summary := CoverageGraphSummary{
		Entities:       len(graph.Entities),
		Edges:          len(graph.Edges),
		SpecTreeDigest: graph.SpecTreeDigest,
	}
	for _, entity := range graph.Entities {
		switch entity.Kind {
		case "capability":
			if entity.ReferenceOnly {
				summary.ReferenceOnlyCapabilities++
			} else {
				summary.OwnedCapabilities++
			}
		case "requirement":
			summary.DistinctRequirementReferences++
		case "open-decision":
			summary.DistinctOpenDecisionReferences++
		case "permission":
			summary.DistinctPermissionReferences++
		case "screen":
			if entity.ReferenceOnly {
				summary.ReferenceOnlyScreens++
			} else {
				summary.OwnedScreens++
			}
		}
	}
	return summary
}
