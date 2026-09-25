package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type SpecDocument struct {
	Path          string   `json:"path"`
	SHA256        string   `json:"sha256"`
	ID            string   `json:"id,omitempty"`
	Type          string   `json:"type,omitempty"`
	Domain        string   `json:"domain,omitempty"`
	Product       string   `json:"product,omitempty"`
	Module        string   `json:"module,omitempty"`
	Status        string   `json:"status,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	Updated       string   `json:"updated,omitempty"`
	SourceOfTruth string   `json:"source_of_truth,omitempty"`
	Requirements  []string `json:"requirements,omitempty"`
	OpenDecisions []string `json:"open_decisions,omitempty"`
	Permissions   []string `json:"permissions,omitempty"`
	References    []string `json:"references,omitempty"`
	Active        bool     `json:"active"`
	Canonical     bool     `json:"canonical"`
}

type SpecInventory struct {
	SchemaVersion int            `json:"schema_version"`
	SpecRoot      string         `json:"spec_root"`
	TreeDigest    string         `json:"tree_digest"`
	Files         int            `json:"files"`
	Documents     []SpecDocument `json:"documents"`
}

type SpecIndexSummary struct {
	Files                    int    `json:"files"`
	ActiveCanonicalDocuments int    `json:"active_canonical_documents"`
	TreeDigest               string `json:"tree_digest"`
	Output                   string `json:"output"`
	Mode                     string `json:"mode"`
}

type SpecBaseline struct {
	SchemaVersion             int    `json:"schema_version"`
	CompilerSchemaVersion     int    `json:"compiler_schema_version"`
	SpecRoot                  string `json:"spec_root"`
	ProductSpecBaselineCommit string `json:"product_spec_baseline_commit"`
	TreeDigest                string `json:"tree_digest"`
	Files                     int    `json:"files"`
	ActiveCanonicalDocuments  int    `json:"active_canonical_documents"`
}

type SpecBaselineSummary struct {
	Files                     int    `json:"files"`
	ActiveCanonicalDocuments  int    `json:"active_canonical_documents"`
	TreeDigest                string `json:"tree_digest"`
	ProductSpecBaselineCommit string `json:"product_spec_baseline_commit"`
	Output                    string `json:"output"`
	Mode                      string `json:"mode"`
}

var cmrdIDPattern = regexp.MustCompile(`\b(?:CAP-[A-Z0-9]+-[0-9]{3}|REQ-[A-Z0-9]+-[0-9]{3}|OPEN-[0-9]{3}|ADR-[0-9]{4}|DEP-(?:[A-Z0-9]+-)?[0-9]{3})\b`)
var permissionPattern = regexp.MustCompile(`\bperm\.[a-zA-Z0-9.*_-]+(?:\.[a-zA-Z0-9.*_-]+)*\b`)

func runSpecIndex(root, specRel, output string, check bool) (SpecIndexSummary, error) {
	inventory, err := buildSpecInventory(root, specRel)
	if err != nil {
		return SpecIndexSummary{}, err
	}
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		return SpecIndexSummary{}, err
	}
	data = append(data, '\n')

	outputPath, err := resolveRepoPath(root, output, !check)
	if err != nil {
		return SpecIndexSummary{}, err
	}
	mode := "write"
	if check {
		mode = "check"
		existing, err := readRepoFile(root, outputPath)
		if err != nil {
			return SpecIndexSummary{}, fmt.Errorf("read generated spec index: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return SpecIndexSummary{}, fmt.Errorf("generated spec index is stale: %s", outputPath)
		}
	} else {
		outputPath, err = writeRepoFile(root, outputPath, data)
		if err != nil {
			return SpecIndexSummary{}, err
		}
	}

	activeCanonical := 0
	for _, doc := range inventory.Documents {
		if doc.Active && doc.Canonical {
			activeCanonical++
		}
	}
	return SpecIndexSummary{
		Files:                    inventory.Files,
		ActiveCanonicalDocuments: activeCanonical,
		TreeDigest:               inventory.TreeDigest,
		Output:                   outputPath,
		Mode:                     mode,
	}, nil
}

func buildSpecInventory(root, specRel string) (SpecInventory, error) {
	specRoot, err := resolveRepoPath(root, specRel, false)
	if err != nil {
		return SpecInventory{}, err
	}
	var docs []SpecDocument
	// #nosec G703 -- specRoot is repository-confined by resolveRepoPath; symlink entries are rejected below.
	err = filepath.WalkDir(specRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("product spec contains symlink: %s", path)
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			return nil
		}
		content, err := readRepoFile(root, path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		doc := parseSpecDocument(rel, content)
		docs = append(docs, doc)
		return nil
	})
	if err != nil {
		return SpecInventory{}, err
	}
	sort.Slice(docs, func(i, j int) bool { return docs[i].Path < docs[j].Path })
	if err := rejectDuplicateCanonicalIDs(docs); err != nil {
		return SpecInventory{}, err
	}

	h := sha256.New()
	for _, doc := range docs {
		h.Write([]byte(doc.Path))
		h.Write([]byte{0})
		h.Write([]byte(doc.SHA256))
		h.Write([]byte{0})
	}
	return SpecInventory{
		SchemaVersion: 1,
		SpecRoot:      filepath.ToSlash(specRel),
		TreeDigest:    hex.EncodeToString(h.Sum(nil)),
		Files:         len(docs),
		Documents:     docs,
	}, nil
}

func parseSpecDocument(path string, content []byte) SpecDocument {
	sum := sha256.Sum256(content)
	meta := parseFrontMatter(string(content))
	refs := extractReferences(string(content))

	status := strings.ToLower(meta.scalar["status"])
	active := !strings.Contains("/"+path, "/99-archive/") && status != "deprecated" && status != "archived"
	canonical := strings.EqualFold(meta.scalar["source-of-truth"], "canonical")

	reqs := append([]string(nil), meta.lists["requirements"]...)
	reqs = append(reqs, meta.lists["requirement_ids"]...)
	openDecisions := append([]string(nil), meta.lists["open_decisions"]...)
	permissions := append([]string(nil), meta.lists["permissions"]...)

	return SpecDocument{
		Path:          path,
		SHA256:        hex.EncodeToString(sum[:]),
		ID:            meta.scalar["id"],
		Type:          meta.scalar["type"],
		Domain:        meta.scalar["domain"],
		Product:       meta.scalar["product"],
		Module:        meta.scalar["module"],
		Status:        meta.scalar["status"],
		Owner:         meta.scalar["owner"],
		Updated:       meta.scalar["updated"],
		SourceOfTruth: meta.scalar["source-of-truth"],
		Requirements:  uniqueSorted(reqs),
		OpenDecisions: uniqueSorted(openDecisions),
		Permissions:   uniqueSorted(permissions),
		References:    refs,
		Active:        active,
		Canonical:     canonical,
	}
}

type frontMatter struct {
	scalar map[string]string
	lists  map[string][]string
}

func parseFrontMatter(content string) frontMatter {
	out := frontMatter{scalar: map[string]string{}, lists: map[string][]string{}}
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return out
	}
	currentList := ""
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			break
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "- ") && currentList != "" {
			value := cleanYAMLScalar(strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
			if value != "" {
				out.lists[currentList] = append(out.lists[currentList], value)
			}
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			currentList = ""
			continue
		}
		key := strings.TrimSpace(line[:idx])
		rawValue := strings.TrimSpace(line[idx+1:])
		if values, ok, err := parseInlineYAMLList(rawValue); ok {
			if err != nil {
				currentList = ""
				continue
			}
			out.lists[key] = append(out.lists[key], values...)
			currentList = ""
			continue
		}
		value := cleanYAMLScalar(rawValue)
		if value == "" {
			currentList = key
			if _, ok := out.lists[key]; !ok {
				out.lists[key] = nil
			}
			continue
		}
		currentList = ""
		out.scalar[key] = value
	}
	return out
}

func parseInlineYAMLList(value string) ([]string, bool, error) {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "[") {
		return nil, false, nil
	}
	if !strings.HasSuffix(value, "]") {
		return nil, true, fmt.Errorf("unterminated inline YAML list")
	}
	body := strings.TrimSpace(value[1 : len(value)-1])
	if body == "" {
		return nil, true, nil
	}
	var values []string
	var current strings.Builder
	var quote rune
	for _, r := range body {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
		case r == ',':
			item := cleanYAMLScalar(current.String())
			if item == "" {
				return nil, true, fmt.Errorf("empty item in inline YAML list")
			}
			values = append(values, item)
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	if quote != 0 {
		return nil, true, fmt.Errorf("unterminated quote in inline YAML list")
	}
	item := cleanYAMLScalar(current.String())
	if item == "" {
		return nil, true, fmt.Errorf("empty item in inline YAML list")
	}
	values = append(values, item)
	return values, true, nil
}

func cleanYAMLScalar(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func extractReferences(content string) []string {
	seen := map[string]struct{}{}
	for _, match := range cmrdIDPattern.FindAllString(content, -1) {
		seen[match] = struct{}{}
	}
	for _, match := range permissionPattern.FindAllString(content, -1) {
		seen[match] = struct{}{}
	}
	refs := make([]string, 0, len(seen))
	for ref := range seen {
		refs = append(refs, ref)
	}
	sort.Strings(refs)
	return refs
}

func rejectDuplicateCanonicalIDs(docs []SpecDocument) error {
	seen := map[string]string{}
	for _, doc := range docs {
		if !doc.Active || !doc.Canonical || doc.ID == "" {
			continue
		}
		if previous, ok := seen[doc.ID]; ok {
			return fmt.Errorf("duplicate active canonical document id %q in %s and %s", doc.ID, previous, doc.Path)
		}
		seen[doc.ID] = doc.Path
	}
	return nil
}

func uniqueSorted(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	out := values[:0]
	var last string
	for i, value := range values {
		if i == 0 || value != last {
			out = append(out, value)
			last = value
		}
	}
	return out
}

func runSpecBaseline(root, specRel, baselineCommit, output string, check bool) (SpecBaselineSummary, error) {
	inventory, err := buildSpecInventory(root, specRel)
	if err != nil {
		return SpecBaselineSummary{}, err
	}
	activeCanonical := 0
	for _, doc := range inventory.Documents {
		if doc.Active && doc.Canonical {
			activeCanonical++
		}
	}
	baseline := SpecBaseline{
		SchemaVersion:             1,
		CompilerSchemaVersion:     2,
		SpecRoot:                  inventory.SpecRoot,
		ProductSpecBaselineCommit: baselineCommit,
		TreeDigest:                inventory.TreeDigest,
		Files:                     inventory.Files,
		ActiveCanonicalDocuments:  activeCanonical,
	}
	data, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		return SpecBaselineSummary{}, err
	}
	data = append(data, '\n')

	outputPath, err := resolveRepoPath(root, output, !check)
	if err != nil {
		return SpecBaselineSummary{}, err
	}
	mode := "write"
	if check {
		mode = "check"
		existing, err := readRepoFile(root, outputPath)
		if err != nil {
			return SpecBaselineSummary{}, fmt.Errorf("read spec baseline: %w", err)
		}
		if !bytes.Equal(existing, data) {
			return SpecBaselineSummary{}, fmt.Errorf("spec baseline is stale: %s", outputPath)
		}
	} else {
		outputPath, err = writeRepoFile(root, outputPath, data)
		if err != nil {
			return SpecBaselineSummary{}, err
		}
	}

	return SpecBaselineSummary{
		Files:                     inventory.Files,
		ActiveCanonicalDocuments:  activeCanonical,
		TreeDigest:                inventory.TreeDigest,
		ProductSpecBaselineCommit: baselineCommit,
		Output:                    outputPath,
		Mode:                      mode,
	}, nil
}
