package workitemtrackingprocess

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
)

// retryOnUnexpectedException retries the given function when Azure DevOps returns
// an "UnexpectedException" (TF401349) error. This error typically occurs due to
// eventual consistency issues where a resource cannot be modified yet because
// dependent resources haven't been fully processed.
func retryOnUnexpectedException(ctx context.Context, timeout time.Duration, f func() error) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err := f()
		if err == nil {
			return nil
		}

		var wrappedErr azuredevops.WrappedError
		if errors.As(err, &wrappedErr) && wrappedErr.TypeKey != nil && *wrappedErr.TypeKey == "UnexpectedException" {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
}