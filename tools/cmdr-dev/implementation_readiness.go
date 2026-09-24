package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const capabilityDependencyRegisterPath = "cmdr-product-spec/00-governance/dependency-register.md"

var readinessDependencyIDPattern = regexp.MustCompile("^DEP-[A-Z0-9-]+$")
var capabilitySelectorTokenPattern = regexp.MustCompile("CAP-[A-Z0-9]+-[0-9]{3}([.][.][0-9]{3}|(/[0-9]{3})*)?")
var openDecisionTokenPattern = regexp.MustCompile("OPEN-[0-9]{3}")

type CapabilityRegistryRecord struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	DeliveryStatus string   `json:"delivery_status"`
	DeliveryMode   string   `json:"delivery_mode"`
	CanonicalFile  string   `json:"canonical_file,omitempty"`
	OpenDecisions  []string `json:"open_decisions,omitempty"`
	SourcePath     string   `json:"source_path"`
}

type DependencyReadinessEvidence struct {
	ID                   string   `json:"id"`
	Dependent            string   `json:"dependent"`
	Status               string   `json:"status"`
	Blocking             string   `json:"blocking"`
	OpenDecisions        []string `json:"open_decisions,omitempty"`
	AffectedCapabilities []string `json:"affected_capabilities,omitempty"`
	SourcePath           string   `json:"source_path"`
}

type ImplementationReadinessBlocker struct {
	Kind          string   `json:"kind"`
	ID            string   `json:"id"`
	Status        string   `json:"status"`
	Blocking      string   `json:"blocking"`
	OpenDecisions []string `json:"open_decisions,omitempty"`
	SourcePaths   []string `json:"source_paths"`
}

type ImplementationReadinessRecord struct {
	Capability          string                           `json:"capability"`
	Name                string                           `json:"name"`
	State               string                           `json:"state"`
	DeliveryStatus      string                           `json:"delivery_status"`
	DeliveryMode        string                           `json:"delivery_mode"`
	RegisterPath        string                           `json:"register_path"`
	CanonicalFile       string                           `json:"canonical_file,omitempty"`
	ImplementedWorkUnit string                           `json:"implemented_work_unit,omitempty"`
	Blockers            []ImplementationReadinessBlocker `json:"blockers,omitempty"`
}

type ImplementationReadinessSummary struct {
	Total                        int `json:"total"`
	Implemented                  int `json:"implemented"`
	Ready                        int `json:"ready"`
	Blocked                      int `json:"blocked"`
	Proposed                     int `json:"proposed"`
	DecisionBlocked              int `json:"decision_blocked"`
	DependencyBlocked            int `json:"dependency_blocked"`
	UnscopedBlockingDependencies int `json:"unscoped_blocking_dependencies"`
}

type ImplementationReadinessProgram struct {
	SchemaVersion                int                             `json:"schema_version"`
	ProgramKind                  string                          `json:"program_kind"`
	ProductSpecBaseline          string                          `json:"product_spec_baseline"`
	SpecTreeDigest               string                          `json:"spec_tree_digest"`
	Records                      []ImplementationReadinessRecord `json:"records"`
	UnscopedBlockingDependencies []DependencyReadinessEvidence   `json:"unscoped_blocking_dependencies,omitempty"`
	Summary                      ImplementationReadinessSummary  `json:"summary"`
	Status                       string                          `json:"status"`
}

func runImplementationReadiness(root string, state CurrentState, graph WorkGraph) (ImplementationReadinessProgram, error) {
	ledgerAudit, err := runImplementationLedgerAudit(root, state, graph)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}
	var ledger ImplementationLedger
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(implementationLedgerPath)), &ledger); err != nil {
		return ImplementationReadinessProgram{}, err
	}

	records, err := loadCapabilityRegistryRecords(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}
	activeCapabilities, err := loadActiveCapabilityIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}
	if len(records) != len(activeCapabilities) {
		return ImplementationReadinessProgram{}, fmt.Errorf("capability metadata coverage mismatch: records=%d active=%d", len(records), len(activeCapabilities))
	}
	if len(records) != state.ProductSpec.KnownCapabilities {
		return ImplementationReadinessProgram{}, fmt.Errorf("capability readiness denominator mismatch: expected %d, got %d", state.ProductSpec.KnownCapabilities, len(records))
	}
	for id := range activeCapabilities {
		if _, ok := records[id]; !ok {
			return ImplementationReadinessProgram{}, fmt.Errorf("active capability %s has no readiness metadata record", id)
		}
	}

	activeOpen, err := loadActiveOpenDecisionIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}
	dependencies, err := loadCapabilityDependencyEvidence(root, records)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}

	implemented := map[string]string{}
	for _, entry := range ledger.Entries {
		if _, exists := implemented[entry.Capability]; exists {
			return ImplementationReadinessProgram{}, fmt.Errorf("duplicate implementation ledger capability %s", entry.Capability)
		}
		implemented[entry.Capability] = entry.WorkUnit
	}

	program, err := compileImplementationReadiness(records, activeOpen, dependencies, implemented, state.ProductSpec.BaselineCommit, ledgerAudit.SpecTreeDigest)
	if err != nil {
		return ImplementationReadinessProgram{}, err
	}
	if program.Summary.Implemented != ledgerAudit.VerifiedClaims {
		return program, fmt.Errorf("readiness implemented count %d does not match verified ledger claims %d", program.Summary.Implemented, ledgerAudit.VerifiedClaims)
	}
	if program.Summary.Total != state.ProductSpec.KnownCapabilities {
		return program, fmt.Errorf("readiness output does not cover the full active capability denominator")
	}
	program.Status = "PASS"
	return program, nil
}

