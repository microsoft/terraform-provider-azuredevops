package sdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
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
