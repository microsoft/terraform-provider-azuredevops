package azuredevops

import (
	"context"
)

type Auth struct {
	AuthString string
}

type AuthProvider interface {
	GetAuth(ctx context.Context) (string, error)
}
