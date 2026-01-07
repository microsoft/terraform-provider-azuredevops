package sdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

type WithMeta struct {
	meta.Meta
}

func (r *WithMeta) SetMeta(m meta.Meta) {
	r.Meta = m
}

type WithMetadata struct{}

func (WithMetadata) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {
	panic("This should have been implemented by the wrapper")
}
