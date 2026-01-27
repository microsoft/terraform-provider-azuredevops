package ctxutil

import (
	"context"
	"time"
)

func UntilDeadline(ctx context.Context) time.Duration {
	t, ok := ctx.Deadline()
	if !ok {
		return 0
	}
	return time.Until(t)
}