func loadCapabilityRegistryRecords(root, specRel string) (map[string]CapabilityRegistryRecord, error) {
	dirRel := filepath.ToSlash(filepath.Join(specRel, "00-governance", "registers"))
	out := map[string]CapabilityRegistryRecord{}
	err := walkRepoDir(root, dirRel, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "capability-register-") || !strings.HasSuffix(name, ".md") {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		source := filepath.ToSlash(rel)
		data, err := readRepoFile(root, source)
		if err != nil {
			return err
		}
		records, err := parseCapabilityRegistryContent(source, string(data))
		if err != nil {
			return err
		}
		for _, record := range records {
			if existing, ok := out[record.ID]; ok {
				return fmt.Errorf("capability %s appears in multiple registry shards: %s and %s", record.ID, existing.SourcePath, record.SourcePath)
			}
			out[record.ID] = record
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("capability readiness parser found no capability records")
	}
	return out, nil
}

func parseCapabilityRegistryContent(source, content string) ([]CapabilityRegistryRecord, error) {
	var columns map[string]int
	var records []CapabilityRegistryRecord
	started := false
	for _, line := range strings.Split(content, "\n") {
		cells, ok := markdownTableCells(line)
		if !ok {
			if started && len(records) > 0 {
				break
			}
			continue
		}
		if markdownSeparatorRow(cells) {
			continue
		}
		if columns == nil {
			candidate := tableColumnMap(cells)
			hasID := tableColumn(candidate, "id", "capability id") >= 0
			hasDelivery := tableColumn(candidate, "delivery status") >= 0 || tableColumn(candidate, "delivery") >= 0
			if hasID && hasDelivery {
				columns = candidate
				started = true
			}
			continue
		}
		id := strings.Trim(tableCell(cells, columns, "id", "capability id"), "`*")
		if !capabilityIDPattern.MatchString(id) {
			continue
		}

		deliveryStatus := strings.ToLower(strings.TrimSpace(tableCell(cells, columns, "delivery status")))
		deliveryMode := strings.ToLower(strings.TrimSpace(tableCell(cells, columns, "delivery mode")))
		if deliveryStatus == "" {
			combined := strings.ToLower(strings.TrimSpace(tableCell(cells, columns, "delivery")))
			parts := strings.Split(combined, "/")
			if len(parts) != 2 {
				return nil, fmt.Errorf("%s capability %s has unsupported combined delivery value %q", source, id, combined)
			}
			deliveryStatus = strings.TrimSpace(parts[0])
			deliveryMode = strings.TrimSpace(parts[1])
		}
		if deliveryStatus != "defined" && deliveryStatus != "proposed" {
			return nil, fmt.Errorf("%s capability %s has unsupported delivery status %q", source, id, deliveryStatus)
		}
		if deliveryMode == "" {
			return nil, fmt.Errorf("%s capability %s has no delivery mode", source, id)
		}

		record := CapabilityRegistryRecord{
			ID:             id,
			Name:           strings.TrimSpace(tableCell(cells, columns, "name", "capability")),
			Status:         strings.ToLower(strings.TrimSpace(tableCell(cells, columns, "status"))),
			DeliveryStatus: deliveryStatus,
			DeliveryMode:   deliveryMode,
			CanonicalFile:  strings.Trim(strings.TrimSpace(tableCell(cells, columns, "canonical file")), "`"),
			OpenDecisions:  extractOpenDecisionIDs(tableCell(cells, columns, "open")),
			SourcePath:     source,
		}
		if record.Name == "" {
			return nil, fmt.Errorf("%s capability %s has no name", source, id)
		}
		records = append(records, record)
	}
	if columns == nil {
		return nil, fmt.Errorf("%s contains no capability register table", source)
	}
	return records, nil
}

func loadCapabilityDependencyEvidence(root string, known map[string]CapabilityRegistryRecord) ([]DependencyReadinessEvidence, error) {
	data, err := readRepoFile(root, capabilityDependencyRegisterPath)
	if err != nil {
		return nil, err
	}
	return parseDependencyRegisterContent(capabilityDependencyRegisterPath, string(data), known)
}

func parseDependencyRegisterContent(source, content string, known map[string]CapabilityRegistryRecord) ([]DependencyReadinessEvidence, error) {
	var columns map[string]int
	var out []DependencyReadinessEvidence
	started := false
	for _, line := range strings.Split(content, "\n") {
		cells, ok := markdownTableCells(line)
		if !ok {
			if started && len(out) > 0 {
				break
			}
			continue
		}
		if markdownSeparatorRow(cells) {
			continue
		}
		if columns == nil {
			candidate := tableColumnMap(cells)
			if tableColumn(candidate, "id") >= 0 && tableColumn(candidate, "dependent") >= 0 &&
				tableColumn(candidate, "status") >= 0 && tableColumn(candidate, "blocking") >= 0 {
				columns = candidate
				started = true
			}
			continue
		}
		id := strings.Trim(tableCell(cells, columns, "id"), "`* ")
		if !readinessDependencyIDPattern.MatchString(id) {
			continue
		}
		dependent := strings.TrimSpace(tableCell(cells, columns, "dependent"))
		evidence := DependencyReadinessEvidence{
			ID:                   id,
			Dependent:            dependent,
			Status:               strings.ToLower(strings.TrimSpace(tableCell(cells, columns, "status"))),
			Blocking:             strings.TrimSpace(tableCell(cells, columns, "blocking")),
			OpenDecisions:        extractOpenDecisionIDs(strings.Join(cells, " ")),
			AffectedCapabilities: expandCapabilitySelectors(dependent, known),
			SourcePath:           source,
		}
		if evidence.Dependent == "" || evidence.Status == "" || evidence.Blocking == "" {
			return nil, fmt.Errorf("dependency %s has incomplete readiness metadata", id)
		}
		out = append(out, evidence)
	}
	if columns == nil {
		return nil, fmt.Errorf("%s contains no dependency register table", source)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func compileImplementationReadiness(records map[string]CapabilityRegistryRecord, activeOpen map[string][]string, dependencies []DependencyReadinessEvidence, implemented map[string]string, baseline, digest string) (ImplementationReadinessProgram, error) {
	program := ImplementationReadinessProgram{
		SchemaVersion:       1,
		ProgramKind:         "capability-implementation-readiness",
		ProductSpecBaseline: baseline,
		SpecTreeDigest:      digest,
	}
	dependencyBlockers := map[string][]ImplementationReadinessBlocker{}
	for _, dependency := range dependencies {
		if !dependencyBlocksImplementation(dependency) {
			continue
		}
		if len(dependency.AffectedCapabilities) == 0 {
			program.UnscopedBlockingDependencies = append(program.UnscopedBlockingDependencies, dependency)
			continue
		}
		for _, capability := range dependency.AffectedCapabilities {
			var activeRefs []string
			for _, id := range dependency.OpenDecisions {
				if _, ok := activeOpen[id]; ok {
					activeRefs = append(activeRefs, id)
				}
			}
			sort.Strings(activeRefs)
			dependencyBlockers[capability] = append(dependencyBlockers[capability], ImplementationReadinessBlocker{
				Kind:          "dependency",
				ID:            dependency.ID,
				Status:        dependency.Status,
				Blocking:      dependency.Blocking,
				OpenDecisions: activeRefs,
				SourcePaths:   []string{dependency.SourcePath},
			})
		}
	}

	ids := make([]string, 0, len(records))
	for id := range records {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		record := records[id]
		var blockers []ImplementationReadinessBlocker
		for _, openID := range record.OpenDecisions {
			sources, active := activeOpen[openID]
			if !active {
				continue
			}
			sourcePaths := append([]string{record.SourcePath}, sources...)
			sourcePaths = sortedUniqueStrings(sourcePaths)
			blockers = append(blockers, ImplementationReadinessBlocker{
				Kind:        "open-decision",
				ID:          openID,
				Status:      "open",
				Blocking:    "capability register OPEN reference",
				SourcePaths: sourcePaths,
			})
		}
		blockers = append(blockers, dependencyBlockers[id]...)
		sort.Slice(blockers, func(i, j int) bool {
			if blockers[i].Kind != blockers[j].Kind {
				return blockers[i].Kind < blockers[j].Kind
			}
			return blockers[i].ID < blockers[j].ID
		})

		readiness := ImplementationReadinessRecord{
			Capability:     id,
			Name:           record.Name,
			DeliveryStatus: record.DeliveryStatus,
			DeliveryMode:   record.DeliveryMode,
			RegisterPath:   record.SourcePath,
			CanonicalFile:  record.CanonicalFile,
			Blockers:       blockers,
		}
		workUnit, isImplemented := implemented[id]
		switch {
		case record.DeliveryStatus == "proposed":
			if isImplemented {
				return program, fmt.Errorf("proposed capability %s cannot have a verified implementation claim", id)
			}
			readiness.State = "PROPOSED"
			program.Summary.Proposed++
		case isImplemented:
			readiness.State = "IMPLEMENTED"
			readiness.ImplementedWorkUnit = workUnit
			program.Summary.Implemented++
		case len(blockers) > 0:
			readiness.State = "BLOCKED"
			program.Summary.Blocked++
		default:
			readiness.State = "READY"
			program.Summary.Ready++
		}
		if readiness.State == "BLOCKED" {
			hasDecision, hasDependency := false, false
			for _, blocker := range blockers {
				switch blocker.Kind {
				case "open-decision":
					hasDecision = true
				case "dependency":
					hasDependency = true
				default:
					return program, fmt.Errorf("capability %s has unknown blocker kind %q", id, blocker.Kind)
				}
			}
			if hasDecision {
				program.Summary.DecisionBlocked++
			}
			if hasDependency {
				program.Summary.DependencyBlocked++
			}
		}
		program.Records = append(program.Records, readiness)
	}
	program.Summary.Total = len(program.Records)
	program.Summary.UnscopedBlockingDependencies = len(program.UnscopedBlockingDependencies)
	if program.Summary.Total != program.Summary.Implemented+program.Summary.Ready+program.Summary.Blocked+program.Summary.Proposed {
		return program, fmt.Errorf("readiness state accounting mismatch")
	}
	return program, nil
}

func dependencyBlocksImplementation(dependency DependencyReadinessEvidence) bool {
	if strings.EqualFold(strings.TrimSpace(dependency.Status), "active") {
		return false
	}
	blocking := strings.ToLower(strings.TrimSpace(dependency.Blocking))
	return strings.Contains(blocking, "yes") ||
		strings.Contains(blocking, "before ") ||
		strings.Contains(blocking, "blocking") ||
		strings.Contains(blocking, "required")
}

func expandCapabilitySelectors(text string, known map[string]CapabilityRegistryRecord) []string {
	set := map[string]bool{}
	for _, token := range capabilitySelectorTokenPattern.FindAllString(text, -1) {
		dash := strings.LastIndex(token, "-")
		if dash < 0 || dash == len(token)-1 {
			continue
		}
		prefix, selector := token[:dash+1], token[dash+1:]
		switch {
		case strings.Contains(selector, ".."):
			parts := strings.SplitN(selector, "..", 2)
			start, errStart := strconv.Atoi(parts[0])
			end, errEnd := strconv.Atoi(parts[1])
			if errStart != nil || errEnd != nil || start > end {
				continue
			}
			for n := start; n <= end; n++ {
				id := fmt.Sprintf("%s%03d", prefix, n)
				if _, ok := known[id]; ok {
					set[id] = true
				}
			}
		case strings.Contains(selector, "/"):
			for _, part := range strings.Split(selector, "/") {
				n, err := strconv.Atoi(part)
				if err != nil {
					continue
				}
				id := fmt.Sprintf("%s%03d", prefix, n)
				if _, ok := known[id]; ok {
					set[id] = true
				}
			}
		default:
			if _, ok := known[token]; ok {
				set[token] = true
			}
		}
	}
	var ids []string
	for id := range set {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func markdownTableCells(line string) ([]string, bool) {
	line = strings.TrimSpace(line)
	if len(line) < 2 || !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
		return nil, false
	}
	raw := strings.Split(line[1:len(line)-1], "|")
	cells := make([]string, len(raw))
	for i, value := range raw {
		cells[i] = strings.TrimSpace(value)
	}
	return cells, true
}

func markdownSeparatorRow(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		value := strings.Trim(strings.TrimSpace(cell), ":")
		if len(value) < 3 || strings.Trim(value, "-") != "" {
			return false
		}
	}
	return true
}

func tableColumnMap(cells []string) map[string]int {
	out := map[string]int{}
	for i, cell := range cells {
		key := strings.ToLower(strings.Trim(strings.TrimSpace(cell), "`*"))
		out[key] = i
	}
	return out
}

func tableColumn(columns map[string]int, names ...string) int {
	for _, name := range names {
		if index, ok := columns[strings.ToLower(name)]; ok {
			return index
		}
	}
	return -1
}

func tableCell(cells []string, columns map[string]int, names ...string) string {
	index := tableColumn(columns, names...)
	if index < 0 || index >= len(cells) {
		return ""
	}
	return cells[index]
}

func extractOpenDecisionIDs(value string) []string {
	ids := openDecisionTokenPattern.FindAllString(value, -1)
	return sortedUniqueStrings(ids)
}

func sortedUniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	set := map[string]bool{}
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			set[value] = true
		}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
