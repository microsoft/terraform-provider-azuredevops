package framework

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

// ResourceTimeout specifies the default timeout value for each operation.
type ResourceTimeout struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

type ResourceIdentity interface {
	// Convert from an import id to the identity
	FromId(id string) error
	// Fields returns each identity field value, together with its path and its corresponding state path.
	Fields() []IdentityField
}

type IdentityField struct {
	PathIdentity path.Path
	PathState    path.Path
	Value        attr.Value
}

// Resource interface defines the mandatory methods that a resource requires to implement.
// Some of the method can be implemented by embedding a utility struct (see the comments).
type Resource interface {
	// ResourceType returns the resource type
	ResourceType() string

	// Identity returns a ResourceIdentity.
	Identity() ResourceIdentity

	// Resource implements the framework Resource interface.
	//
	// For Create, the implement doesn't need to handle the protocol response if the req.Plan contains all the information
	//             for a follow-up Read(). Otherwise, the implement shall set those Read() related attributes to the state.
	// For Update/Delete, the implement doesn't need to handle the protocol response, which is done by wrapper.
	// For Read, the implement must handle the protocol response (e.g. set the state).
	//
	// NOTE: Since the Metadata() is implemented by the wrapper, the implement struct shall not implement it.
	// 		 Instead, it is supposed to embed the WithMetadata to meet the interface requirement.
	//
	resource.ResourceWithIdentity

	// SetMeta sets the provider meta to the resource.
	// This is implemented by the ImplMeta, the implement struct shall simply embed it.
	SetMeta(meta.Meta)

	// Log logs a message
	// This is implemented by the ImplLog, the implement struct shall simply embed it.
	Info(ctx context.Context, msg string, additionalFields ...map[string]any)
	Warn(ctx context.Context, msg string, additionalFields ...map[string]any)
	Error(ctx context.Context, msg string, additionalFields ...map[string]any)
}

// ResourceWithTimeout is an opt-in interface that can implement customized timeout.
type ResourceWithTimeout interface {
	resource.Resource

	// Timeout returns the timeout for each operation.
	Timeout() ResourceTimeout
}

// ResourceWithCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithCreatePoll interface {
	resource.Resource

	CreatePollOption(ctx context.Context) retry.RetryOption

	// CreatePollCheck checks the post create state being read against the plan.
	// If the expected state is not met, return false, which will retry the poll.
	CreatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) bool

	// CreatePollRetryableDiag tells whether a diagnostics returned by the post create read is retryable.
	CreatePollRetryableDiag(diag.Diagnostic) bool
}

// ResourceWithUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithUpdatePoll interface {
	resource.Resource

	UpdatePollOption(ctx context.Context) retry.RetryOption

	// UpdatePollCheck checks the post update state being read against the plan.
	// If the expected state is not met, return false, which will retry the poll.
	UpdatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) bool

	// UpdatePollRetryableDiag tells whether a diagnostics returned by the post update read is retryable.
	UpdatePollRetryableDiag(diag.Diagnostic) bool
}

// ResourceWithDeletePoll is an opt-in interface that introduces a read after delete to ensure the resource has been consistently removed.
type ResourceWithDeletePoll interface {
	resource.Resource

	DeletePollOption(ctx context.Context) retry.RetryOption

	// DeletePollRetryableDiag tells whether a diagnostics returned by the post delete read is retryable.
	DeletePollRetryableDiag(diag.Diagnostic) bool

	// DeletePollTerminalDiag represents the terminal
	DeletePollTerminalDiag(diag.Diagnostic) bool
}

// Additionally, a resource can opt-in any of the following interfaces.
//
// resource.ResourceWithConfigValidators
// resource.ResourceWithModifyPlan
// resource.ResourceWithMoveState
// resource.ResourceWithUpgradeState
// resource.ResourceWithValidateConfig
// resource.ResourceWithUpgradeIdentity
