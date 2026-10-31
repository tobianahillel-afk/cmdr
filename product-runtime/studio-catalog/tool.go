// Package studiocatalog implements the genuine Go runtime for CMDR Studio Tool Definition
// and Functional Capability Contract conforming strictly to all 27 sections of CAP-STD-003,
// tenant-scope-contract.md, error-contract.md, permission-catalog.md, and REQ-PROD-006,
// REQ-PROD-016, REQ-OBJ-009, REQ-AI-002.
package studiocatalog

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// TenantRef represents a tenant identifier string under tenant-scope-contract.md.
type TenantRef = string

// -------------------------------------------------------------------------
// CAP-STD-003 §11: Domain Models, Enums, and Value Types
// -------------------------------------------------------------------------

// ExecutionNature defines the operational nature of a Tool under CAP-STD-003 §3, §11, §25.
// Tri-modal execution classification: deterministic, ai-capable, manual.
// An ai-capable tool NEVER morphs into an Agent (REQ-AI-002).
type ExecutionNature string

const (
	ExecutionNatureDeterministic ExecutionNature = "deterministic"
	ExecutionNatureAICapable     ExecutionNature = "ai-capable"
	ExecutionNatureManual        ExecutionNature = "manual"
)

// IsValid checks if an ExecutionNature is valid.
func (en ExecutionNature) IsValid() bool {
	switch en {
	case ExecutionNatureDeterministic, ExecutionNatureAICapable, ExecutionNatureManual:
		return true
	default:
		return false
	}
}

// ToolFunctionalState represents the 6 strict functional states under CAP-STD-003 §14.
// Critical rule: States 'published' and 'deployed' must NEVER be inferred or claimed in STD-1.
type ToolFunctionalState string

const (
	ToolStateDraft              ToolFunctionalState = "draft"
	ToolStateReviewReady        ToolFunctionalState = "review-ready"
	ToolStateAvailableReference ToolFunctionalState = "available-reference"
	ToolStateDeprecated         ToolFunctionalState = "deprecated"
	ToolStateDisabled           ToolFunctionalState = "disabled"
	ToolStateSuperseded         ToolFunctionalState = "superseded"
)

// IsValid checks if a ToolFunctionalState is one of the 6 canonical states.
func (tfs ToolFunctionalState) IsValid() bool {
	switch tfs {
	case ToolStateDraft, ToolStateReviewReady, ToolStateAvailableReference,
		ToolStateDeprecated, ToolStateDisabled, ToolStateSuperseded:
		return true
	default:
		return false
	}
}

// ActionClass represents the action classification (Class 0, 1, 2, 3, 4) under action-classification.md.
type ActionClass int

const (
	ActionClass0Observation ActionClass = 0 // Read-only discovery and deep inspection
	ActionClass1Assessment  ActionClass = 1 // Deterministic no-effect contract validation
	ActionClass2Preparation ActionClass = 2 // Reversible draft authoring and mutation
	ActionClass3Effect      ActionClass = 3 // External authority required
	ActionClass4Destructive ActionClass = 4 // Heightened governance required
)

// SideEffectContext defines whether a tool causes external side-effects under CAP-STD-003 §11.
type SideEffectContext struct {
	HasSideEffects bool   `json:"has_side_effects"`
	MutationScope  string `json:"mutation_scope,omitempty"` // "read-only", "local-session", "remote-target"
	IsReversible   bool   `json:"is_reversible"`
	Description    string `json:"description,omitempty"`
}

// RuntimeAvailabilityStatus represents the decoupled availability projection under CAP-STD-003 §8, §22, §25.
// Invariant: Tool definition validity is decoupled from runtime daemon/provider availability.
type RuntimeAvailabilityStatus struct {
	IsAvailable  bool      `json:"is_available"`
	HealthStatus string    `json:"health_status"` // "available", "unavailable", "partial", "stale"
	CheckedAt    time.Time `json:"checked_at"`
	Message      string    `json:"message,omitempty"`
}

// ProviderReference defines an external integration, model provider, or runtime host dependency under CAP-STD-003 §8, §18.
type ProviderReference struct {
	ProviderID   string                    `json:"provider_id"`
	ProviderType string                    `json:"provider_type"` // e.g. "model-provider", "integration", "runtime-host"
	Freshness    string                    `json:"freshness,omitempty"`
	Availability RuntimeAvailabilityStatus `json:"availability"`
}

