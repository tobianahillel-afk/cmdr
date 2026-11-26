package harness

import (
	"fmt"
	"strings"
)

// WorkflowStep represents a parsed step inside a GitHub Actions workflow job.
type WorkflowStep struct {
	Name            string
	Run             string
	Uses            string
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

func isDisabledIfExpr(expr string) bool {
	norm := strings.ToLower(strings.TrimSpace(unquoteYAMLValue(expr)))
	norm = strings.Join(strings.Fields(norm), " ")
	switch norm {
	case "false", "${{ false }}", "${{false}}", "!true", "${{ !true }}", "${{!true}}":
		return true
	default:
		return false
	}
}

func containsMainlineHealthGateInvocation(runScript string) bool {
	for _, rawLine := range strings.Split(runScript, "\n") {
		line := strings.TrimSpace(stripYAMLComment(rawLine))
		if line == "" {
			continue
		}
		if strings.Contains(line, "mainline-health-gate") {
			return true
		}
	}
	return false
}

// ValidateHealthGateWorkflow parses GitHub Actions workflow YAML using pure Go stdlib
// and verifies triggers on push/pull_request to main plus an active, uncommented step
// invoking `mainline-health-gate` without `continue-on-error: true` or `if: false`.
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

	// 2. Parse `steps:` blocks across all jobs.
	var steps []WorkflowStep
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

	applyStepField := func(step *WorkflowStep, k, v string, keyIndent int) {
		switch k {
		case "name":
			step.Name = unquoteYAMLValue(v)
		case "uses":
			step.Uses = unquoteYAMLValue(v)
		case "if":
			step.IfExpr = v
		case "continue-on-error":
			if strings.EqualFold(unquoteYAMLValue(v), "true") {
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

		if !inSteps {
			if k, v, ok := splitYAMLKeyValue(trimmed); ok && k == "steps" && v == "" {
				inSteps = true
				stepsIndent = indent
				stepItemIndent = -1
			}
			continue
		}

		isListDash := trimmed == "-" || strings.HasPrefix(trimmed, "- ")
		if isListDash && (stepItemIndent == -1 || indent == stepItemIndent) {
			flushStep()
			stepItemIndent = indent
			currentStep = &WorkflowStep{Line: lineIdx + 1}

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
	}
	flushStep()

	// 3. Evaluate health-gate steps.
	for _, step := range steps {
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
		res.Errors = append(res.Errors, "missing mandatory mainline-health-gate step")
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
