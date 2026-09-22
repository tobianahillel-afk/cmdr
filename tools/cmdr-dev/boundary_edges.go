package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type LocalPackage struct {
	Boundary string `json:"boundary"`
	Ecosystem string `json:"ecosystem"`
	Name string `json:"name"`
	Manifest string `json:"manifest"`
	Directory string `json:"directory"`
}

type BoundaryDependencyEdge struct {
	From string `json:"from"`
	To string `json:"to"`
	Kind string `json:"kind"`
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	Dependency string `json:"dependency"`
}

type BoundaryEdgeAuditSummary struct {
	RuntimeBoundaries int `json:"runtime_boundaries"`
	LocalPackages int `json:"local_packages"`
	ObservedEdges int `json:"observed_edges"`
	CrossBoundaryEdges int `json:"cross_boundary_edges"`
}

type localDependencyRef struct {
	Name string
	Spec string
}

func runBoundaryEdgeAudit(root string) (BoundaryEdgeAuditSummary, error) {
	var registry ArchitectureRegistry
	if err := decodeStrict(filepath.Join(root, filepath.FromSlash(architectureRegistryPath)), &registry); err != nil {
		return BoundaryEdgeAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return BoundaryEdgeAuditSummary{}, err
	}

	packages, runtimeBoundaries, err := discoverLocalPackages(root, registry)
	if err != nil {
		return BoundaryEdgeAuditSummary{}, err
	}
	edges, err := discoverBoundaryEdges(root, registry, packages)
	if err != nil {
		return BoundaryEdgeAuditSummary{}, err
	}
	if err := validateBoundaryEdges(registry, edges); err != nil {
		return BoundaryEdgeAuditSummary{}, err
	}
	summary := BoundaryEdgeAuditSummary{
		RuntimeBoundaries: runtimeBoundaries,
		LocalPackages: len(packages),
		ObservedEdges: len(edges),
	}
	for _, edge := range edges {
		if edge.From != edge.To {
			summary.CrossBoundaryEdges++
		}
	}
	return summary, nil
}

func discoverLocalPackages(root string, registry ArchitectureRegistry) ([]LocalPackage, int, error) {
	var packages []LocalPackage
	runtimeBoundaries := 0
	for _, boundary := range registry.Boundaries {
		if boundary.Kind != "product-runtime" {
			continue
		}
		runtimeBoundaries++
		for _, pattern := range boundary.Roots {
			scanRoot, err := runtimeScanRoot(root, pattern)
			if err != nil {
				return nil, 0, err
			}
			err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if entry.IsDir() {
					if path != scanRoot && ignoredDependencyDir(entry.Name()) {
						return filepath.SkipDir
					}
					return nil
				}
				if entry.Name() != "package.json" && entry.Name() != "go.mod" {
					return nil
				}
				rel, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				rel = filepath.ToSlash(rel)
				dir := filepath.ToSlash(filepath.Dir(rel))
				var name, ecosystem string
				switch entry.Name() {
				case "package.json":
					name, _, err = parseNodeLocalRefs(path)
					ecosystem = "npm"
				case "go.mod":
					name, _, err = parseGoLocalRefs(path)
					ecosystem = "go"
				}
				if err != nil {
					return err
				}
				if name == "" {
					return fmt.Errorf("%s has no local package/module identity", rel)
				}
				packages = append(packages, LocalPackage{
					Boundary: boundary.ID,
					Ecosystem: ecosystem,
					Name: name,
					Manifest: rel,
					Directory: dir,
				})
				return nil
			})
			if err != nil {
				return nil, 0, err
			}
		}
	}
	sort.Slice(packages, func(i, j int) bool {
		if packages[i].Ecosystem != packages[j].Ecosystem {
			return packages[i].Ecosystem < packages[j].Ecosystem
		}
		if packages[i].Name != packages[j].Name {
			return packages[i].Name < packages[j].Name
		}
		return packages[i].Manifest < packages[j].Manifest
	})
	seen := map[string]string{}
	for _, pkg := range packages {
		key := pkg.Ecosystem + ":" + pkg.Name
		if previous, exists := seen[key]; exists {
			return nil, 0, fmt.Errorf("duplicate local package identity %s in %s and %s", key, previous, pkg.Manifest)
		}
		seen[key] = pkg.Manifest
	}
	return packages, runtimeBoundaries, nil
}