// ToolOperation represents a discrete, bounded operation supported by the Tool under CAP-STD-003 §4, §11.
type ToolOperation struct {
	OperationID       string            `json:"operation_id"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Parameters        []Parameter       `json:"parameters"`
	OutputParameters  []Parameter       `json:"output_parameters,omitempty"`
	IsIdempotent      bool              `json:"is_idempotent"`
	SideEffectContext SideEffectContext `json:"side_effect_context"`
}

// ToolParameter aliases Parameter from model.go for consistent API ergonomics.
type ToolParameter = Parameter

// ToolDependency aliases Dependency from model.go.
type ToolDependency = Dependency

// ToolConsumer aliases Consumer from model.go.
type ToolConsumer = Consumer

// ToolDefinition represents the complete canonical Tool contract under CAP-STD-003.
type ToolDefinition struct {
	ID                  string                    `json:"id"`
	TenantID            string                    `json:"tenant_id"`
	EnvironmentID       string                    `json:"environment_id,omitempty"`
	Name                string                    `json:"name"`
	DisplayName         string                    `json:"display_name"`
	Description         string                    `json:"description"`
	Purpose             string                    `json:"purpose"`
	Owner               Owner                     `json:"owner"`
	Version             VersionInfo               `json:"version"`
	FunctionalState     ToolFunctionalState       `json:"functional_state"`
	ExecutionNature     ExecutionNature           `json:"execution_nature"`
	SupportedOperations []string                  `json:"supported_operations"`
	Operations          []ToolOperation           `json:"operations"`
	Parameters          []Parameter               `json:"parameters"`
	OutputParameters    []Parameter               `json:"output_parameters,omitempty"`
	Dependencies        []Dependency              `json:"dependencies"`
	ProviderRefs        []ProviderReference       `json:"provider_refs"`
	Consumers           []Consumer                `json:"consumers"`
	ActionClass         ActionClass               `json:"action_class"`
	RiskClassification  RiskClassification        `json:"risk_classification"`
	SideEffects         SideEffectContext         `json:"side_effects"`
	Limitations         []string                  `json:"limitations"`
	RuntimeAvailability RuntimeAvailabilityStatus `json:"runtime_availability"`
	Provenance          Provenance                `json:"provenance"`
	AuditTrail          []AuditRecord             `json:"audit_trail,omitempty"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

// Clone creates an immutable deep copy of a ToolDefinition.
func (t *ToolDefinition) Clone() *ToolDefinition {
	if t == nil {
		return nil
	}
	cp := *t
	cp.Owner.Normalize()

	if t.SupportedOperations != nil {
		cp.SupportedOperations = make([]string, len(t.SupportedOperations))
		copy(cp.SupportedOperations, t.SupportedOperations)
	}
	if t.Operations != nil {
		cp.Operations = make([]ToolOperation, len(t.Operations))
		for i, op := range t.Operations {
			opCp := op
			if op.Parameters != nil {
				opCp.Parameters = make([]Parameter, len(op.Parameters))
				copy(opCp.Parameters, op.Parameters)
			}
			if op.OutputParameters != nil {
				opCp.OutputParameters = make([]Parameter, len(op.OutputParameters))
				copy(opCp.OutputParameters, op.OutputParameters)
			}
			cp.Operations[i] = opCp
		}
	}
	if t.Parameters != nil {
		cp.Parameters = make([]Parameter, len(t.Parameters))
		copy(cp.Parameters, t.Parameters)
	}
	if t.OutputParameters != nil {
		cp.OutputParameters = make([]Parameter, len(t.OutputParameters))
		copy(cp.OutputParameters, t.OutputParameters)
	}
	if t.Dependencies != nil {
		cp.Dependencies = make([]Dependency, len(t.Dependencies))
		copy(cp.Dependencies, t.Dependencies)
	}
	if t.ProviderRefs != nil {
		cp.ProviderRefs = make([]ProviderReference, len(t.ProviderRefs))
		copy(cp.ProviderRefs, t.ProviderRefs)
	}
	if t.Consumers != nil {
		cp.Consumers = make([]Consumer, len(t.Consumers))
		copy(cp.Consumers, t.Consumers)
	}
	if t.Limitations != nil {
		cp.Limitations = make([]string, len(t.Limitations))
		copy(cp.Limitations, t.Limitations)
	}
	if t.Provenance.ProvenanceChain != nil {
		cp.Provenance.ProvenanceChain = make([]string, len(t.Provenance.ProvenanceChain))
		copy(cp.Provenance.ProvenanceChain, t.Provenance.ProvenanceChain)
	}
	if t.AuditTrail != nil {
		cp.AuditTrail = make([]AuditRecord, len(t.AuditTrail))
		copy(cp.AuditTrail, t.AuditTrail)
	}
	return &cp
}

// ToAsset projects a ToolDefinition into a canonical unified catalog Asset.
func (t *ToolDefinition) ToAsset() *Asset {
	if t == nil {
		return nil
	}
	owner := t.Owner
	owner.Normalize()

	tags := []string{"tool", string(t.ExecutionNature), string(t.FunctionalState)}
	for _, op := range t.SupportedOperations {
		tags = append(tags, "op:"+op)
	}

	return &Asset{
		ID:                 t.ID,
		TenantID:           t.TenantID,
		EnvironmentID:      t.EnvironmentID,
		Type:               AssetKindTool,
		Name:               t.Name,
		DisplayName:        t.DisplayName,
		Description:        t.Description,
		Tags:               tags,
		Owner:              owner,
		Version:            t.Version,
		LifecycleState:     LifecycleState(t.FunctionalState),
		FunctionalState:    FunctionalState(t.FunctionalState),
		DeliveryMode:       DeliveryModePlanned,
		ExecutionClass:     fmt.Sprintf("class-%d", t.ActionClass),
		RiskClassification: t.RiskClassification,
		IsAIPowered:        t.ExecutionNature == ExecutionNatureAICapable,
		HasNonAIFallback:   t.ExecutionNature == ExecutionNatureDeterministic,
		Parameters:         t.Parameters,
		Dependencies:       t.Dependencies,
		Consumers:          t.Consumers,
		Provenance:         t.Provenance,
		AuditTrail:         t.AuditTrail,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
}

// -------------------------------------------------------------------------
// Additional Error Codes for Tool Contract (error-contract.md)
// -------------------------------------------------------------------------

const (
	ErrorCodeInvalidContract        ErrorCode = "invalid-contract"
	ErrorCodeUnsupportedOperation    ErrorCode = "unsupported-operation"
	ErrorCodeIncompatibleDependency  ErrorCode = "incompatible-dependency"
	ErrorCodeUnavailableRuntime      ErrorCode = "unavailable-runtime"
)

// -------------------------------------------------------------------------
// Request / Response Structs for Tool Operations
// -------------------------------------------------------------------------

// InspectToolRequest parameters for deep tool contract inspection (Class 0).
type InspectToolRequest struct {
	TenantRef string `json:"tenant_ref"`
	ToolID    string `json:"tool_id"`
}

// InspectToolResponse represents the result of Class 0 inspection under CAP-STD-003 §12, §16.
type InspectToolResponse struct {
	Tool                *ToolDefinition           `json:"tool"`
	RuntimeAvailability RuntimeAvailabilityStatus `json:"runtime_availability"`
	IsExecutable        bool                      `json:"is_executable"` // Strictly false under REQ-PROD-006
	CorrelationID       string                    `json:"correlation_id"`
	InspectedAt         time.Time                 `json:"inspected_at"`
}

// ValidateToolContractRequest parameters for deterministic contract validation (Class 1).
type ValidateToolContractRequest struct {
	TenantRef TenantRef       `json:"tenant_ref"`
	ToolID    string          `json:"tool_id,omitempty"`
	Tool      *ToolDefinition `json:"tool,omitempty"`
}

// ValidationResult represents the output of Class 1 deterministic validation under CAP-STD-003 §12, §13, §25.
type ValidationResult struct {
	Valid                   bool              `json:"valid"`
	Issues                  []ValidationIssue `json:"issues"`
	CompletenessScore       float64           `json:"completeness_score"` // 0.0 to 1.0
	CheckedRules            []string          `json:"checked_rules"`
	IsAICapable             bool              `json:"is_ai_capable"`
	NonAgentConfirmed       bool              `json:"non_agent_confirmed"`
	DeterministicValidated  bool              `json:"deterministic_validated"`
	DecoupledRuntimeStatus  string            `json:"decoupled_runtime_status"`
	CorrelationID           string            `json:"correlation_id"`
	ValidatedAt             time.Time         `json:"validated_at"`
}

// ValidateToolContractResponse is an alias for ValidationResult for consistent naming.
type ValidateToolContractResponse = ValidationResult

// CreateDraftToolRequest parameters for creating a new draft tool definition (Class 2).
type CreateDraftToolRequest struct {
	TenantRef           TenantRef           `json:"tenant_ref"`
	EnvironmentID       string              `json:"environment_id,omitempty"`
	ID                  string              `json:"id,omitempty"`
	Name                string              `json:"name"`
	DisplayName         string              `json:"display_name,omitempty"`
	Description         string              `json:"description,omitempty"`
	Purpose             string              `json:"purpose"`
	Owner               Owner               `json:"owner"`
	ExecutionNature     ExecutionNature     `json:"execution_nature"`
	SupportedOperations []string            `json:"supported_operations"`
	Operations          []ToolOperation     `json:"operations,omitempty"`
	Parameters          []Parameter         `json:"parameters,omitempty"`
	OutputParameters    []Parameter         `json:"output_parameters,omitempty"`
	Dependencies        []Dependency        `json:"dependencies,omitempty"`
	ProviderRefs        []ProviderReference `json:"provider_refs,omitempty"`
	Consumers           []Consumer          `json:"consumers,omitempty"`
	ActionClass         ActionClass         `json:"action_class"`
	RiskClassification  RiskClassification  `json:"risk_classification"`
	SideEffects         SideEffectContext   `json:"side_effects"`
	Limitations         []string            `json:"limitations,omitempty"`
	Reason              string              `json:"reason,omitempty"`
}

// UpdateDraftToolRequest parameters for updating an existing draft tool definition (Class 2).
type UpdateDraftToolRequest struct {
	TenantRef           TenantRef           `json:"tenant_ref"`
	ToolID              string              `json:"tool_id"`
	Name                string              `json:"name,omitempty"`
	DisplayName         string              `json:"display_name,omitempty"`
	Description         string              `json:"description,omitempty"`
	Purpose             string              `json:"purpose,omitempty"`
	Owner               *Owner              `json:"owner,omitempty"`
	ExecutionNature     ExecutionNature     `json:"execution_nature,omitempty"`
	FunctionalState     ToolFunctionalState `json:"functional_state,omitempty"`
	SupportedOperations []string            `json:"supported_operations,omitempty"`
	Operations          []ToolOperation     `json:"operations,omitempty"`
	Parameters          []Parameter         `json:"parameters,omitempty"`
	OutputParameters    []Parameter         `json:"output_parameters,omitempty"`
	Dependencies        []Dependency        `json:"dependencies,omitempty"`
	ProviderRefs        []ProviderReference `json:"provider_refs,omitempty"`
	Consumers           []Consumer          `json:"consumers,omitempty"`
	ActionClass         *ActionClass        `json:"action_class,omitempty"`
	RiskClassification  RiskClassification  `json:"risk_classification,omitempty"`
	SideEffects         *SideEffectContext  `json:"side_effects,omitempty"`
	Limitations         []string            `json:"limitations,omitempty"`
	Reason              string              `json:"reason,omitempty"`
}

// -------------------------------------------------------------------------
// Tool Store Interface & In-Memory Implementation
// -------------------------------------------------------------------------

// ToolStore defines thread-safe storage operations for ToolDefinitions.
type ToolStore interface {
	Get(ctx context.Context, tenantID, toolID string) (*ToolDefinition, error)
	Put(ctx context.Context, tool *ToolDefinition) error
	Delete(ctx context.Context, tenantID, toolID string) error
	List(ctx context.Context, tenantID string) ([]*ToolDefinition, error)
}

// InMemoryToolStore implements thread-safe in-memory storage for ToolDefinitions.
type InMemoryToolStore struct {
	mu    sync.RWMutex
	tools map[string]*ToolDefinition // key: tenantID + ":" + toolID
}

// NewInMemoryToolStore initializes an empty InMemoryToolStore.
func NewInMemoryToolStore() *InMemoryToolStore {
	return &InMemoryToolStore{
		tools: make(map[string]*ToolDefinition),
	}
}

func toolKey(tenantID, toolID string) string {
	return tenantID + ":" + toolID
}

// Get retrieves a tool by tenant and tool ID.
func (s *InMemoryToolStore) Get(ctx context.Context, tenantID, toolID string) (*ToolDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tools[toolKey(tenantID, toolID)]
	if !ok {
		return nil, nil
	}
	return t.Clone(), nil
}

// Put saves or updates a tool definition.
func (s *InMemoryToolStore) Put(ctx context.Context, tool *ToolDefinition) error {
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[toolKey(tool.TenantID, tool.ID)] = tool.Clone()
	return nil
}

// Delete removes a tool by tenant and tool ID.
func (s *InMemoryToolStore) Delete(ctx context.Context, tenantID, toolID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tools, toolKey(tenantID, toolID))
	return nil
}

// List returns all tools for a tenant.
func (s *InMemoryToolStore) List(ctx context.Context, tenantID string) ([]*ToolDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []*ToolDefinition
	for _, t := range s.tools {
		if t.TenantID == tenantID {
			list = append(list, t.Clone())
		}
	}
	return list, nil
}

// -------------------------------------------------------------------------
// ToolService Interface
// -------------------------------------------------------------------------

// ToolService defines the canonical functional contract service for CMDR Studio Tools under CAP-STD-003.
type ToolService interface {
	// Class 0: Inspect Tool Contract
	InspectTool(ctx context.Context, auth *AuthContext, tenantRef TenantRef, toolID string) (*InspectToolResponse, error)

	// Class 1: Deterministic Contract Validation (REQ-PROD-016)
	ValidateToolContract(ctx context.Context, auth *AuthContext, req *ValidateToolContractRequest) (*ValidationResult, error)

	// Class 2: Draft Creation and Mutation
	CreateDraftTool(ctx context.Context, auth *AuthContext, req *CreateDraftToolRequest) (*ToolDefinition, error)
	UpdateDraftTool(ctx context.Context, auth *AuthContext, req *UpdateDraftToolRequest) (*ToolDefinition, error)

	// Operational Decoupled Runtime Health Updates
	SetRuntimeAvailability(ctx context.Context, tenantRef TenantRef, toolID string, status RuntimeAvailabilityStatus) error

	// REQ-PROD-006: Permanent Execution Blocking (all return 400 execution-blocked)
	Invoke(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Execute(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Run(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Call(ctx context.Context, auth *AuthContext, toolID string, payload any) error
	Trigger(ctx context.Context, auth *AuthContext, toolID string, payload any) error
}

// DefaultToolService implements ToolService.
type DefaultToolService struct {
	store ToolStore
}

// NewDefaultToolService constructs a new DefaultToolService with the provided store.
func NewDefaultToolService(store ToolStore) *DefaultToolService {
	if store == nil {
		store = NewInMemoryToolStore()
	}
	return &DefaultToolService{
		store: store,
	}
}

// ensureAuth safely extracts an AuthContext even if nil was passed.
func ensureAuth(auth *AuthContext) AuthContext {
	if auth == nil {
		return AuthContext{
			PrincipalID: "anonymous",
			Permissions: make(map[string]bool),
		}
	}
	return *auth
}

// -------------------------------------------------------------------------
// Class 0: Inspect Tool Implementation (CAP-STD-003 §12, §16)
// -------------------------------------------------------------------------

// InspectTool performs deep inspection of the Tool functional contract.
// - Enforces tenant isolation and anti-oracle uniform 404 (never leaks existence).
// - Enforces perm.studio.library.read.
// - Enforces REQ-OBJ-009 owner unknown normalization.
// - Decouples runtime availability: offline provider still returns 200 with complete valid definition.
func (s *DefaultToolService) InspectTool(ctx context.Context, auth *AuthContext, tenantRef TenantRef, toolID string) (*InspectToolResponse, error) {
	corrID := GenerateCorrelationID()
	ac := ensureAuth(auth)

	// 1. Validate tenant scoping
	if errEnv := ValidateTenantScope(ac, tenantRef, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Validate read permission (governed strictly by perm.studio.library.read)
	if !ac.CanReadAssetKind(AssetKindTool) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Validate tool ID parameter
	trimmedToolID := strings.TrimSpace(toolID)
	if trimmedToolID == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Tool ID is required for inspection", corrID, 400, false)
	}

	// 4. Retrieve tool from store
	tool, err := s.store.Get(ctx, tenantRef, trimmedToolID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to retrieve tool: %v", err), corrID, 500, true)
	}

	// 5. Anti-oracle uniform 404: if not found in this tenant, return 404
	if tool == nil {
		return nil, NewAssetNotFoundError(trimmedToolID, tenantRef, corrID)
	}

	// 6. REQ-OBJ-009: Normalize owner
	tool.Owner.Normalize()

	// 7. Assemble Class 0 inspection response
	resp := &InspectToolResponse{
		Tool:                tool,
		RuntimeAvailability: tool.RuntimeAvailability,
		IsExecutable:        false, // REQ-PROD-006: definition inspection never confers execution rights
		CorrelationID:       corrID,
		InspectedAt:         time.Now().UTC(),
	}

	return resp, nil
}

// -------------------------------------------------------------------------
// Class 1: Validate Tool Contract Implementation (CAP-STD-003 §12, §13, §25)
// -------------------------------------------------------------------------

// ValidateToolContract evaluates a Tool definition against the 100% deterministic local rule checklist.
// Operates without AI models, network, or external daemons (REQ-PROD-016, REQ-AI-001).
// Strictly enforces anti-agent boundary defense (REQ-AI-002).
func (s *DefaultToolService) ValidateToolContract(ctx context.Context, auth *AuthContext, req *ValidateToolContractRequest) (*ValidationResult, error) {
	corrID := GenerateCorrelationID()
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewError(ErrorCodeInvalidContract, "Validation request cannot be nil", corrID, 400, false)
	}

	// 1. Validate tenant scoping
	if errEnv := ValidateTenantScope(ac, req.TenantRef, corrID); errEnv != nil {
		return nil, errEnv
	}

	// 2. Validate read permission
	if !ac.CanReadAssetKind(AssetKindTool) {
		return nil, NewPermissionDeniedError(PermStudioLibraryRead, corrID)
	}

	// 3. Resolve target tool definition
	var target *ToolDefinition
	if req.Tool != nil {
		target = req.Tool.Clone()
		if target.TenantID == "" {
			target.TenantID = req.TenantRef
		}
	} else if strings.TrimSpace(req.ToolID) != "" {
		existing, err := s.store.Get(ctx, req.TenantRef, strings.TrimSpace(req.ToolID))
		if err != nil {
			return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to lookup tool: %v", err), corrID, 500, true)
		}
		if existing == nil {
			return nil, NewAssetNotFoundError(req.ToolID, req.TenantRef, corrID)
		}
		target = existing
	} else {
		return nil, NewError(ErrorCodeInvalidContract, "Either Tool object or ToolID must be provided for contract validation", corrID, 400, false)
	}

	// 4. Execute deterministic validation checklist
	result := DeterministicValidateToolContract(target, corrID)
	return result, nil
}

// DeterministicValidateToolContract executes the 100% local rule checklist without any AI models (REQ-PROD-016).
func DeterministicValidateToolContract(tool *ToolDefinition, correlationID string) *ValidationResult {
	if correlationID == "" {
		correlationID = GenerateCorrelationID()
	}

	var issues []ValidationIssue
	checkedRules := []string{
		"CHK-TOOL-IDENTITY",
		"CHK-EXECUTION-NATURE",
		"CHK-OPERATIONS-SCOPE",
		"CHK-PARAMETERS-TYPING",
		"CHK-ANTI-AGENT-BOUNDARY",
		"CHK-OWNER-NORMALIZATION",
		"CHK-ACTION-RISK-CLASS",
		"CHK-SIDE-EFFECTS-PROFILE",
		"CHK-DEPENDENCY-INTEGRITY",
		"CHK-DECOUPLED-RUNTIME",
	}

	totalChecks := 10
	passedChecks := 0

	// 1. CHK-TOOL-IDENTITY
	identityPass := true
	if strings.TrimSpace(tool.ID) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "id",
			Code:     "missing-tool-id",
			Severity: "error",
			Message:  "Tool ID is mandatory and cannot be blank",
		})
		identityPass = false
	}
	if strings.TrimSpace(tool.Name) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "name",
			Code:     "missing-tool-name",
			Severity: "error",
			Message:  "Tool name is mandatory and cannot be blank",
		})
		identityPass = false
	}
	if strings.TrimSpace(tool.Purpose) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "purpose",
			Code:     "missing-tool-purpose",
			Severity: "error",
			Message:  "Tool functional purpose statement is mandatory under CAP-STD-003 §3",
		})
		identityPass = false
	}
	if identityPass {
		passedChecks++
	}

	// 2. CHK-EXECUTION-NATURE
	if tool.ExecutionNature.IsValid() {
		passedChecks++
	} else {
		issues = append(issues, ValidationIssue{
			Field:    "execution_nature",
			Code:     "invalid-execution-nature",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid execution nature %q; must be 'deterministic', 'ai-capable', or 'manual'", tool.ExecutionNature),
		})
	}

	// 3. CHK-OPERATIONS-SCOPE
	hasOperations := len(tool.SupportedOperations) > 0 || len(tool.Operations) > 0
	if hasOperations {
		passedChecks++
	} else {
		issues = append(issues, ValidationIssue{
			Field:    "supported_operations",
			Code:     "unsupported-operation",
			Severity: "error",
			Message:  "Tool must declare at least one supported operation",
		})
	}

	// 4. CHK-PARAMETERS-TYPING
	paramPass := true
	for i, param := range tool.Parameters {
		if strings.TrimSpace(param.Name) == "" {
			issues = append(issues, ValidationIssue{
				Field:    fmt.Sprintf("parameters[%d].name", i),
				Code:     "missing-parameter-name",
				Severity: "error",
				Message:  "Parameter name cannot be blank",
			})
			paramPass = false
		}
		if strings.TrimSpace(param.Type) == "" {
			issues = append(issues, ValidationIssue{
				Field:    fmt.Sprintf("parameters[%d].type", i),
				Code:     "missing-parameter-type",
				Severity: "error",
				Message:  "Parameter type cannot be blank",
			})
			paramPass = false
		}
	}
	for i, op := range tool.Operations {
		for j, param := range op.Parameters {
			if strings.TrimSpace(param.Name) == "" || strings.TrimSpace(param.Type) == "" {
				issues = append(issues, ValidationIssue{
					Field:    fmt.Sprintf("operations[%d].parameters[%d]", i, j),
					Code:     "invalid-operation-parameter",
					Severity: "error",
					Message:  "Operation parameter must declare name and type",
				})
				paramPass = false
			}
		}
	}
	if paramPass {
		passedChecks++
	}

	// 5. CHK-ANTI-AGENT-BOUNDARY (REQ-AI-002, CAP-STD-003 §1, §25)
	// AI-capable tools must NEVER morph into an Agent (no planning loops, prompt goals, agent team roles)
	nonAgentConfirmed := true
	forbiddenAgentKeywords := []string{
		"planning_loop", "autonomous_loop", "agent_team", "goal_prompt",
		"agent_role", "loop_autonomy", "subagent_spawn", "self_direct",
	}

	combinedText := strings.ToLower(fmt.Sprintf("%s %s %s %v", tool.Name, tool.Description, tool.Purpose, tool.Limitations))
	for _, kw := range forbiddenAgentKeywords {
		if strings.Contains(combinedText, kw) {
			issues = append(issues, ValidationIssue{
				Field:    "anti_agent_boundary",
				Code:     "anti-agent-boundary-violation",
				Severity: "error",
				Message:  fmt.Sprintf("Tool definition contains agent-exclusive keyword %q. AI-capable tools must remain bounded capability contracts and cannot assume Agent roles (REQ-AI-002).", kw),
			})
			nonAgentConfirmed = false
			break
		}
	}
	if nonAgentConfirmed {
		passedChecks++
	}

	// 6. CHK-OWNER-NORMALIZATION (REQ-OBJ-009)
	// Missing or unresolvable owner normalizes to "owner unknown", is_unknown: true
	if tool.Owner.IsUnknown || tool.Owner.OwnerName == "owner unknown" || strings.TrimSpace(tool.Owner.OwnerID) == "" {
		issues = append(issues, ValidationIssue{
			Field:    "owner",
			Code:     "owner-normalized",
			Severity: "info",
			Message:  "Owner normalized to 'owner unknown' under REQ-OBJ-009",
		})
	}
	passedChecks++

	// 7. CHK-ACTION-RISK-CLASS
	validActionClass := tool.ActionClass >= ActionClass0Observation && tool.ActionClass <= ActionClass4Destructive
	if validActionClass {
		passedChecks++
	} else {
		issues = append(issues, ValidationIssue{
			Field:    "action_class",
			Code:     "invalid-action-class",
			Severity: "error",
			Message:  fmt.Sprintf("Invalid ActionClass %d; must be 0 (Observation), 1 (Assessment), 2 (Preparation), 3 (Effect), or 4 (Destructive)", tool.ActionClass),
		})
	}

	// 8. CHK-SIDE-EFFECTS-PROFILE
	sideEffectsPass := true
	if tool.RiskClassification == RiskClassificationHigh || tool.RiskClassification == RiskClassificationCritical {
		if !tool.SideEffects.HasSideEffects && strings.TrimSpace(tool.SideEffects.Description) == "" {
			issues = append(issues, ValidationIssue{
				Field:    "side_effects",
				Code:     "undeclared-high-risk-side-effects",
				Severity: "warning",
				Message:  "High or critical risk tools should explicitly declare side effect context and reversibility",
			})
			sideEffectsPass = false
		}
	}
	if sideEffectsPass {
		passedChecks++
	}

	// 9. CHK-DEPENDENCY-INTEGRITY
	depPass := true
	for i, dep := range tool.Dependencies {
		if strings.TrimSpace(dep.RefID) == "" {
			issues = append(issues, ValidationIssue{
				Field:    fmt.Sprintf("dependencies[%d].ref_id", i),
				Code:     "missing-dependency-ref",
				Severity: "error",
				Message:  "Dependency ref_id cannot be blank",
			})
			depPass = false
		}
		if strings.TrimSpace(dep.Kind) == "" {
			issues = append(issues, ValidationIssue{
				Field:    fmt.Sprintf("dependencies[%d].kind", i),
				Code:     "missing-dependency-kind",
				Severity: "error",
				Message:  "Dependency kind cannot be blank",
			})
			depPass = false
		}
	}
	if depPass {
		passedChecks++
	}

	// 10. CHK-DECOUPLED-RUNTIME (CAP-STD-003 §8, §25 Given 1)
	// Tool definition validity is strictly decoupled from runtime availability!
	// An offline runtime daemon reports availability separately but does NOT cause contract validation to fail.
	decoupledStatus := "available"
	if !tool.RuntimeAvailability.IsAvailable || tool.RuntimeAvailability.HealthStatus == "unavailable" {
		decoupledStatus = "unavailable"
		issues = append(issues, ValidationIssue{
			Field:    "runtime_availability",
			Code:     "runtime-provider-offline",
			Severity: "info",
			Message:  "Runtime daemon is offline; tool contract specification remains valid under decoupled availability model (CAP-STD-003 §25)",
		})
	}
	passedChecks++

	// Determine overall validity: valid if there are zero errors (warnings and infos are acceptable)
	isValid := true
	for _, issue := range issues {
		if issue.Severity == "error" {
			isValid = false
			break
		}
	}

	score := float64(passedChecks) / float64(totalChecks)

	return &ValidationResult{
		Valid:                  isValid,
		Issues:                 issues,
		CompletenessScore:      score,
		CheckedRules:           checkedRules,
		IsAICapable:            tool.ExecutionNature == ExecutionNatureAICapable,
		NonAgentConfirmed:      nonAgentConfirmed,
		DeterministicValidated: true,
		DecoupledRuntimeStatus: decoupledStatus,
		CorrelationID:          correlationID,
		ValidatedAt:            time.Now().UTC(),
	}
}

