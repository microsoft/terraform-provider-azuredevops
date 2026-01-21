package workitemtrackingprocess

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// retryOnCondition executes f and retries if isRetryable returns true for the error.
func retryOnCondition(ctx context.Context, timeout time.Duration, f func() error, isRetryable func(error) bool) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err := f()
		if err == nil {
			return nil
		}
		if isRetryable(err) {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
}

// RetryOnUnexpectedException retries the given function when Azure DevOps returns
// an "UnexpectedException" (TF401349) error. This error typically occurs due to
// eventual consistency issues where a resource cannot be modified yet because
// dependent resources haven't been fully processed.
func RetryOnUnexpectedException(ctx context.Context, timeout time.Duration, f func() error) error {
	return retryOnCondition(ctx, timeout, f, func(err error) bool {
		return utils.ResponseWasTypeKey(err, "UnexpectedException")
	})
}

// retryOnNotFound retries the given function when Azure DevOps returns
// a 404 Not Found status code. This handles eventual consistency where
// a newly created resource may not immediately be visible.
func retryOnNotFound(ctx context.Context, timeout time.Duration, f func() error) error {
	return retryOnCondition(ctx, timeout, f, func(err error) bool {
		return utils.ResponseWasStatusCode(err, http.StatusNotFound)
	})
}

// retryOnContributionNotFound retries the given function when Azure DevOps returns
// a VS403120 error ("Contribution does not exist"). This handles eventual consistency
// where a newly installed extension's contributions may not be immediately available.
func retryOnContributionNotFound(ctx context.Context, timeout time.Duration, f func() error) error {
	return retryOnCondition(ctx, timeout, f, func(err error) bool {
		return utils.ResponseContainsStatusMessage(err, "VS403120")
	})
}
