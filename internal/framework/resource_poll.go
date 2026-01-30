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
	WriteOperationCreate WriteOperation = "create"
	WriteOperationUpdate WriteOperation = "update"
)

// ResourceWithCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithCreatePoll interface {
	resource.Resource

	CreatePollOption(ctx context.Context) retry.RetryOption

	// CreatePollCheck checks the post create state being read against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	CreatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// CreatePollRetryableDiag tells whether the diagnostics returned by the post create read is retryable.
	CreatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithUpdatePoll interface {
	resource.Resource

	UpdatePollOption(ctx context.Context) retry.RetryOption

	// UpdatePollCheck checks the post update state being read against the plan.
	// If the expected state is not met, return error, which will retry the poll.
	UpdatePollCheck(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) error

	// UpdatePollRetryableDiag tells whether the diagnostics returned by the post update read is retryable.
	UpdatePollRetryableDiags(diag.Diagnostics) bool
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

