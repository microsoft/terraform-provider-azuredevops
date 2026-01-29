package retry

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/ctxutil"
)

type RetryOption struct {
	Delay                     time.Duration // Wait this time before starting checks
	Timeout                   time.Duration // The amount of time to wait before timeout
	MinTimeout                time.Duration // Smallest time to wait before refreshes
	PollInterval              time.Duration // Override MinTimeout/backoff and only poll this often
	ContinuousTargetOccurence int           // Number of times the Target state has to occur continuously
}

func NewSimpleRetryOption(ctx context.Context, streak int) RetryOption {
	return RetryOption{
		Timeout:                   ctxutil.UntilDeadline(ctx),
		MinTimeout:                time.Second,
		ContinuousTargetOccurence: streak,
	}
}

// RetryContext is a basic wrapper around StateChangeConf that will just retry
// a function until it no longer returns an error.
//
// Cancellation from the passed in context will propagate through to the
// underlying StateChangeConf
func RetryContext(ctx context.Context, opt RetryOption, f RetryFunc) error {
	// These are used to pull the error out of the function; need a mutex to
	// avoid a data race.
	var resultErr error
	var resultErrMu sync.Mutex

	c := &StateChangeConf{
		Pending:                   []string{"retryableerror"},
		Target:                    []string{"success"},
		Delay:                     opt.Delay,
		Timeout:                   opt.Timeout,
		MinTimeout:                opt.MinTimeout,
		PollInterval:              opt.PollInterval,
		ContinuousTargetOccurence: opt.ContinuousTargetOccurence,
		Refresh: func() (interface{}, string, error) {
			rerr := f()

			resultErrMu.Lock()
			defer resultErrMu.Unlock()

			if rerr == nil {
				resultErr = nil
				return 42, "success", nil
			}

			resultErr = rerr.Err

			if rerr.Retryable {
				return 42, "retryableerror", nil
			}
			return nil, "quit", rerr.Err
		},
	}

	_, waitErr := c.WaitForStateContext(ctx)

	// Need to acquire the lock here to be able to avoid race using resultErr as
	// the return value
	resultErrMu.Lock()
	defer resultErrMu.Unlock()

	// resultErr may be nil because the wait timed out and resultErr was never
	// set; this is still an error
	if resultErr == nil {
		return waitErr
	}
	// resultErr takes precedence over waitErr if both are set because it is
	// more likely to be useful
	return resultErr
}

// RetryFunc is the function retried until it succeeds.
type RetryFunc func() *RetryError

// RetryError is the required return type of RetryFunc. It forces client code
// to choose whether or not a given error is retryable.
type RetryError struct {
	Err       error
	Retryable bool
}

func (e *RetryError) Unwrap() error {
	return e.Err
}

// RetryableError is a helper to create a RetryError that's retryable from a
// given error. To prevent logic errors, will return an error when passed a
// nil error.
func RetryableError(err error) *RetryError {
	if err == nil {
		return &RetryError{
			Err: errors.New("empty retryable error received. " +
				"This is a bug with the Terraform provider and should be " +
				"reported as a GitHub issue in the provider repository."),
			Retryable: false,
		}
	}
	return &RetryError{Err: err, Retryable: true}
}

// NonRetryableError is a helper to create a RetryError that's _not_ retryable
// from a given error. To prevent logic errors, will return an error when
// passed a nil error.
func NonRetryableError(err error) *RetryError {
	if err == nil {
		return &RetryError{
			Err: errors.New("empty non-retryable error received. " +
				"This is a bug with the Terraform provider and should be " +
				"reported as a GitHub issue in the provider repository."),
			Retryable: false,
		}
	}
	return &RetryError{Err: err, Retryable: false}
}
