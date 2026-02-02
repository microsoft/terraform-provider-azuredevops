package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/retry"
)

type WriteOperation string

const (
	WriteOperationCreate     WriteOperation = "create"
	WriteOperationPostCreate WriteOperation = "post-create"
	WriteOperationUpdate     WriteOperation = "update"
	WriteOperationPostUpdate WriteOperation = "post-update"
)

type PollChecker struct {
	AttrPath path.Path
	Target   attr.Value
}

// ResourceWithCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithCreatePoll interface {
	resource.Resource

	CreatePollOption(ctx context.Context) retry.RetryOption

	// CreatePollCheckers returns poll checkers that check the state being read after create.
	CreatePollCheckers() []PollChecker

	// CreatePollRetryableDiag tells whether the diagnostics returned by the read after create is retryable.
	CreatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithPostCreatePoll is an opt-in interface that makes the read after create retryable on certain conditions.
type ResourceWithPostCreatePoll interface {
	resource.Resource

	PostCreatePollOption(ctx context.Context) retry.RetryOption

	// PostCreatePollCheckers returns poll checkers that check the state being read after post create.
	PostCreatePollCheckers() []PollChecker

	// PostCreatePollRetryableDiag tells whether the diagnostics returned by the read after post create is retryable.
	PostCreatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithUpdatePoll interface {
	resource.Resource

	UpdatePollOption(ctx context.Context) retry.RetryOption

	// UpdatePollCheckers returns poll checkers that check the state being read after update.
	UpdatePollCheckers() []PollChecker

	// UpdatePollRetryableDiag tells whether the diagnostics returned by the read after update is retryable.
	UpdatePollRetryableDiags(diag.Diagnostics) bool
}

// ResourceWithPostUpdatePoll is an opt-in interface that makes the read after update retryable on certain conditions.
type ResourceWithPostUpdatePoll interface {
	resource.Resource

	PostUpdatePollOption(ctx context.Context) retry.RetryOption

	// PostUpdatePollCheckers returns poll checkers that check the state being read after update.
	PostUpdatePollCheckers() []PollChecker

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
