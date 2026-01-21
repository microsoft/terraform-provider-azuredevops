package workitemtracking

import (
	"context"
	"time"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/workitemtrackingprocess"
)

func RetryOnUnexpectedException(ctx context.Context, timeout time.Duration, f func() error) error {
	return workitemtrackingprocess.RetryOnUnexpectedException(ctx, timeout, f)
}