func discoverBoundaryEdges(root string, registry ArchitectureRegistry, packages []LocalPackage) ([]BoundaryDependencyEdge, error) {
	packageIndex := map[string]LocalPackage{}
	for _, pkg := range packages {
		packageIndex[pkg.Ecosystem+":"+pkg.Name] = pkg
	}
	var edges []BoundaryDependencyEdge
	for _, pkg := range packages {
		full := filepath.Join(root, filepath.FromSlash(pkg.Manifest))
		switch pkg.Ecosystem {
		case "npm":
			_, refs, err := parseNodeLocalRefs(full)
			if err != nil {
				return nil, err
			}
			for _, ref := range refs {
				target, err := resolveNodeLocalTarget(root, registry, pkg, ref, packageIndex)
				if err != nil {
					return nil, err
				}
				edges = append(edges, BoundaryDependencyEdge{
					From: pkg.Boundary, To: target.Boundary, Kind: "package-local",
					SourcePath: pkg.Manifest, TargetPath: target.Directory, Dependency: ref.Name,
				})
			}
		case "go":
			_, refs, err := parseGoLocalRefs(full)
			if err != nil {
				return nil, err
			}
			for _, ref := range refs {
				targetDir := filepath.ToSlash(filepath.Clean(filepath.Join(pkg.Directory, filepath.FromSlash(ref.Spec))))
				targetBoundary, ok := boundaryForRepoPath(registry, targetDir)
				if !ok {
					return nil, fmt.Errorf("%s local Go replacement %s targets unowned path %s", pkg.Manifest, ref.Name, targetDir)
				}
				edges = append(edges, BoundaryDependencyEdge{
					From: pkg.Boundary, To: targetBoundary.ID, Kind: "package-local",
					SourcePath: pkg.Manifest, TargetPath: targetDir, Dependency: ref.Name,
				})
			}
		}
	}
	return dedupeBoundaryEdges(edges), nil
}

