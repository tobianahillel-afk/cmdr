package harness

import (
	"fmt"
	"strings"
)

// WorkflowStep represents a parsed step inside a GitHub Actions workflow job.
type WorkflowStep struct {
	Job             string
	Name            string
	Run             string
	Uses            string
	IfExpr          string
	ContinueOnError bool
	Line            int
}

type workflowJob struct {
	Name            string
	IfExpr          string
	ContinueOnError bool
	Line            int
}

// WorkflowValidationResult captures the structured validation outcome for the mainline health gate workflow.
type WorkflowValidationResult struct {
	Valid              bool
	HasHealthGateStep  bool
	HasContinueOnError bool
	HasDisabledIf      bool
	HasPushMain        bool
	HasPRMain          bool
	HealthGateSteps    []WorkflowStep
	Errors             []string
}

// stripYAMLComment strips unquoted `# ...` comments from a single YAML line while preserving `#` inside quotes.
func stripYAMLComment(line string) string {
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(line); i++ {
		c := line[i]
		if inDouble {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inDouble = false
			}
			continue
		}
		if inSingle {
			if c == '\'' {
				if i+1 < len(line) && line[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}

		switch c {
		case '\'':
			inSingle = true
		case '"':
			inDouble = true
		case '#':
			if i == 0 || line[i-1] == ' ' || line[i-1] == '\t' {
				return strings.TrimRight(line[:i], " \r")
			}
		}
	}
	return strings.TrimRight(line, " \r")
}

func leadingSpaces(s string) int {
	n := 0
	for n < len(s) && s[n] == ' ' {
		n++
	}
	return n
}

func unquoteYAMLValue(val string) string {
	v := strings.TrimSpace(val)
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}

func splitYAMLKeyValue(s string) (key string, val string, ok bool) {
	idx := strings.IndexByte(s, ':')
	if idx <= 0 {
		return "", "", false
	}
	key = strings.TrimSpace(s[:idx])
	rest := s[idx+1:]
	if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return "", "", false
	}
	return key, strings.TrimSpace(rest), true
}

func inlineListHasMain(val string) bool {
	trimmed := strings.TrimSpace(val)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return false
	}
	inner := trimmed[1 : len(trimmed)-1]
	for _, part := range strings.Split(inner, ",") {
		if unquoteYAMLValue(part) == "main" {
			return true
		}
	}
	return false
}

func isBlockScalarHeader(val string) bool {
	v := strings.TrimSpace(val)
	switch v {
	case "|", "|-", "|+", ">", ">-", ">+":
		return true
	default:
		return false
	}
}

func normalizeWorkflowExpr(expr string) string {
	s := strings.ToLower(strings.TrimSpace(unquoteYAMLValue(expr)))
	if strings.HasPrefix(s, "${{") && strings.HasSuffix(s, "}}") {
		s = strings.TrimSpace(s[3 : len(s)-2])
	}
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, "! ", "!")
	s = strings.ReplaceAll(s, "( ", "(")
	s = strings.ReplaceAll(s, " )", ")")
	return s
}

func isDisabledIfExpr(expr string) bool {
	s := normalizeWorkflowExpr(expr)
	if s == "" {
		return false
	}
	switch s {
	case "false", "(false)", "!true", "!(true)", "false || false":
		return true
	}
	if strings.HasPrefix(s, "false && ") || strings.HasSuffix(s, " && false") || strings.Contains(s, " && false && ") {
		return true
	}
	if strings.HasPrefix(s, "!true && ") || strings.HasSuffix(s, " && !true") || strings.Contains(s, " && !true && ") {
		return true
	}
	return false
}

func isTruthyContinueOnError(val string) bool {
	switch normalizeWorkflowExpr(val) {
	case "true", "!false":
		return true
	default:
		return false
	}
}

