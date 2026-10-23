package studiocatalog

import (
	"strings"
)

// Canonical permission identifiers from cmdr-product-spec/14-security-permissions-and-trust/permission-catalog.md
const (
	PermStudioLibraryRead         = "perm.studio.library.read"
	PermStudioSkillRead           = "perm.studio.skill.read"
	PermStudioSkillManage         = "perm.studio.skill.manage"
	PermStudioWorkflowRead        = "perm.studio.workflow.read"
	PermStudioWorkflowManage      = "perm.studio.workflow.manage"
	PermStudioAutomationAgentRead = "perm.studio.automation-agent.read"
	PermStudioAutomationAgentManage = "perm.studio.automation-agent.manage"
	PermStudioAgentTeamRead       = "perm.studio.agent-team.read"
	PermStudioAgentTeamManage     = "perm.studio.agent-team.manage"
	PermCmdrStudioVersionRead     = "perm.cmdr-studio.version.read"
	PermCmdrStudioVersionManage   = "perm.cmdr-studio.version.manage"
	PermStudioBuilderRead         = "perm.studio.builder.read"
	PermStudioControlRoomRead     = "perm.studio.control-room.read"
)

// AuthContext holds caller identity, tenant authorization context, and permissions.
type AuthContext struct {
	PrincipalID          string          `json:"principal_id"`
	TenantID             string          `json:"tenant_id"`
	AuthorizedTenantSet  []string        `json:"authorized_tenant_set"` // ATS for MSSP aggregation
	Permissions          map[string]bool `json:"permissions"`
}

// NewAuthContext constructs an AuthContext with normalized permissions.
func NewAuthContext(principalID, tenantID string, authorizedTenants []string, permissions ...string) AuthContext {
	perms := make(map[string]bool)
	for _, p := range permissions {
		perms[strings.TrimSpace(p)] = true
	}
	return AuthContext{
		PrincipalID:         principalID,
		TenantID:            tenantID,
		AuthorizedTenantSet: authorizedTenants,
		Permissions:         perms,
	}
}

// HasPermission checks if the context contains the required permission, checking canonical aliases.
func (ac *AuthContext) HasPermission(permission string) bool {
	if ac.Permissions == nil {
		return false
	}
	if ac.Permissions[permission] {
		return true
	}

	// Normalize between perm.studio.* and perm.cmdr-studio.* namespaces
	alias := resolvePermissionAlias(permission)
	if alias != "" && ac.Permissions[alias] {
		return true
	}

	return false
}

// resolvePermissionAlias resolves dual-namespace identifiers.
func resolvePermissionAlias(perm string) string {
	switch perm {
	case PermStudioSkillRead:
		return "perm.cmdr-studio.skill.read"
	case "perm.cmdr-studio.skill.read":
		return PermStudioSkillRead
	case PermStudioSkillManage:
		return "perm.cmdr-studio.skill.manage"
	case "perm.cmdr-studio.skill.manage":
		return PermStudioSkillManage

	case PermStudioWorkflowRead:
		return "perm.cmdr-studio.workflow.read"
	case "perm.cmdr-studio.workflow.read":
		return PermStudioWorkflowRead
	case PermStudioWorkflowManage:
		return "perm.cmdr-studio.workflow.manage"
	case "perm.cmdr-studio.workflow.manage":
		return PermStudioWorkflowManage

	case PermStudioAutomationAgentRead:
		return "perm.cmdr-studio.automation-agent.read"
	case "perm.cmdr-studio.automation-agent.read":
		return PermStudioAutomationAgentRead
	case PermStudioAutomationAgentManage:
		return "perm.cmdr-studio.automation-agent.manage"
	case "perm.cmdr-studio.automation-agent.manage":
		return PermStudioAutomationAgentManage

	case PermStudioAgentTeamRead:
		return "perm.cmdr-studio.agent-team.read"
	case "perm.cmdr-studio.agent-team.read":
		return PermStudioAgentTeamRead
	case PermStudioAgentTeamManage:
		return "perm.cmdr-studio.agent-team.manage"
	case "perm.cmdr-studio.agent-team.manage":
		return PermStudioAgentTeamManage

	default:
		return ""
	}
}

// CanReadAssetKind checks whether the caller has the specific read permission for an asset kind.
// Under canonical permission catalog, Tool has no dedicated atomic permission and is governed by PermStudioLibraryRead.
func (ac *AuthContext) CanReadAssetKind(kind AssetKind) bool {
	// Baseline check: Studio Library read is required for all
	if !ac.HasPermission(PermStudioLibraryRead) {
		return false
	}

	switch kind {
	case AssetKindSkill:
		return ac.HasPermission(PermStudioSkillRead)
	case AssetKindWorkflow:
		return ac.HasPermission(PermStudioWorkflowRead)
	case AssetKindAutomationAgent:
		return ac.HasPermission(PermStudioAutomationAgentRead)
	case AssetKindAgentTeam:
		return ac.HasPermission(PermStudioAgentTeamRead)
	case AssetKindTool:
		// Governed strictly by perm.studio.library.read
		return true
	default:
		return false
	}
}