// -------------------------------------------------------------------------
// Class 2: Create Draft Tool Implementation (CAP-STD-003 §10, §12)
// -------------------------------------------------------------------------

// CreateDraftTool creates a new mutable Tool functional definition in state 'draft'.
// - Governed by perm.studio.library.read or perm.studio.builder.read.
// - Rejects cross-tenant mutation with 403 cross-tenant-mutation-forbidden.
// - Normalizes missing/blank owner to 'owner unknown' (REQ-OBJ-009).
// - Enforces initial functional state strictly as 'draft'.
func (s *DefaultToolService) CreateDraftTool(ctx context.Context, auth *AuthContext, req *CreateDraftToolRequest) (*ToolDefinition, error) {
	corrID := GenerateCorrelationID()
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewError(ErrorCodeInvalidMetadata, "CreateDraftTool request cannot be nil", corrID, 400, false)
	}

	// 1. Validate tenant scoping
	if errEnv := ValidateTenantScope(ac, req.TenantRef, corrID); errEnv != nil {
		return nil, errEnv
	}

	// Cross-tenant mutation is strictly forbidden even in MSSP mode
	if ac.TenantID != "" && req.TenantRef != ac.TenantID {
		return nil, NewError(ErrorCodeCrossTenantMutationForbidden,
			fmt.Sprintf("Cross-tenant draft creation from %q to %q is strictly forbidden", ac.TenantID, req.TenantRef),
			corrID, 403, false)
	}

	// 2. Validate authoring permission
	if !ac.CanManageAssetKind(AssetKindTool) {
		return nil, NewPermissionDeniedError(PermStudioArtifactManage, corrID)
	}

	// 3. Validate mandatory fields
	toolID := strings.TrimSpace(req.ID)
	if toolID == "" {
		toolID = fmt.Sprintf("tool-%s", GenerateCorrelationID()[5:17])
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, NewInvalidMetadataError("Tool name is required", corrID)
	}

	purpose := strings.TrimSpace(req.Purpose)
	if purpose == "" {
		return nil, NewInvalidMetadataError("Tool purpose statement is required under CAP-STD-003 §3", corrID)
	}

	execNature := req.ExecutionNature
	if execNature == "" {
		execNature = ExecutionNatureDeterministic
	}
	if !execNature.IsValid() {
		return nil, NewInvalidMetadataError(fmt.Sprintf("Invalid execution nature %q", execNature), corrID)
	}

	// 4. Normalize Owner (REQ-OBJ-009)
	owner := req.Owner
	owner.Normalize()

	displayName := strings.TrimSpace(req.DisplayName)
	if displayName == "" {
		displayName = name
	}

	now := time.Now().UTC()
	actor := ac.PrincipalID
	if actor == "" {
		actor = "automation-designer"
	}

	reason := req.Reason
	if reason == "" {
		reason = "Initial draft creation"
	}

	audit := []AuditRecord{
		{
			Timestamp:   now,
			Actor:       actor,
			Reason:      reason,
			Version:     "0.1.0",
			TenantID:    req.TenantRef,
			Action:      "create_draft",
			NewMetadata: fmt.Sprintf("Created tool draft %s (%s)", name, toolID),
		},
	}

	supportedOps := req.SupportedOperations
	if len(supportedOps) == 0 && len(req.Operations) > 0 {
		for _, op := range req.Operations {
			supportedOps = append(supportedOps, op.Name)
		}
	}

	tool := &ToolDefinition{
		ID:                  toolID,
		TenantID:            req.TenantRef,
		EnvironmentID:       req.EnvironmentID,
		Name:                name,
		DisplayName:         displayName,
		Description:         req.Description,
		Purpose:             purpose,
		Owner:               owner,
		Version:             VersionInfo{VersionID: "v0.1.0", SemVer: "0.1.0", CreatedAt: now, IsImmutable: false},
		FunctionalState:     ToolStateDraft, // Strictly 'draft'
		ExecutionNature:     execNature,
		SupportedOperations: supportedOps,
		Operations:          req.Operations,
		Parameters:          req.Parameters,
		OutputParameters:    req.OutputParameters,
		Dependencies:        req.Dependencies,
		ProviderRefs:        req.ProviderRefs,
		Consumers:           req.Consumers,
		ActionClass:         req.ActionClass,
		RiskClassification:  req.RiskClassification,
		SideEffects:         req.SideEffects,
		Limitations:         req.Limitations,
		RuntimeAvailability: RuntimeAvailabilityStatus{
			IsAvailable:  true,
			HealthStatus: "available",
			CheckedAt:    now,
			Message:      "Draft tool registered",
		},
		Provenance: Provenance{
			CreatedBy:       actor,
			CreatedAt:       now,
			UpdatedBy:       actor,
			UpdatedAt:       now,
			ProvenanceChain: []string{fmt.Sprintf("draft:%s:%s", req.TenantRef, toolID)},
		},
		AuditTrail: audit,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 5. Persist
	if err := s.store.Put(ctx, tool); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to store tool: %v", err), corrID, 500, true)
	}

	return tool.Clone(), nil
}