func splitShellSegments(line string) []string {
	var segments []string
	var cur strings.Builder
	inSingle := false
	inDouble := false
	escaped := false

	for i := 0; i < len(line); i++ {
		c := line[i]
		if inDouble {
			cur.WriteByte(c)
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inDouble = false
			}
			continue
		}
		if inSingle {
			cur.WriteByte(c)
			if c == '\'' {
				inSingle = false
			}
			continue
		}
		switch c {
		case '\'':
			inSingle = true
			cur.WriteByte(c)
		case '"':
			inDouble = true
			cur.WriteByte(c)
		case ';':
			segments = append(segments, cur.String())
			cur.Reset()
		case '&', '|':
			if i+1 < len(line) && line[i+1] == c {
				segments = append(segments, cur.String())
				cur.Reset()
				i++
			} else {
				cur.WriteByte(c)
			}
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		segments = append(segments, cur.String())
	}
	return segments
}

func isMainlineHealthGateSegment(seg string) bool {
	tokens := strings.Fields(seg)
	for len(tokens) > 0 && strings.Contains(tokens[0], "=") && !strings.HasPrefix(tokens[0], "-") {
		tokens = tokens[1:]
	}
	if len(tokens) < 2 {
		return false
	}
	cmd := strings.ReplaceAll(tokens[0], "\\", "/")
	if idx := strings.LastIndexByte(cmd, '/'); idx >= 0 {
		cmd = cmd[idx+1:]
	}
	if cmd == "cmdr-dev" || cmd == "cmdr-dev.exe" {
		return tokens[1] == "mainline-health-gate"
	}
	if cmd == "go" && tokens[1] == "run" && len(tokens) >= 4 {
		for _, tok := range tokens[2:] {
			if tok == "mainline-health-gate" {
				return true
			}
		}
	}
	return false
}

func containsMainlineHealthGateInvocation(runScript string) bool {
	for _, rawLine := range strings.Split(runScript, "\n") {
		line := strings.TrimSpace(stripYAMLComment(rawLine))
		if line == "" {
			continue
		}
		for _, seg := range splitShellSegments(line) {
			if isMainlineHealthGateSegment(seg) {
				return true
			}
		}
	}
	return false
}

