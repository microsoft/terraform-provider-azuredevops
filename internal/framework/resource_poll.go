package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

type WriteOperation string

const (
	WriteOperationCreate     WriteOperation = "create"
	WriteOperationPostCreate WriteOperation = "post-create"
	WriteOperationUpdate     WriteOperation = "update"
	WriteOperationPostUpdate WriteOperation = "post-update"
)

// ResourceWithCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithCreatePoll interface {
	resource.Resource

	CreatePollOption(ctx context.Context) retry.RetryOption

	// CreatePollCheck checks the state being read after create against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	CreatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// CreatePollRetryableDiag tells whether the diagnostics returned by the read after create is retryable.
	CreatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithPostCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithPostCreatePoll interface {
	resource.Resource

	PostCreatePollOption(ctx context.Context) retry.RetryOption

	// PostCreatePollCheck checks the state being read after post create against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	PostCreatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// PostCreatePollRetryableDiag tells whether the diagnostics returned by the read after post create is retryable.
	PostCreatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithUpdatePoll interface {
	resource.Resource

	UpdatePollOption(ctx context.Context) retry.RetryOption

	// UpdatePollCheck checks the state being read after update against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	UpdatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// UpdatePollRetryableDiag tells whether the diagnostics returned by the read after update is retryable.
	UpdatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithPostUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithPostUpdatePoll interface {
	resource.Resource

	PostUpdatePollOption(ctx context.Context) retry.RetryOption

	// PostUpdatePollCheck checks the state being read after post update against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	PostUpdatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// PostUpdatePollRetryableDiag tells whether the diagnostics returned by the read after post update is retryable.
	PostUpdatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithDeletePoll is an opt-in interface that introduces a read after delete to ensure the resource has been consistently removed.
type ResourceWithDeletePoll interface {
	resource.Resource

	DeletePollOption(ctx context.Context) retry.RetryOption

	// DeletePollRetryableDiag tells whether the diagnostics returned by the post delete read is retryable.
	DeletePollRetryableDiags(diag.Diagnostics) bool

	// DeletePollTerminalDiags represents the terminal diags returned by the post delete read.
	DeletePollTerminalDiags(diag.Diagnostics) bool
}