// -------------------------------------------------------------------------
// Class 2: Update Draft Tool Implementation (CAP-STD-003 §10, §12)
// -------------------------------------------------------------------------

// UpdateDraftTool updates an existing draft Tool functional definition.
// - Immutability enforcement: Only tools in state 'draft' may be mutated.
//   Attempting to mutate a non-draft tool returns 400 immutable-asset-mutation.
// - Rejects states 'published' and 'deployed' (forbidden under STD-1).
// - Rejects cross-tenant mutation with 403 cross-tenant-mutation-forbidden.
func (s *DefaultToolService) UpdateDraftTool(ctx context.Context, auth *AuthContext, req *UpdateDraftToolRequest) (*ToolDefinition, error) {
	corrID := GenerateCorrelationID()
	ac := ensureAuth(auth)

	if req == nil {
		return nil, NewError(ErrorCodeInvalidMetadata, "UpdateDraftTool request cannot be nil", corrID, 400, false)
	}

	// 1. Validate tenant scoping
	if errEnv := ValidateTenantScope(ac, req.TenantRef, corrID); errEnv != nil {
		return nil, errEnv
	}

	// Cross-tenant mutation check
	if ac.TenantID != "" && req.TenantRef != ac.TenantID {
		return nil, NewError(ErrorCodeCrossTenantMutationForbidden,
			fmt.Sprintf("Cross-tenant draft update from %q to %q is strictly forbidden", ac.TenantID, req.TenantRef),
			corrID, 403, false)
	}

	// 2. Validate authoring permission
	if !ac.CanManageAssetKind(AssetKindTool) {
		return nil, NewPermissionDeniedError(PermStudioArtifactManage, corrID)
	}

	// 3. Lookup existing tool
	toolID := strings.TrimSpace(req.ToolID)
	if toolID == "" {
		return nil, NewError(ErrorCodeMissingAssetID, "Tool ID is required for update", corrID, 400, false)
	}

	existing, err := s.store.Get(ctx, req.TenantRef, toolID)
	if err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to lookup tool: %v", err), corrID, 500, true)
	}
	if existing == nil {
		return nil, NewAssetNotFoundError(toolID, req.TenantRef, corrID)
	}

	// 4. Immutability enforcement: only 'draft' tools may be edited
	if existing.FunctionalState != ToolStateDraft {
		return nil, NewImmutableAssetMutationError(toolID, string(existing.FunctionalState), corrID)
	}

	// 5. State transition check: 'published' and 'deployed' are forbidden under STD-1
	if req.FunctionalState != "" {
		stateStr := string(req.FunctionalState)
		if stateStr == "published" || stateStr == "deployed" {
			return nil, NewInvalidMetadataError(
				fmt.Sprintf("Functional state %q is forbidden in STD-1 (CAP-STD-003 §14). Use 'available-reference'.", stateStr),
				corrID,
			)
		}
		if !req.FunctionalState.IsValid() {
			return nil, NewInvalidMetadataError(fmt.Sprintf("Invalid functional state %q", req.FunctionalState), corrID)
		}
	}

	// 6. Apply updates
	now := time.Now().UTC()
	actor := ac.PrincipalID
	if actor == "" {
		actor = "automation-designer"
	}

	oldDesc := existing.Description

	if strings.TrimSpace(req.Name) != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}
	if strings.TrimSpace(req.DisplayName) != "" {
		existing.DisplayName = strings.TrimSpace(req.DisplayName)
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if strings.TrimSpace(req.Purpose) != "" {
		existing.Purpose = strings.TrimSpace(req.Purpose)
	}
	if req.Owner != nil {
		owner := *req.Owner
		owner.Normalize()
		existing.Owner = owner
	}
	if req.ExecutionNature != "" {
		if !req.ExecutionNature.IsValid() {
			return nil, NewInvalidMetadataError(fmt.Sprintf("Invalid execution nature %q", req.ExecutionNature), corrID)
		}
		existing.ExecutionNature = req.ExecutionNature
	}
	if req.FunctionalState != "" {
		existing.FunctionalState = req.FunctionalState
	}
	if req.SupportedOperations != nil {
		existing.SupportedOperations = req.SupportedOperations
	}
	if req.Operations != nil {
		existing.Operations = req.Operations
	}
	if req.Parameters != nil {
		existing.Parameters = req.Parameters
	}
	if req.OutputParameters != nil {
		existing.OutputParameters = req.OutputParameters
	}
	if req.Dependencies != nil {
		existing.Dependencies = req.Dependencies
	}
	if req.ProviderRefs != nil {
		existing.ProviderRefs = req.ProviderRefs
	}
	if req.Consumers != nil {
		existing.Consumers = req.Consumers
	}
	if req.ActionClass != nil {
		existing.ActionClass = *req.ActionClass
	}
	if req.RiskClassification != "" {
		existing.RiskClassification = req.RiskClassification
	}
	if req.SideEffects != nil {
		existing.SideEffects = *req.SideEffects
	}
	if req.Limitations != nil {
		existing.Limitations = req.Limitations
	}

	existing.UpdatedAt = now
	existing.Provenance.UpdatedBy = actor
	existing.Provenance.UpdatedAt = now

	reason := req.Reason
	if reason == "" {
		reason = "Draft tool update"
	}

	existing.AuditTrail = append(existing.AuditTrail, AuditRecord{
		Timestamp:   now,
		Actor:       actor,
		Reason:      reason,
		OldMetadata: oldDesc,
		NewMetadata: existing.Description,
		Version:     existing.Version.SemVer,
		TenantID:    req.TenantRef,
		Action:      "update_draft",
	})

	// 7. Persist
	if err := s.store.Put(ctx, existing); err != nil {
		return nil, NewError(ErrorCodeInternalError, fmt.Sprintf("Failed to update tool: %v", err), corrID, 500, true)
	}

	return existing.Clone(), nil
}