// ValidateTenantScope enforces tenant scoping rules under tenant-scope-contract.md.
func ValidateTenantScope(tenantRef string, authorizedTenants []string, correlationID string) *ErrorEnvelope {
	trimmed := strings.TrimSpace(tenantRef)
	if trimmed == "" {
		return NewError(
			ErrorCodeMissingTenant,
			"Tenant identifier is required for all Studio catalog operations",
			correlationID,
			400,
			false,
		)
	}

	if trimmed == "*" || trimmed == "%" || strings.EqualFold(trimmed, "all") {
		return NewError(
			ErrorCodeWildcardTenant,
			"Wildcard tenant selectors ('*', '%', 'all') are strictly forbidden",
			correlationID,
			400,
			false,
		)
	}

	// If ATS is specified, requested tenant must be in the set
	if len(authorizedTenants) > 0 {
		found := false
		for _, t := range authorizedTenants {
			if t == trimmed {
				found = true
				break
			}
		}
		if !found {
			return NewError(
				ErrorCodeUnauthorizedTenantInSet,
				"Requested tenant is not in caller's Authorized Tenant Set (ATS)",
				correlationID,
				403,
				false,
			)
		}
	}

	return nil
}

// ValidateMSSPTenantSet enforces that all requested tenants are within the Authorized Tenant Set.
// ADR-0008 approved decision D4: requested ⊆ authorized.
func ValidateMSSPTenantSet(requestedTenants []string, authorizedTenants []string, correlationID string) *ErrorEnvelope {
	if len(requestedTenants) == 0 {
		return NewError(
			ErrorCodeMissingTenant,
			"At least one tenant identifier must be requested",
			correlationID,
			400,
			false,
		)
	}

	authMap := make(map[string]bool)
	for _, at := range authorizedTenants {
		authMap[at] = true
	}

	for _, rt := range requestedTenants {
		rtTrimmed := strings.TrimSpace(rt)
		if rtTrimmed == "" {
			return NewError(ErrorCodeMissingTenant, "Empty tenant ref in requested tenant set", correlationID, 400, false)
		}
		if rtTrimmed == "*" || rtTrimmed == "%" || strings.EqualFold(rtTrimmed, "all") {
			return NewError(ErrorCodeWildcardTenant, "Wildcards are forbidden in MSSP tenant set", correlationID, 400, false)
		}
		if !authMap[rtTrimmed] {
			return NewError(
				ErrorCodeUnauthorizedTenantInSet,
				"Requested tenant "+rtTrimmed+" is not in caller's Authorized Tenant Set",
				correlationID,
				403,
				false,
			)
		}
	}

	return nil
}

// SanitizeAssetForInspect creates a sanitized 200 OK stub when an asset exists within the caller's tenant
// but the caller lacks the specific asset read permission (Existence Oracle Elimination).
func SanitizeAssetForInspect(original *Asset) *Asset {
	if original == nil {
		return nil
	}
	stub := &Asset{
		ID:                 original.ID,
		TenantID:           original.TenantID,
		EnvironmentID:      original.EnvironmentID,
		Type:               original.Type,
		Name:               original.Name,
		DisplayName:        original.DisplayName,
		Description:        "", // Redacted
		Tags:               []string{},
		Owner: Owner{
			OwnerID:     "redacted",
			OwnerName:   "redacted",
			OwnerType:   "unknown",
			IsUnknown:   true,
			DisplayName: "owner unknown",
		},
		Contributors:       []string{},
		Version: VersionInfo{
			VersionID:   original.Version.VersionID,
			SemVer:      original.Version.SemVer,
			CreatedAt:   original.Version.CreatedAt,
			IsImmutable: original.Version.IsImmutable,
		},
		LifecycleState:     original.LifecycleState,
		FunctionalState:    FunctionalStatePermissionFiltered,
		ExecutionClass:     original.ExecutionClass,
		RiskClassification: original.RiskClassification,
		IsAIPowered:        original.IsAIPowered,
		HasNonAIFallback:   original.HasNonAIFallback,
		Parameters:         []Parameter{},
		Dependencies:       []Dependency{},
		Consumers:          []Consumer{},
		Supersession: Supersession{
			IsSuperseded: original.Supersession.IsSuperseded,
		},
		Provenance: Provenance{
			CreatedAt: original.Provenance.CreatedAt,
			UpdatedAt: original.Provenance.UpdatedAt,
		},
		CreatedAt: original.CreatedAt,
		UpdatedAt: original.UpdatedAt,
	}
	return stub
}