func parseNodeLocalRefs(path string) (string, []localDependencyRef, error) {
	var pkg struct {
		Name string `json:"name"`
		Dependencies map[string]string `json:"dependencies"`
		OptionalDependencies map[string]string `json:"optionalDependencies"`
		PeerDependencies map[string]string `json:"peerDependencies"`
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", nil, err
	}
	var refs []localDependencyRef
	merged := map[string]string{}
	for _, values := range []map[string]string{pkg.Dependencies, pkg.OptionalDependencies, pkg.PeerDependencies} {
		for name, version := range values {
			merged[name] = version
		}
	}
	for name, spec := range merged {
		if isLocalNodeDependency(spec) {
			refs = append(refs, localDependencyRef{Name: name, Spec: spec})
		}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
	return pkg.Name, refs, nil
}

func parseGoLocalRefs(path string) (string, []localDependencyRef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	lines := strings.Split(string(data), "\n")
	module := ""
	required := map[string]bool{}
	inRequire := false
	for _, raw := range lines {
		line := strings.TrimSpace(strings.SplitN(raw, "//", 2)[0])
		if strings.HasPrefix(line, "module ") {
			module = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}
		if strings.HasPrefix(line, "require ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "require "))
		} else if !inRequire {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			required[fields[0]] = true
		}
	}

	var refs []localDependencyRef
	inReplace := false
	for _, raw := range lines {
		line := strings.TrimSpace(strings.SplitN(raw, "//", 2)[0])
		if line == "replace (" {
			inReplace = true
			continue
		}
		if inReplace && line == ")" {
			inReplace = false
			continue
		}
		if strings.HasPrefix(line, "replace ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		} else if !inReplace && !strings.Contains(line, " => ") {
			continue
		}
		parts := strings.Split(line, " => ")
		if len(parts) != 2 {
			continue
		}
		left := strings.Fields(parts[0])
		right := strings.Fields(parts[1])
		if len(left) == 0 || len(right) == 0 || !required[left[0]] {
			continue
		}
		if strings.HasPrefix(right[0], "./") || strings.HasPrefix(right[0], "../") {
			refs = append(refs, localDependencyRef{Name: left[0], Spec: right[0]})
		}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
	return module, refs, nil
}

func resolveNodeLocalTarget(root string, registry ArchitectureRegistry, source LocalPackage, ref localDependencyRef, packages map[string]LocalPackage) (LocalPackage, error) {
	spec := strings.TrimSpace(ref.Spec)
	if strings.HasPrefix(spec, "workspace:") {
		target, ok := packages["npm:"+ref.Name]
		if !ok {
			return LocalPackage{}, fmt.Errorf("%s workspace dependency %s has no registered local package", source.Manifest, ref.Name)
		}
		return target, nil
	}
	for _, prefix := range []string{"file:", "link:"} {
		spec = strings.TrimPrefix(spec, prefix)
	}
	targetDir := filepath.ToSlash(filepath.Clean(filepath.Join(source.Directory, filepath.FromSlash(spec))))
	boundary, ok := boundaryForRepoPath(registry, targetDir)
	if !ok {
		return LocalPackage{}, fmt.Errorf("%s local dependency %s targets unowned path %s", source.Manifest, ref.Name, targetDir)
	}
	return LocalPackage{Boundary: boundary.ID, Ecosystem: "npm", Name: ref.Name, Directory: targetDir}, nil
}

func boundaryForRepoPath(registry ArchitectureRegistry, path string) (ArchitectureBoundary, bool) {
	path = filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	var found *ArchitectureBoundary
	for i := range registry.Boundaries {
		boundary := &registry.Boundaries[i]
		if boundary.Kind != "product-runtime" {
			continue
		}
		for _, pattern := range boundary.Roots {
			if repoPathOwnedByPattern(pattern, path) {
				if found != nil && found.ID != boundary.ID {
					return ArchitectureBoundary{}, false
				}
				copy := *boundary
				found = &copy
			}
		}
	}
	if found == nil {
		return ArchitectureBoundary{}, false
	}
	return *found, true
}

func repoPathOwnedByPattern(pattern, path string) bool {
	pattern = filepath.ToSlash(strings.TrimSpace(pattern))
	path = filepath.ToSlash(strings.TrimSpace(path))
	base := strings.TrimSuffix(pattern, "**")
	base = strings.TrimSuffix(base, "/")
	return path == base || strings.HasPrefix(path, base+"/")
}

func validateBoundaryEdges(registry ArchitectureRegistry, edges []BoundaryDependencyEdge) error {
	index := map[string]ArchitectureBoundary{}
	for _, boundary := range registry.Boundaries {
		index[boundary.ID] = boundary
	}
	for _, edge := range edges {
		source, ok := index[edge.From]
		if !ok {
			return fmt.Errorf("edge source boundary %s is unknown", edge.From)
		}
		if _, ok := index[edge.To]; !ok {
			return fmt.Errorf("edge target boundary %s is unknown", edge.To)
		}
		if edge.From == edge.To {
			continue
		}
		allowed := false
		for _, target := range source.MayDependOn {
			if target == edge.To {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("forbidden architecture edge %s -> %s from %s", edge.From, edge.To, edge.SourcePath)
		}
	}
	return nil
}

func dedupeBoundaryEdges(edges []BoundaryDependencyEdge) []BoundaryDependencyEdge {
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From { return edges[i].From < edges[j].From }
		if edges[i].To != edges[j].To { return edges[i].To < edges[j].To }
		if edges[i].SourcePath != edges[j].SourcePath { return edges[i].SourcePath < edges[j].SourcePath }
		return edges[i].Dependency < edges[j].Dependency
	})
	if len(edges) == 0 { return nil }
	out := edges[:0]
	var previous BoundaryDependencyEdge
	for i, edge := range edges {
		if i == 0 || edge != previous {
			out = append(out, edge)
			previous = edge
		}
	}
	return out
}