// -------------------------------------------------------------------------
// Operational Decoupled Runtime Health Updates
// -------------------------------------------------------------------------

// SetRuntimeAvailability updates runtime availability status without modifying immutable contract definition.
func (s *DefaultToolService) SetRuntimeAvailability(ctx context.Context, tenantRef TenantRef, toolID string, status RuntimeAvailabilityStatus) error {
	tool, err := s.store.Get(ctx, tenantRef, toolID)
	if err != nil {
		return err
	}
	if tool == nil {
		return NewAssetNotFoundError(toolID, tenantRef, GenerateCorrelationID())
	}
	tool.RuntimeAvailability = status
	return s.store.Put(ctx, tool)
}

// -------------------------------------------------------------------------
// REQ-PROD-006: Permanent Execution Blocking Implementation
// -------------------------------------------------------------------------

// Invoke is permanently blocked under REQ-PROD-006.
func (s *DefaultToolService) Invoke(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Invoke", GenerateCorrelationID())
}

// Execute is permanently blocked under REQ-PROD-006.
func (s *DefaultToolService) Execute(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Execute", GenerateCorrelationID())
}

// Run is permanently blocked under REQ-PROD-006.
func (s *DefaultToolService) Run(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Run", GenerateCorrelationID())
}

// Call is permanently blocked under REQ-PROD-006.
func (s *DefaultToolService) Call(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Call", GenerateCorrelationID())
}

// Trigger is permanently blocked under REQ-PROD-006.
func (s *DefaultToolService) Trigger(ctx context.Context, auth *AuthContext, toolID string, payload any) error {
	return NewExecutionBlockedError("Trigger", GenerateCorrelationID())
}
