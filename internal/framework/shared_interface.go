package framework

import (
	"context"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

type Logger interface {
	// This is implemented by the ImplLog, the implement struct shall simply embed it.
	Info(ctx context.Context, msg string, additionalFields ...map[string]any)
	Warn(ctx context.Context, msg string, additionalFields ...map[string]any)
	Error(ctx context.Context, msg string, additionalFields ...map[string]any)
}

type ResourceTyper interface {
	// ResourceType returns the type of a resource/data source.
	ResourceType() string
}

// BaseResource is the base interface that shall be implemented by both Resource and DataSource.
type BaseResource interface {
	Logger
	ResourceTyper

	// SetMeta sets the provider meta to the resource.
	// This is implemented by the ImplMeta, the implement struct shall simply embed it.
	SetMeta(meta.Meta)
}
