package framework

import (
	"context"

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

type ImplMetadata struct{}

func (ImplMetadata) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {
	panic("This should have been implemented by the wrapper")
}

type ImplLog[T Resource] struct{ r T }

func (l ImplLog[T]) Log(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.SubsystemInfo(ctx, l.r.ResourceType(), msg, additionalFields...)
}
