package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

type ImplSetMeta struct {
	meta.Meta
}

func (r *ImplSetMeta) SetMeta(m meta.Meta) {
	r.Meta = m
}

type ImplLog[T ResourceTyper] struct{ r T }

func (l ImplLog[T]) Info(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.SubsystemInfo(ctx, l.r.ResourceType(), msg, additionalFields...)
}

func (l ImplLog[T]) Warn(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.SubsystemWarn(ctx, l.r.ResourceType(), msg, additionalFields...)
}

func (l ImplLog[T]) Error(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.SubsystemError(ctx, l.r.ResourceType(), msg, additionalFields...)
}

type ImplIdentity[T any] struct{ model T }

func (i ImplIdentity[T]) IdentityModel() any {
	return i.model
}

// Resource Specific

type ImplResourceMetadata struct{}

func (ImplResourceMetadata) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
}

// Data Source Specific

type ImplDataSourceMetadata struct{}

func (ImplDataSourceMetadata) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
}