// ValidateWorkflowRequiredJob checks that the workflow YAML contains no tabs, has a top-level
// `jobs:` block, and defines an uncommented active job matching requiredJob.
func ValidateWorkflowRequiredJob(content string, requiredJob string) error {
	if strings.TrimSpace(requiredJob) == "" {
		return fmt.Errorf("requiredJob must be non-empty")
	}
	if strings.Contains(content, "\t") {
		return fmt.Errorf("workflow contains forbidden TAB characters in YAML")
	}

	rawLines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	inJobs := false
	jobsIndent := -1
	jobItemIndent := -1
	inBlockScalar := false
	blockScalarIndent := -1

	for _, raw := range rawLines {
		if inBlockScalar {
			trimmedRaw := strings.TrimSpace(raw)
			if trimmedRaw == "" {
				continue
			}
			if leadingSpaces(raw) > blockScalarIndent {
				continue
			}
			inBlockScalar = false
			blockScalarIndent = -1
		}

		clean := stripYAMLComment(raw)
		trimmed := strings.TrimSpace(clean)
		if trimmed == "" {
			continue
		}
		indent := leadingSpaces(clean)

		if k, v, ok := splitYAMLKeyValue(trimmed); ok && isBlockScalarHeader(v) {
			inBlockScalar = true
			blockScalarIndent = indent
			_ = k
		}

		if !inJobs {
			if indent == 0 {
				if k, v, ok := splitYAMLKeyValue(trimmed); ok && k == "jobs" && v == "" {
					inJobs = true
					jobsIndent = indent
					jobItemIndent = -1
				}
			}
			continue
		}

		if indent <= jobsIndent {
			inJobs = false
			jobItemIndent = -1
			continue
		}

		if jobItemIndent == -1 || indent == jobItemIndent {
			if k, _, ok := splitYAMLKeyValue(trimmed); ok {
				jobItemIndent = indent
				if unquoteYAMLValue(k) == requiredJob {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("workflow missing required job %q", requiredJob)
}

// ValidateWorkflowSyntax performs pure-stdlib structural YAML syntax validation for GitHub Actions workflows.
func ValidateWorkflowSyntax(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("workflow YAML is empty")
	}
	if strings.Contains(content, "\t") {
		return fmt.Errorf("workflow contains forbidden TAB characters in YAML")
	}

	rawLines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	bracketDepth := 0
	braceDepth := 0
	inBlockScalar := false
	blockScalarIndent := -1

	// Track active sequence-item indentation to catch misaligned mapping keys (e.g. `- name: ...` followed by dedented `run:`).
	type seqFrame struct {
		seqIndent     int
		dashIndent    int
		itemKeyIndent int
	}
	var seqStack []seqFrame

	for lineIdx, raw := range rawLines {
		lineNum := lineIdx + 1
		if inBlockScalar {
			trimmedRaw := strings.TrimSpace(raw)
			if trimmedRaw == "" {
				continue
			}
			if leadingSpaces(raw) > blockScalarIndent {
				continue
			}
			inBlockScalar = false
			blockScalarIndent = -1
		}

		clean := stripYAMLComment(raw)
		trimmed := strings.TrimSpace(clean)
		if trimmed == "" {
			continue
		}
		indent := leadingSpaces(clean)

		// Pop sequence frames whose scope we exited.
		for len(seqStack) > 0 {
			top := seqStack[len(seqStack)-1]
			if indent < top.seqIndent || (indent == top.seqIndent && !strings.HasPrefix(trimmed, "-") && trimmed != "-") {
				seqStack = seqStack[:len(seqStack)-1]
				continue
			}
			break
		}

		// Verify unclosed quotes on this non-block-scalar line and update bracket/brace balance outside quotes.
		inSingle := false
		inDouble := false
		escaped := false
		for i := 0; i < len(trimmed); i++ {
			c := trimmed[i]
			if inDouble {
				if escaped {
					escaped = false
					continue
				}
				if c == '\\' {
					escaped = true
					continue
				}
				if c == '"' {
					inDouble = false
				}
				continue
			}
			if inSingle {
				if c == '\'' {
					if i+1 < len(trimmed) && trimmed[i+1] == '\'' {
						i++
						continue
					}
					inSingle = false
				}
				continue
			}
			switch c {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '[':
				bracketDepth++
			case ']':
				bracketDepth--
				if bracketDepth < 0 {
					return fmt.Errorf("line %d: unexpected closing ']'", lineNum)
				}
			case '{':
				braceDepth++
			case '}':
				braceDepth--
				if braceDepth < 0 {
					return fmt.Errorf("line %d: unexpected closing '}'", lineNum)
				}
			}
		}
		if inSingle || inDouble {
			return fmt.Errorf("line %d: unclosed string quote in YAML", lineNum)
		}

		// Check sequence item alignment before updating stack.
		isListDash := trimmed == "-" || strings.HasPrefix(trimmed, "- ")
		if len(seqStack) > 0 && !isListDash && bracketDepth == 0 && braceDepth == 0 {
			top := &seqStack[len(seqStack)-1]
			if top.itemKeyIndent > top.seqIndent && indent > top.seqIndent && indent < top.itemKeyIndent {
				return fmt.Errorf("line %d: misaligned YAML indentation (%d spaces, expected >= %d)", lineNum, indent, top.itemKeyIndent)
			}
		}

		if isListDash {
			if len(seqStack) > 0 && seqStack[len(seqStack)-1].dashIndent == -1 && indent > seqStack[len(seqStack)-1].seqIndent {
				seqStack[len(seqStack)-1].dashIndent = indent
			} else if len(seqStack) == 0 || indent > seqStack[len(seqStack)-1].dashIndent {
				seqStack = append(seqStack, seqFrame{seqIndent: indent, dashIndent: indent, itemKeyIndent: -1})
			} else if indent == seqStack[len(seqStack)-1].dashIndent {
				seqStack[len(seqStack)-1].itemKeyIndent = -1
			}
			afterDashRaw := strings.TrimPrefix(trimmed, "-")
			afterDash := strings.TrimSpace(afterDashRaw)
			if afterDash != "" {
				spacesAfterDash := len(afterDashRaw) - len(strings.TrimLeft(afterDashRaw, " "))
				if k, v, ok := splitYAMLKeyValue(afterDash); ok {
					seqStack[len(seqStack)-1].itemKeyIndent = indent + 1 + spacesAfterDash
					_ = k
					if isBlockScalarHeader(v) {
						inBlockScalar = true
						blockScalarIndent = indent + 1 + spacesAfterDash
					}
				}
			}
			continue
		}

		if k, v, ok := splitYAMLKeyValue(trimmed); ok {
			if len(seqStack) > 0 && indent > seqStack[len(seqStack)-1].dashIndent && seqStack[len(seqStack)-1].itemKeyIndent == -1 {
				seqStack[len(seqStack)-1].itemKeyIndent = indent
			}
			if v == "" {
				seqStack = append(seqStack, seqFrame{seqIndent: indent, dashIndent: -1, itemKeyIndent: -1})
			} else if isBlockScalarHeader(v) {
				inBlockScalar = true
				blockScalarIndent = indent
			}
			_ = k
			continue
		}

		// If not inside an open multi-line flow collection, a non-key/non-list line is invalid YAML mapping syntax.
		if bracketDepth == 0 && braceDepth == 0 {
			return fmt.Errorf("line %d: invalid YAML mapping or sequence syntax: %q", lineNum, trimmed)
		}
	}

	if bracketDepth != 0 {
		return fmt.Errorf("unclosed '[' in YAML workflow")
	}
	if braceDepth != 0 {
		return fmt.Errorf("unclosed '{' in YAML workflow")
	}
	return nil
}

// ValidateHealthGateWorkflow parses GitHub Actions workflow YAML using pure Go stdlib
// and verifies triggers on push/pull_request to main plus an active, uncommented step
// in job "full-kernel" invoking `mainline-health-gate` without `continue-on-error: true` or `if: false`.
func ValidateHealthGateWorkflow(content string) WorkflowValidationResult {
	var res WorkflowValidationResult

	if strings.Contains(content, "\t") {
		res.Errors = append(res.Errors, "workflow contains forbidden TAB characters in YAML")
		return res
	}

	rawLines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	// 1. Parse `on:` section for push/pull_request targeting `main`.
	inOn := false
	onIndent := -1
	activeEvent := ""
	eventIndent := -1
	inBranches := false
	branchesIndent := -1

	for _, raw := range rawLines {
		clean := stripYAMLComment(raw)
		trimmed := strings.TrimSpace(clean)
		if trimmed == "" {
			continue
		}
		indent := leadingSpaces(clean)

		if !inOn {
			if indent == 0 {
				if k, _, ok := splitYAMLKeyValue(trimmed); ok && k == "on" {
					inOn = true
					onIndent = indent
				}
			}
			continue
		}

		if indent <= onIndent {
			inOn = false
			inBranches = false
			activeEvent = ""
			continue
		}

		if activeEvent != "" && indent <= eventIndent {
			activeEvent = ""
			inBranches = false
		}
		if inBranches && indent <= branchesIndent {
			inBranches = false
		}

		if inBranches && strings.HasPrefix(trimmed, "- ") {
			branchVal := unquoteYAMLValue(strings.TrimPrefix(trimmed, "- "))
			if branchVal == "main" {
				if activeEvent == "push" {
					res.HasPushMain = true
				} else if activeEvent == "pull_request" {
					res.HasPRMain = true
				}
			}
			continue
		}

		if k, v, ok := splitYAMLKeyValue(trimmed); ok {
			if k == "push" || k == "pull_request" {
				activeEvent = k
				eventIndent = indent
				inBranches = false
				continue
			}
			if activeEvent != "" && k == "branches" {
				if inlineListHasMain(v) {
					if activeEvent == "push" {
						res.HasPushMain = true
					} else if activeEvent == "pull_request" {
						res.HasPRMain = true
					}
					inBranches = false
				} else {
					inBranches = true
					branchesIndent = indent
				}
				continue
			}
		}
	}

	// 2. Parse `jobs:` and `steps:` blocks with job ownership tracking.
	var steps []WorkflowStep
	jobsMap := make(map[string]workflowJob)

	inJobs := false
	jobsIndent := -1
	jobItemIndent := -1
	jobFieldIndent := -1
	var currentJob *workflowJob

	inSteps := false
	stepsIndent := -1
	stepItemIndent := -1
	stepKeyIndent := -1
	var currentStep *WorkflowStep
	inRunBlock := false
	runKeyIndent := -1

	flushStep := func() {
		if currentStep != nil {
			steps = append(steps, *currentStep)
			currentStep = nil
		}
		stepKeyIndent = -1
		inRunBlock = false
		runKeyIndent = -1
	}

	flushJob := func() {
		if currentJob != nil {
			jobsMap[currentJob.Name] = *currentJob
			currentJob = nil
		}
		jobFieldIndent = -1
	}

	applyStepField := func(step *WorkflowStep, k, v string, keyIndent int) {
		switch k {
		case "name":
			step.Name = unquoteYAMLValue(v)
		case "uses":
			step.Uses = unquoteYAMLValue(v)
		case "if":
			step.IfExpr = v
		case "continue-on-error":
			if isTruthyContinueOnError(v) {
				step.ContinueOnError = true
			}
		case "run":
			if isBlockScalarHeader(v) {
				inRunBlock = true
				runKeyIndent = keyIndent
			} else {
				cmd := unquoteYAMLValue(v)
				if cmd != "" {
					if step.Run != "" {
						step.Run += "\n"
					}
					step.Run += cmd
				}
			}
		}
	}

	for lineIdx, raw := range rawLines {
		// Handle multi-line `run: |` / `run: >` bodies before stripping YAML comments,
		// while respecting shell `#` comments on each body line.
		if inSteps && currentStep != nil && inRunBlock {
			trimmedRaw := strings.TrimSpace(raw)
			if trimmedRaw == "" {
				continue
			}
			indent := leadingSpaces(raw)
			if indent > runKeyIndent {
				cleanShell := strings.TrimSpace(stripYAMLComment(trimmedRaw))
				if cleanShell != "" {
					if currentStep.Run != "" {
						currentStep.Run += "\n"
					}
					currentStep.Run += cleanShell
				}
				continue
			}
			inRunBlock = false
			runKeyIndent = -1
		}

		clean := stripYAMLComment(raw)
		trimmed := strings.TrimSpace(clean)
		if trimmed == "" {
			continue
		}
		indent := leadingSpaces(clean)

		if inSteps {
			if indent < stepsIndent || (indent == stepsIndent && !strings.HasPrefix(trimmed, "-")) {
				flushStep()
				inSteps = false
				stepItemIndent = -1
			}
		}

		if inSteps {
			isListDash := trimmed == "-" || strings.HasPrefix(trimmed, "- ")
			if isListDash && (stepItemIndent == -1 || indent == stepItemIndent) {
				flushStep()
				stepItemIndent = indent
				jobName := ""
				if currentJob != nil {
					jobName = currentJob.Name
				}
				currentStep = &WorkflowStep{Job: jobName, Line: lineIdx + 1}

				afterDashRaw := strings.TrimPrefix(trimmed, "-")
				afterDash := strings.TrimSpace(afterDashRaw)
				if afterDash != "" {
					spacesAfterDash := len(afterDashRaw) - len(strings.TrimLeft(afterDashRaw, " "))
					stepKeyIndent = indent + 1 + spacesAfterDash
					if k, v, ok := splitYAMLKeyValue(afterDash); ok {
						applyStepField(currentStep, k, v, stepKeyIndent)
					}
				}
				continue
			}

			if currentStep != nil && indent > stepItemIndent {
				if stepKeyIndent == -1 {
					stepKeyIndent = indent
				}
				if indent == stepKeyIndent {
					if k, v, ok := splitYAMLKeyValue(trimmed); ok {
						applyStepField(currentStep, k, v, indent)
					}
				}
			}
			continue
		}

		if !inJobs {
			if indent == 0 {
				if k, v, ok := splitYAMLKeyValue(trimmed); ok && k == "jobs" && v == "" {
					inJobs = true
					jobsIndent = indent
					jobItemIndent = -1
				}
			}
			continue
		}

		if indent <= jobsIndent {
			flushJob()
			inJobs = false
			jobItemIndent = -1
			continue
		}

		if jobItemIndent == -1 || indent == jobItemIndent {
			if k, _, ok := splitYAMLKeyValue(trimmed); ok {
				flushJob()
				jobItemIndent = indent
				currentJob = &workflowJob{Name: unquoteYAMLValue(k), Line: lineIdx + 1}
			}
			continue
		}

		if currentJob != nil && indent > jobItemIndent {
			if jobFieldIndent == -1 {
				jobFieldIndent = indent
			}
			if indent == jobFieldIndent {
				if k, v, ok := splitYAMLKeyValue(trimmed); ok {
					switch k {
					case "steps":
						if v == "" {
							inSteps = true
							stepsIndent = indent
							stepItemIndent = -1
						}
					case "if":
						currentJob.IfExpr = v
					case "continue-on-error":
						if isTruthyContinueOnError(v) {
							currentJob.ContinueOnError = true
						}
					}
				}
			}
		}
	}
	flushStep()
	flushJob()

	// 3. Evaluate job "full-kernel" and its health-gate steps.
	fullKernelJob, hasFullKernel := jobsMap["full-kernel"]
	if !hasFullKernel {
		res.Errors = append(res.Errors, "missing required job \"full-kernel\"")
	} else {
		if fullKernelJob.ContinueOnError {
			res.HasContinueOnError = true
			res.Errors = append(res.Errors, fmt.Sprintf("line %d: job \"full-kernel\" has continue-on-error: true; must fail closed", fullKernelJob.Line))
		}
		if isDisabledIfExpr(fullKernelJob.IfExpr) {
			res.HasDisabledIf = true
			res.Errors = append(res.Errors, fmt.Sprintf("line %d: job \"full-kernel\" is disabled via if: %s", fullKernelJob.Line, strings.TrimSpace(fullKernelJob.IfExpr)))
		}
	}

	for _, step := range steps {
		if step.Job != "full-kernel" {
			continue
		}
		if !containsMainlineHealthGateInvocation(step.Run) {
			continue
		}
		res.HealthGateSteps = append(res.HealthGateSteps, step)
		if step.ContinueOnError {
			res.HasContinueOnError = true
			res.Errors = append(res.Errors, fmt.Sprintf("line %d: mainline-health-gate step has continue-on-error: true; must fail closed", step.Line))
		}
		if isDisabledIfExpr(step.IfExpr) {
			res.HasDisabledIf = true
			res.Errors = append(res.Errors, fmt.Sprintf("line %d: mainline-health-gate step is disabled via if: %s", step.Line, strings.TrimSpace(step.IfExpr)))
		}
	}

	res.HasHealthGateStep = len(res.HealthGateSteps) > 0
	if !res.HasHealthGateStep {
		res.Errors = append(res.Errors, "missing mandatory mainline-health-gate step in job \"full-kernel\"")
	}
	if !res.HasPushMain {
		res.Errors = append(res.Errors, "workflow missing push trigger targeting main branch")
	}
	if !res.HasPRMain {
		res.Errors = append(res.Errors, "workflow missing pull_request trigger targeting main branch")
	}

	res.Valid = len(res.Errors) == 0
	return res
}
